package temple

import "fmt"

// Showy is shown only by an Example in this file, which is no test
// file: go test runs no Example here, so it promises nothing, and
// nothing outside uses Showy.
func Showy() string { return "showy" }

// Example_showy checks its output, but go test never runs an Example
// outside a test file, so it is no promise. Its name shares no word
// with Showy's hit, so a test that looks for that hit cannot find this
// one instead.
func Example_showy() {
	fmt.Println(Showy())
	// Output: showy
}
