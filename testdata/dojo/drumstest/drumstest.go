// Package drumstest is made for tests, as net/http/httptest is, and only
// tests import it.
package drumstest

// Beat is what a test hears from the drum.
func Beat() int { return 4 }
