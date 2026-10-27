package libninja_test

import (
	"strings"
	"testing"
)

// bareIn is what the number scroll saw in a function whose body is the
// given statements, with the given imports.
func bareIn(test *testing.T, statements string) string {
	test.Helper()
	body := "import (\n" + imports + ")\n\n" + used +
		"// Run does the statements.\nfunc Run(name string, data []byte) {\n" +
		statements + "}\n"
	return strings.Join(hitsIn(test, "number", body), "|")
}

// Calls that ask for a file mode say so by the place of the argument.
// The literal is let by only in that place.
func TestAModeLiteralInAModeArgumentIsKnown(test *testing.T) {
	for _, statement := range []string{
		"_ = os.WriteFile(name, data, 0644)",
		"_ = os.Mkdir(name, 0o755)",
		"_ = os.MkdirAll(name, 0o750)",
		"_, _ = os.OpenFile(name, os.O_RDONLY, 0600)",
		"_ = os.Chmod(name, 0o600)",
		"_ = ioutil.WriteFile(name, data, 0o644)",
		"var file *os.File\n\t_ = file.Chmod(0o600)",
		"var root *os.Root\n\t_ = root.Mkdir(name, 0o700)",
		"var root *os.Root\n\t_ = root.MkdirAll(name, 0o700)",
		"var root *os.Root\n\t_, _ = root.OpenFile(name, os.O_RDONLY, 0o600)",
		"var root *os.Root\n\t_ = root.WriteFile(name, data, 0o600)",
		"var root *os.Root\n\t_ = root.Chmod(name, 0o600)",
		"_ = fs.FileMode(0o644)",
		"_ = os.FileMode(0o644)",
		"_ = os.WriteFile(name, data, (0o644))",
	} {
		got := bareIn(test, "\t"+statement+"\n")
		if got != "" {
			test.Errorf("%s: %s", statement, got)
		}
	}
}

// The same literals are bare anywhere else. That means another argument
// of the same call, another call, a declaration, or a function of the
// module that only shares a name with a mode call.
func TestAModeLiteralAnywhereElseIsStillBare(test *testing.T) {
	for _, statement := range []string{
		"_, _ = os.OpenFile(name, 0644, os.ModePerm)",
		"_ = os.WriteFile(name, data, os.FileMode(0644)+0o1)",
		"_ = os.Mkdir(name, os.ModePerm|0o700)",
		"time.Sleep(0o644)",
		"mode := 0o644\n\t_ = mode",
		"_ = os.FileMode(0).String() + strconv.Itoa(0o644)",
	} {
		got := bareIn(test, "\t"+statement+"\n")
		if got == "" {
			test.Errorf("let by: %s", statement)
		}
	}
}

// The module's own call, named like one of os's, takes no exception.
// Only os and io/fs are known by their package.
func TestAModeCallOfTheModulesOwnIsNotKnown(test *testing.T) {
	body := "// WriteFile is the module's own.\n" +
		"func WriteFile(name string, data []byte, mode int) {}\n\n" +
		"// Run writes.\nfunc Run() { WriteFile(\"a\", nil, 0o644) }\n"
	got := strings.Join(hitsIn(test, "number", body), "|")
	if got != "0o644 is a bare number." {
		test.Errorf("the module's own WriteFile: %q", got)
	}
}

// imports is what every statement above may use. used declares a use of
// each, so none is unused.
const imports = "\t\"io/fs\"\n\t\"io/ioutil\"\n\t\"os\"\n\t\"regexp\"\n" +
	"\t\"strconv\"\n\t\"time\"\n"

const used = "var _, _, _, _, _, _ = fs.ModePerm, ioutil.Discard, " +
	"os.ModePerm, regexp.MustCompile, strconv.Itoa, time.Second\n\n"

// Another call's several results, handed on, leave fewer arguments
// written than places. There is no literal to let by and no error.
func TestAModeCallGivenAnotherCallsResultsHasNoLiteral(test *testing.T) {
	body := "import \"os\"\n\n" +
		"func three() (string, []byte, os.FileMode) { return \"\", nil, os.ModePerm }\n\n" +
		"// Run writes what three gives.\nfunc Run() { _ = os.WriteFile(three()) }\n"
	if seen := hitsIn(test, "number", body); seen != nil {
		test.Errorf("results handed on: %v", seen)
	}
}

