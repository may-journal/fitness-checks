package architecture

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// maxLabelIDs is the most requirement IDs one journey label may name, so a
// step that fulfils more is split.
const maxLabelIDs = 8

// labelID is a requirement number, or an acceptance ID, written in a
// journey label after its callout number.
var labelID = regexp.MustCompile(`\b\d{4}(?:\.\d+)?\b`)

// entityCode is a Mermaid character code such as #9829; in a label, which
// names no requirement.
var entityCode = regexp.MustCompile(`#\w+;`)

const journeyView = "journey"

// journeyStart records where the journey begins: its first diagram's fence
// line, or the view's first line when it has no diagram.
func (m *Model) journeyStart(at Location) {
	if m.journeyAt == (Location{}) {
		m.journeyAt = at
	}
}

// JourneyFiles counts the files the journey check judges: journey views and
// requirement docs.
func (m *Model) JourneyFiles() int {
	count := len(m.Requirements)
	for _, file := range m.Files {
		if viewFile(file) == journeyView {
			count++
		}
	}
	return count
}

// JourneyErrors compares requirement IDs named in journey labels with
// requirement docs, both ways, and requires each behavior node to name one.
// Its notes are the steps that read as a doc they do not name.
func (m *Model) JourneyErrors() (errs, notes []string) {
	if m.journeyAt == (Location{}) {
		return nil, nil
	}
	j := &journey{m: m, named: map[string]bool{}, nodes: map[string]bool{}}
	j.labels()
	j.docs()
	j.behavior()
	j.swapped()
	return j.errs, j.notes
}

// labels reads every journey label: elements first, then edges.
func (j *journey) labels() {
	for _, e := range j.m.Elements {
		if e.View == journeyView {
			j.label(e.ID, e.Name, e.At)
		}
	}
	for _, r := range j.m.Relationships {
		if r.View == journeyView {
			j.label("", r.Label, r.At)
		}
	}
}

// journey holds one comparison: the docs some label names and the elements
// whose own label names one.
type journey struct {
	m     *Model
	named map[string]bool
	nodes map[string]bool
	errs  []string
	notes []string
}

func (j *journey) fail(at Location, format string, args ...any) {
	j.errs = append(j.errs, at.String()+": "+fmt.Sprintf(format, args...))
}

// label reads the requirement IDs one label names for its owner element. A
// label that names any ID covers its node, so a mistyped ID fails once.
func (j *journey) label(owner, text string, at Location) {
	ids := uniqueIDs(labelID.FindAllString(entityCode.ReplaceAllString(text, ""), -1))
	if len(ids) > maxLabelIDs {
		j.fail(at, "label names %d requirements; split the step so each names at most %d", len(ids), maxLabelIDs)
	}
	for _, id := range ids {
		j.named[id] = j.valid(id, at)
	}
	if owner != "" && len(ids) > 0 {
		j.nodes[owner] = true
	}
}

