package libninja

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// deepest is how many calls may stand inside one another, through their
// arguments: f(g(x)) is fine and f(g(h(x))) is a hit. That is the maintainer's
// line from 10-18, as picked-numbers.yml holds it.
var deepest = picked("deepest")

// nests is every call that holds calls inside its arguments past
// deepest. Conversions are calls; a function literal starts its own
// count; a call deep inside one already found is not found again.
func (read reading) nests() []finding {
	var found []finding
	for _, each := range read.files() {
		inside := map[*ast.CallExpr]bool{}
		for node := range ast.Preorder(each.file) {
			call, ok := node.(*ast.CallExpr)
			if !ok || inside[call] {
				continue
			}
			path := deepestPath(call)
			if len(path) <= deepest {
				continue
			}
			markInside(call, inside)
			seen := fmt.Sprintf("%s is %d calls deep.",
				spoken(path), len(path))
			found = append(found, finding{nested,
				read.address(each.pkg.PkgPath),
				read.at(call), seen, ""})
		}
	}
	return found
}

// deepestPath is the calls from this one down its deepest line of calls
// through arguments, outermost first.
func deepestPath(call *ast.CallExpr) []*ast.CallExpr {
	var below []*ast.CallExpr
	for _, arg := range call.Args {
		for _, inner := range callsIn(arg) {
			if path := deepestPath(inner); len(path) > len(below) {
				below = path
			}
		}
	}
	return append([]*ast.CallExpr{call}, below...)
}

// callsIn is the outermost calls in an expression. It does not look
// inside a call, which deepestPath does. It does not look inside a
// function literal, which counts its own.
func callsIn(expr ast.Expr) []*ast.CallExpr {
	var found outermost
	ast.Inspect(expr, found.visit)
	return found.calls
}

// outermost is the outermost calls a walk has found.
type outermost struct {
	calls []*ast.CallExpr
}

// visit keeps a call and goes no further into it, and goes no further
// into a function literal; anything else it walks on into.
func (found *outermost) visit(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.CallExpr:
		found.calls = append(found.calls, node)
		return false
	case *ast.FuncLit:
		return false
	}
	return true
}

// markInside marks every call in a call's arguments to any depth short
// of a function literal, as found with it.
func markInside(call *ast.CallExpr, inside map[*ast.CallExpr]bool) {
	for _, arg := range call.Args {
		for _, inner := range callsIn(arg) {
			inside[inner] = true
			markInside(inner, inside)
		}
	}
}

// spoken is a line of calls as written, each holding the next:
// append(fmt.Sprintf(shown(...))).
func spoken(path []*ast.CallExpr) string {
	var text strings.Builder
	for _, call := range path {
		text.WriteString(types.ExprString(call.Fun))
		text.WriteString("(")
	}
	text.WriteString("...")
	closing := strings.Repeat(")", len(path))
	text.WriteString(closing)
	return text.String()
}
