package village_test

import "testing"

// visitor is a type the external test declares, and builds itself.
type visitor struct {
	name string
}

func TestAVisitorComes(test *testing.T) {
	if (visitor{name: "tony"}).name == "" {
		test.Fatal("no one came")
	}
}
