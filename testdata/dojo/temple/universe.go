package temple

import "errors"

// Steel is another name for Sword, for whoever wants it.
type Steel = Sword

// rack holds swords by reference; the literal leaves their type out.
var rack = []*Sword{{Edge: 1}}

// failing is what a failure says, through the universe's own Error.
func failing() string { return errors.New("dull").Error() }

//ninja:vow nosuch a vow that names no scroll
var spare = rack
