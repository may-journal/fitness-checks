package requirements

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Test0003_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"semantic-commit rejects a subject with no scope", func(t *testing.T) {
			t.Parallel()
			out, code := semanticCommit(t, "--message", "feat: add endpoint")
			sees(t, out, code, 1, `Commit message: "feat: add endpoint" — use type(scope): description `+semanticCommitTypes)
		}},
		{"semantic-commit rejects an empty message", func(t *testing.T) {
			t.Parallel()
			// happyRepo's HEAD subject is valid, so an empty message must not fall back to it.
			out, code := semanticCommit(t, "--message=")
			sees(t, out, code, 1, "No commit message to validate; use type(scope): description")
			noRepo := example(t, "happyRepo", nil)
			mustDo(t, os.RemoveAll(filepath.Join(noRepo, ".git")))
			out, code = fitness(t, noRepo, nil, "semantic-commit")
			sees(t, out, code, 1, "No commit message to validate; use type(scope): description")
		}},
		{"semantic-commit judges only the first line", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			messages := []string{"Merge branch 'feature' into main", "chore(deps): bump tools\n\nNot a semantic line at all"}
			for _, kind := range []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert"} {
				messages = append(messages, kind+"(api): change the endpoint")
			}
			for _, message := range messages {
				out, code := fitness(t, repo, nil, "semantic-commit", "--message", message)
				sees(t, out, code, 0, "All 1 checks passed")
			}
		}},
		{"semantic-commit rejects an unlisted type or breaking mark", func(t *testing.T) {
			t.Parallel()
			for _, subject := range []string{`Revert "feat(api): add endpoint"`, "feat(api)!: drop the old endpoint"} {
				out, code := semanticCommit(t, "--message", subject)
				sees(t, out, code, 1, `Commit message: "`+subject+`" — use type(scope): description`)
			}
		}},
		{"semantic-commit reads HEAD without a fetched origin", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			commit(t, repo, nil, "Add notes\n\nBody text")
			out, code := fitness(t, repo, nil, "semantic-commit")
			sees(t, out, code, 1, `Commit message: "Add notes" — use type(scope): description `+semanticCommitTypes)
		}},
		{"semantic-commit quotes only commits origin lacks", func(t *testing.T) {
			t.Parallel()
			for _, base := range defaultNames {
				repo := forkedRepo(t, base)
				out, code := fitness(t, repo, nil, "semantic-commit")
				sees(t, out, code, 0, "All 1 checks passed")
				commit(t, repo, nil, "Add notes")
				commit(t, repo, nil, "docs(readme): note more")
				out, code = fitness(t, repo, nil)
				quotesOnly(t, out, code, "Add notes")
				if !strings.Contains(row(out, "semantic-commit"), "failed") {
					t.Errorf("semantic-commit did not fail the full run on %s:\n%s", base, out)
				}
				git(t, repo, "checkout", "-q", "-B", base, "origin/"+base)
				commit(t, repo, nil, "Add base notes")
				out, code = fitness(t, repo, nil, "semantic-commit")
				quotesOnly(t, out, code, "Add base notes")
			}
		}},
		{"hooks accept a valid subject on a fresh branch", func(t *testing.T) {
			t.Parallel()
			for _, installerHook := range []bool{false, true} {
				repo, env := hookedFork(t, installerHook)
				out, code := user(t, repo, env, "commit", "--allow-empty", "-m", "docs(readme): touch nothing")
				sees(t, out, code, 0, "docs(readme): touch nothing")
				if strings.Contains(out, inherited) {
					t.Errorf("a hook judged the parent commit:\n%s", out)
				}
			}
		}},
		{"hooks refuse a subject with no scope", func(t *testing.T) {
			t.Parallel()
			repo, env := hookedFork(t, false)
			out, code := user(t, repo, env, "commit", "--allow-empty", "-m", "Add notes")
			quotesOnly(t, out, code, "Add notes")
			repo, env = hookedFork(t, true)
			out, code = user(t, repo, env, "commit", "--allow-empty", "-m", "Add notes")
			sees(t, out, code, 0, "Add notes")
			out, code = user(t, repo, env, "commit", "--allow-empty", "-m", "docs(readme): note more")
			quotesOnly(t, out, code, "Add notes")
		}},
	})
}

