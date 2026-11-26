package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func Test0003_1(t *testing.T) {
	t.Parallel()
	t.Run("0003.1.1 a section past its word limit fails naming the section with its word count", func(t *testing.T) {
		t.Parallel()
		body := readme + "\n## Notes\n"
		for i := 0; i < 16; i++ {
			body += "\n- " + strings.Repeat("item word ", 10) + "\n"
		}
		out, code := proseBudget(t, map[string]string{"README.md": body})
		acceptance.Sees(t, out, code, 1, `README.md: section "Notes" has 320 prose words (max 300)`)
	})
	t.Run("0003.1.2 the changelog is judged like any doc", func(t *testing.T) {
		t.Parallel()
		changelog := "# Changelog\n\n## Changes\n\n### 2026.10.08.1400\n\n- Docs: " + acceptance.LongSentence
		out, code := proseBudget(t, map[string]string{"CHANGELOG.md": changelog})
		acceptance.Sees(t, out, code, 1, "CHANGELOG.md: a list item has 31 words (max 23)")
	})
	t.Run("0003.1.3 a body past a limit fails", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, map[string]string{"body.md": "## Background\n\n" + acceptance.LongSentence}, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, "a sentence has 30 words (max 23)")
	})
	t.Run("0003.1.4 a list with too many items fails", func(t *testing.T) {
		t.Parallel()
		nested := readme + "\n- Top item\n" + strings.Repeat("    - Child item\n", 5)
		out, code := proseBudget(t, map[string]string{"README.md": nested})
		acceptance.Sees(t, out, code, 1, "README.md: a list has more than 4 items at level 2")
	})
	t.Run("0003.1.5 a long sentence or paragraph fails with the count and the limit", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, map[string]string{"README.md": readme + "\n" + acceptance.LongSentence})
		acceptance.Sees(t, out, code, 1, "README.md: a sentence has 30 words (max 23)")
		out, code = proseBudget(t, map[string]string{"README.md": readme + "\nOne. Two. Three. Four. Five.\n"})
		acceptance.Sees(t, out, code, 1, "README.md: a paragraph has 5 sentences (max 4)")
	})
	t.Run("0003.1.6 fenced code and a Measurement dash pass", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, map[string]string{"README.md": readme + "\n```\n" + acceptance.LongSentence + "```\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// Three paragraphs, the default limit, so only a counted - would push
		// the Measurement over.
		out, code = proseBudget(t, map[string]string{"docs/requirements/0001-adds-numbers.md": acceptance.RequirementDoc("0001.1")})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.1.7 a ticked Plan task within the word limit passes", func(t *testing.T) {
		t.Parallel()
		// The task holds 23 words, the default limit, so only a counted box
		// would push it over.
		task := "- [x] " + strings.Repeat("word ", 22) + "end"
		out, code := proseBudget(t, map[string]string{"plan.md": planBody(task)}, "--body-file", "plan.md")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0003_2(t *testing.T) {
	t.Parallel()
	t.Run("0003.2.1 a repo lowers a prose limit", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, map[string]string{
			".fitnessrc.json": `{"proseBudget": {"maxSentenceWords": 5}}` + "\n",
			"README.md":       readme + "\nThis sentence has exactly seven short words.\n",
		})
		acceptance.Sees(t, out, code, 1, "README.md: a sentence has 7 words (max 5)")
	})
	t.Run("0003.2.2 a repo sets the comment word numbers", func(t *testing.T) {
		t.Parallel()
		long := "\t// Hands the number back to the caller without any change at all.\n\n\t// Callers rely on getting their own number back from this function.\n"
		files := func(cfg string) map[string]string {
			f := map[string]string{"app/run.go": inFunctions(long, long, long)["app/run.go"]}
			if cfg != "" {
				f[".fitnessrc.json"] = cfg + "\n"
			}
			return f
		}
		out, code := proseBudget(t, files(""))
		acceptance.Sees(t, out, code, 1, "overlong", "(max 3 words per token)")
		out, code = proseBudget(t, files(`{"prose": {"comments": {"maxWordsPerToken": 20}}}`))
		acceptance.Sees(t, out, code, 0)
		out, code = proseBudget(t, files(`{"prose": {"comments": {"maxWordsPerToken": 2}}}`))
		acceptance.Sees(t, out, code, 1, "(max 2 words per token)")
		out, code = proseBudget(t, files(`{"prose": {"comments": {"minBodyComments": 3}}}`))
		acceptance.Sees(t, out, code, 0)
	})
}
