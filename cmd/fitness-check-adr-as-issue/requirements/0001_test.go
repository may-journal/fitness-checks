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
func TestMain(m *testing.M) { acceptance.Main(m) }

func Test0001_1(t *testing.T) {
	t.Parallel()
	// 0001.1.4.1: a Markdown file sits in an adr folder.
	inFolder := func(t *testing.T) (string, int) {
		return adrAsIssue(t, map[string]string{
			"docs/architecture/adr/0001-use-sqlite.md": adrFile,
			"docs/decisions/ADR/0002-sync.md":          adrFile,
			"docs/architecture/adr/README.md":          "# ADRs\n\n- [Use SQLite](0001-use-sqlite.md)\n",
		})
	}
	t.Run("0001.1.4.1.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := inFolder(t)
		acceptance.Sees(t, out, code, 1, inADRFolder("docs/decisions/ADR/0002-sync.md"), inADRFolder("docs/architecture/adr/0001-use-sqlite.md"))
	})
	t.Run("0001.1.4.1.2 it names the rule it broke", func(t *testing.T) {
		t.Parallel()
		out, code := inFolder(t)
		acceptance.Sees(t, out, code, 1, "a Markdown file in an `adr` folder is an ADR")
	})
	t.Run("0001.1.4.1.3 it names `fitness adr-migrate`", func(t *testing.T) {
		t.Parallel()
		out, code := inFolder(t)
		acceptance.Sees(t, out, code, 1, "`fitness adr-migrate` opens one per file")
	})
	t.Run("0001.1.4.1.4 the folder's README passes", func(t *testing.T) {
		t.Parallel()
		out, code := inFolder(t)
		acceptance.Sees(t, out, code, 1)
		if strings.Contains(out, "adr/README.md:") {
			t.Errorf("the adr folder's README lists ADRs, yet it was reported:\n%s", out)
		}
	})
	// 0001.1.4.2: a file's first top-level heading starts with the word ADR.
	titled := func(t *testing.T) (string, int) {
		return adrAsIssue(t, map[string]string{
			"docs/storage.md":    "# ADR 12: Use SQLite\n\nWe sync offline.\n",
			"docs/adr-checks.md": "# adr-checks\n\nPairs each decision with a check.\n",
			"docs/guide.md":      "# Guide\n\n## ADR 3\n\nA later heading.\n",
		})
	}
	t.Run("0001.1.4.2.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := titled(t)
		acceptance.Sees(t, out, code, 1, "docs/storage.md: a title starting with ADR makes the file an ADR; "+toADRIssue)
	})
	t.Run("0001.1.4.2.2 a later ADR heading passes", func(t *testing.T) {
		t.Parallel()
		out, code := titled(t)
		acceptance.Sees(t, out, code, 1)
		if strings.Contains(out, "docs/guide.md:") {
			t.Errorf("docs/guide.md has no title starting with the word ADR, yet it was reported:\n%s", out)
		}
	})
	t.Run("0001.1.4.2.3 a word such as `adr-checks` passes", func(t *testing.T) {
		t.Parallel()
		out, code := titled(t)
		acceptance.Sees(t, out, code, 1)
		if strings.Contains(out, "docs/adr-checks.md:") {
			t.Errorf("docs/adr-checks.md has no title starting with the word ADR, yet it was reported:\n%s", out)
		}
	})
	// 0001.1.4.3: a file holds the three decision sections.
	sections := func(t *testing.T) (string, int) {
		return adrAsIssue(t, acceptance.With(acceptance.HeldCheck(), map[string]string{
			"docs/writing-checks.md": "# Writing checks\n\n```markdown\n" + acceptance.HeldDecision + "```\n",
		}))
	}
	t.Run("0001.1.4.3.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := sections(t)
		acceptance.Sees(t, out, code, 1, "fitness/adr-0001/README.md: `## Context`, `## Decision` and `## Consequences` sections make the file an ADR; "+toADRIssue)
	})
	t.Run("0001.1.4.3.2 the same sections in a code block pass", func(t *testing.T) {
		t.Parallel()
		out, code := sections(t)
		acceptance.Sees(t, out, code, 1)
		if strings.Contains(out, "docs/writing-checks.md:") {
			t.Errorf("an ADR in a code block is an example, yet it was reported:\n%s", out)
		}
	})
	// 0001.1.4.4: a file only points at an issue.
	linked := func(t *testing.T) (string, int) {
		return adrAsIssue(t, map[string]string{
			"docs/architecture/adr/0001-use-sqlite.md": "https://github.com/may-journal/example/issues/41\n",
			"fitness/adr-0001/README.md":               "# 0001 — Small files\n\nSmall files read fast.\n\nhttps://github.com/may-journal/example/issues/5\n\n## Check\n\nFails a long Markdown file.\n",
			"fitness/adr-0001/check":                   acceptance.VerdictScript(acceptance.Passing),
			".github/ISSUE_TEMPLATE/adr.md":            "---\nname: ADR\nlabels: ADR\n---\n\n## Status\n\n## Context\n\n## Decision\n\n## Consequences\n",
			"docs/notes.md":                            "# Notes\n",
		})
	}
	t.Run("0001.1.4.4.1 a file of one issue URL passes", func(t *testing.T) {
		t.Parallel()
		out, code := linked(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "0001-use-sqlite.md:") {
			t.Errorf("a file of one issue URL was reported:\n%s", out)
		}
	})
	t.Run("0001.1.4.4.2 a check README that links its issue passes", func(t *testing.T) {
		t.Parallel()
		out, code := linked(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "fitness/adr-0001/README.md:") {
			t.Errorf("a check README that links its issue was reported:\n%s", out)
		}
	})
	t.Run("0001.1.4.4.3 a file under `.github/` passes", func(t *testing.T) {
		t.Parallel()
		out, code := linked(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, ".github/ISSUE_TEMPLATE/adr.md:") {
			t.Errorf("a file under .github/ was reported:\n%s", out)
		}
	})
	t.Run("0001.1.4.4.4 it counts the paths it judged", func(t *testing.T) {
		t.Parallel()
		out, code := linked(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// The example repo's three files and the five written here are judged.
		acceptance.PassedWith(t, out, "adr-as-issue", "8")
	})
	// 0001.1.4.5: an allow entry names waiting ADR files with its reason.
	allowed := func(t *testing.T) (string, int) {
		return adrAsIssue(t, acceptance.With(acceptance.HeldCheck(), map[string]string{
			".fitnessrc.json": acceptance.AllowADRs(`{"path": "docs/architecture/adr/*.md", "reason": "the migration waits on repo access"}, ` +
				`{"path": "fitness/adr-*/README.md", "reason": "the checks move with them"}`),
			"docs/architecture/adr/0001-use-sqlite.md": adrFile,
			"docs/architecture/adr/0002-sync.md":       adrFile,
			"docs/decisions/adr/0003-theme.md":         adrFile,
		}))
	}
	t.Run("0001.1.4.5.1 those files pass", func(t *testing.T) {
		t.Parallel()
		out, code := allowed(t)
		acceptance.Sees(t, out, code, 1)
		for _, kept := range []string{"0001-use-sqlite.md", "0002-sync.md", "fitness/adr-0001/README.md"} {
			if strings.Contains(out, kept+":") {
				t.Errorf("%s is allowed, yet it was reported:\n%s", kept, out)
			}
		}
	})
	t.Run("0001.1.4.5.2 every other ADR file still fails", func(t *testing.T) {
		t.Parallel()
		out, code := allowed(t)
		acceptance.Sees(t, out, code, 1, inADRFolder("docs/decisions/adr/0003-theme.md"))
	})
	// 0001.1.4.6: an allow entry is invalid.
	invalid := func(t *testing.T) (string, int) {
		return adrAsIssue(t, map[string]string{
			".fitnessrc.json": acceptance.AllowADRs(`{"path": "docs/architecture/adr/*.md"}, {"path": "docs/old/adr/*.md", "reason": "moved"}, ` +
				`{"path": "docs/[adr/*.md", "reason": "broken"}`),
			"docs/architecture/adr/0001-use-sqlite.md": adrFile,
		})
	}
	t.Run("0001.1.4.6.1 an entry that lacks a reason fails naming its path", func(t *testing.T) {
		t.Parallel()
		out, code := invalid(t)
		acceptance.Sees(t, out, code, 1, `adrAsIssue.allow path "docs/architecture/adr/*.md" needs a reason`)
	})
	t.Run("0001.1.4.6.2 an entry that is no glob fails naming its path", func(t *testing.T) {
		t.Parallel()
		out, code := invalid(t)
		acceptance.Sees(t, out, code, 1, `adrAsIssue.allow path "docs/[adr/*.md" is not a glob: syntax error in pattern`)
	})
	t.Run("0001.1.4.6.3 an entry that matches no ADR file fails naming its path", func(t *testing.T) {
		t.Parallel()
		out, code := invalid(t)
		acceptance.Sees(t, out, code, 1, `adrAsIssue.allow path "docs/old/adr/*.md" matches no ADR file; remove it`)
	})
	t.Run("0001.1.4.7.1 it fails", func(t *testing.T) {
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
