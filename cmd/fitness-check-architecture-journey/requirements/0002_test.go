package requirements

import (
	"fmt"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// ajSection is a new section a pull request appends to the journey view.
const ajSection = "\n## Rerun the suites after a check release\n\nIt reruns them.\n"

// ajAddedSection is the error for the section ajSection adds to the view of files.
func ajAddedSection(files map[string]string) string {
	line := strings.Count(files[ajView], "\n") + 2
	return fmt.Sprintf("architecture/05-user-journey.md:%d: the journey doc gains the new section `## Rerun the suites after a check release`", line)
}

func Test0002_1(t *testing.T) {
	t.Parallel()
	t.Run("0002.1.1 an added `##` section fails at its heading", func(t *testing.T) {
		t.Parallel()
		t.Run("a pull request adding a section fails at its heading", func(t *testing.T) {
			t.Parallel()
			files := ajFixture()
			out, code := ajPullRequest(t, files, map[string]string{ajView: files[ajView] + ajSection})
			sees(t, out, code, 1, ajAddedSection(files))
		})
		t.Run("a section with its own diagram fails telling to extend a flow", func(t *testing.T) {
			t.Parallel()
			files := ajFixture()
			own := ajSection + "\n```mermaid\nflowchart TB\nwho((1 Who)):::persona\nrerun[(2 Rerun<br/>Req 0001)]:::system\nwho --> rerun\n```\n"
			out, code := ajPullRequest(t, files, map[string]string{ajView: files[ajView] + own})
			sees(t, out, code, 1, ajAddedSection(files), "extend an existing flow with the requirement's numbered nodes and callout rows, and cite the requirement there")
		})
		t.Run("a commit adding a section fails at its heading", func(t *testing.T) {
			t.Parallel()
			files := ajFixture()
			repo := example(t, "happyRepo", files)
			acceptance.Write(t, repo, map[string]string{ajView: files[ajView] + ajSection})
			acceptance.Git(t, repo, "add", "-A")
			out, code := fitness(t, repo, nil, ajCheck)
			sees(t, out, code, 1, ajAddedSection(files))
		})
	})
	t.Run("0002.1.2 changes that add no section pass", func(t *testing.T) {
		t.Parallel()
		t.Run("a node added to an existing diagram passes", func(t *testing.T) {
			t.Parallel()
			files := ajFixture()
			grown := ajEdit(ajFixture(), "user --> open", "rerun[(4 Rerun<br/>Req 0001)]:::system\nuser --> open\ndone --> rerun")
			out, code := ajPullRequest(t, files, map[string]string{ajView: grown[ajView]})
			sees(t, out, code, 0)
		})
		t.Run("a subsection passes", func(t *testing.T) {
			t.Parallel()
			files := ajFixture()
			out, code := ajPullRequest(t, files, map[string]string{ajView: files[ajView] + "\n### Rerun\n\nIt reruns them.\n"})
			sees(t, out, code, 0)
		})
		t.Run("a section heading inside a code fence passes", func(t *testing.T) {
			t.Parallel()
			files := ajFixture()
			out, code := ajPullRequest(t, files, map[string]string{ajView: files[ajView] + "\n```text\n## Rerun the suites after a check release\n```\n"})
			sees(t, out, code, 0)
		})
		t.Run("a journey doc new in the pull request keeps its sections", func(t *testing.T) {
			t.Parallel()
			files := ajFixture()
			journey := files[ajView]
			files[ajView] = ""
			repo := example(t, "happyRepo", files)
			acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
			acceptance.Commit(t, repo, map[string]string{ajView: journey + ajSection}, "feat(app): add the journey")
			out, code := fitness(t, repo, acceptance.PullRequestEnv, ajCheck)
			sees(t, out, code, 0)
		})
	})
}
