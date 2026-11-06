# changelog

## unreleased

- practice: an export used from outside only by an Example whose
  output is checked (`// Output:`) is carried, not a practice sword: the
  Example is the promise a library makes to callers outside the module,
  shown by godoc and run by go test.

  The Example must be the name's own package's, or its _test package's:
  another package's Example is not that name's documentation.

  An Example that checks no output carries nothing, and the hit says
  so: "is shown in an example that checks no output."

  A package with `internal` among its path's elements gets no such
  promise, since go lets only its own tree import it, and its hit says
  so: "is shown in an example, and the package is internal: no one
  outside its tree can be promised it."

  The maintainer's ruling, 2026-11-05, after celestial's solver, a
  library whose callers were not yet written.

- new scroll, the defensive feint (`feint`): a return that gives an
  error that is not nil, and beside it a struct that carries nothing
  (`T{}`, or a variable declared and never set before the return), is a
  hit.

  a struct with fields set beside an error is just as much a hit; its
  message says to check the error first. the three kinds are worded
  apart and are equal. basic types and nils are let be.
  a vow names it as `feint`.

- breaking: `Load` returns `*Module`, and nil with the error on a
  refusal, where it returned a `Module` beside the error. a caller that
  used the zero `Module` after an error now meets nil. calls on the
  result (`module.Throw()`) read the same.
