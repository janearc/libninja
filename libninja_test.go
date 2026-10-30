package libninja_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/janearc/libninja"
)

// thrown is a module's report through the real interface. It is loaded
// and thrown and read back from its JSON. Each fixture is a module of
// its own, read without the workspace. libninja itself is read with the
// workspace, as it is built.
func thrown(test *testing.T, dir string) (libninja.Report,
	map[string][]map[string]any) {
	test.Helper()
	if dir != self {
		test.Setenv("GOWORK", "off")
	}
	return reported(test, dir)
}

// failer is what a test and a rapid property both have, to stop on.
type failer interface{ Fatal(args ...any) }

// reported is a module's report and its facts, or the test stops.
func reported(fails failer, dir string) (libninja.Report,
	map[string][]map[string]any) {
	module, err := libninja.Load(dir)
	if err != nil {
		fails.Fatal(err)
	}
	report := module.Throw()
	return report, factsOf(fails, report)
}

// factsOf is a report's hits and vows as its JSON gives them back. When
// the JSON cannot be read, the test stops.
func factsOf(fails failer,
	report libninja.Report) map[string][]map[string]any {
	var written bytes.Buffer
	if err := report.JSON(&written); err != nil {
		fails.Fatal(err)
	}
	var parts map[string]json.RawMessage
	if err := json.Unmarshal(written.Bytes(), &parts); err != nil {
		fails.Fatal(err)
	}
	facts := map[string][]map[string]any{}
	for _, part := range []string{"hits", "vows"} {
		var found []map[string]any
		if err := json.Unmarshal(parts[part], &found); err != nil {
			fails.Fatal(err)
		}
		facts[part] = found
	}
	return facts
}

// holds is whether a finding of the scroll in the file saw something
// that contains said.
func holds(found []map[string]any, scroll, file, said string) bool {
	for _, each := range found {
		if each["scroll"] == scroll && each["file"] == file &&
			strings.Contains(each["seen"].(string), said) {
			return true
		}
	}
	return false
}

// expected is a hit a test looks for: by which scroll, in which file and
// saying what.
type expected struct {
	scroll, file, said string
}

func TestTheDojoBreaksEveryScroll(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	for _, want := range []expected{
		{"gate", "temple/temple.go", "Sword.Edge stands open; touched " +
			"from village."},
		{"stranger", "village/village.go", "Builds temple.Sword by literal."},
		{"well", "temple/temple.go", "Hidden can be changed by anyone."},
		{"waste", "temple/temple.go", "Idle is exported, and nothing " +
			"outside temple uses it."},
		{"practice", "temple/temple.go", "Practice is used from outside " +
			"only by tests."},
		{"dojo", "temple/inside_test.go", "package temple, inside the wall"},
		{"blade", "temple/count.go", "The body of count is forged again " +
			"at temple/count.go:21."},
		{"blade", "village/first.go", "temple.Hone(temple.Forge(...)) is " +
			"forged here"},
		{"stance", "temple/long.go", "long is 54 lines."},
		{"stance", "temple/screens.go", "fiftyOne is 51 lines."},
		{"name", "temple/temple.go", "q stands alone."},
		{"name", "temple/inside_test.go", "t stands alone."},
		{"age", "cmd/gate/main.go", "declares no build and built"},
		{"age", "cmd/gate/main.go", "does not answer --age"},
		{"buried", "village/bury.go", "\"/tmp/dojo.txt\" is one machine's " +
			"path, handed to os.WriteFile."},
		{"literal", "village/nameless.go", "written here with no name"},
		{"index", "village/place.go", "[0] says a place"},
		{"index", "village/place.go", "[1] says a place"},
		{"vow", "village/village.go", "The vow to break stranger gives no " +
			"reason."},
		{"feint", "temple/feint.go", "feintThing{} is an empty " +
			"value handed back with an error."},
		{"feint", "temple/feint.go", "made is a feintThing " +
			"declared and never set"},
		{"feint", "temple/feint.go", "feintThing{...} has " +
			"fields set"},
	} {
		if !holds(facts["hits"], want.scroll, want.file, want.said) {
			test.Errorf("no %s in %s saying %q", want.scroll, want.file,
				want.said)
		}
	}
}

func TestAFunctionOfOneScreenIsNoStance(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if holds(facts["hits"], "stance", "temple/screens.go", "fifty is") {
		test.Error("fifty is exactly the longest stance, and was " +
			"called one past it")
	}
}

func TestAMainsTestMayStandInside(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if holds(facts["hits"], "dojo", "cmd/gate/main_test.go", "") {
		test.Error("a main's test, which can stand nowhere else, was " +
			"sent outside the wall")
	}
}

func TestAPromiseKeptStaysExported(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if holds(facts["hits"], "waste", "temple/temple.go", "Sword.String") {
		test.Error("String satisfies fmt.Stringer, and was called waste")
	}
}

