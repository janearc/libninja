# changelog

## unreleased

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
