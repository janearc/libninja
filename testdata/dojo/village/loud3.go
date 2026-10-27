package village

import (
	"strings"

	"example.com/dojo/temple"
)

// loudThree ties a belt by a shouted colour, as loudOne does.
func loudThree() temple.Belt { return temple.Tie(strings.ToUpper("red")) }

var _ = loudThree
