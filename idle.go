package libninja

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// reach is how far a name is used. The least is only inside its own
// package. Next is from outside by tests alone. The most is from outside
// by real code.
type reach int

const (
	inside reach = iota
	testsOnly
	outside
)

// idle is every exported name no real code outside its package uses.
// The ones only tests use from outside are practice swords; the rest
// are wasted power.
//
// Fields are the gate's and variables the well's, so
// neither is counted twice; methods that satisfy an interface keep a
// promise and stay; and what a tool generated is the tool's.
//
// Support for tests is the work of a package named for another package
// and test, as net/http/httptest does. Tests alone may use its names.
// The package is the fault when real code imports it. It is the fault
// when no test does.
func (read reading) idle() []finding {
	var found []finding
	importers := read.importers()
	for _, pkg := range read.plain() {
		support := read.madeForTests(pkg)
		dead := false
		if support {
			var hits []finding
			hits, dead = read.supportHits(pkg,
				importers[pkg.PkgPath])
			found = append(found, hits...)
		}
		if dead {
			// said once for the package, not again for each name
			continue
		}
		for _, obj := range exported(pkg.Types) {
			if _, ok := obj.(*types.Func); ok &&
				read.promises[read.key(obj)] {
				continue
			}
			if _, ok := obj.(*types.Var); ok ||
				read.madeByATool(obj.Pos()) {
				continue
			}
			student := read.address(pkg.PkgPath)
			at := read.at(obj)
			name := shown(obj)
			switch read.reachOf(obj) {
			case inside:
				found = append(found, finding{waste,
					student, at,
					fmt.Sprintf("%s is exported, "+
						"and nothing outside "+
						"%s uses it.",
						name, student), ""})
			case testsOnly:
				if support {
					continue
				}
				found = append(found, finding{practice,
					student, at,
					fmt.Sprintf("%s is used from outside "+
						"only by tests.", name), ""})
			}
		}
	}
	return found
}

// madeForTests is whether a package says by its name that it supports
// tests. Go's own are another package's name and test. These are
// httptest for http, iotest for io and fstest for fs. The other package
// must be one the module holds or this one imports.
//
// contest and latest end in test and are named for no package.
func (read reading) madeForTests(pkg *packages.Package) bool {
	stem, ok := strings.CutSuffix(pkg.Name, "test")
	if !ok || stem == "" {
		return false
	}
	for _, each := range read.plain() {
		if each.Name == stem {
			return true
		}
	}
	for _, each := range pkg.Imports {
		if each.Name == stem {
			return true
		}
	}
	return false
}

// importer is a place a package is imported. It holds where, the package
// the master addresses, and whether the file is a test's.
type importer struct {
	at      token.Position
	student string
	test    bool
}

// importers is every import of the module's files, keyed by the path
// imported.
func (read reading) importers() map[string][]importer {
	found := map[string][]importer{}
	for _, each := range read.files() {
		for _, spec := range each.file.Imports {
			// the parser has read this path, so it unquotes
			path, _ := strconv.Unquote(spec.Path.Value)
			found[path] = append(found[path], importer{
				read.at(spec), read.address(each.pkg.PkgPath),
				each.test()})
		}
	}
	return found
}

// supportHits is where a package made for tests is at fault: each
// import by a file that is no test's, and the package whole when no test
// imports it. It also says whether no test does.
func (read reading) supportHits(pkg *packages.Package,
	imports []importer) ([]finding, bool) {
	var found []finding
	tested := false
	for _, each := range imports {
		if each.test {
			tested = true
			continue
		}
		found = append(found, finding{practice, each.student, each.at,
			fmt.Sprintf("%s is made for tests, and real code "+
				"imports it.", pkg.Name), ""})
	}
	if !tested {
		first := read.firstOf(pkg)
		found = append(found, finding{practice,
			read.address(pkg.PkgPath), read.at(first.Name),
			fmt.Sprintf("%s is made for tests, and no test "+
				"imports it.", pkg.Name), ""})
	}
	return found, !tested
}

// exported is every exported name a package declares, with the
// exported methods of its types.
func exported(pkg *types.Package) []types.Object {
	var found []types.Object
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)
		if obj.Exported() {
			found = append(found, obj)
		}
	}
	for _, typeName := range typeNames(pkg) {
		// typeNames leaves aliases out. Packages of a module declare
		// no basic type, so each one is a named type
		declared := typeName.Type().(*types.Named)
		for method := range declared.Methods() {
			if method.Exported() {
				found = append(found, method)
			}
		}
	}
	return found
}

// reachOf is how far a name is used. Types are reached by name or by
// any of their methods and fields. They are also reached by an exported
// function of their package that takes one or gives one back.
func (read reading) reachOf(obj types.Object) reach {
	home := obj.Pkg().Path()
	keys := []string{read.key(obj)}
	if typeName, ok := obj.(*types.TypeName); ok {
		keys = append(keys, read.membersOf(typeName)...)
		keys = append(keys, read.carriersOf(typeName)...)
	}
	far := inside
	for _, key := range keys {
		for _, each := range read.uses[key] {
			switch {
			case !each.test && each.from != home:
				return outside
			case each.test && (each.from != home || each.external):
				far = testsOnly
			}
		}
	}
	return far
}

// carriersOf is where each exported function of a type's package is
// declared that takes the type or gives it back. Code outside that calls
// one uses the type, though it never writes the type's name.
func (read reading) carriersOf(typeName *types.TypeName) []string {
	var found []string
	scope := typeName.Pkg().Scope()
	for _, name := range scope.Names() {
		function, ok := scope.Lookup(name).(*types.Func)
		if ok && function.Exported() &&
			carries(function.Signature(), typeName.Type()) {
			found = append(found, read.key(function))
		}
	}
	return found
}

// carries is whether a signature takes or gives a type, itself or by a
// pointer to it.
func carries(signature *types.Signature, carried types.Type) bool {
	for _, tuple := range []*types.Tuple{signature.Params(),
		signature.Results()} {
		for each := range tuple.Variables() {
			held := each.Type()
			if pointer, ok := held.(*types.Pointer); ok {
				held = pointer.Elem()
			}
			if types.Identical(held, carried) {
				return true
			}
		}
	}
	return false
}

// membersOf is where each method and field of a type is declared.
func (read reading) membersOf(typeName *types.TypeName) []string {
	var found []string
	declared, ok := typeName.Type().(*types.Named)
	if !ok {
		return nil
	}
	for method := range declared.Methods() {
		found = append(found, read.key(method))
	}
	switch under := declared.Underlying().(type) {
	case *types.Struct:
		for field := range under.Fields() {
			found = append(found, read.key(field))
		}
	case *types.Interface:
		for method := range under.Methods() {
			found = append(found, read.key(method))
		}
	}
	return found
}

// shown is a name as the master says it: a method with its type.
func shown(obj types.Object) string {
	method, ok := obj.(*types.Func)
	if !ok || method.Signature().Recv() == nil {
		return obj.Name()
	}
	const dot = "."
	held := method.Signature().Recv().Type().String()
	held = held[strings.LastIndex(held, dot)+len(dot):]
	return held + dot + obj.Name()
}

// firstOf is the file of a package that sorts first by name. The master
// addresses the whole package there. Every package that is built has
// a file.
func (read reading) firstOf(pkg *packages.Package) *ast.File {
	return slices.MinFunc(pkg.Syntax, read.byName)
}

// byName orders two files by their names.
func (read reading) byName(one, other *ast.File) int {
	return cmp.Compare(read.at(one).Filename, read.at(other).Filename)
}
