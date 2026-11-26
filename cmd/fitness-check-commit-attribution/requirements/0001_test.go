// Package requirements proves the commit-attribution check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo with commits made over it.
package requirements

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1.1.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := commitAttribution(t, []string{"feat(api): add endpoint\n\nNo trailers here.\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "commit-attribution")
		folder := acceptance.Example(t, "happyRepo", nil)
		acceptance.MustDo(t, os.RemoveAll(filepath.Join(folder, ".git")))
		out, code = acceptance.Fitness(t, folder, nil, "commit-attribution")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "commit-attribution")
	})
	t.Run("0001.1.1.2.1 it passes counting one file", func(t *testing.T) {
		t.Parallel()
		out, code := commitAttribution(t, []string{trailers})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
	})
	t.Run("0001.1.1.3.1 it fails with one error per missing trailer as `commit message missing \"<key>:\" trailer`", func(t *testing.T) {
		t.Parallel()
		out, code := commitAttribution(t, []string{trailers, "feat(api): add paging\n\nJust a body.\n"})
		acceptance.Sees(t, out, code, 1, missing("AI-Models"), missing("AI-Tools"))
	})
	t.Run("0001.1.1.4.1 it counts as missing", func(t *testing.T) {
		t.Parallel()
		out, code := commitAttribution(t, []string{trailers, "feat(api): add paging\n\nAI-Tools: Claude Code\nAI-Models:   \n"})
		acceptance.Sees(t, out, code, 1, missing("AI-Models"))
		out, code = commitAttribution(t, []string{trailers, "feat(api): add paging\n\nAI-Tools: Claude Code\nWe set AI-Models: Opus 4.8 later.\n"})
		acceptance.Sees(t, out, code, 1, missing("AI-Models"))
	})
	t.Run("0001.1.1.5.1 it passes without trailers", func(t *testing.T) {
		t.Parallel()
		for _, subject := range []string{"Merge branch 'feature' into main", `Revert "feat(api): add endpoint"`} {
			out, code := commitAttribution(t, []string{trailers, subject})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		}
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1.1.1 it fails for the proposed message instead of the last commit", func(t *testing.T) {
		t.Parallel()
		out, code := commitAttribution(t, []string{trailers}, "--message=feat(api): add paging\n\nJust a body.\n")
		acceptance.Sees(t, out, code, 1, missing("AI-Models"), missing("AI-Tools"))
	})
	t.Run("0001.2.1.2.1 it fails asking for both trailers", func(t *testing.T) {
		t.Parallel()
		out, code := commitAttribution(t, []string{trailers}, "--message=")
		acceptance.Sees(t, out, code, 1, `No commit message to validate; add "AI-Tools:" and "AI-Models:" trailers to disclose AI usage`)
	})
}

// trailers is a commit message carrying both AI trailers.
const trailers = "feat(api): add endpoint\n\nAI-Tools: Claude Code\nAI-Models: Opus 4.8\n"

// missing is the error for a missing trailer key.
func missing(key string) string {
	return `commit message missing "` + key + `:" trailer`
}

// commitAttribution commits each message on happyRepo in order, the last at
// HEAD, then runs commit-attribution with args.
func commitAttribution(t *testing.T, messages []string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", nil)
	for _, m := range messages {
		acceptance.Commit(t, repo, nil, m)
	}
	return acceptance.Fitness(t, repo, nil, append([]string{"commit-attribution"}, args...)...)
}
