package libninja

import (
	_ "embed"

	"github.com/janearc/puzzlemouth/numbers"
)

// pickedFile is the picked numbers' file, as a refusal names it.
const pickedFile = "picked-numbers.yml"

// pickedRaw is the picked numbers as the file holds them.
//
//go:embed picked-numbers.yml
var pickedRaw []byte

// picked is a picked number by its name, as a whole count. The file is
// libninja's own and fixed when it is built, so a name it does not hold
// or a number that is not whole is a mistake in libninja's code.
//
// TestEveryPickedNumberIsInTheFile finds such a mistake before a build
// and says which. Here it would only be a count of nothing.
func picked(name string) int {
	chosen, _ := numbers.Read(pickedFile, pickedRaw)
	number, _ := chosen.Get(name)
	return int(number.Value())
}
