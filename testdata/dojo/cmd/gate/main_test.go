package main

import "testing"

// a main's test can only stand inside: nothing outside imports a main.
func TestTheGateOpens(test *testing.T) {}
