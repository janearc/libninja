package village

import "example.com/dojo/temple"

// borrowed is temple's sword under another name. Building it by literal
// is still building in temple's house.
type borrowed = temple.Sword

// disguise builds the sword under its borrowed name.
func disguise() borrowed { return borrowed{Edge: 3} }
