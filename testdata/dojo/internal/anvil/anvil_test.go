package anvil_test

import (
	"fmt"

	"example.com/dojo/internal/anvil"
)

// ExampleStrike checks what it prints, which would promise Strike were
// anvil not internal.
func ExampleStrike() {
	fmt.Println(anvil.Strike())
	// Output: clang
}
