package requirements

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// Paths and edits the pairing tests share.
const (
	pairedTest = "reqs/0001_test.go"
	twoTest    = "reqs/0001_0002_test.go"
	swiftTest  = "Tests/AdderTests/AdderTests.swift"
)

// edited is body with one more function, a code change that keeps Go and
// Swift valid.
func edited(body string) string {
	return body + "\nfunc oneMore() {}\n"
}

// withComment is body with one more comment line, a change to no code.
func withComment(body string) string {
	return body + "\n// One more line.\n"
}

// reworded is a requirement doc with one more word in its Why, changing what
// it says, not only its formatting.
func reworded(doc string) string {
	return strings.Replace(doc, "by hand.", "by hand at all.", 1)
}

// unformatted is the requirement doc before a formatter rewrote its quotes,
// wrapping, blank lines and nesting.
var unformatted = strings.NewReplacer(
	"two numbers without", "two numbers\nwithout",
	`['../../.fitnessrc.json']`, `["../../.fitnessrc.json"]`,
	"\n\n-\n\n", "\n-\n",
	"\n   1. GIVEN", "\n    1. GIVEN",
	"\n      1. WHEN", "\n        1. WHEN",
	"\n         1. THEN", "\n            1. THEN",
)

// hooked is files over happyRepo with the installer's pre-commit hook alone.
func hooked(t *testing.T, files map[string]string) string {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	cmd := exec.Command(acceptance.Installer(t), "--install-hook")
	cmd.Dir, cmd.Env = repo, acceptance.UserEnv()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("--install-hook: %v\n%s", err, out)
	}
	return repo
}

// stagedRun writes files over a fresh copy of given, stages them, and runs
// the requirements check on the staged changes.
func stagedRun(t *testing.T, given, files map[string]string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", given)
	acceptance.Write(t, repo, files)
	acceptance.Git(t, repo, "add", "-A")
	return acceptance.Fitness(t, repo, nil, "requirements")
}

