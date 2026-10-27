package village_test

import (
	"testing"

	"example.com/dojo/drums"
	"example.com/dojo/drumstest"
	"example.com/dojo/fmttest"
)

// TestTheDrumsBeat uses a package made for tests, one that real code
// leaks, and an ordinary package that only tests use.
func TestTheDrumsBeat(test *testing.T) {
	if drumstest.Beat()+fmttest.Seep()+drums.Roll() == 0 {
		test.Fatal("the drums are silent")
	}
}
