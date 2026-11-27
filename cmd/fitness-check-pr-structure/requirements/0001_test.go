// Package requirements proves the pr-structure check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
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
	t.Run("0001.1.1 a summary plus background plus changelog bullet passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.1.1", "f8f421a8", func(_ req.A) {
			out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog)
			acceptance.Sees(t, out, code, 0, "pr-structure")
		})
	})
	t.Run("0001.1.2 no blockquote summary fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureBackground+prStructureChangelog)
		req.Then(t, "0001.1.1.2.1", "07f1343a", func(a req.A) {
			a.Equal(1, code)
		})
		req.Then(t, "0001.1.1.2.2", "d44ca635", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "Missing pitch: add a one-line blockquote (> …) before the first ## section")
		})
	})
	t.Run("0001.1.3 a placeholder summary fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, "> REPLACE-ME\n\n"+prStructureBackground+prStructureChangelog)
		req.Then(t, "0001.1.1.3.1", "d2cd8eb7", func(a req.A) {
			a.Equal(1, code)
			a.Contains(out, "Pitch is still the template placeholder")
			a.True(!strings.Contains(out, "Missing pitch"), "a placeholder pitch was reported as missing, want it named as the placeholder:\n%s", out)
		})
		req.Then(t, "0001.1.1.3.2", "42205152", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "write the real one-line value")
		})
	})
	t.Run("0001.1.4 no Background section fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.4.1", "159ed4f1", func(_ req.A) {
			out, code := prStructure(t, prStructureSummary+prStructureChangelog)
			acceptance.Sees(t, out, code, 1, "Missing `## Background` section")
		})
	})
	t.Run("0001.1.5 a Changelog with no bullet fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.5.1", "9379e388", func(_ req.A) {
			out, code := prStructure(t, prStructureSummary+prStructureBackground+"## Changelog\n\nNothing yet.\n")
			acceptance.Sees(t, out, code, 1, "`## Changelog` has no list item (- …)")
		})
	})
	t.Run("0001.1.6 an extra section fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.6.1", "703fd84a", func(_ req.A) {
			out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog+"\n## Testing\n\nRan it.\n")
			acceptance.Sees(t, out, code, 1, "Unexpected `## Testing` section — a PR has only Background and Changelog")
		})
	})
	t.Run("0001.1.7 an extra heading inside fenced code passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.7.1", "07b8e753", func(_ req.A) {
			fence := strings.Repeat("`", 3)
			sample := "\n" + fence + "markdown\n## Testing\n" + fence + "\n"
			out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog+sample)
			acceptance.Sees(t, out, code, 0, "pr-structure")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "pr-structure")
	req.Then(t, "0001.2.1.1.1", "07b8e753", func(_ req.A) {
		acceptance.Sees(t, out, code, 0, "pr-structure")
	})
	req.Then(t, "0001.2.1.1.2", "41c3dea8", func(_ req.A) {
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