func Test0003_2(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"plan-trailer passes a message with no plan reference", func(t *testing.T) {
			t.Parallel()
			out, code := planTrailer(t, "feat(app): add a feature\n\nSome body.\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"plan-trailer passes a correct trailer", func(t *testing.T) {
			t.Parallel()
			out, code := planTrailer(t, "feat(app): add a feature\n\nSome body.\n\nPlan #63\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"plan-trailer rejects each misspelled plan line", func(t *testing.T) {
			t.Parallel()
			out, code := planTrailer(t, "feat(app): add a feature\n\nPlan 63\nplan #63\nPlan: 63\n")
			sees(t, out, code, 1,
				`plan trailer "Plan 63" `+planTrailerFix,
				`plan trailer "plan #63" `+planTrailerFix,
				`plan trailer "Plan: 63" `+planTrailerFix)
		}},
		{"plan-trailer passes body prose starting with Plan", func(t *testing.T) {
			t.Parallel()
			out, code := planTrailer(t, "feat(app): add a feature\n\nPlan the rollout in 3 steps.\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"plan-trailer reads HEAD when no message is given", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			commit(t, repo, nil, "feat(app): add a feature\n\nPlan: 63")
			out, code := fitness(t, repo, nil, "plan-trailer")
			sees(t, out, code, 1, `plan trailer "Plan: 63" `+planTrailerFix)
		}},
		{"plan-trailer passes a plan of another repo", func(t *testing.T) {
			t.Parallel()
			out, code := planTrailer(t, "feat(app): add a feature\n\nPlan may-journal/.github#144\n")
			sees(t, out, code, 0, "All 1 checks passed")
			out, code = planTrailer(t, "feat(app): add a feature\n\nPlan https://github.com/may-journal/.github/issues/144\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"plan-trailer rejects misspelled other repo plan lines", func(t *testing.T) {
			t.Parallel()
			out, code := planTrailer(t, "feat(app): add a feature\n\nPlan may-journal/.github 144\nPlan: may-journal/.github#144\n")
			sees(t, out, code, 1,
				`plan trailer "Plan may-journal/.github 144" `+planTrailerFix,
				`plan trailer "Plan: may-journal/.github#144" `+planTrailerFix)
		}},
	})
}

func Test0003_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"commit-attribution passes a repo with no trailers adopted", func(t *testing.T) {
			t.Parallel()
			out, code := commitAttribution(t, []string{"feat(api): add endpoint\n\nNo trailers here.\n"})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "commit-attribution")

			folder := example(t, "happyRepo", nil)
			mustDo(t, os.RemoveAll(filepath.Join(folder, ".git")))
			out, code = fitness(t, folder, nil, "commit-attribution")
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "commit-attribution")
		}},
		{"commit-attribution passes a commit with both trailers", func(t *testing.T) {
			t.Parallel()
			out, code := commitAttribution(t, []string{commitAttributionTrailers})
			sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		}},
		{"commit-attribution rejects an adopting repo missing trailers", func(t *testing.T) {
			t.Parallel()
			out, code := commitAttribution(t, []string{commitAttributionTrailers, "feat(api): add paging\n\nJust a body.\n"})
			sees(t, out, code, 1, commitAttributionMissing("AI-Models"), commitAttributionMissing("AI-Tools"))
		}},
		{"commit-attribution judges the proposed message", func(t *testing.T) {
			t.Parallel()
			out, code := commitAttribution(t, []string{commitAttributionTrailers}, "--message=feat(api): add paging\n\nJust a body.\n")
			sees(t, out, code, 1, commitAttributionMissing("AI-Models"), commitAttributionMissing("AI-Tools"))
		}},
		{"commit-attribution rejects an empty proposed message", func(t *testing.T) {
			t.Parallel()
			out, code := commitAttribution(t, []string{commitAttributionTrailers}, "--message=")
			sees(t, out, code, 1, `No commit message to validate; add "AI-Tools:" and "AI-Models:" trailers to disclose AI usage`)
		}},
		{"commit-attribution exempts merges and reverts", func(t *testing.T) {
			t.Parallel()
			for _, subject := range []string{"Merge branch 'feature' into main", `Revert "feat(api): add endpoint"`} {
				out, code := commitAttribution(t, []string{commitAttributionTrailers, subject})
				sees(t, out, code, 0, "All 1 checks passed")
			}
		}},
		{"commit-attribution rejects an empty models trailer", func(t *testing.T) {
			t.Parallel()
			out, code := commitAttribution(t, []string{commitAttributionTrailers, "feat(api): add paging\n\nAI-Tools: Claude Code\nAI-Models:   \n"})
			sees(t, out, code, 1, commitAttributionMissing("AI-Models"))
		}},
		{"commit-attribution rejects a trailer mid sentence", func(t *testing.T) {
			t.Parallel()
			out, code := commitAttribution(t, []string{commitAttributionTrailers, "feat(api): add paging\n\nAI-Tools: Claude Code\nWe set AI-Models: Opus 4.8 later.\n"})
			sees(t, out, code, 1, commitAttributionMissing("AI-Models"))
		}},
	})
}

