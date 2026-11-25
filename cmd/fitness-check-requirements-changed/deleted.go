package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/may-journal/fitness-checks/internal/packages"
)

// leftBehind fails the deletion of file when a requirement doc names it and neither that doc nor its tests changed.
func leftBehind(root, file string, tests testSet, changed map[string]bool) []string {
	docs := namingDocs(root, file, packages.Owned(tests.roots, packages.Of(tests.roots, file), tests.docs))
	proving := provingTests(docs, tests)
	if len(docs) == 0 || slices.ContainsFunc(slices.Concat(docs, proving), func(f string) bool { return changed[f] }) {
		return nil
	}
	return []string{fmt.Sprintf("%s was deleted, leaving %s unchanged; change or delete them with it", file, behind(docs, proving))}
}

// namingDocs lists the requirement docs among candidates whose text names
// file by its whole path.
func namingDocs(root, file string, candidates []string) []string {
	named := regexp.MustCompile(`(^|[^\w./-])` + regexp.QuoteMeta(file) + `\b`)
	var docs []string
	for _, d := range candidates {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(d)))
		if err == nil && named.Match(data) {
			docs = append(docs, d)
		}
	}
	return docs
}

// provingTests lists the test files that prove each doc's requirement.
func provingTests(docs []string, tests testSet) []string {
	var files []string
	for _, d := range docs {
		for _, f := range tests.numbers[docNumber(tests.roots, d)] {
			if !slices.Contains(files, f) {
				files = append(files, f)
			}
		}
	}
	return files
}

// behind names the docs and tests a deletion leaves behind.
func behind(docs, tests []string) string {
	text := "requirement doc " + joinList(docs)
	if len(tests) > 0 {
		text += " and test " + joinList(tests)
	}
	return text
}
