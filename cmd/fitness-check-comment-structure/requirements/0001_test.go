// Package requirements proves the comment-structure check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo with the comment saved over it.
package requirements

import (
	"os"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1.1.1 it fails naming the count", func(t *testing.T) {
		t.Parallel()
		out, code := commentStructure(t, overBudget, nil)
		acceptance.Sees(t, out, code, 1, "(comment): 120 words outside tables, code, and diagrams")
	})
	t.Run("0001.1.1.1.2 it fails naming the budget", func(t *testing.T) {
		t.Parallel()
		out, code := commentStructure(t, overBudget, nil)
		acceptance.Sees(t, out, code, 1, "the budget is 100")
	})
	t.Run("0001.1.1.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := commentStructure(t, overBudget, map[string]string{".fitnessrc.json": `{"comment": {"maxWords": 150}}`})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	tables := func(t *testing.T) (string, int) {
		table := "| A | B |\n| - | - |\n| 1 | 2 |\n"
		return commentStructure(t, table+"\nAnd more.\n\n"+table+"\n"+numbered, nil)
	}
	t.Run("0001.1.1.3.1 it fails naming the count", func(t *testing.T) {
		t.Parallel()
		out, code := tables(t)
		acceptance.Sees(t, out, code, 1, "(comment): 2 tables besides callout tables")
	})
	t.Run("0001.1.1.3.2 it fails naming the budget", func(t *testing.T) {
		t.Parallel()
		out, code := tables(t)
		acceptance.Sees(t, out, code, 1, "the budget is 1")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1.1.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := commentStructure(t, "Flow.\n\n```mermaid\nflowchart TD\n  A --> B\n```\n", nil)
		acceptance.Sees(t, out, code, 1, "(comment): line 3: mermaid diagram has no numbered callouts")
	})
	t.Run("0001.2.1.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		legend := strings.Replace(numbered, "\n\n|", "\n\n```mermaid\nflowchart LR\n  k[Step]:::s\n  classDef s fill:#eef\n```\n\n|", 1)
		out, code := commentStructure(t, legend, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	note := func(t *testing.T) (string, int) {
		return commentStructure(t, "The parser works.\n\nEdit 2: it handles tabs too.\n", nil)
	}
	t.Run("0001.3.1.1.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := note(t)
		acceptance.Sees(t, out, code, 1, `(comment): line 3:`)
		out, code = commentStructure(t, "The cache is updated after every push.\n", nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.3.1.1.2 it fails naming the note", func(t *testing.T) {
		t.Parallel()
		out, code := note(t)
		acceptance.Sees(t, out, code, 1, `"Edit 2:": change-log note; edit the text in place and drop the note`)
	})
	frame := func(t *testing.T) (string, int) {
		return commentStructure(t, "Done.\n\n<system-reminder>\n", nil)
	}
	t.Run("0001.3.1.2.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := frame(t)
		acceptance.Sees(t, out, code, 1, `(comment): line 3:`)
	})
	t.Run("0001.3.1.2.2 it fails naming the frame", func(t *testing.T) {
		t.Parallel()
		out, code := frame(t)
		acceptance.Sees(t, out, code, 1, `"<system-reminder": agent frame; post only the finished text`)
	})
	model := func(t *testing.T) (string, int) {
		return commentStructure(t, "Done.\n\nRan on claude-opus-4-1 today.\n", nil)
	}
	t.Run("0001.3.1.3.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := model(t)
		acceptance.Sees(t, out, code, 1, `(comment): line 3:`)
	})
	t.Run("0001.3.1.3.2 it fails naming the identifier", func(t *testing.T) {
		t.Parallel()
		out, code := model(t)
		acceptance.Sees(t, out, code, 1, `"claude-opus-4-1": model identifier; name the agent, not the model`)
	})
	marker := func(t *testing.T) (string, int) {
		return commentStructure(t, "Done.\n\n<!-- fitness:comment-check:31:abc -->\n", nil)
	}
	t.Run("0001.3.1.4.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := marker(t)
		acceptance.Sees(t, out, code, 1, `(comment): line 3:`)
	})
	t.Run("0001.3.1.4.2 it fails naming the marker", func(t *testing.T) {
		t.Parallel()
		out, code := marker(t)
		acceptance.Sees(t, out, code, 1, `"<!-- fitness:": carries a fitness marker; only fitness writes those, so remove it`)
	})
	t.Run("0001.3.1.5.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := commentStructure(t, "Pinned to `claude-opus-4-1` in config.\n", nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_4(t *testing.T) {
	t.Parallel()
	t.Run("0001.4.1.1.1 the check is skipped", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", nil)
		for _, body := range []string{planBody("- [ ] Write the tests."), epicBody("1. #22\n2. #23, after #22")} {
			out, code := acceptance.Piped(t, repo, body, nil, "--blob", "-")
			acceptance.Sees(t, out, code, 0, "Skipped (does not apply to this text): comment-structure")
		}
	})
	t.Run("0001.4.1.2.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "comment-structure")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "comment-structure")
	})
}

// commentStructure runs comment-structure on comment, saved over happyRepo
// with files, as --body-file.
func commentStructure(t *testing.T, comment string, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.With(files, map[string]string{"comment.md": comment})), nil,
		"comment-structure", "--body-file", "comment.md")
}

// words is n words of plain prose.
func words(n int) string {
	return strings.TrimSpace(strings.Repeat("word ", n)) + "."
}

// numbered is a numbered diagram and its matching callout table.
const numbered = "```mermaid\nflowchart TD\n  A[1 Start]:::s --> B[2 End]:::s\n  classDef s fill:#eef\n```\n\n" +
	"| # | Element | Description | Why |\n| - | - | - | - |\n| 1 | Start | First. | Begins. |\n| 2 | End | Last. | Ends. |\n"

// overBudget has 120 prose words, and more inside a table and code.
var overBudget = words(120) + "\n\n| A | B |\n| - | - |\n| " + words(30) + " | x |\n\n```\n" + words(30) + "\n```\n"

// planBody is a Plan that follows the template, with one task.
func planBody(task string) string {
	return "> One clear pitch for this Plan.\n\n## Background\n\nThis Plan is a test fixture.\n\n## What needs to happen\n\n" + task + "\n"
}

// epicBody is an Epic that follows the template, listing items.
func epicBody(items string) string {
	return "> Every fitness run finishes in a fraction of the time.\n\n## Background\n\nThis epic is a test fixture.\n\n## What needs to happen\n\n" + items + "\n"
}
