package main

import (
	"fmt"
	"strings"

	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-runner/go/pkg/repocheck"
	"github.com/may-journal/fitness-runner/go/pkg/walkfs"
)

// reqPrefix starts the path of every file in a requirement's check folder,
// fitness/req-NNNN/.
const reqPrefix = repocheck.Folder + "/req-"

// reqCheckNumber is the requirement number of a file in a fitness/req-NNNN
// folder, or "".
func reqCheckNumber(file string) string {
	rest, ok := strings.CutPrefix(file, reqPrefix)
	n, _, nested := strings.Cut(rest, "/")
	if !ok || !nested || len(n) != 4 || strings.Trim(n, "0123456789") != "" {
		return ""
	}
	return n
}

// reqChecks lists the tracked files of every requirement check folder, and
// the requirement numbers they prove in each package root among roots.
func reqChecks(root string, roots []string) (files []string, proven map[string]map[string]bool) {
	proven = map[string]map[string]bool{}
	for _, f := range walkfs.FilesByExt(root, "") {
		owner := packages.Of(roots, f)
		if n := reqCheckNumber(packages.Rel(owner, f)); n != "" {
			files = append(files, f)
			if proven[owner] == nil {
				proven[owner] = map[string]bool{}
			}
			proven[owner][n] = true
		}
	}
	return files, proven
}

// orphanChecks fails each requirement check folder of package root pkg
// whose number none of its docs defines.
func orphanChecks(pkg string, docs []requirement.Document, proven map[string]bool) []string {
	defined := map[string]bool{}
	for _, d := range docs {
		defined[d.ID] = true
	}
	var errs []string
	for _, n := range sortedKeys(proven) {
		if !defined[n] {
			errs = append(errs, fmt.Sprintf("%s proves requirement %s, which no doc defines; write the doc or delete the check", packages.Join(pkg, reqPrefix+n), n))
		}
	}
	return errs
}
