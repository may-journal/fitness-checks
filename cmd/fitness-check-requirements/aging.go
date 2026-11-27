package main

import (
	"cmp"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-runner/go/pkg/gitx"
)

// agedReleases is how many released versions must be newer than a Removed line before it may go.
const agedReleases = 2

// changelogDir holds the entry files whose version headings are the released versions.
const changelogDir = "changelog/"

// versionHeading matches an entry file's version heading.
var versionHeading = regexp.MustCompile(`(?m)^### (\d{4}\.\d{2}\.\d{2}\.\d{4})\s*$`)

// removedVersion matches the version a Removed name holds.
var removedVersion = regexp.MustCompile(`^Removed in (\d{4}\.\d{2}\.\d{2}\.\d{4}):`)

// drop is a Removed line that a change deleted, with the ID it held.
type drop struct{ id, version string }

// entry is one top-level numbered item of a Requirements section.
type entry struct {
	n       int
	name    string
	removed bool
}

// docChange is one numbered doc's items before and after a diff, with the step IDs after it.
type docChange struct {
	doc      string
	old, now []entry
	steps    []string
}

// whenLine matches a step line on the new side of a diff, or in its context.
var whenLine = regexp.MustCompile(`^[ +]\s+(?:(\d+)\.|-) (?:When|WHEN) `)

// changeReader reads a diff line by line into the changes of each numbered doc.
type changeReader struct {
	s   sides
	out []docChange
}

// docChanges reads each numbered doc's items from a diff.
func docChanges(diff string) []docChange {
	var r changeReader
	for _, line := range strings.Split(diff, "\n") {
		r.read(line)
	}
	return r.out
}

// read folds one diff line into the doc it is in.
func (r *changeReader) read(line string) {
	_, header := diffPath(line)
	r.s.follow(line)
	if header && r.s.doc != "" {
		r.out = append(r.out, docChange{doc: r.s.doc})
	}
	if c := r.current(); c != nil {
		c.add(r.s, line)
	}
}

// current is the change of the doc the reader is in, or nil.
func (r *changeReader) current() *docChange {
	if r.s.doc == "" || len(r.out) == 0 {
		return nil
	}
	return &r.out[len(r.out)-1]
}

// add files a numbered item or a step that line holds.
func (c *docChange) add(s sides, line string) {
	if id, name, sign := s.numbered(line); id != "" {
		c.item(id, name, sign)
		return
	}
	if id := newStep(s, line); id != "" {
		c.steps = append(c.steps, id)
	}
}

// item files an item on the side or sides its sign names.
func (c *docChange) item(id, name, sign string) {
	n, _ := strconv.Atoi(strings.TrimPrefix(id, c.doc+"."))
	e := entry{n: n, name: name, removed: requirement.IsRemoved(name)}
	if sign != "+" {
		c.old = append(c.old, e)
	}
	if sign != "-" {
		c.now = append(c.now, e)
	}
}

// newStep is the ID of the step a line holds on the new side, or empty.
func newStep(s sides, line string) string {
	m := whenLine.FindStringSubmatch(line)
	if m == nil || !s.inReqs || s.new == "" || !newSide(s.visible(line[:1])) {
		return ""
	}
	return s.new + "." + cmp.Or(m[1], "1")
}

// newSide reports whether a line with visible sign v is on the new side.
func newSide(v string) bool {
	return v == " " || v == "+"
}

// markers keeps the Removed items of es.
func markers(es []entry) []entry {
	return slices.DeleteFunc(slices.Clone(es), func(e entry) bool { return !e.removed })
}

// matchIn is the index in left of the line e became, by its text or else its number, or -1.
func matchIn(left []entry, e entry) int {
	if i := slices.IndexFunc(left, func(l entry) bool { return l.name == e.name }); i >= 0 {
		return i
	}
	return slices.IndexFunc(left, func(l entry) bool { return l.n == e.n })
}

// unmatchedMarkers is the Removed lines of c that its new side no longer holds.
func unmatchedMarkers(c docChange) []entry {
	left := markers(c.now)
	var gone []entry
	for _, e := range markers(c.old) {
		if i := matchIn(left, e); i >= 0 {
			left = slices.Delete(left, i, i+1)
			continue
		}
		gone = append(gone, e)
	}
	return gone
}

// dropsOf is the Removed lines c deletes while later items move up, and the lowest number among them.
// A change that keeps its item count is no drop, as when a Removed line returns as an acceptance.
func dropsOf(c docChange) (drops []entry, from int) {
	drops = unmatchedMarkers(c)
	if len(drops) == 0 || len(c.now) != len(c.old)-len(drops) {
		return nil, 0
	}
	from = drops[0].n
	for _, d := range drops {
		from = min(from, d.n)
	}
	return drops, from
}

// compact folds the dropped Removed lines of a diff into h.
// Items from the lowest dropped number on take the numbers they hold now.
func (h history) compact(live map[string]int, diff string) {
	for _, c := range docChanges(diff) {
		if drops, from := dropsOf(c); len(drops) > 0 {
			h.note(c.doc, drops)
			h.renumber(live, c, from)
		}
	}
}

// note records the dropped lines of doc.
func (h history) note(doc string, drops []entry) {
	for _, d := range drops {
		h.drops[drop{fmt.Sprintf("%s.%d", doc, d.n), versionOf(d.name)}] = true
	}
}

// versionOf is the version a Removed name holds, or empty.
func versionOf(name string) string {
	m := removedVersion.FindStringSubmatch(name)
	if m == nil {
		return ""
	}
	return m[1]
}

// pastNumber reports whether id is doc's ID, or a step of it, numbered from or higher.
func pastNumber(doc string, from int, id string) bool {
	rest, ok := strings.CutPrefix(id, doc+".")
	num, _, _ := strings.Cut(rest, ".")
	n, _ := strconv.Atoi(num)
	return ok && n >= from
}

