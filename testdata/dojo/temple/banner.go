package temple

// Banner's fields stand open, and no one outside touches them: Title
// says nothing more, and Motto's tag says a reader needs it.
type Banner struct {
	Title string
	Motto string `yaml:"motto"`
}

// Raise is a method no one outside calls.
func (banner Banner) Raise() string { return banner.Title }

// Belts is a map by name: built by literal outside, it is no stranger's
// house, since it is no struct.
type Belts map[string]int
