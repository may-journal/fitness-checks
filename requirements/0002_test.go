package requirements

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func Test0002_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"describe patterns compile for every check", func(t *testing.T) {
			t.Parallel()
			all, err := describedChecks()
			mustDo(t, err)
			if len(all) == 0 {
				t.Fatal("the build made no check programs")
			}
			for name, d := range all {
				if _, err := checkkit.PatternsOf(d); err != nil {
					t.Errorf("%s: %v", name, err)
				}
				if r := checkkit.UnknownRequirement(d); r != "" {
					t.Errorf("%s requires %q, which fitness-runner does not define", name, r)
				}
			}
		}},
		{"file checks reach only their files", func(t *testing.T) {
			t.Parallel()
			patternReaches(t, "gofmt", "pkg/run.go", true)
			patternReaches(t, "gofmt", "vendor/x/run.go", false)
			patternReaches(t, "gofmt", "README.md", false)
			patternReaches(t, "markdown-links", "docs/guide.md", true)
			patternReaches(t, "markdown-links", "pkg/run.go", false)
			patternReaches(t, "workflow-reporting", ".github/workflows/ci.yml", true)
			patternReaches(t, "workflow-reporting", "docs/ci.yml", false)
			patternReaches(t, "workflow-scripts", ".github/workflows/ci.yml", true)
			patternReaches(t, "workflow-scripts", "tools/setup/action.yml", true)
			patternReaches(t, "workflow-scripts", "docs/ci.yml", false)
			eslint := patternsOf(t, "eslint")
			if !eslint.Met([]string{"web/package.json"}) || eslint.Met([]string{"README.md"}) {
				t.Errorf("eslint must need a package.json and nothing else")
			}
		}},
		{"body and message checks reach their paths", func(t *testing.T) {
			t.Parallel()
			for check, path := range map[string]string{
				"pr-structure":       checkkit.PullBodyPath,
				"pr-closes-issue":    checkkit.PullBodyPath,
				"plan-structure":     checkkit.IssueBodyPath("Plan"),
				"epic-structure":     checkkit.IssueBodyPath("Epic"),
				"issue-checklist":    checkkit.IssueBodyPath("Plan"),
				"comment-structure":  checkkit.CommentBodyPath,
				"semantic-commit":    checkkit.CommitMessagePath,
				"plan-trailer":       checkkit.CommitMessagePath,
				"commit-attribution": checkkit.CommitMessagePath,
				"issue-link-once":    checkkit.CommitMessagePath,
			} {
				patternReaches(t, check, path, true)
				patternReaches(t, check, "pkg/run.go", false)
				patternReaches(t, check, checkkit.BlobPath, !strings.HasPrefix(path, ".git/"))
				if all, err := describedChecks(); err == nil && all[check].MinFitness != "1.75.0" {
					t.Errorf("%s needs fitness %q, want 1.75.0, the first with virtual paths", check, all[check].MinFitness)
				}
			}
			patternReaches(t, "issue-checklist", checkkit.IssueBodyPath(""), true)
			patternReaches(t, "pr-structure", checkkit.IssueBodyPath("Plan"), false)
			patternReaches(t, "plan-structure", checkkit.IssueBodyPath("Epic"), false)
			patternReaches(t, "semantic-commit", checkkit.HistoryPath, true)
			patternReaches(t, "semantic-commit", checkkit.PullBodyPath, false)
			patternReaches(t, "comment-structure", checkkit.PullBodyPath, false)
		}},
		{"prose checks reach the right files", func(t *testing.T) {
			t.Parallel()
			for _, check := range []string{"prose-budget", "text-readability", "markdown-no-bold-italic", "mermaid-syntax", "open-questions"} {
				for _, path := range []string{"docs/guide.md", checkkit.PullBodyPath, checkkit.CommentBodyPath, checkkit.BlobPath} {
					patternReaches(t, check, path, true)
				}
				patternReaches(t, check, "web/app.css", false)
			}
			for _, check := range []string{"prose-budget", "text-readability", "no-contrastive-reframing"} {
				patternReaches(t, check, "pkg/run.go", true)
				patternReaches(t, check, "scripts/deploy", true)
				patternReaches(t, check, checkkit.CommitMessagePath, false)
			}
			patternReaches(t, "no-contrastive-reframing", checkkit.PullBodyPath, false)
		}},
		{"release build rejects bad pattern", func(t *testing.T) {
			t.Parallel()
			program := filepath.Join(t.TempDir(), "describe-patterns")
			build := exec.Command("go", "build", "-o", program, "./scripts/describe-patterns")
			build.Dir = ".."
			if out, err := build.CombinedOutput(); err != nil {
				t.Fatalf("go build: %v\n%s", err, out)
			}
			out, code := patternsExit(t, program, `{"name":"broken","files":"(oops"}`)
			if code != 1 || !strings.Contains(out, `fitness-check-x-v1.0.0: files pattern "(oops"`) {
				t.Errorf("a broken files pattern exited %d: %s", code, out)
			}
			all, err := describedChecks()
			mustDo(t, err)
			gofmt, err := json.Marshal(all["gofmt"])
			mustDo(t, err)
			if out, code := patternsExit(t, program, string(gofmt)); code != 0 {
				t.Errorf("gofmt's describe line exited %d: %s", code, out)
			}
		}},
	})
}

