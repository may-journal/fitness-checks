package main

import (
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-runner/go/pkg/comments"
	"github.com/may-journal/fitness-runner/go/pkg/lang"
	"github.com/may-journal/fitness-runner/go/pkg/lang/langs"
	"github.com/may-journal/fitness-runner/go/pkg/walkfs"
)

// docsDir holds the requirement docs, one per file.
const docsDir = "docs/requirements/"

// numberedDoc matches a requirement doc file: 0014-kebab-title.md.
var numberedDoc = regexp.MustCompile(`^(\d{4})-.+\.md$`)

// docNumber returns the requirement number of a doc in its package root's
// docs folder, qualified by that root, or "". The repo root's numbers stay
// bare.
func docNumber(roots []string, f string) string {
	owner := packages.Of(roots, f)
	if !strings.HasPrefix(packages.Rel(owner, f), docsDir) {
		return ""
	}
	if m := numberedDoc.FindStringSubmatch(path.Base(f)); m != nil {
		return packages.Join(owner, m[1])
	}
	return ""
}

// requirementDocs lists the requirement docs of every package root.
func requirementDocs(root string, roots []string) []string {
	var docs []string
	for _, f := range walkfs.FilesByExt(root, ".md") {
		if docNumber(roots, f) != "" {
			docs = append(docs, f)
		}
	}
	return docs
}

// productExts are the endings of product code: Go, Swift, and shell scripts.
var productExts = []string{".go", ".swift", ".sh"}

// notProductDirs are folders whose files are vendored, fixtures, or build
// output, never product code.
var notProductDirs = []string{"vendor", "third_party", "Pods", "Carthage", "node_modules", "testdata", "dist", ".build", "DerivedData"}

// generatedGo matches the line Go's convention puts in a generated file.
var generatedGo = regexp.MustCompile(`(?m)^// Code generated .* DO NOT EDIT\.$`)

// testSet holds what the language modules know of the repo's tests: files that declare tests, and test files per requirement number.
// It keeps the package roots and requirement docs the tests pair with, since a test pairs only within its package.
type testSet struct {
	files   map[string]bool
	numbers map[string][]string
	roots   []string
	docs    []string
}

// scanTests reads every language module's tests in the scan the
// requirements check reads, shared by the runner when both run.
func scanTests(root string, roots, docs []string) (testSet, error) {
	s := testSet{files: map[string]bool{}, numbers: map[string][]string{}, roots: roots, docs: docs}
	scan, err := lang.ScanTests(root, langs.All)
	if err != nil {
		return s, err
	}
	for _, f := range scan.Files {
		s.files[f] = true
		s.prove(f, lang.TestFileNumbers(langs.All, f)...)
	}
	for _, t := range scan.Tests {
		s.prove(t.File, idNumber(t.ID)...)
	}
	return s, nil
}

// idNumber is the requirement number an acceptance ID carries, if any.
func idNumber(id string) []string {
	if len(id) < 4 {
		return nil
	}
	return []string{id[:4]}
}

// prove records that file proves each of numbers in its package root.
func (s testSet) prove(file string, numbers ...string) {
	owner := packages.Of(s.roots, file)
	for _, n := range numbers {
		n = packages.Join(owner, n)
		if !slices.Contains(s.numbers[n], file) {
			s.numbers[n] = append(s.numbers[n], file)
		}
	}
}

// isTest reports whether f is a test: a module's test file, a Go test file, or a Swift file in a `…Tests` folder.
func (s testSet) isTest(f string) bool {
	return s.files[f] || lang.IsTestSource(langs.All, f) || swiftTestTarget(f)
}

// swiftTestTarget reports whether a Swift file sits in a folder, or has a
// name, that ends in Tests, as Xcode names test targets.
func swiftTestTarget(f string) bool {
	if path.Ext(f) != ".swift" {
		return false
	}
	return slices.ContainsFunc(strings.Split(strings.TrimSuffix(f, ".swift"), "/"), func(part string) bool {
		return strings.HasSuffix(part, "Tests")
	})
}

// anyTest reports whether a test in package root pkg changed.
func (s testSet) anyTest(pkg string, changed map[string]bool) bool {
	for f := range changed {
		if s.isTest(f) && packages.Of(s.roots, f) == pkg {
			return true
		}
	}
	return false
}

// anyDoc reports whether a requirement doc in package root pkg changed.
func (s testSet) anyDoc(pkg string, changed map[string]bool) bool {
	for f := range changed {
		if docNumber(s.roots, f) != "" && packages.Of(s.roots, f) == pkg {
			return true
		}
	}
	return false
}

// documented reports whether package root pkg holds requirement docs.
func (s testSet) documented(pkg string) bool {
	return len(packages.Owned(s.roots, pkg, s.docs)) > 0
}

// productChanges keeps the changes to product code: Go, Swift, and shell
// files that are not tests, vendored, generated, or build output.
func productChanges(root string, tests testSet, changes []change) []change {
	var candidates []string
	for _, c := range changes {
		if isProductPath(c.path) && !tests.isTest(c.path) && !isGeneratedGo(root, c.path) {
			candidates = append(candidates, c.path)
		}
	}
	kept := walkfs.WithoutGenerated(root, candidates)
	build := lang.BuildOutput(langs.All, root)
	return slices.DeleteFunc(slices.Clone(changes), func(c change) bool {
		return !slices.Contains(kept, c.path) || slices.Contains(build, c.path)
	})
}

// isProductPath reports whether f has a product ending outside every
// vendored, fixture, and build output folder.
func isProductPath(f string) bool {
	if !slices.Contains(productExts, path.Ext(f)) {
		return false
	}
	return !slices.ContainsFunc(strings.Split(path.Dir(f), "/"), func(dir string) bool {
		return slices.Contains(notProductDirs, dir)
	})
}

// isGeneratedGo reports whether a Go file's code view keeps the line that
// marks it generated.
func isGeneratedGo(root, f string) bool {
	if path.Ext(f) != ".go" {
		return false
	}
	data, err := comments.CodeFile(root, f)
	return err == nil && generatedGo.Match(data)
}
