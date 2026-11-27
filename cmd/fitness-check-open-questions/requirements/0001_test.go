// Package requirements proves the open-questions check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
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
	t.Run("0001.1.1 a heading or bold or label line or HTML heading fails naming the file with its line", func(t *testing.T) {
		t.Parallel()
		doc := "# Sync\n\n## Open Questions\n\n- Which wins?\n\n**unresolved questions**\n\n- When?\n"
		out, code := openQuestions(t, map[string]string{"docs/sync.md": doc})
		doc = "# Sync\n\nQuestions for the owner: which store?\n\n- **Open questions:** none yet\n"
		labelOut, labelCode := openQuestions(t, map[string]string{"docs/sync.md": doc})
		page := "<!doctype html>\n<main>\n  <h1>Sync</h1>\n  <h2 class=\"x\">Open <em>questions</em></h2>\n  <ul><li>Which wins?</li></ul>\n</main>\n"
		pageOut, pageCode := openQuestions(t, map[string]string{"site/index.html": page})
		req.Then(t, "0001.1.2.1.1", "176ff4d0", func(_ req.A) {
			acceptance.Sees(t, out, code, 1,
				`docs/sync.md:3: a heading opens a list of open questions: "## Open Questions"`,
				`docs/sync.md:7: a bold line opens a list of open questions: "**unresolved questions**"`)
			acceptance.Sees(t, labelOut, labelCode, 1,
				`docs/sync.md:3: a label line opens a list of open questions: "Questions for the owner: which store?"`,
				`docs/sync.md:5: a label line opens a list of open questions: "- **Open questions:** none yet"`)
			acceptance.Sees(t, pageOut, pageCode, 1, `site/index.html:4: a heading opens a list of open questions: "Open questions"`)
		})
		req.Then(t, "0001.1.2.1.2", "e1dd5f6b", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, `"## Open Questions" — `+askInChat)
			acceptance.Sees(t, pageOut, pageCode, 1, `"Open questions" — `+askInChat)
		})
		out, code = openQuestions(t, map[string]string{"docs/sync.md": "# Sync\n\n## Decisions\n\n- The server wins.\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = openQuestions(t, map[string]string{"site/index.html": "<main>\n  <h1>Sync</h1>\n  <p>We have no open questions.</p>\n  <pre>Open questions:</pre>\n</main>\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2 open questions in a sentence or code or a comment pass", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "07b8e753", func(_ req.A) {
			{
				doc := "# Sync\n\nOpen questions go to chat, and we have no open questions left.\n\n**No open questions remain**\n"
				out, code := openQuestions(t, map[string]string{"docs/sync.md": doc})
				acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			}
			{
				doc := "# Sync\n\nUse `## Open questions` nowhere.\n\n```\nOpen questions:\n```\n\n<!-- Open questions: none -->\n"
				out, code := openQuestions(t, map[string]string{"docs/sync.md": doc})
				acceptance.Sees(t, out, code, 0, "All 1 checks passed")
				out, code = openQuestions(t, map[string]string{"docs/sync.md": doc + "\n" + questionsLabel})
				acceptance.Sees(t, out, code, 1, "docs/sync.md:11: "+questionsFinding)
			}
		})
	})
	t.Run("0001.1.3 a body label line fails naming the body", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.3.1", "dd1c247e", func(_ req.A) {
			repo := acceptance.Example(t, "happyRepo", map[string]string{"body.md": "> Pitch.\n\n" + questionsLabel})
			out, code := acceptance.Fitness(t, repo, nil, "open-questions", "--body-file", "body.md")
			acceptance.Sees(t, out, code, 1, questionsFinding)
		})
	})
}

// openQuestions runs `fitness-install -- open-questions` on happyRepo with
// files written over it.
func openQuestions(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "open-questions")
}

// askInChat is the fix every open-questions finding names.
const askInChat = "ask the owner in chat, or turn each unknown into a checklist step"

// questionsLabel is a label line that opens a list of open questions, and
// the finding it earns.
const (
	questionsLabel   = "Open questions:\n\n- Which sync wins?\n"
	questionsFinding = `a label line opens a list of open questions: "Open questions:"`
)