func Test0002_2(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"repo with nothing to judge passes", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			out, _ := fitness(t, repo, nil)
			for _, check := range []string{"eslint", "prettier", "swiftlint", "localized-text", "gofmt", "go-test", "requirements", "requirements-changed"} {
				passedWithNoFiles(t, out, check)
			}
			fixtures := example(t, "happyRepo", map[string]string{"tool/testdata/bad.go": "package bad\nfunc  X( ) {}\n"})
			out, _ = fitness(t, fixtures, nil, "gofmt")
			passedWithNoFiles(t, out, "gofmt")
		}},
		{"result rows list every check once", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			out, _ := fitness(t, repo, nil)
			cells := rowCells(out)
			want := strings.Fields(catalogRows)
			if len(cells) != len(want) {
				t.Fatalf("fitness printed %d rows, want %d:\n%s", len(cells), len(want), out)
			}
			for i, check := range want {
				if !names(cells[i], check) {
					t.Errorf("row %d is %q, want %s", i+1, cells[i], check)
				}
			}
		}},
		{"local-branch declares itself and exits alone", func(t *testing.T) {
			t.Parallel()
			bin := goBuild(t, "./cmd/fitness-check-local-branch") + "/fitness-check-local-branch"
			out, code := execute(t, func() *exec.Cmd { return exec.Command(bin, "--describe") })
			sees(t, out, code, 0, `"name":"local-branch"`, `"scope":"github"`, `"protocol":1`)
			out, code = execute(t, func() *exec.Cmd { return exec.Command(bin, "--root", t.TempDir()) })
			sees(t, out, code, 3, "runs only inside fitness")
		}},
	})
}

func Test0002_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"unreadable file is skipped and named as info", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{
				"README.md": readme + "\nSee [the guide](gone.md).\n",
				"locked.md": "# Locked\n",
			})
			mustDo(t, os.Chmod(filepath.Join(repo, "locked.md"), 0o000))
			out, code := fitness(t, repo, nil, "markdown-links")
			sees(t, out, code, 1, "locked.md: skipped, cannot be read (permission denied)", "README.md:9: broken relative link: gone.md")
			if !strings.Contains(row(out, "markdown-links"), "failed") {
				t.Errorf("the readable file's finding must still fail, got row %q", row(out, "markdown-links"))
			}
		}},
		{"service that cannot answer passes with info", func(t *testing.T) {
			t.Parallel()
			out, code, _ := refCheck(t, refWorkflow("may-journal/fitness-runner/.github/workflows/ci-reusable.yml@v1"), []response{
				{Stderr: "error connecting to api.github.com\n", Exit: 1},
			})
			sees(t, out, code, 0, "workflow refs were not verified: gh api repos/may-journal/fitness-runner/git/ref/tags/v1: error connecting to api.github.com")
			if strings.Contains(row(out, "workflow-ref-exists"), "failed") {
				t.Errorf("a check whose service cannot answer must not fail, got row %q", row(out, "workflow-ref-exists"))
			}
		}},
		{"text checks skip unreadable files", func(t *testing.T) {
			t.Parallel()
			for _, c := range lockedChecks {
				t.Run(c.check, func(t *testing.T) {
					t.Parallel()
					repo := example(t, "happyRepo", c.files)
					mustDo(t, os.Chmod(filepath.Join(repo, filepath.FromSlash(c.locked)), 0o000))
					out, code := fitness(t, repo, nil, c.check)
					sees(t, out, code, 1, c.locked+": skipped, cannot be read (permission denied)", c.finding)
					if !strings.Contains(row(out, c.check), "failed") {
						t.Errorf("the readable file's finding must still fail, got row %q", row(out, c.check))
					}
				})
			}
		}},
		{"go parse error fails beside unreadable file", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"go.mod": "module example.com/m\n\ngo 1.22\n", "open.go": "package m\n\nfunc Open( {\n", "locked.go": "package m\n"})
			mustDo(t, os.Chmod(filepath.Join(repo, "locked.go"), 0o000))
			out, code := fitness(t, repo, nil, "go-complexity")
			sees(t, out, code, 1, "locked.go: skipped, cannot be read (permission denied)", "open.go: ")
			if !strings.Contains(row(out, "go-complexity"), "failed") {
				t.Errorf("a file that does not parse must still fail, got row %q", row(out, "go-complexity"))
			}
		}},
		{"unreadable requirement doc is named as info", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", goRepo())
			mustDo(t, os.Chmod(filepath.Join(repo, "docs", "requirements", "0001-adds-numbers.md"), 0o000))
			out, code := fitness(t, repo, nil, "requirements")
			showsInfo(t, out, code, "requirements", "docs/requirements/0001-adds-numbers.md: skipped, cannot be read (permission denied)")
		}},
		{"history git cannot replay passes with info", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", goRepo())
			head, err := exec.Command("git", "-C", repo, "rev-parse", "HEAD").Output()
			mustDo(t, err)
			keep := strings.TrimSpace(string(head))
			mustDo(t, filepath.WalkDir(filepath.Join(repo, ".git", "objects"), func(p string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() && filepath.Base(filepath.Dir(p))+filepath.Base(p) != keep {
					err = os.Remove(p)
				}
				return err
			}))
			out, code := fitness(t, repo, nil, "requirements")
			showsInfo(t, out, code, "requirements", "requirement IDs were not checked against git history")
		}},
	})
}

