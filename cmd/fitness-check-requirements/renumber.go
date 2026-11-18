package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// renumberFile is the checked-in map from each moved ID to its new home.
const renumberFile = docsDir + "renumber.txt"

// renumberEntry matches one `old -> new` line: a doc number to a doc
// number, or an acceptance ID to an acceptance ID.
var renumberEntry = regexp.MustCompile(`^(\d{4})(\.\d+)?\s*->\s*(\d{4})(\.\d+)?$`)

// byteOrderMark is the UTF-8 mark some editors save at a file's start.
const byteOrderMark = "\xef\xbb\xbf"

// move is one map entry and the line that states it.
type move struct {
	from, to string
	line     int
}

// renumbering is what the map says moved: each old ID's entry.
type renumbering struct {
	moves []move
	from  map[string]bool
}

// readRenumbering reads and parses the map. A repo with no map moved
// nothing, and a map it cannot read fails as a finding.
func readRenumbering(root string) (renumbering, []string) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(renumberFile)))
	if errors.Is(err, fs.ErrNotExist) {
		data, err = nil, nil
	}
	if err != nil {
		return renumbering{from: map[string]bool{}}, []string{renumberFile + " cannot be read: " + err.Error()}
	}
	return parseRenumbering(strings.TrimPrefix(string(data), byteOrderMark))
}

// parseRenumbering reads each entry of the map's text, skipping blank lines
// and `#` comments, and returns an error for each line it rejects.
func parseRenumbering(text string) (renumbering, []string) {
	r := renumbering{from: map[string]bool{}}
	var errs []string
	for i, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if msg := r.add(line, i+1); msg != "" {
			errs = append(errs, requirement.At(renumberFile, i+1, msg))
		}
	}
	return r, errs
}

// add records the entry on line n. It returns why it rejects a malformed
// line or an old ID listed twice.
func (r *renumbering) add(line string, n int) string {
	m := renumberEntry.FindStringSubmatch(line)
	if m == nil || (m[2] == "") != (m[4] == "") {
		return "write each entry as `old -> new`, both doc numbers such as 0142 -> 0005 or both acceptance IDs such as 0142.3 -> 0005.2, with comments on their own `#` line"
	}
	from, to := m[1]+m[2], m[3]+m[4]
	if from == to {
		return from + " maps to itself; list only IDs that moved"
	}
	if r.from[from] {
		return from + " is listed twice; a moved ID has one new home"
	}
	r.from[from] = true
	r.moves = append(r.moves, move{from: from, to: to, line: n})
	return ""
}

// covered is 1 when the map lists id, or the doc an acceptance id belongs
// to, and 0 otherwise. One removal was then a move, so the number is free
// again until a later commit really deletes it.
func (r renumbering) covered(id string) int {
	if r.from[id] || (strings.Contains(id, ".") && r.from[id[:4]]) {
		return 1
	}
	return 0
}

// targetErrors fails each entry whose new ID no doc defines, so every moved
// requirement stays traceable to where it went.
func (r renumbering) targetErrors(docs []requirement.Document) []string {
	defined := map[string]bool{}
	for _, d := range docs {
		for _, id := range d.IDs() {
			defined[id.ID] = true
		}
	}
	var errs []string
	for _, m := range r.moves {
		if !defined[m.to] {
			errs = append(errs, requirement.At(renumberFile, m.line, m.from+" -> "+m.to+" points at "+m.to+", which no doc defines; point it at the ID's new home"))
		}
	}
	return errs
}
