package libninja

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

// numbers is every number written bare: anywhere but the value of a
// declaration that names it and says why, by a comment on its line or
// the line above. 0 and 1 are numbers too. the maintainer's ruling, 10-18: a
// number in code says what it is and where it came from.
//
// Any number in an index is the index scroll's and one in a table is the
// table scroll's, so neither is found twice.
func (read reading) numbers() []finding {
	var found []finding
	for _, each := range read.files() {
		marks := markNumbers(read.fset, each)
		for node := range ast.Preorder(each.file) {
			lit, ok := node.(*ast.BasicLit)
			if !ok || !isNumber(lit) || marks.spoken[lit] {
				continue
			}
			seen := fmt.Sprintf("%s is a bare number.", lit.Value)
			found = append(found, finding{bare,
				read.address(each.pkg.PkgPath), read.at(lit),
				seen, ""})
		}
	}
	return found
}

// isNumber is whether a literal is a number: whole, real or imaginary.
func isNumber(lit *ast.BasicLit) bool {
	switch lit.Kind {
	case token.INT, token.FLOAT, token.IMAG:
		return true
	}
	return false
}

// numberMarks is a file's numbers that are accounted for: named in a
// declaration that says why, or another scroll's to find. trailing is
// the lines a comment ends after code on; alone, the lines a comment
// standing by itself ends on.
type numberMarks struct {
	fset            *token.FileSet
	info            *types.Info
	trailing, alone map[int]bool
	spoken          map[*ast.BasicLit]bool
}

// markNumbers is the numbers in a file that are not bare.
func markNumbers(fset *token.FileSet, each source) numberMarks {
	marks := numberMarks{fset: fset, info: each.pkg.TypesInfo,
		trailing: map[int]bool{},
		alone:    map[int]bool{}, spoken: map[*ast.BasicLit]bool{}}
	body := textOf(each.name)
	for _, group := range each.file.Comments {
		end := fset.Position(group.End()).Line
		if writtenAlone(fset, body, group.Pos()) {
			marks.alone[end] = true
		} else {
			marks.trailing[end] = true
		}
	}
	for node := range ast.Preorder(each.file) {
		marks.mark(node)
	}
	return marks
}

// mark marks the numbers a node accounts for.
func (marks numberMarks) mark(node ast.Node) {
	switch node := node.(type) {
	case *ast.ValueSpec:
		if marks.says(node) {
			marks.values(node.Values)
		}
	case *ast.AssignStmt:
		if node.Tok == token.DEFINE && marks.says(node) {
			marks.values(node.Rhs)
		}
	case *ast.CompositeLit:
		if tabled(node) {
			marks.all(node)
		}
	case *ast.CallExpr:
		for _, known := range knownBy(marks.info, node) {
			marks.spoken[known] = true
		}
	default:
		for _, bound := range literalBounds(node) {
			marks.spoken[bound] = true
		}
	}
}

// says is whether a declaration says why: a comment after it on its own
// line, or one standing alone on the line above it.
func (marks numberMarks) says(node ast.Node) bool {
	line := marks.fset.Position(node.Pos()).Line
	return marks.trailing[line] || marks.alone[line-above]
}

// above is how far the line over a declaration is from it, defined as
// one line up.
const above = 1

// values marks every number in the values a declaration names, and not
// those in its type: [32]byte is a bare 32 whatever it is called.
func (marks numberMarks) values(named []ast.Expr) {
	for _, value := range named {
		marks.all(value)
	}
}

// all marks every number inside a node as spoken for.
func (marks numberMarks) all(node ast.Node) {
	for inner := range ast.Preorder(node) {
		if lit, ok := inner.(*ast.BasicLit); ok {
			marks.spoken[lit] = true
		}
	}
}