func TestAPromiseKeptThroughEmbeddingStaysExported(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if holds(facts["hits"], "waste", "temple/blade.go", "steel.Cut") {
		test.Error("steel.Cut answers Cutter through Blade, and was " +
			"called waste")
	}
}

func TestAnExternalTestBuildsItsOwnAtHome(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if holds(facts["hits"], "stranger", "village/own_test.go", "") {
		test.Error("an external test building its own type was called " +
			"a stranger")
	}
}

func TestADirectoryInsideMeansTheWholeModule(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo/temple")
	if holds(facts["hits"], "waste", "temple/temple.go", "Forge") {
		test.Error("read from temple alone, Forge, which village calls, " +
			"was called waste")
	}
	if !holds(facts["hits"], "stranger", "village/village.go", "") {
		test.Error("read from temple, the village was not read at all")
	}
}

func TestABorrowedNameIsStillAnothersHouse(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "stranger", "village/alias.go",
		"Builds temple.Sword by literal.") {
		test.Error("borrowed{...} builds temple.Sword through an alias, " +
			"and the stranger was not seen")
	}
}

func TestATypeCarriedByAUsedFunctionIsInUse(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	for _, scroll := range []string{"waste", "practice"} {
		if holds(facts["hits"], scroll, "temple/carried.go", "Belt is") {
			test.Errorf("the village calls Tie, which gives a Belt, and "+
				"Belt was called %s", scroll)
		}
	}
	if !holds(facts["hits"], "practice", "temple/carried.go",
		"Sash is used from outside only by tests.") {
		test.Error("only a test calls Wrap, which gives a Sash, and Sash " +
			"was not called practice")
	}
}

// nameLines is the lines a file's one-letter names stand alone on.
func nameLines(facts map[string][]map[string]any, file string) []int {
	var lines []int
	for _, each := range facts["hits"] {
		if each["scroll"] == "name" && each["file"] == file {
			lines = append(lines, int(each["line"].(float64)))
		}
	}
	return lines
}

func TestATestsTMayStandAlone(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	insideAgain := int(fixture(test, "inside-again"))
	drawnAsR := int(fixture(test, "drawn-as-r"))
	inside := nameLines(facts, "temple/inside_test.go")
	if !slices.Equal(inside, []int{insideAgain}) {
		test.Errorf("in inside_test.go only TestInsideAgain's int t, "+
			"at %d, stands alone, not TestInside's *testing.T: %v",
			insideAgain, inside)
	}
	outside := nameLines(facts, "temple/outside_test.go")
	if !slices.Equal(outside, []int{drawnAsR}) {
		test.Errorf("in outside_test.go only the r at %d stands alone, "+
			"not a test's t or a property's: %v", drawnAsR, outside)
	}
	if !holds(facts["hits"], "name", "temple/drawn.go", "t stands alone.") {
		test.Error("a rapid t outside a test file was let stand alone")
	}
}

// decoyed is a hit the decoy scroll makes in the dojo: where, and what
// the master says.
type decoyed struct{ file, seen string }

func TestAMockIsUsedOrItIsADecoy(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	for _, want := range []decoyed{
		{"mocks/bells.go", "MockGong is a mock of Gong that no test makes."},
		{"temple/outside_test.go",
			"fakeGong stands in for Gong, which has a mock: MockGong."},
		{"temple/outside_test.go",
			"gongSound stands in for Gong, which has a mock: MockGong."},
		{"temple/bells.go", "Tests swap tide, first at inside_test.go:19."},
		{"temple/bells.go",
			"Tests swap Weather, first at outside_test.go:72."},
		{"temple/bells.go", "Tests swap Yard, first at outside_test.go:73."},
		{"temple/bells.go",
			"Tests swap Shelves, first at outside_test.go:74."},
		{"temple/bells.go",
			"Tests swap Counter, first at outside_test.go:75."},
	} {
		if !holds(facts["hits"], "decoy", want.file, want.seen) {
			test.Errorf("the decoy at %s was not seen: %s", want.file,
				want.seen)
		}
	}
	for _, each := range facts["hits"] {
		seen, _ := each["seen"].(string)
		for _, unseen := range []string{"MockBell", "worn", "rounds",
			"tolled", "Plot", "Stray", "Lone"} {
			if each["scroll"] == "decoy" && strings.Contains(seen, unseen) {
				test.Errorf("%s is no decoy, and was called one: %s",
					unseen, seen)
			}
		}
	}
}

func TestAMainThatAnswersItsAgeIsLetBe(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	for _, each := range facts["hits"] {
		if each["scroll"] == "age" && each["file"] == "cmd/asked/main.go" {
			test.Errorf("cmd/asked declares its stamps and answers "+
				"--age, and was told: %s", each["seen"])
		}
	}
}

