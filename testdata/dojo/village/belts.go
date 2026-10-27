package village

import "example.com/dojo/temple"

// belts builds temple's map by literal, which is no struct and no
// stranger's house.
func belts() temple.Belts { return temple.Belts{"black": 1} }

var _ = belts
