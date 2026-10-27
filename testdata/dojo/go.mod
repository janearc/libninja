module example.com/dojo

go 1.26

require (
	github.com/janearc/puzzlemouth v0.0.0
	pgregory.net/rapid v0.0.0
)

// puzzlemouth and rapid here are the dojo's stand-ins, at their own
// paths.
replace (
	github.com/janearc/puzzlemouth => ../puzzlemouth
	pgregory.net/rapid => ../rapidstandin
)