func uniqueIDs(ids []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

// valid requires a label ID to name a requirement doc that exists. An
// acceptance ID names a test's promise, not a doc.
func (j *journey) valid(id string, at Location) bool {
	switch {
	case acceptanceID.MatchString(id):
		j.fail(at, "%s is an acceptance ID; name requirement doc %s", id, id[:4])
	case !j.m.hasDoc(id):
		j.fail(at, "no requirement doc %s in docs/requirements", id)
	default:
		return true
	}
	return false
}

func (m *Model) hasDoc(id string) bool {
	for _, doc := range m.Requirements {
		if doc.ID == id {
			return true
		}
	}
	return false
}

// docs requires every requirement doc in some journey label. A doc without
// a valid ID belongs to the requirements check.
func (j *journey) docs() {
	for _, doc := range j.m.Requirements {
		if doc.ID != "" && !j.named[doc.ID] {
			j.fail(j.m.journeyAt, "requirement %s is missing from the journey diagram; name it in the label of the step that fulfils it", doc.ID)
		}
	}
}

// behavior requires every node that promises behavior to name a
// requirement in its own label or an enclosing subgraph's. Personas and
// layout subgraphs promise none.
func (j *journey) behavior() {
	for _, e := range j.m.Elements {
		if e.View == journeyView && promises(e.Kind) && !j.covered(e) {
			j.fail(e.At, "journey node %s names no requirement; add the IDs it fulfils to its label or type it as a persona", e.ID)
		}
	}
}

// covered reports a requirement ID in the element's label or in the label
// of a subgraph around it.
func (j *journey) covered(e Element) bool {
	for id := e.ID; id != ""; id = j.parent(id) {
		if j.nodes[id] {
			return true
		}
	}
	return false
}

func (j *journey) parent(id string) string {
	for _, e := range j.m.Elements {
		if e.View == journeyView && e.ID == id {
			return e.Parent
		}
	}
	return ""
}

func promises(kind string) bool {
	return kind != "subgraph" && !strings.HasPrefix(strings.ToLower(kind), "person")
}

// labelWord is one word of a journey label or a requirement doc's file name.
var labelWord = regexp.MustCompile(`[a-z]{3,}`)

// fillerWords carry no meaning in a title, so they never tie a step to a doc.
var fillerWords = []string{"and", "the", "for", "with", "its", "own", "are", "req"}

// titleWords are the meaningful words of text, lowercased, each without a
// plural `s`, so `Events` and `event` are one word.
func titleWords(text string) []string {
	var out []string
	for _, w := range labelWord.FindAllString(strings.ToLower(text), -1) {
		if slices.Contains(fillerWords, w) {
			continue
		}
		if w = strings.TrimSuffix(w, "s"); !slices.Contains(out, w) {
			out = append(out, w)
		}
	}
	return out
}

// docWords are the meaningful words of a requirement doc's file name, which
// the requirements check holds to the doc's title.
func docWords(doc string) []string {
	return titleWords(strings.TrimSuffix(path.Base(doc), ".md"))
}

// misread is a step whose label reads as docs it leaves out.
type misread struct {
	step   Element
	ids    []string
	readAs []string
}

// swapped fails two steps that read as each other's docs, since their labels
// are crossed. One such step alone may only be worded loosely, so it adds a note.
func (j *journey) swapped() {
	found := j.misreads()
	for _, f := range found {
		text := fmt.Sprintf("journey step %s reads as requirement %s yet names %s", f.step.ID, strings.Join(f.readAs, ", "), strings.Join(f.ids, ", "))
		if crossed(f, found) {
			j.fail(f.step.At, "%s; name the doc whose behavior the step is, or reword the step", text)
			continue
		}
		j.notes = append(j.notes, f.step.At.String()+": "+text+"; check that the step names the doc whose behavior it is")
	}
}

// misreads lists the behavior steps whose label reads as docs it leaves out.
func (j *journey) misreads() []misread {
	var found []misread
	for _, e := range j.m.Elements {
		if e.View != journeyView || !promises(e.Kind) {
			continue
		}
		ids := uniqueIDs(labelID.FindAllString(entityCode.ReplaceAllString(e.Name, ""), -1))
		if others := j.readsAs(titleWords(e.Name), ids); len(others) > 0 {
			found = append(found, misread{e, ids, others})
		}
	}
	return found
}

// crossed reports another step that names a doc f reads as, and reads as a
// doc f names.
func crossed(f misread, all []misread) bool {
	return slices.ContainsFunc(all, func(o misread) bool {
		return shares(o.ids, f.readAs) && shares(f.ids, o.readAs)
	})
}

// readsAs lists each doc the label leaves out while holding every title
// word of it. It lists none for a label that fits a doc it names.
func (j *journey) readsAs(words, ids []string) []string {
	if j.fits(words, ids) {
		return nil
	}
	var others []string
	for _, doc := range j.m.Requirements {
		if !slices.Contains(ids, doc.ID) && holdsAll(words, docWords(doc.Path)) {
			others = append(others, doc.ID)
		}
	}
	return others
}

// fits reports a label that names no doc, or shares a title word with a doc
// it names.
func (j *journey) fits(words, ids []string) bool {
	return len(ids) == 0 || slices.ContainsFunc(j.m.Requirements, func(doc requirement.Document) bool {
		return slices.Contains(ids, doc.ID) && shares(words, docWords(doc.Path))
	})
}

// shares reports a word that both lists hold.
func shares(a, b []string) bool {
	return slices.ContainsFunc(a, func(w string) bool { return slices.Contains(b, w) })
}

// holdsAll reports a title with words, every one of them among words.
func holdsAll(words, title []string) bool {
	return len(title) > 0 && !slices.ContainsFunc(title, func(w string) bool { return !slices.Contains(words, w) })
}
