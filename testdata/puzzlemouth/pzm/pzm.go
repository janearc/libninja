// Package pzm stands in for puzzlemouth's pzm in the dojo. A tool made
// with New answers --age, as every tool the real pzm makes does.
package pzm

// Tool is a command made with pzm.
type Tool struct{ options []Option }

// Option is a part of how a tool is made: here, a stamp game writes.
type Option struct{ stamp string }

// New makes a tool. The real one answers --age with its stamps.
func New(name, purpose string, options ...Option) *Tool {
	_, _ = name, purpose
	return &Tool{options: options}
}

// Build is the build game stamps, given to the tool.
func Build(stamp string) Option { return Option{stamp: stamp} }

// Built is when game stamped it.
func Built(stamp string) Option { return Option{stamp: stamp} }
