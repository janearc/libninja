package libninja

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

// lineWidth is how wide the master's prose is wrapped, picked.
var lineWidth = picked("line-width")

// Speak writes the report in the master's voice. Each package that broke
// a scroll is named. Under it comes each scroll it broke, with its law
// and every place. The way and the film's line close each scroll. The
// vows come next, and the tally last.
func (report Report) Speak(out io.Writer) error {
	var text strings.Builder
	students := report.students()
	fmt.Fprintf(&text, "shuriken reads %s: %d hits in %d packages, "+
		"%s.\n\n", report.module, len(report.hits),
		len(students), report.carrying())
	for _, student := range students {
		report.lesson(&text, student)
	}
	report.spokenVows(&text)
	report.tally(&text)
	_, err := io.WriteString(out, text.String())
	return err
}

// students is every package with a hit, in the order they are taught.
func (report Report) students() []string {
	var found []string
	for _, each := range report.hits {
		if !slices.Contains(found, each.student) {
			found = append(found, each.student)
		}
	}
	return found
}

// lesson is what the master says to one package.
func (report Report) lesson(text *strings.Builder, student string) {
	fmt.Fprintf(text, "%s, stay here. Sit down.\n", student)
	var current scroll
	for _, each := range report.hits {
		if each.student != student {
			continue
		}
		if each.scroll != current {
			report.close(text, current)
			current = each.scroll
			fmt.Fprintf(text, "\n  %s (%s)\n", current.title,
				current.name)
			text.WriteString(wrap(current.law, "  "))
			text.WriteString("\n")
		}
		fmt.Fprintf(text, "    %s:%d  %s\n", report.relative(each.at),
			each.at.Line, each.seen)
	}
	report.close(text, current)
	text.WriteString("\nUnderstand?\n\n")
}

// close ends a scroll's part of a lesson with the way and the film.
func (report Report) close(text *strings.Builder, done scroll) {
	if done.name == "" {
		return
	}
	text.WriteString("\n")
	text.WriteString(wrap(done.way, "  "))
	if done.filmed != "" {
		text.WriteString(wrap(done.filmed, "  "))
	}
}

// spokenVows is every vow kept, with its reason.
func (report Report) spokenVows(text *strings.Builder) {
	if none(len(report.vows)) {
		return
	}
	text.WriteString("the vows\n\n")
	for _, each := range report.vows {
		fmt.Fprintf(text, "  %s:%d  %s: %s\n", report.relative(each.at),
			each.at.Line, each.scroll.name, each.vowed)
	}
	text.WriteString("\n")
}

// tally is how many times each scroll was broken, and the last word.
func (report Report) tally(text *strings.Builder) {
	for _, each := range scrolls() {
		var count int
		for _, hit := range report.hits {
			if hit.scroll == each {
				count++
			}
		}
		if positive(count) {
			fmt.Fprintf(text, "  %-36s %5d\n", each.title, count)
		}
	}
	if report.Clean() {
		text.WriteString(
			"A ninja must always be calm. Now, go practice.\n")
		return
	}
	text.WriteString("\nKeep in mind, you must always be a good ninja.\n")
}

// wrap is prose wrapped to lineWidth, each line indented.
func wrap(prose, indent string) string {
	const space = " "
	var lines []string
	var line string
	for word := range strings.FieldsSeq(prose) {
		if line != "" &&
			len(indent)+len(line)+len(space)+len(word) > lineWidth {
			lines = append(lines, indent+line)
			line = ""
		}
		if line != "" {
			line += space
		}
		line += word
	}
	return strings.Join(append(lines, indent+line), "\n") + "\n"
}

// carrying is how many vows the module carries, and how many to each
// scroll in the order taught: "3 vows (well 2, nest 1)". Nothing to
// count says "0 vows". Every vow written is counted, whether it covers a
// hit or not. No count of hits is read clean while vows stand under it.
func (report Report) carrying() string {
	said := fmt.Sprintf("%d vows", len(report.carried))
	var counts []string
	for start := nothing; start < len(report.carried); {
		name := report.carried[start]
		end := start
		for end < len(report.carried) && report.carried[end] == name {
			end++
		}
		scroll := vowedScroll(name)
		counts = append(counts, fmt.Sprintf("%s %d", scroll, end-start))
		start = end
	}
	if none(len(counts)) {
		return said
	}
	return said + " (" + strings.Join(counts, ", ") + ")"
}

// unnamed is what the master calls a vow that names no scroll at all.
const unnamed = "unnamed"

// vowedScroll is a vow's scroll as the master says it.
func vowedScroll(name string) string {
	if name == "" {
		return unnamed
	}
	return name
}
