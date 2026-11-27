// Package requirements proves the changelog-updated check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo, staging a commit as a pre-commit
// hook would see it.
package requirements

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of every repo the tests make.
func TestMain(m *testing.M) { acceptance.Main(m) }

// The shared helpers from acceptance, named as the root requirement tests
// name them.
var (
	entryFileText     = acceptance.EntryFileText
	bareEntryText     = acceptance.BareEntryText
	withEntries       = acceptance.WithEntries
	example           = acceptance.Example
	commit            = acceptance.Commit
	write             = acceptance.Write
	git               = acceptance.Git
	fitness           = acceptance.Fitness
	read              = acceptance.Read
	sees              = acceptance.Sees
	passedWithNoFiles = acceptance.PassedWithNoFiles
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a merge, cherry-pick, revert or nothing staged skips", func(t *testing.T) {
		t.Parallel()
		t.Run("a merge, cherry-pick or revert in progress skips", func(t *testing.T) {
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
			// A conflicting cherry-pick stops with nothing about its change staged.
			git(t, repo, "merge", "--abort")
			commit(t, repo, map[string]string{"src.txt": "other\n"}, "fix: other")
			if err := exec.Command("git", "-C", repo, "cherry-pick", "side").Run(); err == nil {
				t.Fatal("cherry-pick must stop on its conflict")
			}
			out, code = fitness(t, repo, nil, "changelog-updated")
			sees(t, out, code, 0, "All 1 checks passed")
			// Reverting a commit that added an entry stages its deletion.
			git(t, repo, "cherry-pick", "--abort")
			commit(t, repo, changelogStagedEntry(changelogStamp, "Feat: added runner feature for validation."), "feat: add runner")
			git(t, repo, "revert", "--no-commit", "HEAD")
			out, code = fitness(t, repo, nil, "changelog-updated")
			sees(t, out, code, 0, "All 1 checks passed")
		})
		t.Run("nothing staged skips", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			write(t, repo, map[string]string{"src.txt": changelogFeature})
			out, code := fitness(t, repo, nil, "changelog-updated")
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "changelog-updated")
		})
	})
	t.Run("0001.1.2 changes with no entry file fail", func(t *testing.T) {
		t.Parallel()
		out, code := changelogStaged(t, map[string]string{"src.txt": changelogFeature})
		sees(t, out, code, 1, "Stage a changelog/<yyyy-mm-dd-HHMM>-<change>.md entry file that mentions your staged changes")
	})
	t.Run("0001.1.3 a new entry passes whatever the changelog", func(t *testing.T) {
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
	})
	t.Run("0001.1.4 an entry sharing few words fails with the count", func(t *testing.T) {
		t.Parallel()
		out, code := changelogStaged(t, changelogStagedEntry(changelogStamp, "Fix: minor tweak."))
		sees(t, out, code, 1,
			"Changelog entry file additions should mention at least 3 words from your staged changes (found 0: )",
			"e.g. use words like: added, runner, feature, for, validation")
		// The suggestion stops at ten words, and CHANGELOG.md counts as the diff.
		repo := example(t, "happyRepo", map[string]string{"CHANGELOG.md": "# Changelog\n"})
		files := withEntries(map[string]string{"2026-10-12-0900-runner-feature": bareEntryText(changelogStamp, "Fix: minor tweak.")})
		files["CHANGELOG.md"] = "# Changelog\n\nalpha bravo charlie delta echo foxtrot golf hotel india juliet kilo lima\n"
		write(t, repo, files)
		git(t, repo, "add", "-A")
		out, code = fitness(t, repo, nil, "changelog-updated")
		sees(t, out, code, 1, "(found 0: )", "e.g. use words like: alpha, bravo, charlie, delta, echo, foxtrot, golf, hotel, india, juliet")
		if strings.Contains(read(out), "kilo") {
			t.Errorf("the suggestion must stop at ten words:\n%s", out)
		}
	})
	t.Run("0001.1.5 only added words of three or more letters count", func(t *testing.T) {
		t.Parallel()
		repo := example(t, "happyRepo", map[string]string{"src.txt": "Runner feature validation.\n"})
		files := withEntries(map[string]string{
			"2026-10-12-0900-runner-feature": bareEntryText(changelogStamp, "Feat: RUNNER FEATURE VALIDATION ADDED UI DB."),
		})
		files["src.txt"] = "Runner feature validation.\nAdded UI to DB.\n"
		write(t, repo, files)
		git(t, repo, "add", "-A")
		out, code := fitness(t, repo, nil, "changelog-updated")
		sees(t, out, code, 1, "(found 1: added)", "e.g. use words like: added")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a new entry name without its date fails", func(t *testing.T) {
		t.Parallel()
		files := withEntries(map[string]string{"runner-feature": entryFileText(changelogStamp, "Feat: added runner feature for validation.")})
		files["src.txt"] = changelogFeature
		out, code := changelogStaged(t, files)
		sees(t, out, code, 1, "changelog/runner-feature.md must be named changelog/2026-10-12-0900-runner-feature.md: its heading's date, then the change")
	})
	t.Run("0001.2.2 a commit adding another entry after an earlier one on the branch fails", func(t *testing.T) {
		t.Parallel()
		repo := changelogBranch(t, map[string]string{"changelog/2026-01-02-0900-runner-feature.md": changelogFirstEntry})
		out, code := changelogSecond(t, repo, changelogOption())
		sees(t, out, code, 1, "this branch already has its entry changelog/2026-01-02-0900-runner-feature.md",
			"instead of a new entry: changelog/2026-10-12-0900-runner-option.md")
	})
	t.Run("0001.2.3 a stacked branch counts only its own commits", func(t *testing.T) {
		t.Parallel()
		repo := changelogBranch(t, map[string]string{"changelog/2026-01-02-0900-runner-feature.md": changelogFirstEntry})
		git(t, repo, "update-ref", "refs/remotes/origin/work", "HEAD")
		git(t, repo, "checkout", "-q", "-b", "stacked")
		out, code := changelogSecond(t, repo, changelogOption())
		sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.4 the default branch or a repo without origin judges each commit alone", func(t *testing.T) {
		t.Parallel()
		first := map[string]string{"changelog/2026-01-02-0900-runner-feature.md": changelogFirstEntry, "src.txt": changelogFeature}
		repo := example(t, "happyRepo", nil)
		git(t, repo, "checkout", "-q", "-B", "main")
		git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
		commit(t, repo, first, "feat: add runner")
		out, code := changelogSecond(t, repo, changelogOption())
		sees(t, out, code, 0, "All 1 checks passed")
		repo = example(t, "happyRepo", nil)
		git(t, repo, "checkout", "-q", "-b", "work")
		commit(t, repo, first, "feat: add runner")
		out, code = changelogSecond(t, repo, changelogOption())
		sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.5 editing or deleting a released entry fails", func(t *testing.T) {
		t.Parallel()
		const released = "changelog/2026-10-01-0900-runner-feature.md"
		const want = released + " shipped in v1.0.0, so its notes are published; leave it as it is and add a new changelog/<yyyy-mm-dd-HHMM>-<change>.md entry file instead"
		old := entryFileText("2026.10.01.0900", "Feat: added runner feature for validation.")
		for _, edit := range []string{old + "- Feat: added runner option for validation.\n", ""} {
			repo := example(t, "happyRepo", map[string]string{released: old})
			git(t, repo, "tag", "v1.0.0")
			write(t, repo, map[string]string{
				released:  edit,
				"src.txt": changelogFeature + "Added runner option for validation.\n",
			})
			git(t, repo, "add", "-A")
			out, code := fitness(t, repo, nil, "changelog-updated")
			sees(t, out, code, 1, want)
		}
	})
	t.Run("0001.2.6 dropping front matter from a released entry passes", func(t *testing.T) {
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
	})
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