// prRun commits files with message on a pull request branch off a fresh copy
// of given and runs the requirements check on it.
func prRun(t *testing.T, given, files map[string]string, message string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", given)
	acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	acceptance.Write(t, repo, files)
	acceptance.Git(t, repo, "add", "-A")
	acceptance.Git(t, repo, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", message)
	return acceptance.Fitness(t, repo, acceptance.PullRequestEnv, "requirements")
}

// userSteps runs each git command as a user, failing the test on the first
// one that fails.
func userSteps(t *testing.T, repo string, steps ...[]string) {
	t.Helper()
	for _, step := range steps {
		if out, code := acceptance.User(t, repo, nil, step...); code != 0 {
			t.Fatalf("git %v: exit %d\n%s", step, code, out)
		}
	}
}

// replayed replays an unpaired --no-verify test change over the same line
// changed with its doc, returning what the first failing commit printed.
func replayed(t *testing.T, replay []string, finish ...[]string) (out string, code int) {
	t.Helper()
	repo := hooked(t, acceptance.PairedRepo())
	base := acceptance.TestFile("reqs", "Test0001_1")
	acceptance.Write(t, repo, map[string]string{pairedTest: strings.Replace(base, "unnamed", "unnamed one", 1)})
	userSteps(t, repo, []string{"checkout", "-q", "-b", "side"}, []string{"add", "-A"},
		[]string{"commit", "-q", "--no-verify", "-m", "test(reqs): reword the test"}, []string{"checkout", "-q", "-"})
	acceptance.Write(t, repo, map[string]string{pairedTest: strings.Replace(base, "unnamed", "unnamed two", 1), acceptance.PairedDoc: reworded(acceptance.RequirementDoc("0001.1"))})
	userSteps(t, repo, []string{"add", "-A"}, []string{"commit", "-q", "-m", "test(reqs): check sums again"})
	if _, code := acceptance.User(t, repo, nil, replay...); code == 0 {
		t.Fatalf("git %v did not conflict", replay)
	}
	userSteps(t, repo, []string{"checkout", "--theirs", pairedTest}, []string{"add", "-A"})
	for _, args := range finish {
		if out, code = acceptance.User(t, repo, nil, append([]string{"commit"}, args...)...); code != 0 {
			break
		}
	}
	return out, code
}

// mergedBase is a branch that merged its moved base, whose new commit
// changes a requirement doc alone. It returns the repo, with origin/main at
// the base's new tip.
func mergedBase(t *testing.T) string {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", acceptance.PairedRepo())
	acceptance.Git(t, repo, "checkout", "-q", "-b", "base")
	acceptance.Commit(t, repo, map[string]string{acceptance.PairedDoc: reworded(acceptance.RequirementDoc("0001.1"))}, "docs(requirements): reword 0001")
	acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	acceptance.Git(t, repo, "checkout", "-q", "-b", "work", "HEAD~1")
	acceptance.Commit(t, repo, map[string]string{"README.md": "# App\n\nShort and clean.\n\nMore words.\n"}, "docs(readme): add words")
	acceptance.Git(t, repo, "-c", "user.name=t", "-c", "user.email=t@t", "merge", "-q", "--no-edit", "base")
	return repo
}

// twoRepo is pairedRepo whose one test file proves requirements 0001 and 0002.
// Requirement 0002 covers its own area, so the two docs never overlap.
func twoRepo() map[string]string {
	doubles := "1. Doubles\n   1. GIVEN one number\n      1. WHEN I double it\n         1. THEN I get twice that number\n"
	return acceptance.With(acceptance.PairedRepo(), map[string]string{
		"reqs/0001_test.go":                    "",
		twoTest:                                acceptance.CaseFiles("reqs", acceptance.CaseFunc("Test0001_1", "0001.1.1.1.1 I get their sum"), acceptance.CaseFunc("Test0002_1", "0002.1.1.1.1 I get twice that number")),
		"docs/requirements/0002-adds-again.md": strings.Replace(acceptance.DocWith(doubles), "0001 Adds Numbers", "0002 Adds Again", 1),
	})
}

// renamedPairs proves a renamed doc keeping its number and wording needs no test
// change, unlike one with new wording.
func renamedPairs(t *testing.T) {
	t.Helper()
	const renamed = "docs/requirements/0001-sums-numbers.md"
	moved := map[string]string{acceptance.PairedDoc: "", renamed: acceptance.RequirementDoc("0001.1")}
	out, code := stagedRun(t, acceptance.PairedRepo(), moved)
	acceptance.Sees(t, out, code, 0)
	out, code = prRun(t, acceptance.PairedRepo(), moved, "docs(requirements): rename the sum requirement")
	acceptance.Sees(t, out, code, 0)
	out, code = stagedRun(t, acceptance.PairedRepo(), map[string]string{acceptance.PairedDoc: "", renamed: reworded(acceptance.RequirementDoc("0001.1"))})
	acceptance.Sees(t, out, code, 1, renamed+" changed without its requirement test")
}

// commentedPairs proves a test or doc changed only in comments needs no
// partner change. A code change beside the comments still needs one.
func commentedPairs(t *testing.T) {
	t.Helper()
	test := acceptance.TestFile("reqs", "Test0001_1")
	rewritten := strings.Replace(withComment(test), "func Test0001_1", "// Test0001_1 proves the sum.\nfunc Test0001_1", 1)
	noted := map[string]string{acceptance.PairedDoc: acceptance.RequirementDoc("0001.1") + "\n<!-- Sums stay small. -->\n"}
	for _, files := range []map[string]string{{pairedTest: rewritten}, noted} {
		out, code := stagedRun(t, acceptance.PairedRepo(), files)
		acceptance.Sees(t, out, code, 0)
		out, code = prRun(t, acceptance.PairedRepo(), files, "docs(reqs): reword comments")
		acceptance.Sees(t, out, code, 0)
	}
	unpaired := pairedTest + " changed without its requirement doc " + acceptance.PairedDoc + "; change both together"
	out, code := stagedRun(t, acceptance.PairedRepo(), map[string]string{pairedTest: edited(rewritten)})
	acceptance.Sees(t, out, code, 1, unpaired)
	out, code = prRun(t, acceptance.PairedRepo(), map[string]string{pairedTest: edited(rewritten)}, "test(reqs): check sums twice")
	acceptance.Sees(t, out, code, 1, unpaired)
	out, code = stagedRun(t, acceptance.PairedRepo(), acceptance.With(noted, map[string]string{pairedTest: edited(rewritten)}))
	acceptance.Sees(t, out, code, 1, unpaired)
	out, code = stagedRun(t, acceptance.PairedRepo(), map[string]string{pairedTest: rewritten, acceptance.PairedDoc: reworded(acceptance.RequirementDoc("0001.1"))})
	acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+" changed without its requirement test "+pairedTest+"; change both together")
}