func Test0002_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"read-repo-first passes with zero files", func(t *testing.T) {
			t.Parallel()
			out, code := readRepoFirst(t)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "read-repo-first")
		}},
		{"asks whether decisions and checks were read", func(t *testing.T) {
			t.Parallel()
			out, code := readRepoFirst(t)
			sees(t, out, code, 0, "Did you familiarize yourself with the decisions logged in the repo, "+
				`specifically all "Fitness Checks" that are enabled via @mayjournal/fitness?`)
		}},
		{"warns never to skip hooks", func(t *testing.T) {
			t.Parallel()
			out, code := readRepoFirst(t)
			sees(t, out, code, 0, "NOTE: Do not under any circumstance use `--no-verify`")
		}},
		{"table lists only enabled checks", func(t *testing.T) {
			t.Parallel()
			out, code := readRepoFirstChosen(t, "duplicate-code")
			sees(t, out, code, 0, "Check │ Src", "duplicate-code │ go/cmd/fitness-check-duplicate-code")
			if strings.Contains(out, "fitness-check-gofmt") {
				t.Errorf("table lists gofmt, which this run did not enable:\n%s", out)
			}
		}},
		{"long cells are shortened with ellipsis", func(t *testing.T) {
			t.Parallel()
			out, code := readRepoFirstChosen(t, "markdown-filename-kebab-case")
			sees(t, out, code, 0, "markdown-filename-kebab-c… │ go/cmd/fitness-check-markdown-file"+"n…")
		}},
		{"no color when color is disabled", func(t *testing.T) {
			t.Parallel()
			colored := readRepoFirstRaw(t, "NO_COLOR=")
			plain := readRepoFirstRaw(t, "NO_COLOR=1")
			if !ansi.MatchString(colored) {
				t.Errorf("without NO_COLOR the table should be colored:\n%q", colored)
			}
			if ansi.MatchString(plain) || !strings.Contains(plain, "│ Check") {
				t.Errorf("with NO_COLOR the table should print plain:\n%q", plain)
			}
		}},
	})
}

func Test0002_5(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"fitness key fails once per file", func(t *testing.T) {
			t.Parallel()
			out, code := noFitnessFront(t, map[string]string{
				"docs/guide.md":                    fitnessBlock + "# Guide\n",
				"changelog/2026-10-11-0900-one.md": entryFileText("2026.10.11.0900", "Feat: one change."),
				".github/PULL_REQUEST_TEMPLATE.md": "<!--\n---\nfitnessFunctions: ['pr-structure']\n---\n-->\n<!-- One sentence. -->\n> REPLACE-ME\n",
			})
			sees(t, out, code, 1)
			flagged(t, out, "docs/guide.md", "changelog/2026-10-11-0900-one.md", ".github/PULL_REQUEST_TEMPLATE.md")
		}},
		{"crlf and byte order mark still fail", func(t *testing.T) {
			t.Parallel()
			out, code := noFitnessFront(t, map[string]string{
				"docs/crlf.md":                  "---\r\nrelatedConfigurations: ['../.fitnessrc.json']\r\n---\r\n\r\n# Crlf\r\n",
				"docs/bom.md":                   byteOrderMark + "---\nfitnessFunctions: ['prose-budget']\n---\n\n# Bom\n",
				".github/ISSUE_TEMPLATE/bug.md": "---\nname: Bug\nrelatedConfigurations: ['../../.fitnessrc.json']\nlabels: Bug\n---\n\n## Steps\n",
			})
			sees(t, out, code, 1)
			flagged(t, out, "docs/crlf.md", "docs/bom.md", ".github/ISSUE_TEMPLATE/bug.md")
		}},
		{"other front matter and fixtures pass", func(t *testing.T) {
			t.Parallel()
			out, code := noFitnessFront(t, map[string]string{
				"docs/notes.md":                     "---\ntitle: Notes\n---\n\n# Notes\n",
				"docs/ruled.md":                     "---\n\nrelatedConfigurations: a note set between two rules.\n\n---\n\n# Ruled\n",
				".github/ISSUE_TEMPLATE/bug.md":     "---\nname: Bug\nlabels: Bug\n---\n\n## Steps\n",
				"tests/examples/demoRepo/README.md": fitnessBlock + "# Demo\n",
				"go/tool/testdata/guide.md":         fitnessBlock + "# Guide\n",
			})
			sees(t, out, code, 0, "All 1 checks passed")
			if strings.Contains(out, "strip-front-matter") {
				t.Errorf("a file without fitness front matter, or a fixture, was flagged:\n%s", out)
			}
		}},
	})
}

