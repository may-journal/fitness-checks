// Package requirements proves the changelog-bullets check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
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
	withEntries       = acceptance.WithEntries
	upstreamChangelog = acceptance.UpstreamChangelog
	example           = acceptance.Example
	fitness           = acceptance.Fitness
	read              = acceptance.Read
	sees              = acceptance.Sees
	passedWith        = acceptance.PassedWith
	passedWithNoFiles = acceptance.PassedWithNoFiles
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 it passes counting the entry files", func(t *testing.T) {
		t.Parallel()
		files := withEntries(map[string]string{
			"2026-10-08-1400-entry":  changelogBulletsDoc("Feat: one.", "Fix: two.", "Docs: three.", "Style: four."),
			"2026-10-08-1400-later":  changelogBulletsDoc("Refactor: one.", "Perf: two.", "Test: three.", "Build: four.", "Ci: five."),
			"2026-10-08-1400-latest": changelogBulletsDoc("Chore: one.", "Revert: two.", "feat(release)!: three.", "fix!: four."),
		})
		out, code := fitness(t, example(t, "happyRepo", files), nil, "changelog-bullets")
		sees(t, out, code, 0, "All 1 checks passed")
		passedWith(t, out, "changelog-bullets", "3")
	})
	t.Run("0001.1.2.1.2 types read `Feat:` or `feat(scope)!:` alike", func(t *testing.T) {
		t.Parallel()
		out, code := changelogBullets(t, changelogBulletsDoc("Feat: one.", "feat: two.", "feat(release)!: three.", "fix!: four.", "fix(ci): five."))
		sees(t, out, code, 0, "All 1 checks passed")
	})
	few := func(t *testing.T) (string, int) {
		return changelogBullets(t, changelogBulletsDoc("Docs: add the readme.", "Docs: add the changelog."))
	}
	t.Run("0001.1.2.2.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := few(t)
		sees(t, out, code, 1, bulletsFile+` section`)
	})
	t.Run("0001.1.2.2.2 it fails naming the section", func(t *testing.T) {
		t.Parallel()
		out, code := few(t)
		sees(t, out, code, 1, `section "2026.10.08.1400"`)
	})
	t.Run("0001.1.2.2.3 it fails naming its count", func(t *testing.T) {
		t.Parallel()
		out, code := few(t)
		sees(t, out, code, 1, `has 2 bullets; keep entries to 3-5 bullets`)
		out, code = changelogBullets(t, "### 2026.10.08.1400\n\n* Fix: one.\n* Fix: two.\n* Fix: three.\n")
		sees(t, out, code, 1, `"2026.10.08.1400" has 0 bullets`)
	})
	t.Run("0001.1.2.2.4 a section of more than 5 fails the same way", func(t *testing.T) {
		t.Parallel()
		out, code := changelogBullets(t, changelogBulletsDoc("Feat: one.", "Fix: two.", "Fix: three.", "Test: four.", "Docs: five.", "Docs: six."))
		sees(t, out, code, 1, bulletsFile+` section "2026.10.08.1400" has 6 bullets`)
	})
	long := func(t *testing.T) (string, int) {
		long := "Feat: " + strings.Repeat("x", 359)
		wrapped := "Fix: " + strings.Repeat("x", 200) + "\n  " + strings.Repeat("y", 200)
		return changelogBullets(t, changelogBulletsDoc(long, wrapped, "Docs: three."))
	}
	t.Run("0001.1.2.3.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := long(t)
		sees(t, out, code, 1, bulletsFile+":3: bullet is", bulletsFile+":4: bullet is")
	})
	t.Run("0001.1.2.3.2 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := long(t)
		sees(t, out, code, 1, bulletsFile+":3:", bulletsFile+":4:")
	})
	t.Run("0001.1.2.3.3 it fails naming the length", func(t *testing.T) {
		t.Parallel()
		out, code := long(t)
		sees(t, out, code, 1, "bullet is 365 characters; keep each under 365", "bullet is 406 characters; keep each under 365")
	})
	t.Run("0001.1.2.4.1 it fails that bullet alone naming its line", func(t *testing.T) {
		t.Parallel()
		out, code := changelogBullets(t, changelogBulletsDoc("Feat: one.", "fix(release)!: two.", "chore(ci): three.", "Revert: four.", "Added the readme."))
		sees(t, out, code, 1, bulletsFile+":7: bullet must start with a semantic type")
		if n := strings.Count(read(out), "must start with a semantic type"); n != 1 {
			t.Errorf("only the untyped bullet must fail, got %d failures:\n%s", n, out)
		}
	})
	t.Run("0001.1.2.4.2 a commit subject of such a type passes", func(t *testing.T) {
		t.Parallel()
		out, code := changelogBullets(t, changelogBulletsDoc("fix(release)!: two.", "chore(ci): three.", "Revert: four."))
		sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.5.1 each error names its own file", func(t *testing.T) {
		t.Parallel()
		files := withEntries(map[string]string{
			"short-entry":   "### 2026.10.11.0900\n\n- Feat: one.\n- Fix: two.\n",
			"untyped-entry": "### 2026.10.11.1000\n\n- Feat: one.\n- Fix: two.\n- no type here.\n",
		})
		out, code := fitness(t, example(t, "happyRepo", files), nil, "changelog-bullets")
		sees(t, out, code, 1, `changelog/short-entry.md section "2026.10.11.0900" has 2 bullets`,
			"changelog/untyped-entry.md:5: bullet must start with a semantic type")
	})
	t.Run("0001.1.2.6.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		repo := example(t, "happyRepo", map[string]string{"CHANGELOG.md": upstreamChangelog})
		out, code := fitness(t, repo, nil, "changelog-bullets")
		sees(t, out, code, 0, "All 1 checks passed")
		passedWithNoFiles(t, out, "changelog-bullets")
	})
}

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
