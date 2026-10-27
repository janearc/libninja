package libninja

import (
	"go/ast"
	"go/types"
	"slices"
	"strings"
)

// composed is each call of a call that a call makes, when one of the two
// comes from outside the standard library. Each is said once, as the
// master says it: "bands.Split(swatch.Blackbody(...))".
func composed(info *types.Info, call *ast.CallExpr) []string {
	outer := callee(info, call)
	if outer == nil {
		return nil
	}
	var found []string
	for _, argument := range call.Args {
		inner, ok := ast.Unparen(argument).(*ast.CallExpr)
		if !ok {
			continue
		}
		called := callee(info, inner)
		if called == nil ||
			standard(outer.Pkg().Path()) &&
				standard(called.Pkg().Path()) {
			continue
		}
		said := short(outer) + "(" + short(called) + "(...))"
		if !slices.Contains(found, said) {
			found = append(found, said)
		}
	}
	return found
}

// callee is the package-level function a call calls. It is nil when the
// call is of a method or a builtin. So it is for a conversion or a
// function value.
func callee(info *types.Info, call *ast.CallExpr) *types.Func {
	var ident *ast.Ident
	switch called := ast.Unparen(call.Fun).(type) {
	case *ast.Ident:
		ident = called
	case *ast.SelectorExpr:
		ident = called.Sel
	default:
		return nil
	}
	function, ok := info.Uses[ident].(*types.Func)
	if !ok || function.Pkg() == nil || function.Signature().Recv() != nil {
		return nil
	}
	return function
}

// short is a function as code calls it: its package's name and its
// own.
func short(function *types.Func) string {
	return function.Pkg().Name() + "." + function.Name()
}

// standard is whether a package path is the standard library's, whose
// first element has no dot.
func standard(pkgPath string) bool {
	first, _, _ := strings.Cut(pkgPath, "/")
	return !strings.Contains(first, ".")
}

// fileCount is how many files the places are in.
func fileCount(places []forged) int {
	files := map[string]bool{}
	for _, place := range places {
		files[place.at.Filename] = true
	}
	return len(files)
}
