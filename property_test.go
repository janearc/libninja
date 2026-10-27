package libninja_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/janearc/puzzlemouth/numbers"

	"github.com/janearc/libninja"
)

// heldTo is a number the tests are held to, by its name in
// testdata/numbers.yaml. It can be read before any test runs.
func heldTo(name string) float64 {
	return heldIn("testdata/numbers.yaml", name)
}

// heldIn is a number by its name in a numbers file. When the file or the
// name is not there, the run stops and says which.
func heldIn(file, name string) float64 {
	raw, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	held, err := numbers.Read(file, raw)
	if err != nil {
		panic(err)
	}
	number, err := held.Get(name)
	if err != nil {
		panic(err)
	}
	return number.Value()
}

// askedFlags is the flags a run of the tests was given by name.
type askedFlags struct{ names map[string]bool }

// note keeps one flag's name.
func (asked askedFlags) note(given *flag.Flag) { asked.names[given.Name] = true }

// TestMain runs each property on as many drawn modules as
// testdata/numbers.yaml picks, unless the run asks for more. Each module
// is written and loaded, so rapid's own count would make every mutants
// run far longer. -rapid.checks or RAPID_CHECKS asks for another count.
func TestMain(suite *testing.M) {
	flag.Parse()
	asked := askedFlags{map[string]bool{}}
	flag.Visit(asked.note)
	if !asked.names["rapid.checks"] && os.Getenv("RAPID_CHECKS") == "" {
		checks := int(heldTo("rapid-checks"))
		_ = flag.Set("rapid.checks", strconv.Itoa(checks))
	}
	os.Exit(suite.Run())
}

// drawnModule is a module of drawn files, written into a folder of its
// own under base.
type drawnModule struct{ base string }

// write puts the files in a fresh folder with a go.mod, and is the folder.
func (module drawnModule) write(t *rapid.T, files map[string]string) string {
	dir, err := os.MkdirTemp(module.base, "drawn")
	if err != nil {
		t.Fatal(err)
	}
	files["go.mod"] = "module example.com/drawn\n\ngo 1.26\n"
	mode := os.FileMode(heldTo("open-file"))
	for name, text := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(text), mode); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// lined is source written a line at a time, each line's number known as
// it is added.
type lined struct{ lines []string }

// add puts a line at the end and is its number, counted from one.
func (source *lined) add(text string) int {
	source.lines = append(source.lines, text)
	return len(source.lines)
}

// text is the source whole.
func (source *lined) text() string {
	return strings.Join(source.lines, "\n") + "\n"
}

// seenAt is one scroll's findings in one part of the facts, in order.
// Each is its file and line and what was seen.
func seenAt(facts map[string][]map[string]any, part, scroll string) []string {
	found := seenAs(facts, part, scroll)
	slices.Sort(found)
	return found
}

// spokenAs stops a property unless the master says every hit as the
// facts list it, and calls the module calm exactly when it is clean.
func spokenAs(t *rapid.T, report libninja.Report,
	facts map[string][]map[string]any) {
	var spoken bytes.Buffer
	if err := report.Speak(&spoken); err != nil {
		t.Fatal(err)
	}
	for _, each := range facts["hits"] {
		said := fmt.Sprintf("    %s:%v  %s\n", each["file"], each["line"],
			each["seen"])
		if !strings.Contains(spoken.String(), said) {
			t.Fatalf("the master never says %q", said)
		}
	}
	calm := strings.HasSuffix(spoken.String(), "Now, go practice.\n")
	if report.Clean() != calm {
		t.Fatalf("clean is %v, and the master calm is %v", report.Clean(),
			calm)
	}
}

// ends is a function's first line and its closing brace. Any function
// written over lines has those two, as go has defined it.
const ends = 2

// atLeastOne is the fewest things a drawn module holds, defined so: a
// module of nothing tests nothing.
const atLeastOne = 1

// stanceCase draws functions of many lengths into modules under base.
type stanceCase struct{ module drawnModule }

