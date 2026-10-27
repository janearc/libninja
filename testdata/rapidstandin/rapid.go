// Package rapid stands in for pgregory.net/rapid in the dojo: its T is
// what a property's t is typed.
package rapid

// T is what a property is given.
type T struct{}

// Check runs a property once.
func Check(_ any, property func(*T)) { property(&T{}) }
