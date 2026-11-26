// Package requirements proves the no-fitness-front-matter check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// fix is what the check prints after each file that holds a key.
const fix = ": front matter holds fitnessFunctions or relatedConfigurations, which no check reads; run `fitness strip-front-matter`"

// block is a plain block holding both fitness keys.
const block = "---\nfitnessFunctions:\n  - prose-budget\nrelatedConfigurations: ['../.fitnessrc.json']\n---\n\n"

// byteOrderMark is the UTF-8 BOM some editors put before a file's first line.
const byteOrderMark = string(rune(0xFEFF))

// front writes files over happyRepo and runs the check there.
func front(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "no-fitness-front-matter")
}

// flagged asserts out names each file once with the fix.
func flagged(t *testing.T, out string, files ...string) {
	t.Helper()
	for _, f := range files {
		if n := strings.Count(acceptance.Read(out), acceptance.Read(f+fix)); n != 1 {
			t.Errorf("%s is named %d times, want once:\n%s", f, n, out)
		}
	}
}

// keyFiles holds a fitness key in a doc, a changelog entry and a PR template.
func keyFiles() map[string]string {
	return map[string]string{
		"docs/guide.md":                    block + "# Guide\n",
		"changelog/2026-10-11-0900-one.md": "---\nrelatedConfigurations: ['../.fitnessrc.json']\n---\n\n### 2026.10.11.0900\n\n- Feat: one change.\n",
		".github/PULL_REQUEST_TEMPLATE.md": "<!--\n---\nfitnessFunctions: ['pr-structure']\n---\n-->\n<!-- One sentence. -->\n> REPLACE-ME\n",
	}
}

// edgeFiles holds a key behind CRLF, a byte order mark and other keys.
func edgeFiles() map[string]string {
	return map[string]string{
		"docs/crlf.md":                  "---\r\nrelatedConfigurations: ['../.fitnessrc.json']\r\n---\r\n\r\n# Crlf\r\n",
		"docs/bom.md":                   byteOrderMark + "---\nfitnessFunctions: ['prose-budget']\n---\n\n# Bom\n",
		".github/ISSUE_TEMPLATE/bug.md": "---\nname: Bug\nrelatedConfigurations: ['../../.fitnessrc.json']\nlabels: Bug\n---\n\n## Steps\n",
	}
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 it fails once per file", func(t *testing.T) {
		t.Parallel()
		out, code := front(t, keyFiles())
		acceptance.Sees(t, out, code, 1)
		flagged(t, out, "docs/guide.md", "changelog/2026-10-11-0900-one.md", ".github/PULL_REQUEST_TEMPLATE.md")
	})
	t.Run("0001.1.2.1.2 the failure names `fitness strip-front-matter`", func(t *testing.T) {
		t.Parallel()
		out, code := front(t, keyFiles())
		acceptance.Sees(t, out, code, 1, "fitness strip-front-matter")
	})
	t.Run("0001.1.2.2.1 it still fails", func(t *testing.T) {
		t.Parallel()
		out, code := front(t, edgeFiles())
		acceptance.Sees(t, out, code, 1)
		flagged(t, out, "docs/crlf.md", "docs/bom.md")
	})
	t.Run("0001.1.2.2.2 other keys beside the fitness key do not save it", func(t *testing.T) {
		t.Parallel()
		out, code := front(t, edgeFiles())
		acceptance.Sees(t, out, code, 1)
		flagged(t, out, ".github/ISSUE_TEMPLATE/bug.md")
	})
	t.Run("0001.1.2.3.1 the file passes", func(t *testing.T) {
		t.Parallel()
		out, code := front(t, map[string]string{
			"docs/notes.md":                 "---\ntitle: Notes\n---\n\n# Notes\n",
			".github/ISSUE_TEMPLATE/bug.md": "---\nname: Bug\nlabels: Bug\n---\n\n## Steps\n",
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "strip-front-matter") {
			t.Errorf("a file without fitness front matter was flagged:\n%s", out)
		}
	})
	t.Run("0001.1.2.4.1 the file passes", func(t *testing.T) {
		t.Parallel()
		out, code := front(t, map[string]string{
			"docs/ruled.md": "---\n\nrelatedConfigurations: a note set between two rules.\n\n---\n\n# Ruled\n",
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "strip-front-matter") {
			t.Errorf("a leading rule was flagged:\n%s", out)
		}
	})
	t.Run("0001.1.2.5.1 it is never judged", func(t *testing.T) {
		t.Parallel()
		out, code := front(t, map[string]string{
			"tests/examples/demoRepo/README.md": block + "# Demo\n",
			"go/tool/testdata/guide.md":         block + "# Guide\n",
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "strip-front-matter") {
			t.Errorf("a fixture was flagged:\n%s", out)
		}
	})
}
