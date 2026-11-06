package temple

// Elsewhere is called from outside only by an Example in the village's
// tests, which checks its output. Another package's Example is not
// temple's documentation, so it promises nothing: a practice sword.
func Elsewhere() string { return "elsewhere" }
