package libninja_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/janearc/puzzlemouth/numbers"
)

// pickedAsked is every name libninja's own code asks picked for, by the
// file that asks.
func pickedAsked(test *testing.T) map[string]string {
	test.Helper()
	files, err := filepath.Glob("*.go")
	if err != nil {
		test.Fatal(err)
	}
	asked := map[string]string{}
	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(fset, file, nil,
			parser.SkipObjectResolution)
		if err != nil {
			test.Fatal(err)
		}
		for node := range ast.Preorder(parsed) {
			for _, name := range pickedNames(node) {
				asked[name] = file
			}
		}
	}
	return asked
}

// pickedNames is the names a call of picked asks for, if a node is one.
func pickedNames(node ast.Node) []string {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}
	if called, ok := call.Fun.(*ast.Ident); !ok || called.Name != "picked" {
		return nil
	}
	var names []string
	for _, arg := range call.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			name, _ := strconv.Unquote(lit.Value)
			names = append(names, name)
		}
	}
	return names
}

func TestEveryPickedNumberIsInTheFile(test *testing.T) {
	raw, err := os.ReadFile("picked-numbers.yml")
	if err != nil {
		test.Fatal(err)
	}
	chosen, err := numbers.Read("picked-numbers.yml", raw)
	if err != nil {
		test.Fatal(err)
	}
	asked := pickedAsked(test)
	if _, found := asked["longest-stance"]; !found {
		test.Fatal("the stance scroll's call of picked was not found, " +
			"so no call was read")
	}
	for name, file := range asked {
		number, err := chosen.Get(name)
		if err != nil {
			test.Errorf("%s asks for %s: %v", file, name, err)
			continue
		}
		if value := number.Value(); value != math.Trunc(value) {
			test.Errorf("%s asks for %s, which is %v, not a whole count",
				file, name, value)
		}
	}
}
