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
	t.Run("0001.1.1 a summary plus background plus changelog bullet passes", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog)
		acceptance.Sees(t, out, code, 0, "pr-structure")
	})
	t.Run("0001.1.2 no blockquote summary fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureBackground+prStructureChangelog)
		acceptance.Sees(t, out, code, 1, "Missing pitch: add a one-line blockquote (> …) before the first ## section")
	})
	t.Run("0001.1.3 a placeholder summary fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, "> REPLACE-ME\n\n"+prStructureBackground+prStructureChangelog)
		acceptance.Sees(t, out, code, 1, "Pitch is still the template placeholder; write the real one-line value")
		if strings.Contains(out, "Missing pitch") {
			t.Errorf("a placeholder pitch was reported as missing, want it named as the placeholder:\n%s", out)
		}
	})
	t.Run("0001.1.4 no Background section fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureSummary+prStructureChangelog)
		acceptance.Sees(t, out, code, 1, "Missing `## Background` section")
	})
	t.Run("0001.1.5 a Changelog with no bullet fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureSummary+prStructureBackground+"## Changelog\n\nNothing yet.\n")
		acceptance.Sees(t, out, code, 1, "`## Changelog` has no list item (- …)")
	})
	t.Run("0001.1.6 an extra section fails", func(t *testing.T) {
		t.Parallel()
		out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog+"\n## Testing\n\nRan it.\n")
		acceptance.Sees(t, out, code, 1, "Unexpected `## Testing` section — a PR has only Background and Changelog")
	})
	t.Run("0001.1.7 an extra heading inside fenced code passes", func(t *testing.T) {
		t.Parallel()
		fence := strings.Repeat("`", 3)
		sample := "\n" + fence + "markdown\n## Testing\n" + fence + "\n"
		out, code := prStructure(t, prStructureSummary+prStructureBackground+prStructureChangelog+sample)
		acceptance.Sees(t, out, code, 0, "pr-structure")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "pr-structure")
	acceptance.Sees(t, out, code, 0, "pr-structure")
	acceptance.PassedWithNoFiles(t, out, "pr-structure")
}

// prStructure runs `fitness-install -- pr-structure --body-file pr.md` on
// happyRepo with body as the PR description.
func prStructure(t *testing.T, body string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{"pr.md": body})
	return acceptance.Fitness(t, repo, nil, "pr-structure", "--body-file", "pr.md")
}
