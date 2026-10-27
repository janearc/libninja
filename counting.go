package libninja

// nothing is the count of an empty thing, and where a count starts:
// defined so.
const nothing = 0

// lineOrigin is where a file's lines and columns are counted from, as
// go/token has defined them: one.
const lineOrigin = 1

// bothEnds is what a span of lines counts beyond its last less its
// first, as a span counted at both ends is defined: its first line too.
const bothEnds = 1

// positive is whether a count is more than nothing.
func positive(count int) bool {
	return count > nothing
}

// none is whether a count is nothing.
func none(count int) bool {
	return count == nothing
}

// spanned is how many lines run from first to last, both counted.
func spanned(first, last int) int {
	return last - first + bothEnds
}
