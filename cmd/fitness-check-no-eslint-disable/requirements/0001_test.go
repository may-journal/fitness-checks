// Package requirements proves the no-eslint-disable check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// next is a next-line directive comment.
const next = "// eslint-disable-next-line no-console\n"

// exts are the extensions the check scans.
var exts = []string{".cts", ".mts", ".ts", ".cjs", ".js", ".mjs", ".tsx"}

// run runs the check in repo.
func run(t *testing.T, repo string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, repo, nil, "no-eslint-disable")
}

// check runs the check on happyRepo with files written over it.
func check(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return run(t, acceptance.Example(t, "happyRepo", files))
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 every directive form is listed with its file line and directive", func(t *testing.T) {
		t.Parallel()
		legacy := "export const x = 1;\n" + next + "console.log(x);\n" +
			"const y = 2; // eslint-disable-line\n/* eslint-disable */\n"
		out, code := check(t, map[string]string{"src/legacy.ts": legacy})
		req.Then(t, "0001.1.2.1.1", "a767c42a", func(_ req.A) {
			acceptance.Sees(t, out, code, 1,
				"src/legacy.ts:2: eslint-disable-next-line",
				"src/legacy.ts:4: eslint-disable-line",
				"src/legacy.ts:5: eslint-disable")
		})
		req.Then(t, "0001.1.2.1.2", "db673b63", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "src/legacy.ts:2:", "src/legacy.ts:4:", "src/legacy.ts:5:")
		})
		req.Then(t, "0001.1.2.1.3", "f1876531", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, ": eslint-disable-next-line", ": eslint-disable-line", ": eslint-disable")
		})
		req.Then(t, "0001.1.2.1.4", "0dc3958d", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "src/legacy.ts:2: eslint-disable-next-line", "src/legacy.ts:4: eslint-disable-line")
		})
	})
	t.Run("0001.1.2 files without directives pass and count", func(t *testing.T) {
		t.Parallel()
		out, code := check(t, map[string]string{
			"src/a.ts": "export const a = 1;\n",
			"src/b.js": "export const b = 2;\n",
		})
		req.Then(t, "0001.1.2.2.1", "07b8e753", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.2.2", "8b2d3768", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "2 files scanned")
		})
	})
	t.Run("0001.1.3 a string holding directive text fails in every extension", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.3.1", "b39697d1", func(_ req.A) {
			files := map[string]string{}
			var want []string
			for _, ext := range exts {
				files["src/f"+ext] = "export const hint = 'add eslint-disable-next-line to skip';\n"
				want = append(want, "src/f"+ext+":1: eslint-disable-next-line")
			}
			out, code := check(t, files)
			acceptance.Sees(t, out, code, 1, append(want, "7 files scanned")...)
		})
	})
	t.Run("0001.1.4 a Markdown file holding directive text passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.4.1", "dbfdaa5e", func(_ req.A) {
			out, code := check(t, map[string]string{
				"docs/lint.md": "Never write " + next,
				"src/a.ts":     "export const a = 1;\n",
			})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		})
	})
	t.Run("0001.1.5 a repo with no source passes with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.5.1", "ebe8880b", func(_ req.A) {
			out, code := check(t, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "no-eslint-disable")
		})
	})
	t.Run("0001.1.6 node_modules and dist are scanned", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "5fe9535b", func(_ req.A) {
			out, code := check(t, map[string]string{
				"node_modules/pkg/index.js": "/* eslint-disable */\n",
				"dist/out.js":               "/* eslint-disable */\n",
				"src/clean.ts":              "export const a = 1;\n",
			})
			acceptance.Sees(t, out, code, 1,
				"node_modules/pkg/index.js:1: eslint-disable",
				"dist/out.js:1: eslint-disable",
				"3 files scanned")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a directive file that is untracked passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.1.1", "07b8e753", func(_ req.A) {
			repo := acceptance.Example(t, "happyRepo", map[string]string{"src/a.ts": "export const a = 1;\n"})
			acceptance.Write(t, repo, map[string]string{"src/scratch.ts": next})
			out, code := run(t, repo)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		})
	})
	t.Run("0001.2.2 a scan of staged files reads only the staged file", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.2.1", "75833105", func(_ req.A) {
			repo := acceptance.Example(t, "happyRepo", map[string]string{"src/old.ts": next})
			acceptance.Write(t, repo, map[string]string{"src/new.ts": "export const a = 1;\n"})
			acceptance.Git(t, repo, "add", "src/new.ts")
			out, code := run(t, repo)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		})
	})
}
