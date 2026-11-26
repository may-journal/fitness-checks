// Package gherkin is migration 0003.2-gherkin. It moves a requirement doc's
// acceptances between bullets and numbered GIVEN, WHEN and THEN items.
package gherkin

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// ID names the migration, and the folder that holds it.
const ID = "0003.2-gherkin"

// rewrite replaces the Requirements section of the doc at path with what
// render makes of its acceptances. A doc it can't read or won't move stays
// as it is.
func rewrite(root, path string, toGherkin bool, render func(requirement.Document) string) error {
	file := filepath.Join(root, filepath.FromSlash(path))
	b, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	d, errs := requirement.ParseWith(path, string(b), requirement.MaxSteps(root))
	if !moves(d, errs, toGherkin) {
		return nil
	}
	return os.WriteFile(file, []byte(replaced(string(b), render(d))), 0o644)
}

// moves reports whether a passing doc with acceptances moves to the other
// form.
func moves(d requirement.Document, errs []string, toGherkin bool) bool {
	return len(errs) == 0 && len(d.Acceptances) > 0 && d.Gherkin != toGherkin
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