func Test0002_6(t *testing.T) {
	runCases(t, []edgeCase{
		{"large diff starts git few times", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", with(goRepo(), counters("")))
			write(t, repo, largeDiff())
			git(t, repo, "add", "-A")
			path, calls := loggedGit(t)
			out, code := fitness(t, repo, []string{path}, "--policy=external", "--checks="+strings.Join(largeDiffChecks, ","))
			passesBatched(t, out, code, largeDiffChecks, calls)
			pull := example(t, "happyRepo", with(goRepo(), counters("")))
			git(t, pull, "update-ref", "refs/remotes/origin/main", "HEAD")
			commit(t, pull, largeDiff(), "fix(counters): add an offset")
			path, calls = loggedGit(t)
			out, code = fitness(t, pull, append([]string{path}, prEnv...), "--policy=external", "--checks="+strings.Join(largeDiffChecks[:2], ","))
			passesBatched(t, out, code, largeDiffChecks[:2], calls)
		}},
		{"hundreds of docs start git few times", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", proseRepo())
			path, calls := loggedGit(t)
			out, code := fitness(t, repo, []string{path}, "--policy=external", "--checks="+strings.Join(proseChecks, ","), "--all")
			t.Log("\n" + out)
			sees(t, out, code, 0, "All 5 checks passed")
			n := occurrences(calls, "\n")
			t.Logf("git started %d times", n)
			if n > maxProseGitCalls {
				t.Errorf("the prose checks started git %d times on %d files, want at most %d, none per file", n, 2*proseFiles, maxProseGitCalls)
			}
		}},
		{"long branch needs one git log call", func(t *testing.T) {
			t.Parallel()
			repo, pr := longBranch(t)
			rp := standIn(t, map[string][]response{"gh": {pr, ghIssue(t, 11, planBody("- [x] Write the tests."), "Plan")}})
			defer func() {
				if rp.called("gh", "repo view") {
					t.Errorf("local-branch asked gh for the repo, which the PR answer names: %q", rp.calls("gh"))
				}
			}()
			gitPath, calls := loggedGit(t)
			gitDir, _, _ := strings.Cut(strings.TrimPrefix(gitPath, "PATH="), string(os.PathListSeparator))
			env := append([]string{"PATH=" + gitDir + string(os.PathListSeparator) + strings.TrimPrefix(rp.env[0], "PATH=")}, rp.env[1:]...)
			out, code := fitness(t, repo, env, "--policy=external", "--checks="+strings.Join(commitChecks, ","))
			sees(t, out, code, 0, "All 3 checks passed")
			// The fast total bounds the time. A long branch costs git log calls, one per
			// check at any length.
			if logs := gitLogCalls(t, calls); logs > len(commitChecks) {
				t.Errorf("the commit checks ran git log %d times on a branch of %d commits, want at most one each", logs, longBranchCommits)
			}
		}},
		{"large tree passes with few git calls", func(t *testing.T) {
			// Not parallel: a 50,000-file tree built beside other tests would time their disk load.
			repo := generatedTree(t, largeTree(), largeTreeFiles)
			write(t, repo, map[string]string{"dist/index.js": "export const built = true;\n"})
			path, calls := loggedGit(t)
			out, code := fitness(t, repo, []string{path}, "--all", "--policy=external", "--checks="+strings.Join(largeTreeChecks, ","))
			sees(t, out, code, 0, "All 2 checks passed")
			fewGitCalls(t, calls)
		}},
		{"nested dist file fails by name", func(t *testing.T) {
			// Not parallel, like the "large tree passes with few git calls" case.
			repo := generatedTree(t, largeTree(), largeTreeFiles)
			write(t, repo, map[string]string{nestedDist: "export const built = true;\n"})
			git(t, repo, "add", "-f", nestedDist)
			git(t, repo, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "build(tree): track a build")
			path, calls := loggedGit(t)
			out, code := fitness(t, repo, []string{path}, "--all", "--policy=external", "--checks="+largeTreeChecks[0])
			sees(t, out, code, 1, nestedDist)
			fewGitCalls(t, calls)
		}},
	})
}

// cspell:ignore zorbleflux

