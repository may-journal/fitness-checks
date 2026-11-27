// Package requirements proves the semantic-commit check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

// types is the type list the check prints in every failure.
const types = "(types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert)"

// inherited is an upstream subject with no scope, as a fork's main can hold.
const inherited = "test: migrate from jest to vite (#75)"

// noMessage is what an empty message or a missing repo fails with.
const noMessage = "No commit message to validate; use type(scope): description"

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a listed type with a scope or a merge passes on its first line", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", nil)
		messages := []string{"Merge branch 'feature' into main"}
		for _, kind := range []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert"} {
			messages = append(messages, kind+"(api): change the endpoint")
		}
		req.Then(t, "0001.1.1.1.1", "07b8e753", func(_ req.A) {
			for _, message := range messages {
				messageSees(t, repo, message, 0, "All 1 checks passed")
			}
		})
		req.Then(t, "0001.1.1.1.2", "ba400078", func(_ req.A) {
			messageSees(t, repo, "chore(deps): bump tools\n\nNot a semantic line at all", 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.2 a subject with no scope fails with the format then the types", func(t *testing.T) {
		t.Parallel()
		out, code := semanticCommit(t, "--message", "feat: add endpoint")
		req.Then(t, "0001.1.1.2.1", "2221fc49", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, `Commit message: "feat: add endpoint"`)
		})
		req.Then(t, "0001.1.1.2.2", "398d6a4b", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, `Commit message: "feat: add endpoint" — use type(scope): description`)
		})
		req.Then(t, "0001.1.1.2.3", "8f78d607", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, types)
		})
	})
	t.Run("0001.1.3 an unlisted type or a breaking mark fails naming the subject with the format to use", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.3.1", "74436bea", func(_ req.A) {
			for _, subject := range []string{`Revert "feat(api): add endpoint"`, "feat(api)!: drop the old endpoint"} {
				messageSees(t, acceptance.Example(t, "happyRepo", nil), subject, 1, `Commit message: "`+subject+`" — use type(scope): description`)
			}
		})
	})
	t.Run("0001.1.4 an empty message or no repo fails with a no-message error", func(t *testing.T) {
		t.Parallel()
		// happyRepo's HEAD subject is valid, so an empty message must not fall back to it.
		req.Then(t, "0001.1.1.4.1", "04d08f23", func(_ req.A) {
			out, code := semanticCommit(t, "--message=")
			acceptance.Sees(t, out, code, 1, noMessage)
			noRepo := acceptance.Example(t, "happyRepo", nil)
			acceptance.MustDo(t, os.RemoveAll(filepath.Join(noRepo, ".git")))
			out, code = acceptance.Fitness(t, noRepo, nil, "semantic-commit")
			acceptance.Sees(t, out, code, 1, noMessage)
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a repo with no origin refs judges the HEAD subject", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.1.1", "4867e93e", func(_ req.A) {
			repo := acceptance.Example(t, "happyRepo", nil)
			acceptance.Commit(t, repo, nil, "Add notes\n\nBody text")
			repoSees(t, repo, 1, `Commit message: "Add notes" — use type(scope): description `+types)
		})
	})
	t.Run("0001.2.2 only commits origin lacks are judged whatever the default branch is called", func(t *testing.T) {
		t.Parallel()
		for _, base := range []string{"main", "master", "develop"} {
			repo := forkedRepo(t, base)
			acceptance.Commit(t, repo, nil, "Add notes")
			acceptance.Commit(t, repo, nil, "docs(readme): note more")
			out, code := acceptance.Fitness(t, repo, nil)
			row := acceptance.Row(out, "semantic-commit")
			acceptance.Git(t, repo, "checkout", "-q", "-B", base, "origin/"+base)
			acceptance.Commit(t, repo, nil, "Add base notes")
			out2, code2 := acceptance.Fitness(t, repo, nil, "semantic-commit")
			req.Then(t, "0001.2.1.2.1", "19b7f3c8", func(a req.A) {
				quotesOnly(t, out, code, "Add notes")
				a.True(strings.Contains(row, "failed"), "semantic-commit did not fail the full run on %s:\n%s", base, out)
				quotesOnly(t, out2, code2, "Add base notes")
			})
			req.Then(t, "0001.2.1.2.2", "577cc591", func(a req.A) {
				a.True(!strings.Contains(out, inherited), "semantic-commit judged inherited history:\n%s", out)
				a.True(!strings.Contains(out2, inherited), "semantic-commit judged inherited history:\n%s", out2)
			})
		}
	})
	t.Run("0001.2.3 a branch with no commits of its own passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.3.1", "07b8e753", func(_ req.A) {
			for _, base := range []string{"main", "master", "develop"} {
				repoSees(t, forkedRepo(t, base), 0, "All 1 checks passed")
			}
		})
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 a valid subject commits without judging the parent", func(t *testing.T) {
		t.Parallel()
		repo, env := hookedFork(t)
		out, code := acceptance.User(t, repo, env, "commit", "--allow-empty", "-m", "docs(readme): touch nothing")
		req.Then(t, "0001.3.1.1.1", "97c828eb", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "docs(readme): touch nothing")
		})
		req.Then(t, "0001.3.1.1.2", "7a84ba46", func(a req.A) {
			a.True(!strings.Contains(out, inherited), "a hook judged the parent commit:\n%s", out)
		})
	})
	t.Run("0001.3.2 an earlier bad subject fails the next commit", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.3.1.2.1", "f88bf2bd", func(_ req.A) {
			repo, env := hookedFork(t)
			out, code := acceptance.User(t, repo, env, "commit", "--allow-empty", "-m", "Add notes")
			acceptance.Sees(t, out, code, 0, "Add notes")
			out, code = acceptance.User(t, repo, env, "commit", "--allow-empty", "-m", "docs(readme): note more")
			quotesOnly(t, out, code, "Add notes")
		})
	})
}

