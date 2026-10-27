package libninja_test

import (
	"bytes"
	"strings"
	"testing"
)

// spokenOf is what the master says of a house with one more file.
// That file is more.go, holding the given body below its package line.
func spokenOf(test *testing.T, body string) string {
	test.Helper()
	dir := house(test, ".", map[string]string{
		".editorconfig": "root = true\n\n[*.go]\nindent_style = tab\n",
		"more.go":       "package house\n\n" + body,
	})
	report, _ := thrown(test, dir)
	var spoken bytes.Buffer
	if err := report.Speak(&spoken); err != nil {
		test.Fatal(err)
	}
	return spoken.String()
}

// The master says how many vows the module carries beside the hits.
// The count goes scroll by scroll in the order taught. Every vow that
// covers nothing is carried too.
func TestTheHeaderCountsTheVowsByScroll(test *testing.T) {
	body := "// First is open.\nvar First = map[string]string{} " +
		"//ninja:vow well the lamp is read by the wire\n\n" +
		"// Second is open.\nvar Second = map[string]string{} " +
		"//ninja:vow well the second lamp too, for the wire\n\n" +
		"// calm is calm.\nfunc calm() {} " +
		"//ninja:vow nest nothing here is nested at all\n"
	said := spokenOf(test, body)
	want := "2 hits in 1 packages, 3 vows (well 2, nest 1).\n"
	if !strings.Contains(said, want) {
		test.Errorf("the header did not say %q:\n%s", want, said)
	}
}

// With no vows the master says none, and names no scroll.
func TestTheHeaderOfAModuleWithoutVowsNamesNone(test *testing.T) {
	said := spokenOf(test, "// Calm is calm.\nfunc Calm() {}\n")
	if !strings.Contains(said, "0 vows.\n") ||
		strings.Contains(said, "vows (") {
		test.Errorf("a module with no vows said:\n%s", said)
	}
}

// Fewer than three words of reason say nothing a reader can use. The
// vow covers nothing and is a hit of its own.
func TestAVowWithTooShortAReasonIsAHitAndCoversNothing(test *testing.T) {
	for _, reason := range []string{"ok", "the wire", "it"} {
		body := "// Open is open.\nvar Open = map[string]string{} " +
			"//ninja:vow well " + reason + "\n"
		vows := hitsIn(test, "vow", body)
		wells := hitsIn(test, "well", body)
		if strings.Join(vows, "|") != "The vow to break well gives too "+
			"little reason. Say why, in 3 words at least." ||
			len(wells) != one {
			test.Errorf("reason %q: vow hits %q, well hits %q", reason,
				vows, wells)
		}
	}
}

// Three words is the picked floor. The vow covers the finding and is no
// hit.
func TestAVowWithThreeWordsOfReasonCovers(test *testing.T) {
	body := "// Open is open.\nvar Open = map[string]string{} " +
		"//ninja:vow well the wire reads\n"
	if vows, wells := hitsIn(test, "vow", body),
		hitsIn(test, "well", body); vows != nil || wells != nil {
		test.Errorf("three words: vow hits %q, well hits %q", vows, wells)
	}
}

// No reason at all is still said to be none.
func TestAVowWithNoReasonStillGivesNone(test *testing.T) {
	body := "// Open is open.\nvar Open = map[string]string{} " +
		"//ninja:vow well\n"
	vows := hitsIn(test, "vow", body)
	if strings.Join(vows, "|") !=
		"The vow to break well gives no reason." {
		test.Errorf("no reason: %q", vows)
	}
}
