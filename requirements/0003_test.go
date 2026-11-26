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
		{"semantic-commit passes only a listed type and scope, so a subject with no scope fails", func(t *testing.T) {
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

func Test0003_5(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"changelog rejects a file name that lacks the date", func(t *testing.T) {
			t.Parallel()
			out, code := changelogRun(t, withEntries(map[string]string{
				"undated-change":               entryFileText("2026.10.11.0900", "Feat: undated change."),
				"2026-10-11-0800-moved-change": entryFileText("2026.10.11.1000", "Fix: moved change."),
			}))
			sees(t, out, code, 1,
				"changelog/undated-change.md must be named changelog/2026-10-11-0900-undated-change.md: its heading's date, then the change",
				"changelog/2026-10-11-0800-moved-change.md must be named changelog/2026-10-11-1000-moved-change.md")
		}},
		{"changelog rejects two dated or undated headings", func(t *testing.T) {
			t.Parallel()
			out, code := changelogRun(t, withEntries(map[string]string{
				"2026-10-11-0900-good-change":  entryFileText("2026.10.11.0900", "Feat: good change."),
				"2026-10-11-1000-two-headings": entryFileText("2026.10.11.1000", "Fix: one.") + "\n" + entryFileText("2026.10.11.1100", "Fix: two."),
				"2026-10-11-1200-undated":      "### Added\n\n- Feat: one.\n- Fix: two.\n- Docs: three.\n",
			}))
			sees(t, out, code, 1,
				"changelog/2026-10-11-1000-two-headings.md must hold one ### yyyy.mm.dd.HHMM heading and its bullets, nothing else",
				"changelog/2026-10-11-1200-undated.md must hold one ### yyyy.mm.dd.HHMM heading")
		}},
		{"changelog judges only entry files and passes upstream history", func(t *testing.T) {
			t.Parallel()
			files := withEntries(map[string]string{
				"2026-10-10-0900-older-change": entryFileText("2026.10.10.0900", "Fix: older change."),
				"2026-10-11-0900-newer-change": bareEntryText("2026.10.11.0900", "Feat: newer change."),
			})
			files["CHANGELOG.md"] = upstreamChangelog
			files["package.json"] = `{"name": "app", "private": true, "version": "2.0.1"}` + "\n"
			files["package-lock.json"] = `{"name": "app", "version": "2.0.0"}` + "\n"
			out, code := changelogRun(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "changelog", "2")
		}},
		{"changelog passes a repo with no entry files", func(t *testing.T) {
			t.Parallel()
			for _, changelog := range []string{upstreamChangelog, ""} {
				out, code := changelogRun(t, map[string]string{"CHANGELOG.md": changelog})
				sees(t, out, code, 0, "All 1 checks passed")
				passedWithNoFiles(t, out, "changelog")
			}
		}},
		{"changelog-updated rejects a new entry name without its date", func(t *testing.T) {
			t.Parallel()
			files := withEntries(map[string]string{"runner-feature": entryFileText(changelogStamp, "Feat: added runner feature for validation.")})
			files["src.txt"] = changelogFeature
			out, code := changelogStaged(t, files)
			sees(t, out, code, 1, "changelog/runner-feature.md must be named changelog/2026-10-12-0900-runner-feature.md: its heading's date, then the change")
		}},
	})
}

func Test0003_6(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"changelog-bullets rejects a bullet over the length limit", func(t *testing.T) {
			t.Parallel()
			long := "Feat: " + strings.Repeat("x", 359)
			wrapped := "Fix: " + strings.Repeat("x", 200) + "\n  " + strings.Repeat("y", 200)
			out, code := changelogBullets(t, changelogBulletsDoc(long, wrapped, "Docs: three."))
			sees(t, out, code, 1,
				bulletsFile+":3: bullet is 365 characters; keep each under 365",
				bulletsFile+":4: bullet is 406 characters; keep each under 365")
		}},
		{"changelog-bullets rejects too few or too many bullets", func(t *testing.T) {
			t.Parallel()
			out, code := changelogBullets(t, changelogBulletsDoc("Docs: add the readme.", "Docs: add the changelog."))
			sees(t, out, code, 1, bulletsFile+` section "2026.10.08.1400" has 2 bullets; keep entries to 3-5 bullets`)
			out, code = changelogBullets(t, changelogBulletsDoc("Feat: one.", "Fix: two.", "Fix: three.", "Test: four.", "Docs: five.", "Docs: six."))
			sees(t, out, code, 1, `"2026.10.08.1400" has 6 bullets`)
			out, code = changelogBullets(t, "### 2026.10.08.1400\n\n* Fix: one.\n* Fix: two.\n* Fix: three.\n")
			sees(t, out, code, 1, `"2026.10.08.1400" has 0 bullets`)
		}},
		{"changelog-bullets rejects a file with two bullets", func(t *testing.T) {
			t.Parallel()
			files := withEntries(map[string]string{
				"short-entry":   "### 2026.10.11.0900\n\n- Feat: one.\n- Fix: two.\n",
				"untyped-entry": "### 2026.10.11.1000\n\n- Feat: one.\n- Fix: two.\n- no type here.\n",
			})
			out, code := fitness(t, example(t, "happyRepo", files), nil, "changelog-bullets")
			sees(t, out, code, 1, `changelog/short-entry.md section "2026.10.11.0900" has 2 bullets`,
				"changelog/untyped-entry.md:5: bullet must start with a semantic type")
		}},
		{"changelog-bullets passes upstream lists and no entry files", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"CHANGELOG.md": upstreamChangelog})
			out, code := fitness(t, repo, nil, "changelog-bullets")
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "changelog-bullets")
		}},
		{"changelog-bullets fails only the untyped bullet", func(t *testing.T) {
			t.Parallel()
			out, code := changelogBullets(t, changelogBulletsDoc("Feat: one.", "fix(release)!: two.", "chore(ci): three.", "Revert: four.", "Added the readme."))
			sees(t, out, code, 1, bulletsFile+":7: bullet must start with a semantic type")
			if n := strings.Count(read(out), "must start with a semantic type"); n != 1 {
				t.Errorf("only the untyped bullet must fail, got %d failures:\n%s", n, out)
			}
		}},
	})
}

