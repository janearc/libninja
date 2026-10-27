package libninja

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Module is a Go module read whole with its tests. It is ready for the
// scrolls.
type Module struct {
	read reading
}

// reading is a module read once for every scroll. It holds the packages
// and every use of every name, keyed by where the name is declared. It
// holds the methods that keep a promise to an interface, keyed the same
// way. It holds the vows, and which files are generated.
type reading struct {
	dir       string
	module    string
	fset      *token.FileSet
	pkgs      []*packages.Package
	uses      map[string][]use
	promises  map[string]bool
	vows      map[string]vow
	generated map[string]bool
}

// use is one place a name is used. It holds where, and from which
// package with any _test taken off. It holds whether the use is in a
// test file, and whether it is from an external test package.
type use struct {
	at       token.Position
	from     string
	test     bool
	external bool
}

// loadMode is what the scrolls need of each package. That is its files
// and syntax, and its types with the uses of every name.
const loadMode = packages.NeedName | packages.NeedFiles |
	packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo |
	packages.NeedImports | packages.NeedModule

// Load reads every package of the module dir is in, with its tests.
// Any directory inside a module means the whole module, since what a
// name reaches is only known by reading all of it.
//
// It refuses a directory in no module. It also refuses a module that
// does not compile, since the scrolls read types.
func Load(dir string) (Module, error) {
	whole, err := moduleRoot(dir)
	if err != nil {
		return Module{}, err
	}
	fset := token.NewFileSet()
	config := &packages.Config{Mode: loadMode, Dir: whole, Tests: true,
		Fset: fset}
	loaded, err := packages.Load(config, "./...")
	if err != nil {
		return Module{}, err
	}
	kept, err := compiled(loaded)
	if err != nil {
		return Module{}, err
	}
	module := moduleOf(kept)
	if module == "" {
		return Module{}, fmt.Errorf("libninja: no module under %s",
			whole)
	}
	read := reading{dir: whole, module: module, fset: fset, pkgs: kept}
	read.markGenerated()
	read.indexUses()
	read.gatherPromises()
	read.gatherVows()
	return Module{read}, nil
}

// compiled is the packages that hold code. The test mains go generates
// are left out. The first package that does not compile ends it, with an
// error that names the package.
func compiled(loaded []*packages.Package) ([]*packages.Package, error) {
	var kept []*packages.Package
	for _, pkg := range loaded {
		if strings.HasSuffix(pkg.ID, ".test") {
			continue
		}
		var broken []error
		for _, said := range pkg.Errors {
			broken = append(broken, said)
		}
		if err := errors.Join(broken...); err != nil {
			return nil, fmt.Errorf(
				"libninja: %s does not compile: %w",
				pkg.ID, err)
		}
		kept = append(kept, pkg)
	}
	return kept, nil
}

// moduleOf is the path of the module the packages were read from, or
// none. They are all of one module, so the first that knows it says.
func moduleOf(kept []*packages.Package) string {
	for _, pkg := range kept {
		if pkg.Module != nil {
			return pkg.Module.Path
		}
	}
	return ""
}

// plain is each package as it is built, without its tests.
func (read reading) plain() []*packages.Package {
	var found []*packages.Package
	for _, pkg := range read.pkgs {
		if pkg.ID == pkg.PkgPath &&
			!strings.HasSuffix(pkg.PkgPath, "_test") {
			found = append(found, pkg)
		}
	}
	return found
}

// inModule is whether a package path is this module's.
func (read reading) inModule(pkgPath string) bool {
	pkgPath = strings.TrimSuffix(pkgPath, "_test")
	return pkgPath == read.module ||
		strings.HasPrefix(pkgPath, read.module+"/")
}

// address is what the master calls a package: its path inside the
// module, or the module's last word for the root package.
func (read reading) address(pkgPath string) string {
	pkgPath = strings.TrimSuffix(pkgPath, "_test")
	inside := strings.TrimPrefix(strings.TrimPrefix(pkgPath, read.module),
		"/")
	if inside == "" {
		return path.Base(read.module)
	}
	return inside
}

// key is where something is declared. That place is the same in every
// build of its package, with tests or without.
func (read reading) key(thing placed) string {
	return read.at(thing).String()
}

// placed is anything that knows where it begins: a node of the syntax,
// or a name the types know.
type placed interface {
	Pos() token.Pos
}

// at is where something in the module begins.
func (read reading) at(thing placed) token.Position {
	return read.fset.Position(thing.Pos())
}

// source is one file of the module, with the package it was read in.
type source struct {
	pkg  *packages.Package
	file *ast.File
	name string
}

// test is whether the file is a test's.
func (each source) test() bool {
	return strings.HasSuffix(each.name, "_test.go")
}

// files is every file of the module once, though a package built with
// its tests holds its plain files a second time. Generated files are
// left out.
func (read reading) files() []source {
	walked := map[string]bool{}
	var found []source
	for _, pkg := range read.pkgs {
		for _, file := range pkg.Syntax {
			// the file's own name, not where a //line says it is
			name := read.fset.File(file.Pos()).Name()
			if walked[name] || ast.IsGenerated(file) {
				continue
			}
			walked[name] = true
			found = append(found, source{pkg, file, name})
		}
	}
	return found
}

// indexUses records every use of every name the module declares, keyed
// by the name's declaration. Files in both a package and its test build
// are read twice, so a use may be listed twice. The scrolls ask only
// whether a use of some kind is there.
func (read *reading) indexUses() {
	read.uses = map[string][]use{}
	for _, pkg := range read.pkgs {
		for ident, obj := range pkg.TypesInfo.Uses {
			if obj.Pkg() == nil ||
				!read.inModule(obj.Pkg().Path()) {
				continue
			}
			at := read.at(ident)
			declared := read.key(obj)
			read.uses[declared] = append(read.uses[declared], use{
				at: at,
				from: strings.TrimSuffix(pkg.PkgPath,
					"_test"),
				test: strings.HasSuffix(at.Filename,
					"_test.go"),
				external: strings.HasSuffix(pkg.PkgPath,
					"_test")})
		}
	}
}

// markGenerated records every generated file, which no scroll reads:
// what a tool wrote is the tool's to answer for.
func (read *reading) markGenerated() {
	read.generated = map[string]bool{}
	for _, pkg := range read.pkgs {
		for _, file := range pkg.Syntax {
			if ast.IsGenerated(file) {
				read.generated[read.fset.Position(
					file.Package).Filename] = true
			}
		}
	}
}

// madeByATool is whether something is declared in a generated file.
func (read reading) madeByATool(pos token.Pos) bool {
	return read.generated[read.fset.Position(pos).Filename]
}

// moduleRoot is the nearest directory at or above dir that holds a
// go.mod. The walk ends when it has asked the root. Abs fails only when
// the working directory cannot be read; the walk then starts from
// nothing and finds no module.
func moduleRoot(dir string) (string, error) {
	start, _ := filepath.Abs(dir)
	for at, last := start, ""; at != last; at, last = filepath.Dir(at), at {
		_, missing := os.Stat(filepath.Join(at, "go.mod"))
		if missing == nil {
			return at, nil
		}
	}
	return "", fmt.Errorf("libninja: %s is in no module", dir)
}
