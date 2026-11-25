package architecture

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-runner/go/pkg/walkfs"
)

// Markers are Mermaid comments around the nodes generated from the requirement docs.
// A journey view without them stays hand-written.
const (
	regionStart = "%% journey-nodes:start"
	regionEnd   = "%% journey-nodes:end"
)

// regionFix is the command that rewrites a stale region.
const regionFix = "fitness architecture-journey --fix"

// firstRegionNumber is the callout number of the first generated node, after the persona's 1.
const firstRegionNumber = 2

// regionEscapes are the characters a Mermaid label cannot hold bare.
var regionEscapes = strings.NewReplacer(
	"#", "#35;", `"`, "#34;", "(", "#40;", ")", "#41;", "[", "#91;", "]", "#93;",
	"{", "#123;", "}", "#125;", "<", "#60;", ">", "#62;", "|", "#124;", ";", "#59;",
)

// regionLines writes one system node per requirement doc, in ID order, from its ID and title.
func regionLines(docs []requirement.Document, indent string) []string {
	var out []string
	for _, doc := range sortedDocs(docs) {
		node := fmt.Sprintf("req%s[(%d %s<br/>Req %s)]:::system", doc.ID, firstRegionNumber+len(out), regionEscapes.Replace(doc.Title), doc.ID)
		out = append(out, indent+node)
	}
	return out
}

// sortedDocs lists the docs with an ID and a title, by ID.
func sortedDocs(docs []requirement.Document) []requirement.Document {
	var out []requirement.Document
	for _, doc := range docs {
		if doc.ID != "" && doc.Title != "" {
			out = append(out, doc)
		}
	}
	slices.SortFunc(out, func(a, b requirement.Document) int { return strings.Compare(a.ID, b.ID) })
	return out
}

// region holds the line numbers around a view's marked lines, and the markers' indent.
type region struct {
	from, to int
	indent   string
}

// marked lists the indexes of the lines that hold marker.
func marked(lines []string, marker string) []int {
	var out []int
	for i, l := range lines {
		if strings.TrimSpace(l) == marker {
			out = append(out, i)
		}
	}
	return out
}

// repeatProblem names a marker that appears twice, and the index of its second line.
func repeatProblem(starts, ends []int) (string, int) {
	if len(starts) > 1 {
		return "has a second `" + regionStart + "` marker", starts[1]
	}
	if len(ends) > 1 {
		return "has a second `" + regionEnd + "` marker", ends[1]
	}
	return "", 0
}

// pairProblem names what is wrong with the markers' pairing, and the index it concerns.
func pairProblem(starts, ends []int) (string, int) {
	if problem, i := repeatProblem(starts, ends); problem != "" {
		return problem, i
	}
	switch {
	case len(starts) == 0:
		return "has `" + regionEnd + "` with no start marker", ends[0]
	case len(ends) == 0 || ends[0] < starts[0]:
		return "has `" + regionStart + "` with no end marker after it", starts[0]
	}
	return "", 0
}

// findRegion locates the marked region in a view's lines. A view without markers has none, and a broken pair is a problem.
func findRegion(lines []string) (r region, found bool, problem string, at int) {
	starts, ends := marked(lines, regionStart), marked(lines, regionEnd)
	if len(starts)+len(ends) == 0 {
		return r, false, "", 0
	}
	if problem, i := pairProblem(starts, ends); problem != "" {
		return r, false, problem, i + 1
	}
	line := lines[starts[0]]
	return region{starts[0] + 1, ends[0], line[:len(line)-len(strings.TrimLeft(line, " \t"))]}, true, "", 0
}

// journeyViews lists the journey view files under root.
func journeyViews(files []string) []string {
	var out []string
	for _, file := range files {
		if viewFile(file) == journeyView {
			out = append(out, file)
		}
	}
	return out
}

// loadDocs reads the requirement docs under root and the files around them.
func loadDocs(root string) ([]requirement.Document, []string) {
	files := walkfs.FilesByExt(root, "")
	var paths []string
	for _, file := range files {
		if requirementFile(file) {
			paths = append(paths, file)
		}
	}
	docs, _, _ := requirement.Load(root, paths)
	return docs, files
}

// JourneyRegionErrors fails each journey view whose marked region differs
// from the nodes its requirement docs generate. A view without markers is
// hand-written and passes.
func (m *Model) JourneyRegionErrors() []string {
	if m.repo.root == "" || len(m.Files) == 0 {
		return nil
	}
	var errs []string
	for _, file := range journeyViews(m.Files) {
		body, err := os.ReadFile(filepath.Join(m.repo.root, filepath.FromSlash(file)))
		if err != nil {
			errs = append(errs, file+":1: "+err.Error())
			continue
		}
		errs = append(errs, regionErrors(file, strings.Split(string(body), "\n"), m.Requirements)...)
	}
	return errs
}

func regionErrors(file string, lines []string, docs []requirement.Document) []string {
	r, found, problem, at := findRegion(lines)
	switch {
	case problem != "":
		return []string{fmt.Sprintf("%s:%d: journey view %s; mark one region with both markers, then run `%s`", file, at, problem, regionFix)}
	case !found:
		return nil
	}
	if slices.Equal(trimmed(lines[r.from:r.to]), trimmed(regionLines(docs, r.indent))) {
		return nil
	}
	return []string{fmt.Sprintf("%s:%d: the generated journey nodes differ from the requirement docs; run `%s` to rewrite them", file, r.from, regionFix)}
}

// trimmed drops each line's indent and trailing space, so only the nodes count.
func trimmed(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = strings.TrimSpace(l)
	}
	return out
}

// RewriteJourneyRegion rewrites the marked region of each journey view under
// root from the requirement docs, leaving every other line as written. A
// view without markers is left alone.
func RewriteJourneyRegion(root string) error {
	docs, files := loadDocs(root)
	for _, file := range journeyViews(files) {
		full := filepath.Join(root, filepath.FromSlash(file))
		body, err := os.ReadFile(full)
		if err != nil {
			return err
		}
		lines := strings.Split(string(body), "\n")
		r, found, _, _ := findRegion(lines)
		if !found {
			continue
		}
		next := append(slices.Clone(lines[:r.from]), regionLines(docs, r.indent)...)
		next = append(next, lines[r.to:]...)
		if err := os.WriteFile(full, []byte(strings.Join(next, "\n")), 0o644); err != nil {
			return err
		}
	}
	return nil
}