// check holds that a function is a stance exactly when it runs past
// the longest stance.
func (each stanceCase) check(t *rapid.T) {
	longest := int(heldTo("drawn-longest-function"))
	most := int(heldTo("drawn-most-functions"))
	stance := int(heldIn("picked-numbers.yml", "longest-stance"))
	near := int(heldTo("drawn-near-stance"))
	// half the lengths are drawn from anywhere and half from beside the
	// stance. So a line either way of it is drawn in every run
	anywhere := rapid.IntRange(ends, longest)
	beside := rapid.IntRange(stance-near, stance+near)
	lengths := rapid.SliceOfN(rapid.OneOf(anywhere, beside), atLeastOne,
		most).Draw(t, "lengths")
	source := lined{[]string{"package drawn"}}
	var want []string
	for index, length := range lengths {
		name := "step" + string(rune('a'+index))
		source.add("")
		first := source.add("func " + name + "() {")
		for range length - ends {
			source.add("\t_ = 0")
		}
		source.add("}")
		if length > stance {
			want = append(want, fmt.Sprintf("drawn.go:%d %s is %d "+
				"lines.", first, name, length))
		}
	}
	slices.Sort(want)
	dir := each.module.write(t, map[string]string{"drawn.go": source.text()})
	report, facts := reported(t, dir)
	if got := seenAt(facts, "hits", "stance"); !slices.Equal(got, want) {
		t.Fatalf("stances %q, not %q", got, want)
	}
	spokenAs(t, report, facts)
}

// Whatever lengths are drawn, a function is a stance exactly when it
// runs past the longest stance. Its lines are counted from its first to
// its brace.
func TestAFunctionPastAScreenIsAStance(test *testing.T) {
	test.Setenv("GOWORK", "off")
	rapid.Check(test, stanceCase{drawnModule{test.TempDir()}}.check)
}

// spot is where a drawn one-letter name is made. It is a local, a test's
// *testing.T or a test file's parameter of another type.
type spot struct {
	letter rune
	place  string
}

// drawSpot draws one name's letter and where it is made.
func drawSpot(t *rapid.T) spot {
	places := []string{"local", "testing", "int"}
	return spot{rune(rapid.IntRange('a', 'z').Draw(t, "letter")),
		rapid.SampledFrom(places).Draw(t, "place")}
}

// spotted writes a drawn name where its spot says. It is where the name
// stands, and whether the rule lets it stand there alone.
func spotted(each spot, index int, code, tests *lined) (string, bool) {
	name := string(each.letter)
	alone := !strings.ContainsRune("ikxyz", each.letter)
	switch each.place {
	case "local":
		code.add("")
		line := code.add(fmt.Sprintf("func use%d() { %s := 0; _ = %s }",
			index, name, name))
		return fmt.Sprintf("drawn.go:%d", line), alone
	case "testing":
		tests.add("")
		line := tests.add(fmt.Sprintf(
			"func TestUse%d(%s *testing.T) { _ = %s }", index, name, name))
		return fmt.Sprintf("drawn_test.go:%d", line), alone && name != "t"
	}
	tests.add("")
	line := tests.add(fmt.Sprintf("func use%d(%s int) int { return %s }",
		index, name, name))
	return fmt.Sprintf("drawn_test.go:%d", line), alone
}

// nameCase draws one-letter names into modules under base.
type nameCase struct{ module drawnModule }

// check holds that a one-letter name stands alone unless it is a test's
// t or one of the letters ikxyz.
func (each nameCase) check(t *rapid.T) {
	most := int(heldTo("drawn-most-names"))
	spots := rapid.SliceOfN(rapid.Custom(drawSpot), atLeastOne, most).
		Draw(t, "spots")
	code := lined{[]string{"package drawn"}}
	tests := lined{[]string{"package drawn_test", "", `import "testing"`,
		"", "func TestDrawn(test *testing.T) {}"}}
	var want []string
	for index, drawn := range spots {
		at, alone := spotted(drawn, index, &code, &tests)
		if alone {
			want = append(want, at+" "+string(drawn.letter)+
				" stands alone.")
		}
	}
	slices.Sort(want)
	dir := each.module.write(t, map[string]string{"drawn.go": code.text(),
		"drawn_test.go": tests.text()})
	report, facts := reported(t, dir)
	if got := seenAt(facts, "hits", "name"); !slices.Equal(got, want) {
		t.Fatalf("names %q, not %q", got, want)
	}
	spokenAs(t, report, facts)
}

