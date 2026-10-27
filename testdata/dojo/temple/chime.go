package temple

// Anything is what every value is. Every mock satisfies it, and it is
// the mock of none.
type Anything interface{}

// Chime is rung by hand.
type Chime interface{ Chime() string }

// MockChime is a mock of Chime interface.
//
// It is written by hand, outside a test, with mockgen's words above it.
// The decoy scroll asks only the files mockgen wrote, so it is no mock.
type MockChime struct{}

// NewMockChime makes the hand-written chime.
func NewMockChime() *MockChime { return &MockChime{} }

// Chime is the hand-written chime's sound.
func (*MockChime) Chime() string { return "ting" }

// Drum is beaten. Its mock in mocks has no constructor.
type Drum interface{ Beat() string }

// brass is the temple's own gong, the real one: no stand-in.
type brass struct{}

// Strike is the brass gong's sound.
func (brass) Strike() string { return "bong" }

// season is a package variable the temple's own code changes, which no
// test swaps.
var season = "spring"

// turn moves the season on, in code that is not a test.
func turn() string {
	season = "autumn"
	return season
}

var _, _ = brass{}, turn

// Strikes is how many times the gong was struck. A test bumps it with
// ++, which writes it as surely as an assignment does.
var Strikes int
