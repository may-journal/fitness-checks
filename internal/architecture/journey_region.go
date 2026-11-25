package architecture

import (
	"errors"
	"fmt"
	"github.com/may-journal/fitness-checks/internal/packages"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// Markers are Mermaid comments around nodes generated from the requirement docs, and the start marker lists the doc IDs its region holds.
// A journey view without markers stays hand-written.
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

// regionLines writes one system node per listed requirement doc, in ID order, from its ID and title.
// A node's callout number is its place among all docs, so the regions number on without a gap.
func regionLines(docs []requirement.Document, ids []string, indent string) []string {
	var out []string
	for i, doc := range sortedDocs(docs) {
		if slices.Contains(ids, doc.ID) {
			node := fmt.Sprintf("req%s[(%d %s<br/>Req %s)]:::system", doc.ID, firstRegionNumber+i, regionEscapes.Replace(doc.Title), doc.ID)
			out = append(out, indent+node)
		}
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

// region holds the line numbers around one marked region's lines, the markers' indent and the doc IDs it lists.
type region struct {
	from, to int
	indent   string
	ids      []string
}

// findRegions locates every marked region in a view's lines. A view without markers has none, and a broken pair is a problem.
func findRegions(lines []string) (regions []region, problem string, at int) {
	open := -1
	for i, l := range lines {
		kind := markerKind(l)
		if problem = pairingProblem(kind, open); problem != "" {
			return nil, problem, i + 1
		}
		if kind == regionEnd {
			regions = append(regions, regionAt(lines[open], open, i))
		}
		open = nextOpen(kind, open, i)
	}
	if open >= 0 {
		return nil, "has `" + regionStart + "` with no end marker", open + 1
	}
	return regions, "", 0
}

// nextOpen is the index of the open start marker after a marker of kind at index i, or -1.
func nextOpen(kind string, open, i int) int {
	switch kind {
	case regionStart:
		return i
	case regionEnd:
		return -1
	}
	return open
}

// markerKind names the marker a line holds, or nothing.
func markerKind(line string) string {
	text := strings.TrimSpace(line)
	switch {
	case text == regionStart || strings.HasPrefix(text, regionStart+" "):
		return regionStart
	case text == regionEnd:
		return regionEnd
	}
	return ""
}

// pairingProblem names what is wrong with a marker met while a region is open at index open, or -1.
func pairingProblem(kind string, open int) string {
	switch {
	case kind == regionStart && open >= 0:
		return "has `" + regionStart + "` before the last region's end marker"
	case kind == regionEnd && open < 0:
		return "has `" + regionEnd + "` with no start marker"
	}
	return ""
}

// regionAt reads the region between the start line at index start and the end line at index end.
func regionAt(line string, start, end int) region {
	return region{
		from:   start + 1,
		to:     end,
		indent: line[:len(line)-len(strings.TrimLeft(line, " \t"))],
		ids:    strings.Fields(strings.TrimPrefix(strings.TrimSpace(line), regionStart)),
	}
}

// journeyViews lists the journey view files of the package root dir.
func journeyViews(dir string, files []string) []string {
	var out []string
	for _, file := range files {
		if viewFile(packages.Rel(dir, file)) == journeyView {
			out = append(out, file)
		}
	}
	return out
}

// loadDocs reads the requirement docs of the package root dir and the files
// around them.
func loadDocs(root, dir string) ([]requirement.Document, []string, []string) {
	files := packageFiles(root, dir)
	var paths []string
	for _, file := range files {
		if requirementFile(packages.Rel(dir, file)) {
			paths = append(paths, file)
		}
	}
	docs, _, skipped := requirement.Load(root, paths)
	return docs, files, skipped
}

// JourneyRegionErrors fails each journey view whose marked region differs
// from the nodes its requirement docs generate. A view without markers is
// hand-written and passes.
func (m *Model) JourneyRegionErrors() []string {
	if m.repo.root == "" || len(m.Files) == 0 {
		return nil
	}
	var errs []string
	for _, file := range journeyViews(m.dir, m.Files) {
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
	regions, problem, at := findRegions(lines)
	if problem != "" {
		return []string{fmt.Sprintf("%s:%d: journey view %s; mark each region with both markers, then run `%s`", file, at, problem, regionFix)}
	}
	if len(regions) == 0 {
		return nil
	}
	errs := coverageErrors(file, regions, docs)
	for _, r := range regions {
		if !slices.Equal(trimmed(lines[r.from:r.to]), trimmed(regionLines(docs, r.ids, r.indent))) {
			errs = append(errs, fmt.Sprintf("%s:%d: the generated journey nodes differ from the requirement docs; run `%s` to rewrite them", file, r.from, regionFix))
		}
	}
	return errs
}

// coverageErrors fails a doc that no region lists, a doc two regions list and an ID that is no doc.
func coverageErrors(file string, regions []region, docs []requirement.Document) []string {
	var errs []string
	known := map[string]bool{}
	for _, doc := range sortedDocs(docs) {
		known[doc.ID] = true
		errs = append(errs, listingErrors(file, regions, doc.ID)...)
	}
	for _, r := range regions {
		for _, id := range r.ids {
			if !known[id] {
				errs = append(errs, fmt.Sprintf("%s:%d: the journey-nodes marker lists %s, which is no requirement doc", file, r.from, id))
			}
		}
	}
	return errs
}

// listingErrors fails a doc that no region lists or that two regions list.
func listingErrors(file string, regions []region, id string) []string {
	var in []int
	for _, r := range regions {
		if slices.Contains(r.ids, id) {
			in = append(in, r.from)
		}
	}
	switch {
	case len(in) == 0:
		return []string{fmt.Sprintf("%s:%d: requirement doc %s is listed in no journey-nodes region; add it to a `%s` marker", file, regions[0].from, id, regionStart)}
	case len(in) > 1:
		return []string{fmt.Sprintf("%s:%d: requirement doc %s is listed in two journey-nodes regions; list it in one", file, in[1], id)}
	}
	return nil
}

// trimmed drops each line's indent and trailing space, so only the nodes count.
func trimmed(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = strings.TrimSpace(l)
	}
	return out
}

// RewriteJourneyRegion rewrites the marked region of each journey view in the package root dir from its docs.
// It writes nothing when a doc fails to load or lacks a title, and leaves a view without markers alone.
func RewriteJourneyRegion(root, dir string) error {
	docs, files, skipped := loadDocs(root, dir)
	if err := untitled(docs, skipped); err != nil {
		return err
	}
	for _, file := range journeyViews(dir, files) {
		if err := rewriteView(filepath.Join(root, filepath.FromSlash(file)), docs); err != nil {
			return err
		}
	}
	return nil
}

// untitled fails on a doc that did not load or lacks its ID or title.
func untitled(docs []requirement.Document, skipped []string) error {
	for _, doc := range docs {
		if doc.ID == "" || doc.Title == "" {
			skipped = append(skipped, doc.Path+": needs `# NNNN Title` before the region can be written")
		}
	}
	if len(skipped) > 0 {
		return errors.New(skipped[0])
	}
	return nil
}

// rewriteView writes each region of one view from docs, unless it already holds them.
func rewriteView(full string, docs []requirement.Document) error {
	body, err := os.ReadFile(full)
	if err != nil {
		return err
	}
	lines := strings.Split(string(body), "\n")
	regions, _, _ := findRegions(lines)
	var next []string
	last := 0
	for _, r := range regions {
		next = append(next, lines[last:r.from]...)
		next = append(next, regionLines(docs, r.ids, r.indent)...)
		last = r.to
	}
	next = append(next, lines[last:]...)
	if slices.Equal(next, lines) {
		return nil
	}
	return os.WriteFile(full, []byte(strings.Join(next, "\n")), 0o644)
}