func Test0003_7(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"changelog-updated rejects an entry sharing few words", func(t *testing.T) {
			t.Parallel()
			out, code := changelogStaged(t, changelogStagedEntry(changelogStamp, "Fix: minor tweak."))
			sees(t, out, code, 1,
				"Changelog entry file additions should mention at least 3 words from your staged changes (found 0: )",
				"e.g. use words like: added, runner, feature, for, validation")
		}},
		{"changelog-updated passes a merge replaying an old entry", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"src.txt": "base\n"})
			git(t, repo, "checkout", "-q", "-b", "side")
			commit(t, repo, map[string]string{
				"src.txt":      changelogFeature,
				"CHANGELOG.md": changelogEntry("2026.07.07.0850", "Feat: added runner feature for validation."),
			}, "feat: add runner")
			git(t, repo, "checkout", "-q", "-")
			commit(t, repo, map[string]string{"README.md": readme + "\nMore.\n"}, "docs: more")
			git(t, repo, "-c", "user.name=t", "-c", "user.email=t@t", "merge", "-q", "--no-commit", "--no-ff", "side")
			out, code := fitness(t, repo, nil, "changelog-updated")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"changelog-updated rejects a second entry file on a branch", func(t *testing.T) {
			t.Parallel()
			first := entryFileText(changelogFirst, "Feat: added runner feature for validation.")
			repo := changelogBranch(t, nil, map[string]string{"changelog/2026-01-02-0900-runner-feature.md": first})
			second := withEntries(map[string]string{"2026-10-12-0900-runner-option": entryFileText(changelogStamp, "Feat: added runner option for validation.")})
			out, code := changelogSecond(t, repo, second)
			sees(t, out, code, 1, "this branch already has its entry changelog/2026-01-02-0900-runner-feature.md",
				"instead of a new entry: changelog/2026-10-12-0900-runner-option.md")
		}},
		{"changelog-updated rejects changes with no entry file", func(t *testing.T) {
			t.Parallel()
			out, code := changelogStaged(t, map[string]string{"src.txt": changelogFeature})
			sees(t, out, code, 1, "Stage a changelog/<yyyy-mm-dd-HHMM>-<change>.md entry file that mentions your staged changes")
		}},
		{"changelog-updated passes a new entry whatever the changelog", func(t *testing.T) {
			t.Parallel()
			for _, changelog := range []string{
				changelogEntry(changelogStamp, "Feat: added runner feature for validation."),
				strings.Replace(changelogBase, "- Docs: add the readme.", "- Docs: add the readme file.", 1),
				"",
			} {
				files := changelogStagedEntry(changelogStamp, "Feat: added runner feature for validation.")
				files["CHANGELOG.md"] = changelog
				out, code := changelogStaged(t, files)
				sees(t, out, code, 0, "All 1 checks passed")
			}
			repo := example(t, "happyRepo", map[string]string{"CHANGELOG.md": ""})
			write(t, repo, changelogStagedEntry(changelogStamp, "Feat: added runner feature for validation."))
			git(t, repo, "add", "-A")
			out, code := fitness(t, repo, nil, "changelog-updated")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"changelog-updated rejects edits to a released entry", func(t *testing.T) {
			t.Parallel()
			const released = "changelog/2026-10-01-0900-runner-feature.md"
			old := entryFileText("2026.10.01.0900", "Feat: added runner feature for validation.")
			repo := example(t, "happyRepo", map[string]string{released: old})
			git(t, repo, "tag", "v1.0.0")
			write(t, repo, map[string]string{
				released:  old + "- Feat: added runner option for validation.\n",
				"src.txt": changelogFeature + "Added runner option for validation.\n",
			})
			git(t, repo, "add", "-A")
			out, code := fitness(t, repo, nil, "changelog-updated")
			sees(t, out, code, 1, released+" shipped in v1.0.0, so its notes are published; leave it as it is and add a new changelog/<yyyy-mm-dd-HHMM>-<change>.md entry file instead")
		}},
		{"changelog-updated passes dropping front matter from a released entry", func(t *testing.T) {
			t.Parallel()
			const released = "changelog/2026-10-01-0900-old-feature.md"
			item := "Feat: added old feature."
			repo := example(t, "happyRepo", map[string]string{released: entryFileText("2026.10.01.0900", item)})
			git(t, repo, "tag", "v1.0.0")
			files := changelogStagedEntry(changelogStamp, "Feat: added runner feature for validation.")
			files[released] = bareEntryText("2026.10.01.0900", item)
			write(t, repo, files)
			git(t, repo, "add", "-A")
			out, code := fitness(t, repo, nil, "changelog-updated")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
	})
}

