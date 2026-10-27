package libninja_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/janearc/libninja"
)

// golden is one of the dojo's reports as a file holds it, and how the
// report is written.
type golden struct {
	file  string
	write func(libninja.Report, *bytes.Buffer) error
}

// goldens are the dojo's report as the master speaks it, and as JSON.
var goldens = []golden{{"testdata/dojo.spoken", spoken},
	{"testdata/dojo.json", listed}}

// spoken writes a report as the master speaks it.
func spoken(report libninja.Report, out *bytes.Buffer) error {
	return report.Speak(out)
}

// listed writes a report as JSON.
func listed(report libninja.Report, out *bytes.Buffer) error {
	return report.JSON(out)
}

// writeGoldens names the setting that writes the goldens again, on
// purpose.
const writeGoldens = "WRITE_GOLDENS"

// The dojo's report is said exactly as its goldens hold it. That covers
// each hit's place and words, the order taught and the wrapping. It
// covers the vows and the tally too. With WRITE_GOLDENS set, the test
// writes them again and stops. So new files are read before they stay.
func TestTheDojoIsSpokenAsWritten(test *testing.T) {
	report, _ := thrown(test, "testdata/dojo")
	writing := os.Getenv(writeGoldens) != ""
	mode := os.FileMode(heldTo("open-file"))
	for _, each := range goldens {
		var got bytes.Buffer
		if err := each.write(report, &got); err != nil {
			test.Fatal(err)
		}
		if writing {
			if err := os.WriteFile(each.file, got.Bytes(), mode); err != nil {
				test.Fatal(err)
			}
			continue
		}
		want, err := os.ReadFile(each.file)
		if err != nil {
			test.Fatal(err)
		}
		said, held := got.String(), string(want)
		if said != held {
			parted := firstDifference(said, held)
			test.Errorf("%s is not what the dojo says now, from line %d. "+
				"if the change is meant, run with %s=1 and read the new "+
				"file before keeping it", each.file, parted, writeGoldens)
		}
	}
	if writing {
		test.Fatal("the goldens were written. read them, then run again " +
			"without " + writeGoldens)
	}
}

// The dojo is spoken alike wherever its tree sits. Its copy in the
// system's temp has nothing of the repository above it. That copy says
// what the dojo in the repository says, in the same order.
func TestTheDojoIsSpokenAlikeWhereverItSits(test *testing.T) {
	copied := test.TempDir()
	if err := os.CopyFS(copied, os.DirFS("testdata")); err != nil {
		test.Fatal(err)
	}
	here, _ := thrown(test, "testdata/dojo")
	there, _ := thrown(test, filepath.Join(copied, "dojo"))
	var said, copySaid bytes.Buffer
	if err := here.Speak(&said); err != nil {
		test.Fatal(err)
	}
	if err := there.Speak(&copySaid); err != nil {
		test.Fatal(err)
	}
	if said.String() != copySaid.String() {
		parted := firstDifference(said.String(), copySaid.String())
		test.Errorf("the dojo's copy in %s is spoken otherwise, from "+
			"line %d", copied, parted)
	}
}

// lineOrigin is the number of a text's first line, defined so.
const lineOrigin = 1

// firstDifference is the number of the first line where two texts part.
func firstDifference(one, other string) int {
	ones, others := strings.Split(one, "\n"), strings.Split(other, "\n")
	shorter := min(len(ones), len(others))
	for place := range shorter {
		if ones[place] != others[place] {
			return place + lineOrigin
		}
	}
	return shorter + lineOrigin
}
