package libninja

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

// handBack is what the scroll tells the student to do, whichever kind
// of struct came back.
const handBack = "On a refusal return nothing usable: a pointer, nil; " +
	"a slice, nil; or change the result."

// feints is every return that hands a struct back beside an error
// that is not nil. The struct is a literal with nothing in it, or a
// variable that was declared and never set before the return; those
// are the empty feints.
//
// A literal with fields set is found as well, and as much: it looks
// like real data, and is often built from a call that failed. Its
// message says to check the error first.
//
// the maintainer's ruling, 10-30: on a refusal, return nothing usable.
// A zero struct beside an error is one a caller will use unchecked.
// Basic types are left alone, and so is every nil.
func (read reading) feints() []finding {
	var found []finding
	for _, each := range read.files() {
		info := each.pkg.TypesInfo
		for node := range ast.Preorder(each.file) {
			body, sig := bodyOf(info, node)
			if body == nil || !refuses(sig) {
				continue
			}
			judge := newJudge(each.pkg.Types, info, body)
			for _, ret := range ownReturns(body) {
				for _, seen := range judge.sayAll(ret, sig) {
					found = append(found, finding{feint,
						read.address(each.pkg.PkgPath),
						read.at(ret), seen, ""})
				}
			}
		}
	}
	return found
}

// bodyOf is the body of a function declaration or literal and its
// signature, or no body for any other node.
func bodyOf(info *types.Info, node ast.Node) (*ast.BlockStmt,
	*types.Signature) {
	switch node := node.(type) {
	case *ast.FuncDecl:
		made, ok := info.Defs[node.Name].(*types.Func)
		if !ok {
			return nil, nil
		}
		return node.Body, made.Signature()
	case *ast.FuncLit:
		sig, _ := info.TypeOf(node).(*types.Signature)
		return node.Body, sig
	}
	return nil, nil
}

// refuses is whether a signature can refuse: it gives results, and the
// last is the error type itself.
func refuses(sig *types.Signature) bool {
	if sig == nil || sig.Results().Len() < minimumResults {
		return false
	}
	results := sig.Results()
	last := results.At(results.Len() - lastBack)
	return types.Identical(last.Type(), errorType)
}

// errorType is the type error, the one every refusal gives.
var errorType = types.Universe.Lookup("error").Type()

// minimumResults is the fewest results a refusing function has: the
// thing and the error, defined so.
const minimumResults = 2

// lastBack is how far back from the end the last of a list is, defined
// so.
const lastBack = 1

// returns is the return statements a walk finds in a function.
type returns struct {
	found []*ast.ReturnStmt
}

// visit keeps a return and goes no further into a function literal,
// which is judged on its own results.
func (all *returns) visit(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.FuncLit:
		return false
	case *ast.ReturnStmt:
		all.found = append(all.found, node)
	}
	return true
}

// ownReturns is the returns of a body that are the function's own.
func ownReturns(body *ast.BlockStmt) []*ast.ReturnStmt {
	var all returns
	ast.Inspect(body, all.visit)
	return all.found
}

// judge knows, for one function, which variables were declared with no
// value and where each was touched, so a return can be judged.
type judge struct {
	here     *types.Package
	info     *types.Info
	declared map[types.Object]bool
	touched  map[types.Object][]token.Pos
}

// newJudge reads a function's body for the variables it declares bare
// and the places that could set them.
func newJudge(here *types.Package, info *types.Info,
	body *ast.BlockStmt) *judge {
	read := &judge{here: here, info: info,
		declared: map[types.Object]bool{},
		touched:  map[types.Object][]token.Pos{}}
	for node := range ast.Preorder(body) {
		read.note(node)
	}
	return read
}

// note records what a node declares bare, or touches.
func (read *judge) note(node ast.Node) {
	switch node := node.(type) {
	case *ast.ValueSpec:
		if len(node.Values) != nothing {
			return
		}
		for _, name := range node.Names {
			read.declared[read.info.Defs[name]] = true
		}
	case *ast.AssignStmt:
		for _, side := range node.Lhs {
			read.touch(side, node.Pos())
		}
	case *ast.IncDecStmt:
		read.touch(node.X, node.Pos())
	case *ast.RangeStmt:
		read.touch(node.Key, node.Pos())
		read.touch(node.Value, node.Pos())
	case *ast.UnaryExpr:
		if node.Op == token.AND {
			read.touch(node.X, node.Pos())
		}
	case *ast.CallExpr:
		read.touchByMethod(node)
	}
}

