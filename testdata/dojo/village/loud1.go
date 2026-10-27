package village

import (
	"strings"

	"example.com/dojo/temple"
)

// loudOne ties a belt by a shouted colour. The call of a call is forged
// in three files, and temple's half is outside the standard library.
func loudOne() temple.Belt { return temple.Tie(strings.ToUpper("black")) }

var _ = loudOne
