package village

// Twice runs a function with no name, two times.
func Twice() int {
	count := 0
	add := func() { count++ }
	add()
	add()
	return count
}
