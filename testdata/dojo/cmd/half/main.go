// Command half uses pzm, but makes no tool with it: it takes only the
// stamp of when it was built, so nothing answers --age for it.
package main

import "github.com/janearc/puzzlemouth/pzm"

// the stamps game build writes.
var (
	build = "dev"
	built = ""
)

func main() {
	_ = build
	_ = pzm.Built(built)
}
