// Package requirements proves the adr-as-issue check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a file in an adr folder fails", func(t *testing.T) {
		t.Parallel()
		out, code := adrAsIssue(t, map[string]string{
			"docs/architecture/adr/0001-use-sqlite.md": adrFile,
			"docs/decisions/ADR/0002-sync.md":          adrFile,
			"docs/architecture/adr/README.md":          "# ADRs\n\n- [Use SQLite](0001-use-sqlite.md)\n",
		})
		acceptance.Sees(t, out, code, 1, inADRFolder("docs/decisions/ADR/0002-sync.md"), inADRFolder("docs/architecture/adr/0001-use-sqlite.md"))
		if strings.Contains(out, "adr/README.md:") {
			t.Errorf("the adr folder's README lists ADRs, yet it was reported:\n%s", out)
		}
	})
	t.Run("0001.1.2 a file titled as an ADR fails", func(t *testing.T) {
		t.Parallel()
		out, code := adrAsIssue(t, map[string]string{
			"docs/storage.md":    "# ADR 12: Use SQLite\n\nWe sync offline.\n",
			"docs/adr-checks.md": "# adr-checks\n\nPairs each decision with a check.\n",
			"docs/guide.md":      "# Guide\n\n## ADR 3\n\nA later heading.\n",
		})
		acceptance.Sees(t, out, code, 1, "docs/storage.md: a title starting with ADR makes the file an ADR; "+toADRIssue)
		for _, kept := range []string{"docs/adr-checks.md", "docs/guide.md"} {
			if strings.Contains(out, kept+":") {
				t.Errorf("%s has no title starting with the word ADR, yet it was reported:\n%s", kept, out)
			}
		}
	})
	t.Run("0001.1.3 a file holding the decision sections fails", func(t *testing.T) {
		t.Parallel()
		out, code := adrAsIssue(t, acceptance.With(acceptance.HeldCheck(), map[string]string{
			"docs/writing-checks.md": "# Writing checks\n\n```markdown\n" + acceptance.HeldDecision + "```\n",
		}))
		acceptance.Sees(t, out, code, 1, "fitness/adr-0001/README.md: `## Context`, `## Decision` and `## Consequences` sections make the file an ADR; "+toADRIssue)
		if strings.Contains(out, "docs/writing-checks.md:") {
			t.Errorf("an ADR in a code block is an example, yet it was reported:\n%s", out)
		}
	})
	t.Run("0001.1.4 every ADR kept as a linked issue passes", func(t *testing.T) {
		t.Parallel()
		out, code := adrAsIssue(t, map[string]string{
			"docs/architecture/adr/0001-use-sqlite.md": "https://github.com/may-journal/example/issues/41\n",
			"fitness/adr-0001/README.md":               "# 0001 — Small files\n\nSmall files read fast.\n\nhttps://github.com/may-journal/example/issues/5\n\n## Check\n\nFails a long Markdown file.\n",
			"fitness/adr-0001/check":                   acceptance.VerdictScript(acceptance.Passing),
			".github/ISSUE_TEMPLATE/adr.md":            "---\nname: ADR\nlabels: ADR\n---\n\n## Status\n\n## Context\n\n## Decision\n\n## Consequences\n",
			"docs/notes.md":                            "# Notes\n",
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// The example repo's three files and the five written here are judged.
		acceptance.PassedWith(t, out, "adr-as-issue", "8")
	})
	t.Run("0001.1.5 allowed ADR files with reasons wait", func(t *testing.T) {
		t.Parallel()
		out, code := adrAsIssue(t, acceptance.With(acceptance.HeldCheck(), map[string]string{
			".fitnessrc.json": acceptance.AllowADRs(`{"path": "docs/architecture/adr/*.md", "reason": "the migration waits on repo access"}, ` +
				`{"path": "fitness/adr-*/README.md", "reason": "the checks move with them"}`),
			"docs/architecture/adr/0001-use-sqlite.md": adrFile,
			"docs/architecture/adr/0002-sync.md":       adrFile,
			"docs/decisions/adr/0003-theme.md":         adrFile,
		}))
		acceptance.Sees(t, out, code, 1, inADRFolder("docs/decisions/adr/0003-theme.md"))
		for _, kept := range []string{"0001-use-sqlite.md", "0002-sync.md", "fitness/adr-0001/README.md"} {
			if strings.Contains(out, kept+":") {
				t.Errorf("%s is allowed, yet it was reported:\n%s", kept, out)
			}
		}
	})
	t.Run("0001.1.6 an allow entry without reason, glob or match fails", func(t *testing.T) {
		t.Parallel()
		out, code := adrAsIssue(t, map[string]string{
			".fitnessrc.json": acceptance.AllowADRs(`{"path": "docs/architecture/adr/*.md"}, {"path": "docs/old/adr/*.md", "reason": "moved"}, ` +
				`{"path": "docs/[adr/*.md", "reason": "broken"}`),
			"docs/architecture/adr/0001-use-sqlite.md": adrFile,
		})
		acceptance.Sees(t, out, code, 1,
			`adrAsIssue.allow path "docs/architecture/adr/*.md" needs a reason`,
			`adrAsIssue.allow path "docs/old/adr/*.md" matches no ADR file; remove it`,
			`adrAsIssue.allow path "docs/[adr/*.md" is not a glob: syntax error in pattern`,
			inADRFolder("docs/architecture/adr/0001-use-sqlite.md"),
		)
	})
	t.Run("0001.1.7 a repo git cannot list fails", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", nil)
		acceptance.MustDo(t, os.WriteFile(filepath.Join(repo, ".git", "index"), []byte("not an index\n"), 0o644))
		out, code := acceptance.Fitness(t, repo, nil, "adr-as-issue")
		acceptance.Sees(t, out, code, 1, "cannot list the repo's files: git cannot list the tracked files")
	})
}

// adrAsIssue runs `fitness-install -- adr-as-issue` on happyRepo with files
// written over it.
func adrAsIssue(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "adr-as-issue")
}

// toADRIssue is the fix every ADR file's error names.
const toADRIssue = "open an `ADR` issue instead; `fitness adr-migrate` opens one per file and points each link at it"

// adrFile is an ADR as a file in an `adr` folder holds it.
const adrFile = "# ADR 0001 — Use SQLite\n\nStatus: accepted\n\n## Context\n\nWe sync offline.\n"

// inADRFolder is the error line for an ADR file in an `adr` folder.
func inADRFolder(path string) string {
	return path + ": a Markdown file in an `adr` folder is an ADR; " + toADRIssue
}
