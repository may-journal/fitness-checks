// Package requirements proves the markdown-no-bold-italic check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
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
		stars, starsCode := readme(t, "Some **strong** and *soft* words.\n")
		snake, snakeCode := readme(t, "Some __strong__ and _soft_ words.\n")
		req.Then(t, "0001.1.2.1.1", "98b7a18a", func(a req.A) {
			acceptance.Sees(t, stars, starsCode, 1,
				`README.md: disallowed **bold**`+rule+`"**strong**"`,
				`README.md: disallowed *italic*`+rule+`"*soft*"`)
			acceptance.Sees(t, snake, snakeCode, 1,
				`README.md: disallowed __bold__`+rule+`"__strong__"`,
				`README.md: disallowed _italic_`+rule+`"_soft_"`)
			a.Equal(2, strings.Count(stars, "disallowed"))
			a.Equal(2, strings.Count(snake, "disallowed"))
		})
		req.Then(t, "0001.1.2.1.2", "6f2cb2c1", func(a req.A) {
			acceptance.Sees(t, stars, starsCode, 1, `README.md: disallowed **bold**`, `README.md: disallowed *italic*`)
			acceptance.Sees(t, snake, snakeCode, 1, `README.md: disallowed __bold__`, `README.md: disallowed _italic_`)
		})
		req.Then(t, "0001.1.2.1.3", "62245b6c", func(a req.A) {
			acceptance.Sees(t, stars, starsCode, 1, `"**strong**"`, `"*soft*"`)
			acceptance.Sees(t, snake, snakeCode, 1, `"__strong__"`, `"_soft_"`)
		})
	})
	t.Run("0001.1.2 code and links pass", func(t *testing.T) {
		t.Parallel()
		code, codeCode := readme(t, "Use `**strong**` here.\n\n```\n_soft_ and **strong**\n```\n")
		link, linkCode := readme(t, "See [the _dev_ path](https://example.com/_dev_/to_page_).\n")
		req.Then(t, "0001.1.2.2.1", "dbfdaa5e", func(a req.A) {
			acceptance.Sees(t, code, codeCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.2.2", "8fc6f86d", func(a req.A) {
			acceptance.Sees(t, link, linkCode, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.3 asterisk list markers pass", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.3.1", "c777fdc5", func(_ req.A) {
			out, code := readme(t, "* first item\n* second item\n* third item\n")
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.4 a body is judged as the description", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.4.1", "7b3c06c4", func(_ req.A) {
			out, code := emphasis(t, map[string]string{"body.md": "Hello _there_.\n"}, "--body-file", "body.md")
			acceptance.Sees(t, out, code, 1, `(description): disallowed _italic_`+rule+`"_there_"`)
		})
	})
	t.Run("0001.1.5 only bold outside a release scope fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.5.1", "bbfbe6cf", func(_ req.A) {
			changelog := "# Changelog\n\n## Changes\n\n## [1.1.1](https://github.com/o/r/compare/v1.1.0...v1.1.1) (2026-10-10)\n\n" +
				"### Features\n\n* **install:** change **everything** ([#239](https://github.com/o/r/issues/239))\n\n" +
				"### 2026.10.09.1500\n\n- Fix: one.\n"
			out, code := emphasis(t, map[string]string{"CHANGELOG.md": changelog})
			acceptance.Sees(t, out, code, 1, `CHANGELOG.md: disallowed **bold**`+rule+`"**everything**"`)
			if n := strings.Count(out, "disallowed"); n != 1 {
				t.Errorf("only the bold outside a scope must fail, got %d findings:\n%s", n, out)
			}
		})
	})
	t.Run("0001.1.6 a repo with no Markdown passes with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "ebe8880b", func(_ req.A) {
			out, code := emphasis(t, map[string]string{"README.md": "", "CHANGELOG.md": ""})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "markdown-no-bold-italic")
		})
	})
}
