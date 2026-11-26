package main

import (
	"cmp"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-checks/internal/requirement"
)

// history is what past commits say about requirement IDs.
type history struct {
	// removed holds every ID a commit took out entirely.
	removed map[string]bool
	// homes maps each numbered acceptance, by its doc number and name, to the
	// ID it was added at.
	homes map[string]string
	// held holds every acceptance ID a commit added outside an HTML comment.
	// A numbered doc that drops one must keep it as a Removed line.
	held map[string]bool
}

// wholeDoc is the diff context that keeps each changed doc whole in one hunk,
// so a step line's acceptance is always above it.
const wholeDoc = "-U100000"

// historyLine matches an added or removed title or acceptance line in a diff.
var historyLine = regexp.MustCompile(`^([+-])(?:# (\d{4}) |- (\d{4}\.\d+)\s*$)`)

// readHistories replays one git log of every root's doc commits into each root's record, empty with a note when git fails.
func readHistories(root string, roots, paths []string) (map[string]history, string) {
	hists := map[string]history{}
	dirs := make([]string, len(roots))
	for i, r := range roots {
		hists[r] = history{removed: map[string]bool{}, homes: map[string]string{}, held: map[string]bool{}}
		dirs[i] = packages.Join(r, docsDir)
	}
	log, note := docLog(root, dirs)
	logs := splitLog(roots, log)
	for _, r := range roots {
		// A restructure the docs judged now keep none of restarts the record, as one more commit would.
		if h := hists[r]; h.replay(logs[r]).keepsNone(packages.Owned(roots, r, paths)) {
			h.restart()
		}
	}
	return hists, note
}

// docLog is the patch log of dirs oldest first, or empty with a note when git fails with a history to read.
// Its fixed prefixes, unquoted paths and paths relative to root keep each diff header matching a package root.
func docLog(root string, dirs []string) (log, note string) {
	out, err := exec.Command("git", append([]string{"-C", root, "-c", "core.quotePath=false", "log", "--reverse", "--no-color",
		"--no-ext-diff", "--no-renames", "--src-prefix=a/", "--dst-prefix=b/", "--relative", wholeDoc, "--format=%x00", "-p", "--"}, dirs...)...).Output()
	if err == nil {
		return string(out), ""
	}
	if why := gitFailure(err); !noHistory.MatchString(why) {
		return "", "requirement IDs were not checked against git history: " + why
	}
	return "", ""
}

// splitLog divides a history log among roots by the root that owns each file's diff. Every root's log keeps every
// commit, so each replays the same commits in order.
func splitLog(roots []string, log string) map[string]string {
	commits := map[string][]string{}
	for _, commit := range strings.Split(log, "\x00") {
		for r, part := range splitCommit(roots, commit) {
			commits[r] = append(commits[r], part)
		}
	}
	logs := map[string]string{}
	for r, parts := range commits {
		logs[r] = strings.Join(parts, "\x00")
	}
	return logs
}

// splitCommit divides one commit's patch among roots by the root that owns each file's diff, giving every root a part.
func splitCommit(roots []string, commit string) map[string]string {
	builders := map[string]*strings.Builder{}
	for _, r := range roots {
		builders[r] = &strings.Builder{}
	}
	owner := packages.Root
	for _, line := range strings.SplitAfter(commit, "\n") {
		if path, ok := diffPath(strings.TrimSuffix(line, "\n")); ok {
			owner = packages.Of(roots, path)
		}
		builders[owner].WriteString(line)
	}
	parts := map[string]string{}
	for r, b := range builders {
		parts[r] = b.String()
	}
	return parts
}

// replay folds each commit of a history log into h, clearing the record after
// each restructure, and returns the docs that remain.
func (h history) replay(log string) docFiles {
	live, docs := map[string]int{}, docFiles{}
	for _, commit := range strings.Split(log, "\x00") {
		restructure := docs.restructures(commit)
		net, names := commitNet(commit)
		h.apply(live, net)
		h.settle(names)
		h.keep(names["+"])
		h.keep(names[heldSide])
		if restructure {
			h.restart()
		}
	}
	return docs
}

// restart forgets the record, as a restructure does.
func (h history) restart() {
	clear(h.removed)
	clear(h.homes)
	clear(h.held)
}

// keep marks each ID that ids map to as held.
func (h history) keep(ids map[string]string) {
	for _, id := range ids {
		h.held[id] = true
	}
}

