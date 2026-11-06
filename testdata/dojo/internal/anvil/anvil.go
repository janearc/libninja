// Package anvil is internal: go lets only the dojo's own tree import
// it, so no Example can promise its names to anyone else.
package anvil

// Strike is called from outside only by an Example whose output is
// checked. Being internal, that promises nothing, and it is a practice
// sword as if no Example showed it.
func Strike() string { return "clang" }
