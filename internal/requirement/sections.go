package requirement

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var acceptID = regexp.MustCompile(`^(\d{4})\.[1-9]\d*$`)

// numberedItem is a numbered list line, "2. When ...", with its number and text.
var numberedItem = regexp.MustCompile(`^(\d{1,3})\. (.*)$`)

// chainKeywords open the three nested lines under each acceptance ID.
var chainKeywords = []string{"Given", "When", "Then"}

// stepShape is what a numbered step that is not a When fails with.
const stepShape = "each step is a numbered When nested under the Given, with one Then nested under the When"

// sourceError is what a line starting "Source:" fails with. A line
// number in a doc goes stale as soon as the code moves, so docs cite none.
const sourceError = "delete this `Source:` line, or run `fitness requirements --fix`; a Measurement is a numerator, `-` and a denominator"

// checkMeasurement requires a numerator, a "-" line, and a denominator, and
// fails a line starting "Source:", which older docs carry. It
// returns that line's number, or 0, so the fix can delete it.
func checkMeasurement(p *problems, s section) int {
	lines, source := withoutSource(s.lines)
	if source > 0 {
		p.add(source, sourceError)
	}
	if len(lines) > 0 && !isRatio(lines) {
		p.add(s.n, "Measurement needs three lines: numerator, `-`, denominator")
	}
	return source
}

// withoutSource drops the first line starting "Source:" and returns its
// number, or 0 when there is none.
func withoutSource(lines []line) ([]line, int) {
	for i, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l.text), "Source:") {
			rest := append(append([]line{}, lines[:i]...), lines[i+1:]...)
			return rest, l.n
		}
	}
	return lines, 0
}

// isRatio reports whether lines are a numerator, a "-" line, and a
// denominator.
func isRatio(lines []line) bool {
	return len(lines) == 3 && strings.TrimSpace(lines[1].text) == "-"
}

// DropSource deletes d's Source line from its file under root, with the
// blank line before it, so the Measurement ends at its denominator.
func DropSource(root string, d Document) error {
	if d.SourceLine == 0 {
		return nil
	}
	file := filepath.Join(root, filepath.FromSlash(d.Path))
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	lines := strings.SplitAfter(string(b), "\n")
	from := d.SourceLine - 1
	if from > 0 && strings.TrimSpace(lines[from-1]) == "" {
		from--
	}
	kept := append(lines[:from:from], lines[d.SourceLine:]...)
	return os.WriteFile(file, []byte(strings.Join(kept, "")), 0o644)
}

// item is one list line: its indent width, its text after the "- " or
// "N. " marker, and for a numbered line its number.
type item struct {
	line
	indent   int
	numbered bool
	num      int
}

// parseRequirements groups list items under top-level acceptance IDs and
// checks each one's Given, When, Then chain.
func parseRequirements(p *problems, docID string, s section) []IDLine {
	if len(s.lines) == 0 {
		return nil
	}
	var out []IDLine
	seen := map[string]bool{}
	for _, g := range groupItems(p, docID, s.lines) {
		a := checkAcceptance(p, docID, g)
		if seen[a.ID] {
			p.add(a.Line, a.ID+" appears twice; each acceptance needs its own ID")
		}
		seen[a.ID] = true
		out = append(out, a)
	}
	return out
}

// groupItems splits list items into groups, each a top-level item and the
// items nested under it.
func groupItems(p *problems, docID string, lines []line) [][]item {
	var groups [][]item
	for _, l := range lines {
		it, ok := toItem(p, docID, l)
		if !ok {
			continue
		}
		if it.indent == 0 || len(groups) == 0 {
			groups = append(groups, nil)
		}
		groups[len(groups)-1] = append(groups[len(groups)-1], it)
	}
	return groups
}

// toItem parses a "- text" or "N. text" line. Any other line is not an
// acceptance.
func toItem(p *problems, docID string, l line) (item, bool) {
	trimmed := strings.TrimLeft(l.text, " \t")
	indent := len(strings.ReplaceAll(l.text[:len(l.text)-len(trimmed)], "\t", "    "))
	if text, ok := strings.CutPrefix(trimmed, "- "); ok {
		return item{line: line{l.n, strings.TrimSpace(text)}, indent: indent}, true
	}
	if m := numberedItem.FindStringSubmatch(trimmed); m != nil {
		n, _ := strconv.Atoi(m[1])
		return item{line: line{l.n, strings.TrimSpace(m[2])}, indent: indent, numbered: true, num: n}, true
	}
	p.add(l.n, idMessage(docID))
	return item{}, false
}