// The -1 that asks a regexp for all its matches is the call's own
// word for all, in the place the count goes.
func TestMinusOneForAllMatchesIsKnown(test *testing.T) {
	for _, statement := range []string{
		"_ = regexp.MustCompile(name).FindAllString(name, -1)",
		"_ = regexp.MustCompile(name).FindAll(data, -1)",
		"_ = regexp.MustCompile(name).FindAllIndex(data, -1)",
		"_ = regexp.MustCompile(name).FindAllStringIndex(name, -1)",
		"_ = regexp.MustCompile(name).FindAllSubmatch(data, -1)",
		"_ = regexp.MustCompile(name).FindAllSubmatchIndex(data, -1)",
		"_ = regexp.MustCompile(name).FindAllStringSubmatch(name, -1)",
		"_ = regexp.MustCompile(name).FindAllStringSubmatchIndex(name, -1)",
		"_ = regexp.MustCompile(name).Split(name, -1)",
		"_ = regexp.MustCompile(name).FindAllString(name, (-1))",
	} {
		if got := bareIn(test, "\t"+statement+"\n"); got != "" {
			test.Errorf("%s: %s", statement, got)
		}
	}
}

// -1 is bare elsewhere, and so is any other count where all goes.
func TestMinusOneElsewhereIsStillBare(test *testing.T) {
	for _, statement := range []string{
		"_ = regexp.MustCompile(name).FindAllString(name, 3)",
		"_ = regexp.MustCompile(name).FindAllString(name, 1)",
		"_ = regexp.MustCompile(name).FindAllString(name, -2)",
		"_ = strconv.Itoa(-1)",
		"time.Sleep(-1)",
		"_ = regexp.MustCompile(name).NumSubexp() - 1",
	} {
		if got := bareIn(test, "\t"+statement+"\n"); got == "" {
			test.Errorf("let by: %s", statement)
		}
	}
}

// Its place in strconv's calls tells a precision or a bit size.
func TestAStrconvPrecisionAndBitSizeAreKnown(test *testing.T) {
	for _, statement := range []string{
		"_ = strconv.FormatFloat(value, 'f', -1, 64)",
		"_ = strconv.FormatFloat(value, 'f', 3, 32)",
		"_ = strconv.AppendFloat(data, value, 'f', -1, 64)",
		"_ = strconv.AppendFloat(data, value, 'g', 4, 32)",
		"_, _ = strconv.ParseFloat(name, 64)",
		"_, _ = strconv.ParseInt(name, base, 64)",
		"_, _ = strconv.ParseUint(name, base, 32)",
	} {
		got := bareIn(test, "\tvar value float64\n\tvar base int\n"+
			"\t_, _ = value, base\n\t"+statement+"\n")
		if got != "" {
			test.Errorf("%s: %s", statement, got)
		}
	}
}

// The same calls hold bare numbers in their other places. These are the
// value, the base and the count of another function. Bit sizes are bare
// outside strconv.
func TestAStrconvNumberInAnotherPlaceIsStillBare(test *testing.T) {
	for statement, want := range map[string]string{
		"_ = strconv.FormatFloat(2.5, 'f', -1, 64)":      "2.5 is a bare number.",
		"_, _ = strconv.ParseInt(name, 10, 64)":          "10 is a bare number.",
		"_ = strconv.Itoa(64)":                           "64 is a bare number.",
		"_ = strconv.AppendFloat(data, 2.5, 'g', 4, 32)": "2.5 is a bare number.",
		"time.Sleep(64)":                                 "64 is a bare number.",
	} {
		if got := bareIn(test, "\t"+statement+"\n"); got != want {
			test.Errorf("%s: %q, want %q", statement, got, want)
		}
	}
}

// Mkdir's mode is its second argument, and pair gives both. The call
// has one argument written, as many as its place is far along. It has
// no literal to let by and must not be read past its arguments.
func TestACallGivenResultsUpToItsPlaceHasNoLiteral(test *testing.T) {
	body := "import \"os\"\n\n" +
		"func pair() (string, os.FileMode) { return \"\", os.ModePerm }\n\n" +
		"// Run makes what pair names.\nfunc Run() { _ = os.Mkdir(pair()) }\n"
	if seen := hitsIn(test, "number", body); seen != nil {
		test.Errorf("results handed on: %v", seen)
	}
}

// Written out, the mode of the same call is still let by. Numbers
// beside it are not. The place is read, not the call.
func TestTheSameCallWrittenOutKeepsItsPlace(test *testing.T) {
	body := "import \"os\"\n\n" +
		"// Run makes a directory.\nfunc Run() { _ = os.Mkdir(\"a\", 0o700) }\n\n" +
		"// Again makes it with another mode.\nfunc Again(mode os.FileMode) { _ = os.Mkdir(\"a\", mode+5) }\n"
	got := strings.Join(hitsIn(test, "number", body), "|")
	if got != "5 is a bare number." {
		test.Errorf("written out: %q", got)
	}
}
