package village

import (
	"strings"

	"example.com/dojo/temple"
)

// loudTwo ties a belt by a shouted colour, as loudOne does.
func loudTwo() temple.Belt { return temple.Tie(strings.ToUpper("white")) }

var _ = loudTwo
