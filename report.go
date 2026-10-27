package libninja

import (
	"cmp"
	"encoding/json"
	"go/token"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

// finding is one place a scroll is broken. It holds the scroll and the
// package the master addresses. It holds where, and what he saw. When a
// vow covers it, it holds the vow's reason too.
type finding struct {
	scroll  scroll
	student string
	at      token.Position
	seen    string
	vowed   string
}

// Report is what the scrolls found: the hits, and the vows that cover
// the rest.
type Report struct {
	dir, module string
	hits, vows  []finding
	carried     []string
}

// Throw reads the module with every scroll.
func (module Module) Throw() Report {
	read := module.read
	var found []finding
	for _, throw := range []func() []finding{read.openGates,
		read.strangers, read.wells, read.idle, read.dojos, read.decoys,
		read.blades, read.compositions, read.stances, read.names,
		read.ages, read.buriedPaths, read.editors, read.nests,
		read.literals,
		read.indexes, read.numbers, read.tables, read.kernels,
		read.sourcelessNumbers, read.emptyVows, read.crowd} {
		found = append(found, throw()...)
	}
	report := Report{dir: read.dir, module: read.module,
		carried: read.carriedScrolls()}
	used := map[string]bool{}
	for _, each := range found {
		if written, key, ok := read.vowFor(each); ok {
			used[key] = true
			each.vowed = written.reason
			report.vows = append(report.vows, each)
		} else {
			report.hits = append(report.hits, each)
		}
	}
	report.hits = append(report.hits, read.idleVows(used)...)
	order := teaching{read.dir}
	slices.SortFunc(report.hits, order.taught)
	slices.SortFunc(report.vows, order.taught)
	return report
}

// teaching is the order findings are spoken in, for a module in dir.
type teaching struct{ dir string }

// taught orders two findings. They go by package first. Within a
// package they go by scroll, and within a scroll by place. Each file is
// taken as the report says it, from the module's root. So the order is
// the same wherever the module sits.
func (order teaching) taught(one, other finding) int {
	first, then := one.scroll.taught(), other.scroll.taught()
	file := relativeTo(order.dir, one.at)
	otherFile := relativeTo(order.dir, other.at)
	return cmp.Or(cmp.Compare(one.student, other.student),
		cmp.Compare(first, then), cmp.Compare(file, otherFile),
		cmp.Compare(one.at.Line, other.at.Line),
		cmp.Compare(one.at.Column, other.at.Column))
}

// Clean is whether no scroll is broken without a vow.
func (report Report) Clean() bool {
	return none(len(report.hits))
}

// JSON writes the findings as facts, with none of the master's words.
// It writes the module, then the hits and the vows. Each holds its
// scroll and its package. Each holds its place relative to the module,
// and what was seen.
func (report Report) JSON(out io.Writer) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(map[string]any{"module": report.module,
		"hits": report.listed(report.hits),
		"vows": report.listed(report.vows)})
}

// listed is findings as plain values for JSON.
func (report Report) listed(found []finding) []map[string]any {
	listed := []map[string]any{}
	for _, each := range found {
		fact := map[string]any{"scroll": each.scroll.name,
			"package": each.student,
			"file":    report.relative(each.at),
			"line":    each.at.Line, "column": each.at.Column,
			"seen": each.seen}
		if each.vowed != "" {
			fact["vow"] = each.vowed
		}
		listed = append(listed, fact)
	}
	return listed
}

// relative is a place's file as the module's own path to it.
func (report Report) relative(at token.Position) string {
	return relativeTo(report.dir, at)
}

// relativeTo is a place's file as the path to it from dir. When a //line
// directive puts the place outside dir, the path is whole.
func relativeTo(dir string, at token.Position) string {
	return strings.TrimPrefix(at.Filename, dir+string(filepath.Separator))
}