// semanticCommit runs semantic-commit with args on happyRepo, whose HEAD
// subject is "docs(readme): add readme".
func semanticCommit(t *testing.T, args ...string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, append([]string{"semantic-commit"}, args...)...)
}

// forkedRepo is happyRepo on a fresh branch over inherited history that
// origin/<base> holds. The inherited commit has a bad subject, committed
// without hooks as a fetch brings it, and origin is an empty bare repo.
func forkedRepo(t *testing.T, base string) string {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", nil)
	out, err := exec.Command("git", "-C", repo, "-c", "user.name=t", "-c", "user.email=t@t", "commit-tree", "HEAD^{tree}", "-p", "HEAD", "-m", inherited).Output()
	acceptance.MustDo(t, err)
	acceptance.Git(t, repo, "reset", "-q", "--hard", strings.TrimSpace(string(out)))
	acceptance.Remote(t, repo)
	acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/"+base, "HEAD")
	acceptance.Git(t, repo, "checkout", "-q", "-b", "work")
	return repo
}

// hookedFork is a forked repo with only the installer's pre-commit hook.
func hookedFork(t *testing.T) (string, []string) {
	t.Helper()
	repo := forkedRepo(t, "main")
	cmd := exec.Command(acceptance.Installer(t), "--install-hook")
	cmd.Dir, cmd.Env = repo, acceptance.UserEnv()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("--install-hook: %v\n%s", err, out)
	}
	return repo, nil
}

// quotesOnly asserts semantic-commit failed on subject and never quoted the
// inherited one.
func quotesOnly(t *testing.T, out string, code int, subject string) {
	t.Helper()
	acceptance.Sees(t, out, code, 1, `Commit message: "`+subject+`" — use type(scope): description `+types)
	if strings.Contains(out, inherited) {
		t.Errorf("semantic-commit judged inherited history:\n%s", out)
	}
}

// messageSees runs semantic-commit on repo with message and expects code and want.
func messageSees(t *testing.T, repo, message string, wantCode int, want string) {
	t.Helper()
	out, code := acceptance.Fitness(t, repo, nil, "semantic-commit", "--message", message)
	acceptance.Sees(t, out, code, wantCode, want)
}

// repoSees runs semantic-commit on repo's own commits and expects code and want.
func repoSees(t *testing.T, repo string, wantCode int, want string) {
	t.Helper()
	out, code := acceptance.Fitness(t, repo, nil, "semantic-commit")
	acceptance.Sees(t, out, code, wantCode, want)
}
