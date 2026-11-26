// Package requirements proves the pr-structure check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

// prStructureSummary is a real one-line blockquote summary.
const prStructureSummary = "> Ship the thing so a reader gets the value.\n\n"

// prStructureBackground is a filled `## Background` section.
const prStructureBackground = "## Background\n\nWhy this exists.\n\n"

// prStructureChangelog is a `## Changelog` section with one bullet.
const prStructureChangelog = "## Changelog\n\n- Feat: the new behavior\n"

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1.1.1 `pr-structure` passes it", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog)
		acceptance.Sees(t, out, code, 0, "pr-structure")
	})
	noSummary := prStructureBackground + prStructureChangelog
	t.Run("0001.1.1.2.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, noSummary)
		acceptance.Sees(t, out, code, 1)
	})
	t.Run("0001.1.1.2.2 it asks for a one-line summary", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, noSummary)
		acceptance.Sees(t, out, code, 1, "Missing pitch: add a one-line blockquote (> …) before the first ## section")
	})
	placeholder := "> REPLACE-ME\n\n" + prStructureBackground + prStructureChangelog
	t.Run("0001.1.1.3.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, placeholder)
		acceptance.Sees(t, out, code, 1)
	})
	t.Run("0001.1.1.3.2 it asks for the real summary", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, placeholder)
		acceptance.Sees(t, out, code, 1, "Pitch is still the template placeholder; write the real one-line value")
	})
	t.Run("0001.1.1.4.1 it fails naming the missing section", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureSummary+prStructureChangelog)
		acceptance.Sees(t, out, code, 1, "Missing `## Background` section")
	})
	t.Run("0001.1.1.5.1 it fails naming the empty section", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureSummary+prStructureBackground+"## Changelog\n\nNothing yet.\n")
		acceptance.Sees(t, out, code, 1, "`## Changelog` has no list item (- …)")
	})
	t.Run("0001.1.1.6.1 it fails naming the unexpected section", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog+"\n## Testing\n\nRan it.\n")
		acceptance.Sees(t, out, code, 1, "Unexpected `## Testing` section — a PR has only Background and Changelog")
	})
	t.Run("0001.1.1.7.1 it passes", func(t *testing.T) {
		t.Parallel()
		fence := strings.Repeat("`", 3)
		sample := "\n" + fence + "markdown\n## Testing\n" + fence + "\n"
		out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog+sample)
		acceptance.Sees(t, out, code, 0, "pr-structure")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	run := func(t *testing.T) (string, int) {
		return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "pr-structure")
	}
	t.Run("0001.2.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := run(t)
		acceptance.Sees(t, out, code, 0, "pr-structure")
	})
	t.Run("0001.2.1.1.2 it judges zero files", func(t *testing.T) {
		t.Parallel()
		out, code := run(t)
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWithNoFiles(t, out, "pr-structure")
	})
}

// prStructure runs `fitness-install -- pr-structure --body-file pr.md` on
// happyRepo with body as the PR description.
func prStructure(t *testing.T, body string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{"pr.md": body})
	return acceptance.Fitness(t, repo, nil, "pr-structure", "--body-file", "pr.md")
}
