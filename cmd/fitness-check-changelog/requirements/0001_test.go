// Package requirements proves the changelog check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
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
	upstreamChangelog = acceptance.UpstreamChangelog
	example           = acceptance.Example
	fitness           = acceptance.Fitness
	sees              = acceptance.Sees
	passedWith        = acceptance.PassedWith
	passedWithNoFiles = acceptance.PassedWithNoFiles
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.3.1.1 it passes counting the entry files", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, withEntries(map[string]string{
			"2026-10-10-0900-older-change": entryFileText("2026.10.10.0900", "Fix: older change."),
			"2026-10-11-0900-newer-change": bareEntryText("2026.10.11.0900", "Feat: newer change."),
			"2026-10-12-0900-third-change": bareEntryText("2026.10.12.0900", "Docs: third change."),
		}))
		sees(t, out, code, 0, "All 1 checks passed")
		passedWith(t, out, "changelog", "3")
	})
	t.Run("0001.1.3.1.2 front matter is optional", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, withEntries(map[string]string{
			"2026-10-10-0900-older-change": entryFileText("2026.10.10.0900", "Fix: older change."),
			"2026-10-11-0900-newer-change": bareEntryText("2026.10.11.0900", "Feat: newer change."),
		}))
		sees(t, out, code, 0, "All 1 checks passed")
		passedWith(t, out, "changelog", "2")
	})
	headings := func(t *testing.T) (string, int) {
		return changelogRun(t, withEntries(map[string]string{
			"2026-10-11-0900-good-change":  entryFileText("2026.10.11.0900", "Feat: good change."),
			"2026-10-11-1000-two-headings": entryFileText("2026.10.11.1000", "Fix: one.") + "\n" + entryFileText("2026.10.11.1100", "Fix: two."),
			"2026-10-11-1200-undated":      "### Added\n\n- Feat: one.\n- Fix: two.\n- Docs: three.\n",
		}))
	}
	t.Run("0001.1.3.2.1 it fails asking for one `### yyyy.mm.dd.HHMM` heading with its bullets", func(t *testing.T) {
		t.Parallel()
		out, code := headings(t)
		sees(t, out, code, 1, "changelog/2026-10-11-1000-two-headings.md must hold one ### yyyy.mm.dd.HHMM heading and its bullets, nothing else")
	})
	t.Run("0001.1.3.2.2 an undated heading fails the same way", func(t *testing.T) {
		t.Parallel()
		out, code := headings(t)
		sees(t, out, code, 1, "changelog/2026-10-11-1200-undated.md must hold one ### yyyy.mm.dd.HHMM heading")
	})
	names := func(t *testing.T) (string, int) {
		return changelogRun(t, withEntries(map[string]string{
			"undated-change":               entryFileText("2026.10.11.0900", "Feat: undated change."),
			"2026-10-11-0800-moved-change": entryFileText("2026.10.11.1000", "Fix: moved change."),
		}))
	}
	t.Run("0001.1.3.3.1 it fails naming the file it should be `changelog/<yyyy-mm-dd-HHMM>-<change>.md`", func(t *testing.T) {
		t.Parallel()
		out, code := names(t)
		sees(t, out, code, 1, "changelog/undated-change.md must be named changelog/2026-10-11-0900-undated-change.md: its heading's date, then the change")
	})
	t.Run("0001.1.3.3.2 a name that carries another date fails the same way", func(t *testing.T) {
		t.Parallel()
		out, code := names(t)
		sees(t, out, code, 1, "changelog/2026-10-11-0800-moved-change.md must be named changelog/2026-10-11-1000-moved-change.md")
	})
	t.Run("0001.1.3.4.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, map[string]string{"CHANGELOG.md": ""})
		sees(t, out, code, 0, "All 1 checks passed")
		passedWithNoFiles(t, out, "changelog")
	})
	t.Run("0001.1.3.4.2 a `CHANGELOG.md` makes no difference", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, map[string]string{"CHANGELOG.md": upstreamChangelog})
		sees(t, out, code, 0, "All 1 checks passed")
		passedWithNoFiles(t, out, "changelog")
	})
	upstream := func(t *testing.T) (string, int) {
		files := withEntries(map[string]string{
			"2026-10-10-0900-older-change": entryFileText("2026.10.10.0900", "Fix: older change."),
			"2026-10-11-0900-newer-change": bareEntryText("2026.10.11.0900", "Feat: newer change."),
		})
		files["CHANGELOG.md"] = upstreamChangelog
		files["package.json"] = `{"name": "app", "private": true, "version": "2.0.1"}` + "\n"
		files["package-lock.json"] = `{"name": "app", "version": "2.0.0"}` + "\n"
		return changelogRun(t, files)
	}
	t.Run("0001.1.3.5.1 it is not judged", func(t *testing.T) {
		t.Parallel()
		out, code := upstream(t)
		sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "CHANGELOG.md") {
			t.Errorf("an upstream CHANGELOG.md must not be judged:\n%s", out)
		}
	})
	t.Run("0001.1.3.5.2 a moved package version is not judged either", func(t *testing.T) {
		t.Parallel()
		out, code := upstream(t)
		sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "package") {
			t.Errorf("a moved package version must not be judged:\n%s", out)
		}
	})
	t.Run("0001.1.3.5.3 the entry files still pass", func(t *testing.T) {
		t.Parallel()
		out, code := upstream(t)
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
