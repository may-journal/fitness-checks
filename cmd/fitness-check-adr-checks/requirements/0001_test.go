// Package requirements proves the adr-checks check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := adrChecks(t, acceptance.Example(t, "happyRepo", nil))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "adr-checks")
	})
	t.Run("0001.1.2.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := adrChecks(t, acceptance.Example(t, "happyRepo", adrRepo()))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "adr-checks", "1")
	})
	t.Run("0001.1.2.2.2 the paired check runs", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", adrRepo()), nil, "check", "adr-0001")
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWith(t, out, "adr-0001", "1")
	})
	unlinked := func(t *testing.T, text string) (string, int) {
		return adrChecks(t, acceptance.Example(t, "happyRepo", acceptance.With(adrRepo(), map[string]string{adrPath: text})))
	}
	const noLink = adrPath + ": links no ADR issue; put the full URL of the `ADR` issue this check proves here"
	t.Run("0001.1.2.3.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := unlinked(t, "# Small files\n\nKeeps every Markdown file small.\n\n## Check\n\nFails a long Markdown file.\n")
		acceptance.Sees(t, out, code, 1, noLink)
	})
	t.Run("0001.1.2.3.2 a link only inside a sentence fails too", func(t *testing.T) {
		t.Parallel()
		out, code := unlinked(t, "# Small files\n\nKeeps every Markdown file small, as "+adrIssue+" plans.\n\n## Check\n\nFails a long Markdown file.\n")
		acceptance.Sees(t, out, code, 1, noLink)
	})
	t.Run("0001.1.2.3.3 an `adrAsIssue` allow glob lets it wait", func(t *testing.T) {
		t.Parallel()
		waits := acceptance.With(acceptance.HeldCheck(), map[string]string{
			".fitnessrc.json": acceptance.AllowADRs(`{"path": "fitness/adr-*/README.md", "reason": "the checks move with them"}`),
		})
		out, code := adrChecks(t, acceptance.Example(t, "happyRepo", waits))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.4.1 the run fails before any check starts", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", map[string]string{adrPath: adrText("")}), nil)
		acceptance.Sees(t, out, code, 1, "fitness setup: local check fitness/adr-0001 has no program; add an executable check or a main.go")
	})
	t.Run("0001.1.2.4.2 a superseded ADR is exempt", func(t *testing.T) {
		t.Parallel()
		// Another check keeps the fitness folder from being empty.
		retired := map[string]string{
			adrPath:                        adrText("Status: Superseded by a later decision.\n\n"),
			"fitness/house-rule/README.md": "# House rule\n\nKeeps the house in order.\n",
			"fitness/house-rule/check":     acceptance.VerdictScript(acceptance.Passing),
		}
		out, code := adrChecks(t, acceptance.Example(t, "happyRepo", retired))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "adr-checks", "1")
	})
	t.Run("0001.1.2.5.1 it fails as a broken scan", func(t *testing.T) {
		t.Parallel()
		empty := `{"ok": true, "errors": [], "filesChecked": 0}`
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", adrCheck("0001", empty)), nil, "check", "adr-0001")
		acceptance.Sees(t, out, code, 1, "adr-0001 judged no files, so its scan is broken")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	const want = adrPath + " changed its ADR issue without its check fitness/adr-0001; change both together"
	staged := func(t *testing.T) (string, int) {
		repo := acceptance.Example(t, "happyRepo", adrRepo())
		acceptance.Write(t, repo, map[string]string{adrPath: strings.Replace(adrText(""), adrIssue, adrIssue+"6", 1)})
		acceptance.Git(t, repo, "add", "-A")
		return adrChecks(t, repo)
	}
	inPullRequest := func(t *testing.T) (string, int) {
		repo := acceptance.Example(t, "happyRepo", adrRepo())
		acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
		acceptance.Commit(t, repo, map[string]string{adrPath: strings.Replace(adrText(""), adrIssue, adrIssue+"7", 1)}, "docs(adr): reword 0001")
		return acceptance.Fitness(t, repo, acceptance.PullRequestEnv, "adr-checks")
	}
	t.Run("0001.2.1.1.1 it fails naming the README", func(t *testing.T) {
		t.Parallel()
		out, code := staged(t)
		acceptance.Sees(t, out, code, 1, adrPath+" changed its ADR issue")
	})
	t.Run("0001.2.1.1.2 it fails naming the check", func(t *testing.T) {
		t.Parallel()
		out, code := staged(t)
		acceptance.Sees(t, out, code, 1, "without its check fitness/adr-0001; change both together")
	})
	t.Run("0001.2.1.1.3 it fails the same way staged or in a pull request", func(t *testing.T) {
		t.Parallel()
		out, code := staged(t)
		acceptance.Sees(t, out, code, 1, want)
		out, code = inPullRequest(t)
		acceptance.Sees(t, out, code, 1, want)
	})
	t.Run("0001.2.1.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", adrRepo())
		acceptance.Write(t, repo, map[string]string{"fitness/adr-0001/check": acceptance.VerdictScript(`{"ok": true, "errors": [], "filesChecked": 2}`)})
		acceptance.Git(t, repo, "add", "-A")
		out, code := adrChecks(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.2.2 a change to only the README's own text passes too", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", adrRepo())
		acceptance.Write(t, repo, map[string]string{adrPath: strings.Replace(adrText(""), "Small files read fast.", "Small files read fast and diff well.", 1)})
		acceptance.Git(t, repo, "add", "-A")
		out, code := adrChecks(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.3.1 it passes without a code change", func(t *testing.T) {
		t.Parallel()
		unlinked := strings.Replace(adrText(""), adrIssue+"\n\n", "", 1)
		repo := acceptance.Example(t, "happyRepo", acceptance.With(adrRepo(), map[string]string{
			adrPath:           unlinked,
			".fitnessrc.json": acceptance.AllowADRs(`{"path": "fitness/adr-*/README.md", "reason": "the issue is not open yet"}`),
		}))
		acceptance.Write(t, repo, map[string]string{adrPath: adrText(""), ".fitnessrc.json": "{}\n"})
		acceptance.Git(t, repo, "add", "-A")
		out, code := adrChecks(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.4.1 it judges no changes", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", adrRepo())
		acceptance.Commit(t, repo, map[string]string{adrPath: strings.Replace(adrText(""), adrIssue, adrIssue+"6", 1)}, "docs(adr): relink 0001")
		out, code := adrChecks(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

// adrChecks runs `fitness-install -- adr-checks` in repo.
func adrChecks(t *testing.T, repo string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, repo, nil, "adr-checks")
}

// adrPath is the one ADR adrRepo holds: the README of its check.
const adrPath = "fitness/adr-0001/README.md"

// adrIssue is the URL of the ADR issue the README of ADR 0001 links.
const adrIssue = "https://github.com/may-journal/example/issues/5"

// adrText is the README of ADR 0001: its title, summary, status, the ADR
// issue it links, and what its check fails.
func adrText(status string) string {
	return "# 0001 — Small files\n\nSmall files read fast.\n\n" + status + adrIssue + "\n\n## Check\n\nFails a long Markdown file.\n"
}

// adrRepo is a repo with ADR 0001 as the README of its check.
func adrRepo() map[string]string {
	return adrCheck("0001", acceptance.Passing)
}

// adrCheck is the check folder fitness/adr-<number>, whose README is ADR
// 0001 and whose program prints verdict.
func adrCheck(number, verdict string) map[string]string {
	folder := "fitness/adr-" + number + "/"
	return map[string]string{
		folder + "README.md": adrText(""),
		folder + "check":     acceptance.VerdictScript(verdict),
	}
}