func TestAModuleThatCannotBeReadSaysWhy(test *testing.T) {
	test.Setenv("GOWORK", "off")
	// notcompiling has two errors: go list's and the type checker's. go
	// list's comes first. importsnothing has one and is refused for it.
	for dir, says := range map[string]string{
		test.TempDir():              "is in no module",
		"testdata/modules/badgomod": "errors parsing go.mod",
		"testdata/modules/nopackages": "libninja: no module under " +
			"/",
		"testdata/modules/notcompiling": "libninja: " +
			"example.com/notcompiling does not compile: -: # " +
			"example.com/notcompiling\n",
		"testdata/modules/importsnothing": "libninja: " +
			"example.com/importsnothing does not compile: ",
	} {
		_, err := libninja.Load(dir)
		if err == nil || !strings.Contains(err.Error(), says) {
			test.Errorf("%s: %v, not %q", dir, err, says)
		}
	}
}

func TestADirectoryThatIsGoneSaysWhy(test *testing.T) {
	gone := test.TempDir()
	test.Chdir(gone)
	if err := os.Remove(gone); err != nil {
		test.Fatal(err)
	}
	_, err := libninja.Load(self)
	if err == nil || err.Error() != "libninja: . is in no module" {
		test.Errorf("a directory that is gone: %v", err)
	}
}

func TestACleanModuleIsCalledCalm(test *testing.T) {
	report, _ := thrown(test, "testdata/modules/clean")
	var spoken bytes.Buffer
	if err := report.Speak(&spoken); err != nil {
		test.Fatal(err)
	}
	if !report.Clean() ||
		!strings.Contains(spoken.String(), "Now, go practice.") ||
		strings.Contains(spoken.String(), "the vows") {
		test.Errorf("a clean module, clean %v, said:\n%s", report.Clean(),
			spoken.String())
	}
}

func TestAVowNamingNoScrollIsAHit(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "vow", "temple/universe.go",
		`The vow names "nosuch", which is no scroll.`) {
		test.Error("a vow naming nosuch was not called out")
	}
	if !holds(facts["hits"], "vow", "temple/bare.go",
		`The vow names "", which is no scroll.`) {
		test.Error("a bare vow was not called out")
	}
}

func TestAVowWhereTheSourceCannotBeReadStandsAlone(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "vow", "/dojo-elsewhere/generated.y",
		"The vow to break stranger covers nothing") {
		test.Error("a vow in a file a //line directive moves out of " +
			"reach was not read as standing alone, said whole")
	}
	if !holds(facts["hits"], "vow", "temple/nowhere.go",
		"The vow to break stranger covers nothing") {
		test.Error("a vow a //line directive puts past the end of its " +
			"line was not read as standing alone")
	}
}

// The vows standing.go writes are moved by //line to notes.txt. Each is
// judged where it is written. Three stand alone above their vars. One
// comes after its var's code. So all four vars are vowed.
func TestAVowIsJudgedWhereItIsWritten(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	vowed := seenAs(facts, "vows", "well")
	for _, each := range []string{"temple/notes.txt:2 DirectedOne",
		"temple/notes.txt:4 DirectedTwo",
		"temple/notes.txt:6 DirectedThree",
		"temple/notes.txt:7 DirectedFour"} {
		if !slices.Contains(vowed, each+" can be changed by anyone.") {
			test.Errorf("%s was not vowed: %v", each, vowed)
		}
	}
}

// seenAs is one scroll's findings in one part of the facts, each as its
// file and line and what was seen.
func seenAs(facts map[string][]map[string]any, part, scroll string) []string {
	var found []string
	for _, each := range facts[part] {
		if each["scroll"] == scroll {
			found = append(found, fmt.Sprintf("%s:%v %s", each["file"],
				each["line"], each["seen"]))
		}
	}
	return found
}

func TestAGenericPromiseKeptStaysExported(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if holds(facts["hits"], "waste", "temple/generic.go", "String") {
		test.Error("Box's String keeps fmt.Stringer's promise, and was " +
			"called waste")
	}
}

func TestAnExternalTestIsStillOnlyATest(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "practice", "temple/generic.go",
		"Drill is used from outside only by tests.") {
		test.Error("only temple's external test calls Drill, and it " +
			"was not called practice")
	}
}

func TestGeneratedCodeIsTheToolsToAnswerFor(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	for _, scroll := range []string{"gate", "well", "waste"} {
		if holds(facts["hits"], scroll, "temple/made.go", "") {
			test.Errorf("a generated file broke %s", scroll)
		}
	}
}

func TestAVariableIsTheWellsAlone(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "well", "temple/blade.go", "Spare") {
		test.Error("Spare is an open well, and was not called one")
	}
	if holds(facts["hits"], "waste", "temple/blade.go", "Spare") {
		test.Error("Spare is a well, and was counted as waste as well")
	}
}

