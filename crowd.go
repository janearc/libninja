package libninja

import (
	"cmp"
	"fmt"
	"go/ast"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// testsPerVow is how many test functions a module has at least to each
// vow it carries, picked.
var testsPerVow = picked("tests-per-vow")

// testPrefix is how go names a test function, as go test defines it.
const testPrefix = "Test"

// crowd is the module's vows when they are more than one to every
// testsPerVow tests. Then the code is running from a fight, or a scroll
// is wrong and belongs with the master; either way a vow is the
// cheapest way to hide it.
//
// the maintainer's ruling, 10-25: one vow in twenty tests is the most.
func (read reading) crowd() []finding {
	tests := read.testCount()
	if none(len(read.vows)) || len(read.vows)*testsPerVow <= tests {
		return nil
	}
	first := read.firstVow()
	seen := fmt.Sprintf("The module carries %d vows over %d tests, "+
		"more than one in %d. Bring your vows to the master.",
		len(read.vows), tests, testsPerVow)
	return []finding{{crowdedVows, first.student, first.at, seen, ""}}
}

// firstVow is the vow written first, by file and line. The crowd is
// said there every run.
func (read reading) firstVow() vow {
	var all []vow
	for _, written := range read.vows {
		all = append(all, written)
	}
	slices.SortFunc(all, teaching{read.dir}.written)
	return all[nothing]
}

// written orders two vows by file and line. Each file is as the report
// says it, from the module's root.
func (order teaching) written(one, other vow) int {
	file, otherFile := relativeTo(order.dir, one.at),
		relativeTo(order.dir, other.at)
	return cmp.Or(cmp.Compare(file, otherFile),
		cmp.Compare(one.at.Line, other.at.Line))
}

// testCount is how many test functions the module's test files declare.
// Each is a function, not a method, named Test and then anything but a
// lower case letter. It takes one *testing.T.
func (read reading) testCount() int {
	var count int
	for _, each := range read.files() {
		if !strings.HasSuffix(each.name, "_test.go") {
			continue
		}
		for _, declared := range each.file.Decls {
			function, ok := declared.(*ast.FuncDecl)
			if ok && isTest(function) {
				count++
			}
		}
	}
	return count
}

// isTest is whether a declaration is a test function as go test finds
// one.
func isTest(function *ast.FuncDecl) bool {
	if function.Recv != nil || function.Type.TypeParams != nil {
		return false
	}
	rest, named := strings.CutPrefix(function.Name.Name, testPrefix)
	if !named {
		return false
	}
	letter, _ := utf8.DecodeRuneInString(rest)
	params := function.Type.Params.List
	if unicode.IsLower(letter) || len(params) != lonely {
		return false
	}
	param := params[nothing]
	return len(param.Names) <= lonely && isTestingT(param.Type)
}

// lonely is how many parameters a test function takes, defined so by go
// test.
const lonely = 1

// isTestingT is whether a type is written *testing.T.
func isTestingT(kind ast.Expr) bool {
	star, ok := kind.(*ast.StarExpr)
	if !ok {
		return false
	}
	selector, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && pkg.Name == "testing" && selector.Sel.Name == "T"
}
