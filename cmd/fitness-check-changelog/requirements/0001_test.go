// Package requirements proves the changelog check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
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
	t.Run("0001.1.1 entry files with a dated heading under a dated name pass, counting them", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, withEntries(map[string]string{
			"2026-10-10-0900-older-change": entryFileText("2026.10.10.0900", "Fix: older change."),
			"2026-10-11-0900-newer-change": bareEntryText("2026.10.11.0900", "Feat: newer change."),
			"2026-10-12-0900-third-change": bareEntryText("2026.10.12.0900", "Docs: third change."),
		}))
		req.Then(t, "0001.1.3.1.1", "17360cd9", func(a req.A) {
			sees(a.T(), out, code, 0, "All 1 checks passed")
			passedWith(a.T(), out, "changelog", "3")
		})
		req.Then(t, "0001.1.3.1.2", "da6aa874", func(a req.A) {
			// Two of the three entry files carry no front matter.
			sees(a.T(), out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.2 two dated or undated headings fail", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, withEntries(map[string]string{
			"2026-10-11-0900-good-change":  entryFileText("2026.10.11.0900", "Feat: good change."),
			"2026-10-11-1000-two-headings": entryFileText("2026.10.11.1000", "Fix: one.") + "\n" + entryFileText("2026.10.11.1100", "Fix: two."),
			"2026-10-11-1200-undated":      "### Added\n\n- Feat: one.\n- Fix: two.\n- Docs: three.\n",
		}))
		req.Then(t, "0001.1.3.2.1", "c191df5f", func(a req.A) {
			sees(a.T(), out, code, 1, "changelog/2026-10-11-1000-two-headings.md must hold one ### yyyy.mm.dd.HHMM heading and its bullets, nothing else")
		})
		req.Then(t, "0001.1.3.2.2", "64c245ff", func(a req.A) {
			sees(a.T(), out, code, 1, "changelog/2026-10-11-1200-undated.md must hold one ### yyyy.mm.dd.HHMM heading")
		})
	})
	t.Run("0001.1.3 a name that lacks its heading's date fails", func(t *testing.T) {
		t.Parallel()
		out, code := changelogRun(t, withEntries(map[string]string{
			"undated-change":               entryFileText("2026.10.11.0900", "Feat: undated change."),
			"2026-10-11-0800-moved-change": entryFileText("2026.10.11.1000", "Fix: moved change."),
		}))
		req.Then(t, "0001.1.3.3.1", "b10e3ae0", func(a req.A) {
			sees(a.T(), out, code, 1, "changelog/undated-change.md must be named changelog/2026-10-11-0900-undated-change.md: its heading's date, then the change")
		})
		req.Then(t, "0001.1.3.3.2", "848bcb6d", func(a req.A) {
			sees(a.T(), out, code, 1, "changelog/2026-10-11-0800-moved-change.md must be named changelog/2026-10-11-1000-moved-change.md")
		})
	})
	t.Run("0001.1.4 a repo with no entry files passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.3.4.1", "857dd667", func(a req.A) {
			out, code := changelogRun(t, map[string]string{"CHANGELOG.md": ""})
			sees(a.T(), out, code, 0, "All 1 checks passed")
			passedWithNoFiles(a.T(), out, "changelog")
		})
		req.Then(t, "0001.1.3.4.2", "28de5c62", func(a req.A) {
			out, code := changelogRun(t, map[string]string{"CHANGELOG.md": upstreamChangelog})
			sees(a.T(), out, code, 0, "All 1 checks passed")
			passedWithNoFiles(a.T(), out, "changelog")
		})
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
		req.Then(t, "0001.1.3.5.1", "ebb2c064", func(a req.A) {
			sees(a.T(), out, code, 0, "All 1 checks passed")
			a.True(!strings.Contains(out, "CHANGELOG.md"), "the upstream history was judged:\n%s", out)
		})
		req.Then(t, "0001.1.3.5.2", "ba2ffd97", func(a req.A) {
			sees(a.T(), out, code, 0, "All 1 checks passed")
			a.True(!strings.Contains(out, "package"), "a moved package version was judged:\n%s", out)
		})
		req.Then(t, "0001.1.3.5.3", "68390325", func(a req.A) {
			passedWith(a.T(), out, "changelog", "2")
		})
	})
}

// changelogRun runs `fitness-install -- changelog` on happyRepo with files
// written over it.
func changelogRun(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "changelog")
}
