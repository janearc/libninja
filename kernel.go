package libninja

import (
	"fmt"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

// metalSuffix is how a Metal kernel's source file is named.
const metalSuffix = ".metal"

// kernels is every number written in a Metal kernel of the module.
//
// the maintainer's ruling, 10-18: the kernel holds no numbers of its own. Every
// number it uses is declared once in go and named. The prelude hands
// each one to the kernel, so a number typed into a .metal file is one
// go cannot see.
func (read reading) kernels() []finding {
	var found []finding
	for _, file := range read.metalFiles() {
		// when a kernel cannot be read, it holds no number the
		// scroll can see
		body, _ := os.ReadFile(file)
		rel, _ := filepath.Rel(read.dir, filepath.Dir(file))
		pkgPath := path.Join(read.module, filepath.ToSlash(rel))
		student := read.address(pkgPath)
		for _, number := range numbersIn(string(body)) {
			at := token.Position{Filename: file, Line: number.line,
				Column: number.column}
			seen := fmt.Sprintf("%s is the kernel's own number.",
				number.text)
			found = append(found, finding{kernelHeld, student, at,
				seen, ""})
		}
	}
	return found
}

// metalFiles is every .metal file under the module's root. It leaves out
// testdata, hidden directories and modules of their own.
func (read reading) metalFiles() []string {
	gathered := metalGathering{root: read.dir}
	filepath.WalkDir(read.dir, gathered.visit)
	return gathered.files
}

// metalGathering is a walk for .metal files from a module's root.
type metalGathering struct {
	root  string
	files []string
}

// visit keeps a .metal file and passes by what is not the module's.
func (gathered *metalGathering) visit(at string, entry fs.DirEntry,
	err error) error {
	if err != nil {
		return nil
	}
	if entry.IsDir() {
		if at != gathered.root && notTheModules(at, entry.Name()) {
			return filepath.SkipDir
		}
		return nil
	}
	if strings.HasSuffix(at, metalSuffix) {
		gathered.files = append(gathered.files, at)
	}
	return nil
}

// notTheModules is whether a directory below a module's root holds
// something other than the module's own code.
func notTheModules(at, name string) bool {
	if name == "testdata" || strings.HasPrefix(name, ".") {
		return true
	}
	_, err := os.Stat(filepath.Join(at, "go.mod"))
	return err == nil
}

// metalNumber is a number written in a kernel's source, and where.
type metalNumber struct {
	text         string
	line, column int
}

// numbersIn is every number written in a kernel's source. It leaves out
// comments, strings and the digits inside names like float3.
func numbersIn(source string) []metalNumber {
	reading := metalReading{text: []rune(source), line: lineOrigin,
		column: lineOrigin}
	for !reading.done() {
		reading.step()
	}
	return reading.found
}

// metalReading is a kernel's source being read, one token at a time.
type metalReading struct {
	text         []rune
	at           int
	line, column int
	found        []metalNumber
}

// done is whether the source is read to its end.
func (reading *metalReading) done() bool {
	return !reading.within(here)
}

// peek is the character a little way ahead, or none past the end.
func (reading *metalReading) peek(ahead int) rune {
	if !reading.within(ahead) {
		return unicode.ReplacementChar
	}
	return reading.text[reading.at+ahead]
}

// within is whether a place a little way ahead is in the source. done
// and peek ask it alike, so the end is one bound.
func (reading *metalReading) within(ahead int) bool {
	return reading.at+ahead < len(reading.text)
}

// here is how far ahead the character being read is, defined as no way
// at all.
const here = 0

// next is how far ahead the character after this one is, defined as one
// character.
const next = 1

// advance moves past one character, keeping the line and column.
func (reading *metalReading) advance() {
	if reading.peek(here) == '\n' {
		reading.line++
		reading.column = lineOrigin
	} else {
		reading.column++
	}
	reading.at++
}

// step reads one thing: a comment or a string or a name or a number or
// one character of anything else.
func (reading *metalReading) step() {
	this, after := reading.peek(here), reading.peek(next)
	switch {
	case this == '/' && after == '/':
		reading.until("\n")
	case this == '/' && after == '*':
		reading.until("*/")
	case this == '"' || this == '\'':
		reading.quoted(this)
	case this == '_' || unicode.IsLetter(this):
		reading.name()
	case unicode.IsDigit(this) || this == '.' && unicode.IsDigit(after):
		reading.number()
	default:
		reading.advance()
	}
}

// until reads up to and past the closing of a comment.
func (reading *metalReading) until(end string) {
	closing := []rune(end)
	for !reading.done() && !reading.reads(closing) {
		reading.advance()
	}
	for range closing {
		reading.advance()
	}
}

// reads is whether the source reads closing from here.
func (reading *metalReading) reads(closing []rune) bool {
	for ahead, want := range closing {
		if reading.peek(ahead) != want {
			return false
		}
	}
	return true
}

// quoted reads past a string or a character, escapes and all.
func (reading *metalReading) quoted(mark rune) {
	reading.advance()
	for !reading.done() && reading.peek(here) != mark {
		if reading.peek(here) == '\\' {
			reading.advance()
		}
		reading.advance()
	}
	reading.advance()
}

// name reads past a name, digits in it and all.
func (reading *metalReading) name() {
	for !reading.done() && namePart(reading.peek(here)) {
		reading.advance()
	}
}

// number reads a number with its point, exponent and suffix. It keeps
// the number.
func (reading *metalReading) number() {
	start, line, column := reading.at, reading.line, reading.column
	var last rune
	for !reading.done() {
		this := reading.peek(here)
		signed := (this == '+' || this == '-') &&
			strings.ContainsRune("eEpP", last)
		if !namePart(this) && this != '.' && !signed {
			break
		}
		last = this
		reading.advance()
	}
	written := string(reading.text[start:reading.at])
	reading.found = append(reading.found,
		metalNumber{written, line, column})
}

// namePart is whether a character can be part of a name.
func namePart(this rune) bool {
	return this == '_' || unicode.IsLetter(this) || unicode.IsDigit(this)
}
