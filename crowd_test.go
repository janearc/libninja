package libninja_test

import (
	"fmt"
	"strings"
	"testing"
)

// crowdIn is what the crowd scroll saw in a house. The house has the
// given count of vows, each a well, and the given count of test
// functions. Any other code goes in its test file below them.
func crowdIn(test *testing.T, vows, tests int, other string) []string {
	test.Helper()
	var code, checks strings.Builder
	code.WriteString("package house\n")
	for each := range vows {
		fmt.Fprintf(&code, "\n// Open%d is open.\nvar Open%d = "+
			"map[string]string{} //ninja:vow well the wire reads it\n", each, each)
	}
	checks.WriteString("package house\n\nimport \"testing\"\n\nvar _ testing.TB\n")
	for each := range tests {
		fmt.Fprintf(&checks, "\nfunc Test%d(t *testing.T) {}\n", each)
	}
	checks.WriteString(other)
	dir := house(test, ".", map[string]string{
		".editorconfig": "root = true\n\n[*.go]\nindent_style = tab\n",
		"more.go":       code.String(),
		"more_test.go":  checks.String(),
	})
	_, facts := thrown(test, dir)
	var seen []string
	for _, each := range facts["hits"] {
		if each["scroll"] == "crowd" {
			seen = append(seen, each["file"].(string)+" "+
				each["seen"].(string))
		}
	}
	return seen
}

// said is the master's words when vows crowd tests.
func said(vows, tests int) string {
	return fmt.Sprintf("more.go The module carries %d vows over %d tests, "+
		"more than one in 20. Bring your vows to the master.",
		vows, tests)
}

// How many vows a house is drawn with, in the tests below. Each is a
// count and no more, defined so.
const (
	// oneVow is a single vow, defined so.
	oneVow = 1
	// twoVows is a pair of vows, defined so.
	twoVows = 2
	// noVows is a house that has made none, defined so.
	noVows = 0
	// noTests is a house that tests nothing, defined so.
	noTests = 0
	// oneTest is a single test function, defined so.
	oneTest = 1
)

// ceiling is how many tests a vow needs to stand over, held in
// testdata/numbers.yaml.
func ceiling() int {
	return int(heldTo("tests-per-vow"))
}

// One vow to twenty tests is the most there may be. So is two to forty.
func TestOneVowInTwentyTestsIsNotCrowded(test *testing.T) {
	for _, vows := range []int{oneVow, twoVows} {
		tests := vows * ceiling()
		if seen := crowdIn(test, vows, tests, ""); seen != nil {
			test.Errorf("%d vows, %d tests: %v", vows, tests, seen)
		}
	}
}

// One test fewer, and the vows crowd the tests.
func TestOneTestFewerThanTwentyToAVowIsCrowded(test *testing.T) {
	for _, vows := range []int{oneVow, twoVows} {
		tests := vows*ceiling() - oneTest
		seen := crowdIn(test, vows, tests, "")
		if strings.Join(seen, "|") != said(vows, tests) {
			test.Errorf("%d vows, %d tests: %q", vows, tests, seen)
		}
	}
}

// What is not a test function is not counted. That means a main for
// tests, an example, a benchmark, a method, or a function that takes
// another type.
func TestOnlyTestFunctionsAreCounted(test *testing.T) {
	other := "\nfunc TestMain(m *testing.M) {}\n" +
		"\nfunc ExampleHouse() {}\n" +
		"\nfunc BenchmarkHouse(b *testing.B) {}\n" +
		"\nfunc TestHelp(t *testing.T, more int) {}\n" +
		"\nfunc TestNone() {}\n" +
		"\nfunc TestOther(other *testing.B) {}\n" +
		"\nfunc TestAnother(other int) {}\n" +
		"\nfunc TestPointer(other *int) {}\n" +
		"\nfunc Testing(t *testing.T) {}\n" +
		"\nfunc helper(t *testing.T) {}\n" +
		"\ntype holder struct{}\n\n" +
		"func (holder) TestMethod(t *testing.T) {}\n"
	tests := ceiling() - oneTest
	seen := crowdIn(test, oneVow, tests, other)
	if strings.Join(seen, "|") != said(oneVow, tests) {
		test.Errorf("what is not a test was counted: %q", seen)
	}
}

// No vows are never crowded, with tests or without.
func TestAModuleWithoutVowsIsNotCrowded(test *testing.T) {
	for _, tests := range []int{noTests, ceiling()} {
		if seen := crowdIn(test, noVows, tests, ""); seen != nil {
			test.Errorf("no vows, %d tests: %v", tests, seen)
		}
	}
}

// Every vow with no test to stand over is crowded.
func TestAVowOverNoTestsIsCrowded(test *testing.T) {
	seen := crowdIn(test, oneVow, noTests, "")
	if strings.Join(seen, "|") != said(oneVow, noTests) {
		test.Errorf("one vow, no tests: %q", seen)
	}
}

// crowdSeen is what the crowd scroll said of a house with its files and
// its own tests below them: one vow, and the given count of tests.
func crowdSeen(test *testing.T, files map[string]string) []string {
	test.Helper()
	files[".editorconfig"] = "root = true\n\n[*.go]\nindent_style = tab\n"
	files["more.go"] = "package house\n\n// Open is open.\nvar Open = " +
		"map[string]string{} //ninja:vow well the wire reads it\n"
	_, facts := thrown(test, house(test, ".", files))
	var seen []string
	for _, each := range facts["hits"] {
		if each["scroll"] == "crowd" {
			seen = append(seen, each["seen"].(string))
		}
	}
	return seen
}

// Test in a file no test owns is not a test. go test reads only the
// files that end in _test.go.
func TestATestFunctionInAFileOfNoTestIsNotCounted(test *testing.T) {
	body := "package house\n\nimport \"testing\"\n\n" +
		"// TestNot is no test, being in a file go test never reads.\n" +
		"func TestNot(t *testing.T) {}\n"
	seen := crowdSeen(test, map[string]string{"notatest.go": body})
	if strings.Join(seen, "|") != crowded(noTests) {
		test.Errorf("a Test in a plain file was counted: %q", seen)
	}
}

// The same function in a _test file is a test that counts.
func TestATestFunctionInATestFileIsCounted(test *testing.T) {
	body := "package house\n\nimport \"testing\"\n\n" +
		"func TestYes(t *testing.T) {}\n"
	seen := crowdSeen(test, map[string]string{"yes_test.go": body})
	if strings.Join(seen, "|") != crowded(oneTest) {
		test.Errorf("a Test in a _test file: %q", seen)
	}
}

// Taking a testing.T without being named Test makes no test, whatever
// the first letter. Taking another package's T makes none either.
func TestAFunctionNotNamedTestIsNotCounted(test *testing.T) {
	body := "package house\n\nimport (\n\t\"testing\"\n\n" +
		"\t\"example.com/house/other\"\n)\n\n" +
		"var _ testing.TB\n\n" +
		"func Helper(t *testing.T) {}\n\n" +
		"func TestOther(t *other.T) {}\n"
	files := map[string]string{
		"help_test.go":   body,
		"other/other.go": "package other\n\n// T is not testing's.\ntype T struct{}\n",
	}
	seen := crowdSeen(test, files)
	if strings.Join(seen, "|") != crowded(noTests) {
		test.Errorf("what is not named Test was counted: %q", seen)
	}
}

// crowded is what the crowd scroll says of crowdSeen's house. The house
// carries one vow, over the given count of tests.
func crowded(tests int) string {
	return fmt.Sprintf("The module carries %d vows over %d tests, more "+
		"than one in 20. Bring your vows to the master.", oneVow, tests)
}
