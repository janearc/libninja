package libninja

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// pzmPath is puzzlemouth's tool package, whose tools answer --age.
const pzmPath = "github.com/janearc/puzzlemouth/pzm"

// ages is every main that game cannot stamp or that does not answer
// --age: one that declares no build and built or one whose source
// never mentions the flag and makes no pzm tool.
func (read reading) ages() []finding {
	var found []finding
	for _, pkg := range read.plain() {
		if pkg.Name != "main" {
			continue
		}
		at, answers := read.mainOf(pkg)
		var lacks []string
		if !declaresStamp(pkg) {
			lacks = append(lacks,
				"declares no build and built, so game "+
					"cannot stamp it")
		}
		if !answers {
			lacks = append(lacks, "does not answer --age")
		}
		if none(len(lacks)) {
			continue
		}
		found = append(found, finding{trueAge,
			read.address(pkg.PkgPath), at,
			"This main " + strings.Join(lacks, ", and ") + ".", ""})
	}
	return found
}

// declaresStamp is whether a package declares build and built as
// strings at the top level, the two names game build stamps.
func declaresStamp(pkg *packages.Package) bool {
	for _, name := range []string{"build", "built"} {
		obj, ok := pkg.Types.Scope().Lookup(name).(*types.Var)
		if !ok ||
			!types.Identical(obj.Type(), types.Typ[types.String]) {
			return false
		}
	}
	return true
}

// mainOf is where a package's func main is and whether it answers
// --age. It answers when any of its source names the flag as a string:
// "age", "-age" or "--age". It also answers when it makes a pzm tool,
// which answers it for every tool.
func (read reading) mainOf(pkg *packages.Package) (token.Position, bool) {
	var at token.Position
	answers := false
	for _, file := range pkg.Syntax {
		if at.Filename == "" {
			at = read.at(file)
		}
		for node := range ast.Preorder(file) {
			switch node := node.(type) {
			case *ast.FuncDecl:
				if node.Name.Name == "main" &&
					node.Recv == nil {
					at = read.at(node)
				}
			case *ast.BasicLit:
				if namesTheFlag(node) {
					answers = true
				}
			case *ast.SelectorExpr:
				if makesATool(pkg, node) {
					answers = true
				}
			}
		}
	}
	return at, answers
}

// namesTheFlag is whether a literal is the age flag as a string. Only a
// string literal is written in double quotes, so its value says enough.
func namesTheFlag(lit *ast.BasicLit) bool {
	switch lit.Value {
	case `"age"`, `"-age"`, `"--age"`:
		return true
	}
	return false
}

// makesATool is whether a selector names pzm's New, the maker of a tool.
func makesATool(pkg *packages.Package, node *ast.SelectorExpr) bool {
	called, ok := pkg.TypesInfo.Uses[node.Sel].(*types.Func)
	return ok && called.Pkg() != nil &&
		called.Pkg().Path() == pzmPath && called.Name() == "New"
}
