package libninja

import (
	"fmt"
	"go/ast"
	"go/token"
)

// indexes is every index or slice bound written as a bare number:
// words[0] says a place and not what lives there. Any index a range or a
// search gives is a name already.
//
// the maintainer's ruling, 10-18: a stack is used through stack operations; a
// fixed set of parts is a struct with named fields; a place that must
// be named is declared along with why.
func (read reading) indexes() []finding {
	var found []finding
	for _, each := range read.files() {
		for node := range ast.Preorder(each.file) {
			for _, bound := range literalBounds(node) {
				seen := fmt.Sprintf("[%s] says a place, "+
					"not what is there.", bound.Value)
				found = append(found, finding{place,
					read.address(each.pkg.PkgPath),
					read.at(bound), seen, ""})
			}
		}
	}
	return found
}

// literalBounds is the bare numbers an index or a slice is written with.
func literalBounds(node ast.Node) []*ast.BasicLit {
	var bounds []ast.Expr
	switch node := node.(type) {
	case *ast.IndexExpr:
		bounds = []ast.Expr{node.Index}
	case *ast.SliceExpr:
		bounds = []ast.Expr{node.Low, node.High, node.Max}
	}
	var found []*ast.BasicLit
	for _, bound := range bounds {
		lit, ok := ast.Unparen(bound).(*ast.BasicLit)
		if ok && lit.Kind == token.INT {
			found = append(found, lit)
		}
	}
	return found
}
