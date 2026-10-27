package village

import "example.com/dojo/fmttest"

// leak lets a package made for tests into real code.
func leak() int { return fmttest.Seep() }
