package libninja

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"
)

// smallestBlade is the fewest syntax nodes a body has before two
// bodies alike count as one forged twice; below it, alike is usually
// only short. Picked, as picked-numbers.yml holds it.
var smallestBlade = picked("smallest-blade")

// fewestCopies is how many files must make the same call of a call
// before it is a blade forged in each: the rule of three, picked.
var fewestCopies = picked("fewest-copies")

// forged is one place a blade is made. It holds where it is and which
// package made it. It also holds what the blade is called there.
type forged struct {
	at      token.Position
	student string
	name    string
}

// blades is every function body that has the same shape as another,
// names of its own locals aside.
func (read reading) blades() []finding {
	alike := map[[sha256.Size]byte][]forged{}
	for _, each := range read.files() {
		if each.test() {
			continue
		}
		for _, decl := range each.file.Decls {
			function, ok := decl.(*ast.FuncDecl)
			if !ok || function.Body == nil {
				continue
			}
			shape, nodes := fingerprint(each.pkg.TypesInfo,
				function.Body)
			if nodes >= smallestBlade {
				sum := sha256.Sum256([]byte(shape))
				alike[sum] = append(alike[sum], forged{
					read.at(function),
					read.address(each.pkg.PkgPath),
					function.Name.Name})
			}
		}
	}
	var found []finding
	for _, places := range alike {
		found = append(found, read.copies(places,
			"The body of %s is forged "+
				"again at %s.")...)
	}
	return found
}

// twice is how many times a blade is forged before it is copies, as
// copies are defined. Once is a blade and twice is one forged again.
const twice = 2

// copies is a finding at every place a blade is forged, when it is
// forged more than once. Each finding names the other places by file and
// line, as every place in a report is named.
func (read reading) copies(places []forged, said string) []finding {
	if len(places) < twice {
		return nil
	}
	var found []finding
	for _, place := range places {
		var others []string
		for _, other := range places {
			if other != place {
				file := relativeTo(read.dir, other.at)
				at := fmt.Sprintf("%s:%d", file, other.at.Line)
				others = append(others, at)
			}
		}
		slices.Sort(others)
		elsewhere := strings.Join(others, ", ")
		found = append(found, finding{blade, place.student, place.at,
			fmt.Sprintf(said, place.name, elsewhere), ""})
	}
	return found
}

// fingerprint is a body's shape and how many nodes it has. The shape
// holds each node's kind, and each operator and literal. Names from
// outside the body stand as themselves. Names made inside it stand by
// the order they appear, so renaming a local changes nothing.
func fingerprint(info *types.Info, body *ast.BlockStmt) (string, int) {
	var shape strings.Builder
	locals := map[types.Object]int{}
	var nodes int
	for node := range ast.Preorder(body) {
		nodes++
		fmt.Fprintf(&shape, "%T", node)
		switch part := node.(type) {
		case *ast.Ident:
			shape.WriteString(identShape(info, part, locals))
		case *ast.BasicLit:
			shape.WriteString(part.Value)
		case *ast.BinaryExpr:
			shape.WriteString(part.Op.String())
		case *ast.UnaryExpr:
			shape.WriteString(part.Op.String())
		case *ast.AssignStmt:
			shape.WriteString(part.Tok.String())
		case *ast.IncDecStmt:
			shape.WriteString(part.Tok.String())
		case *ast.BranchStmt:
			shape.WriteString(part.Tok.String())
		}
		shape.WriteByte(';')
	}
	return shape.String(), nodes
}

// identShape is a name as a fingerprint keeps it.
func identShape(info *types.Info, ident *ast.Ident,
	locals map[types.Object]int) string {
	obj := info.ObjectOf(ident)
	switch {
	case obj == nil:
		return ident.Name
	case obj.Pkg() == nil || obj.Parent() == nil:
		return obj.Name()
	case obj.Parent() == obj.Pkg().Scope():
		return obj.Pkg().Path() + "." + obj.Name()
	}
	if imported, ok := obj.(*types.PkgName); ok {
		return imported.Imported().Path()
	}
	if _, known := locals[obj]; !known {
		locals[obj] = len(locals)
	}
	return fmt.Sprintf("local%d", locals[obj])
}

// compositions is every call of a call that is made in at least
// fewestCopies files, where f or g comes from outside the standard
// library. In f(g(...)), such a call is a blade the owner of f or g
// should forge once.
func (read reading) compositions() []finding {
	made := map[string][]forged{}
	for _, each := range read.files() {
		if each.test() {
			continue
		}
		for node := range ast.Preorder(each.file) {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				continue
			}
			info := each.pkg.TypesInfo
			for _, said := range composed(info, call) {
				made[said] = append(made[said], forged{
					read.at(call),
					read.address(each.pkg.PkgPath),
					said})
			}
		}
	}
	var found []finding
	for _, places := range made {
		if fileCount(places) >= fewestCopies {
			found = append(found, read.copies(places,
				"%s is forged here, "+
					"and again at %s.")...)
		}
	}
	return found
}