// catalogRows is the result row order: sorted repo-wide checks, then each
// sorted Go, Swift, Node and Java module's checks.
const catalogRows = `
adr-as-issue adr-checks
architecture-actors architecture-code architecture-components
architecture-containers architecture-contract architecture-journey
build-output-untracked changelog changelog-bullets changelog-updated
comment-structure commit-attribution cspell doc-template duplicate-code
epic-structure gitignore-why issue-checklist issue-link-once local-branch
markdown-filename-camel-case
markdown-filename-kebab-case markdown-links
markdown-no-bold-italic mermaid-callout-why mermaid-callouts
mermaid-diagram-prose mermaid-diagram-table-gap mermaid-legend
mermaid-level-bleed mermaid-syntax no-contrastive-reframing no-fitness-front-matter
no-plans-dir open-questions plan-structure
plan-trailer pr-closes-issue pr-structure prose-budget read-repo-first
release-changelog repeated-string-literals requirements requirements-changed
semantic-commit swift-storybook-requirements text-readability workflow-ref-exists workflow-reporting
workflow-scripts
go-complexity go-test go-vet gofmt
localized-text swiftlint
dependency-currency eslint no-eslint-disable node-version prettier
vitest-coverage-exclude vitest-coverage-full
checkstyle java-test java-version spotless`

// allPassed is the summary of a run where every installed check passed.
func allPassed() string {
	return fmt.Sprintf("All %d checks passed", len(strings.Fields(catalogRows)))
}

// passedOf is the summary count of a run where all but notPassed installed
// checks passed.
func passedOf(notPassed int) string {
	n := len(strings.Fields(catalogRows))
	return fmt.Sprintf("%d of %d checks passed", n-notPassed, n)
}

// rowCells returns the check cell of each results-table row, top to bottom.
func rowCells(out string) []string {
	var cells []string
	for _, line := range strings.Split(out, "\n") {
		if m := resultRow.FindStringSubmatch(line); m != nil {
			cells = append(cells, m[1])
		}
	}
	return cells
}

// unknownWord is a word no dictionary holds, built so cspell never reads it here.
var unknownWord = strings.Repeat("q", 3) + "xw" + strings.Repeat("v", 3)

// lockedChecks are the file-reading checks, each paired with one readable file
// it must still judge and one it cannot read.
var lockedChecks = []struct {
	check, finding, locked string
	files                  map[string]string
}{
	{"no-eslint-disable", "src/open.ts:1: eslint-disable", "src/locked.ts", map[string]string{
		"src/open.ts": "/* eslint-disable */\n", "src/locked.ts": "export const a = 1;\n"}},
	{"repeated-string-literals", "shared-literal-value", "src/locked.ts", map[string]string{
		"src/a.ts":      "export const a = 'shared-literal-value';\nexport const b = 'shared-literal-value';\n",
		"src/b.ts":      "export const c = 'shared-literal-value';\n",
		"src/locked.ts": "export const d = 1;\n"}},
	{"cspell", "Unknown word (" + unknownWord + ")", "locked.txt", map[string]string{
		"open.txt": unknownWord + "\n", "locked.txt": "hello\n"}},
	{"prose-budget", "docs/open.md", "docs/locked.md", map[string]string{
		"docs/open.md":   "# Open\n\n" + strings.Repeat("This sentence has some words in it. ", 40) + "\n",
		"docs/locked.md": "# Locked\n\nShort.\n"}},
	{"text-readability", "docs/open.md", "docs/locked.md", map[string]string{
		"docs/open.md":   "# Open\n\n" + strings.Repeat("Notwithstanding interdepartmental organizational responsibilities institutional administrators characteristically misunderstand straightforwardly communicated recommendations notwithstanding counterrevolutionary disproportionately uncharacteristically incomprehensible institutionalization. ", 10) + "\n",
		"docs/locked.md": "# Locked\n\nShort.\n"}},
	{"go-complexity", "open.go", "locked.go", map[string]string{
		"go.mod":    "module example.com/m\n\ngo 1.22\n",
		"open.go":   "package m\n\nfunc Open(a int) int {\n\tif a > 1 {\n\t\treturn 1\n\t}\n\tif a > 2 {\n\t\treturn 2\n\t}\n\tif a > 3 {\n\t\treturn 3\n\t}\n\tif a > 4 {\n\t\treturn 4\n\t}\n\tif a > 5 {\n\t\treturn 5\n\t}\n\tif a > 6 {\n\t\treturn 6\n\t}\n\treturn 0\n}\n",
		"locked.go": "package m\n"}},
}

// readRepoFirst runs `fitness-install -- <args>` on happyRepo; with no args
// it runs read-repo-first alone.
func readRepoFirst(t *testing.T, args ...string) (string, int) {
	t.Helper()
	if len(args) == 0 {
		args = []string{"read-repo-first"}
	}
	return fitness(t, example(t, "happyRepo", nil), nil, args...)
}

// readRepoFirstChosen runs read-repo-first beside the checks a user chose.
func readRepoFirstChosen(t *testing.T, checks string) (string, int) {
	t.Helper()
	return readRepoFirst(t, "--policy=external", "--checks=read-repo-first,"+checks)
}

