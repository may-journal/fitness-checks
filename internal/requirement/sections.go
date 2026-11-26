package requirement

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var acceptID = regexp.MustCompile(`^(\d{4})\.[1-9]\d*$`)

// numberedItem is a numbered list line, "2. When ...", with its number and text.
var numberedItem = regexp.MustCompile(`^(\d{1,3})\. (.*)$`)

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

// NumberedForm is where the requirements README shows the numbered
// acceptance form, with GIVEN, WHEN and THEN in caps.
const NumberedForm = "https://github.com/may-journal/fitness-checks/blob/main/cmd/fitness-check-requirements/README.md#template"

// topLevelMessage is what a Requirements line that is no list item fails with.
const topLevelMessage = "top-level items must be numbered acceptance names, `1. Name`"

// parseRequirements reads a Requirements section of numbered acceptances. A
// section that opens with any other list, such as the retired bullet form,
// fails at each acceptance.
func parseRequirements(p *problems, docID string, s section) []IDLine {
	if len(s.lines) == 0 {
		return nil
	}
	if hasGherkin(s) {
		return parseGherkin(p, docID, s)
	}
	return bulletAcceptances(p, groupItems(p, s.lines, topLevelMessage))
}

// bulletAcceptances fails each acceptance of a section not in the numbered
// form at its first line. It keeps each `- NNNN.N` ID and its numbered step
// IDs, so tests and views still name them while the doc moves.
func bulletAcceptances(p *problems, groups [][]item) []IDLine {
	var out []IDLine
	for _, g := range groups {
		top := g[0]
		p.add(top.n, "this is the bullet acceptance form, which `requirements` no longer reads; "+
			"number each acceptance and write its GIVEN, WHEN and THEN items in caps, as "+NumberedForm+" shows")
		if !top.numbered && acceptID.MatchString(top.text) {
			out = append(out, IDLine{ID: top.text, Line: top.n, Steps: bulletSteps(top.text, g[1:])})
		}
	}
	return out
}

// bulletSteps is the numbered `N. When` steps of bullet acceptance id, each
// with its ID, such as 0001.1.2.
func bulletSteps(id string, items []item) []Step {
	var steps []Step
	for _, it := range items {
		if it.numbered && strings.HasPrefix(it.text, "When ") {
			steps = append(steps, Step{ID: fmt.Sprintf("%s.%d", id, it.num), Line: it.n})
		}
	}
	return steps
}

// groupItems splits list items into groups, each a top-level item and the
// items nested under it. A line that is no list item fails with msg.
func groupItems(p *problems, lines []line, msg string) [][]item {
	var groups [][]item
	for _, l := range lines {
		it, ok := listItem(l)
		if !ok {
			p.add(l.n, msg)
			continue
		}
		if it.indent == 0 || len(groups) == 0 {
			groups = append(groups, nil)
		}
		groups[len(groups)-1] = append(groups[len(groups)-1], it)
	}
	return groups
}

// listItem parses a "- text" or "N. text" line, and reports whether it is one.
func listItem(l line) (item, bool) {
	trimmed := strings.TrimLeft(l.text, " \t")
	indent := len(strings.ReplaceAll(l.text[:len(l.text)-len(trimmed)], "\t", "    "))
	if text, ok := strings.CutPrefix(trimmed, "- "); ok {
		return item{line: line{l.n, strings.TrimSpace(text)}, indent: indent}, true
	}
	if m := numberedItem.FindStringSubmatch(trimmed); m != nil {
		n, _ := strconv.Atoi(m[1])
		return item{line: line{l.n, strings.TrimSpace(m[2])}, indent: indent, numbered: true, num: n}, true
	}
	return item{}, false
}
