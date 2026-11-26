// Package requirements proves the mermaid-syntax check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

// cspell:ignore Descr mindmap xychart

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// syntax runs the check on happyRepo with files written over it.
func syntax(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, append([]string{"mermaid-syntax"}, args...)...)
}

// The shared mermaid fixtures from acceptance.
var (
	diagram = acceptance.Diagram
	doc     = acceptance.DiagramDoc
)

// verifyRejects fails when the README holds one diagram of the given lines.
func verifyRejects(t *testing.T, lines []string, want ...string) {
	t.Helper()
	out, code := syntax(t, doc(diagram(lines...)))
	acceptance.Sees(t, out, code, 1, want...)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	verifyPasses := func(t *testing.T, known ...[]string) {
		var diagrams []string
		for _, lines := range known {
			diagrams = append(diagrams, diagram(lines...))
		}
		out, code := syntax(t, doc(diagrams...))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	}
	t.Run("0001.1.1.1.1 the flowchart passes", func(t *testing.T) {
		t.Parallel()
		verifyPasses(t,
			[]string{"flowchart LR", `A["Start (now)"]:::s --> B{Ready?} -->|yes| C((Done))`, "B -.-> D & E",
				"subgraph S [Group]", "direction TB", "F -- note --> G", "end", "classDef s fill:#eef"},
			[]string{"graph TD;", "A-->B;"},
			[]string{"flowchart-elk TD", "accDescr {", "free text: here", "}", "A-- yes -->B"})
	})
	t.Run("0001.1.1.1.2 the sequence diagram passes", func(t *testing.T) {
		t.Parallel()
		verifyPasses(t,
			[]string{"sequenceDiagram", "participant A as Client", "loop every minute", "A->>+B: ping", "alt ok",
				"B-->>-A: pong", "else down", "Note over A,B: retry", "end", "end"},
			[]string{"sequenceDiagram", "par_over A", "A->>B: x", "and B", "B->>A: y", "end"})
	})
	t.Run("0001.1.1.1.3 the C4 diagram passes", func(t *testing.T) {
		t.Parallel()
		verifyPasses(t, []string{"C4Context", `Person(user, "1 User")`, `Rel(user, app, "2")`})
	})
	t.Run("well formed diagrams of the other known types pass", func(t *testing.T) {
		t.Parallel()
		verifyPasses(t, []string{"mindmap", "root", "  a)cloud("}, []string{"xychart", "bar [1, 2]"})
	})
	openEdge := []string{"flowchart TD", "A -->", "subgraph S", "B --> C"}
	t.Run("0001.1.1.2.1 it fails naming each line", func(t *testing.T) {
		t.Parallel()
		verifyRejects(t, openEdge,
			"README.md: line 11: an edge has no target node",
			"README.md: line 12: subgraph is never closed with end")
	})
	t.Run("0001.1.1.2.2 it names the diagram line", func(t *testing.T) {
		t.Parallel()
		verifyRejects(t, openEdge, "(diagram at line 9)")
	})
	t.Run("0001.1.1.3.1 it fails naming each line", func(t *testing.T) {
		t.Parallel()
		verifyRejects(t, []string{"sequenceDiagram", "A->>B", "loop every minute", "A->>B: hi"},
			`README.md: line 11: message "A->>B" needs a target and a colon`,
			"README.md: line 12: loop is never closed with end (diagram at line 9)")
	})
	t.Run("0001.1.1.4.1 it fails naming the type", func(t *testing.T) {
		t.Parallel()
		verifyRejects(t, []string{"graf TD", "A --> B"}, `README.md: line 10: unknown diagram type "graf" (diagram at line 9)`)
	})
	t.Run("0001.1.1.5.1 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		verifyRejects(t, []string{"C4Context", `Person(user, "1 User")`, `Rel(user, app, "2"`},
			"README.md: line 12: '(' is never closed (diagram at line 9)")
	})
	t.Run("0001.1.1.6.1 it fails asking for double quotes", func(t *testing.T) {
		t.Parallel()
		verifyRejects(t, []string{"flowchart TD", "A[Start (x)] --> B"},
			`README.md: line 11: node A: label "Start (x)" has brackets or quotes; wrap it in double quotes (diagram at line 9)`)
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.2.1.1 it fails naming `(description)`", func(t *testing.T) {
		t.Parallel()
		out, code := syntax(t, map[string]string{"body.md": diagram("flowchart TD", "A[Start (x)] --> B")},
			"--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, `(description): line 3: node A: label "Start (x)"`)
	})
	t.Run("0001.2.2.2.1 the check passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := syntax(t, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "mermaid-syntax")
	})
	judged := func(t *testing.T) (string, int, string) {
		bad := diagram("flowchart TD", "A[Start (x)] --> B")
		files := doc(bad, "> ```mermaid\n> flowchart TD\n> A[Start (x)] --> B\n> ```\n")
		files["docs/guide.markdown"] = "# Guide\n\n" + bad
		files["docs/flow.mmd"] = "flowchart TD\nA[Start (x)] --> B\n"
		files["docs/reply.md"] = "> ```mermaid\n> flowchart TD\n> A --> B\n\nAfter the quote.\n\n" + bad
		files["docs/how-to.md"] = "# How\n\nQuote a diagram like this:\n\n````markdown\n> ```mermaid\n> flowchart TD\n> A[Start (x)] --> B\n> ```\n````\n"
		out, code := syntax(t, files)
		return out, code, `node A: label "Start (x)" has brackets or quotes; wrap it in double quotes`
	}
	t.Run("0001.2.2.3.1 each real diagram is judged", func(t *testing.T) {
		t.Parallel()
		out, code, label := judged(t)
		acceptance.Sees(t, out, code, 1, "README.md: line 11: "+label+" (diagram at line 9)",
			"README.md: line 16: "+label+" (diagram at line 14)",
			"docs/guide.markdown: line 5: "+label+" (diagram at line 3)",
			"docs/flow.mmd: line 2: "+label,
			"docs/reply.md: line 9: "+label+" (diagram at line 7)")
	})
	t.Run("0001.2.2.3.2 a code block example is not judged", func(t *testing.T) {
		t.Parallel()
		out, code, label := judged(t)
		acceptance.Sees(t, out, code, 1)
		for _, unwanted := range []string{"docs/flow.mmd: line 2: " + label + " (diagram", "docs/how-to.md", "docs/reply.md: line 5"} {
			if strings.Contains(acceptance.Read(out), unwanted) {
				t.Errorf("a whole-file diagram names no fence, and a quoted example in a code block is code:\n%s", out)
			}
		}
	})
}
