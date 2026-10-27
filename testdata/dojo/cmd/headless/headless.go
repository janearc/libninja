// Command headless is a main package with no func main. go's type
// checker lets it by, so the age scroll places it at its package clause.
package main

// motto is all it holds.
var motto = "no head"

var _ = motto
