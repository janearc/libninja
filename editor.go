package libninja

import (
	"bytes"
	"fmt"
	"go/token"
	"os"
	"path/filepath"

	"github.com/editorconfig/editorconfig-core-go/v2"
)

// editorconfigName is the file EditorConfig reads, in a directory and
// every directory above it.
const editorconfigName = ".editorconfig"

// editorFiles is the project settings an editor reads where they sit, so
// being in the project is all it takes to be the project's. Each is named
// as its editor's own documentation or source names it. Those were read
// on 18 October 2026.
//
// Left out are files that sit elsewhere and name the project by path: a
// .code-workspace or a .sublime-project. Also left out are settings
// carried in each file: modelines and emacs's local variables. Cursor is
// left out too. Its documentation names no settings file of its own.
var editorFiles = []string{
	".idea",                                       // the jetbrains editors
	".vscode/settings.json",                       // vs code
	".zed/settings.json",                          // zed
	".helix/config.toml", ".helix/languages.toml", // helix
	".dir-locals.el", ".dir-locals-2.el", // emacs
	".nvim.lua", ".nvimrc", ".exrc", // neovim
	".vimrc", "_vimrc", ".gvimrc", // vim
}

// editors is the scroll of the unmarked house: a module no editor is told
// how to keep. It looks from the module's root up to its repository's.
//
// Each .editorconfig counts when one of its sections covers the module's
// go files, as EditorConfig's own code matches them. Project settings of
// an editor count by being there.
func (read reading) editors() []finding {
	rooted := false
	for _, dir := range upToRepository(read.dir) {
		if present(dir, editorFiles) {
			return nil
		}
		if rooted {
			continue
		}
		covers, root := read.editorconfigCovers(dir)
		if covers {
			return nil
		}
		rooted = root
	}
	return []finding{{editorless, read.address(read.module),
		read.atGoMod(), fmt.Sprintf("No editor is told how to keep "+
			"%s: no .editorconfig covers its go files, and no "+
			"editor's project settings are in it.",
			read.address(read.module)), ""}}
}

// editorconfigCovers is whether an .editorconfig in dir has a section
// that covers one of the module's go files, as EditorConfig's code
// matches its globs. It is also whether the file is a root. EditorConfig
// stops looking past a root.
func (read reading) editorconfigCovers(dir string) (covers, root bool) {
	// what cannot be read is read as empty. It covers nothing and is no
	// root.
	body, _ := os.ReadFile(filepath.Join(dir, editorconfigName))
	config, err := editorconfig.Parse(bytes.NewReader(body))
	if err != nil {
		return false, false
	}
	for _, path := range read.goFiles() {
		// dir is the module's root or above it, so every go file of the
		// module is below dir and Rel cannot refuse
		rel, _ := filepath.Rel(dir, path)
		said, err := config.GetDefinitionForFilename(
			filepath.ToSlash(rel))
		if err == nil && positive(len(said.Raw)) {
			return true, config.Root
		}
	}
	return false, config.Root
}

// goFiles is every go file of the module's packages.
func (read reading) goFiles() []string {
	var files []string
	for _, pkg := range read.pkgs {
		files = append(files, pkg.GoFiles...)
	}
	return files
}

// present is whether any of the files is in dir.
func present(dir string, files []string) bool {
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			return true
		}
	}
	return false
}

// upToRepository is dir and each directory above it up to the one that
// holds .git. It is dir alone when none does, since a project's settings
// are its own and not those of the disk above it. The walk goes on while
// the directory changes, so it ends at the root.
func upToRepository(dir string) []string {
	var dirs []string
	for at, last := dir, ""; at != last; at, last = filepath.Dir(at), at {
		dirs = append(dirs, at)
		if _, err := os.Stat(filepath.Join(at, ".git")); err == nil {
			return dirs
		}
	}
	return []string{dir}
}

// atGoMod is the module's go.mod, where a finding about the whole module
// is placed.
func (read reading) atGoMod() token.Position {
	return token.Position{Filename: filepath.Join(read.dir, "go.mod"),
		Line: lineOrigin, Column: lineOrigin}
}
