// Package requirements proves the swiftlint check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo, with SwiftLint answering saved responses.
package requirements

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// TestMain lets this test binary stand in for SwiftLint. Run as swiftlint, it
// answers from saved responses instead of running tests.
func TestMain(m *testing.M) {
	acceptance.ServeStandIns([]string{"swiftlint"})
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a violation fails naming file, line, column and rule", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.1.1", "b29b6c4a", func(_ req.A) {
			repo := swiftlintRepo(t, map[string]string{"Short.swift": "let ab = 1\n", "NoNewline.swift": "let value = 1"})
			short := filepath.Join(repo, "Short.swift")
			rp, out, code := swiftlintRun(t, repo, swiftlintFindings(swiftlintShort, short))
			swiftlintAsked(t, rp)
			swiftlintSees(t, out, code, 1, short+
				":1:5 - Variable name 'ab' should be between 3 and 40 characters long (identifier_name)")
			// A whole-file finding has no column.
			whole := filepath.Join(repo, "NoNewline.swift")
			rp, out, code = swiftlintRun(t, repo, swiftlintFindings(swiftlintNoNewline, whole))
			swiftlintAsked(t, rp)
			swiftlintSees(t, out, code, 1, whole+":1 - Files should have a single trailing newline (trailing_newline)")
		})
	})
	t.Run("0001.1.2 clean Swift files pass, counting the files", func(t *testing.T) {
		t.Parallel()
		repo := swiftlintRepo(t, map[string]string{"Clean.swift": "let value = 1\n"})
		rp, out, code := swiftlintRun(t, repo, swiftlintClean)
		req.Then(t, "0001.1.2.2.1", "07b8e753", func(_ req.A) {
			swiftlintAsked(t, rp)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.2.2", "f60f2f46", func(_ req.A) {
			acceptance.PassedWith(t, out, "swiftlint", "1")
		})
	})
	t.Run("0001.1.3 the repo's own config judges", func(t *testing.T) {
		t.Parallel()
		repo := swiftlintRepo(t, map[string]string{
			".swiftlint.yml": "disabled_rules:\n  - identifier_name\n",
			"Short.swift":    "let ab = 1\n",
		})
		rp, out, code := swiftlintRun(t, repo, swiftlintClean)
		req.Then(t, "0001.1.2.3.1", "ee71ae74", func(_ req.A) {
			swiftlintAsked(t, rp)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		// A config excluding every file leaves SwiftLint nothing to lint.
		repo = swiftlintRepo(t, map[string]string{
			".swiftlint.yml": "excluded:\n  - Short.swift\n",
			"Short.swift":    "let ab = 1\n",
		})
		rp, out, code = swiftlintRun(t, repo, swiftlintNothingLeft)
		req.Then(t, "0001.1.2.3.2", "fd41d61f", func(_ req.A) {
			swiftlintAsked(t, rp)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.4 a pull request lints only its changed files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.4.1", "72d79f37", func(_ req.A) {
			for config, want := range map[string]string{
				"":                     "lint --strict --reporter json --quiet --force-exclude Other.swift",
				"included:\n  - Lib\n": swiftlintArgs,
			} {
				files := map[string]string{"Clean.swift": "let value = 1\n", "Other.swift": "let other = 1\n"}
				if config != "" {
					files[".swiftlint.yml"] = config
				}
				repo := swiftlintRepo(t, files)
				acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
				acceptance.Commit(t, repo, map[string]string{"Other.swift": "let other = 2\n"}, "fix(app): change other")
				rp := acceptance.StandIn(t, map[string][]acceptance.Response{"swiftlint": {{Match: []string{want}, Stdout: "[\n\n]\n"}}})
				out, code := acceptance.Fitness(t, repo, append(rp.Env, acceptance.PullRequestEnv...), "swiftlint")
				acceptance.Sees(t, out, code, 0, "All 1 checks passed")
				acceptance.PassedWith(t, out, "swiftlint", map[bool]string{true: "2", false: "1"}[want == swiftlintArgs])
				if calls := rp.Calls("swiftlint"); len(calls) != 1 || calls[0] != want {
					t.Errorf("with config %q, swiftlint calls = %q, want one %q", config, calls, want)
				}
			}
		})
	})
	t.Run("0001.1.5 a repo with no Swift files passes with no files and never runs SwiftLint", func(t *testing.T) {
		t.Parallel()
		rp, out, code := swiftlintRun(t, swiftlintRepo(t, nil))
		req.Then(t, "0001.1.2.5.1", "857dd667", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "swiftlint")
		})
		req.Then(t, "0001.1.2.5.2", "e667d0de", func(a req.A) {
			a.Equal(0, len(rp.Calls("swiftlint")))
		})
	})
	t.Run("0001.1.6 a missing SwiftLint shows the install hint", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "150e1d85", func(_ req.A) {
			repo := swiftlintRepo(t, map[string]string{"Clean.swift": "let value = 1\n"})
			out, code := acceptance.Fitness(t, repo, []string{acceptance.PathWithout(t, "swiftlint")}, "swiftlint")
			acceptance.ShowsInfo(t, out, code, "swiftlint", "SwiftLint not installed: brew install swiftlint")
		})
	})
	t.Run("0001.1.7 a run that fails without a finding shows why", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.7.1", "4ca0535d", func(_ req.A) {
			repo := swiftlintRepo(t, map[string]string{"Clean.swift": "let value = 1\n"})
			rp, out, code := swiftlintRun(t, repo, acceptance.Response{Match: []string{swiftlintArgs}, Stderr: "Fatal error: SourceKit crashed\n", Exit: 1})
			acceptance.ShowsInfo(t, out, code, "swiftlint", "swiftlint lint --strict exited 1 without a finding: Fatal error: SourceKit crashed")
			swiftlintAsked(t, rp)
			out, code = acceptance.Fitness(t, repo, []string{swiftlintScript(t, "kill -9 $$\n")}, "swiftlint")
			acceptance.ShowsInfo(t, out, code, "swiftlint", "swiftlint was stopped by a signal: signal: killed")
			out, code = acceptance.Fitness(t, repo, []string{swiftlintScript(t, "")}, "swiftlint")
			acceptance.ShowsInfo(t, out, code, "swiftlint", "swiftlint could not start: fork/exec")
		})
	})
}

