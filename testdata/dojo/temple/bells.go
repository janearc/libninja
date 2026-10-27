package temple

// Bell rings. A test makes its mock, as tests should.
type Bell interface{ Ring() int }

// Gong is struck. Its mock is made by no test; a test writes its own.
type Gong interface{ Strike() string }

// tide is a package variable a test in the temple swaps.
var tide = 1

// Weather is a package variable a test outside the temple swaps.
var Weather = "clear"

// Tide is the tide, so the village can read it.
func Tide() int { return tide }

// Garden is a place with a gate.
type Garden struct{ Gate string }

// Yard, Shelves and Counter are package variables a test reaches into:
// a field, an element and what a pointer holds.
var (
	Yard    = Garden{Gate: "shut"}
	Shelves = []int{1}
	Counter = new(int)
)

// Plot is a garden of its own, made fresh for whoever asks.
func Plot() *Garden { return &Garden{Gate: "shut"} }
