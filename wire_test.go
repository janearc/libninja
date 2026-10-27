package libninja_test

import (
	"strings"
	"testing"
)

// one is the count of a single finding, defined so.
const one = 1

// openField is whether the gate scroll found exactly one open field,
// and it is the one named.
func openField(seen []string, name string) bool {
	return len(seen) == one && strings.HasPrefix(strings.Join(seen, ""), name)
}

// encoding/json cannot reach a private field, so a field with a json
// tag is exported for it. It is not surface to guard.
func TestAJSONTaggedFieldMayStandOpen(test *testing.T) {
	for _, tag := range []string{`json:"name"`, `json:"name,omitempty"`,
		`json:",omitempty"`, `yaml:"name" json:"name"`} {
		body := "// Wire is read by json.\ntype Wire struct {\n" +
			"\tName string `" + tag + "`\n}\n"
		if seen := hitsIn(test, "gate", body); seen != nil {
			test.Errorf("tag %s: %v", tag, seen)
		}
	}
}

// The near misses stay open. They are no tag, a tag for another reader,
// a json tag that says the field is not written, and a key that only
// ends in json.
func TestAFieldWithoutAJSONTagStillStandsOpen(test *testing.T) {
	for _, field := range []string{"Name string",
		"Name string `yaml:\"name\"`",
		"Name string `json:\"-\"`",
		"Name string `xjson:\"name\"`"} {
		body := "// Wire is read by someone.\ntype Wire struct {\n\t" +
			field + "\n}\n"
		if seen := hitsIn(test, "gate", body); !openField(seen,
			"Wire.Name stands open") {
			test.Errorf("field %s: %v", field, seen)
		}
	}
}

// Of two fields, only the json-tagged one is let by.
func TestOnlyTheJSONTaggedFieldOfTwoIsLetBy(test *testing.T) {
	body := "// Wire is read by json.\ntype Wire struct {\n" +
		"\tName string `json:\"name\"`\n\tOpen string\n}\n"
	if seen := hitsIn(test, "gate", body); !openField(seen,
		"Wire.Open stands open") {
		test.Errorf("one tagged, one not: %v", seen)
	}
}
