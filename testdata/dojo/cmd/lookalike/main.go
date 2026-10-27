// Command lookalike calls a New that is not pzm's, so it makes no tool
// and does not answer --age.
package main

import (
	"errors"
	"fmt"
)

// the stamps game build writes.
var (
	build = "dev"
	built = ""
)

// options is a type with a method named main.
type options struct{}

func main() { fmt.Println(errors.New(build + built)) }

// main is a method too, which is not where the command starts.
func (options) main() {}

// later comes after main, and is not it either.
func later() {}

var _, _ = later, options.main
