package village

// Head is the first of the words and the rest after it, each by its
// place.
func Head(words []string) (string, []string) {
	return words[0], words[1:]
}
