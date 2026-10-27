package libninja_test

import (
	"strings"
	"testing"
)

// hitsIn is what one scroll saw in a house holding one more file.
// That file is more.go, with the given body below its package line.
func hitsIn(test *testing.T, scroll, body string) []string {
	test.Helper()
	dir := house(test, ".", map[string]string{
		".editorconfig": "root = true\n\n[*.go]\nindent_style = tab\n",
		"more.go":       "package house\n\n" + body,
	})
	_, facts := thrown(test, dir)
	var seen []string
	for _, each := range facts["hits"] {
		if each["scroll"] == scroll && each["file"] == "more.go" {
			seen = append(seen, each["seen"].(string))
		}
	}
	return seen
}

func TestANumberNamedWithWhyIsFine(test *testing.T) {
	body := "// most is the most there may be.\nconst most = 3\n\n" +
		"const least = 1 // the least there may be\n\n" +
		"// Many is whether there are many.\n" +
		"func Many(count int) bool {\n" +
		"\tenough := 2 // two is company\n" +
		"\treturn count > enough && count < most && count > least\n}\n"
	if seen := hitsIn(test, "number", body); seen != nil {
		test.Errorf("named, with why: %v", seen)
	}
}

func TestANumberNamedWithoutWhyIsBare(test *testing.T) {
	body := "const most = 3\n\n" +
		"// Many is whether there are many.\n" +
		"func Many(count int) bool { return count > most }\n"
	seen := hitsIn(test, "number", body)
	if strings.Join(seen, "|") != "3 is a bare number." {
		test.Errorf("named, without why: %v", seen)
	}
}

func TestACommentOnTheLineAboveMustStandAlone(test *testing.T) {
	body := "// Many is whether there are many.\n" +
		"func Many(count int) bool {\n" +
		"\tseen := count // how many came\n" +
		"\tenough := 2\n" +
		"\treturn seen > enough\n}\n"
	seen := hitsIn(test, "number", body)
	if strings.Join(seen, "|") != "2 is a bare number." {
		test.Errorf("the comment above trails other code: %v", seen)
	}
}

func TestZeroAndOneAreNumbersToo(test *testing.T) {
	body := "// Some is whether there is one or more.\n" +
		"func Some(count int) bool { return count > 0 || count == 1 }\n"
	seen := hitsIn(test, "number", body)
	if strings.Join(seen, "|") != "0 is a bare number.|1 is a bare number." {
		test.Errorf("zero and one: %v", seen)
	}
}

func TestAnArrayLengthIsANumber(test *testing.T) {
	body := "// Sum is a sum.\nvar Sum [32]byte\n"
	if seen := hitsIn(test, "number", body); seen == nil {
		test.Error("[32]byte was let by")
	}
}

func TestATableIsFoundOnceAtItsTop(test *testing.T) {
	body := "// weights are lamb's.\n" +
		"var weights = [][]float64{{0.2126, 0.7152}, {-0.0722, 1}}\n"
	seen := hitsIn(test, "table", body)
	if strings.Join(seen, "|") !=
		"Numbers are kept here as a table in code." {
		test.Errorf("a table of rows: %v", seen)
	}
	if bare := hitsIn(test, "number", body); bare != nil {
		test.Errorf("a table's numbers were found again: %v", bare)
	}
}

// Tables handed to a call have no declaration to say why. Their numbers
// are still the table's. They are not bare.
func TestATableInACallHasNoBareNumbers(test *testing.T) {
	body := "// Sum is a sum.\nfunc Sum(row []int) int { return len(row) }\n\n" +
		"// Total sums a table.\nfunc Total() int { return Sum([]int{3, 4}) }\n"
	seen := hitsIn(test, "table", body)
	if strings.Join(seen, "|") !=
		"Numbers are kept here as a table in code." {
		test.Errorf("a table in a call: %v", seen)
	}
	if bare := hitsIn(test, "number", body); bare != nil {
		test.Errorf("a table's numbers in a call were found again: %v",
			bare)
	}
}

func TestOneNumberInALiteralIsNotATable(test *testing.T) {
	body := "// Point is a place.\ntype Point struct{ X, Y int }\n\n" +
		"// Up is a step up.\nvar Up = Point{Y: 1}\n"
	if seen := hitsIn(test, "table", body); seen != nil {
		test.Errorf("one number: %v", seen)
	}
}
