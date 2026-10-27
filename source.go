package libninja

import (
	"go/ast"
	"go/token"
	"regexp"
)

// defined is the word a named number's comment says when the number is
// true by definition, and so may live in code.
var defined = regexp.MustCompile(`(?i)\bdefined\b`)

// sourcelessSeen is what the scroll sees at a named number that is not
// defined.
const sourcelessSeen = "A number named here says why, but not that " +
	"it is defined; measured or picked, it lives in data."

// sourcelessNumbers is every number named in code that its comment does
// not call defined. The maintainer's rulings, 10-18: a measured number
// is pulled from its source into the research data; a picked number is
// config; code holds only what is true by definition and says so.
func (read reading) sourcelessNumbers() []finding {
	var found []finding
	for _, each := range read.files() {
		said := commentsByLine(read.fset, each)
		for node := range ast.Preorder(each.file) {
			at, values, ok := naming(node)
			if !ok || !holdsANumber(values) {
				continue
			}
			line := read.at(at).Line
			why := said.trailing[line] + said.alone[line-above]
			if why == "" || defined.MatchString(why) {
				continue
			}
			found = append(found, finding{sourceless,
				read.address(each.pkg.PkgPath), read.at(at),
				sourcelessSeen, ""})
		}
	}
	return found
}

// naming is a declaration that names values: a const or var spec or a
// := with its values. It is also whether the node is one.
func naming(node ast.Node) (ast.Node, []ast.Expr, bool) {
	switch node := node.(type) {
	case *ast.ValueSpec:
		return node, node.Values, true
	case *ast.AssignStmt:
		return node, node.Rhs, node.Tok == token.DEFINE
	}
	return nil, nil, false
}

// holdsANumber is whether any of the values has a number written in it.
func holdsANumber(values []ast.Expr) bool {
	for _, value := range values {
		for inner := range ast.Preorder(value) {
			lit, ok := inner.(*ast.BasicLit)
			if ok && isNumber(lit) {
				return true
			}
		}
	}
	return false
}

// remarks is a file's comments by the line they end on. The trailing ones
// come after code on that line. The alone ones stand by themselves.
type remarks struct {
	trailing, alone map[int]string
}

// commentsByLine is a file's comments, by the line each ends on.
func commentsByLine(fset *token.FileSet, each source) remarks {
	said := remarks{trailing: map[int]string{}, alone: map[int]string{}}
	body := textOf(each.name)
	for _, group := range each.file.Comments {
		end := fset.Position(group.End()).Line
		if writtenAlone(fset, body, group.Pos()) {
			said.alone[end] = group.Text()
		} else {
			said.trailing[end] = group.Text()
		}
	}
	return said
}
