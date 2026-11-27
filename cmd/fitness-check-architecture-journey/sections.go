package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-runner/go/pkg/gitx"
)

const (
	sectionMarker = "## "
	sectionFix    = "extend an existing flow with the requirement's numbered nodes and callout rows, and cite the requirement there"
)

// sectionErrors fails each `##` section a journey doc of the package root dir holds that its base lacks.
// The base is a pull request's merge base or else the last commit, and a doc new in the change keeps its sections.
func sectionErrors(root, dir string, files []string) []string {
	if gitx.Replaying(root) {
		return nil
	}
	base := gitx.ScopeBase(root)
	var errs []string
	for _, file := range files {
		if isJourneyDoc(packages.Rel(dir, file)) {
			errs = append(errs, docSections(root, base, file)...)
		}
	}
	return errs
}

// isJourneyDoc reports whether rel, inside a package root, is its journey view.
func isJourneyDoc(rel string) bool {
	folder := path.Dir(rel)
	return (folder == "architecture" || folder == "docs/architecture") && path.Base(rel) == "05-user-journey.md"
}

// docSections names the new sections of file, or none when base lacks the file.
func docSections(root, base, file string) []string {
	before, existed := gitx.Show(root, base, file)
	body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file)))
	if !existed || err != nil {
		return nil
	}
	return newSections(file, before, string(body))
}

// newSections names each section of body that before lacks, at its line.
// Two sections of one title count against the copies before holds.
func newSections(file, before, body string) []string {
	had := map[string]int{}
	for _, s := range sections(before) {
		had[s.title]++
	}
	var errs []string
	for _, s := range sections(body) {
		if had[s.title] > 0 {
			had[s.title]--
			continue
		}
		errs = append(errs, fmt.Sprintf("%s:%d: the journey doc gains the new section `%s%s`; %s", file, s.line, sectionMarker, s.title, sectionFix))
	}
	return errs
}

// section is one `##` heading of a doc and its line.
type section struct {
	title string
	line  int
}

// sections lists the `##` headings of text outside code fences.
func sections(text string) []section {
	var out []section
	fenced := false
	for i, line := range strings.Split(text, "\n") {
		if isFence(line) {
			fenced = !fenced
		} else if title, ok := strings.CutPrefix(line, sectionMarker); ok && !fenced {
			out = append(out, section{strings.TrimSpace(title), i + 1})
		}
	}
	return out
}

// isFence reports a line that opens or closes a code fence.
func isFence(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
}
