// Package village visits the temple.
package village

import "example.com/dojo/temple"

// Visit builds a sword in the temple's house, and takes its edge.
func Visit() float64 {
	sword := temple.Forge()
	//ninja:vow stranger
	made := temple.Sword{Edge: 2}
	_ = temple.Tie("black")
	return sword.Edge + made.Edge + float64(temple.Hidden)
}
