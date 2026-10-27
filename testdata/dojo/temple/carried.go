package temple

// Belt is what a student wears out of the temple. No one outside writes
// its name: Tie hands one to the village.
type Belt struct{ colour string }

// Tie is a belt of a colour.
func Tie(colour string) Belt { return Belt{colour: colour} }

// Sash is worn only in practice. Wrap hands one out, and only a test
// calls Wrap.
type Sash struct{ knots int }

// Wrap is a sash with its knots.
func Wrap(knots int) Sash { return Sash{knots: knots} }
