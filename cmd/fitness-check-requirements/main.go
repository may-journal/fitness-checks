package main

import (
	"errors"
	"slices"
	"sort"
	"strings"

	"github.com/may-journal/fitness-checks/internal/config"
	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/lang"
	"github.com/may-journal/fitness-runner/go/pkg/lang/langs"
	"github.com/may-journal/fitness-runner/go/pkg/walkfs"
)

func main() {
	program.Configured(docSource, config.Keys(defaults, configSource), checkkit.Check{
		Describe: checkkit.Describe{
			Name:     "requirements",
			Scope:    checkkit.ScopeFiles,
			Tests:    true,
			Requires: []checkkit.Requirement{checkkit.RequiresRepo},
		},
		Run: run,
	})
}

// docsDir holds the requirement docs, one per file.
const docsDir = "docs/requirements/"

// fixFlag makes the check delete every leftover Source line before it
// judges, so `fitness requirements --fix` converts a repo's docs.
const fixFlag = "--fix"

func run(root string, args []string) (checkkit.Result, error) {
	roots := packages.Roots()
	paths := docPaths(root, roots)
	if err := dropSources(root, paths, args); err != nil {
		return checkkit.Result{}, err
	}
	tests, files, written, err := collectTests(root, langs.All)
	if err != nil || len(written) == 0 {
		return checkkit.Pass(0), err
	}
	checks, proven := reqChecks(root, roots)
	errs, notes, err := packageErrors(root, roots, paths, tests, written, proven)
	if err != nil {
		return checkkit.Result{}, err
	}
	pairs, err := pairErrors(root, roots, paths, append(files, checks...), tests)
	res := result(len(paths)+len(files), append(errs, pairs...))
	res.Info = notes
	return res, err
}

// pkg is one package root with the docs, tests and requirement checks it
// holds, judged like a repo of its own.
type pkg struct {
	root   string
	paths  []string
	tests  []reqTest
	proven map[string]bool
	// alone is true for the repo root of a repo that names no packages.
	alone bool
}

// split gives each package root the docs, tests and checks it owns.
func split(roots, paths []string, tests []reqTest, proven map[string]map[string]bool) []pkg {
	var out []pkg
	for _, r := range roots {
		p := pkg{root: r, paths: packages.Owned(roots, r, paths), proven: proven[r], alone: len(roots) == 1}
		for _, t := range tests {
			if packages.Of(roots, t.File) == r {
				p.tests = append(p.tests, t)
			}
		}
		out = append(out, p)
	}
	return out
}

// packageErrors judges each package root's docs and tests on their own, then
// the test tools of every test.
func packageErrors(root string, roots, paths []string, tests []reqTest, written []lang.Module, proven map[string]map[string]bool) (errs, notes []string, err error) {
	limit, err := maxAcceptances(root)
	if err != nil {
		return nil, nil, err
	}
	whole := true
	for _, p := range split(roots, paths, tests, proven) {
		found, skipped, judged := docErrors(root, p, written, limit)
		errs, notes, whole = append(errs, found...), append(notes, skipped...), whole && judged
	}
	if !whole {
		return errs, notes, nil
	}
	tooled, err := toolErrors(root, tests, written)
	return append(errs, tooled...), notes, err
}

// docErrors judges one package's doc shapes, IDs, shared areas and test ownership, noting each unreadable doc and reporting a whole read.
// A doc its proven numbers name needs no test, since its fitness/req-NNNN check proves it.
func docErrors(root string, p pkg, written []lang.Module, limit int) (errs, notes []string, whole bool) {
	docs, errs, skipped := requirement.Load(root, p.paths)
	notes = skipped
	if len(skipped) > 0 {
		// A doc left out would leave its tests unowned and its IDs unknown.
		return errs, append(notes, "ownership and history were not judged without every doc"), false
	}
	errs = append(errs, missingDocs(p)...)
	hist, historyNote := readHistory(root, packages.Join(p.root, docsDir), p.paths)
	if historyNote != "" {
		notes = append(notes, historyNote)
	}
	errs = append(errs, idErrors(docs, hist)...)
	errs = append(errs, capErrors(docs, limit)...)
	errs = append(errs, overlapErrors(docs)...)
	errs = append(errs, ownershipErrors(docs, p.tests, written, p.proven)...)
	errs = append(errs, stepCaseErrors(docs, p.tests)...)
	return append(errs, orphanChecks(p.root, docs, p.proven)...), notes, true
}

