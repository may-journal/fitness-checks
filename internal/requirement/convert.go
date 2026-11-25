package requirement

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ToGherkin is the up of migration 0003.6: it rewrites bullet acceptances as
// one all-caps Gherkin block. Commas stay for a person to split.
func ToGherkin(root, path string) error {
	return rewriteSection(root, path, false, gherkinBlock)
}

// ToBullets is the down of migration 0003.6: it rewrites the Gherkin block as
// bullets. Bullets hold no AND lines, so each joins its line with "and".
func ToBullets(root, path string) error {
	return rewriteSection(root, path, true, func(_ string, section []string) string {
		return bulletsFrom(sectionLines(section))
	})
}

// rewriteSection replaces the Requirements section of the doc at path with
// what render makes of it. A doc that fails, or is not in the form gherkin
// names, stays as it is.
func rewriteSection(root, path string, gherkin bool, render func(docID string, section []string) string) error {
	file := filepath.Join(root, filepath.FromSlash(path))
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	d, errs := ParseWith(path, string(b), MaxSteps(root))
	lines := strings.SplitAfter(string(b), "\n")
	from, to := requirementsSpan(lines)
	if unchanged(d, errs, from) || inGherkin(lines[from:to]) != gherkin {
		return nil
	}
	out := strings.Join(lines[:from+1], "") + "\n" + render(d.ID, lines[from+1:to]) + tailAfter(lines[to:])
	return os.WriteFile(file, []byte(out), 0o644)
}

// unchanged reports whether a doc stays as it is: it fails, holds no
// acceptance or has no Requirements section.
func unchanged(d Document, errs []string, from int) bool {
	return len(errs) > 0 || len(d.Acceptances) == 0 || from < 0
}

// inGherkin reports whether section lines hold a Gherkin block.
func inGherkin(lines []string) bool {
	return strings.Contains(strings.Join(lines, ""), gherkinFence)
}

// requirementsSpan is the index of the Requirements heading and of the line
// past its section, or -1 when the doc has none.
func requirementsSpan(lines []string) (int, int) {
	from := -1
	for i, l := range lines {
		switch {
		case strings.TrimSpace(l) == "## Requirements":
			from = i
		case from >= 0 && strings.HasPrefix(l, "## "):
			return from, i
		}
	}
	return from, len(lines)
}

// tailAfter is what follows the section, set apart by a blank line.
func tailAfter(rest []string) string {
	if len(rest) == 0 {
		return ""
	}
	return "\n" + strings.Join(rest, "")
}

// gherkinBlock renders the bullet acceptances in section as SCENARIOs, each
// keyword in caps and every step at one indent under its SCENARIO.
func gherkinBlock(docID string, section []string) string {
	var scenarios []string
	for _, g := range groupItems(&problems{}, docID, sectionLines(section)) {
		s := scenarioWord + " " + g[0].text + "\n"
		for _, it := range g[1:] {
			word, rest, _ := strings.Cut(it.text, " ")
			s += "  " + strings.ToUpper(word) + " " + rest + "\n"
		}
		scenarios = append(scenarios, s)
	}
	return gherkinFence + "\n" + strings.Join(scenarios, "\n") + fence + "\n"
}

// sectionLines numbers the non-blank lines of a section.
func sectionLines(section []string) []line {
	var out []line
	for i, l := range section {
		if t := strings.TrimRight(l, " \t\r\n"); strings.TrimSpace(t) != "" {
			out = append(out, line{i + 1, t})
		}
	}
	return out
}

// scenarioLines is one SCENARIO's ID, its GIVEN, and its WHEN and THEN lines
// in order, each with its AND lines joined on.
type scenarioLines struct {
	id    string
	given string
	steps []string
}

// bulletsFrom renders each SCENARIO in a Gherkin block as a bullet acceptance.
func bulletsFrom(block []line) string {
	var all []*scenarioLines
	for _, l := range block {
		all = addLine(all, l)
	}
	var out string
	for _, s := range all {
		out += s.bullets()
	}
	return out
}

// addLine adds one block line to the SCENARIOs read so far.
func addLine(all []*scenarioLines, l line) []*scenarioLines {
	word, rest, _ := splitWord(l.text)
	upper, rest := strings.ToUpper(word), strings.TrimSpace(rest)
	if upper == scenarioWord {
		return append(all, &scenarioLines{id: rest})
	}
	if len(all) > 0 {
		all[len(all)-1].add(upper, word, rest)
	}
	return all
}

// add files one keyword line under its SCENARIO.
func (s *scenarioLines) add(upper, word, rest string) {
	switch upper {
	case "GIVEN":
		s.given = rest
	case "AND", "BUT":
		s.join(strings.ToLower(word) + " " + rest)
	case "WHEN", "THEN":
		s.steps = append(s.steps, stepWords[upper]+" "+rest)
	}
}

// join adds an AND or BUT line to the line it extends.
func (s *scenarioLines) join(extra string) {
	if len(s.steps) == 0 {
		s.given += " " + extra
		return
	}
	s.steps[len(s.steps)-1] += " " + extra
}

// bullets renders s as one Given with a single step, or a numbered chain.
func (s *scenarioLines) bullets() string {
	out := "- " + s.id + "\n  - Given " + s.given + "\n"
	if len(s.steps) == 2 {
		return out + "    - " + s.steps[0] + "\n      - " + s.steps[1] + "\n"
	}
	for i := 0; i+1 < len(s.steps); i += 2 {
		out += fmt.Sprintf("    %d. %s\n       - %s\n", i/2+1, s.steps[i], s.steps[i+1])
	}
	return out
}

// BulletText is doc with its Gherkin block read back as bullets. A doc moved
// to Gherkin then compares with its bullet form word for word.
func BulletText(doc string) string {
	lines := strings.SplitAfter(doc, "\n")
	from, to := requirementsSpan(lines)
	if from < 0 || !inGherkin(lines[from:to]) {
		return doc
	}
	return strings.Join(lines[:from+1], "") + bulletsFrom(sectionLines(lines[from+1:to])) + strings.Join(lines[to:], "")
}
