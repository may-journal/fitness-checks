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
	t.Run("0001.1.1 a repo with no fitness folder passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := adrChecks(t, acceptance.Example(t, "happyRepo", nil))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "adr-checks")
	})
	t.Run("0001.1.2 a check linking its ADR issue passes and runs", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", adrRepo())
		out, code := adrChecks(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "adr-checks", "1")
		out, code = acceptance.Fitness(t, repo, nil, "check", "adr-0001")
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWith(t, out, "adr-0001", "1")
	})
	t.Run("0001.1.3 a README linking no issue fails unless allowed to wait", func(t *testing.T) {
		t.Parallel()
		// One README has no issue URL, the other one only inside a sentence.
		for _, text := range []string{
			"# Small files\n\nKeeps every Markdown file small.\n\n## Check\n\nFails a long Markdown file.\n",
			"# Small files\n\nKeeps every Markdown file small, as " + adrIssue + " plans.\n\n## Check\n\nFails a long Markdown file.\n",
		} {
			out, code := adrChecks(t, acceptance.Example(t, "happyRepo", acceptance.With(adrRepo(), map[string]string{adrPath: text})))
			acceptance.Sees(t, out, code, 1, adrPath+": links no ADR issue; put the full URL of the `ADR` issue this check proves here")
		}
		waits := acceptance.With(acceptance.HeldCheck(), map[string]string{
			".fitnessrc.json": acceptance.AllowADRs(`{"path": "fitness/adr-*/README.md", "reason": "the checks move with them"}`),
		})
		out, code := adrChecks(t, acceptance.Example(t, "happyRepo", waits))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 an ADR without a check program fails unless superseded", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", map[string]string{adrPath: adrText("")}), nil)
		acceptance.Sees(t, out, code, 1, "fitness setup: local check fitness/adr-0001 has no program; add an executable check or a main.go")
		// Another check keeps the fitness folder from being empty.
		retired := map[string]string{
			adrPath:                        adrText("Status: Superseded by a later decision.\n\n"),
			"fitness/house-rule/README.md": "# House rule\n\nKeeps the house in order.\n",
			"fitness/house-rule/check":     acceptance.VerdictScript(acceptance.Passing),
		}
		out, code = adrChecks(t, acceptance.Example(t, "happyRepo", retired))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "adr-checks", "1")
	})
	t.Run("0001.1.5 a paired check that judged no files fails", func(t *testing.T) {
		t.Parallel()
		empty := `{"ok": true, "errors": [], "filesChecked": 0}`
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", adrCheck("0001", empty)), nil, "check", "adr-0001")
		acceptance.Sees(t, out, code, 1, "adr-0001 judged no files, so its scan is broken")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a README pointed at another issue alone fails naming the README and the check", func(t *testing.T) {
		t.Parallel()
		want := adrPath + " changed its ADR issue without its check fitness/adr-0001; change both together"
		repo := acceptance.Example(t, "happyRepo", adrRepo())
		acceptance.Write(t, repo, map[string]string{adrPath: strings.Replace(adrText(""), adrIssue, adrIssue+"6", 1)})
		acceptance.Git(t, repo, "add", "-A")
		out, code := adrChecks(t, repo)
		acceptance.Sees(t, out, code, 1, want)
		repo = acceptance.Example(t, "happyRepo", adrRepo())
		acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
		acceptance.Commit(t, repo, map[string]string{adrPath: strings.Replace(adrText(""), adrIssue, adrIssue+"7", 1)}, "docs(adr): reword 0001")
		out, code = acceptance.Fitness(t, repo, acceptance.PullRequestEnv, "adr-checks")
		acceptance.Sees(t, out, code, 1, want)
		// The same change with the check's program passes.
		repo = acceptance.Example(t, "happyRepo", adrRepo())
		acceptance.Write(t, repo, map[string]string{
			adrPath:                  strings.Replace(adrText(""), adrIssue, adrIssue+"6", 1),
			"fitness/adr-0001/check": acceptance.VerdictScript(`{"ok": true, "errors": [], "filesChecked": 2}`),
		})
		acceptance.Git(t, repo, "add", "-A")
		out, code = adrChecks(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2 the check's code or README text may change alone", func(t *testing.T) {
		t.Parallel()
		for path, text := range map[string]string{
			"fitness/adr-0001/check": acceptance.VerdictScript(`{"ok": true, "errors": [], "filesChecked": 2}`),
			adrPath:                  strings.Replace(adrText(""), "Small files read fast.", "Small files read fast and diff well.", 1),
		} {
			repo := acceptance.Example(t, "happyRepo", adrRepo())
			acceptance.Write(t, repo, map[string]string{path: text})
			acceptance.Git(t, repo, "add", "-A")
			out, code := adrChecks(t, repo)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		}
	})
	t.Run("0001.2.3 a README linking its first issue needs no code change", func(t *testing.T) {
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
	t.Run("0001.2.4 a full scan judges no changes", func(t *testing.T) {
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
