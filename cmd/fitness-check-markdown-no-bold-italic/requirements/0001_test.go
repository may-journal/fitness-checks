// Package requirements proves the markdown-no-bold-italic check's
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

// rule is the advice every finding carries.
const rule = " (use only when explicitly required): "

// emphasis runs the check on happyRepo with files written over it.
func emphasis(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, append([]string{"markdown-no-bold-italic"}, args...)...)
}

// readme runs the check with text as the README.
func readme(t *testing.T, text string) (string, int) {
	t.Helper()
	return emphasis(t, map[string]string{"README.md": "# App\n\n" + text})
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 emphasis fails once per occurrence", func(t *testing.T) {
		t.Parallel()
		out, code := readme(t, "Some **strong** and *soft* words.\n")
		acceptance.Sees(t, out, code, 1,
			`README.md: disallowed **bold**`+rule+`"**strong**"`,
			`README.md: disallowed *italic*`+rule+`"*soft*"`)
		out, code = readme(t, "Some __strong__ and _soft_ words.\n")
		acceptance.Sees(t, out, code, 1,
			`README.md: disallowed __bold__`+rule+`"__strong__"`,
			`README.md: disallowed _italic_`+rule+`"_soft_"`)
	})
	t.Run("0001.1.2 code and links pass", func(t *testing.T) {
		t.Parallel()
		out, code := readme(t, "Use `**strong**` here.\n\n```\n_soft_ and **strong**\n```\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = readme(t, "See [the _dev_ path](https://example.com/_dev_/to_page_).\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.3 asterisk list markers pass", func(t *testing.T) {
		t.Parallel()
		out, code := readme(t, "* first item\n* second item\n* third item\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 a body is judged as the description", func(t *testing.T) {
		t.Parallel()
		out, code := emphasis(t, map[string]string{"body.md": "Hello _there_.\n"}, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, `(description): disallowed _italic_`+rule+`"_there_"`)
	})
	t.Run("0001.1.5 only bold outside a release scope fails", func(t *testing.T) {
		t.Parallel()
		changelog := "# Changelog\n\n## Changes\n\n## [1.1.1](https://github.com/o/r/compare/v1.1.0...v1.1.1) (2026-10-10)\n\n" +
			"### Features\n\n* **install:** change **everything** ([#239](https://github.com/o/r/issues/239))\n\n" +
			"### 2026.10.09.1500\n\n- Fix: one.\n"
		out, code := emphasis(t, map[string]string{"CHANGELOG.md": changelog})
		acceptance.Sees(t, out, code, 1, `CHANGELOG.md: disallowed **bold**`+rule+`"**everything**"`)
		if n := strings.Count(out, "disallowed"); n != 1 {
			t.Errorf("only the bold outside a scope must fail, got %d findings:\n%s", n, out)
		}
	})
	t.Run("0001.1.6 a repo with no Markdown passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := emphasis(t, map[string]string{"README.md": "", "CHANGELOG.md": ""})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "markdown-no-bold-italic")
	})
}
