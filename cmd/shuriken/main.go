// Command shuriken throws libninja's scrolls at a Go module and says
// what they hit, in the master's voice.
//
//	shuriken throw [dir]
//	shuriken facts [dir]
//
// dir is the module's directory, the current one if none is given. A
// directory inside a module means the whole module it is in, and it may
// be written the go way, ./..., which means the same thing.
//
// facts writes the findings as JSON, for game and other programs. Either
// fails when a scroll is broken without a vow, or when the module cannot
// be read, and says which. It is on pzm, which answers --age and --who.
package main

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/janearc/libninja"
	"github.com/janearc/puzzlemouth/pzm"
	"github.com/janearc/puzzlemouth/verb"
)

// the four stamps game build writes into a command's main package
var build, built, who, sessions string

// afterName is where the arguments begin in os.Args, as go has defined
// it: the program's own name comes first.
const afterName = 1

// exit ends the process with a code. It is os.Exit; a test holds it to
// run main and keep going.
//
//ninja:vow decoy main ends the process, and no package imports main
var exit = os.Exit

// main runs shuriken on what it was asked.
func main() {
	exit(tool().Run(os.Args[afterName:], os.Stdout, os.Stderr))
}

// tool is shuriken and its two verbs.
func tool() *pzm.Tool {
	made := pzm.New("shuriken", "throw libninja's scrolls at a go module",
		pzm.Build(build), pzm.Built(built), pzm.Who(who),
		pzm.Sessions(sessions))
	// Both verbs read and type-check the whole module, which holds a
	// core while it runs. Neither writes, nor reaches past this machine.
	reads := verb.NewFacts(verb.WritesNothing, verb.LeavesAlone,
		verb.ThisMachine, verb.NoOne, verb.FullCPU)
	made.Verb("throw", 't', "say what the scrolls hit, in the master's "+
		"voice", throwPage, "Example_throw", reads, throw)
	made.Verb("facts", 'f', "write what the scrolls hit as json facts",
		factsPage, "Example_facts", reads, facts)
	return made
}

const throwPage = `# throw

reads the whole module the directory is in, the current one if none is
named, throws every scroll at it, and says what each hit and why, with
the vows that cover the rest. it fails when a scroll is broken without
a vow. it reads the module and writes nothing.`

const factsPage = `# facts

the same throw, written as json for game and other programs: each hit
and each vow with its scroll, file, line and what it saw. it fails as
throw does.`

// errHit is a throw that hit, told apart from one that could not read.
var errHit = errors.New("a scroll is broken without a vow")

// throw says what the scrolls hit.
func throw(call pzm.Call) error {
	return thrown(call, libninja.Report.Speak)
}

// facts writes what the scrolls hit as json.
func facts(call pzm.Call) error {
	return thrown(call, libninja.Report.JSON)
}

// thrown reads the module the call names and throws every scroll at it.
// It writes the report as asked and fails when a scroll hit.
func thrown(call pzm.Call,
	write func(libninja.Report, io.Writer) error) error {
	dir := dirOf(call.Args())
	module, err := libninja.Load(dir)
	if err != nil {
		return err
	}
	report := module.Throw()
	if err := write(report, call.Out()); err != nil {
		return err
	}
	if !report.Clean() {
		return errHit
	}
	return nil
}

// named is the place of the one argument a verb takes, as shuriken's
// verbs are defined: the directory.
const named = 0

// dirOf is the directory the arguments name, the current one if none;
// a go pattern's /... names the same module.
func dirOf(args []string) string {
	if len(args) <= named {
		return "."
	}
	pattern := strings.TrimSuffix(args[named], "...")
	dir := strings.TrimSuffix(pattern, "/")
	if dir == "" {
		return "."
	}
	return dir
}