// readRepoFirstRaw runs read-repo-first on happyRepo with env added and
// returns its banner, up to the note, exactly as printed with color codes.
func readRepoFirstRaw(t *testing.T, env ...string) string {
	t.Helper()
	cmd := exec.Command(installer(t), "--", "read-repo-first")
	cmd.Dir, cmd.Env = example(t, "happyRepo", nil), append(userEnv(), env...)
	out, err := cmd.CombinedOutput()
	mustDo(t, err)
	banner, _, _ := strings.Cut(string(out), "NOTE:")
	return banner
}

// noFitnessFront writes files over happyRepo and runs
// `fitness no-fitness-front-matter` there.
func noFitnessFront(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "no-fitness-front-matter")
}

// frontFix is what the check prints after each file that holds a key.
const frontFix = ": front matter holds fitnessFunctions or relatedConfigurations, which no check reads; run `fitness strip-front-matter`"

// byteOrderMark is the UTF-8 BOM some editors put before a file's first line.
const byteOrderMark = string(rune(0xFEFF))

// fitnessBlock is a plain block holding both fitness keys.
const fitnessBlock = "---\nfitnessFunctions:\n  - prose-budget\nrelatedConfigurations: ['../.fitnessrc.json']\n---\n\n"

// flagged asserts out names each file once with the fix.
func flagged(t *testing.T, out string, files ...string) {
	t.Helper()
	for _, f := range files {
		if n := strings.Count(read(out), read(f+frontFix)); n != 1 {
			t.Errorf("%s is named %d times, want once:\n%s", f, n, out)
		}
	}
}

// largeDiffFiles is how many product files the large diff edits, so with its
// test, doc and changelog entry it changes 266 files.
const largeDiffFiles = 263

// largeDiffChecks are the checks a large staged diff must keep fast.
var largeDiffChecks = []string{"requirements-changed", "requirements", "changelog-updated"}

// counter is product file i of the large diff, returning i plus step.
func counter(i int, step string) string {
	return fmt.Sprintf("package counters\n\n// Counter%03d returns its number.\nfunc Counter%03d() int {\n\treturn %d%s\n}\n", i, i, i, step)
}

// counters lists the large diff's product files, each returning its number
// plus step.
func counters(step string) map[string]string {
	files := map[string]string{}
	for i := range largeDiffFiles {
		files[fmt.Sprintf("counters/counter%03d.go", i)] = counter(i, step)
	}
	return files
}

// largeDiff is every counter's code changed, with goRepo's test, doc and a
// changelog entry that names the change.
func largeDiff() map[string]string {
	return with(counters(" + 1 // counter offset"), map[string]string{
		"adder_test.go": edited(goRepo()["adder_test.go"]),
		pairedDoc:       reworded(goRepo()[pairedDoc]),
		"changelog/2026-11-08-1400-counter-offset.md": bareEntryText("2026.11.08.1400", "Fix: every counter will return its number plus an offset, by hand."),
	})
}

// maxLargeDiffGitCalls bounds the git calls of a whole run on the large diff,
// runner included, far below one per changed file.
const maxLargeDiffGitCalls = 40

// passesBatched asserts checks passed the large diff, and the run started git
// a few times in all, not once per file.
func passesBatched(t *testing.T, out string, code int, checks []string, calls string) {
	t.Helper()
	sees(t, out, code, 0, fmt.Sprintf("All %d checks passed", len(checks)))
	n := occurrences(calls, "\n")
	t.Logf("git started %d times", n)
	if n == 0 || n > maxLargeDiffGitCalls {
		t.Errorf("the run started git %d times on a %d-file diff, want 1 to %d, none per file", n, len(largeDiff()), maxLargeDiffGitCalls)
	}
}

// proseChecks are the checks that read every doc and comment.
var proseChecks = []string{"text-readability", "prose-budget", "cspell", "mermaid-callouts", "repeated-string-literals"}

// proseFiles is how many docs, and how many commented Go files, proseRepo
// holds.
const proseFiles = 150

// maxProseGitCalls bounds the git calls of the whole run, runner included,
// far below one per file; a run makes about seven.
const maxProseGitCalls = 20

// proseRepo is goRepo with proseFiles short docs and as many commented Go
// files, each with its own string literal.
func proseRepo() map[string]string {
	files := goRepo()
	for i := range proseFiles {
		files[fmt.Sprintf("docs/notes/note%03d.md", i)] = fmt.Sprintf("# Note %d\n\nThis note keeps the plan for the garden. We water the plants each morning.\n", i)
		files[fmt.Sprintf("garden/plant%03d.go", i)] = fmt.Sprintf("package garden\n\n// Plant%03d names one plant in the garden.\nconst Plant%03d = \"plant %d\"\n", i, i, i)
	}
	return files
}