func Test0003_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"issue-link-once rejects a repeated link by number or URL", func(t *testing.T) {
			t.Parallel()
			repo := iloOnGitHub(t, "topic")
			iloCommit(t, repo, "", "feat(app): start\n\nPlan #12")
			out, code := iloCheck(t, repo, "fix(app): more\n\nSee https://github.com/"+iloRepo+"/issues/12")
			sees(t, out, code, 1, iloRepo+`#12 is already linked by commit`, `("feat(app): start")`)
			repo = iloOnGitHub(t, "topic")
			iloCommit(t, repo, "", "feat(app): begin\n\nPlan may-journal/.github#144")
			out, code = iloCheck(t, repo, "fix(app): more\n\nPlan https://github.com/may-journal/.github/issues/144")
			sees(t, out, code, 1, `may-journal/.github#144 is already linked by commit`, `("feat(app): begin")`)
		}},
		{"issue-link-once passes a link to another repo", func(t *testing.T) {
			t.Parallel()
			repo := iloOnGitHub(t, "topic")
			iloCommit(t, repo, "", "feat(app): start\n\nPlan #12")
			out, code := iloCheck(t, repo, "fix(app): more\n\nSee may-journal/fitness-runner#12")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"issue-link-once ignores comments and scissors lines", func(t *testing.T) {
			t.Parallel()
			repo := iloBranch(t, "topic")
			iloCommit(t, repo, "", "feat(app): start\n\nPlan #12")
			message := "fix(app): more\n# Plan #12\n# ------------------------ >8 ------------------------\nPlan #12\n"
			out, code := iloCheck(t, repo, message)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"issue-link-once rejects a link the open PR holds", func(t *testing.T) {
			t.Parallel()
			repo := iloBranch(t, "topic")
			git(t, repo, "remote", "set-url", "origin", "https://github.com/acme/app.git")
			rp := standIn(t, map[string][]response{"gh": {
				{Match: []string{"pr view", "--json number,state,body"}, Stdout: `{"body":"Adds the export screen.\n\nCloses #12","number":7,"state":"OPEN"}` + "\n"},
			}})
			out, code := fitness(t, repo, rp.env, "issue-link-once", "--message", "fix(app): more\n\nFixes #12")
			sees(t, out, code, 1, "acme/app#12 is already linked by open PR #7")
		}},
		{"issue-link-once passes one early link with main merged", func(t *testing.T) {
			t.Parallel()
			repo := iloBranch(t, "topic")
			iloCommit(t, repo, "2026-09-30T10:00:00Z", "feat(app): start\n\nPlan #12")
			git(t, repo, "switch", "-q", "main")
			iloCommit(t, repo, "2026-09-30T13:00:00Z", "feat(app): more on main\n\nPlan #12")
			git(t, repo, "push", "-q", "origin", "main")
			git(t, repo, "switch", "-q", "topic")
			user(t, repo, nil, "merge", "-q", "--no-edit", "main")
			out, code := iloPullRequest(t, repo)
			sees(t, out, code, 0, "All 1 checks passed · 2 files scanned")
		}},
		{"issue-link-once rejects a second commit repeating a link", func(t *testing.T) {
			t.Parallel()
			repo := iloBranch(t, "topic")
			iloCommit(t, repo, "2026-09-30T10:00:00Z", "feat(app): start\n\nPlan #12")
			iloCommit(t, repo, "2026-09-30T11:00:00Z", "fix(app): more\n\nPlan #12")
			out, code := iloPullRequest(t, repo)
			sees(t, out, code, 1, "commit "+iloHead(t, repo)[:7]+` ("fix(app): more"): #12 is already linked by commit`, "--force-with-lease")
		}},
		{"issue-link-once rejects a later commit repeating the PR link", func(t *testing.T) {
			t.Parallel()
			repo := iloBranch(t, "topic")
			iloCommit(t, repo, "2026-09-30T13:00:00Z", "fix(app): more\n\nFixes #12")
			out, code := iloPullRequest(t, repo)
			sees(t, out, code, 1, "#12 is already linked by open PR #9", "--force-with-lease")
		}},
		{"issue-link-once passes with no message or PR", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", nil), nil, "issue-link-once")
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "issue-link-once")
		}},
	})
}

