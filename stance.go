package libninja

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// longestStance is the most lines a function may take from its first
// line to its closing brace. That is about one screen. It is the maintainer's
// test of whether it can be read at once, as picked-numbers.yml holds it.
var longestStance = picked("longest-stance")

// standAlone is the letters that may be a whole name. They are i the
// iterator and k, and the coordinates x y z. Not r, g and b: code that
// names a colour's channels is doing colour arithmetic by hand. That is
// libtheme's work, and belongs there under whole names.
const standAlone = "ikxyz"

// oneLetter is how long a name of one letter is, defined so.
const oneLetter = 1

// stances is every function longer than a screen.
func (read reading) stances() []finding {
	var found []finding
	for _, each := range read.files() {
		for _, decl := range each.file.Decls {
			function, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			first := read.at(function)
			last := read.fset.Position(function.End()).Line
			lines := spanned(first.Line, last)
			if lines > longestStance {
				found = append(found, finding{stance,
					read.address(each.pkg.PkgPath), first,
					fmt.Sprintf("%s is %d lines.",
						function.Name.Name,
						lines), ""})
			}
		}
	}
	return found
}

// names is every name of one letter that is not allowed to stand
// alone.
func (read reading) names() []finding {
	var found []finding
	for _, each := range read.files() {
		for node := range ast.Preorder(each.file) {
			ident, ok := node.(*ast.Ident)
			if !ok || len(ident.Name) != oneLetter ||
				ident.Name == "_" ||
				strings.Contains(standAlone, ident.Name) ||
				each.pkg.TypesInfo.Defs[ident] == nil ||
				testsT(each, ident) {
				continue
			}
			found = append(found, finding{truename,
				read.address(each.pkg.PkgPath),
				read.at(ident),
				fmt.Sprintf("%s stands alone.",
					ident.Name), ""})
		}
	}
	return found
}

// testsT is whether a name is a test's t. That is t in a test file,
// typed *testing.T or rapid's *T. Go's own tests name it so, as rapid's
// properties do. That t may stand alone and no other t may.
func testsT(each source, ident *ast.Ident) bool {
	defined := each.pkg.TypesInfo.Defs[ident]
	typed := types.TypeString(defined.Type(), nil)
	return ident.Name == "t" && each.test() &&
		(typed == "*testing.T" || typed == "*pgregory.net/rapid.T")
}

// dojos is every test file that sits inside its package instead of
// beside it. The tests of a main are let be: no package outside can
// import a main, so inside is the only place they can stand.
func (read reading) dojos() []finding {
	var found []finding
	for _, each := range read.files() {
		named := each.file.Name.Name
		if !each.test() || strings.HasSuffix(named, "_test") ||
			named == "main" {
			continue
		}
		found = append(found, finding{dojo,
			read.address(each.pkg.PkgPath), read.at(each.file),
			fmt.Sprintf("This test is package %s, inside the wall.",
				named), ""})
	}
	return found
}
