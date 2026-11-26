package main

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-checks/internal/requirement"
)

// checkPrefix starts the folder of each check a monorepo builds.
const checkPrefix = "fitness-check-"

// nameLike matches a hyphenated name, or any name in backticks, as a doc writes a check's name.
var nameLike = regexp.MustCompile("`([a-z][a-z0-9-]*)`|\\b([a-z][a-z0-9]*(?:-[a-z0-9]+)+)\\b")

// checkNames lists the checks beside the repo root of a monorepo, one per package folder, or none for another root.
func checkNames(roots []string, root string) map[string]bool {
	if root != packages.Root || len(roots) == 1 {
		return nil
	}
	names := map[string]bool{}
	for _, r := range roots {
		if r != packages.Root {
			names[strings.TrimPrefix(path.Base(r), checkPrefix)] = true
		}
	}
	return names
}

// scopeErrors flags a root doc acceptance of a monorepo that names a single check, since that check's own doc holds it.
func scopeErrors(docs []requirement.Document, names map[string]bool) []string {
	var errs []string
	for _, d := range docs {
		for _, a := range d.Acceptances {
			if named := checkNamed(a, names); named != "" {
				errs = append(errs, fmt.Sprintf("%s:%d: acceptance %s names the %s check; move it to that check's own requirement doc", d.Path, a.Line, a.ID, named))
			}
		}
	}
	return errs
}

// checkNamed is the first check an acceptance's name, GIVEN, WHEN or THEN lines name, or empty.
func checkNamed(a requirement.IDLine, names map[string]bool) string {
	text := []string{a.Name, a.Given}
	for _, s := range a.Steps {
		text = append(text, s.When, s.Then)
	}
	for _, m := range nameLike.FindAllStringSubmatch(strings.Join(text, "\n"), -1) {
		if name := m[1] + m[2]; names[name] {
			return name
		}
	}
	return ""
}
