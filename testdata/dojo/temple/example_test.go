package temple_test

import (
	"fmt"

	"example.com/dojo/temple"
)

// ExampleShown checks what it prints, so it is Shown's promise.
func ExampleShown() {
	fmt.Println(temple.Shown())
	// Output: shown
}

// ExampleUnshown checks nothing, so it promises nothing.
func ExampleUnshown() {
	_ = temple.Unshown()
}
