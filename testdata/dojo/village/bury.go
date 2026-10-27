package village

import "os"

// bury writes where only this machine can find it.
func bury() error {
	return os.WriteFile("/tmp/dojo.txt", nil, 0o644)
}

var _ = bury