// longBranchCommits is how many commits longBranch adds to the branch.
const longBranchCommits = 1000

// commitChecks are the checks that read a branch's commits.
var commitChecks = []string{"commit-attribution", "semantic-commit", "local-branch"}

// longBranch is branchRepo with longBranchCommits commits that one
// fast-import call writes, each with a conventional subject and AI trailers.
// Its open PR is based before them and names its base repo, as GitHub does.
func longBranch(t *testing.T) (string, response) {
	t.Helper()
	repo := branchRepo(t, nil)
	pr := branchAnswer(t, repo, "test(app): add one file", prClosing(11))
	var answer map[string]any
	mustDo(t, json.Unmarshal([]byte(pr.Stdout), &answer))
	answer["repo"] = exampleRepo
	data, err := json.Marshal(answer)
	mustDo(t, err)
	pr.Stdout = string(data) + "\n"
	var stream strings.Builder
	for i := range longBranchCommits {
		msg := fmt.Sprintf("test(app): add step %d\n\nAI-Tools: Claude Code\nAI-Models: Claude\n", i)
		fmt.Fprintf(&stream, "commit refs/heads/work\ncommitter t <t@t> %d +0000\ndata %d\n%s\n", 1700000000+i, len(msg), msg)
		if i == 0 {
			stream.WriteString("from refs/remotes/origin/main\n\n")
		}
	}
	gitInput(t, repo, stream.String(), "fast-import", "--quiet")
	return repo, pr
}

// gitLogCalls counts the git log calls loggedGit kept. None fails the test,
// since then nothing read the branch.
func gitLogCalls(t *testing.T, calls string) int {
	t.Helper()
	raw, _ := os.ReadFile(calls)
	n := 0
	for _, line := range strings.Split(string(raw), "\n") {
		if slices.Contains(strings.Fields(line), "log") {
			n++
		}
	}
	if n == 0 {
		t.Fatalf("no check ran git log through the logged git:\n%s", raw)
	}
	return n
}

// largeTreeFiles is how many generated files the large tree adds beside a
// Java project and a components view.
const largeTreeFiles = 50000

// largeTreeChecks are the repo-state checks a large tree must keep fast.
var largeTreeChecks = []string{"build-output-untracked", "architecture-components"}

// nestedDist is a build output file deep in the large tree.
const nestedDist = "tree/0049/dist/app.js"

// largeTree is aoFixture with a Java project, a TypeScript source the store
// owns and an ignored dist folder.
func largeTree() map[string]string {
	return with(aoFixture(), map[string]string{
		"pom.xml":            javaPom(""),
		".gitignore":         "dist/\n",
		"src/store/index.ts": "import { save } from './save';\n\nexport const ready = save;\n",
		"src/store/save.ts":  "export const save = true;\n",
	})
}

// maxLargeTreeGitCalls bounds the git calls of a whole run on the large tree,
// runner included, far below one per tracked file.
const maxLargeTreeGitCalls = 20

// fewGitCalls asserts the run started git a few times in all, with one
// `git ls-files` for build-output-untracked, not one per file.
func fewGitCalls(t *testing.T, calls string) {
	t.Helper()
	n := occurrences(calls, "\n")
	t.Logf("git started %d times", n)
	if n == 0 || n > maxLargeTreeGitCalls {
		t.Errorf("the run started git %d times on %d tracked files, want 1 to %d, none per file", n, largeTreeFiles, maxLargeTreeGitCalls)
	}
}

