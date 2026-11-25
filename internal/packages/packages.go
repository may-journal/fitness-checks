// Package packages names the package roots of a monorepo run and the root
// that owns each file. A check judges each package like a repo of its own.
package packages

import (
	"slices"
	"strings"

	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

// Root is the repo root, which owns every file no package root holds.
const Root = "."

// Roots lists the roots this run judges, repo-relative: the repo root first,
// then each package root the packages key of .fitnessrc.json names. Without
// that key the repo root is the only one, and every check judges as before.
func Roots() []string {
	roots := []string{Root}
	for _, r := range checkkit.PackageRoots() {
		if r != Root && !slices.Contains(roots, r) {
			roots = append(roots, r)
		}
	}
	return roots
}

// Of is the deepest of roots that holds file, or the repo root.
func Of(roots []string, file string) string {
	owner := Root
	for _, r := range roots {
		if r != Root && strings.HasPrefix(file, r+"/") && len(r) > len(owner) {
			owner = r
		}
	}
	return owner
}

// Rel is file's path inside root.
func Rel(root, file string) string {
	if root == Root {
		return file
	}
	return strings.TrimPrefix(file, root+"/")
}

// Join is rel inside root, as a repo-relative path.
func Join(root, rel string) string {
	if root == Root {
		return rel
	}
	return root + "/" + rel
}

// Owned keeps the files whose owner among roots is root.
func Owned(roots []string, root string, files []string) []string {
	var out []string
	for _, f := range files {
		if Of(roots, f) == root {
			out = append(out, f)
		}
	}
	return out
}
