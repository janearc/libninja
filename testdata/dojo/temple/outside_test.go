package temple_test

import (
	"testing"

	"pgregory.net/rapid"

	"example.com/dojo/mocks"
	"example.com/dojo/temple"
)

// TestDrill is temple's external test, and the only caller of Drill.
func TestDrill(test *testing.T) {
	if temple.Drill() != 5 {
		test.Error("the drill is off")
	}
}

// TestTheSashIsWrapped is the only caller of Wrap.
func TestTheSashIsWrapped(test *testing.T) {
	_ = temple.Wrap(3)
}

// TestTheBeltTies names its testing.T t, as go's own tests do.
func TestTheBeltTies(t *testing.T) {
	_ = temple.Tie("white")
	t.Log("tied")
}

// TestTheBeltIsDrawn is a property. Its t is rapid's, and stands alone;
// the r of the second is not a t, and may not.
func TestTheBeltIsDrawn(test *testing.T) {
	rapid.Check(test, func(t *rapid.T) { _ = temple.Tie("red") })
	rapid.Check(test, func(r *rapid.T) { _ = temple.Tie("blue") })
}

// fakeGong is a gong the test carves for itself, though mockgen made one.
type fakeGong struct{}

// Strike is the carved gong's sound.
func (fakeGong) Strike() string { return "bong" }

// worn is mockgen's bell, worn inside a test's own type: still the mock.
type worn struct{ *mocks.MockBell }

// rounds is the test's own count, which it may change as it likes.
var rounds = 0

// gongSound is a gong carved as a function: still a stand-in.
type gongSound func() string

// Strike is whatever the function sounds.
func (sound gongSound) Strike() string { return sound() }

// wornByValue holds mockgen's bell itself, not a pointer to it.
type wornByValue struct{ mocks.MockBell }

// plain is a struct with no name of its own.
type plain = struct{}

// tolled holds a type that is no mock and a struct that has no name.
type tolled struct {
	temple.Belt
	plain
}

// TestTheBellsRing makes the bell's mock, carves a gong and swaps the
// temple's weather, and reaches into its yard, shelves and counter.
func TestTheBellsRing(test *testing.T) {
	bell := worn{mocks.NewMockBell()}
	var gong temple.Gong = fakeGong{}
	temple.Weather = "rain"
	temple.Yard.Gate = "open"
	temple.Shelves[0] = 2
	*temple.Counter = 3
	temple.Plot().Gate = "open"
	rounds = bell.Ring()
	_, _, _, _ = wornByValue{}, tolled{}, recorded{}, belted{}
	if gong.Strike() != "bong" || rounds != 0 {
		test.Fatal("the bells are off")
	}
}

// recorded holds mockgen's recorder, which is no mock, and strikes by
// its own method: a stand-in for Gong.
type recorded struct{ *mocks.MockBellMockRecorder }

// Strike is the recorded gong's own sound.
func (recorded) Strike() string { return "tap" }

// belted holds a belt, which is no mock, and strikes by its own method:
// a stand-in for Gong.
type belted struct{ temple.Belt }

// Strike is the belted gong's own sound.
func (belted) Strike() string { return "thud" }

// TestTheGongIsCounted bumps the temple's count of strikes, a swap.
func TestTheGongIsCounted(test *testing.T) { temple.Strikes++ }
