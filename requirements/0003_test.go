package requirements

import (
	"os/exec"
	"strings"
	"testing"
)

func Test0003_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"hooks accept a valid subject on a fresh branch", func(t *testing.T) {
			t.Parallel()
			repo, env := hookedFork(t)
			out, code := user(t, repo, env, "commit", "--allow-empty", "-m", "docs(readme): touch nothing")
			sees(t, out, code, 0, "docs(readme): touch nothing")
			if strings.Contains(out, inherited) {
				t.Errorf("a hook judged the parent commit:\n%s", out)
			}
		}},
		{"hooks refuse a subject with no scope", func(t *testing.T) {
			t.Parallel()
			repo, env := hookedFork(t)
			out, code := user(t, repo, env, "commit", "--allow-empty", "-m", "Add notes")
			quotesOnly(t, out, code, "Add notes")
		}},
	})
}

// semanticCommitTypes is the type list the check prints in every failure.
const semanticCommitTypes = "(types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert)"

// inherited is an upstream subject with no scope, as a fork's main can hold.
const inherited = "test: migrate from jest to vite (#75)"

// fork commits an inherited bad subject to repo without hooks, as a fetch
// does, records it as origin/main, and branches fresh.
func fork(t *testing.T, repo string) {
	t.Helper()
	out, err := exec.Command("git", "-C", repo, "-c", "user.name=t", "-c", "user.email=t@t", "commit-tree", "HEAD^{tree}", "-p", "HEAD", "-m", inherited).Output()
	mustDo(t, err)
	git(t, repo, "reset", "-q", "--hard", strings.TrimSpace(string(out)))
	remote(t, repo)
	git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	git(t, repo, "checkout", "-q", "-b", "work")
}

// hookedFork is a happyRepo fork with the shared hooks.
func hookedFork(t *testing.T) (string, []string) {
	t.Helper()
	repo, env := hooked(t, nil)
	fork(t, repo)
	return repo, env
}

// quotesOnly asserts semantic-commit failed on subject and never quoted the
// inherited one.
func quotesOnly(t *testing.T, out string, code int, subject string) {
	t.Helper()
	sees(t, out, code, 1, `Commit message: "`+subject+`" — use type(scope): description `+semanticCommitTypes)
	if strings.Contains(out, inherited) {
		t.Errorf("semantic-commit judged inherited history:\n%s", out)
	}
}
