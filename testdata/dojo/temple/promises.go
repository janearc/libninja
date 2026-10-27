package temple

// Cord has methods named as fmt.Stringer's and error's are, of other
// types, so it keeps neither promise and both are waste.
type Cord struct{}

// String is no Stringer's: it gives a number.
func (Cord) String() int { return 0 }

// Error is no error's: it gives a number.
func (Cord) Error() int { return 0 }

// Knot is an error only by its pointer, which keeps the promise.
type Knot struct{}

// Error answers error for a *Knot.
func (*Knot) Error() string { return "knot" }

// Drawer draws a number.
type Drawer interface{ Draw() int }

// Quiver is generic, so its Draw is matched to Drawer's by name alone.
type Quiver[Held any] struct{ held Held }

// Draw keeps Drawer's promise by its name.
func (quiver Quiver[Held]) Draw() Held { return quiver.held }

// Holder is generic, so no interface in sight is it: Hold keeps nothing.
type Holder[Held any] interface{ Hold() Held }

// Crate is generic, and its Hold is asked for by no interface in sight.
type Crate[Held any] struct{ held Held }

// Hold is waste: only Holder asks for it, and Holder is not in sight.
func (crate Crate[Held]) Hold() Held { return crate.held }
