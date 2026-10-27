// Package temple breaks the scrolls, so the tests can see them hit.
package temple

// Sword's edge stands open.
type Sword struct {
	Edge  float64
	sharp bool
}

// Wire is open on purpose, and says so.
type Wire struct {
	//ninja:vow gate yaml reads it; this is the wire
	Name string `yaml:"name"`
}

// Hidden is a variable anyone may change.
var Hidden = 3

// Forge is a sword. The village calls it.
func Forge() Sword { return Sword{Edge: 1, sharp: true} }

// Hone is a sword made sharp. The village calls it.
func Hone(sword Sword) Sword {
	sword.sharp = true
	return sword
}

// Idle is called by no one outside.
func Idle() int { return 1 }

// Practice is called from outside only by a test.
func Practice() int { return 2 }

// String keeps fmt.Stringer's promise, so it stays exported.
func (sword Sword) String() string { return "sword" }

func one(q int) int { return q }
