// Package requirements proves the pr-closes-issue check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 any closing keyword in any reference form passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.1.1", "1c7a9eef", func(_ req.A) {
			files := map[string]string{}
			for _, keyword := range prClosesKeywords {
				files[keyword+".md"] = "> pitch\n\n" + strings.ToUpper(keyword[:1]) + keyword[1:] + ": #12\n"
			}
			repo := acceptance.Example(t, "happyRepo", files)
			for _, keyword := range prClosesKeywords {
				out, code := acceptance.Fitness(t, repo, nil, "pr-closes-issue", "--body-file", keyword+".md")
				acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			}
			for _, body := range []string{"Closes May-Journal/.github#144\n", "Fixes https://github.com/may-journal/.github/issues/144\n"} {
				out, code := prCloses(t, body)
				acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			}
		})
	})
	t.Run("0001.1.2 a reference only or spaced number fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "62f6ad5c", func(_ req.A) {
			out, code := prCloses(t, "This addresses #12, is part of #6, and follows #7. Closes # 12, and closes #\n12.\n")
			acceptance.Sees(t, out, code, 1, prClosesNoKeyword)
		})
	})
	t.Run("0001.1.3 an implemented issue left open fails naming each", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.3.1", "924a4d11", func(_ req.A) {
			out, code := prCloses(t, "Implements #5.\n\nPlan #63 for context.\n\nCloses #6.\n")
			acceptance.Sees(t, out, code, 1,
				"PR says it implements #5 but no closing keyword (close/fix/resolve + #5) closes it",
				"PR says it implements #63 but no closing keyword (close/fix/resolve + #63) closes it")
		})
	})
	t.Run("0001.1.4 a keyword only inside code or a comment fails", func(t *testing.T) {
		t.Parallel()
		out, code := prCloses(t, "Merging `Closes #12` closed it.\n\n```\nFixes #12\n```\n\n<!-- Resolves #12 -->\n")
		req.Then(t, "0001.1.1.4.1", "88f49e4e", func(a req.A) {
			acceptance.Sees(t, out, code, 1, prClosesNoKeyword)
			a.Equal(1, strings.Count(out, "has no GitHub closing keyword"))
			a.True(!strings.Contains(out, "says it implements"), "a hidden keyword was reported as an unclosed implemented issue:\n%s", out)
		})
	})
	t.Run("0001.1.5 a required issue left open fails naming it", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.5.1", "c80ab29e", func(_ req.A) {
			out, code := prCloses(t, "Closes #12.\n", "--require-close=12,40")
			acceptance.Sees(t, out, code, 1, "PR closes a Plan that closes #40, so the PR must also close #40")
		})
	})
	t.Run("0001.1.6 every required issue closed passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.6.1", "07b8e753", func(_ req.A) {
			out, code := prCloses(t, "Closes #12 and closes #40.\n", "--require-close=12,40")
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "pr-closes-issue")
	req.Then(t, "0001.2.1.1.1", "07b8e753", func(_ req.A) {
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	req.Then(t, "0001.2.1.1.2", "41c3dea8", func(_ req.A) {
		acceptance.PassedWithNoFiles(t, out, "pr-closes-issue")
	})
}

// prCloses runs `fitness-install -- pr-closes-issue --body-file body.md
// <args>` on happyRepo with body as the PR description.
func prCloses(t *testing.T, body string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{"body.md": body})
	return acceptance.Fitness(t, repo, nil, append([]string{"pr-closes-issue", "--body-file", "body.md"}, args...)...)
}

// prClosesNoKeyword is the failure a body with no closing keyword gets.
const prClosesNoKeyword = "PR description has no GitHub closing keyword (close/fix/resolve + #NN)"

// prClosesKeywords is GitHub's fixed set of issue-closing keywords.
var prClosesKeywords = strings.Fields("close closes closed fix fixes fixed resolve resolves resolved")
