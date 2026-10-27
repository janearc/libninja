package temple

// Rope has a vow at the end of one line, which covers that line only:
// Loose, on the next, is still open.
type Rope struct {
	Knot  string //ninja:vow gate the boat reads the knot
	Loose string
}

//ninja:vow name nothing here is one letter long
func calm() int { return 0 }
