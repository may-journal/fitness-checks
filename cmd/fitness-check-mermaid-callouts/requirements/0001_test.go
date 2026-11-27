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
	"github.com/may-journal/fitness-checks/req"
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
	t.Run("0001.1.1 callouts mapping to rows in both directions pass in every diagram kind", func(t *testing.T) {
		t.Parallel()
		out, code := callouts(t, doc(
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
		req.Then(t, "0001.1.2.1.1", "dbfdaa5e", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.1.2", "1a4969c2", func(a req.A) {
			a.True(regexp.MustCompile(`│ passed\s+│ 1\s+│`).MatchString(acceptance.Row(out, "mermaid-callouts")), "mermaid-callouts must pass judging 1 file, got row %q", acceptance.Row(out, "mermaid-callouts"))
		})
	})
	t.Run("0001.1.2 a diagram with no numbered callouts fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "96b80c92", func(_ req.A) {
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
	})
	t.Run("0001.1.3 a callout that does not pair with its row fails", func(t *testing.T) {
		t.Parallel()
		out, code := callouts(t, map[string]string{
			"docs/flow.md": strings.Join([]string{
				diagram("classDiagram", `class a["1 A"]`, `class b["1 B"]`, "a --> b : A", "b --> a : C"),
				table("1", "A", "A", "B"),
			}, "\n"),
			"docs/unpaired.md": strings.Join([]string{
				table("A"),
				diagram(`Rel(a, b, "A")`),
			}, "\n"),
		})
		req.Then(t, "0001.1.2.3.1", "fffd5dbe", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/flow.md: diagram callout C has no matching table row (line 1)")
		})
		req.Then(t, "0001.1.2.3.2", "c6ae1eb6", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/flow.md: callout table row B has no matching diagram callout (line 9)")
		})
		req.Then(t, "0001.1.2.3.3", "f78ea230", func(a req.A) {
			acceptance.Sees(t, out, code, 1,
				"docs/flow.md: diagram callout 1 appears 2 times (line 1)",
				"docs/flow.md: callout table row A appears 2 times (line 9)")
		})
		req.Then(t, "0001.1.2.3.4", "9791176a", func(a req.A) {
			acceptance.Sees(t, out, code, 1,
				"docs/unpaired.md: callout table at line 1 has no preceding mermaid diagram",
				"docs/unpaired.md: mermaid diagram at line 5 has callouts but no associated callout table")
		})
	})
	t.Run("0001.1.4 a numbered connector asks for a letter", func(t *testing.T) {
		t.Parallel()
		out, code := callouts(t, doc(
			diagram(`Container(app, "1", "UI")`, `Rel(app, api, "2")`),
			table("1", "2"),
		))
		req.Then(t, "0001.1.2.4.1", "723f578a", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/flow.md: connector callout 2 must be a letter; numbers are for nodes")
			a.True(!strings.Contains(out, "has no matching"), "a numbered connector must fail only for its form, not as an orphan:\n%s", out)
		})
		req.Then(t, "0001.1.2.4.2", "e11719d8", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "run `fitness mermaid-callouts --fix` (line 1)")
		})
	})
	t.Run("0001.1.5 a connector letter breaking the doubling fails", func(t *testing.T) {
		t.Parallel()
		out, code := callouts(t, doc(
			diagram("flowchart LR", "a -->|Z| b", "b -->|AA| c", "c -->|AB| a"),
			table("Z", "AA", "AB"),
		))
		req.Then(t, "0001.1.2.5.1", "c20ccfb7", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/flow.md: connector callout AB breaks the letter doubling; use A to Z, then AA to ZZ, then AAA (line 1)")
			a.True(!strings.Contains(out, "callout AA "), "AA follows the doubling and must pass:\n%s", out)
		})
		req.Then(t, "0001.1.2.5.2", "4f891d85", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/flow.md: callout table row AB breaks the letter doubling; use A to Z, then AA to ZZ, then AAA (line 8)")
			a.True(!strings.Contains(out, "row AA "), "AA follows the doubling and must pass:\n%s", out)
		})
	})
	t.Run("0001.1.6 a PR body callout with no row fails naming the description", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "41f704b4", func(_ req.A) {
			body := doc(diagram(`Rel(a, b, "A")`, `Rel(b, c, "B")`), table("A"))["docs/flow.md"]
			out, code := callouts(t, map[string]string{"body.md": body}, "--body-file", "body.md")
			acceptance.Sees(t, out, code, 1, "(description): diagram callout B has no matching table row (line 1)")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 fix gives connectors the next free letter", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", numbered())
		out, code := acceptance.Fitness(t, repo, nil, "mermaid-callouts")
		acceptance.Sees(t, out, code, 1, "connector callout 2 must be a letter")
		fixOut, fixCode := acceptance.Fitness(t, repo, nil, "mermaid-callouts", "--fix")
		want := doc(
			diagram(
				"flowchart LR",
				`app[1 App] -->|B| api["3 API"]`,
				"api -->|A| db((4 Store))",
				"db -- C --> app",
			),
			"| # | Description |\n| --- | --- |\n| 1 | step 1 |\n| B | step 2 |\n| 3 | step 3 |\n| 4 | step 4 |\n| A | step A |\n| C | step 6 |\n",
			diagram(`Container(web, "5", "UI")`, `Rel(web, api, "A", "HTTPS")`),
			"| # | Description |\n| --- | --- |\n| 5  | step 5 |\n| A  | step 12 |\n",
		)["docs/flow.md"]
		got := flow(t, repo)
		again, againCode := acceptance.Fitness(t, repo, nil, "mermaid-callouts")
		req.Then(t, "0001.2.2.1.1", "a06a2531", func(a req.A) {
			acceptance.Sees(t, fixOut, fixCode, 0, "docs/flow.md: converted connector callouts 2 to B, 6 to C, 12 to A", "--fix changed 1 file")
		})
		req.Then(t, "0001.2.2.1.2", "34e48052", func(a req.A) {
			a.Equal(want, got)
		})
		req.Then(t, "0001.2.2.1.3", "20e5b014", func(a req.A) {
			acceptance.Sees(t, again, againCode, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.2.2 a second fix run changes nothing", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.2.2.1", "ba4676f2", func(_ req.A) {
			repo := acceptance.Example(t, "happyRepo", numbered())
			out, code := acceptance.Fitness(t, repo, nil, "mermaid-callouts", "--fix")
			acceptance.Sees(t, out, code, 0, "--fix changed 1 file")
			fixed := flow(t, repo)
			out, code = acceptance.Fitness(t, repo, nil, "mermaid-callouts", "--fix")
			acceptance.Sees(t, out, code, 0, "--fix changed 0 files")
			if strings.Contains(out, "converted") {
				t.Errorf("a second --fix must list no file, got:\n%s", out)
			}
			if got := flow(t, repo); got != fixed {
				t.Errorf("a second --fix must change nothing, got:\n%s\nwant:\n%s", got, fixed)
			}
		})
	})
	t.Run("0001.2.3 a numbered connector sharing a node number stays unfixed", func(t *testing.T) {
		t.Parallel()
		files := doc(
			diagram(`Container(app, "2", "UI")`, `Rel(app, api, "2")`),
			table("2"),
		)
		repo := acceptance.Example(t, "happyRepo", files)
		out, code := acceptance.Fitness(t, repo, nil, "mermaid-callouts", "--fix")
		req.Then(t, "0001.2.2.3.1", "9f34c72c", func(a req.A) {
			acceptance.Sees(t, out, code, 1,
				"docs/flow.md: --fix left the mermaid diagram at line 1 unchanged: callout 2 marks a node too",
				"--fix changed 0 files")
			a.Equal(files["docs/flow.md"], flow(t, repo))
		})
		req.Then(t, "0001.2.2.3.2", "00eed45a", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/flow.md: connector callout 2 must be a letter")
		})
	})
	t.Run("0001.2.4 a file with no diagram and no table is skipped", func(t *testing.T) {
		t.Parallel()
		out, code := callouts(t, nil)
		req.Then(t, "0001.2.2.4.1", "cee2af4c", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.2.2.4.2", "218895d7", func(a req.A) {
			acceptance.PassedWithNoFiles(t, out, "mermaid-callouts")
		})
	})
}
