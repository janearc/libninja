package libninja

import "go/ast"

// row is how many numbers side by side make a table, picked.
var row = picked("row")

// tables is every table of numbers written in code: a composite literal
// holding a row of numbers among its own elements, found once at its
// outermost.
//
// the maintainer's ruling, 10-18: data lives in data files and is pinned by its
// source. The code reads it. Neither a name nor a comment makes a table
// into anything but data in code.
func (read reading) tables() []finding {
	var found []finding
	for _, each := range read.files() {
		inside := map[*ast.CompositeLit]bool{}
		for node := range ast.Preorder(each.file) {
			composite, ok := node.(*ast.CompositeLit)
			if !ok || inside[composite] ||
				!holdsATable(composite) {
				continue
			}
			for inner := range ast.Preorder(composite) {
				nested, ok := inner.(*ast.CompositeLit)
				if ok {
					inside[nested] = true
				}
			}
			const seen = "Numbers are kept here as a table in code."
			found = append(found, finding{tableInCode,
				read.address(each.pkg.PkgPath),
				read.at(composite), seen, ""})
		}
	}
	return found
}

// holdsATable is whether a composite literal is a table or holds one.
func holdsATable(composite *ast.CompositeLit) bool {
	for inner := range ast.Preorder(composite) {
		nested, ok := inner.(*ast.CompositeLit)
		if ok && tabled(nested) {
			return true
		}
	}
	return false
}

// tabled is whether a composite literal's own elements are a row of
// numbers or more, a minus sign or parentheses seen through.
func tabled(composite *ast.CompositeLit) bool {
	var count int
	for _, element := range composite.Elts {
		if pair, ok := element.(*ast.KeyValueExpr); ok {
			element = pair.Value
		}
		if numeral(element) {
			count++
		}
	}
	return count >= row
}

// numeral is whether an expression is a number written out, a minus
// sign or parentheses seen through.
func numeral(expr ast.Expr) bool {
	expr = ast.Unparen(expr)
	if sign, ok := expr.(*ast.UnaryExpr); ok {
		expr = ast.Unparen(sign.X)
	}
	lit, ok := expr.(*ast.BasicLit)
	return ok && isNumber(lit)
}
