package libninja_test

import (
	"os"
	"strings"
	"testing"
)

func TestTheFixturesModesAreReadAsWritten(test *testing.T) {
	dir := os.FileMode(fixture(test, "open-dir")).String()
	file := os.FileMode(fixture(test, "open-file")).String()
	if dir != "-rwxr-xr-x" || file != "-rw-r--r--" {
		test.Errorf("0755 read as %s, 0644 as %s", dir, file)
	}
}

func TestADefinedNumberMayLiveInCode(test *testing.T) {
	body := "// first is where lines start, as go/token has Defined it.\n" +
		"const first = 1\n\n" +
		"// Up is a step up.\nfunc Up(line int) int {\n" +
		"\tstep := 1 // one line, defined so\n\treturn line - step + first\n}\n"
	if seen := hitsIn(test, "source", body); seen != nil {
		test.Errorf("defined numbers: %v", seen)
	}
}

func TestAPickedNumberInCodeIsAHit(test *testing.T) {
	body := "// widest is how wide a line may run, picked by eye.\n" +
		"const widest = 80\n\n" +
		"// Wide is whether a line is too wide.\n" +
		"func Wide(width int) bool {\n" +
		"\tslack := 2 // a little room\n\treturn width > widest+slack\n}\n"
	const hit = "A number named here says why, but not that it is " +
		"defined; measured or picked, it lives in data."
	seen := strings.Join(hitsIn(test, "source", body), "|")
	if seen != hit+"|"+hit {
		test.Errorf("picked numbers in code: %v", seen)
	}
}

func TestANameWithoutANumberIsNotTheSourceScrolls(test *testing.T) {
	body := "// greeting is what is said first.\n" +
		"const greeting = \"hello\"\n"
	if seen := hitsIn(test, "source", body); seen != nil {
		test.Errorf("a name holding no number: %v", seen)
	}
}
