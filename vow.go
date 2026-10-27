package libninja

import (
	"bytes"
	"cmp"
	"fmt"
	"go/token"
	"os"
	"slices"
	"strings"
)

// vowMark is how a vow begins.
const vowMark = "//ninja:vow"

// fewestReasonWords is how many words a reason has at least, picked.
var fewestReasonWords = picked("fewest-reason-words")

// vow is a scroll broken on purpose. It holds which scroll and why. It
// holds where it is written and in which package. It holds whether it
// stands alone on its line: alone, it covers the line below it. At the
// end of a line of code, it covers that line and no other.
type vow struct {
	scroll  string
	reason  string
	at      token.Position
	student string
	alone   bool
}

// gatherVows finds every vow in the module, keyed by its file and
// line.
func (read *reading) gatherVows() {
	read.vows = map[string]vow{}
	for _, each := range read.files() {
		read.vowsIn(each)
	}
}

// vowsIn gathers the vows written in one file.
func (read *reading) vowsIn(each source) {
	body := textOf(each.name)
	for _, group := range each.file.Comments {
		for _, comment := range group.List {
			if !strings.HasPrefix(comment.Text, vowMark) {
				continue
			}
			said := strings.TrimPrefix(comment.Text, vowMark)
			words := strings.Join(strings.Fields(said), " ")
			written := vow{at: read.at(comment),
				student: read.address(each.pkg.PkgPath)}
			at := comment.Pos()
			written.alone = writtenAlone(read.fset, body, at)
			scroll, reason, _ := strings.Cut(words, " ")
			written.scroll, written.reason = scroll, reason
			key := lineOf(written.at.Filename, written.at.Line)
			read.vows[key] = written
		}
	}
}

// textOf is a file's bytes. When the file cannot be read it is empty.
// Every comment in it then stands alone, the old reading.
func textOf(filename string) []byte {
	body, _ := os.ReadFile(filename)
	return body
}

// newline ends a line.
const newline = "\n"

// writtenAlone is whether only space comes before a comment on the line
// it is written on. The line is read from body, the file it is written
// in. Directives such as //line move where a comment is said to be. They
// do not move where it is written, so that is where it is judged.
func writtenAlone(fset *token.FileSet, body []byte, at token.Pos) bool {
	end := min(fset.File(at).Offset(at), len(body))
	start := bytes.LastIndex(body[:end], []byte(newline)) + len(newline)
	return strings.TrimSpace(string(body[start:end])) == ""
}

// lineOf is a file and line as one key.
func lineOf(filename string, line int) string {
	return fmt.Sprintf("%s:%d", filename, line)
}

// vowFor is the vow that covers a finding, and where it is written. The
// vow names the finding's scroll with a reason. It sits at the end of
// the finding's line or alone on the line above.
func (read reading) vowFor(found finding) (vow, string, bool) {
	for _, line := range []int{found.at.Line, found.at.Line - above} {
		key := lineOf(found.at.Filename, line)
		written, ok := read.vows[key]
		if !ok || written.alone != (line < found.at.Line) {
			continue
		}
		kept := written.scroll == found.scroll.name
		if kept && reasoned(written.reason) {
			return written, key, true
		}
	}
	return vow{}, "", false
}

// idleVows is every vow with a scroll and a reason that covers no
// finding, by the keys of the vows that did.
func (read reading) idleVows(used map[string]bool) []finding {
	var found []finding
	for key, written := range read.vows {
		_, known := named(written.scroll)
		if used[key] || !known ||
			!reasoned(written.reason) {
			continue
		}
		found = append(found, finding{emptyVow, written.student,
			written.at,
			fmt.Sprintf("The vow to break %s covers "+
				"nothing here. Put it down.",
				written.scroll), ""})
	}
	return found
}

// emptyVows is every vow that names no scroll or gives no reason.
func (read reading) emptyVows() []finding {
	var found []finding
	for _, written := range read.vows {
		_, known := named(written.scroll)
		switch {
		case !known:
			found = append(found, finding{emptyVow, written.student,
				written.at,
				fmt.Sprintf("The vow names %q, which is "+
					"no scroll.", written.scroll), ""})
		case written.reason == "":
			found = append(found, finding{emptyVow, written.student,
				written.at,
				fmt.Sprintf("The vow to break %s gives "+
					"no reason.", written.scroll), ""})
		case !reasoned(written.reason):
			found = append(found, finding{emptyVow, written.student,
				written.at,
				fmt.Sprintf("The vow to break %s gives too "+
					"little reason. Say why, in %d words "+
					"at least.", written.scroll,
					fewestReasonWords), ""})
		}
	}
	return found
}

// carriedScrolls is the scroll each vow in the module names. There is
// one for every vow written, in the order the scrolls are taught. Any
// name that is no scroll comes last, as it was written.
func (read reading) carriedScrolls() []string {
	var names []string
	for _, written := range read.vows {
		names = append(names, written.scroll)
	}
	slices.SortFunc(names, taughtFirst)
	return names
}

// taught is where a scroll comes in the teaching, or past the last for
// a name that is no scroll.
func taught(name string) int {
	each, known := named(name)
	if !known {
		return len(scrolls())
	}
	return each.taught()
}

// taughtFirst orders two scroll names by the order of teaching, and
// by their letters within one that is no scroll.
func taughtFirst(one, other string) int {
	place, otherPlace := taught(one), taught(other)
	return cmp.Or(cmp.Compare(place, otherPlace), cmp.Compare(one, other))
}

// reasoned is whether a reason has words enough to say anything.
func reasoned(reason string) bool {
	return len(strings.Fields(reason)) >= fewestReasonWords
}
