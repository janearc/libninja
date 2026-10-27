package village_test

import (
	"fmt"
	"testing"

	"example.com/dojo/temple"
)

func TestPractice(test *testing.T) {
	if temple.Practice() != 2 {
		test.Fatal("practice is not two")
	}
}

// TestTheVillageForges composes a call in a test file, as two other
// test files do; no scroll counts a test's compositions.
func TestTheVillageForges(test *testing.T) { _ = fmt.Sprint(temple.Forge()) }
