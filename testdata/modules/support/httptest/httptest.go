// Package httptest supports tests of net/http, and says so by import.
package httptest

import _ "net/http" // the package this one is named for

// Server is a server for a test.
func Server() int { return 3 }
