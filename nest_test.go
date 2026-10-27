package libninja_test

import "testing"

// nestsIn is what the nest scroll saw in a house holding one more file.
// That file is nest.go, with the given body below its package line.
func nestsIn(test *testing.T, body string) []map[string]any {
	test.Helper()
	dir := house(test, ".", map[string]string{
		".editorconfig": "root = true\n\n[*.go]\nindent_style = tab\n",
		"nest.go":       "package house\n\n" + body,
	})
	_, facts := thrown(test, dir)
	var nests []map[string]any
	for _, each := range facts["hits"] {
		if each["scroll"] == "nest" {
			nests = append(nests, each)
		}
	}
	return nests
}

// alone is what the one finding saw, or nothing when there is not
// exactly one.
func alone(found []map[string]any) any {
	var only map[string]any
	for _, each := range found {
		if only != nil {
			return nil
		}
		only = each
	}
	if only == nil {
		return nil
	}
	return only["seen"]
}

// three is a call three deep, in a function the other tests share.
const three = "func one(a int) int { return a }\n\n" +
	"func three() int { return one(one(one(1))) }\n"

func TestThreeCallsDeepIsAHit(test *testing.T) {
	found := nestsIn(test, three)
	if alone(found) != "one(one(one(...))) is 3 calls deep." {
		test.Errorf("three deep: %v", found)
	}
}

func TestTwoCallsDeepIsFine(test *testing.T) {
	body := "func one(a int) int { return a }\n\n" +
		"func two() int { return one(one(1)) }\n"
	if found := nestsIn(test, body); found != nil {
		test.Errorf("two deep: %v", found)
	}
}

func TestAConversionIsACall(test *testing.T) {
	body := "func one(a float64) float64 { return a }\n\n" +
		"func three() int { return int(one(float64(1))) }\n"
	if found := nestsIn(test, body); alone(found) == nil {
		test.Errorf("a conversion three deep: %v", found)
	}
}

func TestAFunctionLiteralStartsItsOwnCount(test *testing.T) {
	body := "func one(a int) int { return a }\n\n" +
		"func run(do func() int) int { return do() }\n\n" +
		"func closed() int {\n" +
		"\treturn one(run(func() int { return one(one(1)) }))\n}\n"
	if found := nestsIn(test, body); found != nil {
		test.Errorf("two deep, then two inside a literal: %v", found)
	}
}

func TestAChainIsNotNesting(test *testing.T) {
	body := "type chain struct{}\n\n" +
		"func (chain) next() chain { return chain{} }\n\n" +
		"func one(c chain) chain { return c }\n\n" +
		"func chained() chain {\n" +
		"\treturn one(chain{}.next().next().next())\n}\n"
	if found := nestsIn(test, body); found != nil {
		test.Errorf("a chain handed to one call: %v", found)
	}
}

func TestADeepLineIsFoundOnceAtItsTop(test *testing.T) {
	body := "func one(a int) int { return a }\n\n" +
		"func five() int { return one(one(one(one(one(1))))) }\n"
	found := nestsIn(test, body)
	if alone(found) != "one(one(one(one(one(...))))) is 5 calls deep." {
		test.Errorf("five deep: %v", found)
	}
}

// When two arguments go equally deep, the line said is the first.
func TestOfTwoLinesEquallyDeepTheFirstIsSaid(test *testing.T) {
	body := "func one(a int) int { return a }\n\n" +
		"func two(a, b int) int { return a + b }\n\n" +
		"func first(a int) int { return a }\n\n" +
		"func second(a int) int { return a }\n\n" +
		"func both() int { return two(first(one(1)), second(one(1))) }\n"
	found := nestsIn(test, body)
	if alone(found) != "two(first(one(...))) is 3 calls deep." {
		test.Errorf("two lines equally deep: %v", found)
	}
}