func TestAVowAtTheEndOfALineCoversOnlyThatLine(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if holds(facts["hits"], "gate", "temple/vows.go", "Rope.Knot") {
		test.Error("Rope.Knot is vowed at the end of its line, and was " +
			"still a hit")
	}
	if !holds(facts["hits"], "gate", "temple/vows.go", "Rope.Loose") {
		test.Error("Rope.Loose is open on the line after a trailing " +
			"vow, and the vow covered it")
	}
}

func TestAVowThatCoversNothingIsAHit(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "vow", "temple/vows.go", "The vow to break "+
		"name covers nothing here.") {
		test.Error("a vow over calm, which breaks no scroll, was not " +
			"called out")
	}
}

func TestAVowWithAReasonCoversItsFinding(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if holds(facts["hits"], "gate", "temple/temple.go", "Wire.Name") {
		test.Error("Wire.Name is vowed, and was still a hit")
	}
	if !holds(facts["vows"], "gate", "temple/temple.go", "Wire.Name") {
		test.Error("Wire.Name's vow is not among the vows")
	}
}

func TestTheMasterAddressesEachStudent(test *testing.T) {
	report, _ := thrown(test, "testdata/dojo")
	var spoken bytes.Buffer
	if err := report.Speak(&spoken); err != nil {
		test.Fatal(err)
	}
	for _, said := range []string{"village, stay here. Sit down.",
		"temple, stay here. Sit down.", `"Correct your stance."`,
		"Understand?", "Keep in mind, you must always be a good ninja."} {
		if !strings.Contains(spoken.String(), said) {
			test.Errorf("the master never says %q", said)
		}
	}
	if report.Clean() {
		test.Error("the dojo breaks scrolls, and was called clean")
	}
}

// self is libninja's own directory, where its tests run.
const self = "."

func TestLibninjaKeepsItsOwnScrolls(test *testing.T) {
	report, facts := thrown(test, self)
	if report.Clean() {
		return
	}
	for _, each := range facts["hits"] {
		test.Errorf("%s:%v %s: %s", each["file"], each["line"],
			each["scroll"], each["seen"])
	}
}

// inDir is whether a practice or a waste hit is in a file under a
// directory.
func inDir(facts map[string][]map[string]any, dir string) bool {
	for _, each := range facts["hits"] {
		if (each["scroll"] == "practice" || each["scroll"] == "waste") &&
			strings.HasPrefix(each["file"].(string), dir+"/") {
			return true
		}
	}
	return false
}

func TestAPackageMadeForTestsIsNoPractice(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if inDir(facts, "drumstest") {
		test.Error("drumstest is named for tests and only a test " +
			"imports it, and was still called out")
	}
}

func TestAPackageMadeForTestsRealCodeImportsIsLeaked(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "practice", "village/leak.go",
		"fmttest is made for tests, and real code imports it.") {
		test.Error("village/leak.go imports fmttest, and the leak " +
			"was not seen")
	}
	if inDir(facts, "fmttest") {
		test.Error("fmttest's names are used, and none of them was " +
			"the fault")
	}
}

func TestAPackageMadeForTestsNoTestImportsIsDead(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "practice", "ostest/ostest.go",
		"ostest is made for tests, and no test imports it.") {
		test.Error("no test imports ostest, and it was not called dead")
	}
	if holds(facts["hits"], "waste", "ostest/ostest.go", "Rust") {
		test.Error("ostest is called dead once, not once more for " +
			"each of its names")
	}
}

func TestAnOrdinaryPackageOnlyTestsUseIsStillPractice(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	if !holds(facts["hits"], "practice", "drums/drums.go",
		"Roll is used from outside only by tests.") {
		test.Error("drums is no package for tests by its name, and " +
			"Roll was not called practice")
	}
}

func TestAPackageNamedForAnotherAndTestIsMadeForTests(test *testing.T) {
	_, facts := thrown(test, "testdata/modules/support")
	for _, dir := range []string{"gittest", "httptest"} {
		if inDir(facts, dir) {
			test.Errorf("%s is the name of a package and test, and "+
				"only a test imports it, and was called out", dir)
		}
	}
}

func TestAPackageThatOnlyEndsInTestIsNoSupport(test *testing.T) {
	_, facts := thrown(test, "testdata/modules/support")
	for _, expect := range []expected{
		{"practice", "contest/contest.go",
			"Pipe is used from outside only by tests."},
		{"practice", "latest/latest.go",
			"Newest is used from outside only by tests."}} {
		if !holds(facts["hits"], expect.scroll, expect.file, expect.said) {
			test.Errorf("%s ends in test and is named for no package, "+
				"and was let by", expect.file)
		}
	}
}
