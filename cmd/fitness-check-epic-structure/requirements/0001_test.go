// Package requirements proves the epic-structure check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo with the epic body saved over it.
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
	t.Run("0001.1.1 a pitch and Background and a list of references pass", func(t *testing.T) {
		t.Parallel()
		out, code := epicStructure(t, nil, pitch+background+list("1. #22\n2. #23, after #22\n3. org/repo#24"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		repo := acceptance.Example(t, "happyRepo", nil)
		out, code = acceptance.Fitness(t, repo, nil, "epic-structure", "--body", pitch+background+list("- #22"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
	})
	t.Run("0001.1.2 a missing pitch or section fails naming it", func(t *testing.T) {
		t.Parallel()
		out, code := epicStructure(t, nil, background+list("- #1"))
		acceptance.Sees(t, out, code, 1, "Missing pitch: add a one-line blockquote (> …) before the first ## section")
		out, code = epicStructure(t, nil, pitch+list("- #1"))
		acceptance.Sees(t, out, code, 1, "Missing `## Background` section")
		out, code = epicStructure(t, nil, pitch+background)
		acceptance.Sees(t, out, code, 1, "Missing `## What needs to happen` section")
	})
	t.Run("0001.1.3 another section fails naming it", func(t *testing.T) {
		t.Parallel()
		out, code := epicStructure(t, nil, pitch+background+"## Notes\n\nExtra.\n\n"+list("- #1"))
		acceptance.Sees(t, out, code, 1, "Unexpected `## Notes` section — a epic has only Background and What needs to happen")
	})
	t.Run("0001.1.4 the repo's own Epic template judges the sections", func(t *testing.T) {
		t.Parallel()
		body := pitch + background + "## Rollout\n\nIn steps.\n\n" + list("- #1")
		out, code := epicStructure(t, nil, body)
		acceptance.Sees(t, out, code, 1, "Unexpected `## Rollout` section")
		own := "---\nname: Epic\nabout: An epic\nlabels: Epic\n---\n\n> One line pitch.\n\n## Background\n\n<!-- why -->\n\nText\n\n" +
			"## Rollout\n\n<!-- how -->\n\nText\n\n## What needs to happen\n\n<!-- list -->\n\n- #1\n"
		out, code = epicStructure(t, map[string]string{".github/ISSUE_TEMPLATE/epic.md": own}, body)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.5 headings and items in fenced code never count", func(t *testing.T) {
		t.Parallel()
		sample := "```markdown\n## Extra\n- [ ] no\n```\n\n"
		out, code := epicStructure(t, nil, pitch+"## Background\n\nWhy.\n\n"+sample+list("- #1"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a checkbox or emoji or status word or other text fails naming the item", func(t *testing.T) {
		t.Parallel()
		out, code := epicStructure(t, nil, pitch+background+list("- [ ] #1"))
		acceptance.Sees(t, out, code, 1, "item 1 uses a checkbox; list the issue as a plain reference, and GitHub shows whether it is done")
		out, code = epicStructure(t, nil, pitch+background+list("- #1 ✅\n- #2, blocked\n- #5 after #6"))
		acceptance.Sees(t, out, code, 1,
			"item 1 has an emoji; GitHub shows each issue's state, so drop it",
			"item 2 (#2) notes a status; GitHub shows each issue's state, so drop it",
			"item 3 (#5 after #6) is not a plain #N or org/repo#N reference; put any short note after a comma")
	})
	t.Run("0001.2.2 a reference listed twice fails", func(t *testing.T) {
		t.Parallel()
		out, code := epicStructure(t, nil, pitch+background+list("- #3\n- #3"))
		acceptance.Sees(t, out, code, 1, "#3 is listed more than once; list each issue once")
	})
	t.Run("0001.2.3 an item listed before the item its note names fails", func(t *testing.T) {
		t.Parallel()
		out, code := epicStructure(t, nil, pitch+background+list("- org/repo#3\n- #1, after #2\n- #2"))
		acceptance.Sees(t, out, code, 1, "item 2 (#1) is listed before #2, which it runs after; list them in run order")
	})
	t.Run("0001.2.4 a table fails", func(t *testing.T) {
		t.Parallel()
		out, code := epicStructure(t, nil, pitch+background+"| a | b |\n| - | - |\n| 1 | 2 |\n\n"+list("- #1"))
		acceptance.Sees(t, out, code, 1, "tables are not allowed; list each issue as a plain reference")
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "epic-structure")
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	acceptance.PassedWithNoFiles(t, out, "epic-structure")
}

// pitch and background open every epic body these tests judge.
const (
	pitch      = "> Every run finishes in a fraction of the time.\n\n"
	background = "## Background\n\nThis epic is a test fixture.\n\n"
)

// list is the What needs to happen section holding items.
func list(items string) string {
	return "## What needs to happen\n\n" + items + "\n"
}

// epicStructure runs `fitness-install -- epic-structure --body-file` on
// happyRepo with files and the body written over it.
func epicStructure(t *testing.T, files map[string]string, body string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", acceptance.With(files, map[string]string{"epic.md": body}))
	return acceptance.Fitness(t, repo, nil, "epic-structure", "--body-file", "epic.md")
}
