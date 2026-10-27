// Command tool is made with pzm's New, so pzm answers --age for it. It
// names no flag of its own.
package main

import "github.com/janearc/puzzlemouth/pzm"

// the stamps game build writes.
var (
	build = "dev"
	built = ""
)

func main() { _ = pzm.New("tool", "stands in", pzm.Build(build), pzm.Built(built)) }
