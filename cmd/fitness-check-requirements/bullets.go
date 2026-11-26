package main

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// asBullets is doc, found at path, with its numbered acceptances read back as
// bullets, and whether it held them. A doc moved to the numbered form by hand
// then compares with its bullet form word for word.
func asBullets(path, doc string) (string, bool) {
	d, _ := requirement.ParseWith(path, doc, math.MaxInt)
	if !d.Gherkin {
		return doc, false
	}
	return replaced(doc, bullets(d)), true
}

// bullets renders d's acceptances as one Given with a single step or a
// numbered chain.
func bullets(d requirement.Document) string {
	var out string
	for _, a := range d.Acceptances {
		out += "- " + a.ID + "\n  - Given " + a.Given + "\n"
		if len(a.Steps) == 1 {
			out += "    - When " + a.Steps[0].When + "\n      - Then " + a.Steps[0].Then + "\n"
			continue
		}
		for i, st := range a.Steps {
			out += fmt.Sprintf("    %d. When %s\n       - Then %s\n", i+1, st.When, st.Then)
		}
	}
	return out
}

// replaced is doc with the body of its Requirements section, comments
// included, swapped for section. Only that section changes.
func replaced(doc, section string) string {
	lines := strings.SplitAfter(doc, "\n")
	from, to := span(lines)
	return strings.Join(lines[:from+1], "") + "\n" + section + tailAfter(lines[to:])
}

// span is the index of the Requirements heading and of the line past its
// section.
func span(lines []string) (int, int) {
	from := slices.IndexFunc(lines, func(l string) bool { return strings.TrimSpace(l) == "## Requirements" })
	next := slices.IndexFunc(lines[from+1:], func(l string) bool { return strings.HasPrefix(l, "## ") })
	if next < 0 {
		return from, len(lines)
	}
	return from, from + 1 + next
}

// tailAfter is what follows the section, set apart by a blank line.
func tailAfter(rest []string) string {
	if len(rest) == 0 {
		return ""
	}
	return "\n" + strings.Join(rest, "")
}
