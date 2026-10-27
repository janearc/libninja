package temple

import "testing"

func TestInside(t *testing.T) {
	if one(1) != 1 {
		t.Fatal("one is not one")
	}
}

func TestInsideAgain(test *testing.T) {
	if t := one(2); t != 2 {
		test.Fatal("two is not two")
	}
}

// TestTheTideTurns reaches into the temple to swap its tide.
func TestTheTideTurns(test *testing.T) {
	tide = 2
	(tide) = 3
	if Tide() != 3 {
		test.Fatal("the tide did not turn")
	}
}