// settle folds one commit's acceptance names into h: a name only added takes
// its ID as home, and one only removed loses it. A name that moved keeps its
// home.
func (h history) settle(names map[string]map[string]string) {
	for key := range names["-"] {
		if movedFrom(names, key) == "" {
			delete(h.homes, key)
		}
	}
	for key, to := range names["+"] {
		if from := movedFrom(names, key); from != "" {
			h.homes[key] = cmp.Or(h.homes[key], from)
			continue
		}
		h.homes[key] = to
	}
}

// movedFrom is the ID the name at key moved from in one commit, or empty. A
// name moves when it is removed and added, unless its old ID became Removed.
func movedFrom(names map[string]map[string]string, key string) string {
	from, removed := names["-"][key]
	_, added := names["+"][key]
	if !removed || !added || names[removedSide][from] != "" {
		return ""
	}
	return from
}

// removedSide files the IDs a commit wrote Removed lines at, beside the
// names it removed and added.
const removedSide = "Removed"

// heldSide files the bullet acceptance IDs a commit added outside an HTML
// comment.
const heldSide = "Held"

// homeKey names a numbered acceptance across commits: its doc number and
// name.
func homeKey(doc, name string) string {
	return doc + " " + name
}

// keepsNone reports whether d holds docs and paths keep none of them, as an
// uncommitted restructure does.
func (d docFiles) keepsNone(paths []string) bool {
	return len(d) > 0 && !slices.ContainsFunc(paths, func(p string) bool { return d[p] })
}

// diffHeader starts one file's part of a diff.
const diffHeader = "diff --git "

// diffPath is the path a diff header names. Both halves name one file without rename detection, so the second half
// starts past the middle, quoted alike.
func diffPath(line string) (string, bool) {
	rest, ok := strings.CutPrefix(line, diffHeader)
	if !ok {
		return "", false
	}
	half := rest[len(rest)/2+1:]
	if unquoted, err := strconv.Unquote(half); err == nil {
		half = unquoted
	}
	return strings.TrimPrefix(half, "b/"), true
}

// docFiles holds the requirement doc paths that exist as history replays.
type docFiles map[string]bool