// swiftlintArgs is the strict JSON lint fitness asks SwiftLint for.
const swiftlintArgs = "lint --strict --reporter json --quiet ."

// swiftlintShort is what `swiftlint lint --strict --reporter json` prints for
// `let ab = 1` in file: a warning made an error by --strict, with a column.
const swiftlintShort = `[
  {
    "character" : 5,
    "file" : "%s",
    "line" : 1,
    "reason" : "Variable name 'ab' should be between 3 and 40 characters long",
    "rule_id" : "identifier_name",
    "severity" : "Error",
    "type" : "Identifier Name"
  }
]
`

// swiftlintNoNewline is SwiftLint's JSON for a file missing its trailing
// newline: a whole-file finding, so the column is null.
const swiftlintNoNewline = `[
  {
    "character" : null,
    "file" : "%s",
    "line" : 1,
    "reason" : "Files should have a single trailing newline",
    "rule_id" : "trailing_newline",
    "severity" : "Error",
    "type" : "Trailing Newline"
  }
]
`

// swiftlintClean is SwiftLint's JSON for a run with no findings.
var swiftlintClean = acceptance.Response{Match: []string{swiftlintArgs}, Stdout: "[\n\n]\n"}

// swiftlintNothingLeft is SwiftLint's answer when its config excludes every
// Swift file: an error exit naming no lintable files.
var swiftlintNothingLeft = acceptance.Response{
	Match:  []string{swiftlintArgs},
	Stderr: "No lintable files found at paths: '.'\n",
	Exit:   1,
}

// swiftlintFindings is SwiftLint's strict answer for one finding in path: the
// JSON array and exit 2, since --strict makes warnings errors.
func swiftlintFindings(body, path string) acceptance.Response {
	return acceptance.Response{Match: []string{swiftlintArgs}, Stdout: fmt.Sprintf(body, path), Exit: 2}
}

// swiftlintRepo makes happyRepo with files written over it.
func swiftlintRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	return acceptance.Example(t, "happyRepo", files)
}

// swiftlintRun runs `fitness-install -- swiftlint` on repo with SwiftLint
// answering saved responses, and returns the stand-ins with the result.
func swiftlintRun(t *testing.T, repo string, saved ...acceptance.Response) (acceptance.Replays, string, int) {
	t.Helper()
	rp := acceptance.StandIn(t, map[string][]acceptance.Response{"swiftlint": saved})
	out, code := acceptance.Fitness(t, repo, rp.Env, "swiftlint")
	return rp, out, code
}

// swiftlintScript returns a PATH whose swiftlint is the shell script body,
// ahead of the real tools. An empty body names an interpreter that does not
// exist, so SwiftLint cannot start.
func swiftlintScript(t *testing.T, body string) string {
	t.Helper()
	shebang := "#!/bin/sh\n"
	if body == "" {
		shebang = "#!/nonexistent/interpreter\n"
	}
	dir := t.TempDir()
	acceptance.Write(t, dir, map[string]string{"swiftlint": shebang + body})
	return "PATH=" + dir + string(os.PathListSeparator) + os.Getenv("PATH")
}

// swiftlintAsked asserts fitness ran SwiftLint's strict JSON lint once.
func swiftlintAsked(t *testing.T, rp acceptance.Replays) {
	t.Helper()
	req.Asserted(t)
	if calls := rp.Calls("swiftlint"); len(calls) != 1 || calls[0] != swiftlintArgs {
		t.Errorf("swiftlint calls = %q, want one %q", calls, swiftlintArgs)
	}
}

// swiftlintSees asserts out shows want and exit code code. The table wraps
// long absolute paths mid-word, so spaces are ignored when comparing.
func swiftlintSees(t *testing.T, out string, gotCode, code int, want string) {
	t.Helper()
	acceptance.Sees(t, out, gotCode, code)
	if !strings.Contains(swiftlintSquash(out), swiftlintSquash(want)) {
		t.Errorf("output is missing %q:\n%s", want, out)
	}
}

// swiftlintSquash drops table borders and all whitespace from text.
func swiftlintSquash(text string) string {
	return strings.ReplaceAll(acceptance.Read(text), " ", "")
}
