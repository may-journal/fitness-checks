package main

import (
	"cmp"
	"errors"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

// history is what past commits say about requirement IDs.
type history struct {
	// removed holds every ID a commit took out entirely.
	removed map[string]bool
}

// wholeDoc is the diff context that keeps each changed doc whole in one hunk,
// so a step line's acceptance is always above it.
const wholeDoc = "-U100000"

// historyLine matches an added or removed title or acceptance line in a diff.
var historyLine = regexp.MustCompile(`^([+-])(?:# (\d{4}) |- (\d{4}\.\d+)\s*$)`)

// readHistory replays the commits to the docs in dir oldest first, then paths, the docs judged now, as one more commit.
// Each restructure restarts the record, which is empty with a note when git fails.
func readHistory(root, dir string, paths []string) (h history, note string) {
	h = history{removed: map[string]bool{}}
	out, err := exec.Command("git", "-C", root, "log", "--reverse", "--no-color", "--no-ext-diff",
		"--no-renames", wholeDoc, "--format=%x00", "-p", "--", dir).Output()
	if err != nil {
		if why := gitFailure(err); !noHistory.MatchString(why) {
			return h, "requirement IDs were not checked against git history: " + why
		}
		return h, ""
	}
	if h.replay(string(out)).keepsNone(paths) {
		clear(h.removed)
	}
	return h, ""
}

// replay folds each commit of a history log into h, clearing the record after
// each restructure, and returns the docs that remain.
func (h history) replay(log string) docFiles {
	live, docs := map[string]int{}, docFiles{}
	for _, commit := range strings.Split(log, "\x00") {
		restructure := docs.restructures(commit)
		h.apply(live, commitNet(commit))
		if restructure {
			clear(h.removed)
		}
	}
	return docs
}

// keepsNone reports whether d holds docs and paths keep none of them, as an
// uncommitted restructure does.
func (d docFiles) keepsNone(paths []string) bool {
	return len(d) > 0 && !slices.ContainsFunc(paths, func(p string) bool { return d[p] })
}

// diffFile matches the header that starts one file's part of a diff, capturing its path.
var diffFile = regexp.MustCompile(`^diff --git a/.* b/(.*)$`)

// docFiles holds the requirement doc paths that exist as history replays.
type docFiles map[string]bool

// restructures folds one commit's added and deleted docs into d and reports whether it removed every doc that existed before.
// History runs without rename detection, so a doc renamed away counts as removed.
func (d docFiles) restructures(diff string) bool {
	before, gone, path := len(d), 0, ""
	for _, line := range strings.Split(diff, "\n") {
		if m := diffFile.FindStringSubmatch(line); m != nil {
			path = m[1]
		} else {
			gone += d.fold(path, line)
		}
	}
	return before > 0 && gone == before
}

// fold applies one line of path's diff header to d, and returns 1 when it deletes a doc.
func (d docFiles) fold(path, line string) int {
	switch {
	case strings.HasPrefix(line, "new file mode") && isDoc(path):
		d[path] = true
	case strings.HasPrefix(line, "deleted file mode") && d[path]:
		delete(d, path)
		return 1
	}
	return 0
}

// isDoc reports whether a Markdown path under docsDir is a requirement doc
// rather than the template.
func isDoc(path string) bool {
	return strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, "/template.md")
}

// noHistory matches git's complaint outside a repo or before the first commit,
// where there is no history to read.
var noHistory = regexp.MustCompile(`not a git repository|does not have any commits`)

// gitFailure names why a git command failed: its error text, else the Go
// error.
func gitFailure(err error) string {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if msg := strings.Join(strings.Fields(string(exitErr.Stderr)), " "); msg != "" {
			return msg
		}
	}
	return err.Error()
}

// commitNet counts each ID's added minus removed lines in one commit, so an
// edited or moved line nets to zero. A step's ID comes from the acceptance
// above it, on the side of the diff the step line is on.
func commitNet(diff string) map[string]int {
	net := map[string]int{}
	var under sides
	for _, line := range strings.Split(diff, "\n") {
		if m := historyLine.FindStringSubmatch(line); m != nil {
			net[m[2]+m[3]] += lineSign[m[1]]
		}
		under.follow(line)
		if id := under.step(line); id != "" {
			net[id] += lineSign[line[:1]]
		}
	}
	return net
}

// lineSign is how an added or removed diff line counts toward its ID.
var lineSign = map[string]int{"+": 1, "-": -1}

// acceptanceLine matches an acceptance ID line on either side of a diff, or
// in its context.
var acceptanceLine = regexp.MustCompile(`^([ +-])- (\d{4}\.\d+)\s*$`)

// stepLine matches an added or removed When line: numbered step N, or a
// single step's plain bullet, which is step 1.
var stepLine = regexp.MustCompile(`^[+-]\s+(?:(\d+)\.|-) When `)

// sides is the acceptance each side of a diff is under: the old one for
// removed lines, the new one for added lines.
type sides struct{ old, new string }

// follow moves each side to the acceptance a diff line opens there. A new
// file starts under none.
func (s *sides) follow(line string) {
	if strings.HasPrefix(line, "diff --git ") {
		*s = sides{}
	}
	m := acceptanceLine.FindStringSubmatch(line)
	if m == nil {
		return
	}
	if m[1] != "+" {
		s.old = m[2]
	}
	if m[1] != "-" {
		s.new = m[2]
	}
}

// step is the ID of the step an added or removed When line holds, or "".
func (s sides) step(line string) string {
	m := stepLine.FindStringSubmatch(line)
	under := map[string]string{"+": s.new, "-": s.old}[line[:min(1, len(line))]]
	if m == nil || under == "" {
		return ""
	}
	return under + "." + cmp.Or(m[1], "1")
}

// apply folds one commit's net changes into the live counts, recording IDs
// whose last copy disappears.
func (h history) apply(live, net map[string]int) {
	for id, n := range net {
		before := live[id]
		live[id] += n
		if before > 0 && live[id] <= 0 {
			h.removed[id] = true
		}
	}
}