// changelogRun runs `fitness-install -- changelog` on happyRepo with files
// written over it.
func changelogRun(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "changelog")
}

// entryFileText is an entry file under stamp with three typed bullets, the
// first one item, below the front matter older entries carry.
func entryFileText(stamp, item string) string {
	return "---\nrelatedConfigurations: ['../.fitnessrc.json']\n---\n\n" + bareEntryText(stamp, item)
}

// bareEntryText is entryFileText without front matter, which is optional.
func bareEntryText(stamp, item string) string {
	return "### " + stamp + "\n\n- " + item + "\n- Test: cover it.\n- Docs: describe it.\n"
}

// withEntries is files over happyRepo holding each named entry file.
func withEntries(entries map[string]string) map[string]string {
	files := map[string]string{}
	for name, body := range entries {
		files["changelog/"+name+".md"] = body
	}
	return files
}

// upstreamChangelog is a CHANGELOG.md as an upstream project keeps it:
// undated, typed sections of numbered items under a version heading.
const upstreamChangelog = "# Changelog\n\n## v2.0.0\n\n### Added\n\n1. Catalog-level publishing\n2. Plugin hooks\n\n### Removed\n1. Old lint tools\n"

// changelogBullets runs `fitness-install -- changelog-bullets` on happyRepo
// holding one entry file, body, at changelog/2026-10-08-1400-entry.md.
func changelogBullets(t *testing.T, body string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", withEntries(map[string]string{"2026-10-08-1400-entry": body}))
	return fitness(t, repo, nil, "changelog-bullets")
}

// bulletsFile is the path changelogBullets writes its entry file to.
const bulletsFile = "changelog/2026-10-08-1400-entry.md"

// changelogBulletsDoc is an entry file whose section, opened on line 1,
// holds bullets, so its first bullet sits on line 3.
func changelogBulletsDoc(bullets ...string) string {
	doc := "### 2026.10.08.1400\n\n"
	for _, b := range bullets {
		doc += "- " + b + "\n"
	}
	return doc
}

// changelogStaged commits happyRepo, stages staged over it, and runs
// `fitness-install -- changelog-updated` as a pre-commit hook would.
func changelogStaged(t *testing.T, staged map[string]string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", nil)
	write(t, repo, staged)
	git(t, repo, "add", "-A")
	return fitness(t, repo, nil, "changelog-updated")
}

// changelogBase is happyRepo's changelog.
const changelogBase = "---\nrelatedConfigurations: ['.fitnessrc.json']\n---\n\n# Changelog\n\n## Changes\n\n" +
	"### 2026.10.08.1400\n\n- Docs: add the readme.\n- Docs: add the changelog.\n- Docs: add the license note.\n"

// changelogEntry is happyRepo's changelog with a new section, holding one
// bullet item, on top.
func changelogEntry(stamp, item string) string {
	return strings.Replace(changelogBase, "## Changes\n\n", "## Changes\n\n### "+stamp+"\n\n- "+item+"\n\n", 1)
}

// changelogStamp is the heading of a new entry.
const changelogStamp = "2026.10.12.0900"

// changelogFeature is a staged source change the entry should mention.
const changelogFeature = "Added runner feature for validation.\n"

// changelogStagedEntry is happyRepo's staged feature with its entry file
// under stamp, its first bullet item.
func changelogStagedEntry(stamp, item string) map[string]string {
	files := withEntries(map[string]string{"2026-10-12-0900-runner-feature": entryFileText(stamp, item)})
	files["src.txt"] = changelogFeature
	return files
}

// changelogBranch is happyRepo, with files over it, on a branch off
// origin/main whose first commit added first for the staged feature.
func changelogBranch(t *testing.T, files, first map[string]string) string {
	t.Helper()
	repo := example(t, "happyRepo", files)
	git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	git(t, repo, "checkout", "-q", "-b", "work")
	first["src.txt"] = changelogFeature
	commit(t, repo, first, "feat: add runner")
	return repo
}

// changelogFirst is the stamp of the branch's first entry.
const changelogFirst = "2026.01.02.0900"

// changelogSecond stages a second change in repo, with entry, and runs
// changelog-updated on it.
func changelogSecond(t *testing.T, repo string, entry map[string]string) (string, int) {
	t.Helper()
	entry["src.txt"] = changelogFeature + "Added runner option for validation.\n"
	write(t, repo, entry)
	git(t, repo, "add", "-A")
	return fitness(t, repo, nil, "changelog-updated")
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
