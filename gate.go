package libninja

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

// openGates is every exported field of every struct type a package
// declares, with who outside the package touches it.
func (read reading) openGates() []finding {
	var found []finding
	for _, pkg := range read.plain() {
		for _, typeName := range typeNames(pkg.Types) {
			structure, ok :=
				typeName.Type().Underlying().(*types.Struct)
			if !ok || read.madeByATool(typeName.Pos()) {
				continue
			}
			for index := range structure.NumFields() {
				field := structure.Field(index)
				if field.Exported() &&
					!onTheWire(structure.Tag(index)) {
					tag := structure.Tag(index)
					found = append(found, read.gateOf(pkg,
						typeName, field, tag))
				}
			}
		}
	}
	return found
}

// gateOf is one open field as the master sees it.
func (read reading) gateOf(pkg *packages.Package, typeName *types.TypeName,
	field *types.Var, tag string) finding {
	seen := fmt.Sprintf("%s.%s stands open", typeName.Name(), field.Name())
	if touched := read.outsiders(field); positive(len(touched)) {
		seen += "; touched from " + strings.Join(touched, ", ")
	}
	seen += "."
	if tag != "" {
		seen += " Its tag says a reader needs it. " +
			"If this is the wire, vow it."
	}
	return finding{gate, read.address(pkg.PkgPath),
		read.at(field), seen, ""}
}

// outsiders is every package outside a name's own that uses it. Each is
// given once, by its address.
func (read reading) outsiders(obj types.Object) []string {
	var found []string
	for _, each := range read.uses[read.key(obj)] {
		where := read.address(each.from)
		if each.from != obj.Pkg().Path() &&
			!slices.Contains(found, where) {
			found = append(found, where)
		}
	}
	slices.Sort(found)
	return found
}

// typeNames is every type a package declares by name, aliases left
// out.
func typeNames(pkg *types.Package) []*types.TypeName {
	var found []*types.TypeName
	for _, name := range pkg.Scope().Names() {
		typeName, ok := pkg.Scope().Lookup(name).(*types.TypeName)
		if ok && !typeName.IsAlias() {
			found = append(found, typeName)
		}
	}
	return found
}

// strangers is every struct literal of a type from another package of
// the module.
func (read reading) strangers() []finding {
	var found []finding
	for _, each := range read.files() {
		for node := range ast.Preorder(each.file) {
			switch literal := node.(type) {
			case *ast.CompositeLit:
				built := read.strangerAt(each, literal)
				found = append(found, built...)
			}
		}
	}
	return found
}

// strangerAt is the finding a literal makes when it builds a struct of
// another package of the module, or none.
func (read reading) strangerAt(each source,
	literal *ast.CompositeLit) []finding {
	owner, ok := read.foreignStruct(each.pkg, literal)
	if !ok {
		return nil
	}
	seen := fmt.Sprintf("Builds %s.%s by literal.", owner.Pkg().Name(),
		owner.Name())
	return []finding{{stranger, read.address(each.pkg.PkgPath),
		read.at(literal), seen, ""}}
}

// foreignStruct is the type a literal builds, when it is a struct
// declared in another package of the module. Aliases are seen through:
// a borrowed name is still another's house.
func (read reading) foreignStruct(pkg *packages.Package,
	literal *ast.CompositeLit) (*types.TypeName, bool) {
	built := types.Unalias(pkg.TypesInfo.TypeOf(literal))
	if pointer, ok := built.(*types.Pointer); ok {
		built = types.Unalias(pointer.Elem())
	}
	declared, ok := built.(*types.Named)
	if !ok {
		return nil, false
	}
	owner := declared.Obj()
	if _, isStruct := declared.Underlying().(*types.Struct); !isStruct ||
		owner.Pkg() == nil || !read.inModule(owner.Pkg().Path()) ||
		athome(owner.Pkg().Path(), pkg.PkgPath) {
		return nil, false
	}
	return owner, true
}

// athome is whether a type's package is the one building it. That holds
// for the same package. It holds for the package an external test stands
// beside. It holds for that test building its own.
func athome(owner, building string) bool {
	return owner == building ||
		owner == strings.TrimSuffix(building, "_test")
}

// wells is every exported variable a package declares.
func (read reading) wells() []finding {
	var found []finding
	for _, pkg := range read.plain() {
		for _, name := range pkg.Types.Scope().Names() {
			variable, ok :=
				pkg.Types.Scope().Lookup(name).(*types.Var)
			if ok && variable.Exported() &&
				!read.madeByATool(variable.Pos()) {
				found = append(found, finding{well,
					read.address(pkg.PkgPath),
					read.at(variable),
					fmt.Sprintf(
						"%s can be changed by anyone.",
						name), ""})
			}
		}
	}
	return found
}

// onTheWire is whether a struct tag gives a field to encoding/json: it
// has a json key, and the key does not say the field is left out. The
// field is exported because json can read nothing private.
func onTheWire(tag string) bool {
	said, ok := reflect.StructTag(tag).Lookup("json")
	return ok && said != jsonLeftOut
}

// jsonLeftOut is what a json tag says when the field is not written, as
// encoding/json has defined it.
const jsonLeftOut = "-"
