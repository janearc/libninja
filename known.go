package libninja

import (
	"go/ast"
	"go/token"
	"go/types"
)

// The places of a call's arguments, counted from the first as go counts
// them. Each is named for the position it is, defined so.
const (
	// firstArgument is the argument a call is given first, defined so.
	firstArgument = nothing
	// secondArgument is the argument after the first, defined so.
	secondArgument = firstArgument + lineOrigin
	// thirdArgument is the argument after the second, defined so.
	thirdArgument = secondArgument + lineOrigin
	// fourthArgument is the argument after the third, defined so.
	fourthArgument = thirdArgument + lineOrigin
	// fifthArgument is the argument after the fourth, defined so.
	fifthArgument = fourthArgument + lineOrigin
)

// knownArguments is the places in a call's arguments where a bare
// number is not a guess, by the call's full name. The call is known by
// its package and its type, never by its name alone. The module's own
// function called WriteFile is not os's.
//
// Modes are known by their place in os and io/fs. The system writes
// them that way and a file mode has no other name. The places in strconv tell
// a precision and a bit size. They count digits and bits, and the call
// says which.
func knownArguments(full string) []int {
	switch full {
	case "os.WriteFile", "io/ioutil.WriteFile", "os.OpenFile",
		"(*os.Root).OpenFile", "(*os.Root).WriteFile":
		return []int{thirdArgument}
	case "os.Mkdir", "os.MkdirAll", "os.Chmod", "(*os.Root).Mkdir",
		"(*os.Root).MkdirAll", "(*os.Root).Chmod":
		return []int{secondArgument}
	case "(*os.File).Chmod":
		return []int{firstArgument}
	case "strconv.ParseFloat":
		return []int{secondArgument}
	case "strconv.ParseInt", "strconv.ParseUint":
		return []int{thirdArgument}
	case "strconv.FormatFloat":
		return []int{thirdArgument, fourthArgument}
	case "strconv.AppendFloat":
		return []int{fourthArgument, fifthArgument}
	}
	return nil
}

// fileMode is the type a mode is converted to, as io/fs names it. os's
// FileMode is the same type under another name.
const fileMode = "io/fs.FileMode"

// modeType is whether a type is io/fs.FileMode, seen through aliases.
func modeType(kind types.Type) bool {
	return types.Unalias(kind).String() == fileMode
}

// knownBy is the number literals a call is given where its place knows
// them. Those are a mode where a call takes a mode, a precision or a bit
// size where strconv takes one, the conversion to a file mode and the -1
// that asks a regexp for all its matches.
func knownBy(info *types.Info, call *ast.CallExpr) []*ast.BasicLit {
	var places []int
	only := ""
	if info.Types[call.Fun].IsType() {
		if modeType(info.TypeOf(call.Fun)) {
			places = []int{firstArgument}
		}
	} else if called := calledFunction(info, call); called != nil {
		places = knownArguments(called.FullName())
		if allMatches(called.FullName()) {
			places, only = []int{secondArgument}, minusOne
		}
	}
	var found []*ast.BasicLit
	for _, place := range places {
		if place >= len(call.Args) {
			continue
		}
		lit, text := signed(call.Args[place])
		if lit != nil && (only == "" || text == only) {
			found = append(found, lit)
		}
	}
	return found
}

// calledFunction is the function or method a call calls. It is nil for
// a call of a value, a builtin or a conversion.
func calledFunction(info *types.Info, call *ast.CallExpr) *types.Func {
	var ident *ast.Ident
	switch called := ast.Unparen(call.Fun).(type) {
	case *ast.Ident:
		ident = called
	case *ast.SelectorExpr:
		ident = called.Sel
	}
	function, _ := info.Uses[ident].(*types.Func)
	return function
}

// allMatches is whether a call is a regexp method whose second argument
// is a count of matches, where -1 is the call's own word for all.
func allMatches(full string) bool {
	switch full {
	case "(*regexp.Regexp).FindAll", "(*regexp.Regexp).FindAllIndex",
		"(*regexp.Regexp).FindAllString",
		"(*regexp.Regexp).FindAllStringIndex",
		"(*regexp.Regexp).FindAllStringSubmatch",
		"(*regexp.Regexp).FindAllStringSubmatchIndex",
		"(*regexp.Regexp).FindAllSubmatch",
		"(*regexp.Regexp).FindAllSubmatchIndex",
		"(*regexp.Regexp).Split":
		return true
	}
	return false
}

// minusOne is a number as it is written when it asks for all matches.
const minusOne = "-1"

// signed is the literal in an argument and its text with the minus sign
// before it, or nil when the argument is not a number written out.
func signed(arg ast.Expr) (*ast.BasicLit, string) {
	arg = ast.Unparen(arg)
	sign := ""
	if minus, ok := arg.(*ast.UnaryExpr); ok && minus.Op == token.SUB {
		arg, sign = ast.Unparen(minus.X), "-"
	}
	lit, ok := arg.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return nil, ""
	}
	return lit, sign + lit.Value
}
