package temple

import "pgregory.net/rapid"

// drawn takes a property's t outside a test file, where it may not
// stand alone.
func drawn(t *rapid.T) { _ = t }

var _ = drawn
