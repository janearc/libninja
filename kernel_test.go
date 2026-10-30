package libninja_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/janearc/libninja"
)

// kernelIn is what the kernel scroll saw in a house holding a kernel.
// The kernel is gpu/trace.metal with the given source.
func kernelIn(test *testing.T, source string) []string {
	test.Helper()
	var seen []string
	for _, each := range kernelHits(test, source) {
		seen = append(seen, each["seen"].(string))
	}
	return seen
}

// kernelHits is the kernel scroll's hits in a house holding a kernel.
// The house is loaded first. Its throw then has a deadline. Readings
// that never end fail the test there, instead of hanging it. The throw's
// goroutine runs on until the test binary exits.
func kernelHits(test *testing.T, source string) []map[string]any {
	test.Helper()
	dir := house(test, ".", map[string]string{
		".editorconfig":       "root = true\n\n[*.go]\nindent_style = tab\n",
		"gpu/trace.metal":     source,
		"testdata/skip.metal": "kernel void skip() { int never = 7; }\n",
	})
	test.Setenv("GOWORK", "off")
	module, err := libninja.Load(dir)
	if err != nil {
		test.Fatal(err)
	}
	each := throwing{module, make(chan libninja.Report, heldOne)}
	go each.throw()
	deadline := time.Duration(heldTo("throw-seconds")) * time.Second
	select {
	case report := <-each.done:
		return kernelOnly(factsOf(test, report)["hits"])
	case <-time.After(deadline):
		test.Fatalf("the throw of a kernel did not end in %v", deadline)
	}
	return nil
}

// heldOne is how many reports a throw hands on, defined so: one.
const heldOne = 1

// throwing is a module thrown on a goroutine of its own, so that a test
// can give the throw a deadline.
type throwing struct {
	module *libninja.Module
	done   chan libninja.Report
}

// throw throws the module and hands its report on.
func (each throwing) throw() { each.done <- each.module.Throw() }

// kernelOnly is the hits of the kernel scroll among hits.
func kernelOnly(hits []map[string]any) []map[string]any {
	var kernel []map[string]any
	for _, each := range hits {
		if each["scroll"] == "kernel" {
			kernel = append(kernel, each)
		}
	}
	return kernel
}

// Each number on a later line is placed by its own line and column.
// Lines and columns count from one. Tabs are one column each.
func TestAKernelNumberIsPlacedByItsLineAndColumn(test *testing.T) {
	source := "kernel void trace() {\n\tint a = 7;\n  int b = 8;\n}\n"
	var places []string
	for _, each := range kernelHits(test, source) {
		places = append(places, fmt.Sprintf("%v:%v:%v %v", each["file"],
			each["line"], each["column"], each["seen"]))
	}
	want := "gpu/trace.metal:2:10 7 is the kernel's own number.|" +
		"gpu/trace.metal:3:11 8 is the kernel's own number."
	if said := strings.Join(places, "|"); said != want {
		test.Errorf("where a kernel's numbers are:\n got %s\nwant %s",
			said, want)
	}
}

// The sign of an exponent stays in its number. Between two numbers, the
// minus parts them.
func TestAnExponentsSignStaysInItsNumber(test *testing.T) {
	source := "constant float far = 1e+5f;\nconstant int less = 2-1;\n"
	seen := strings.Join(kernelIn(test, source), "|")
	want := "1e+5f is the kernel's own number.|" +
		"2 is the kernel's own number.|1 is the kernel's own number."
	if seen != want {
		test.Errorf("signs in a kernel:\n got %s\nwant %s", seen, want)
	}
}

func TestEveryNumberInAKernelIsItsOwn(test *testing.T) {
	source := "#include <metal_stdlib>\n" +
		"kernel void trace(device float3 *light [[buffer(0)]]) {\n" +
		"\tfloat edge = 1e-7f;\n" +
		"\tfloat half = .5;\n" +
		"\tlight[0] = float3(0x1F, 2.0h, edge * half);\n}\n"
	seen := strings.Join(kernelIn(test, source), "|")
	want := "0 is the kernel's own number.|" +
		"1e-7f is the kernel's own number.|" +
		".5 is the kernel's own number.|" +
		"0 is the kernel's own number.|" +
		"0x1F is the kernel's own number.|" +
		"2.0h is the kernel's own number."
	if seen != want {
		test.Errorf("a kernel's numbers:\n got %s\nwant %s", seen, want)
	}
}

func TestCommentsStringsAndNamesHoldNoKernelNumbers(test *testing.T) {
	source := "// 42 in a comment\n/* 7 and\n 8 */\n" +
		"constant char *said = \"3 4\";\n" +
		"kernel void trace(device float4x4 *m2, uint3 at) {}\n"
	if seen := kernelIn(test, source); seen != nil {
		test.Errorf("comments, strings and names: %v", seen)
	}
}

func TestAnEscapedQuoteStaysInTheString(test *testing.T) {
	source := "constant char *said = \"3\\\" 4\";\n" +
		"constant char quote = '\\'';\n"
	if seen := kernelIn(test, source); seen != nil {
		test.Errorf("a string with an escaped quote in it: %v", seen)
	}
}

// superuser is root's user id, as unix has defined it: the one user a
// shut directory does not keep out.
const superuser = 0

// opening is a directory to open again once a test is done with it, so
// the test's temporary directory can be taken away.
type opening struct {
	path string
	mode os.FileMode
}

// open sets the directory's mode back.
func (again opening) open() { os.Chmod(again.path, again.mode) }

// go leaves out directories named with _, but they are still walked for
// kernels. One that is shut is passed by, and the rest is read.
func TestADirectoryTheWalkCannotReadIsPassedBy(test *testing.T) {
	if os.Getuid() == superuser {
		test.Skip("root reads a shut directory")
	}
	dir := house(test, ".", map[string]string{
		".editorconfig":      "root = true\n\n[*.go]\nindent_style = tab\n",
		"gpu/trace.metal":    "kernel void trace() { int seen = 7; }\n",
		"_shut/hidden.metal": "kernel void hidden() { int unseen = 8; }\n",
	})
	shut := opening{filepath.Join(dir, "_shut"),
		os.FileMode(fixture(test, "open-dir"))}
	shutMode := os.FileMode(fixture(test, "shut-dir"))
	if err := os.Chmod(shut.path, shutMode); err != nil {
		test.Fatal(err)
	}
	test.Cleanup(shut.open)
	_, facts := thrown(test, dir)
	if !holds(facts["hits"], "kernel", "gpu/trace.metal", "7 is") {
		test.Error("a kernel beside a shut directory was not read")
	}
}
