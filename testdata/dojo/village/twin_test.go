package village_test

import (
	"fmt"
	"testing"

	"example.com/dojo/temple"
)

// TestTheTwinsForge composes the same call as two other test files.
func TestTheTwinsForge(test *testing.T) { _ = fmt.Sprint(temple.Forge()) }
