package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func Test0004_1(t *testing.T) {
	t.Parallel()
	const source, swiftTest = "Sources/App/App.swift", "Tests/AppTests/AppTests.swift"
	swift := acceptance.With(acceptance.GoRepo(), map[string]string{source: "struct App {}\n", swiftTest: "import Testing\n"})
	t.Run("0004.1.1 a Swift source or a shell script changed alone fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := staged(t, swift, map[string]string{source: "struct App {}\nstruct More {}\n"})
		acceptance.Sees(t, out, code, 1, source+" changed without a test and a requirement doc; change both with it")
		out, code = staged(t, given(), map[string]string{"greet.sh": script + "echo more\n"})
		acceptance.Sees(t, out, code, 1, "greet.sh changed without a test and a requirement doc; change both with it")
	})
	t.Run("0004.1.2 a Swift source beside a Swift test plus a doc passes", func(t *testing.T) {
		t.Parallel()
		out, code := staged(t, swift, map[string]string{
			source:    "struct App {}\nstruct More {}\n",
			swiftTest: "import Testing\n\n@Test func more() {}\n",
			pairedDoc: strings.Replace(acceptance.GoRepo()[pairedDoc], "by hand.", "by hand at all.", 1),
		})
		acceptance.Sees(t, out, code, 0, passed)
	})
	t.Run("0004.1.3 a directive or a shebang or a cgo preamble changed among comment edits counts as code", func(t *testing.T) {
		t.Parallel()
		codeWithComments(t)
	})
}

func Test0004_2(t *testing.T) {
	t.Parallel()
	unpaired := map[string]string{"adder.go": adderPlus()}
	t.Run("0004.2.1 a pre-commit hook alone judges the staged files", func(t *testing.T) {
		t.Parallel()
		out, code := hookRun(t, false, "hook", "pre-commit")
		acceptance.Sees(t, out, code, 1, "[requirements-changed] Please fix these items:", "adder.go changed without a test and a requirement doc")
	})
	t.Run("0004.2.2 a commit-msg hook beside it judges them instead", func(t *testing.T) {
		t.Parallel()
		out, _ := hookRun(t, true, "hook", "pre-commit")
		if strings.Contains(out, "[requirements-changed]") {
			t.Errorf("pre-commit judged the staged files beside a commit-msg hook:\n%s", out)
		}
		message := filepath.Join(t.TempDir(), "message")
		acceptance.MustDo(t, os.WriteFile(message, []byte("feat(adder): double numbers\n"), 0o644))
		out, code := hookRun(t, true, "hook", "commit-msg", message)
		acceptance.Sees(t, out, code, 1, "[requirements-changed] Please fix these items:", "adder.go changed without a test and a requirement doc")
	})
	t.Run("0004.2.3 an amend stages product code whose parent changed the test plus doc and passes", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", acceptance.GoRepo())
		acceptance.Commit(t, repo, map[string]string{
			"adder_test.go": edited(acceptance.GoRepo()["adder_test.go"]),
			pairedDoc:       strings.Replace(acceptance.GoRepo()[pairedDoc], "by hand.", "by hand at all.", 1),
		}, "test(adder): check sums twice")
		acceptance.Write(t, repo, unpaired)
		acceptance.Git(t, repo, "add", "-A")
		out, code := acceptance.Fitness(t, repo, []string{"FITNESS_COMMIT_AMEND=1"}, "requirements-changed")
		acceptance.Sees(t, out, code, 0, passed)
		out, code = acceptance.Fitness(t, repo, nil, "requirements-changed")
		acceptance.Sees(t, out, code, 1, "adder.go changed without a test and a requirement doc")
	})
	t.Run("0004.2.4 a pull request splits product code from its test plus doc across commits and passes", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", acceptance.GoRepo())
		acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
		acceptance.Commit(t, repo, unpaired, "feat(adder): double numbers")
		acceptance.Commit(t, repo, map[string]string{
			"adder_test.go": edited(acceptance.GoRepo()["adder_test.go"]),
			pairedDoc:       strings.Replace(acceptance.GoRepo()[pairedDoc], "by hand.", "by hand at all.", 1),
		}, "test(adder): check doubling")
		out, code := acceptance.Fitness(t, repo, acceptance.PullRequestEnv, "requirements-changed")
		acceptance.Sees(t, out, code, 0, passed)
	})
}

// hookRun runs the fitness hook args on a fresh copy of GoRepo whose adder.go
// is staged with no test or doc. A commit-msg hook sits beside the pre-commit
// run when messageHook is set.
func hookRun(t *testing.T, messageHook bool, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", acceptance.GoRepo())
	if messageHook {
		acceptance.MustDo(t, os.WriteFile(filepath.Join(repo, ".git", "hooks", "commit-msg"), []byte("#!/bin/sh\n"), 0o755))
	}
	acceptance.Write(t, repo, map[string]string{"adder.go": adderPlus()})
	acceptance.Git(t, repo, "add", "-A")
	return acceptance.Fitness(t, repo, nil, args...)
}
