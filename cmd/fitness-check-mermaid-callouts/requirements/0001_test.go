// Package requirements proves the mermaid-callouts check's docs/requirements
// as a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// callouts runs the check on happyRepo with files written over it.
func callouts(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, append([]string{"mermaid-callouts"}, args...)...)
}

// diagram fences mermaid lines as a diagram block.
func diagram(lines ...string) string {
	return "```mermaid\n" + strings.Join(lines, "\n") + "\n```\n"
}

// table builds a callout table with one row per label.
func table(labels ...string) string {
	out := "| # | Description |\n| --- | --- |\n"
	for _, n := range labels {
		out += "| " + n + " | step " + n + " |\n"
	}
	return out
}

// doc joins diagram and table blocks into docs/flow.md.
func doc(blocks ...string) map[string]string {
	return map[string]string{"docs/flow.md": strings.Join(blocks, "\n")}
}

// flow reads docs/flow.md from repo.
func flow(t *testing.T, repo string) string {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(repo, "docs", "flow.md"))
	acceptance.MustDo(t, err)
	return string(got)
}

// numbered is a doc whose flowchart and C4 diagrams label connectors with
// numbers beside numbered nodes and a letter in use.
func numbered() map[string]string {
	return doc(
		diagram(
			"flowchart LR",
			`app[1 App] -->|2| api["3 API"]`,
			"api -->|A| db((4 Store))",
			"db -- 6 --> app",
		),
		table("1", "2", "3", "4", "A", "6"),
		diagram(`Container(web, "5", "UI")`, `Rel(web, api, "12", "HTTPS")`),
		"| # | Description |\n| --- | --- |\n| 5  | step 5 |\n| 12 | step 12 |\n",
	)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	valid := func(t *testing.T) (string, int) {
		return callouts(t, doc(
			diagram(
				"flowchart LR",
				`app[1 App] -->|A| api["2 API"]`,
				"api -- B --> db((3 Store))",
				`db -->|"C"| app`,
				"db -->|YES| app -->|go| Q --> db",
			),
			table("1", "2", "3", "A", "B", "C"),
			diagram(`Container(web, "4", "UI")`, `Rel(web, api, "D", "HTTPS")`),
			diagram(`Person(p, "Person")`, `System(s, "System")`, `Rel(p, s, "uses")`),
			table("4", "D"),
			diagram(
				"wardley-beta",
				"  anchor User [0.99, 0.66]",
				"  component code [0.60, 0.35]",
				`  annotation 5,[0.40, 0.35] "Genesis: code"`,
				"  User -> code",
			),
			table("5"),
			diagram("gitGraph", "  commit", "  branch dev", "  commit"),
		))
	}
	t.Run("0001.1.2.1.1 the file passes", func(t *testing.T) {
		t.Parallel()
		out, code := valid(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.1.2 the count includes the file", func(t *testing.T) {
		t.Parallel()
		out, code := valid(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if !regexp.MustCompile(`│ passed\s+│ 1\s+│`).MatchString(acceptance.Row(out, "mermaid-callouts")) {
			t.Errorf("mermaid-callouts must pass judging 1 file, got row %q", acceptance.Row(out, "mermaid-callouts"))
		}
	})
	t.Run("0001.1.2.2.1 it fails with `mermaid diagram has no numbered callouts`", func(t *testing.T) {
		t.Parallel()
		out, code := callouts(t, map[string]string{
			"docs/plain.md": strings.Join([]string{
				diagram("flowchart LR", "Login --> Home"),
				diagram("flowchart LR", "a[1 A] --> b[2 B]"),
				table("1", "2"),
				diagram("flowchart TD", "Fetch -->|then parse| Save"),
			}, "\n"),
		})
		acceptance.Sees(t, out, code, 1,
			"docs/plain.md: line 1: mermaid diagram has no numbered callouts; number its nodes, as in A[1 Step]",
			"docs/plain.md: line 16: mermaid diagram has no numbered callouts")
	})
	unpaired := func(t *testing.T) (string, int) {
		return callouts(t, map[string]string{
			"docs/flow.md": strings.Join([]string{
				diagram("classDiagram", `class a["1 A"]`, `class b["1 B"]`, "a --> b : A", "b --> a : C"),
				table("1", "A", "A", "B"),
			}, "\n"),
			"docs/unpaired.md": strings.Join([]string{
				table("A"),
				diagram(`Rel(a, b, "A")`),
			}, "\n"),
		})
	}
	t.Run("0001.1.2.3.1 a callout with no row fails naming it with its line", func(t *testing.T) {
		t.Parallel()
		out, code := unpaired(t)
		acceptance.Sees(t, out, code, 1, "docs/flow.md: diagram callout C has no matching table row (line 1)")
	})
	t.Run("0001.1.2.3.2 a row with no callout fails naming it with its line", func(t *testing.T) {
		t.Parallel()
		out, code := unpaired(t)
		acceptance.Sees(t, out, code, 1, "docs/flow.md: callout table row B has no matching diagram callout (line 9)")
	})
	t.Run("0001.1.2.3.3 a callout or row that appears twice fails the same way", func(t *testing.T) {
		t.Parallel()
		out, code := unpaired(t)
		acceptance.Sees(t, out, code, 1,
			"docs/flow.md: diagram callout 1 appears 2 times (line 1)",
			"docs/flow.md: callout table row A appears 2 times (line 9)")
	})
	t.Run("0001.1.2.3.4 a callout with no table or a row with no diagram fails the same way", func(t *testing.T) {
		t.Parallel()
		out, code := unpaired(t)
		acceptance.Sees(t, out, code, 1,
			"docs/unpaired.md: callout table at line 1 has no preceding mermaid diagram",
			"docs/unpaired.md: mermaid diagram at line 5 has callouts but no associated callout table")
	})
	numberedConnector := func(t *testing.T) (string, int) {
		return callouts(t, doc(
			diagram(`Container(app, "1", "UI")`, `Rel(app, api, "2")`),
			table("1", "2"),
		))
	}
	t.Run("0001.1.2.4.1 it fails asking for a letter", func(t *testing.T) {
		t.Parallel()
		out, code := numberedConnector(t)
		acceptance.Sees(t, out, code, 1, "docs/flow.md: connector callout 2 must be a letter; numbers are for nodes")
		if strings.Contains(out, "has no matching") {
			t.Errorf("a numbered connector must fail only for its form, not as an orphan:\n%s", out)
		}
	})
	t.Run("0001.1.2.4.2 the failure names `--fix`", func(t *testing.T) {
		t.Parallel()
		out, code := numberedConnector(t)
		acceptance.Sees(t, out, code, 1, "run `fitness mermaid-callouts --fix` (line 1)")
	})
	doubling := func(t *testing.T) (string, int) {
		return callouts(t, doc(
			diagram("flowchart LR", "a -->|Z| b", "b -->|AA| c", "c -->|AB| a"),
			table("Z", "AA", "AB"),
		))
	}
	t.Run("0001.1.2.5.1 it fails for the diagram callout", func(t *testing.T) {
		t.Parallel()
		out, code := doubling(t)
		acceptance.Sees(t, out, code, 1,
			"docs/flow.md: connector callout AB breaks the letter doubling; use A to Z, then AA to ZZ, then AAA (line 1)")
		if strings.Contains(out, "callout AA ") {
			t.Errorf("AA follows the doubling and must pass:\n%s", out)
		}
	})
	t.Run("0001.1.2.5.2 it fails for the table row", func(t *testing.T) {
		t.Parallel()
		out, code := doubling(t)
		acceptance.Sees(t, out, code, 1,
			"docs/flow.md: callout table row AB breaks the letter doubling; use A to Z, then AA to ZZ, then AAA (line 8)")
		if strings.Contains(out, "row AA ") {
			t.Errorf("AA follows the doubling and must pass:\n%s", out)
		}
	})
	t.Run("0001.1.2.6.1 it fails naming `(description)`", func(t *testing.T) {
		t.Parallel()
		body := doc(diagram(`Rel(a, b, "A")`, `Rel(b, c, "B")`), table("A"))["docs/flow.md"]
		out, code := callouts(t, map[string]string{"body.md": body}, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, "(description): diagram callout B has no matching table row (line 1)")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	fixed := func(t *testing.T) (string, string, int) {
		repo := acceptance.Example(t, "happyRepo", numbered())
		out, code := acceptance.Fitness(t, repo, nil, "mermaid-callouts", "--fix")
		return repo, out, code
	}
	t.Run("0001.2.2.1.1 each becomes the next free letter", func(t *testing.T) {
		t.Parallel()
		repo, out, code := fixed(t)
		acceptance.Sees(t, out, code, 0, "docs/flow.md: converted connector callouts 2 to B, 6 to C, 12 to A", "--fix changed 1 file")
		got := flow(t, repo)
		for _, line := range []string{`app[1 App] -->|B| api["3 API"]`, "api -->|A| db((4 Store))", "db -- C --> app", `Rel(web, api, "A", "HTTPS")`} {
			if !strings.Contains(got, line) {
				t.Errorf("--fix must turn each numbered connector into the next free letter; missing %q in:\n%s", line, got)
			}
		}
	})
	t.Run("0001.2.2.1.2 its row follows", func(t *testing.T) {
		t.Parallel()
		repo, out, code := fixed(t)
		acceptance.Sees(t, out, code, 0, "--fix changed 1 file")
		got := flow(t, repo)
		for _, row := range []string{"| B | step 2 |", "| C | step 6 |", "| A  | step 12 |"} {
			if !strings.Contains(got, row) {
				t.Errorf("--fix must turn each connector's row into its letter; missing %q in:\n%s", row, got)
			}
		}
	})
	t.Run("0001.2.2.1.3 the repo passes", func(t *testing.T) {
		t.Parallel()
		repo, _, _ := fixed(t)
		out, code := acceptance.Fitness(t, repo, nil, "mermaid-callouts")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2.2.1 it changes nothing", func(t *testing.T) {
		t.Parallel()
		repo, _, _ := fixed(t)
		once := flow(t, repo)
		out, code := acceptance.Fitness(t, repo, nil, "mermaid-callouts", "--fix")
		acceptance.Sees(t, out, code, 0, "--fix changed 0 files")
		if strings.Contains(out, "converted") {
			t.Errorf("a second --fix must list no file, got:\n%s", out)
		}
		if got := flow(t, repo); got != once {
			t.Errorf("a second --fix must change nothing, got:\n%s\nwant:\n%s", got, once)
		}
	})
	shared := func(t *testing.T) (string, string, int, string) {
		files := doc(
			diagram(`Container(app, "2", "UI")`, `Rel(app, api, "2")`),
			table("2"),
		)
		repo := acceptance.Example(t, "happyRepo", files)
		out, code := acceptance.Fitness(t, repo, nil, "mermaid-callouts", "--fix")
		return repo, out, code, files["docs/flow.md"]
	}
	t.Run("0001.2.2.3.1 `--fix` leaves the diagram unchanged", func(t *testing.T) {
		t.Parallel()
		repo, out, code, original := shared(t)
		acceptance.Sees(t, out, code, 1,
			"docs/flow.md: --fix left the mermaid diagram at line 1 unchanged: callout 2 marks a node too",
			"--fix changed 0 files")
		if got := flow(t, repo); got != original {
			t.Errorf("--fix must leave a diagram it cannot convert byte for byte, got:\n%s", got)
		}
	})
	t.Run("0001.2.2.3.2 the check still fails", func(t *testing.T) {
		t.Parallel()
		_, out, code, _ := shared(t)
		acceptance.Sees(t, out, code, 1, "docs/flow.md: connector callout 2 must be a letter")
	})
	noDiagram := func(t *testing.T) (string, int) { return callouts(t, nil) }
	t.Run("0001.2.2.4.1 it is skipped", func(t *testing.T) {
		t.Parallel()
		out, code := noDiagram(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2.4.2 it is not counted", func(t *testing.T) {
		t.Parallel()
		out, _ := noDiagram(t)
		acceptance.PassedWithNoFiles(t, out, "mermaid-callouts")
	})
}