// Every one-letter name stands alone unless it is a test's t or one of
// the letters ikxyz.
func TestAOneLetterNameStandsAloneByTheRule(test *testing.T) {
	test.Setenv("GOWORK", "off")
	rapid.Check(test, nameCase{drawnModule{test.TempDir()}}.check)
}

// placed is what a drawn vow leaves. It holds the well's hit when the vow
// does not cover it and the well as vowed when it does. It also holds
// the vow's own hit.
type placed struct{ hit, kept, idle []string }

// vowPlaced writes an open well and a vow placed by one way of drawing.
func vowPlaced(source *lined, index int, placing string) placed {
	name := fmt.Sprintf("Open%d", index)
	opened := "var " + name + " = 1"
	seen := name + " can be changed by anyone."
	idleSaid := "%d The vow to break %s covers nothing here. Put it down."
	source.add("")
	var idle []string
	switch placing {
	case "beside":
		at := source.add(opened + " //ninja:vow well drawn right beside")
		return placed{kept: []string{fmt.Sprintf("drawn.go:%d %s", at,
			seen)}}
	case "above":
		source.add("//ninja:vow well drawn right above")
		at := source.add(opened)
		return placed{kept: []string{fmt.Sprintf("drawn.go:%d %s", at,
			seen)}}
	case "bare":
		at := source.add("//ninja:vow well")
		idle = []string{fmt.Sprintf("drawn.go:%d The vow to break well "+
			"gives no reason.", at)}
	case "far":
		at := source.add("//ninja:vow well drawn too far")
		source.add("// between")
		idle = []string{fmt.Sprintf("drawn.go:"+idleSaid, at, "well")}
	case "other":
		at := source.add("//ninja:vow gate drawn for another")
		idle = []string{fmt.Sprintf("drawn.go:"+idleSaid, at, "gate")}
	case "after":
		// at the end of a line of code, a vow covers that line
		at := source.add(fmt.Sprintf("var pad%d = 0 //ninja:vow well "+
			"drawn after code", index))
		idle = []string{fmt.Sprintf("drawn.go:"+idleSaid, at, "well")}
	}
	at := source.add(opened)
	return placed{hit: []string{fmt.Sprintf("drawn.go:%d %s", at, seen)},
		idle: idle}
}

// vowCase draws vows in every placing into modules under base.
type vowCase struct{ module drawnModule }

// check holds that a vow covers its scroll's hit only with a reason, on
// the hit's line or alone on the line above.
func (each vowCase) check(t *rapid.T) {
	most := int(heldTo("drawn-most-vows"))
	placings := rapid.SliceOfN(rapid.SampledFrom([]string{"none",
		"beside", "above", "bare", "far", "other", "after"}), atLeastOne,
		most).Draw(t, "placings")
	source := lined{[]string{"package drawn"}}
	var all placed
	for index, placing := range placings {
		one := vowPlaced(&source, index, placing)
		all.hit = append(all.hit, one.hit...)
		all.kept = append(all.kept, one.kept...)
		all.idle = append(all.idle, one.idle...)
	}
	for _, sorted := range [][]string{all.hit, all.kept, all.idle} {
		slices.Sort(sorted)
	}
	dir := each.module.write(t, map[string]string{"drawn.go": source.text()})
	report, facts := reported(t, dir)
	for _, check := range []struct {
		part, scroll string
		want         []string
	}{{"hits", "well", all.hit}, {"vows", "well", all.kept},
		{"hits", "vow", all.idle}} {
		got := seenAt(facts, check.part, check.scroll)
		if !slices.Equal(got, check.want) {
			t.Fatalf("%s of %s: %q, not %q", check.part, check.scroll,
				got, check.want)
		}
	}
	spokenAs(t, report, facts)
}

// Each vow covers its scroll's hit when it gives a reason and stands on
// the hit's line or alone on the line above. Anywhere else the hit
// stands, and the vow is a hit of its own. So is a vow with no reason.
func TestAVowCoversWhatItStandsBesideByTheRule(test *testing.T) {
	test.Setenv("GOWORK", "off")
	rapid.Check(test, vowCase{drawnModule{test.TempDir()}}.check)
}