// describedChecks holds each check program's --describe metadata, built once.
var describedChecks = sync.OnceValues(func() (map[string]checkkit.Describe, error) {
	dir, err := os.MkdirTemp("", "checks-0169-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	build := exec.Command("go", "build", "-o", dir+string(filepath.Separator), "./cmd/...")
	build.Dir = ".."
	if out, err := build.CombinedOutput(); err != nil {
		return nil, errors.New("go build ./cmd/...: " + string(out))
	}
	return describeAll(dir)
})

// describeAll runs --describe on each check program in dir, by name.
func describeAll(dir string) (map[string]checkkit.Describe, error) {
	entries, err := os.ReadDir(dir)
	all := map[string]checkkit.Describe{}
	for _, e := range entries {
		if name, ok := strings.CutPrefix(e.Name(), checkkit.BinPrefix); ok && err == nil {
			all[name], err = describeLine(filepath.Join(dir, e.Name()))
		}
	}
	return all, err
}

// describeLine decodes the --describe line program prints.
func describeLine(program string) (checkkit.Describe, error) {
	var d checkkit.Describe
	out, err := exec.Command(program, "--describe").Output()
	if err == nil {
		err = json.Unmarshal(out, &d)
	}
	return d, err
}

// patternsOf returns check's compiled patterns, failing the test when it
// is missing or its patterns do not compile.
func patternsOf(t *testing.T, check string) checkkit.Patterns {
	t.Helper()
	all, err := describedChecks()
	mustDo(t, err)
	d, ok := all[check]
	if !ok {
		t.Fatalf("no check program describes %s", check)
	}
	p, err := checkkit.PatternsOf(d)
	mustDo(t, err)
	return p
}

// patternReaches asserts whether check's patterns reach path.
func patternReaches(t *testing.T, check, path string, want bool) {
	t.Helper()
	if got := patternsOf(t, check).Matches(path); got != want {
		t.Errorf("%s reaches %s = %v, want %v", check, path, got, want)
	}
}

// patternsExit pipes line to the describe-patterns program and returns
// what it printed and its exit code.
func patternsExit(t *testing.T, program, line string) (string, int) {
	t.Helper()
	return execute(t, func() *exec.Cmd {
		cmd := exec.Command(program, "fitness-check-x-v1.0.0")
		cmd.Stdin = strings.NewReader(line)
		return cmd
	})
}

func Test0002_7(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"each package owns its docs and tests", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", monorepo(true)), nil, "requirements")
			sees(t, out, code, 0)
			out, code = fitness(t, example(t, "happyRepo", with(monorepo(true), map[string]string{
				packageDoc("a"): requirementDoc("0001.1", "0001.2"),
			})), nil, "requirements")
			sees(t, out, code, 1, packageDoc("a")+":21: 0001.2 has no test")
			out, code = fitness(t, example(t, "happyRepo", with(monorepo(true), map[string]string{
				packageDoc("b"): "", "pkgs/b/requirements/0001_test.go": "",
			})), nil, "requirements")
			sees(t, out, code, 1, "pkgs/b/docs/requirements/ has no requirement docs")
			out, code = fitness(t, example(t, "happyRepo", with(monorepo(true), map[string]string{"pkgs/c/README.md": "# Shared assets\n"})), nil, "requirements")
			sees(t, out, code, 0)
		}},
		{"a number repeats only across packages", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", with(monorepo(true), map[string]string{
				"pkgs/a/docs/requirements/0001-other.md": requirementDoc("0001.1"),
			})), nil, "requirements")
			sees(t, out, code, 1, "requirement 0001 is used by "+packageDoc("a")+", pkgs/a/docs/requirements/0001-other.md")
			if strings.Contains(out, "docs/requirements/0001-adds-numbers.md, pkgs/") {
				t.Errorf("the repo root's 0001 must not clash with a package's:\n%s", out)
			}
		}},
		{"a package change pairs with its own doc and test", func(t *testing.T) {
			t.Parallel()
			edit := map[string]string{"pkgs/a/a.go": monorepo(true)["pkgs/a/a.go"] + "\n// Two is two.\nconst Two = 2\n"}
			pair := func(doc, test string) map[string]string {
				return with(edit, map[string]string{doc: requirementDoc("0001.1", "0001.2"), test: testFile("requirements", "Test0001_1", "Test0001_2")})
			}
			out, code := rcStaged(t, monorepo(true), pair(pairedDoc, "adder_test.go"))
			sees(t, out, code, 1, "pkgs/a/a.go changed without a test and a requirement doc")
			out, code = rcStaged(t, monorepo(true), pair(packageDoc("a"), "pkgs/a/requirements/0001_test.go"))
			sees(t, out, code, 0, rcPassed)
			naming := with(monorepo(true), map[string]string{pairedDoc: strings.Replace(requirementDoc("0001.1"), "by hand.", "by hand in pkgs/a/a.go.", 1)})
			out, code = rcStaged(t, naming, map[string]string{"pkgs/a/a.go": ""})
			sees(t, out, code, 1, "pkgs/a/a.go was deleted, leaving requirement doc "+pairedDoc)
		}},
		{"a repo with no packages key is one package", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", monorepo(false)), nil, "requirements")
			sees(t, out, code, 1, "0001.1 is owned by 3 tests")
		}},
	})
}

// monorepo is goRepo with packages pkgs/a and pkgs/b, each with its own
// requirement 0001 and test, named in .fitnessrc.json when packages is set.
func monorepo(packages bool) map[string]string {
	rc := `{"policy": "external", "checks": ["requirements"]}`
	if packages {
		rc = `{"policy": "external", "checks": ["requirements"], "packages": ["pkgs/*"]}`
	}
	repo := with(goRepo(), map[string]string{".fitnessrc.json": rc + "\n"})
	for _, name := range []string{"a", "b"} {
		repo["pkgs/"+name+"/"+name+".go"] = "// Package " + name + " adds numbers.\npackage " + name + "\n"
		repo[packageDoc(name)] = requirementDoc("0001.1")
		repo["pkgs/"+name+"/requirements/0001_test.go"] = testFile("requirements", "Test0001_1")
	}
	return repo
}

// packageDoc is the path of requirement 0001 in package pkgs/<name>.
func packageDoc(name string) string {
	return "pkgs/" + name + "/docs/requirements/0001-adds-numbers.md"
}
