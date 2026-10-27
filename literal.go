package libninja

import "go/ast"

// literals is every function literal in the module: a function written
// where it is used, with no name.
//
// the maintainer's ruling, 10-18: a goroutine runs a named method. Any api that
// takes a function takes a named one. What the literal would have held
// lives in a type, where it is seen.
func (read reading) literals() []finding {
	var found []finding
	for _, each := range read.files() {
		for node := range ast.Preorder(each.file) {
			literal, ok := node.(*ast.FuncLit)
			if !ok {
				continue
			}
			found = append(found, finding{nameless,
				read.address(each.pkg.PkgPath),
				read.at(literal),
				"A function is written here with no name.", ""})
		}
	}
	return found
}
