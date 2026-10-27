// Package importsnothing imports a package that is not there, which is
// one error and no more.
package importsnothing

import "example.com/importsnothing/nowhere"

var _ = nowhere.X