func Test0003_8(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"release-changelog passes a repo with no manifest", func(t *testing.T) {
			t.Parallel()
			out, code := releaseChangelog(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "release-changelog")
		}},
		{"release-changelog passes an empty bootstrap manifest", func(t *testing.T) {
			t.Parallel()
			out, code := releaseChangelog(t, map[string]string{".release-please-manifest.json": "{}\n"})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"release-changelog rejects a manifest without a root version", func(t *testing.T) {
			t.Parallel()
			releaseChangelogFails(t, []releaseChangelogCase{
				{map[string]string{".release-please-manifest.json": "{\n"}, ".release-please-manifest.json is invalid JSON"},
				{map[string]string{".release-please-manifest.json": "null\n"}, ".release-please-manifest.json must be a JSON object"},
				{map[string]string{".release-please-manifest.json": `{"go":"1.2.3"}`}, `.release-please-manifest.json must contain a root "." release version`},
				{map[string]string{".release-please-manifest.json": `{".":"v1.2.3"}`}, `.release-please-manifest.json root version must be numeric major.minor.patch (got "v1.2.3")`},
			})
		}},
		{"release-changelog rejects a version file that differs", func(t *testing.T) {
			t.Parallel()
			files := releaseChangelogRelease(upstreamChangelog)
			files["version.txt"] = "1.2.4\n"
			out, code := releaseChangelog(t, files)
			sees(t, out, code, 1, `version.txt version "1.2.4" must match .release-please-manifest.json version "1.2.3"`)
		}},
		{"release-changelog rejects a release with no version file", func(t *testing.T) {
			t.Parallel()
			out, code := releaseChangelog(t, map[string]string{".release-please-manifest.json": `{".":"1.2.3"}` + "\n"})
			sees(t, out, code, 1, "version.txt is missing; it must match release 1.2.3")
		}},
		{"release-changelog passes a matching version and bare changelog", func(t *testing.T) {
			t.Parallel()
			for _, changelog := range []string{upstreamChangelog, ""} {
				repo := example(t, "happyRepo", releaseChangelogRelease(changelog))
				write(t, repo, map[string]string{".release-please-manifest.json": `{".": "1.2.3"}` + "\n"})
				git(t, repo, "add", ".release-please-manifest.json")
				out, code := fitness(t, repo, nil, "release-changelog")
				sees(t, out, code, 0, "All 1 checks passed")
			}
		}},
	})
}

// commitAttributionTrailers is a commit message carrying both AI trailers.
const commitAttributionTrailers = "feat(api): add endpoint\n\nAI-Tools: Claude Code\nAI-Models: Opus 4.8\n"

// commitAttributionMissing is the error for a missing trailer key.
func commitAttributionMissing(key string) string {
	return `commit message missing "` + key + `:" trailer`
}

// commitAttribution commits each message on happyRepo in order, the last at
// HEAD, then runs `fitness-install -- commit-attribution <args>`.
func commitAttribution(t *testing.T, messages []string, args ...string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", nil)
	for _, m := range messages {
		commit(t, repo, nil, m)
	}
	return fitness(t, repo, nil, append([]string{"commit-attribution"}, args...)...)
}

// iloRepo is the repo on GitHub that origin points at, so a bare `#NN`
// belongs to it.
const iloRepo = "may-journal/example"

// iloOrigin is iloRepo's URL.
const iloOrigin = "https://github.com/" + iloRepo + ".git"

// iloBranch is happyRepo on branch, off a main pushed to a local origin.
func iloBranch(t *testing.T, branch string) string {
	t.Helper()
	repo := example(t, "happyRepo", nil)
	remote(t, repo)
	git(t, repo, "push", "-q", "origin", "HEAD:main")
	git(t, repo, "switch", "-q", "-c", branch)
	return repo
}

// iloOnGitHub is iloBranch with origin pointed at iloRepo on GitHub.
func iloOnGitHub(t *testing.T, branch string) string {
	t.Helper()
	repo := iloBranch(t, branch)
	git(t, repo, "remote", "set-url", "origin", iloOrigin)
	return repo
}

// iloCommit commits message in repo, authored at date when it is set.
func iloCommit(t *testing.T, repo, date, message string) {
	t.Helper()
	if out, code := user(t, repo, []string{"GIT_AUTHOR_DATE=" + date}, "commit", "-q", "--allow-empty", "-m", message); code != 0 {
		t.Fatalf("commit %q: %s", message, out)
	}
}

// iloCheck runs issue-link-once on the proposed commit message.
func iloCheck(t *testing.T, repo, message string) (string, int) {
	t.Helper()
	return fitness(t, repo, nil, "issue-link-once", "--message", message)
}

// iloHead returns repo's HEAD commit.
func iloHead(t *testing.T, repo string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", repo, "rev-parse", "HEAD").Output()
	mustDo(t, err)
	return strings.TrimSpace(string(out))
}

