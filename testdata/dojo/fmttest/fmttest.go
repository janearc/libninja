// Package fmttest is made for tests, and real code in the village
// imports it.
package fmttest

import _ "fmt" // the package this one is named for

// Seep is what leaks.
func Seep() int { return 6 }