// touchByMethod records a call of a method that takes a pointer, on a
// variable: the method may set it.
func (read *judge) touchByMethod(call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	chosen := read.info.Selections[sel]
	if chosen == nil || chosen.Kind() != types.MethodVal {
		return
	}
	method, ok := chosen.Obj().(*types.Func)
	if !ok || method.Signature().Recv() == nil {
		return
	}
	recv := method.Signature().Recv().Type()
	if _, pointed := recv.(*types.Pointer); pointed {
		read.touch(sel.X, call.Pos())
	}
}

// touch records that the variable at the root of an expression may be
// set at a place: x, x.field, x[i] and (x) all set x.
func (read *judge) touch(expr ast.Expr, at token.Pos) {
	for expr != nil {
		switch each := expr.(type) {
		case *ast.ParenExpr:
			expr = each.X
		case *ast.SelectorExpr:
			expr = each.X
		case *ast.IndexExpr:
			expr = each.X
		case *ast.StarExpr:
			expr = each.X
		case *ast.Ident:
			object := read.info.Uses[each]
			if object == nil {
				object = read.info.Defs[each]
			}
			read.touched[object] = append(read.touched[object], at)
			return
		default:
			return
		}
	}
}

// setBefore is whether a variable may have been set before a place.
func (read *judge) setBefore(object types.Object, at token.Pos) bool {
	for _, each := range read.touched[object] {
		if each < at {
			return true
		}
	}
	return false
}

// sayAll is what a return says against the scroll: one saying for each
// struct it hands back beside an error that is not nil.
func (read *judge) sayAll(ret *ast.ReturnStmt,
	sig *types.Signature) []string {
	count := sig.Results().Len()
	if len(ret.Results) != count {
		return nil
	}
	final := count - lastBack
	if read.info.Types[ast.Unparen(ret.Results[final])].IsNil() {
		return nil
	}
	var said []string
	for _, each := range ret.Results[:final] {
		if seen := read.say(each, ret); seen != "" {
			said = append(said, seen)
		}
	}
	return said
}

// say is what the scroll says of one result, or nothing when it is
// fine.
func (read *judge) say(expr ast.Expr, ret *ast.ReturnStmt) string {
	kind := read.info.TypeOf(expr)
	if kind == nil || !built(kind) {
		return ""
	}
	name := types.TypeString(types.Unalias(kind), read.qualify)
	switch each := ast.Unparen(expr).(type) {
	case *ast.CompositeLit:
		if len(each.Elts) == nothing {
			return fmt.Sprintf(emptyKind, name, handBack)
		}
		return fmt.Sprintf(filledKind, name, handBack)
	case *ast.Ident:
		object := read.info.Uses[each]
		if read.declared[object] && !read.setBefore(object, ret.Pos()) {
			return fmt.Sprintf(bareKind, each.Name, name, handBack)
		}
	}
	return ""
}

// built is whether a type is a struct, by what it is and not what it
// is called.
func built(kind types.Type) bool {
	_, yes := types.Unalias(kind).Underlying().(*types.Struct)
	return yes
}

// emptyKind, filledKind and bareKind are what the scroll says of the
// three kinds of struct. Each is given the type, then the way out; the
// bare kind is given the variable first.
const emptyKind = "%s{} is an empty value handed back with an error. " +
	"A caller may use it unchecked. %s"

const filledKind = "%s{...} has fields set and is handed back with an " +
	"error. It looks like real data, so a caller who skips the check " +
	"uses it. Check the error first and build only on success: " +
	"if err != nil { return nil, err }. %s"

const bareKind = "%s is a %s declared and never set, handed back with " +
	"an error. A caller may use it unchecked. %s"

// qualify names a package for a type as the master says it: by its
// name, and not at all for the package the code is in.
func (read *judge) qualify(pkg *types.Package) string {
	if pkg == read.here {
		return ""
	}
	return pkg.Name()
}