// iloPullRequest runs issue-link-once as CI does on pull request #9, opened
// at noon with a body closing #12, whose head is repo's HEAD.
func iloPullRequest(t *testing.T, repo string) (string, int) {
	t.Helper()
	payload := map[string]any{"pull_request": map[string]any{
		"number": 9, "body": "Closes #12", "created_at": "2026-09-30T12:00:00Z",
		"head": map[string]any{"sha": iloHead(t, repo)},
	}}
	data, err := json.Marshal(payload)
	mustDo(t, err)
	path := filepath.Join(t.TempDir(), "event.json")
	mustDo(t, os.WriteFile(path, data, 0o644))
	env := []string{"GITHUB_BASE_REF=main", "GITHUB_EVENT_PATH=" + path}
	return fitness(t, repo, env, "issue-link-once")
}

// planTrailer runs `fitness-install -- plan-trailer --message <msg>` on
// happyRepo, as the commit-msg hook passes a proposed message.
func planTrailer(t *testing.T, msg string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", nil), nil, "plan-trailer", "--message", msg)
}

// planTrailerFix is the hint every misspelled plan reference carries.
const planTrailerFix = `must read exactly "Plan #<number>", "Plan <owner>/<repo>#<number>" or "Plan https://github.com/<owner>/<repo>/issues/<number>"`

// releaseChangelog runs `fitness-install -- release-changelog` on happyRepo
// with files written over it.
func releaseChangelog(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "release-changelog")
}

// releaseChangelogRelease is a repo releasing 1.2.3 with changelog as its
// CHANGELOG.md, or none when it is empty.
func releaseChangelogRelease(changelog string) map[string]string {
	return map[string]string{
		".release-please-manifest.json": `{".":"1.2.3"}` + "\n",
		"version.txt":                   "1.2.3\n",
		"CHANGELOG.md":                  changelog,
	}
}

// releaseChangelogCase is one Given a test runs, with the text it must show.
type releaseChangelogCase struct {
	files map[string]string
	want  string
}

// releaseChangelogFails runs each case and asserts it fails showing its text.
func releaseChangelogFails(t *testing.T, cases []releaseChangelogCase) {
	t.Helper()
	for _, c := range cases {
		out, code := releaseChangelog(t, c.files)
		sees(t, out, code, 1, c.want)
	}
}

// semanticCommit runs `fitness-install -- semantic-commit <args>` on
// happyRepo, whose HEAD subject is "docs(readme): add readme".
func semanticCommit(t *testing.T, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", nil), nil, append([]string{"semantic-commit"}, args...)...)
}

// semanticCommitTypes is the type list the check prints in every failure.
const semanticCommitTypes = "(types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert)"

// inherited is an upstream subject with no scope, as a fork's main can hold.
const inherited = "test: migrate from jest to vite (#75)"

// defaultNames are fork default branches with origin/HEAD unset: the usual
// two, and one Fitness cannot guess.
var defaultNames = []string{"main", "master", "develop"}

// forkedRepo is happyRepo on a fresh branch over inherited history that
// origin/<base> holds.
func forkedRepo(t *testing.T, base string) string {
	t.Helper()
	repo := example(t, "happyRepo", nil)
	fork(t, repo, base)
	return repo
}

// fork commits an inherited bad subject to repo without hooks, as a fetch
// does, records it as origin/<base>, and branches fresh.
func fork(t *testing.T, repo, base string) {
	t.Helper()
	out, err := exec.Command("git", "-C", repo, "-c", "user.name=t", "-c", "user.email=t@t", "commit-tree", "HEAD^{tree}", "-p", "HEAD", "-m", inherited).Output()
	mustDo(t, err)
	git(t, repo, "reset", "-q", "--hard", strings.TrimSpace(string(out)))
	remote(t, repo)
	git(t, repo, "update-ref", "refs/remotes/origin/"+base, "HEAD")
	git(t, repo, "checkout", "-q", "-b", "work")
}

// hookedFork is a happyRepo fork with the shared hooks, or with only the
// installer's pre-commit hook when installerHook is set.
func hookedFork(t *testing.T, installerHook bool) (string, []string) {
	t.Helper()
	if !installerHook {
		repo, env := hooked(t, nil)
		fork(t, repo, "main")
		return repo, env
	}
	repo := forkedRepo(t, "main")
	cmd := exec.Command(installer(t), "--install-hook")
	cmd.Dir, cmd.Env = repo, userEnv()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("--install-hook: %v\n%s", err, out)
	}
	return repo, nil
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
