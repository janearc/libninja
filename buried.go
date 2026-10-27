package libninja

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// buriedRoots is how a path that is true on one machine begins.
var buriedRoots = []string{"/Users/", "/home/", "/tmp", "/var/",
	"/private/", "/opt/", "/etc/", "~/"}

// buriedIn is whether a string is one machine's path: it begins at one
// of the roots, or it reaches into somebody's scratch.
func buriedIn(text string) bool {
	for _, root := range buriedRoots {
		if strings.HasPrefix(text, root) {
			return true
		}
	}
	return strings.Contains(text, "scratch/")
}

// buriedPaths is every path literal that is one machine's, handed to
// os or filepath in code that is not a test. Messages that name a path
// are not hits. Calls that open one are.
func (read reading) buriedPaths() []finding {
	var found []finding
	for _, each := range read.files() {
		if each.test() {
			continue
		}
		for node := range ast.Preorder(each.file) {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				continue
			}
			called, ok := filesystemCall(each.pkg, call)
			if !ok {
				continue
			}
			found = append(found,
				read.buriedArgs(each.pkg, call, called)...)
		}
	}
	return found
}

// buriedArgs is the findings for one call's path literals.
func (read reading) buriedArgs(pkg *packages.Package, call *ast.CallExpr,
	called string) []finding {
	var found []finding
	for _, arg := range call.Args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			continue
		}
		text, err := strconv.Unquote(lit.Value)
		if err != nil || !buriedIn(text) {
			continue
		}
		found = append(found, finding{buried, read.address(pkg.PkgPath),
			read.at(lit),
			fmt.Sprintf("%q is one machine's path, "+
				"handed to %s.", text,
				called), ""})
	}
	return found
}

// filesystemCall is whether a call reaches os or path/filepath, and
// what it calls there.
func filesystemCall(pkg *packages.Package, call *ast.CallExpr) (string,
	bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	obj := pkg.TypesInfo.Uses[sel.Sel]
	if obj == nil || obj.Pkg() == nil {
		return "", false
	}
	switch obj.Pkg().Path() {
	case "os", "path/filepath", "io/ioutil":
		return obj.Pkg().Name() + "." + obj.Name(), true
	}
	return "", false
}
