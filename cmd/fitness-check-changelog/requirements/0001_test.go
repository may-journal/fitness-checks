// Package requirements proves the changelog check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of every repo the tests make.
func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

// The shared helpers from acceptance, named as the root requirement tests
// name them.
var (
	entryFileText     = acceptance.EntryFileText
	bareEntryText     = acceptance.BareEntryText
	withEntries       = acceptance.WithEntries
	upstreamChangelog = acceptance.UpstreamChangelog
	example           = acceptance.Example
	fitness           = acceptance.Fitness
	sees              = acceptance.Sees
	passedWith        = acceptance.PassedWith
	passedWithNoFiles = acceptance.PassedWithNoFiles
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 dated, well-named entry files pass, counting them", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, withEntries(map[string]string{
			"2026-10-10-0900-older-change": entryFileText("2026.10.10.0900", "Fix: older change."),
			"2026-10-11-0900-newer-change": bareEntryText("2026.10.11.0900", "Feat: newer change."),
			"2026-10-12-0900-third-change": bareEntryText("2026.10.12.0900", "Docs: third change."),
		}))
		sees(t, out, code, 0, "All 1 checks passed")
		passedWith(t, out, "changelog", "3")
	})
	t.Run("0001.1.2 two dated or undated headings fail", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, withEntries(map[string]string{
			"2026-10-11-0900-good-change":  entryFileText("2026.10.11.0900", "Feat: good change."),
			"2026-10-11-1000-two-headings": entryFileText("2026.10.11.1000", "Fix: one.") + "\n" + entryFileText("2026.10.11.1100", "Fix: two."),
			"2026-10-11-1200-undated":      "### Added\n\n- Feat: one.\n- Fix: two.\n- Docs: three.\n",
		}))
		sees(t, out, code, 1,
			"changelog/2026-10-11-1000-two-headings.md must hold one ### yyyy.mm.dd.HHMM heading and its bullets, nothing else",
			"changelog/2026-10-11-1200-undated.md must hold one ### yyyy.mm.dd.HHMM heading")
	})
	t.Run("0001.1.3 a name that lacks its heading's date fails", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, withEntries(map[string]string{
			"undated-change":               entryFileText("2026.10.11.0900", "Feat: undated change."),
			"2026-10-11-0800-moved-change": entryFileText("2026.10.11.1000", "Fix: moved change."),
		}))
		sees(t, out, code, 1,
			"changelog/undated-change.md must be named changelog/2026-10-11-0900-undated-change.md: its heading's date, then the change",
			"changelog/2026-10-11-0800-moved-change.md must be named changelog/2026-10-11-1000-moved-change.md")
	})
	t.Run("0001.1.4 a repo with no entry files passes", func(t *testing.T) {
		t.Parallel()
		for _, changelog := range []string{upstreamChangelog, ""} {
			out, code := changelogRun(t, map[string]string{"CHANGELOG.md": changelog})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "changelog")
		}
	})
	t.Run("0001.1.5 only entry files are judged, so upstream history passes", func(t *testing.T) {
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
	})
}

// changelogRun runs `fitness-install -- changelog` on happyRepo with files
// written over it.
func changelogRun(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "changelog")
}
