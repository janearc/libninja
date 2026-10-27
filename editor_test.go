package libninja_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/janearc/puzzlemouth/numbers"
)

// fixture is a number the tests are held to, by its name in
// testdata/numbers.yaml.
func fixture(test *testing.T, name string) float64 {
	test.Helper()
	raw, err := os.ReadFile("testdata/numbers.yaml")
	if err != nil {
		test.Fatal(err)
	}
	held, err := numbers.Read("testdata/numbers.yaml", raw)
	if err != nil {
		test.Fatal(err)
	}
	number, err := held.Get(name)
	if err != nil {
		test.Fatal(err)
	}
	return number.Value()
}

// house is a module written into a fresh directory with its files, each
// path to its body. It has a .git at top so the house is its own
// repository. The module itself sits at module, below top. It is the
// directory the module is in.
func house(test *testing.T, module string, files map[string]string) string {
	test.Helper()
	top := test.TempDir()
	files[".git/HEAD"] = "ref: refs/heads/main\n"
	files[filepath.Join(module, "go.mod")] = "module example.com/house\n\n" +
		"go 1.24\n"
	files[filepath.Join(module, "house.go")] = "// Package house is kept.\n" +
		"package house\n"
	openDir := os.FileMode(fixture(test, "open-dir"))
	openFile := os.FileMode(fixture(test, "open-file"))
	for path, body := range files {
		at := filepath.Join(top, path)
		if err := os.MkdirAll(filepath.Dir(at), openDir); err != nil {
			test.Fatal(err)
		}
		if err := os.WriteFile(at, []byte(body), openFile); err != nil {
			test.Fatal(err)
		}
	}
	return filepath.Join(top, module)
}

// unmarked is whether the editor scroll found the house unmarked.
func unmarked(test *testing.T, dir string) bool {
	test.Helper()
	_, facts := thrown(test, dir)
	return holds(facts["hits"], "editor", "go.mod", "No editor is told")
}

func TestAHouseNoEditorIsToldOfIsUnmarked(test *testing.T) {
	if !unmarked(test, house(test, ".", map[string]string{})) {
		test.Error("a module with no editor settings was not found")
	}
}

func TestAnEditorconfigCoveringItsGoMarksIt(test *testing.T) {
	dir := house(test, ".", map[string]string{
		".editorconfig": "root = true\n\n[*.go]\nindent_style = tab\n",
	})
	if unmarked(test, dir) {
		test.Error("an .editorconfig for *.go was not seen")
	}
}

func TestAnEditorconfigForOtherFilesLeavesItUnmarked(test *testing.T) {
	dir := house(test, ".", map[string]string{
		".editorconfig": "root = true\n\n[*.md]\nindent_size = 2\n",
	})
	if !unmarked(test, dir) {
		test.Error("an .editorconfig for markdown alone was taken")
	}
}

func TestAnEditorsProjectSettingsMarkIt(test *testing.T) {
	dir := house(test, ".", map[string]string{
		".vscode/settings.json": "{}\n",
	})
	if unmarked(test, dir) {
		test.Error(".vscode/settings.json was not seen")
	}
}

func TestSettingsAboveTheRepositoryAreNotItsOwn(test *testing.T) {
	dir := house(test, "inner", map[string]string{
		"inner/.git/HEAD": "ref: refs/heads/main\n",
		".editorconfig":   "[*.go]\nindent_style = tab\n",
	})
	if !unmarked(test, dir) {
		test.Error("an .editorconfig outside the repository was taken")
	}
}

func TestARootEditorconfigHidesTheOnesAbove(test *testing.T) {
	dir := house(test, "inner", map[string]string{
		".editorconfig":       "[*.go]\nindent_style = tab\n",
		"inner/.editorconfig": "root = true\n\n[*.md]\nindent_size = 2\n",
	})
	if !unmarked(test, dir) {
		test.Error("an .editorconfig above a root one was taken")
	}
}

func TestAnEditorconfigThatCannotBeReadLeavesItUnmarked(test *testing.T) {
	dir := house(test, ".", map[string]string{
		".editorconfig": "[*.go\nindent_style = tab\n",
	})
	if !unmarked(test, dir) {
		test.Error("an .editorconfig that does not parse was taken")
	}
}

func TestWithNoRepositoryOnlyTheModuleIsItsOwn(test *testing.T) {
	dir := house(test, "inner", map[string]string{
		".editorconfig": "[*.go]\nindent_style = tab\n",
	})
	if err := os.RemoveAll(filepath.Join(dir, "..", ".git")); err != nil {
		test.Fatal(err)
	}
	if !unmarked(test, dir) {
		test.Error("with no repository, an .editorconfig above the " +
			"module was taken")
	}
}

func TestAnEditorconfigAboveTheModuleInTheRepositoryMarksIt(
	test *testing.T) {
	dir := house(test, "inner", map[string]string{
		".editorconfig": "root = true\n\n[*.go]\nindent_style = tab\n",
	})
	if unmarked(test, dir) {
		test.Error("the repository's .editorconfig was not seen")
	}
}
