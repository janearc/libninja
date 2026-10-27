// Package drums is no package for tests by its name, and only a test
// calls it.
package drums

// Roll is called from outside only by a test.
func Roll() int { return 8 }
