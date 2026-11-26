// Package requirements proves the issue-link-once check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo with a local origin.
package requirements

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain lets this test binary stand in for gh, which answers the pull
// request the branch has open.
func TestMain(m *testing.M) {
	acceptance.Main(m, "gh")
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a repeated link fails naming the earlier commit", func(t *testing.T) {
		t.Parallel()
		repo := onGitHub(t, "topic")
		commitAs(t, repo, "", "feat(app): start\n\nPlan #12")
		out, code := check(t, repo, "fix(app): more\n\nSee https://github.com/"+ghRepo+"/issues/12")
		acceptance.Sees(t, out, code, 1, ghRepo+`#12 is already linked by commit`, `("feat(app): start")`)
		out, code = check(t, repo, "fix(app): more\n\nPlan #12")
		acceptance.Sees(t, out, code, 1, ghRepo+`#12 is already linked by commit`, `("feat(app): start")`)
		repo = onGitHub(t, "topic")
		commitAs(t, repo, "", "feat(app): begin\n\nPlan may-journal/.github#144")
		out, code = check(t, repo, "fix(app): more\n\nPlan https://github.com/may-journal/.github/issues/144")
		acceptance.Sees(t, out, code, 1, `may-journal/.github#144 is already linked by commit`, `("feat(app): begin")`)
	})
	t.Run("0001.1.2 a link to another repo passes", func(t *testing.T) {
		t.Parallel()
		repo := onGitHub(t, "topic")
		commitAs(t, repo, "", "feat(app): start\n\nPlan #12")
		out, code := check(t, repo, "fix(app): more\n\nSee may-journal/fitness-runner#12")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.3 comment and scissors lines are ignored", func(t *testing.T) {
		t.Parallel()
		repo := branch(t, "topic")
		commitAs(t, repo, "", "feat(app): start\n\nPlan #12")
		message := "fix(app): more\n# Plan #12\n# ------------------------ >8 ------------------------\nPlan #12\n"
		out, code := check(t, repo, message)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 neither a message nor a pull request passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "issue-link-once")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "issue-link-once")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	repo := branch(t, "topic")
	acceptance.Git(t, repo, "remote", "set-url", "origin", "https://github.com/acme/app.git")
	rp := acceptance.StandIn(t, map[string][]acceptance.Response{"gh": {
		{Match: []string{"pr view", "--json number,state,body"}, Stdout: `{"body":"Adds the export screen.\n\nCloses #12","number":7,"state":"OPEN"}` + "\n"},
	}})
	out, code := acceptance.Fitness(t, repo, rp.Env, "issue-link-once", "--message", "fix(app): more\n\nFixes #12")
	acceptance.Sees(t, out, code, 1, "acme/app#12 is already linked by open PR #7")
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 a second commit repeating a link fails", func(t *testing.T) {
		t.Parallel()
		repo := branch(t, "topic")
		commitAs(t, repo, "2026-09-30T10:00:00Z", "feat(app): start\n\nPlan #12")
		commitAs(t, repo, "2026-09-30T11:00:00Z", "fix(app): more\n\nPlan #12")
		out, code := pullRequest(t, repo)
		acceptance.Sees(t, out, code, 1, "commit "+head(t, repo)[:7]+` ("fix(app): more"): #12 is already linked by commit`, "--force-with-lease")
	})
	t.Run("0001.3.2 a later commit repeating the pull request link fails", func(t *testing.T) {
		t.Parallel()
		repo := branch(t, "topic")
		commitAs(t, repo, "2026-09-30T13:00:00Z", "fix(app): more\n\nFixes #12")
		out, code := pullRequest(t, repo)
		acceptance.Sees(t, out, code, 1, "#12 is already linked by open PR #9", "--force-with-lease")
	})
	t.Run("0001.3.3 one early link passes with main merged in", func(t *testing.T) {
		t.Parallel()
		repo := branch(t, "topic")
		commitAs(t, repo, "2026-09-30T10:00:00Z", "feat(app): start\n\nPlan #12")
		acceptance.Git(t, repo, "switch", "-q", "main")
		commitAs(t, repo, "2026-09-30T13:00:00Z", "feat(app): more on main\n\nPlan #12")
		acceptance.Git(t, repo, "push", "-q", "origin", "main")
		acceptance.Git(t, repo, "switch", "-q", "topic")
		acceptance.User(t, repo, nil, "merge", "-q", "--no-edit", "main")
		out, code := pullRequest(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 2 files scanned")
	})
}

// ghRepo is the repo on GitHub that origin points at, so a bare `#NN`
// belongs to it.
const ghRepo = "may-journal/example"

// origin is ghRepo's URL.
const origin = "https://github.com/" + ghRepo + ".git"

// branch is happyRepo on a branch, off a main pushed to a local origin.
func branch(t *testing.T, name string) string {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", nil)
	acceptance.Remote(t, repo)
	acceptance.Git(t, repo, "push", "-q", "origin", "HEAD:main")
	acceptance.Git(t, repo, "switch", "-q", "-c", name)
	return repo
}

// onGitHub is branch with origin pointed at ghRepo on GitHub.
func onGitHub(t *testing.T, name string) string {
	t.Helper()
	repo := branch(t, name)
	acceptance.Git(t, repo, "remote", "set-url", "origin", origin)
	return repo
}

// commitAs commits message in repo, authored at date when it is set.
func commitAs(t *testing.T, repo, date, message string) {
	t.Helper()
	if out, code := acceptance.User(t, repo, []string{"GIT_AUTHOR_DATE=" + date}, "commit", "-q", "--allow-empty", "-m", message); code != 0 {
		t.Fatalf("commit %q: %s", message, out)
	}
}

// check runs issue-link-once on the proposed commit message.
func check(t *testing.T, repo, message string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, repo, nil, "issue-link-once", "--message", message)
}

// head returns repo's HEAD commit.
func head(t *testing.T, repo string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", repo, "rev-parse", "HEAD").Output()
	acceptance.MustDo(t, err)
	return strings.TrimSpace(string(out))
}

// pullRequest runs issue-link-once as CI does on pull request #9, opened at
// noon with a body closing #12, whose head is repo's HEAD.
func pullRequest(t *testing.T, repo string) (string, int) {
	t.Helper()
	payload := map[string]any{"pull_request": map[string]any{
		"number": 9, "body": "Closes #12", "created_at": "2026-09-30T12:00:00Z",
		"head": map[string]any{"sha": head(t, repo)},
	}}
	data, err := json.Marshal(payload)
	acceptance.MustDo(t, err)
	path := filepath.Join(t.TempDir(), "event.json")
	acceptance.MustDo(t, os.WriteFile(path, data, 0o644))
	env := []string{"GITHUB_BASE_REF=main", "GITHUB_EVENT_PATH=" + path}
	return acceptance.Fitness(t, repo, env, "issue-link-once")
}
