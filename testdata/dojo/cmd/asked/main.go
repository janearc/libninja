// Command asked answers --age by its own flag, from the stamps game
// writes, so the age scroll lets it be.
package main

import (
	"flag"
	"fmt"
)

// the stamps game build writes.
var (
	build = "dev"
	built = ""
)

func main() {
	asked := flag.Bool("age", false, "say which commit this is")
	flag.Parse()
	if *asked {
		fmt.Println(build, built)
	}
}
