// Package gittest supports tests of git, as httptest does for http.
package gittest

import "example.com/support/git"

// Repo is a repository for a test.
func Repo() int { return git.Clone() }