func Test0006_1(t *testing.T) {
	t.Parallel()
	// The commit hook marks an amend this way, as it finds git running with --amend.
	amending := []string{"FITNESS_COMMIT_AMEND=1"}
	unpairedDoc := pairedTest + " changed without its requirement doc " + acceptance.PairedDoc + "; change both together"
	t.Run("0006.1.1.1.1 it fails naming both files", func(t *testing.T) {
		t.Parallel()
		repo := hooked(t, acceptance.PairedRepo())
		acceptance.Write(t, repo, map[string]string{pairedTest: edited(acceptance.TestFile("reqs", "Test0001_1"))})
		acceptance.User(t, repo, nil, "add", "-A")
		out, code := acceptance.User(t, repo, nil, "commit", "-m", "test(reqs): check sums twice")
		acceptance.Sees(t, out, code, 1, unpairedDoc)
		out, code = prRun(t, acceptance.PairedRepo(), map[string]string{acceptance.PairedDoc: reworded(acceptance.RequirementDoc("0001.1"))}, "docs(requirements): reword 0001")
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+" changed without its requirement test "+pairedTest+"; change both together")
	})
	t.Run("0006.1.1.2.1 a change to either doc pairs it", func(t *testing.T) {
		t.Parallel()
		doc := "docs/requirements/0002-adds-again.md"
		test := map[string]string{twoTest: edited(twoRepo()[twoTest])}
		out, code := stagedRun(t, twoRepo(), test)
		acceptance.Sees(t, out, code, 1, twoTest+" changed without its requirement doc "+acceptance.PairedDoc+" or "+doc)
		out, code = stagedRun(t, twoRepo(), acceptance.With(test, map[string]string{doc: reworded(twoRepo()[doc])}))
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0006.1.1.3.1 it needs no partner", func(t *testing.T) {
		t.Parallel()
		out, code := stagedRun(t, acceptance.PairedRepo(), map[string]string{"reqs/tools_test.go": edited(acceptance.PairedRepo()["reqs/tools_test.go"])})
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0006.1.1.4.1 only its own changes are judged", func(t *testing.T) {
		t.Parallel()
		keep := []string{"--no-edit"}
		for _, r := range [][][]string{
			{{"merge", "--no-edit", "side"}, keep, {"--amend", "-m", "chore(reqs): merge the reword"}},
			{{"cherry-pick", "side"}, keep},
			{{"revert", "--no-edit", "side"}, {"-m", "revert(reqs): undo the reword"}},
		} {
			out, code := replayed(t, r[0], r[1:]...)
			acceptance.Sees(t, out, code, 0)
		}
		out, code := acceptance.Fitness(t, mergedBase(t), acceptance.PullRequestEnv, "requirements")
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0006.1.1.5.1 it is judged since its parent", func(t *testing.T) {
		t.Parallel()
		repo := hooked(t, acceptance.PairedRepo())
		acceptance.Write(t, repo, map[string]string{pairedTest: edited(acceptance.TestFile("reqs", "Test0001_1"))})
		userSteps(t, repo, []string{"add", "-A"}, []string{"commit", "-q", "--no-verify", "-m", "test(reqs): check sums twice"})
		acceptance.Write(t, repo, map[string]string{acceptance.PairedDoc: reworded(acceptance.RequirementDoc("0001.1"))})
		acceptance.Git(t, repo, "add", "-A")
		out, code := acceptance.Fitness(t, repo, amending, "requirements")
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0006.1.1.6.1 it needs no doc", func(t *testing.T) {
		t.Parallel()
		failing := acceptance.With(swiftRepo(), map[string]string{swiftTest: xcTestFile("test0001_1")})
		fixed := map[string]string{swiftTest: swiftTestingFile("test0001_1")}
		out, code := stagedRun(t, failing, fixed)
		acceptance.Sees(t, out, code, 0)
		out, code = prRun(t, failing, fixed, "test(adder): move to Swift Testing")
		acceptance.Sees(t, out, code, 0)
		repo := hooked(t, failing)
		acceptance.Write(t, repo, fixed)
		userSteps(t, repo, []string{"add", "-A"}, []string{"commit", "-q", "--no-verify", "-m", "test(adder): move to Swift Testing"})
		acceptance.Write(t, repo, map[string]string{swiftTest: edited(fixed[swiftTest])})
		acceptance.Git(t, repo, "add", "-A")
		out, code = acceptance.Fitness(t, repo, amending, "requirements")
		acceptance.Sees(t, out, code, 0)
		out, code = stagedRun(t, swiftRepo(), map[string]string{swiftTest: edited(swiftTestingFile("test0001_1"))})
		acceptance.Sees(t, out, code, 1, swiftTest+" changed without its requirement doc "+acceptance.PairedDoc+"; change both together")
	})
	t.Run("0006.1.1.7.1 it needs no partner", func(t *testing.T) {
		t.Parallel()
		given := acceptance.With(acceptance.PairedRepo(), map[string]string{acceptance.PairedDoc: unformatted.Replace(acceptance.RequirementDoc("0001.1"))})
		if !strings.Contains(given[acceptance.PairedDoc], "\n    1. GIVEN") {
			t.Fatalf("the unformatted doc must nest its numbered items apart from the formatted one:\n%s", given[acceptance.PairedDoc])
		}
		formatted := map[string]string{acceptance.PairedDoc: acceptance.RequirementDoc("0001.1")}
		out, code := stagedRun(t, given, formatted)
		acceptance.Sees(t, out, code, 0)
		out, code = prRun(t, given, formatted, "style(docs): format requirement docs with Biome")
		acceptance.Sees(t, out, code, 0)
		out, code = stagedRun(t, given, acceptance.With(formatted, map[string]string{pairedTest: edited(acceptance.TestFile("reqs", "Test0001_1"))}))
		acceptance.Sees(t, out, code, 0)
		sourced := acceptance.With(acceptance.PairedRepo(), map[string]string{acceptance.PairedDoc: acceptance.WithSource(acceptance.RequirementDoc("0001.1"), "`adder.go:5`")})
		dropped := map[string]string{acceptance.PairedDoc: acceptance.RequirementDoc("0001.1")}
		out, code = stagedRun(t, sourced, dropped)
		acceptance.Sees(t, out, code, 0)
		out, code = prRun(t, sourced, dropped, "docs(requirements): drop the Source line")
		acceptance.Sees(t, out, code, 0)
		out, code = stagedRun(t, sourced, map[string]string{acceptance.PairedDoc: strings.Replace(acceptance.RequirementDoc("0001.1"), "sums asked for", "sums tried", 1)})
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+" changed without its requirement test")
		cited := map[string]string{acceptance.PairedDoc: acceptance.WithSource(acceptance.RequirementDoc("0001.1"), "`adder.go:6`")}
		out, code = stagedRun(t, sourced, cited)
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+" changed without its requirement test")
		commentedPairs(t)
	})
	t.Run("0006.1.1.7.2 a reformat or rename of numbered acceptances needs no partner", func(t *testing.T) {
		t.Parallel()
		renamed := map[string]string{acceptance.PairedDoc: strings.Replace(acceptance.RequirementDoc("0001.1"), "Adds pair 1", "Sums a pair", 1)}
		out, code := stagedRun(t, acceptance.PairedRepo(), renamed)
		acceptance.Sees(t, out, code, 0)
		changedGiven := map[string]string{acceptance.PairedDoc: strings.Replace(acceptance.RequirementDoc("0001.1"), "GIVEN two numbers", "GIVEN three numbers", 1)}
		out, code = stagedRun(t, acceptance.PairedRepo(), changedGiven)
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+" changed without its requirement test")
		renamedPairs(t)
	})
}
