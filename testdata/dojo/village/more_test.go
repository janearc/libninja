package village_test

import (
	"fmt"
	"os"
	"testing"

	"example.com/dojo/temple"
)

// TestTheVillageSpeaks makes in a test what outside one would be a
// forged composition, a buried path and a blade. No scroll reads a test
// for those.
func TestTheVillageSpeaks(test *testing.T) {
	_ = fmt.Sprint(temple.Forge())
	_, _ = os.ReadFile("/tmp/dojo-test.txt")
	if twiceOne(3) != twiceTwo(3) {
		test.Error("the twins differ")
	}
}

func twiceOne(times int) int {
	total := 0
	for step := 0; step < times; step++ {
		total += step * 2
	}
	if total > 10 {
		total -= 3
	}
	total *= times
	total = total - times
	return (total + times)
}

func twiceTwo(times int) int {
	total := 0
	for step := 0; step < times; step++ {
		total += step * 2
	}
	if total > 10 {
		total -= 3
	}
	total *= times
	total = total - times
	return (total + times)
}
