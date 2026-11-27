// Package requirements proves the gitignore-why check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo with a .gitignore written over it.
package requirements

import (
	"os"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a repo with no gitignore passes without files", func(t *testing.T) {
		t.Parallel()
		out, code := gitignoreWhy(t, "")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "gitignore-why")
	})
	t.Run("0001.1.2 patterns that each have a comment above pass counting the file", func(t *testing.T) {
		t.Parallel()
		out, code := gitignoreWhy(t, "# dependencies\nnode_modules/\n\n# build output\ndist\n\n# keep this tracked\n!dist/keep.js\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "gitignore-why", "1")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	out, code := gitignoreWhy(t, "node_modules/\n")
	acceptance.Sees(t, out, code, 1, `.gitignore:1: pattern "node_modules/" `+noComment)
	out, code = gitignoreWhy(t, "# deps\nnode_modules/\ndist\n")
	acceptance.Sees(t, out, code, 1, `.gitignore:3: pattern "dist" `+noComment)
	out, code = gitignoreWhy(t, "# deps\n\nnode_modules/\n")
	acceptance.Sees(t, out, code, 1, `.gitignore:3: pattern "node_modules/" `+noComment)
	out, code = gitignoreWhy(t, "#\nnode_modules/\n")
	acceptance.Sees(t, out, code, 1, `.gitignore:2: pattern "node_modules/" `+noComment)
}

// noComment is the tail of every gitignore-why failure.
const noComment = "has no explanatory # comment on the line above"

// gitignoreWhy runs `fitness-install -- gitignore-why` on happyRepo with the
// given .gitignore; an empty one means the repo has none.
func gitignoreWhy(t *testing.T, gitignore string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{".gitignore": gitignore})
	return acceptance.Fitness(t, repo, nil, "gitignore-why")
}