// restructures folds one commit's added and deleted docs into d and reports whether it removed every doc that existed before.
// History runs without rename detection, so a doc renamed away counts as removed.
func (d docFiles) restructures(diff string) bool {
	before, gone, path := len(d), 0, ""
	for _, line := range strings.Split(diff, "\n") {
		if p, ok := diffPath(line); ok {
			path = p
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
// edited or moved line nets to zero. It also maps each numbered acceptance
// the commit removed or added, by diff sign, to its ID there.
func commitNet(diff string) (map[string]int, map[string]map[string]string) {
	sc := scan{net: map[string]int{}, names: map[string]map[string]string{"+": {}, "-": {}, removedSide: {}, heldSide: {}}}
	for _, line := range strings.Split(diff, "\n") {
		sc.read(line)
	}
	return sc.net, sc.names
}

// scan is one commit's diff so far: each ID's net count, and where each
// numbered acceptance was removed or added.
type scan struct {
	net   map[string]int
	names map[string]map[string]string
	under sides
}

// read counts one diff line. A step's ID comes from the acceptance above it,
// on the side of the diff the step line is on.
func (sc *scan) read(line string) {
	if m := historyLine.FindStringSubmatch(line); m != nil {
		sc.net[m[2]+m[3]] += lineSign[m[1]]
	}
	sc.under.follow(line)
	sc.note(sc.under.numbered(line))
	sc.hold(line)
	if id := sc.under.step(line); id != "" {
		sc.net[id] += lineSign[line[:1]]
	}
}

// hold files a bullet acceptance ID line that a diff line adds to a
// Requirements section outside an HTML comment.
func (sc *scan) hold(line string) {
	if m := acceptanceLine.FindStringSubmatch(line); m != nil && sc.under.owns(m[2]) && sc.under.visible(m[1]) == "+" {
		sc.names[heldSide][m[2]] = m[2]
	}
}

// owns reports whether id is an acceptance of the requirement doc whose
// Requirements section the diff is in.
func (s sides) owns(id string) bool {
	return s.inReqs && s.doc != "" && strings.HasPrefix(id, s.doc+".")
}

// note counts a numbered acceptance line toward its ID and files its name.
// An added Removed line files the ID it deletes.
func (sc *scan) note(id, name, sign string) {
	if counts(name, sign) {
		sc.net[id] += lineSign[sign]
		sc.names[sign][homeKey(sc.under.doc, name)] = id
	}
	if requirement.IsRemoved(name) && sign == "+" {
		sc.names[removedSide][id] = id
	}
}

// counts reports whether a numbered acceptance line with diff sign counts
// toward its ID. A Removed line holds none, so writing one deletes the ID.
func counts(name, sign string) bool {
	return !requirement.IsRemoved(name) && (sign == "+" || sign == "-")
}

// lineSign is how an added or removed diff line counts toward its ID.
var lineSign = map[string]int{"+": 1, "-": -1}

// acceptanceLine matches an acceptance ID line on either side of a diff, or
// in its context.
var acceptanceLine = regexp.MustCompile(`^([ +-])- (\d{4}\.\d+)\s*$`)

// stepLine matches an added or removed When line: numbered step N, or a
// single step's plain bullet, which is step 1. A numbered acceptance writes
// it WHEN.
var stepLine = regexp.MustCompile(`^[+-]\s+(?:(\d+)\.|-) (?:When|WHEN) `)

// numberedLine matches a top-level numbered item on either side of a diff,
// or in its context. In a Requirements section it is a numbered acceptance.
var numberedLine = regexp.MustCompile(`^([ +-])(\d+)\. (.+)$`)

// sides is the acceptance each side of a diff is under: the old one for
// removed lines, the new one for added lines. It also holds the doc's
// number, and whether the diff is in its Requirements section.
type sides struct {
	old, new string
	doc      string
	inReqs   bool
	// hidden marks the old and the new side as inside an HTML comment.
	hidden [2]bool
}

// follow moves each side to the acceptance a diff line opens there. A new
// file starts under none.
func (s *sides) follow(line string) {
	if p, ok := diffPath(line); ok {
		*s = startFile(p)
		return
	}
	s.place(line)
	if id, sign := s.opened(line); id != "" {
		s.open(sign, id)
	}
}

// place follows the section a diff line is in, and the HTML comment each
// side it is on is in.
func (s *sides) place(line string) {
	text := line[min(1, len(line)):]
	if heading, ok := strings.CutPrefix(text, "## "); ok {
		s.inReqs = strings.TrimSpace(heading) == "Requirements"
	}
	if !strings.HasPrefix(line, "+") {
		s.hidden[0] = commentState(text, s.hidden[0])
	}
	if !strings.HasPrefix(line, "-") {
		s.hidden[1] = commentState(text, s.hidden[1])
	}
}

// commentState is whether text leaves its side inside an HTML comment, from
// was.
func commentState(text string, was bool) bool {
	if strings.Contains(text, "-->") && strings.LastIndex(text, "-->") > strings.LastIndex(text, "<!--") {
		return false
	}
	return was || strings.Contains(text, "<!--")
}

// visible is a line's diff sign once comments hide a side. A context line
// hidden on one side reads as added or removed, and one hidden on all as none.
func (s sides) visible(sign string) string {
	old := sign != "+" && !s.hidden[0]
	added := sign != "-" && !s.hidden[1]
	return map[[2]bool]string{{true, true}: " ", {false, true}: "+", {true, false}: "-"}[[2]bool{old, added}]
}

// startFile is the sides of a new file's diff, holding the number of a
// requirement doc.
func startFile(p string) sides {
	if m := numberedDoc.FindStringSubmatch(path.Base(p)); m != nil && isDoc(p) {
		return sides{doc: m[1]}
	}
	return sides{}
}

// opened is the acceptance ID a diff line opens, with its diff sign, or an
// empty ID.
func (s sides) opened(line string) (string, string) {
	if id, _, sign := s.numbered(line); id != "" {
		return id, sign
	}
	if m := acceptanceLine.FindStringSubmatch(line); m != nil {
		return m[2], m[1]
	}
	return "", ""
}

// numbered is the ID, name and diff sign of a numbered acceptance line, or
// all empty. Its ID is the doc's number and the item's.
func (s sides) numbered(line string) (id, name, sign string) {
	m := numberedLine.FindStringSubmatch(line)
	if m == nil || !s.inReqs || s.doc == "" || s.visible(m[1]) == "" {
		return "", "", ""
	}
	n, _ := strconv.Atoi(m[2])
	return fmt.Sprintf("%s.%d", s.doc, n), strings.TrimSpace(m[3]), s.visible(m[1])
}

// open moves the sides a line with diff sign is on to acceptance id.
func (s *sides) open(sign, id string) {
	if sign != "+" {
		s.old = id
	}
	if sign != "-" {
		s.new = id
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
