package temple

// Shown is called from outside only by an Example whose output is
// checked: the Example is its promise, so it is no practice sword.
func Shown() string { return "shown" }

// Unshown is called from outside only by an Example that checks no
// output, which promises nothing: a practice sword, and the master says
// how it is shown.
func Unshown() string { return "unshown" }