// readme is happyRepo's README.md, copied from the root requirement tests.
const readme = "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n"

// changelogBase is a changelog with one dated section of three bullets.
const changelogBase = "---\nrelatedConfigurations: ['.fitnessrc.json']\n---\n\n# Changelog\n\n## Changes\n\n" +
	"### 2026.10.08.1400\n\n- Docs: add the readme.\n- Docs: add the changelog.\n- Docs: add the license note.\n"

// changelogEntry is changelogBase with a new section, holding one bullet
// item, on top.
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

// changelogBranch is happyRepo on a branch off
// origin/main whose first commit added first for the staged feature.
func changelogBranch(t *testing.T, first map[string]string) string {
	t.Helper()
	repo := example(t, "happyRepo", nil)
	git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	git(t, repo, "checkout", "-q", "-b", "work")
	first["src.txt"] = changelogFeature
	commit(t, repo, first, "feat: add runner")
	return repo
}

// changelogFirstEntry is the entry file a branch's first commit adds.
var changelogFirstEntry = entryFileText("2026.01.02.0900", "Feat: added runner feature for validation.")

// changelogOption is a second commit's own new entry file.
func changelogOption() map[string]string {
	return withEntries(map[string]string{"2026-10-12-0900-runner-option": entryFileText(changelogStamp, "Feat: added runner option for validation.")})
}

// changelogSecond stages a second change in repo, with entry, and runs
// changelog-updated on it.
func changelogSecond(t *testing.T, repo string, entry map[string]string) (string, int) {
	t.Helper()
	entry["src.txt"] = changelogFeature + "Added runner option for validation.\n"
	write(t, repo, entry)
	git(t, repo, "add", "-A")
	return fitness(t, repo, nil, "changelog-updated")
}
