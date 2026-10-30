// Package libninja is the way of the ninja: rules for Go written by
// agents, each rule a scroll, each scroll checked by reading the whole
// module at once.
//
// The scrolls forbid what stops agent-built code from growing: fields
// left open, another package's structs built by literal, exported
// variables, exports nobody outside calls, exports only tests call,
// and tests that look inside the package.
//
// Also the same work forged twice, functions too long to see at once,
// one-letter names, a main that cannot say which commit it is, and a
// path that is true on one machine written into the code, and a struct
// handed back beside an error, for a caller to use unchecked.
//
// A scroll is broken only with a vow, written where it is broken, on
// the line above the finding or at the end of its line:
//
//	//ninja:vow gate yaml reads these fields; this is the wire
//
// A vow names its scroll and gives a reason. A vow alone on its line
// covers the line below; one at the end of a line covers that line
// only. A vow without a reason, or with fewer than three words of it,
// covers nothing, and a vow that covers nothing is itself a hit.
//
// A module that carries more than one vow to every twenty tests is a
// hit too: its vows go to the master.
//
// Every vow that covers a finding is printed with it on every run.
// The report counts every vow the module carries, by scroll, beside the
// hits.
//
// shuriken, in cmd/shuriken, is the command that throws them.
package libninja