// renumber forgets what h holds of c's doc from number from on, and records what c holds there.
func (h history) renumber(live map[string]int, c docChange, from int) {
	past := func(id string) bool { return pastNumber(c.doc, from, id) }
	maps.DeleteFunc(h.removed, func(id string, _ bool) bool { return past(id) })
	maps.DeleteFunc(h.held, func(id string, _ bool) bool { return past(id) })
	maps.DeleteFunc(live, func(id string, _ int) bool { return past(id) })
	maps.DeleteFunc(h.homes, func(key, home string) bool { return strings.HasPrefix(key, c.doc+" ") && past(home) })
	h.record(live, c, from)
	for _, id := range c.steps {
		if past(id) {
			live[id]++
		}
	}
}

// record files the items of c numbered from or higher as new.
func (h history) record(live map[string]int, c docChange, from int) {
	for _, e := range c.now {
		id := fmt.Sprintf("%s.%d", c.doc, e.n)
		if e.n < from {
			continue
		}
		h.held[id] = true
		if e.removed {
			h.removed[id] = true
			continue
		}
		h.homes[homeKey(c.doc, e.name)] = id
		live[id] = 1
	}
}

// agingErrors flags each Removed line deleted before two released versions are newer than it.
func agingErrors(root string, docs []requirement.Document, hist history) []string {
	if len(hist.drops) == 0 {
		return nil
	}
	released := releasedVersions(root)
	var errs []string
	for _, d := range docs {
		errs = append(errs, youngErrors(root, d, hist, released)...)
	}
	return errs
}

// youngErrors flags the young lines of d that history shows deleted and the doc does not hold again.
func youngErrors(root string, d requirement.Document, hist history, released []string) []string {
	kept := keptVersions(root, d.Path)
	var errs []string
	for _, x := range slices.SortedFunc(maps.Keys(hist.drops), func(a, b drop) int { return strings.Compare(a.id+a.version, b.id+b.version) }) {
		newer := newerThan(released, x.version)
		if strings.HasPrefix(x.id, d.ID+".") && !kept[x.version] && newer < agedReleases {
			errs = append(errs, requirement.At(d.Path, d.TitleLine, tooYoung(x, newer)))
		}
	}
	return errs
}

// tooYoung is the error for a Removed line deleted with too few newer released versions.
func tooYoung(x drop, newer int) string {
	_, n, _ := strings.Cut(x.id, ".")
	return fmt.Sprintf("%s was Removed in %q and %d released versions are newer; keep its `%s. Removed` line until %d are. "+
		"To drop old markers by hand, delete each `Removed` line with %d or more newer released versions in the changelog, "+
		"move each later acceptance up one number, and rename its test and step cases to match in the same change, "+
		"such as `Test0001_4` to `Test0001_3` and `0001.4.1` to `0001.3.1`",
		x.id, x.version, newer, n, agedReleases, agedReleases)
}

// newerThan counts the released versions newer than version, none for an empty one.
func newerThan(released []string, version string) int {
	if version == "" {
		return 0
	}
	return len(slices.DeleteFunc(slices.Clone(released), func(v string) bool { return v <= version }))
}

// keptVersions files the versions of the Removed lines the doc at path holds.
func keptVersions(root, path string) map[string]bool {
	kept := map[string]bool{}
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		return kept
	}
	for _, c := range docChanges(diffHeader + "a/" + path + " b/" + path + "\n+" + strings.ReplaceAll(string(b), "\n", "\n+")) {
		for _, e := range markers(c.now) {
			kept[versionOf(e.name)] = true
		}
	}
	return kept
}

// releasedVersions lists the versions in the changelog entry files of the base commit.
func releasedVersions(root string) []string {
	base := gitx.ScopeBase(root)
	if base == "" {
		return nil
	}
	text, _ := gitx.ShowAll(root, changelogObjects(root, base))
	var versions []string
	for _, body := range text {
		for _, m := range versionHeading.FindAllStringSubmatch(body, -1) {
			versions = append(versions, m[1])
		}
	}
	return versions
}

// changelogObjects names, as rev:path, the changelog files base holds.
func changelogObjects(root, base string) []string {
	out, err := exec.Command("git", "-C", root, "ls-tree", "-r", "-z", "--name-only", base, "--", changelogDir).Output()
	if err != nil {
		return nil
	}
	var objects []string
	for _, f := range strings.Split(string(out), "\x00") {
		objects = append(objects, base+":"+f)
	}
	return objects
}

// uncommittedDocs reads the docs at paths as a diff from HEAD to the working tree.
// It holds only docs both have that differ, and is empty without a HEAD.
func uncommittedDocs(root string, paths []string) string {
	objects := make([]string, len(paths))
	for i, p := range paths {
		objects[i] = "HEAD:" + p
	}
	was, _ := gitx.ShowAll(root, objects)
	var diff strings.Builder
	for _, p := range paths {
		if old, now, ok := changedDoc(root, p, was["HEAD:"+p]); ok {
			diff.WriteString(docSwap(p, old, now))
		}
	}
	return diff.String()
}

// changedDoc is the doc at p as HEAD held it and as the working tree holds it, and whether they differ.
func changedDoc(root, p, old string) (string, string, bool) {
	now, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(p)))
	return old, string(now), old != "" && err == nil && old != string(now)
}

// docSwap is a diff of the doc at p that removes each old line and adds each new line.
func docSwap(p, old, now string) string {
	return diffHeader + "a/" + p + " b/" + p + "\n-" + strings.ReplaceAll(old, "\n", "\n-") + "\n+" + strings.ReplaceAll(now, "\n", "\n+") + "\n"
}
