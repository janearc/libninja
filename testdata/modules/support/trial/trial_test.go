package trial_test

import (
	"testing"

	"example.com/support/contest"
	"example.com/support/gittest"
	"example.com/support/httptest"
	"example.com/support/latest"
)

// TestEveryPackageIsUsed calls each package beside it.
func TestEveryPackageIsUsed(test *testing.T) {
	if gittest.Repo()+httptest.Server()+contest.Pipe()+latest.Newest() == 0 {
		test.Fatal("nothing was used")
	}
}
