package temple

// Box holds anything, and says so as a string.
type Box[Held any] struct{ held Held }

// String keeps fmt.Stringer's promise for every box, so it stays
// exported though a generic type cannot be asked what it implements.
func (box Box[Held]) String() string { return "box" }

// Drill is used from outside only by temple's own external test.
func Drill() int { return 5 }