// checkAcceptance requires a top-level NNNN.N ID with the doc's prefix,
// nesting one Given and one step, or a numbered chain of steps.
func checkAcceptance(p *problems, docID string, g []item) IDLine {
	checkAcceptanceID(p, docID, g[0])
	if slices.ContainsFunc(g[1:], func(it item) bool { return it.numbered }) {
		return checkChain(p, g)
	}
	return checkSingle(p, g)
}

// checkSingle requires exactly one Given, one When and one Then, each a
// bullet nested under the line above.
func checkSingle(p *problems, g []item) IDLine {
	top := g[0]
	if len(g) != len(chainKeywords)+1 {
		p.add(top.n, top.text+" needs exactly one Given, one When, and one Then; a branch is a new ID")
		return IDLine{ID: top.text, Line: top.n}
	}
	ordered := true
	for i, kw := range chainKeywords {
		ordered = checkChainLine(p, kw, g[i], g[i+1]) && ordered
	}
	a := withGiven(g)
	if ordered {
		a.Steps = []Step{{Line: g[2].n, When: strings.TrimPrefix(g[2].text, "When ")}}
	}
	return a
}

// withGiven is the acceptance g opens, with the Given on its second line.
func withGiven(g []item) IDLine {
	return IDLine{ID: g[0].text, Line: g[0].n, Given: strings.TrimPrefix(g[1].text, "Given "), GivenLine: g[1].n}
}

// checkChainLine requires child to be a bullet that starts with kw and sits
// deeper than parent, and reports whether it does.
func checkChainLine(p *problems, kw string, parent, child item) bool {
	if child.numbered || !nests(kw, parent, child) {
		p.add(child.n, "each acceptance is one Given, then one When, then one Then, each nested under the line above")
		return false
	}
	return true
}

// nests reports whether child starts with kw and sits deeper than parent.
func nests(kw string, parent, child item) bool {
	return strings.HasPrefix(child.text, kw+" ") && child.indent > parent.indent
}

// checkChain requires one Given, then 2 to p.maxSteps numbered steps. Step 2
// of 0003.6 is 0003.6.2.
func checkChain(p *problems, g []item) IDLine {
	if !checkChainLine(p, "Given", g[0], g[1]) {
		return IDLine{ID: g[0].text, Line: g[0].n}
	}
	a := withGiven(g)
	steps, ok := chainSteps(p, a.ID, g[1], g[2:])
	switch {
	case !ok:
	case len(steps) > p.maxSteps:
		p.add(a.Line, fmt.Sprintf("%s chains %d steps, at most %d; move the rest to a new ID", a.ID, len(steps), p.maxSteps))
	case len(steps) == 1:
		p.add(steps[0].Line, "a single step is the plain bullet `- When`; number steps only in a chain of two or more")
	default:
		a.Steps = steps
	}
	return a
}

// chainSteps reads the steps of acceptance id's chain under given, and
// reports whether each one is well formed.
func chainSteps(p *problems, id string, given item, rest []item) ([]Step, bool) {
	var steps []Step
	for ; len(rest) > 0; rest = rest[min(2, len(rest)):] {
		n := len(steps) + 1
		if !checkStep(p, n, given, rest) {
			return steps, false
		}
		steps = append(steps, Step{ID: fmt.Sprintf("%s.%d", id, n), Line: rest[0].n, When: strings.TrimPrefix(rest[0].text, "When ")})
	}
	return steps, true
}

// checkStep requires rest to open with step n: a When numbered n under the
// Given, with one Then bullet under it.
func checkStep(p *problems, n int, given item, rest []item) bool {
	when := rest[0]
	switch {
	case !nests("When", given, when):
		p.add(when.n, stepShape)
	case !numberedAs(when, n):
		p.add(when.n, fmt.Sprintf("number this step `%d.`; a chain numbers its steps `1.` to `%d.` in order", n, p.maxSteps))
	case !thenUnder(when, rest[1:]):
		p.add(when.n, fmt.Sprintf("step %d needs exactly one Then nested under its When", n))
	default:
		return true
	}
	return false
}

// numberedAs reports whether it is numbered n.
func numberedAs(it item, n int) bool {
	return it.numbered && it.num == n
}

// thenUnder reports whether after opens with a Then bullet nested under when.
func thenUnder(when item, after []item) bool {
	return len(after) > 0 && !after[0].numbered && nests("Then", when, after[0])
}

// checkAcceptanceID requires a top-level `NNNN.N` bullet with the doc's ID.
func checkAcceptanceID(p *problems, docID string, top item) {
	m := acceptID.FindStringSubmatch(top.text)
	if top.indent != 0 || top.numbered || m == nil || m[1] != docID {
		p.add(top.n, idMessage(docID))
	}
}

// idMessage says what a top-level Requirements line must be.
func idMessage(docID string) string {
	return "top-level items must be acceptance IDs `" + docID + ".N`"
}
