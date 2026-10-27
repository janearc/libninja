package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Throw says what the scrolls hit; in a kept module, nothing.
func Example_throw() {
	os.Setenv("GOWORK", "off")
	var out bytes.Buffer
	code := tool().Run([]string{"throw", "testdata/kept"}, &out, &out)
	fmt.Print(out.String())
	fmt.Println(code)
	// Output:
	// shuriken reads example.com/kept: 0 hits in 0 packages, 0 vows.
	//
	// A ninja must always be calm. Now, go practice.
	// 0
}

// Facts writes the same as json.
func Example_facts() {
	os.Setenv("GOWORK", "off")
	var out bytes.Buffer
	code := tool().Run([]string{"facts", "testdata/kept"}, &out, &out)
	fmt.Print(out.String())
	fmt.Println(code)
	// Output:
	// {
	//   "hits": [],
	//   "module": "example.com/kept",
	//   "vows": []
	// }
	// 0
}

func TestEveryVerbHasItsExample(test *testing.T) {
	if err := tool().Examples(os.DirFS(".")); err != nil {
		test.Error(err)
	}
}

// verbFailed is pzm's exit code for a verb that ran and failed, as pzm
// has defined it.
const verbFailed = 1

func TestAHitFailsTheThrowAndSaysSo(test *testing.T) {
	test.Setenv("GOWORK", "off")
	var out, errs bytes.Buffer
	code := tool().Run([]string{"throw", "testdata/struck"}, &out, &errs)
	if code != verbFailed ||
		!strings.Contains(out.String(), "q stands alone") ||
		!strings.Contains(errs.String(), "broken without a vow") {
		test.Errorf("a struck module: %d, %q, %q", code, out.String(),
			errs.String())
	}
}

func TestAModuleThatCannotBeReadFailsTheThrow(test *testing.T) {
	test.Setenv("GOWORK", "off")
	var out, errs bytes.Buffer
	code := tool().Run([]string{"throw", test.TempDir()}, &out, &errs)
	if code != verbFailed ||
		!strings.Contains(errs.String(), "in no module") {
		test.Errorf("no module: %d, %q", code, errs.String())
	}
}

// refusing is a writer that takes nothing.
type refusing struct{}

// Write refuses every byte, and writes none.
func (refusing) Write([]byte) (written int, err error) {
	return written, errors.New("full")
}

func TestAReportThatCannotBeWrittenFailsTheThrow(test *testing.T) {
	test.Setenv("GOWORK", "off")
	var errs bytes.Buffer
	code := tool().Run([]string{"throw", "testdata/kept"}, refusing{},
		&errs)
	if code != verbFailed || !strings.Contains(errs.String(), "full") {
		test.Errorf("a report that cannot be written: %d, %q", code,
			errs.String())
	}
}

// heldExit is main's way out, held by a test: the code it was handed.
type heldExit struct{ code *int }

// end keeps the code instead of ending the process.
func (held heldExit) end(code int) { *held.code = code }

// restoring puts main's world back as a test found it.
type restoring struct {
	exit func(int)
	args []string
	out  *os.File
}

// restore puts exit, the arguments and standard output back.
func (was restoring) restore() {
	exit, os.Args, os.Stdout = was.exit, was.args, was.out
}

func TestMainExitsWithTheToolsCode(test *testing.T) {
	place := filepath.Join(test.TempDir(), "said")
	said, err := os.Create(place)
	if err != nil {
		test.Fatal(err)
	}
	was := restoring{exit, os.Args, os.Stdout}
	test.Cleanup(was.restore)
	code := verbFailed
	exit = heldExit{&code}.end
	os.Args = []string{"shuriken", "--age"}
	os.Stdout = said
	main()
	was.restore()
	text, _ := os.ReadFile(said.Name())
	// main hands the tool what follows its own name: the answer to
	// --age, exactly as the tool gives it when asked
	var age bytes.Buffer
	answered := tool().Run([]string{"--age"}, &age, &age)
	if code != answered || string(text) != age.String() ||
		!strings.HasPrefix(age.String(), "shuriken") {
		test.Errorf("main exited %d, saying %q; --age is %d, %q", code,
			text, answered, age.String())
	}
}

func TestAGoPatternNamesTheSameModule(test *testing.T) {
	for given, want := range map[string]string{"": ".", "./...": ".",
		"...": ".", "testdata/kept/...": "testdata/kept",
		"testdata/kept": "testdata/kept"} {
		args := strings.Fields(given)
		if got := dirOf(args); got != want {
			test.Errorf("%q names %q, want %q", given, got, want)
		}
	}
}

// TestMain points HOME at an empty directory, so pzm finds no
// ~/.config/shuriken/config of the person running the tests.
func TestMain(suite *testing.M) {
	home, err := os.MkdirTemp("", "shuriken-home")
	if err != nil {
		panic(err)
	}
	os.Setenv("HOME", home)
	code := suite.Run()
	os.RemoveAll(home)
	os.Exit(code)
}