// dropSources deletes the leftover Source line from each doc at paths when
// args ask for the fix.
func dropSources(root string, paths, args []string) error {
	if !slices.Contains(args, fixFlag) {
		return nil
	}
	docs, _, _ := requirement.Load(root, paths)
	var err error
	for _, d := range docs {
		err = errors.Join(err, requirement.DropSource(root, d))
	}
	return err
}

// result passes with no errors and fails with any.
func result(files int, errs []string) checkkit.Result {
	if len(errs) > 0 {
		return checkkit.Fail(files, errs...)
	}
	return checkkit.Pass(files)
}

// collectTests reads every module's tests and their files in registry order, and lists the modules whose language the repo uses.
// The scan is the one requirements-changed reads, shared by the runner when both run.
func collectTests(root string, modules []lang.Module) ([]reqTest, []string, []lang.Module, error) {
	scan, err := lang.ScanTests(root, modules)
	if err != nil {
		return nil, nil, nil, err
	}
	named := map[string]lang.Module{}
	for _, m := range modules {
		named[m.Name] = m
	}
	var tests []reqTest
	for _, t := range scan.Tests {
		tests = append(tests, reqTest{Test: t.Test, module: named[t.Module]})
	}
	var written []lang.Module
	for _, name := range scan.Written {
		written = append(written, named[name])
	}
	return tests, scan.Files, written, nil
}

// docPaths lists every package root's requirement docs. A template.md
// beside them is the doc-template check's concern, not a requirement.
func docPaths(root string, roots []string) []string {
	var out []string
	for _, f := range walkfs.FilesByExt(root, ".md") {
		if docIn(roots, f) {
			out = append(out, f)
		}
	}
	return out
}

// docIn reports whether f is a requirement doc in the docs folder of the
// package root among roots that owns it.
func docIn(roots []string, f string) bool {
	return strings.HasPrefix(packages.Rel(packages.Of(roots, f), f), docsDir) && isDoc(f)
}

// missingDocs fails a package root with code that has no requirement docs.
// The repo root of a monorepo may leave its behavior to its packages.
func missingDocs(p pkg) []string {
	if len(p.paths) > 0 || p.root == packages.Root && !p.alone {
		return nil
	}
	return []string{packages.Join(p.root, docsDir) + " has no requirement docs; add NNNN-kebab-title.md files from the template in the fitness-check-requirements README"}
}

// idErrors flags a requirement ID that two docs share, and any doc or
// acceptance ID that history shows was removed before.
func idErrors(docs []requirement.Document, hist history) []string {
	var errs []string
	owners := map[string][]string{}
	for _, d := range docs {
		owners[d.ID] = append(owners[d.ID], d.Path)
		errs = append(errs, reusedErrors(d, hist)...)
	}
	for _, id := range sortedKeys(owners) {
		if len(owners[id]) > 1 {
			errs = append(errs, "requirement "+id+" is used by "+strings.Join(owners[id], ", ")+"; give each requirement its own ID")
		}
	}
	return errs
}

// reusedErrors flags the doc's IDs that a past commit deleted.
func reusedErrors(d requirement.Document, hist history) []string {
	var errs []string
	for _, id := range d.IDs() {
		if hist.removed[id.ID] {
			errs = append(errs, requirement.At(d.Path, id.Line, id.ID+" was deleted before; IDs are never reused, so "+reuseFix(id.ID)))
		}
	}
	return errs
}

// reuseFix says how to add back what a deleted ID held. Steps are numbered in
// order, so a deleted step returns only under a new acceptance.
func reuseFix(id string) string {
	if strings.Count(id, ".") == 2 {
		return "chain the step under a new acceptance ID"
	}
	return "take the next free number"
}

// sortedKeys returns m's keys in order, for stable output.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
