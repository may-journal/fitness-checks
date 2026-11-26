// Package requirements proves the no-plans-dir check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// Where the check sends each kind of temporary work.
const (
	toPlanIssue  = "open a `Plan`-labeled GitHub Issue and put the plan there, not a file in the repo"
	toPlanDesign = "post it as a comment on the Plan Issue that it designs, not a file in the repo"
	toPRComment  = "post it as a comment on the pull request it reviews, not a file in the repo"
)

// planChecklist is a Plan's checklist as a Markdown file would hold it.
const planChecklist = "# Ship it\n\n## What needs to happen\n\n- [ ] Write the tests.\n"

// plans runs the check on happyRepo with files written over it.
func plans(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "no-plans-dir")
}

// inPlansFolder is the error line for a file in a `plans` folder.
func inPlansFolder(path string) string {
	return path + ": a `plans` folder holds temporary work; " + toPlanIssue
}

// inOrder asserts each want appears in out after the one before it.
func inOrder(t *testing.T, out string, want ...string) {
	t.Helper()
	text, at := acceptance.Read(out), 0
	for _, w := range want {
		i := strings.Index(text[at:], w)
		if i < 0 {
			t.Errorf("output is missing %q after position %d, or lists it out of order:\n%s", w, at, out)
			return
		}
		at += i + len(w)
	}
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 it fails naming the path", func(t *testing.T) {
		t.Parallel()
		out, code := plans(t, map[string]string{"docs/plans/03-publish.md": "# Publish\n"})
		acceptance.Sees(t, out, code, 1, "docs/plans/03-publish.md: a `plans` folder")
		out, code = plans(t, map[string]string{"docs/plans/archive/old.md": "# Old\n"})
		acceptance.Sees(t, out, code, 1, "docs/plans/archive/old.md: a `plans` folder")
	})
	t.Run("0001.1.2.1.2 it names a Plan issue", func(t *testing.T) {
		t.Parallel()
		out, code := plans(t, map[string]string{"docs/plans/03-publish.md": "# Publish\n"})
		acceptance.Sees(t, out, code, 1, inPlansFolder("docs/plans/03-publish.md"))
		out, code = plans(t, map[string]string{"docs/plans/archive/old.md": "# Old\n"})
		acceptance.Sees(t, out, code, 1, inPlansFolder("docs/plans/archive/old.md"))
	})
	keptFiles := map[string]string{
		"docs/notes.md":         "# Notes\n",
		"assets/screen.png":     "png",
		"docs/planner-guide.md": "# Planner guide\n\nA Plan lists its steps under a heading.\n",
		"go/plan_test.go":       "package plan\n\nconst body = \"## What needs to happen\\n\\n- [ ] Ship it.\\n\"\n",
	}
	t.Run("0001.1.2.2.1 they pass", func(t *testing.T) {
		t.Parallel()
		out, code := plans(t, keptFiles)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.2.2 every file stays in place", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", keptFiles)
		out, code := acceptance.Fitness(t, repo, nil, "no-plans-dir")
		acceptance.Sees(t, out, code, 0)
		for name := range keptFiles {
			if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(name))); err != nil {
				t.Errorf("no-plans-dir must leave %s in place: %v", name, err)
			}
		}
	})
	artifactFolders := map[string]string{
		"src/Plans/next.md":               "# Next\n",
		"design/mockups/home.png":         "png",
		"reports/validation/run.txt":      "ok",
		"review-artifacts/diff.txt":       "diff",
		"assets/icons/validation-ok.png":  "png",
		"docs/validation-checklist.md":    "# Validation\n",
		"tests/fixtures/review/input.txt": "in",
	}
	t.Run("0001.1.2.3.1 it fails in path order", func(t *testing.T) {
		t.Parallel()
		out, code := plans(t, artifactFolders)
		acceptance.Sees(t, out, code, 1)
		inOrder(t, out, "design/mockups/home.png:", "reports/validation/run.txt:", "review-artifacts/diff.txt:", "src/Plans/next.md:")
		for _, kept := range []string{"assets/icons/validation-ok.png", "docs/validation-checklist.md", "tests/fixtures/review/input.txt"} {
			if strings.Contains(out, kept+":") {
				t.Errorf("%s names no artifact folder, yet it was reported:\n%s", kept, out)
			}
		}
	})
	t.Run("0001.1.2.3.2 it names the comment it belongs in", func(t *testing.T) {
		t.Parallel()
		out, code := plans(t, artifactFolders)
		acceptance.Sees(t, out, code, 1,
			"design/mockups/home.png: a `mockups` folder holds temporary work; "+toPlanDesign,
			"reports/validation/run.txt: a `validation` folder holds temporary work; "+toPRComment,
			"review-artifacts/diff.txt: a `review-artifacts` folder holds temporary work; "+toPRComment,
			"src/Plans/next.md: a `plans` folder holds temporary work; "+toPlanIssue,
		)
	})
	artifactNames := map[string]string{
		"docs/launch-plan.md":  "# Launch\n",
		"img/home-mockup.png":  "png",
		"pr-42-before.png":     "png",
		"docs/planner.md":      "# Planner\n",
		"img/print-before.png": "png",
	}
	t.Run("0001.1.2.4.1 it fails naming the rule", func(t *testing.T) {
		t.Parallel()
		out, code := plans(t, artifactNames)
		acceptance.Sees(t, out, code, 1,
			"docs/launch-plan.md: a `*-plan.md` file is a plan;",
			"img/home-mockup.png: a `*-mockup.*` file is a design mockup;",
			"pr-42-before.png: a `pr-<number>-*` file belongs to one pull request;",
		)
		for _, kept := range []string{"docs/planner.md", "img/print-before.png"} {
			if strings.Contains(out, kept+":") {
				t.Errorf("%s matches no artifact name, yet it was reported:\n%s", kept, out)
			}
		}
	})
	t.Run("0001.1.2.4.2 it names the comment it belongs in", func(t *testing.T) {
		t.Parallel()
		out, code := plans(t, artifactNames)
		acceptance.Sees(t, out, code, 1,
			"docs/launch-plan.md: a `*-plan.md` file is a plan; "+toPlanIssue,
			"img/home-mockup.png: a `*-mockup.*` file is a design mockup; "+toPlanDesign,
			"pr-42-before.png: a `pr-<number>-*` file belongs to one pull request; "+toPRComment,
		)
	})
	t.Run("0001.1.2.5.1 it fails naming the Plan issue", func(t *testing.T) {
		t.Parallel()
		out, code := plans(t, map[string]string{
			"docs/ship.md":                   planChecklist,
			".github/ISSUE_TEMPLATE/plan.md": planChecklist,
			"docs/what-needs-to-happen.md":   "# Notes\n\n## What needs to happen\n\nNothing yet.\n",
		})
		acceptance.Sees(t, out, code, 1, "docs/ship.md: a `## What needs to happen` checklist is a plan; "+toPlanIssue)
		for _, kept := range []string{".github/ISSUE_TEMPLATE/plan.md", "docs/what-needs-to-happen.md"} {
			if strings.Contains(out, kept+":") {
				t.Errorf("%s holds no Plan checklist outside .github/, yet it was reported:\n%s", kept, out)
			}
		}
	})
	t.Run("0001.1.2.6.1 it fails before the commit", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{"docs/home.png": "png"})
		acceptance.Write(t, repo, map[string]string{"docs/plans/next.md": "# Next\n"})
		acceptance.Git(t, repo, "add", "docs/plans/next.md")
		acceptance.MustDo(t, os.MkdirAll(filepath.Join(repo, "docs", "mockups"), 0o755))
		acceptance.Git(t, repo, "mv", "docs/home.png", "docs/mockups/home.png")
		out, code := acceptance.Fitness(t, repo, nil, "no-plans-dir")
		acceptance.Sees(t, out, code, 1,
			"docs/mockups/home.png: a `mockups` folder holds temporary work; "+toPlanDesign,
			inPlansFolder("docs/plans/next.md"),
		)
	})
	t.Run("0001.1.2.7.1 it fails saying it cannot list them", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", nil)
		acceptance.MustDo(t, os.WriteFile(filepath.Join(repo, ".git", "index"), []byte("not an index"), 0o644))
		out, code := acceptance.Fitness(t, repo, nil, "no-plans-dir")
		acceptance.Sees(t, out, code, 1, "cannot list the repo's files")
	})
}
