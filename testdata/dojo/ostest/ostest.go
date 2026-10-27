// Package ostest is made for tests, and no test imports it. It is in
// two files, and the master addresses it at the first.
package ostest

import _ "os" // the package this one is named for

// Ash is what is left.
func Ash() int { return 9 }
