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
	t.Run("0001.1.2.1.1 it fails naming each as `file:line[:col] - reason (rule_id)`", func(t *testing.T) {
		t.Parallel()
		repo := swiftlintRepo(t, map[string]string{"Short.swift": "let ab = 1\n", "NoNewline.swift": "let value = 1"})
		short := filepath.Join(repo, "Short.swift")
		rp, out, code := swiftlintRun(t, repo, swiftlintFindings(swiftlintShort, short))
		assertSwiftlintAsked(t, rp)
		assertSwiftlintSees(t, out, code, 1, short+
			":1:5 - Variable name 'ab' should be between 3 and 40 characters long (identifier_name)")
		// A whole-file finding has no column.
		whole := filepath.Join(repo, "NoNewline.swift")
		rp, out, code = swiftlintRun(t, repo, swiftlintFindings(swiftlintNoNewline, whole))
		assertSwiftlintAsked(t, rp)
		assertSwiftlintSees(t, out, code, 1, whole+":1 - Files should have a single trailing newline (trailing_newline)")
	})
	t.Run("0001.1.2.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		repo := swiftlintRepo(t, map[string]string{"Clean.swift": "let value = 1\n"})
		rp, out, code := swiftlintRun(t, repo, swiftlintClean)
		assertSwiftlintAsked(t, rp)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.2.2 it counts the files SwiftLint was asked to lint", func(t *testing.T) {
		t.Parallel()
		repo := swiftlintRepo(t, map[string]string{"Clean.swift": "let value = 1\n"})
		_, out, code := swiftlintRun(t, repo, swiftlintClean)
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWith(t, out, "swiftlint", "1")
	})
	t.Run("0001.1.2.3.1 that config judges", func(t *testing.T) {
		t.Parallel()
		repo := swiftlintRepo(t, map[string]string{
			".swiftlint.yml": "disabled_rules:\n  - identifier_name\n",
			"Short.swift":    "let ab = 1\n",
		})
		rp, out, code := swiftlintRun(t, repo, swiftlintClean)
		assertSwiftlintAsked(t, rp)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.3.2 a config excluding every file passes", func(t *testing.T) {
		t.Parallel()
		// A config excluding every file leaves SwiftLint nothing to lint.
		repo := swiftlintRepo(t, map[string]string{
			".swiftlint.yml": "excluded:\n  - Short.swift\n",
			"Short.swift":    "let ab = 1\n",
		})
		rp, out, code := swiftlintRun(t, repo, swiftlintNothingLeft)
		assertSwiftlintAsked(t, rp)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.4.1 SwiftLint lints only those unless the config lists what to include", func(t *testing.T) {
		t.Parallel()
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
	t.Run("0001.1.2.5.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		_, out, code := swiftlintRun(t, swiftlintRepo(t, nil))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "swiftlint")
	})
	t.Run("0001.1.2.5.2 it never runs SwiftLint", func(t *testing.T) {
		t.Parallel()
		rp, out, code := swiftlintRun(t, swiftlintRepo(t, nil))
		acceptance.Sees(t, out, code, 0)
		if calls := rp.Calls("swiftlint"); len(calls) != 0 {
			t.Errorf("swiftlint ran with no Swift files: %q", calls)
		}
	})
	t.Run("0001.1.2.6.1 it reports info with the hint `brew install swiftlint`", func(t *testing.T) {
		t.Parallel()
		repo := swiftlintRepo(t, map[string]string{"Clean.swift": "let value = 1\n"})
		out, code := acceptance.Fitness(t, repo, []string{acceptance.PathWithout(t, "swiftlint")}, "swiftlint")
		acceptance.ShowsInfo(t, out, code, "swiftlint", "SwiftLint not installed: brew install swiftlint")
	})
	t.Run("0001.1.2.7.1 it reports info quoting what went wrong", func(t *testing.T) {
		t.Parallel()
		repo := swiftlintRepo(t, map[string]string{"Clean.swift": "let value = 1\n"})
		rp, out, code := swiftlintRun(t, repo, acceptance.Response{Match: []string{swiftlintArgs}, Stderr: "Fatal error: SourceKit crashed\n", Exit: 1})
		acceptance.ShowsInfo(t, out, code, "swiftlint", "swiftlint lint --strict exited 1 without a finding: Fatal error: SourceKit crashed")
		assertSwiftlintAsked(t, rp)
		out, code = acceptance.Fitness(t, repo, []string{swiftlintScript(t, "kill -9 $$\n")}, "swiftlint")
		acceptance.ShowsInfo(t, out, code, "swiftlint", "swiftlint was stopped by a signal: signal: killed")
		out, code = acceptance.Fitness(t, repo, []string{swiftlintScript(t, "")}, "swiftlint")
		acceptance.ShowsInfo(t, out, code, "swiftlint", "swiftlint could not start: fork/exec")
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

// assertSwiftlintAsked asserts fitness ran SwiftLint's strict JSON lint once.
func assertSwiftlintAsked(t *testing.T, rp acceptance.Replays) {
	t.Helper()
	if calls := rp.Calls("swiftlint"); len(calls) != 1 || calls[0] != swiftlintArgs {
		t.Errorf("swiftlint calls = %q, want one %q", calls, swiftlintArgs)
	}
}

// assertSwiftlintSees asserts out shows want and exit code code. The table wraps
// long absolute paths mid-word, so spaces are ignored when comparing.
func assertSwiftlintSees(t *testing.T, out string, gotCode, code int, want string) {
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
