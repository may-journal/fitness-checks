package requirements

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func Test0005_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"valid diagrams of every type pass", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidSyntax(t, mermaidLegendDoc(
				mermaidLegendDiagram("flowchart LR", `A["Start (now)"]:::s --> B{Ready?} -->|yes| C((Done))`,
					"B -.-> D & E", "subgraph S [Group]", "direction TB", "F -- note --> G", "end", "classDef s fill:#eef"),
				mermaidLegendDiagram("graph TD;", "A-->B;"),
				mermaidLegendDiagram("sequenceDiagram", "participant A as Client", "loop every minute",
					"A->>+B: ping", "alt ok", "B-->>-A: pong", "else down", "Note over A,B: retry", "end", "end"),
				mermaidLegendDiagram("C4Context", `Person(user, "1 User")`, `Rel(user, app, "2")`),
				mermaidLegendDiagram("flowchart-elk TD", "accDescr {", "free text: here", "}", "A-- yes -->B"),
				mermaidLegendDiagram("sequenceDiagram", "par_over A", "A->>B: x", "and B", "B->>A: y", "end"),
				mermaidLegendDiagram("mindmap", "root", "  a)cloud("),
				mermaidLegendDiagram("xychart", "bar [1, 2]")))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"flowchart with open edge and subgraph fails", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidSyntax(t, mermaidLegendDoc(mermaidLegendDiagram("flowchart TD", "A -->", "subgraph S", "B --> C")))
			sees(t, out, code, 1, "README.md: line 11: an edge has no target node (diagram at line 9)",
				"README.md: line 12: subgraph is never closed with end (diagram at line 9)")
		}},
		{"sequence diagram with open message and loop fails", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidSyntax(t, mermaidLegendDoc(mermaidLegendDiagram("sequenceDiagram", "A->>B", "loop every minute", "A->>B: hi")))
			sees(t, out, code, 1, `README.md: line 11: message "A->>B" needs a target and a colon`,
				"README.md: line 12: loop is never closed with end (diagram at line 9)")
		}},
		{"unknown diagram type fails by name", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidSyntax(t, mermaidLegendDoc(mermaidLegendDiagram("graf TD", "A --> B")))
			sees(t, out, code, 1, `README.md: line 10: unknown diagram type "graf" (diagram at line 9)`)
		}},
		{"c4 diagram with unclosed bracket fails", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidSyntax(t, mermaidLegendDoc(mermaidLegendDiagram("C4Context", `Person(user, "1 User")`, `Rel(user, app, "2"`)))
			sees(t, out, code, 1, "README.md: line 12: '(' is never closed (diagram at line 9)")
		}},
		{"pr body with broken flowchart fails", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidSyntax(t, map[string]string{"body.md": mermaidLegendDiagram("flowchart TD", "A[Start (x)] --> B")},
				"--body-file", "body.md")
			sees(t, out, code, 1, `(description): line 3: node A: label "Start (x)"`)
		}},
		{"docs with no diagrams pass unchecked", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidSyntax(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "mermaid-syntax")
		}},
		// The case puts one bad label in the README on line 11 and quoted on line
		// 16. Two more files hold it too.
		{"bare brackets and odd files fail on real nodes", func(t *testing.T) {
			t.Parallel()
			bad := mermaidLegendDiagram("flowchart TD", "A[Start (x)] --> B")
			files := mermaidLegendDoc(bad, "> ```mermaid\n> flowchart TD\n> A[Start (x)] --> B\n> ```\n")
			files["docs/guide.markdown"] = "# Guide\n\n" + bad
			files["docs/flow.mmd"] = "flowchart TD\nA[Start (x)] --> B\n"
			files["docs/reply.md"] = "> ```mermaid\n> flowchart TD\n> A --> B\n\nAfter the quote.\n\n" + bad
			files["docs/how-to.md"] = "# How\n\nQuote a diagram like this:\n\n````markdown\n> ```mermaid\n> flowchart TD\n> A[Start (x)] --> B\n> ```\n````\n"
			out, code := mermaidSyntax(t, files)
			label := `node A: label "Start (x)" has brackets or quotes; wrap it in double quotes`
			sees(t, out, code, 1, "README.md: line 11: "+label+" (diagram at line 9)",
				"README.md: line 16: "+label+" (diagram at line 14)",
				"docs/guide.markdown: line 5: "+label+" (diagram at line 3)",
				"docs/flow.mmd: line 2: "+label,
				"docs/reply.md: line 9: "+label+" (diagram at line 7)")
			for _, unwanted := range []string{"docs/flow.mmd: line 2: " + label + " (diagram", "docs/how-to.md", "docs/reply.md: line 5"} {
				if strings.Contains(read(out), unwanted) {
					t.Errorf("a whole-file diagram names no fence, and a quoted example in a code block is code:\n%s", out)
				}
			}
		}},
	})
}

func Test0005_2(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"matched callouts in every diagram kind pass", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidCallouts(t, mermaidCalloutsDoc(
				mermaidCalloutsDiagram(
					"flowchart LR",
					`app[1 App] -->|A| api["2 API"]`,
					"api -- B --> db((3 Store))",
					`db -->|"C"| app`,
					"db -->|YES| app -->|go| Q --> db",
				),
				mermaidCalloutsTable("1", "2", "3", "A", "B", "C"),
				mermaidCalloutsDiagram(`Container(web, "4", "UI")`, `Rel(web, api, "D", "HTTPS")`),
				mermaidCalloutsDiagram(`Person(p, "Person")`, `System(s, "System")`, `Rel(p, s, "uses")`),
				mermaidCalloutsTable("4", "D"),
				mermaidCalloutsDiagram(
					"wardley-beta",
					"  anchor User [0.99, 0.66]",
					"  component code [0.60, 0.35]",
					`  annotation 5,[0.40, 0.35] "Genesis: code"`,
					"  User -> code",
				),
				mermaidCalloutsTable("5"),
				mermaidCalloutsDiagram("gitGraph", "  commit", "  branch dev", "  commit"),
			))
			sees(t, out, code, 0, "All 1 checks passed")
			if !regexp.MustCompile(`│ passed\s+│ 1\s+│`).MatchString(row(out, "mermaid-callouts")) {
				t.Errorf("mermaid-callouts must pass judging 1 file, got row %q", row(out, "mermaid-callouts"))
			}
		}},
		{"a missing or repeated callout fails, and so does an unpaired one", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidCallouts(t, map[string]string{
				"docs/flow.md": strings.Join([]string{
					mermaidCalloutsDiagram("classDiagram", `class a["1 A"]`, `class b["1 B"]`, "a --> b : A", "b --> a : C"),
					mermaidCalloutsTable("1", "A", "A", "B"),
				}, "\n"),
				"docs/unpaired.md": strings.Join([]string{
					mermaidCalloutsTable("A"),
					mermaidCalloutsDiagram(`Rel(a, b, "A")`),
				}, "\n"),
				"docs/plain.md": strings.Join([]string{
					mermaidCalloutsDiagram("flowchart LR", "Login --> Home"),
					mermaidCalloutsDiagram("flowchart LR", "a[1 A] --> b[2 B]"),
					mermaidCalloutsTable("1", "2"),
					mermaidCalloutsDiagram("flowchart TD", "Fetch -->|then parse| Save"),
				}, "\n"),
			})
			sees(t, out, code, 1,
				"docs/plain.md: line 1: mermaid diagram has no numbered callouts; number its nodes, as in A[1 Step]",
				"docs/plain.md: line 16: mermaid diagram has no numbered callouts",
				"docs/flow.md: diagram callout 1 appears 2 times (line 1)",
				"docs/flow.md: callout table row A appears 2 times (line 9)",
				"docs/flow.md: diagram callout C has no matching table row (line 1)",
				"docs/flow.md: callout table row B has no matching diagram callout (line 9)",
				"docs/unpaired.md: callout table at line 1 has no preceding mermaid diagram",
				"docs/unpaired.md: mermaid diagram at line 5 has callouts but no associated callout table")
		}},
		{"pr body callout missing its row fails", func(t *testing.T) {
			t.Parallel()
			body := mermaidCalloutsDoc(
				mermaidCalloutsDiagram(`Rel(a, b, "A")`, `Rel(b, c, "B")`),
				mermaidCalloutsTable("A"),
			)["docs/flow.md"]
			out, code := mermaidCallouts(t, map[string]string{"body.md": body}, "--body-file", "body.md")
			sees(t, out, code, 1, "(description): diagram callout B has no matching table row (line 1)")
		}},
	})
}

func Test0005_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"numbered connector asks for a letter", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidCallouts(t, mermaidCalloutsDoc(
				mermaidCalloutsDiagram(`Container(app, "1", "UI")`, `Rel(app, api, "2")`),
				mermaidCalloutsTable("1", "2"),
			))
			sees(t, out, code, 1, "docs/flow.md: connector callout 2 must be a letter; numbers are for nodes; run `fitness mermaid-callouts --fix` (line 1)")
			if strings.Contains(out, "has no matching") {
				t.Errorf("a numbered connector must fail only for its form, not as an orphan:\n%s", out)
			}
		}},
		{"connector letter breaking the doubling fails", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidCallouts(t, mermaidCalloutsDoc(
				mermaidCalloutsDiagram("flowchart LR", "a -->|Z| b", "b -->|AA| c", "c -->|AB| a"),
				mermaidCalloutsTable("Z", "AA", "AB"),
			))
			sees(t, out, code, 1,
				"docs/flow.md: connector callout AB breaks the letter doubling; use A to Z, then AA to ZZ, then AAA (line 1)",
				"docs/flow.md: callout table row AB breaks the letter doubling; use A to Z, then AA to ZZ, then AAA (line 8)")
			if strings.Contains(out, "callout AA ") || strings.Contains(out, "row AA ") {
				t.Errorf("AA follows the doubling and must pass:\n%s", out)
			}
		}},
		{"fix gives connectors the next free letter", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", mermaidCalloutsNumbered())
			out, code := fitness(t, repo, nil, "mermaid-callouts")
			sees(t, out, code, 1, "connector callout 2 must be a letter")
			out, code = fitness(t, repo, nil, "mermaid-callouts", "--fix")
			sees(t, out, code, 0, "docs/flow.md: converted connector callouts 2 to B, 6 to C, 12 to A", "--fix changed 1 file")
			want := mermaidCalloutsDoc(
				mermaidCalloutsDiagram(
					"flowchart LR",
					`app[1 App] -->|B| api["3 API"]`,
					"api -->|A| db((4 Store))",
					"db -- C --> app",
				),
				"| # | Description |\n| --- | --- |\n| 1 | step 1 |\n| B | step 2 |\n| 3 | step 3 |\n| 4 | step 4 |\n| A | step A |\n| C | step 6 |\n",
				mermaidCalloutsDiagram(`Container(web, "5", "UI")`, `Rel(web, api, "A", "HTTPS")`),
				"| # | Description |\n| --- | --- |\n| 5  | step 5 |\n| A  | step 12 |\n",
			)["docs/flow.md"]
			if got := mermaidCalloutsFlow(t, repo); got != want {
				t.Errorf("--fix must turn each numbered connector and its row into the next free letter, got:\n%s\nwant:\n%s", got, want)
			}
			out, code = fitness(t, repo, nil, "mermaid-callouts")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"second fix run changes nothing", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", mermaidCalloutsNumbered())
			out, code := fitness(t, repo, nil, "mermaid-callouts", "--fix")
			sees(t, out, code, 0, "--fix changed 1 file")
			fixed := mermaidCalloutsFlow(t, repo)
			out, code = fitness(t, repo, nil, "mermaid-callouts", "--fix")
			sees(t, out, code, 0, "--fix changed 0 files")
			if strings.Contains(out, "converted") {
				t.Errorf("a second --fix must list no file, got:\n%s", out)
			}
			if got := mermaidCalloutsFlow(t, repo); got != fixed {
				t.Errorf("a second --fix must change nothing, got:\n%s\nwant:\n%s", got, fixed)
			}
		}},
		{"numbered connector sharing a node number stays unfixed", func(t *testing.T) {
			t.Parallel()
			files := mermaidCalloutsDoc(
				mermaidCalloutsDiagram(`Container(app, "2", "UI")`, `Rel(app, api, "2")`),
				mermaidCalloutsTable("2"),
			)
			repo := example(t, "happyRepo", files)
			out, code := fitness(t, repo, nil, "mermaid-callouts", "--fix")
			sees(t, out, code, 1,
				"docs/flow.md: --fix left the mermaid diagram at line 1 unchanged: callout 2 marks a node too",
				"--fix changed 0 files",
				"docs/flow.md: connector callout 2 must be a letter")
			if got := mermaidCalloutsFlow(t, repo); got != files["docs/flow.md"] {
				t.Errorf("--fix must leave a diagram it cannot convert byte for byte, got:\n%s", got)
			}
		}},
	})
}

func Test0005_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"callout table without why column fails", func(t *testing.T) {
			t.Parallel()
			out, code := calloutWhy(t, map[string]string{"docs/design.md": calloutWhyDoc(calloutWhyNoWhy)})
			sees(t, out, code, 1, `docs/design.md: callout table at line 10 is missing a "Why" column`)
		}},
		{"numbered row with empty why fails", func(t *testing.T) {
			t.Parallel()
			out, code := calloutWhy(t, map[string]string{"docs/design.md": calloutWhyDoc(calloutWhyHead, "| [2] | an edge |  |")})
			sees(t, out, code, 1, `docs/design.md: callout table row [2] has an empty "Why" cell (line 10)`)
		}},
		{"short numbered row fails as empty why", func(t *testing.T) {
			t.Parallel()
			out, code := calloutWhy(t, map[string]string{"docs/design.md": calloutWhyDoc(calloutWhyHead, "| 3 | an edge |")})
			sees(t, out, code, 1, `docs/design.md: callout table row 3 has an empty "Why" cell (line 10)`)
		}},
		{"filled why in every numbered row passes", func(t *testing.T) {
			t.Parallel()
			doc := calloutWhyDoc("| # | Description | WHY |\n| --- | --- | --- |", "| 1 | an edge | it carries the request |")
			out, code := calloutWhy(t, map[string]string{"docs/design.md": doc})
			sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		}},
		{"unnumbered row with empty why passes", func(t *testing.T) {
			t.Parallel()
			out, code := calloutWhy(t, map[string]string{"docs/design.md": calloutWhyDoc(calloutWhyHead, "| n/a | an edge |  |")})
			sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
		}},
		{"why column missing with no diagram fails", func(t *testing.T) {
			t.Parallel()
			out, code := calloutWhy(t, map[string]string{"docs/design.md": "# Design\n\n" + calloutWhyNoWhy + "\n"})
			sees(t, out, code, 1, `docs/design.md: callout table at line 3 is missing a "Why" column`)
		}},
		{"docs with only plain tables pass for why", func(t *testing.T) {
			t.Parallel()
			out, code := calloutWhy(t, map[string]string{"docs/plain.md": "# Plain\n\n| Name | Value |\n| --- | --- |\n| a | b |\n"})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "mermaid-callout-why")
		}},
		{"pr body table lacking why fails", func(t *testing.T) {
			t.Parallel()
			out, code := calloutWhy(t, map[string]string{"body.md": "## Design\n\n" + calloutWhyNoWhy + "\n"}, "--body-file", "body.md")
			sees(t, out, code, 1, `(description): callout table at line 3 is missing a "Why" column`)
		}},
		{"rel label with prose fails", func(t *testing.T) {
			t.Parallel()
			body := "Rel(app, api, \"6 Uses over HTTP\")\nBiRel(app, db, \"6 talks both ways\")\nRel_Back(api, app, '6 replies')\nRel(api, db,\n  \"6 wrapped prose\")\nRel(api, db, \"B reads\")"
			out, code := mdpCheck(t, map[string]string{"docs/arch.md": mdpPaired(body)})
			sees(t, out, code, 1, mdpError("6 Uses over HTTP"), mdpError("6 talks both ways"), mdpError("6 replies"), mdpError("6 wrapped prose"), mdpError("B reads"))
		}},
		{"flowchart pipe label with prose fails", func(t *testing.T) {
			t.Parallel()
			out, code := mdpCheck(t, map[string]string{"docs/arch.md": mdpPaired("flowchart LR\napp -->|6 calls the api| api -->|6 reads| db -->|A calls back| app")})
			sees(t, out, code, 1, mdpError("6 calls the api"), mdpError("6 reads"), mdpError("A calls back"))
		}},
		{"colon edge label with prose fails", func(t *testing.T) {
			t.Parallel()
			out, code := mdpCheck(t, map[string]string{"docs/arch.md": mdpPaired("flowchart LR\napp --> api : 6 sends data   \napi --> db : C writes")})
			sees(t, out, code, 1, mdpError("6 sends data"), mdpError("C writes"))
		}},
		{"bare number labels pass", func(t *testing.T) {
			t.Parallel()
			body := "flowchart LR\napp -->|6| api\napi --> db : 6\ndb --> app : uses 6\n%% note : 6 something\nRel(app, api, \"6\")\nclassDef x fill:#6366f1"
			out, code := mdpCheck(t, map[string]string{"docs/arch.md": mdpPaired(body)})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"prose label fails without a callout table", func(t *testing.T) {
			t.Parallel()
			out, code := mdpCheck(t, map[string]string{"docs/arch.md": mdpDiagram(`Rel(app, api, "6 Uses over HTTP")`)})
			sees(t, out, code, 1, mdpError("6 Uses over HTTP"))
		}},
		{"prose in node names passes", func(t *testing.T) {
			t.Parallel()
			body := "Container(app, \"6 App Shell\", \"SwiftUI\")\nRel(app, api, \"6\")\na[\"6 node\"] -->|6| b"
			out, code := mdpCheck(t, map[string]string{"docs/arch.md": mdpPaired(body)})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"pr body prose label fails", func(t *testing.T) {
			t.Parallel()
			files := map[string]string{"body.md": mdpPaired(`Rel(app, api, "6 Uses over HTTP")`)}
			out, code := mdpCheck(t, files, "--body-file", "body.md")
			sees(t, out, code, 1, `has a relationship label with prose ("6 Uses over HTTP")`)
		}},
		{"repo with no diagram passes prose check", func(t *testing.T) {
			t.Parallel()
			out, code := mdpCheck(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "mermaid-diagram-prose")
		}},
		{"paragraph between diagram and table fails", func(t *testing.T) {
			t.Parallel()
			out, code := tableGap(t, map[string]string{"doc.md": tableGapDoc(tableGapDiagram, tableGapProse, tableGapTable)})
			sees(t, out, code, 1, tableGapError("doc.md", "8"))
		}},
		{"paragraph between legend and table fails", func(t *testing.T) {
			t.Parallel()
			doc := tableGapDoc(tableGapDiagram, tableGapLegend, tableGapProse, tableGapTable)
			out, code := tableGap(t, map[string]string{"doc.md": doc})
			sees(t, out, code, 1, tableGapError("doc.md", "13"))
		}},
		{"matching caption alone passes", func(t *testing.T) {
			t.Parallel()
			out, code := tableGap(t, map[string]string{
				"a.md": tableGapDoc(tableGapDiagram, tableGapLegend, tableGapCaption, tableGapTable),
				"b.md": tableGapDoc(tableGapDiagram, "Numbers on classes and relationships match the callout table.", tableGapTable),
				"c.md": tableGapDoc(tableGapDiagram, `Numbers match the callout table. "Screen" is a terminal state.`, tableGapTable),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"paragraph after caption fails", func(t *testing.T) {
			t.Parallel()
			doc := tableGapDoc(tableGapDiagram, tableGapLegend, tableGapCaption, tableGapProse, tableGapTable)
			out, code := tableGap(t, map[string]string{"doc.md": doc})
			sees(t, out, code, 1, tableGapError("doc.md", "15"))
			if strings.Contains(read(out), read(tableGapError("doc.md", "13"))) {
				t.Errorf("the caption at line 13 must not be flagged:\n%s", out)
			}
		}},
		{"prose outside the gap passes", func(t *testing.T) {
			t.Parallel()
			out, code := tableGap(t, map[string]string{
				"a.md": tableGapDoc("Intro prose before the diagram.", tableGapDiagram, tableGapCaption, tableGapTable, "Trailing prose after the table."),
				"b.md": tableGapDoc("Prose above.", "```mermaid\nContainer(app, \"UI\")\n```", "Prose below."),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"markdown with no diagram passes gap check", func(t *testing.T) {
			t.Parallel()
			out, code := tableGap(t, map[string]string{"doc.md": tableGapDoc("Just prose, nothing to see here.")})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "mermaid-diagram-table-gap")
		}},
		{"pr body paragraph in gap fails", func(t *testing.T) {
			t.Parallel()
			body := strings.Join([]string{tableGapDiagram, tableGapProse, tableGapTable}, "\n\n") + "\n"
			out, code := tableGap(t, map[string]string{"body.md": body}, "--body-file", "body.md")
			sees(t, out, code, 1, tableGapError("(description)", "6"))
		}},
		{"flowchart callout with no style class fails", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidLegend(t, mermaidLegendDoc(mermaidLegendDiagram(
				"flowchart TB", "persona((1 Persona))", "classDef persona fill:#eef")))
			sees(t, out, code, 1, "README.md: callout 1 (persona) has no style class (diagram at line 9)")
		}},
		{"numbered callouts without classdef legend fail", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidLegend(t, mermaidLegendDoc(mermaidLegendDiagram(
				"flowchart TB", "persona((1 Persona)):::persona")))
			sees(t, out, code, 1, "README.md: diagram at line 9 has numbered callouts but no classDef legend")
		}},
		{"inline class or style callouts pass", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidLegend(t, mermaidLegendDoc(mermaidLegendDiagram(
				"flowchart TB",
				"persona((1 Persona)):::persona",
				"screen[2 Screen]",
				"db[(3 Store)]",
				"class screen, db thing",
				"dec{4 Choice}",
				"style dec fill:#eef",
				"classDef persona fill:#eef",
				"classDef thing fill:#f8f")))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"non flowchart diagram passes legend check", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidLegend(t, mermaidLegendDoc(mermaidLegendDiagram(
				"block-beta", "a[1 Store]", "b[2 Screen]")))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"labels without separate number pass", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidLegend(t, mermaidLegendDoc(mermaidLegendDiagram(
				"graph LR", "screen[Screen 2]", "persona((1Persona))")))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"broken diagram named by its own line", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidLegend(t, mermaidLegendDoc(
				mermaidLegendDiagram("flowchart TB", "a[1 A]:::x", "classDef x fill:#eef"),
				mermaidLegendDiagram("flowchart TB", "b[2 B]:::x", "classDef y fill:#eef", "c[3 C]")))
			sees(t, out, code, 1, "README.md: callout 3 (c) has no style class (diagram at line 15)")
		}},
		{"pr body unstyled callout fails", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidLegend(t, map[string]string{"body.md": mermaidLegendDiagram(
				"flowchart TB", "persona((1 Persona))", "classDef persona fill:#eef")},
				"--body-file", "body.md")
			sees(t, out, code, 1, "(description): callout 1 (persona) has no style class (diagram at line 1)")
		}},
		{"docs with no diagrams pass legend check", func(t *testing.T) {
			t.Parallel()
			out, code := mermaidLegend(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "mermaid-legend")
		}},
	})
}

func Test0005_5(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"test only repo with no architecture passes", func(t *testing.T) {
			t.Parallel()
			for _, files := range []map[string]string{nil, {
				"vite.config.ts":                        "export default {};\n",
				"tools/check.sh":                        "#!/bin/sh\nexit 0\n",
				"src/store/store_test.go":               "package store\n",
				"src/app.test.ts":                       "export {};\n",
				"src/__tests__/app.ts":                  "export {};\n",
				"pom.xml":                               javaPom(""),
				"src/test/java/AppTest.java":            "class AppTest {}\n",
				"src/integrationTest/java/AppFlow.java": "class AppFlow {}\n",
			}} {
				out, code := acRun(t, files)
				sees(t, out, code, 0)
				passedWithNoFiles(t, out, acCheck)
			}
		}},
		{"renumbered callouts keep diagram ids passing", func(t *testing.T) {
			t.Parallel()
			for _, root := range []string{acRoot, "docs/architecture/"} {
				files := acFixture(root)
				out, code := acRun(t, files)
				sees(t, out, code, 0)
				passedWith(t, out, acCheck, "5")
				files[root+acContext] = strings.NewReplacer("\"1 User\"", "\"21 User\"", "\"2 App\"", "\"22 App\"", "\"3\"", "\"23\"", "| 1 |", "| 21 |", "| 2 |", "| 22 |", "| 3 |", "| 23 |").Replace(files[root+acContext])
				out, code = acRun(t, files)
				sees(t, out, code, 0)
				files[root+acContext] = acView("---\ntitle: Context\n---\nflowchart LR\naccTitle: Context\nuser((1 User)):::persona\nUpdateApp[2 App]:::system\nuser -->|3| UpdateApp\nclassDef persona fill:#fff", "A person.", "The system.", "Uses it.")
				out, code = acRun(t, files)
				sees(t, out, code, 0)
			}
		}},
		{"partial or missing architecture fails for views", func(t *testing.T) {
			t.Parallel()
			for _, file := range []string{"src/main.go", "Sources/App/App.swift", "src/app.ts", "src/app.js", "src/main/java/App.java"} {
				out, code := acRun(t, map[string]string{file: "// The app.\n"})
				sees(t, out, code, 1, "docs/architecture/ has no views, but product code such as "+file+" needs them; add the five views",
					"actors (01-system-context.md), containers (02-containers.md), components (03-components.md), code (04-code.md), and journey (05-user-journey.md)")
			}
			files := acFixture(acRoot)
			delete(files, acRoot+acJourney)
			out, code := acRun(t, files)
			sees(t, out, code, 1, "architecture/01-system-context.md:1: missing architecture view 05-user-journey.md")
			files = acFixture(acRoot)
			files[acRoot+acCode] = "# View\n\nNo diagram yet.\n"
			out, code = acRun(t, files)
			sees(t, out, code, 1, "architecture/04-code.md:1: view has no diagram outside its legends")
		}},
		{"duplicate or unresolved ids name diagram lines", func(t *testing.T) {
			t.Parallel()
			cases := []struct{ view, old, new, want string }{
				{acContext, "System(app, \"2 App\")", "System(user, \"2 App\")", "architecture/01-system-context.md:6: duplicate ID user; first at architecture/01-system-context.md:5"},
				{acContext, "System(app,", "System(2app,", "architecture/01-system-context.md:6: invalid ID \"2app\""},
				{acContext, "Rel(user, app,", "Rel(user, ghost,", "architecture/01-system-context.md:7: unresolved ID ghost"},
				{acCode, "save --> Record", "save --> Missing", "architecture/04-code.md:9: unresolved ID Missing"},
				{acComponent, "Container_Boundary(service, \"1 Service (src/cmd/service)\") {\n", "", "architecture/03-components.md:5: component entry needs an enclosing boundary"},
				{acCode, "classDiagram", "sequenceDiagram", "architecture/04-code.md:3: unsupported Mermaid form \"sequenceDiagram\""},
				{acContext, "Rel(user, app, \"3\")", "Rel(user, app, \"3\")\nShape(user, \"9\")", "architecture/01-system-context.md:8: unsupported C4 macro \"Shape\""},
				{acContainer, "Person(user, \"1 User\")", "Person(user, \"1 Customer\")", "architecture/02-containers.md:5: actor user is named \"Customer\" here and \"User\" at architecture/01-system-context.md:5"},
				{acContainer, "System_Boundary", "Person(member, \"5 User\")\nSystem_Boundary", "architecture/02-containers.md:6: actor name \"User\" is ambiguous: member here and user at architecture/01-system-context.md:5"},
			}
			for _, c := range cases {
				out, code := acRun(t, acEdit(c.view, c.old, c.new))
				sees(t, out, code, 1, c.want)
			}
		}},
		{"two elements naming one path fail", func(t *testing.T) {
			t.Parallel()
			out, code := acRun(t, acEdit(acComponent, "Component(entry, \"2 Entry\", \"Go\")", "Component(entry, \"2 Entry\", \"Go\", \"src/store/store.go\")"))
			sees(t, out, code, 1, "architecture/03-components.md:7: overlapping source claims: store claims src/store and entry claims src/store/store.go at architecture/03-components.md:6")
			out, code = acRun(t, acEdit(acComponent, "\"src/store\"", "\"src/stores\""))
			sees(t, out, code, 1, "architecture/03-components.md:7: unresolved source path src/stores")
		}},
		{"unreadable source link names row", func(t *testing.T) {
			t.Parallel()
			for anchor, want := range map[string]string{"#L900)": "src/cmd/service/main.go#L900", "#line3)": "src/cmd/service/main.go#line3", "#L3-L1)": "src/cmd/service/main.go#L3-L1"} {
				out, code := acRun(t, acEdit(acComponent, acEvidence, anchor))
				sees(t, out, code, 1, "architecture/03-components.md:17: unresolved evidence "+want)
			}
			out, code := acRun(t, acEdit(acComponent, "service/main.go"+acEvidence, "service/gone.go"+acEvidence))
			sees(t, out, code, 1, "architecture/03-components.md:17: unresolved evidence src/cmd/service/gone.go#L3")
		}},
		{"invalid requirement link names the reference", func(t *testing.T) {
			t.Parallel()
			out, code := acRun(t, acEdit(acJourney, "`0001.1`", "`0001.9`"))
			sees(t, out, code, 1, "architecture/05-user-journey.md:28: unresolved requirement 0001.9 in docs/requirements")
			out, code = acRun(t, acEdit(acJourney, "0001-open.md", "0002-open.md"))
			sees(t, out, code, 1, "architecture/05-user-journey.md:29: unresolved requirement 0002 in docs/requirements")
			files := acFixture(acRoot)
			files[acDoc] = strings.ReplaceAll(files[acDoc], "opens\n-\n", "opens\n")
			out, code = acRun(t, files)
			sees(t, out, code, 0)
			if strings.Contains(out, acDoc) {
				t.Errorf("doc shape belongs to the requirements check, got:\n%s", out)
			}
		}},
		{"cited step resolves to its chain", func(t *testing.T) {
			t.Parallel()
			for step, want := range map[string]string{"`0001.1.2`": "", "`0001.1.3`": "architecture/05-user-journey.md:28: unresolved requirement 0001.1.3 in docs/requirements"} {
				files := acEdit(acJourney, "`0001.1`", step)
				files[acDoc] = acChained
				out, code := acRun(t, files)
				sees(t, out, code, min(len(want), 1), want)
			}
		}},
		{"unchanged broken evidence stays reported", func(t *testing.T) {
			t.Parallel()
			drifts := []struct {
				check, want string
				files       map[string]string
				changed     []string
			}{
				{acCheck, "architecture/03-components.md:17: unresolved evidence src/cmd/service/main.go#L900", acEdit(acComponent, acEvidence, "#L900)"), []string{"README.md", "src/store/store.go", acRoot + acContext}},
				{accCheck, "src/cmd/api/main.go:1: Go command api" + accNoOwner, with(acFixture(acRoot), map[string]string{"src/cmd/api/main.go": "package main\n"}), []string{"README.md", "src/store/store.go", acRoot + acContainer}},
				{aoCheck, "src/api/api.go" + aoNoOwner, with(aoFixture(), map[string]string{"src/api/api.go": "package api\n"}), []string{"README.md", aoStore, aoView}},
			}
			for _, d := range drifts {
				repo := example(t, "happyRepo", d.files)
				for _, changed := range d.changed {
					out, code := fitness(t, repo, []string{"FITNESS_CHANGED_FILES=" + changed}, d.check)
					sees(t, out, code, 1, d.want)
				}
			}
		}},
		{"level repeating parent description fails", func(t *testing.T) {
			t.Parallel()
			out, code := levelBleed(t, map[string]string{
				"docs/architecture/01-context.md":    levelBleedDoc("User reads the timeline"),
				"docs/architecture/02-containers.md": levelBleedDoc("User reads the timeline", "New detail"),
			})
			sees(t, out, code, 1, levelBleedRepeat("docs/architecture/02-containers.md", "user reads the timeline", 1))
		}},
		{"distinct or blank descriptions pass", func(t *testing.T) {
			t.Parallel()
			out, code := levelBleed(t, map[string]string{
				"architecture/01-context.md":    levelBleedDoc("The whole system", ""),
				"architecture/02-containers.md": levelBleedDoc("The macOS app", ""),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"case or spacing repeat reported once", func(t *testing.T) {
			t.Parallel()
			out, code := levelBleed(t, map[string]string{
				"architecture/01-context.md": levelBleedDoc("The app"),
				"architecture/02-app.md":     levelBleedDoc("The  APP", "the app", "New detail"),
			})
			want := levelBleedRepeat("architecture/02-app.md", "the app", 1)
			sees(t, out, code, 1, want)
			if n := strings.Count(read(out), read(want)); n != 1 {
				t.Errorf("the repeat must be reported once, got %d:\n%s", n, out)
			}
		}},
		{"levels one and three sharing text pass", func(t *testing.T) {
			t.Parallel()
			out, code := levelBleed(t, map[string]string{
				"architecture/01-context.md":    levelBleedDoc("Shared line"),
				"architecture/02-containers.md": levelBleedDoc("Something else"),
				"architecture/03-components.md": levelBleedDoc("Shared line"),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"levels two and ten compare as adjacent", func(t *testing.T) {
			t.Parallel()
			out, code := levelBleed(t, map[string]string{
				"architecture/2-app.md":   levelBleedDoc("The engine"),
				"architecture/10-deep.md": levelBleedDoc("The engine"),
			})
			sees(t, out, code, 1, levelBleedRepeat("architecture/10-deep.md", "the engine", 2))
		}},
		{"markdown outside numbered level files passes", func(t *testing.T) {
			t.Parallel()
			out, code := levelBleed(t, map[string]string{
				"architecture/overview.md":    levelBleedDoc("Anything"),
				"not-architecture/01-x.md":    levelBleedDoc("Anything"),
				"not-architecture/02-x.md":    levelBleedDoc("Anything"),
				"docs/01-context-notes.md":    levelBleedDoc("Anything"),
				"docs/02-container-notes.md":  levelBleedDoc("Anything"),
				"architecture/sub/01-deep.md": levelBleedDoc("Anything"),
			})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "mermaid-level-bleed")
		}},
		{"description header column is compared", func(t *testing.T) {
			t.Parallel()
			out, code := levelBleed(t, map[string]string{
				"architecture/01-context.md": "| # | Label | DESCRIPTION |\n| --- | --- | --- |\n| 1 | Same label | The system |\n",
				"architecture/02-app.md":     "| # | Label | Description |\n| --- | --- | --- |\n| 1 | Same label | The system |\n",
			})
			sees(t, out, code, 1, levelBleedRepeat("architecture/02-app.md", "the system", 1))
			if strings.Contains(out, "same label") {
				t.Errorf("only the Description column may be compared:\n%s", out)
			}
		}},
		{"missing description header uses second column", func(t *testing.T) {
			t.Parallel()
			out, code := levelBleed(t, map[string]string{
				"architecture/01-context.md": "| # | Text | Notes |\n| --- | --- | --- |\n| 1 | Same | x |\n",
				"architecture/02-app.md":     "| # | Text | Notes |\n| --- | --- | --- |\n| 1 | same | x |\n",
			})
			sees(t, out, code, 1, levelBleedRepeat("architecture/02-app.md", "same", 1))
			if strings.Contains(out, `"x"`) {
				t.Errorf("only the second column may be compared:\n%s", out)
			}
		}},
	})
}

func Test0005_6(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"journey labels naming every doc pass", func(t *testing.T) {
			t.Parallel()
			out, code := ajRun(t, map[string]string{"src/main.go": "package main\n"})
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, ajCheck)
			out, code = ajRun(t, ajFixture())
			sees(t, out, code, 0)
			passedWith(t, out, ajCheck, "2")
			files := ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "open(2 Open the app<br/>Req 0001, 0002)")
			out, code = ajRun(t, ajEdit(files, "-->|saves| done", "-->|saves 0002| done"))
			sees(t, out, code, 0)
			passedWith(t, out, ajCheck, "3")
			out, code = ajRun(t, ajEdit(ajFixture(), ajTagOpen, "open(2 Open #9829; the app<br/>Req 0001)"))
			sees(t, out, code, 0)
		}},
		{"doc named by no journey label fails", func(t *testing.T) {
			t.Parallel()
			out, code := ajRun(t, acFixture(acRoot))
			sees(t, out, code, 1, ajMissingDoc("0001"))
			out, code = ajRun(t, ajSecondDoc(ajFixture()))
			sees(t, out, code, 1, ajMissingDoc("0002"))
			out, code = ajRun(t, ajEdit(ajSecondDoc(ajFixture()), ajTagDone, "done[3 Record saved<br/>Req 0002]"))
			sees(t, out, code, 0)
		}},
		{"label naming acceptance or unknown doc fails", func(t *testing.T) {
			t.Parallel()
			for label, want := range map[string]string{
				"open(2 Open the app<br/>Req 0001.1)":     "architecture/05-user-journey.md:18: 0001.1 is an acceptance ID; name requirement doc 0001",
				"open(2 Open the app<br/>Req 0001, 0009)": "architecture/05-user-journey.md:18: no requirement doc 0009 in docs/requirements",
			} {
				out, code := ajRun(t, ajEdit(ajFixture(), ajTagOpen, label))
				sees(t, out, code, 1, want)
			}
			out, code := ajRun(t, ajEdit(ajFixture(), ajTagOpen, "open(2 Open the app<br/>Req 0009)"))
			sees(t, out, code, 1, "architecture/05-user-journey.md:18: no requirement doc 0009 in docs/requirements")
			if strings.Contains(out, "journey node open names no requirement") {
				t.Errorf("a mistyped ID also reported its step as naming none:\n%s", out)
			}
			repo := example(t, "happyRepo", ajEdit(ajSecondDoc(ajFixture()), ajTagDone, "done[3 Record saved<br/>Req 0001, 0002]"))
			write(t, repo, map[string]string{ajDoc2: ""})
			out, code = fitness(t, repo, nil, ajCheck)
			sees(t, out, code, 1, "architecture/05-user-journey.md:19: no requirement doc 0002 in docs/requirements")
		}},
		{"label citing a step reads as its doc", func(t *testing.T) {
			t.Parallel()
			files := ajEdit(ajEdit(acFixture(acRoot), ajOpen, "open(2 Open the app<br/>Req 0001.1.1)"), ajDone, "done[3 Record saved<br/>Req 0001.1.2]")
			files[ajDoc] = acChained
			out, code := ajRun(t, files)
			sees(t, out, code, 0)
			out, code = ajRun(t, ajEdit(files, "Req 0001.1.2", "Req 0001.1.3"))
			sees(t, out, code, 1, "architecture/05-user-journey.md:19: no step 0001.1.3 in docs/requirements")
		}},
		{"label naming too many docs fails", func(t *testing.T) {
			t.Parallel()
			nine, eight := ajFixture(), ajFixture()
			for _, id := range []string{"0002", "0003", "0004", "0005", "0006", "0007", "0008", "0009"} {
				nine[strings.Replace(ajDoc2, "0002", id, 1)] = strings.ReplaceAll(nine[ajDoc], "0001", id)
				if id != "0009" {
					eight[strings.Replace(ajDoc2, "0002", id, 1)] = strings.ReplaceAll(eight[ajDoc], "0001", id)
				}
			}
			ids := "0001, 0002, 0003, 0004, 0005, 0006, 0007, 0008"
			out, code := ajRun(t, ajEdit(nine, ajTagOpen, "open(2 Open the app<br/>Req "+ids+", 0009)"))
			sees(t, out, code, 1, "architecture/05-user-journey.md:18: label names 9 requirements; split the step so each names at most 8")
			out, code = ajRun(t, ajEdit(eight, ajTagOpen, "open(2 Open the app<br/>Req "+ids+")"))
			sees(t, out, code, 0)
		}},
		{"stale generated journey region fails and fresh one passes", func(t *testing.T) {
			t.Parallel()
			fresh := ajRegion(acFixture(acRoot), "0001", "    req0001[(2 Open<br/>Req 0001)]:::system")
			out, code := ajRun(t, fresh)
			sees(t, out, code, 0)
			stale := ajSecondDoc(ajRegion(acFixture(acRoot), "0001 0002", "    req0001[(2 Open<br/>Req 0001)]:::system"))
			out, code = ajRun(t, stale)
			sees(t, out, code, 1, "the generated journey nodes differ from the requirement docs; run `fitness architecture-journey --fix` to rewrite them")
			retitled := ajRegion(acFixture(acRoot), "0001", "    req0001[(2 Open The App<br/>Req 0001)]:::system")
			out, code = ajRun(t, retitled)
			sees(t, out, code, 1, "architecture/05-user-journey.md:18: the generated journey nodes differ")
		}},
		{"each journey region holds only the docs its marker lists", func(t *testing.T) {
			t.Parallel()
			two := ajSecondDoc(ajRegion(acFixture(acRoot), "0001", "    req0001[(2 Open<br/>Req 0001)]:::system"))
			two = ajEdit(two, "    %% journey-nodes:end\n", "    %% journey-nodes:end\n    %% journey-nodes:start 0002\n    req0002[(3 Open<br/>Req 0002)]:::system\n    %% journey-nodes:end\n")
			out, code := ajRun(t, two)
			sees(t, out, code, 0)
			swapped := ajEdit(ajSecondDoc(ajRegion(acFixture(acRoot), "0001", "    req0001[(2 Open<br/>Req 0001)]:::system")), "%% journey-nodes:start 0001", "%% journey-nodes:start 0002")
			out, code = ajRun(t, swapped)
			sees(t, out, code, 1, "the generated journey nodes differ")
		}},
		{"a doc listed in no journey region or in two fails", func(t *testing.T) {
			t.Parallel()
			unlisted := ajSecondDoc(ajRegion(acFixture(acRoot), "0001", "    req0001[(2 Open<br/>Req 0001)]:::system"))
			out, code := ajRun(t, unlisted)
			sees(t, out, code, 1, "requirement doc 0002 is listed in no journey-nodes region")
			twice := ajEdit(ajRegion(acFixture(acRoot), "0001", "    req0001[(2 Open<br/>Req 0001)]:::system"), "    %% journey-nodes:end\n", "    %% journey-nodes:end\n    %% journey-nodes:start 0001\n    req0001[(2 Open<br/>Req 0001)]:::system\n    %% journey-nodes:end\n")
			out, code = ajRun(t, twice)
			sees(t, out, code, 1, "requirement doc 0001 is listed in two journey-nodes regions")
			unknown := ajRegion(acFixture(acRoot), "0001 0009", "    req0001[(2 Open<br/>Req 0001)]:::system")
			out, code = ajRun(t, unknown)
			sees(t, out, code, 1, "lists 0009, which is no requirement doc")
		}},
		{"journey region with a missing end marker fails", func(t *testing.T) {
			t.Parallel()
			open := ajRegion(acFixture(acRoot), "0001", "")
			open[ajView] = strings.Replace(open[ajView], "    %% journey-nodes:end\n", "", 1)
			out, code := ajRun(t, open)
			sees(t, out, code, 1, "has `%% journey-nodes:start` with no end marker")
		}},
		{"fix rewrites every stale journey region and nothing else", func(t *testing.T) {
			t.Parallel()
			files := ajSecondDoc(ajRegion(acFixture(acRoot), "0001", "    stale"))
			files = ajEdit(files, "    %% journey-nodes:end\n", "    %% journey-nodes:end\n    %% journey-nodes:start 0002\n    stale\n    %% journey-nodes:end\n")
			repo := example(t, "happyRepo", files)
			out, code := fitness(t, repo, nil, ajCheck, "--fix")
			sees(t, out, code, 0)
			body, err := os.ReadFile(filepath.Join(repo, ajView))
			want := strings.Replace(files[ajView], "    stale\n    %% journey-nodes:end\n    %% journey-nodes:start 0002\n    stale", "req0001[(2 Open<br/>Req 0001)]:::system\n    %% journey-nodes:end\n    %% journey-nodes:start 0002\n    req0002[(3 Open<br/>Req 0002)]:::system", 1)
			if err != nil || string(body) != want {
				t.Errorf("--fix changed more than the regions: %v\n%s", err, body)
			}
		}},
		{"ids outside the diagram cover nothing", func(t *testing.T) {
			t.Parallel()
			files := acFixture(acRoot)
			files[ajView] += ajTable + "\nThe open step fulfils 0001.\n"
			out, code := ajRun(t, files)
			sees(t, out, code, 1, ajMissingDoc("0001"), "journey node open names no requirement", "journey node done names no requirement")
		}},
		{"step with no id fails", func(t *testing.T) {
			t.Parallel()
			out, code := ajRun(t, ajEdit(ajFixture(), ajTagDone, ajDone))
			sees(t, out, code, 1, "architecture/05-user-journey.md:19: journey node done names no requirement")
			out, code = ajRun(t, ajEdit(ajFixture(), ajTagOpen, ajOpen))
			sees(t, out, code, 1, "architecture/05-user-journey.md:18: journey node open names no requirement")
			out, code = ajRun(t, ajEdit(ajFixture(), ajTagDone, "subgraph steps[Steps]\n"+ajDone+"\nend"))
			sees(t, out, code, 1, "journey node done names no requirement")
			out, code = ajRun(t, ajEdit(ajFixture(), ajTagDone, "subgraph steps[Steps<br/>Req 0001]\n"+ajDone+"\nend"))
			sees(t, out, code, 0)
		}},
		{"renamed node keeps passing and dropped id fails", func(t *testing.T) {
			t.Parallel()
			out, code := ajRun(t, ajEdit(ajFixture(), ajTagDone, "done[3 Record stored<br/>Req 0001]"))
			sees(t, out, code, 0)
			renamed := ajEdit(ajEdit(ajFixture(), ajTagDone, "stored[3 Record saved<br/>Req 0001]"), "|saves| done", "|saves| stored")
			out, code = ajRun(t, renamed)
			sees(t, out, code, 0)
			dropped := ajEdit(ajEdit(ajFixture(), ajTagDone, "done[3 Record saved<br/>Req 0002]"), ajTagOpen, "open(2 Open the app<br/>Req 0002)")
			out, code = ajRun(t, ajSecondDoc(dropped))
			sees(t, out, code, 1, ajMissingDoc("0001"))
		}},
		{"steps naming each other's docs fail", func(t *testing.T) {
			t.Parallel()
			crossed := ajEdit(ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "open(2 Open the app<br/>Req 0002)"), ajTagDone, "done[3 Save the record<br/>Req 0001]")
			out, code := ajRun(t, crossed)
			sees(t, out, code, 1, "architecture/05-user-journey.md:18: journey step open reads as requirement 0001 yet names 0002; name the doc whose behavior the step is, or reword the step",
				"architecture/05-user-journey.md:19: journey step done reads as requirement 0002 yet names 0001")
			straight := ajEdit(ajEdit(ajSecondDoc(ajFixture()), ajTagDone, "done[3 Save the record<br/>Req 0002]"), "-->|saves| done", "-->|saves 0001| done")
			out, code = ajRun(t, straight)
			sees(t, out, code, 0)
			if strings.Contains(out, "reads as requirement") {
				t.Errorf("a step naming its own doc must add no note:\n%s", out)
			}
			grouped := ajEdit(ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "subgraph open[2 Open the app<br/>Req 0002]\nend"), ajTagDone, "subgraph done[3 Save the record<br/>Req 0001]\nend")
			out, code = ajRun(t, grouped)
			if strings.Contains(out, "reads as requirement") {
				t.Errorf("a layout subgraph is no step to judge:\n%s", out)
			}
		}},
		{"one step reading as another doc adds a note", func(t *testing.T) {
			t.Parallel()
			loose := ajEdit(ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "open(2 Open the app<br/>Req 0002)"), "-->|saves| done", "-->|saves 0001| done")
			out, code := ajRun(t, loose)
			showsInfo(t, out, code, ajCheck, "architecture/05-user-journey.md:18: journey step open reads as requirement 0001 yet names 0002; check that the step names the doc whose behavior it is")
			both := ajEdit(ajSecondDoc(ajFixture()), ajTagOpen, "open(2 Open the app<br/>Req 0001, 0002)")
			out, code = ajRun(t, both)
			sees(t, out, code, 0)
			if strings.Contains(out, "reads as requirement") {
				t.Errorf("a step naming the doc it reads as must add no note:\n%s", out)
			}
		}},
		{"journey drift stays after other side changes", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", ajSecondDoc(ajFixture()))
			for _, changed := range []string{"README.md", ajDoc, ajView} {
				out, code := fitness(t, repo, []string{"FITNESS_CHANGED_FILES=" + changed}, ajCheck)
				sees(t, out, code, 1, ajMissingDoc("0002"))
			}
		}},
		{"requirements with no architecture pass actors", func(t *testing.T) {
			t.Parallel()
			out, code := aaRun(t, map[string]string{"src/main.go": "package main\n", acDoc: acDocWith("no network")})
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, aaCheck)
			out, code = aaRun(t, aaFixture("no network"))
			sees(t, out, code, 1, aaNoActor("no network", " starts with no actor", "User"))
		}},
		{"given starting with actor name or id passes", func(t *testing.T) {
			t.Parallel()
			for _, given := range []string{"User" + aaOpens, "user" + aaOpens, "User"} {
				out, code := aaRun(t, aaFixture(given))
				sees(t, out, code, 0)
				passedWith(t, out, aaCheck, "1")
			}
			for _, given := range []string{"Users open the app", "Admin" + aaOpens, "USER" + aaOpens} {
				out, code := aaRun(t, aaFixture(given))
				sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
			}
		}},
		{"numbered given reads like a bullet given", func(t *testing.T) {
			t.Parallel()
			gherkinDoc := func(given string) string {
				return strings.Replace(acDocWith(given), "- 0001.1\n    - Given "+given+"\n        - When it starts\n            - Then it opens\n",
					"1. Starts\n   1. GIVEN "+given+"\n      1. WHEN it starts\n         1. THEN it opens\n", 1)
			}
			if !strings.Contains(gherkinDoc("User"+aaOpens), "1. GIVEN") {
				t.Fatal("the fixture doc no longer holds the bullets this case numbers")
			}
			files := aaFixture("")
			files[acDoc] = gherkinDoc("User" + aaOpens)
			out, code := aaRun(t, files)
			sees(t, out, code, 0)
			files[acDoc] = gherkinDoc("Admin" + aaOpens)
			out, code = aaRun(t, files)
			sees(t, out, code, 1, acDoc+":16: Given \"Admin"+aaOpens+"\" starts with no actor")
		}},
		{"given naming only a state fails", func(t *testing.T) {
			t.Parallel()
			out, code := aaRun(t, aaFixture("User has an app"))
			sees(t, out, code, 0)
			given := "a user" + aaOpens
			out, code = aaRun(t, aaFixture(given))
			sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
		}},
		{"removed actor fails on unrelated change", func(t *testing.T) {
			t.Parallel()
			given := "Host" + aaOpens
			out, code := aaRun(t, aaWith(given, aaRel, aaHost))
			sees(t, out, code, 0)
			repo := example(t, "happyRepo", aaFixture(given))
			for _, changed := range []string{"README.md", "src/store/store.go", acRoot + acContext} {
				out, code := fitness(t, repo, []string{"FITNESS_CHANGED_FILES=" + changed}, aaCheck)
				sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
			}
		}},
		{"actor id equal to another name fails", func(t *testing.T) {
			t.Parallel()
			out, code := aaRun(t, aaWith("Host"+aaOpens, aaRel, aaHost))
			sees(t, out, code, 0)
			out, code = aaRun(t, aaWith("User"+aaOpens, aaRel, aaRel+"\nPerson(User, \"4 Host\")"))
			sees(t, out, code, 1, aaGiven+"actor \"User\" is ambiguous: user at "+aaContext+":5 and User at "+aaContext+":8")
		}},
		{"several actors joined by and pass", func(t *testing.T) {
			t.Parallel()
			out, code := aaRun(t, aaWith("User and Host open the app", aaRel, aaHost))
			sees(t, out, code, 0)
			given := "User and Admin open the app"
			out, code = aaRun(t, aaWith(given, aaRel, aaHost))
			sees(t, out, code, 1, aaNoActor(given, " names no actor at \"Admin open the app\"", "User, Host"))
			out, code = aaRun(t, aaWith("User and", aaRel, aaHost))
			sees(t, out, code, 1, aaNoActor("User and", " names no actor at \"\"", "User, Host"))
		}},
		{"rejected acceptance reports no actor finding", func(t *testing.T) {
			t.Parallel()
			files := aaFixture("User" + aaOpens)
			files[acDoc] = strings.Replace(files[acDoc], "        - When it starts\n", "", 1)
			out, code := aaRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aaCheck, "1")
			if strings.Contains(out, acDoc) {
				t.Errorf("doc shape belongs to the requirements check, got:\n%s", out)
			}
		}},
		{"system or legend node is not an actor", func(t *testing.T) {
			t.Parallel()
			given := "User" + aaOpens
			out, code := aaRun(t, aaFixture(given))
			sees(t, out, code, 0)
			for _, given := range []string{"App" + aaOpens, "app" + aaOpens, "Person" + aaOpens, "p" + aaOpens} {
				out, code := aaRun(t, aaFixture(given))
				sees(t, out, code, 1, aaNoActor(given, " starts with no actor", "User"))
			}
			out, code = aaRun(t, aaWith(given, "C4Context\nPerson(user, \"1 User\")\nSystem(app, \"2 App\")\nRel(user, app, \"3\")", "flowchart LR\nuser((1 User)):::person\napp[2 App]:::system\nuser -->|3| app"))
			sees(t, out, code, 1, aaGiven+"\""+given+"\" starts with no actor; no actor is declared at "+aaContext+":3")
		}},
	})
}

func Test0005_7(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"native app with node tools passes", func(t *testing.T) {
			t.Parallel()
			out, code := accRun(t, accNative())
			sees(t, out, code, 0)
			passedWith(t, out, accCheck, "1")
			out, code = accRun(t, acFixture(acRoot))
			sees(t, out, code, 0)
			generated := accEdit(accNative(), "\"3 App\"", "\"3 Main\"")
			generated[accPbxproj] = pbxproj([2]string{"App", "application"})
			out, code = accRun(t, generated)
			sees(t, out, code, 0)
			out, code = accRun(t, accEdit(accNative(), "Container(notify, \"5 Notify\", \"Swift\")\n", ""))
			sees(t, out, code, 1, "project.yml:17: extension target Notify"+accNoOwner)
		}},
		{"undeclared entry point or deploy file fails", func(t *testing.T) {
			t.Parallel()
			cases := []struct {
				files map[string]string
				want  string
			}{
				{map[string]string{"src/cmd/api/main.go": "package main\n"}, "src/cmd/api/main.go:1: Go command api" + accNoOwner},
				{map[string]string{accPackage: accService}, "package.json:4: Node service portal" + accNoOwner},
				{map[string]string{accPackage: "{\n  \"name\": \"portal\",\n  \"bin\": \"cli.js\"\n}\n"}, "package.json:3: Node service portal" + accNoOwner},
				{map[string]string{accPackage: accDevTools, "Dockerfile": "FROM node\n"}, "Dockerfile:1: Node service portal" + accNoOwner},
				{accRootCommand, "main.go:1: Go command tool" + accNoOwner},
				{accJavaApp, "api/src/main/java/Api.java:4: Java application api" + accNoOwner},
				{with(accJavaApp, map[string]string{"pom.xml": javaPom(""), accJavaMain: "public class Api {\n  public static void main(\n      String[] args) {}\n}\n"}), accJavaMain + ":2: Java application api" + accNoOwner},
			}
			for _, c := range cases {
				out, code := accRun(t, with(acFixture(acRoot), c.files))
				sees(t, out, code, 1, c.want)
			}
			out, code := accRun(t, accEdit(with(acFixture(acRoot), accRootCommand), "}\nRel", "Container(tool, \"5 Tool\", \"Go\")\n}\nRel"))
			sees(t, out, code, 0)
			out, code = accRun(t, with(acFixture(acRoot), map[string]string{accPackage: accDevTools, "src/gen/gen.go": "//go:build ignore\n\npackage main\n"}))
			sees(t, out, code, 0)
			out, code = accRun(t, accEdit(with(acFixture(acRoot), accJavaApp), "}\nRel", "Container(api, \"5 Api\", \"Java\")\n}\nRel"))
			sees(t, out, code, 0)
			out, code = accRun(t, with(acFixture(acRoot), with(accJavaApp, map[string]string{accJavaMain: "public class Api {\n  private static void main(String[] args) {}\n\n  static void main(String name) {}\n}\n"})))
			sees(t, out, code, 0)
		}},
		{"container technology differing from language fails", func(t *testing.T) {
			t.Parallel()
			out, code := accRun(t, acEdit(acContainer, "\"Go\", \"src/cmd/service\"", "\"Node\", \"src/cmd/service\""))
			sees(t, out, code, 1, "architecture/02-containers.md:7: container service says \"Node\" but Go command service at src/cmd/service/main.go:1 is Go")
			out, code = accRun(t, acEdit(acContainer, "\"Go\", \"src/cmd/service\"", "\"Go 1.26\", \"src/cmd/service\""))
			sees(t, out, code, 0)
		}},
		{"store without path or unnamed platform service fails", func(t *testing.T) {
			t.Parallel()
			files := accEdit(acFixture(acRoot), "}\nRel", "ContainerDb(db, \"5 Records\", \"SQLite\")\n}\nRel")
			out, code := accRun(t, files)
			sees(t, out, code, 1, "architecture/02-containers.md:8: store db names no schema, configuration, or source path")
			out, code = accRun(t, accEdit(files, "\"SQLite\")", "\"SQLite\", \"src/store\")"))
			sees(t, out, code, 0)
			files = with(accNative(), map[string]string{"App/Sync.swift": "import CloudKit\n"})
			out, code = accRun(t, files)
			sees(t, out, code, 1, "App/Sync.swift:1: CloudKit is used here but no external system in the containers view names it")
			out, code = accRun(t, accEdit(files, "}\nRel", "}\nSystem_Ext(icloud, \"7 Sync\", \"Apple\", \"CloudKit records\")\nRel"))
			sees(t, out, code, 0)
		}},
		{"unreadable manifest names the file", func(t *testing.T) {
			t.Parallel()
			out, code := accRun(t, with(acFixture(acRoot), map[string]string{accPackage: "{\n  \"name\": \"portal\",\n"}))
			sees(t, out, code, 1, "package.json:1: cannot read package.json")
			files := accNative()
			files[accProject] = strings.Replace(files[accProject], "    templates: [Extension]\n    sources: Notify\n", "    sources: Notify\n", 1)
			out, code = accRun(t, files)
			sees(t, out, code, 1, "project.yml:17: target Notify has no type")
			files = accNative()
			files[accProject] = strings.Replace(files[accProject], "sources: Notify", "sources: Missing", 1)
			out, code = accRun(t, files)
			sees(t, out, code, 1, "project.yml:17: target Notify names a source Missing that does not exist")
			out, code = accRun(t, accXcode([2]string{"App", "application"}, [2]string{"Widget", ""}))
			sees(t, out, code, 1, "App.xcodeproj/project.pbxproj:11: target Widget has no type")
		}},
		{"root path and generated manifest match products", func(t *testing.T) {
			t.Parallel()
			files := accXcode([2]string{"App", "application"}, [2]string{"Widget", "app-extension"})
			files[".gitattributes"] = "*.pbxproj linguist-generated\n"
			out, code := accRun(t, files)
			sees(t, out, code, 0)
			files = accNative()
			files[".gitattributes"] = "project.yml linguist-generated\n"
			out, code = accRun(t, files)
			sees(t, out, code, 0)
			// A Java build at the repo root beside the service, owned by the root
			// path alone.
			root := accEdit(with(acFixture(acRoot), map[string]string{
				"pom.xml":                  javaPom(""),
				"src/main/java/Start.java": "public class Start {\n  public static void main(String[] args) {}\n}\n",
			}), "}\nRel", "Container(start, \"5 Start\", \"Java\", \".\")\n}\nRel")
			out, code = accRun(t, root)
			sees(t, out, code, 0)
		}},
		{"container with missing target or only source layers fails", func(t *testing.T) {
			t.Parallel()
			want := "architecture/02-containers.md:8: container widget matches no build target, entry point, or deploy file by name"
			out, code := accRun(t, accXcode([2]string{"App", "application"}, [2]string{"Widget", "app-extension"}))
			sees(t, out, code, 0)
			out, code = accRun(t, accXcode([2]string{"App", "application"}))
			sees(t, out, code, 1, want)
			files := accNative()
			files[accProject] = strings.Replace(files[accProject], "  Widget:\n    templates: [Extension]\n    sources: # folders\n      - path: Widget\n        excludes:\n          - \"**/*.md\"\n", "", 1)
			out, code = accRun(t, files)
			sees(t, out, code, 1, want)
			files = accEdit(acFixture(acRoot), "}\nRel", "Container(store, \"5 Store\", \"Go\", \"src/store\")\n}\nRel")
			out, code = accRun(t, files)
			sees(t, out, code, 1, "architecture/02-containers.md:8: container store claims only source layers (src/store): no build target, entry point, or deploy file is under them; show them in the component view")
			files["src/store/store.go"] = "package main\n"
			out, code = accRun(t, files)
			sees(t, out, code, 0)
		}},
		{"published npm library entries need exact naming", func(t *testing.T) {
			t.Parallel()
			library := with(acFixture(acRoot), map[string]string{accPackage: accLibrary, "lib/index.ts": "export const kit = 1;\n"})
			out, code := accRun(t, library)
			sees(t, out, code, 1, "package.json:3: Node library @acme/kit"+accNoOwner)
			for _, kit := range [][2]string{{"Kit dist/esm", "Bundles"}, {"@acme/kit", "Bundles"}, {"Kit", "Builds dist/esm/index.js."}} {
				out, code = accRun(t, accKit(library, kit[0], kit[1]))
				sees(t, out, code, 0)
			}
			out, code = accRun(t, accKit(library, "Kit", "Wraps @acme/kit in dist/esm"))
			sees(t, out, code, 1, "package.json:3: Node library @acme/kit"+accNoOwner, "architecture/02-containers.md:8: container kit matches no build target, entry point, or deploy file by name")
			for _, tools := range []string{
				"{\n  \"name\": \"tools\",\n  \"private\": true,\n  \"main\": \"lib/index.ts\"\n}\n",
				"{\n  \"name\": \"tools\",\n  \"main\": \"index.js\",\n  \"scripts\": {\n    \"test\": \"vitest\"\n  }\n}\n",
			} {
				out, code = accRun(t, with(library, map[string]string{accPackage: tools}))
				sees(t, out, code, 0)
			}
			tracked := with(library, map[string]string{accPackage: "{\n  \"name\": \"kit\",\n  \"exports\": \"./lib/index.ts\"\n}\n"})
			out, code = accRun(t, tracked)
			sees(t, out, code, 1, "package.json:3: Node library kit"+accNoOwner)
			out, code = accRun(t, accKit(tracked, "Kit", "Loads lib/index.ts"))
			sees(t, out, code, 0)
		}},
		{"no architecture passes component check", func(t *testing.T) {
			t.Parallel()
			out, code := aoRun(t, map[string]string{"src/main.go": "package main\n"})
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, aoCheck)
		}},
		{"one source claimed by two components fails", func(t *testing.T) {
			t.Parallel()
			second := "\n```mermaid\nC4Component\nContainer_Boundary(jobs, \"5 Jobs\") {\nComponent(cache, \"6 Cache\", \"Go\", \"src/store\")\n}\n```\n\n| # | Description | Why |\n| --- | --- | --- |\n| 5 | Jobs. | Later. |\n| 6 | Cache. | Later. |\n"
			files := aoFixture()
			files[aoView] += second
			out, code := aoRun(t, files)
			sees(t, out, code, 1, aoStore+":1: components store at "+aoView+":7 and cache at "+aoView+":22 both name this source")
			files = aoEdit(files, "\"Go\", \"src/store\")\n}\n```", "\"Go\", \"src/cache\")\n}\n```")
			files["src/cache/cache.go"] = "package cache\n"
			out, code = aoRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
		}},
		{"edge evidence must come from first component", func(t *testing.T) {
			t.Parallel()
			row := aoView + ":17: "
			cases := []struct{ old, new, want string }{
				{"cmd/service/main.go#L3)", "store/store.go#L1)", "evidence src/store/store.go#L1 is owned by store, not entry"},
				{"#L3)", "#L1)", "evidence src/cmd/service/main.go#L1 names nothing store owns"},
				{"../src/cmd/service/main.go#L3)", "../tools/check.sh)", "unsupported evidence tools/check.sh; Go and Swift sources only"},
			}
			for _, c := range cases {
				out, code := aoRun(t, aoEdit(aoFixture(), c.old, c.new))
				sees(t, out, code, 1, row+c.want)
			}
			files := aoFixture()
			files[aoView] = aoComponents("", "\nSystem_Ext(mailer, \"5 Mailer\", \"SMTP relay\")\nRel(entry, mailer, \"6\")", "Relays mail.", "Sends through the relay, see [source](../src/cmd/service/main.go#L7).")
			files[aoMain] = aoImports + "\nfunc send(c rune) bool { return c == '\"' } // smtp relay\n"
			out, code := aoRun(t, files)
			sees(t, out, code, 1, aoView+":21: evidence src/cmd/service/main.go#L7 does not name Mailer outside comments")
			files[aoMain] = aoImports + "var body = `raw\n` // smtp relay\n"
			out, code = aoRun(t, files)
			sees(t, out, code, 1, aoView+":21: evidence src/cmd/service/main.go#L7 does not name Mailer outside comments")
			files[aoMain] = aoImports + "\nfunc send() string { return \"https://smtp\" }\n"
			out, code = aoRun(t, files)
			sees(t, out, code, 0)
			renamed := aoEdit(aoFixture(), "#L3)", "#L5)")
			renamed[aoStore] = "package records\n\nfunc Save() {}\n"
			renamed["src/store/tools.go"] = "//go:build ignore\n\npackage main\n"
			renamed[aoMain] = strings.Replace(aoImports, "store.Save()", "records.Save()", 1)
			out, code = aoRun(t, renamed)
			sees(t, out, code, 0)
			renamed[aoMain] = strings.Replace(aoImports, "store.Save()", "_ = 1", 1)
			out, code = aoRun(t, renamed)
			sees(t, out, code, 1, row+"evidence src/cmd/service/main.go#L5 names nothing store owns")
		}},
		{"component reached by another runtime fails", func(t *testing.T) {
			t.Parallel()
			out, code := aoRun(t, aoWorker())
			sees(t, out, code, 1, "src/cmd/worker/main.go:3: Go command worker reaches store here, but store is under container service at "+aoView+":5, not worker; a component two runtimes share belongs outside every container boundary")
			shared := aoEdit(aoWorker(), "Component(store, \"3 Store\", \"Go\", \"src/store\")\n}", "}\nSystem_Boundary(shared, \"Shared\") {\nComponent(store, \"3 Store\", \"Go\", \"src/store\")\n}")
			out, code = aoRun(t, shared)
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
			moved := aoEdit(aoEdit(shared, aoEntry, ""), "Component(job,", aoEntry+"Component(job,")
			out, code = aoRun(t, moved)
			sees(t, out, code, 1, aoMain+":1: Go command service is built from entry, which is under container worker at "+aoView+":11, not service")
			files := aoEdit(aoNative(), "Rel(ui, model, \"8\")", "Rel(ui, model, \"8\")\nRel(glance, model, \"9\")")
			files["Widget/Widget.swift"] = "import WidgetKit\n\nlet shared = Store()\n"
			out, code = aoRun(t, files)
			sees(t, out, code, 1, "Widget/Widget.swift:3: extension target Widget reaches model here, but model is under container main at "+aoView+":5, not widget")
			out, code = aoRun(t, aoEdit(files, aoModel+"}", "}\nSystem_Boundary(shared, \"Shared\") {\n"+aoModel+"}"))
			sees(t, out, code, 0)
		}},
		{"swift names reach the right declaration", func(t *testing.T) {
			t.Parallel()
			out, code := aoRun(t, aoResolved("        if save { store.save(all: []) }\n"))
			sees(t, out, code, 1, aoHome+":7: ui names Store.save(all:), owned by sync"+aoNotDrawn+"ui -> sync edge",
				"App/Model/Saving.swift:8: model names Store.save(all:), owned by sync"+aoNotDrawn+"model -> sync edge")
			out, code = aoRun(t, aoResolved("        if save { load(into: store) }\n"))
			sees(t, out, code, 1, aoHome+":7: ui names load(into:), owned by sync"+aoNotDrawn+"ui -> sync edge")
			files := aoResolved("        func load(into store: Store) {}\n        if save { load(into: store) }\n")
			files["App/Model/Saving.swift"] = "extension Store {\n    func save(_ record: String) {}\n}\n"
			out, code = aoRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
		}},
		{"product source with no component owner fails", func(t *testing.T) {
			t.Parallel()
			files := with(aoFixture(), map[string]string{
				"src/api/api.go":        "package api\n",
				"src/api/api.config.ts": "export const port = 8080;\n",
				"index.mjs":             "export const ready = true;\n",
				"rollup.config.mjs":     "export default { input: \"index.mjs\" };\n",
				".eslintrc.cjs":         "module.exports = {};\n",
				"package.json":          "{}\n",
			})
			out, code := aoRun(t, files)
			sees(t, out, code, 1, "src/api/api.go"+aoNoOwner, "src/api/api.config.ts"+aoNoOwner, "index.mjs"+aoNoOwner)
			if strings.Contains(out, "rollup.config.mjs") || strings.Contains(out, ".eslintrc") {
				t.Errorf("a root tool config needs no owner:\n%s", out)
			}
			files = aoEdit(files, "\"src/store\")", "\"src/store\", \"and src/api\")")
			files = aoEdit(files, aoEntry, "Component(entry, \"2 Entry\", \"Go\", \"Loads index.mjs, per package.json\")\n")
			files = aoEdit(files, "| `src/cmd/service` starts it. |", "| `src/cmd/service` starts it, per `package.json`. |")
			files = aoEdit(files, "| Keeps records. |", "| Keeps records, `index.mjs` and `missing.mjs` export them, per `package.json`. |")
			out, code = aoRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
			out, code = fitness(t, example(t, "happyRepo", files), nil, acCheck)
			sees(t, out, code, 0)
			files["src/api/api.go"], files["src/api/api.config.ts"] = "", ""
			out, code = fitness(t, example(t, "happyRepo", files), nil, acCheck)
			sees(t, out, code, 1, aoView+":7: unresolved source path src/api")
			tooling := aoNative()
			tooling[aoFrames] = "import Foundation\n\nlet frameRate = 30\n"
			out, code = aoRun(t, tooling)
			sees(t, out, code, 1, aoFrames+aoNoOwner)
			tooling = aoTooling("Component(frames, \"9 Build log frames\", \"Swift\", \".claude/skills/build-log\")", "Captures build log frames.")
			tooling[aoFrames] = "import Foundation\n\nlet frameRate = 30\n"
			out, code = aoRun(t, tooling)
			sees(t, out, code, 0)
		}},
		{"import or reference missing from the view fails", func(t *testing.T) {
			t.Parallel()
			out, code := aoRun(t, aoFixture())
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
			want := aoMain + ":3: entry imports src/store, owned by store" + aoNotDrawn + "entry -> store edge"
			out, code = aoRun(t, aoEdit(aoFixture(), aoEdge, ""))
			sees(t, out, code, 1, want)
			out, code = aoRun(t, aoEdit(aoFixture(), aoEdge, "Rel(store, entry, \"4\")"))
			sees(t, out, code, 1, want+"; "+aoView+":9 draws store -> entry")
			out, code = aoRun(t, aoNative())
			sees(t, out, code, 0)
			out, code = aoRun(t, aoEdit(aoNative(), "Rel(ui, model, \"8\")", ""))
			sees(t, out, code, 1, aoHome+":4: ui names Store, owned by model"+aoNotDrawn+"ui -> model edge")
			files := with(aoFixture(), map[string]string{
				"src/api/api.go":           "package api\n\nfunc Ping() {}\n",
				"src/cmd/service/linux.go": "//go:build linux\n\npackage main\n\nimport \"example.com/app/src/api\"\n\nfunc ping() { api.Ping() }\n",
			})
			api := "Component(api, \"5 API\", \"Go\", \"src/api\")\n"
			files[aoView] = aoComponents(api, "", "Answers pings.")
			out, code = aoRun(t, files)
			sees(t, out, code, 1, "src/cmd/service/linux.go:5: entry imports src/api under go:build linux, owned by api"+aoNotDrawn+"entry -> api edge")
			files[aoView] = aoComponents(api, "\nRel(entry, api, \"6\")", "Answers pings.", "Pings the API.")
			out, code = aoRun(t, files)
			sees(t, out, code, 0)
			native := aoEdit(aoNative(), "Rel(ui, model, \"8\")", "")
			native[aoHome] = "import SwiftUI\n\nstruct HomeScreen {\n#if DEBUG // tools\n    let store = Store()\n#endif\n}\n"
			out, code = aoRun(t, native)
			sees(t, out, code, 1, aoHome+":5: ui names Store under #if DEBUG, owned by model"+aoNotDrawn+"ui -> model edge")
			native[aoHome] = "import SwiftUI\n\nstruct HomeScreen {\n#if DEBUG\n    let debug = true\n#endif // DEBUG\n    let store = Store()\n}\n"
			out, code = aoRun(t, native)
			sees(t, out, code, 1, aoHome+":7: ui names Store, owned by model"+aoNotDrawn+"ui -> model edge")
			native[aoView] = aoNative()[aoView]
			out, code = aoRun(t, native)
			sees(t, out, code, 0)
			tooling := aoTooling("Component(seed, \"9 Seed\", \"Swift\", \".codex/scripts\")", "Seeds sample records.")
			tooling[aoSeed] = "import Foundation\n\nlet sample = Store()\n"
			out, code = aoRun(t, tooling)
			sees(t, out, code, 1, aoSeed+":3: seed names Store, owned by model"+aoNotDrawn+"seed -> model edge")
			tooling = aoEdit(tooling, aoDrawn, aoDrawn+"\nRel(seed, model, \"10\")")
			tooling[aoView] += "| 10 | Seed fills the store. | Agent tooling meets the product bar. |\n"
			out, code = aoRun(t, tooling)
			sees(t, out, code, 0)
		}},
		{"scripts need no owner but started files do", func(t *testing.T) {
			t.Parallel()
			product := []string{"tools/build.mjs", "scripts/server.mjs", "lib/index.js", "pub/a.mjs"}
			files := with(aoFixture(), map[string]string{
				"package.json":         aoScripts,
				"scripts/analyze.mjs":  "export const size = 1;\n",
				"scripts/lib/stats.ts": "export const mean = 0;\n",
				"tools/gen.ts":         "export const gen = 1;\n",
				"seed.mjs":             "export const rows = [];\n",
			})
			var want []string
			for _, file := range product {
				files[file] = "export const run = 1;\n"
				want = append(want, file+aoNoOwner)
			}
			out, code := aoRun(t, files)
			sees(t, out, code, 1, want...)
			for _, tool := range []string{"scripts/analyze", "scripts/lib", "tools/gen.ts", "seed.mjs"} {
				if strings.Contains(out, tool) {
					t.Errorf("%s needs no owner:\n%s", tool, out)
				}
			}
			files = aoEdit(files, "\"src/store\")", "\"src/store\", \"and "+strings.Join(product, " ")+"\")")
			out, code = aoRun(t, files)
			sees(t, out, code, 0)
			passedWith(t, out, aoCheck, "1")
		}},
	})
}

func Test0005_8(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"no architecture or test only commit passes", func(t *testing.T) {
			t.Parallel()
			out, code := cvRun(t, map[string]string{"src/main.go": "package main\n"})
			sees(t, out, code, 0)
			passedWithNoFiles(t, out, cvCheck)
			repo, env := hooked(t, cvEdit(cvFixture(), cvView, cvEdge, ""))
			write(t, repo, map[string]string{"src/store/store_test.go": testFile("store")})
			user(t, repo, env, "add", "src/store/store_test.go")
			out, _ = user(t, repo, env, "commit", "-m", "test(store): add the store test")
			for _, check := range architectureChecks {
				passedWithNoFiles(t, out, check)
			}
			out, code = fitness(t, repo, nil, "--all", cvCheck)
			sees(t, out, code, 1, cvStore+":9: save references Record, but "+cvView+" "+cvNoEdge+"save -> Record")
		}},
		{"omitted exported type or doubled class fails", func(t *testing.T) {
			t.Parallel()
			files := cvEdit(cvFixture(), cvStore, cvSaves, cvSaves+"\n// Batch is many records.\ntype Batch struct{}\n\ntype batch struct{}\n")
			out, code := cvRun(t, files)
			sees(t, out, code, 1, cvStore+":14: store.Batch is an exported type in src/store, which "+cvView+" shows, but no node names it")
			generator := with(cvFixture(), map[string]string{"src/store/gen.go": "//go:build ignore\n\npackage main\n\n// Tool writes records by hand.\ntype Tool struct{}\n\nfunc main() {}\n"})
			out, code = cvRun(t, generator)
			sees(t, out, code, 0)
			out, code = cvRun(t, cvEdit(generator, "src/store/gen.go", "//go:build ignore\n\npackage main", "package store"))
			sees(t, out, code, 1, "src/store/gen.go:4: store.Tool is an exported type in src/store, which "+cvView+" shows, but no node names it")
			files[cvView] = cvCode("\nclass Batch[\"4 store.Batch\"]", "Many records, "+cvLink)
			out, code = cvRun(t, files)
			sees(t, out, code, 0)
			files[cvView] = cvCode("\nclass again[\"4 store.Record\"]", "The record again, "+cvLink)
			out, code = cvRun(t, files)
			sees(t, out, code, 1, cvView+":10: again and Record at "+cvView+":5 both name store.Record")
		}},
		{"call between classes missing from the view fails", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", cvEdit(cvFixture(), cvView, cvEdge, ""))
			for _, changed := range []string{"", "go.mod", cvStore, cvView} {
				out, code := fitness(t, repo, []string{"FITNESS_CHANGED_FILES=" + changed}, cvCheck)
				sees(t, out, code, 1, cvStore+":9: save references Record, but "+cvView+" "+cvNoEdge+"save -> Record"+cvAllowIt)
			}
			files := cvEdit(cvFixture(), cvView, cvEdge, "class load[\"3 store.Load() string\"]")
			files = cvEdit(files, cvView, "Save returns the record.", "Loads a name, "+cvLink)
			files[cvStore] += "\n// Load names a Record in prose only.\nfunc Load() string { return \"Record\" } // Record\n"
			out, code := cvRun(t, files)
			sees(t, out, code, 1, cvStore+":9: save references Record")
			files[cvView] = strings.Replace(files[cvView], "class load", cvEdge+"\nclass load", 1)
			files[cvView] = strings.Replace(files[cvView], "| 3 | Loads a name, "+cvLink, "| 3 | Save returns the record. | Shared by later checks. |\n| 4 | Loads a name, "+cvLink, 1)
			files[cvView] = strings.Replace(files[cvView], "\"3 store.Load", "\"4 store.Load", 1)
			out, code = cvRun(t, files)
			sees(t, out, code, 0)
			native := cvEdit(cvNative(), cvView, "HomeScreen --> Store : 5", "")
			native[cvHome] = "import SwiftUI\n\nstruct HomeScreen {\n    let title = \"Store\" // Store\n    let store = Store()\n}\n"
			out, code = cvRun(t, native)
			sees(t, out, code, 1, cvHome+":5: HomeScreen calls Store, but "+cvView+" "+cvNoEdge+"HomeScreen -> Store"+cvAllowIt)
			out, code = cvRun(t, cvNative())
			sees(t, out, code, 0)
			native = cvEdit(cvNative(), cvModel, "{ true }\n", "{ true }\n    deinit {\n        _ = HomeScreen()\n    }\n")
			native = cvEdit(native, cvView, "HomeScreen --> Store : 5", "HomeScreen --> Store : 5\nclass saveOne[\"6 Store.save(String) Bool\"]")
			native[cvView] += "| 6 | Saves one, in [Store.swift](../App/Model/Store.swift). | Shared. |\n"
			out, code = cvRun(t, native)
			sees(t, out, code, 1, cvModel+":10: Store calls HomeScreen, but "+cvView+" "+cvNoEdge+"Store -> HomeScreen"+cvAllowIt)
			native = cvEdit(native, cvView, "Store.save(String) Bool\"]", "Store.save(String) Bool\"]\nStore --> HomeScreen : 7")
			native[cvView] += "| 7 | The store makes a screen. | Shared. |\n"
			out, code = cvRun(t, native)
			sees(t, out, code, 0)
		}},
		{"drawn edge missing in source fails", func(t *testing.T) {
			t.Parallel()
			out, code := cvRun(t, cvEdit(cvNative(), cvHome, "let store = Store()", "let count = 0"))
			sees(t, out, code, 1, cvView+":11: "+cvView+" draws association HomeScreen -> Store, but no source shows HomeScreen calls or references or embeds Store")
			out, code = cvRun(t, cvSink("Put"))
			sees(t, out, code, 0)
			out, code = cvRun(t, cvSink("Write"))
			sees(t, out, code, 1, cvView+":12: "+cvView+" draws realization Disk -> Sink, but no source shows Disk conforms to Sink")
			sink := cvEdit(cvSink("Put"), cvView, "Disk is a sink.", "Disk is a sink, see [source](../src/store/store.go#L22).")
			out, code = cvRun(t, sink)
			sees(t, out, code, 0)
			out, code = cvRun(t, cvEdit(sink, cvView, "store.go#L22)", "store.go#L1)"))
			sees(t, out, code, 1, cvView+":22: evidence src/store/store.go#L1 shows no Disk -> Sink relation")
			out, code = cvRun(t, cvEdit(cvNative(), cvModel, "final class Store: Saving {", "final class Store {"))
			sees(t, out, code, 1, cvView+":10: "+cvView+" draws realization Store -> Saving, but no source shows Store conforms to Saving")
			base := cvEdit(cvNative(), cvModel, "final class Store: Saving {", "class Base {}\n\nfinal class Store: Base, Saving {")
			base = cvEdit(base, cvView, "HomeScreen --> Store : 5", "HomeScreen --> Store : 5\nclass Base[\"6 Base\"]\nBase <|-- Store : 7")
			base = cvEdit(base, cvView, "| 5 | The screen makes a store. | Shared by later checks. |\n", "| 5 | The screen makes a store. | Shared by later checks. |\n| 6 | The base, in [Store.swift](../App/Model/Store.swift). | Shared. |\n| 7 | Store extends it. | Shared. |\n")
			out, code = cvRun(t, base)
			sees(t, out, code, 0)
			out, code = cvRun(t, cvEdit(base, cvModel, "Store: Base, Saving", "Store: Saving"))
			sees(t, out, code, 1, cvView+":13: "+cvView+" draws inheritance Store -> Base, but no source shows Store embeds or inherits from Base")
			linked := cvEdit(cvFixture(), cvView, "Save returns the record.", "Save returns the record, see [source](../src/store/store.go#L10).")
			out, code = cvRun(t, linked)
			sees(t, out, code, 0)
			out, code = cvRun(t, cvEdit(linked, cvView, "store.go#L10)", "store.go#L1)"))
			sees(t, out, code, 1, cvView+":16: evidence src/store/store.go#L1 shows no save -> Record relation")
		}},
		{"class naming no symbol needs abstract mark", func(t *testing.T) {
			t.Parallel()
			files := cvFixture()
			files[cvView] = cvCode("\nclass Contract[\"4 Save contract\"]\nsave --> Contract : 5", "What saving promises.", "Save keeps the promise.")
			out, code := cvRun(t, files)
			sees(t, out, code, 1, cvView+":10: Contract \"Save contract\" names no declaration; name a Go, Swift, TypeScript, or Java symbol, or mark it <<abstract>>")
			files[cvView] = cvCode("\nclass Contract[\"4 Save contract\"] {\n<<abstract>>\n+promise: one record\n}\nsave --> Contract : 5", "What saving promises.", "Save keeps the promise.")
			out, code = cvRun(t, files)
			sees(t, out, code, 0)
			files[cvView] = strings.Replace(files[cvView], "Save keeps the promise.", "Save keeps the promise, see [source](../src/store/store.go#L10).", 1)
			out, code = cvRun(t, files)
			sees(t, out, code, 1, cvView+":23: unsupported evidence src/store/store.go#L10; save -> Contract touches an abstract node, which no source declares")
		}},
		{"name declared twice lists each declaration", func(t *testing.T) {
			t.Parallel()
			opener := "class opener[\"4 store.open()\"]"
			both := "Opens storage, in [linux](../src/store/open_linux.go) or [darwin](../src/store/open_darwin.go)."
			files := with(cvFixture(), map[string]string{
				"src/store/open_linux.go":  "//go:build linux\n\npackage store\n\nfunc open() {}\n",
				"src/store/open_darwin.go": "//go:build darwin\n\npackage store\n\nfunc open() {}\n",
			})
			files[cvView] = cvCode("\n"+opener, both)
			out, code := cvRun(t, files)
			sees(t, out, code, 1, cvView+":10: opener: store.open is declared 2 times, at src/store/open_darwin.go:5 under go:build darwin, src/store/open_linux.go:5 under go:build linux; link one file or line range, or state its signature")
			files[cvView] = cvCode("\n"+opener, "Opens storage, in [linux](../src/store/open_linux.go).")
			out, code = cvRun(t, files)
			sees(t, out, code, 0)
			native := cvNative()
			native[cvModel] = strings.Replace(native[cvModel], "    func save(_ id: String) -> Bool { true }\n", "#if DEBUG\n    func save(_ id: String) -> Bool { true }\n#else\n    func save(_ id: String) -> Bool { false }\n#endif\n    @available(iOS 17, *) func save(_ ids: [String]) -> Bool { true }\n", 1)
			native[cvView] = strings.Replace(native[cvView], "HomeScreen --> Store : 5", "HomeScreen --> Store : 5\nclass saveOne[\"6 Store.save(String) Bool\"]\nclass saveAll[\"7 Store.save([String]) Bool\"]", 1)
			native[cvView] += "| 6 | Saves one, in [Store.swift](../App/Model/Store.swift#L7-L12). | Shared. |\n| 7 | Saves many, in [Store.swift](../App/Model/Store.swift). | Shared. |\n"
			out, code = cvRun(t, native)
			sees(t, out, code, 1, cvView+":12: saveOne: Store.save is declared 2 times, at "+cvModel+":9 under #if DEBUG, "+cvModel+":11 under #else of #if DEBUG")
			native[cvView] = strings.Replace(native[cvView], "Store.swift#L7-L12", "Store.swift#L8-L9", 1)
			out, code = cvRun(t, native)
			sees(t, out, code, 0)
		}},
		{"class differing from its source fails", func(t *testing.T) {
			t.Parallel()
			broken := with(cvFixture(), map[string]string{"src/store/testdata/broken.go": "package broken\n\nfunc {\n"})
			out, code := cvRun(t, broken)
			sees(t, out, code, 0)
			passedWith(t, out, cvCheck, "1")
			apart := with(cvFixture(), map[string]string{"src/store/load.go": "package store\n\n// Load reads the record back.\nfunc (r Record) Load() Record { return r }\n"})
			apart[cvView] = cvCode("\nclass load[\"4 store.Record.Load() Record\"]\nload --> Record : 5", "Reads it back, in [load.go](../src/store/load.go).", "Load returns the record.")
			out, code = cvRun(t, apart)
			sees(t, out, code, 0)
			moved := with(cvFixture(), map[string]string{cvStore: cvRecord, "src/store/save.go": "package store\n\n" + cvSaves})
			out, code = cvRun(t, moved)
			sees(t, out, code, 1, cvView+":8: save: store.Save(id string) Record is declared at src/store/save.go:4, not in the linked src/store/store.go")
			out, code = cvRun(t, cvEdit(moved, cvView, "Writes one record, in [store.go](../src/store/store.go)", "Writes one record, in [save.go](../src/store/save.go)"))
			sees(t, out, code, 0)
			renamed := cvEdit(cvFixture(), cvStore, "type Record struct", "type Entry struct")
			renamed = cvEdit(renamed, cvStore, "Record {\n\treturn Record{", "Entry {\n\treturn Entry{")
			out, code = cvRun(t, renamed)
			sees(t, out, code, 1, cvView+":5: Record: src/store/store.go declares no store.Record")
			cases := []struct{ file, old, new, want string }{
				{cvStore, "Save(id string) Record {", "Save(id string, force bool) Record {", ":8: save: store.Save at src/store/store.go:9 has signature (string, bool) Record, not store.Save(id string) Record"},
				{cvView, cvSave, "class save[\"2 store.Save\"]", ":8: save: store.Save at src/store/store.go:9 is a function, not a type"},
				{cvStore, "\tID string\n", "\tKey string\n", ":6: Record: store.Record at src/store/store.go:4 declares no field ID"},
				{cvView, "+ID\n", "+ID: int\n", ":6: Record: field store.Record.ID at src/store/store.go:5 is string, not int"},
				{cvView, "+ID\n", "+ID\n+Load() Record\n", ":7: Record: store.Record at src/store/store.go:4 declares no method Load"},
			}
			for _, c := range cases {
				out, code := cvRun(t, cvEdit(cvFixture(), c.file, c.old, c.new))
				sees(t, out, code, 1, cvView+c.want)
			}
			files := cvEdit(cvFixture(), cvView, "+ID\n", "+ID: string\n+Load() (Record, error)\n")
			files = cvEdit(files, cvStore, cvSaves, cvSaves+"\n// Load reads the record back.\nfunc (r Record) Load() (Record, error) { return r, nil }\n")
			out, code = cvRun(t, files)
			sees(t, out, code, 0)
			embeds := cvEdit(cvFixture(), cvStore, "package store\n", "package store\n\nimport \"sync\"\n")
			embeds = cvEdit(embeds, cvStore, "\tID string\n", "\tID string\n\tsync.Mutex\n")
			out, code = cvRun(t, cvEdit(embeds, cvView, "+ID\n", "+ID\n+Mutex\n"))
			sees(t, out, code, 0)
			out, code = cvRun(t, cvEdit(embeds, cvView, "+ID\n", "+ID\n+Locker\n"))
			sees(t, out, code, 1, cvView+":7: Record: store.Record at src/store/store.go:6 declares no field Locker")
		}},
		{"typescript and java views resolve to declarations", func(t *testing.T) {
			t.Parallel()
			out, code := cvRun(t, cvScript())
			sees(t, out, code, 0)
			passedWith(t, out, cvCheck, "1")
			cases := []struct{ file, old, new, want string }{
				{"src/store/cache.ts", "", "export default class Cache {}\n\nclass Local {}\n", "src/store/cache.ts:1: Cache is an exported type in src/store, which " + cvView + " shows, but no node names it"},
				{cvView, "+save(string) boolean", "+save(number) boolean", cvView + ":11: Store: method Store.save at " + cvTSStore + ":13 has signature (string) boolean, not save(number) boolean"},
				{cvView, "number~\n+save", "string~\n+save", cvView + ":10: Store: field Store.records at " + cvTSStore + ":10 is Map<string, number>, not Map<string, string>"},
				{cvView, "disk, in [disk.ts](../src/store/disk.ts)", "disk, in [store.ts](../src/store/store.ts)", cvView + ":13: DiskStore: DiskStore is declared at " + cvTSDisk + ":3, not in the linked " + cvTSStore},
				{cvView, "DiskStore --> write : 8", "", cvTSDisk + ":6: DiskStore calls write, but " + cvView + " " + cvNoEdge + "DiskStore -> write" + cvAllowIt},
				{cvTSStore, "Store implements Saving", "Store", cvView + ":18: " + cvView + " draws realization Store -> Saving, but no source shows Store conforms to Saving"},
				{cvView, "\nHandler --> Store : 11", "", cvTSView + ":3: Handler references Store, but " + cvView + " " + cvNoEdge + "Handler -> Store" + cvAllowIt},
			}
			for _, c := range cases {
				out, code := cvRun(t, cvEdit(cvScript(), c.file, c.old, c.new))
				sees(t, out, code, 1, c.want)
			}
			out, code = cvRun(t, cvJava())
			sees(t, out, code, 0)
			passedWith(t, out, cvCheck, "1")
			java := []struct{ file, old, new, want string }{
				{cvJavaDir + "Cache.java", "", "package shop;\n\n// public class Ghost {}\npublic class Cache {}\n\nclass Local {}\n", cvJavaDir + "Cache.java:4: shop.Cache is an exported type in src/main/java/shop, which " + cvView + " shows, but no node names it"},
				{cvView, "+save(String) boolean", "+save(int) boolean", cvView + ":11: Store: method shop.Store.save at " + cvJavaDir + "Store.java:10 has signature (String) boolean, not save(int) boolean"},
				{cvView, "+count: int", "+count: long", cvView + ":20: Entry: field shop.Store.Entry.count at " + cvJavaDir + "Store.java:15 is int, not long"},
				{cvView, "+count: int", "+static final int count", cvView + ":20: Entry: member \"+static final int count\" names no field or method"},
				{cvView, "DiskStore --> write : 9\n", "", cvJavaDir + "DiskStore.java:6: DiskStore calls write, but " + cvView + " " + cvNoEdge + "DiskStore -> write" + cvAllowIt},
				{cvJavaDir + "Store.java", "Store implements Saving", "Store", cvView + ":22: " + cvView + " draws realization Store -> Saving, but no source shows Store conforms to Saving"},
			}
			for _, c := range java {
				out, code := cvRun(t, cvEdit(cvJava(), c.file, c.old, c.new))
				sees(t, out, code, 1, c.want)
			}
			shelf := cvEdit(cvJava(), cvJavaDir+"Store.java", "{}\n}\n", "{}\n}\n\nclass Shelf {\n  record Entry(long size) {}\n}\n")
			out, code = cvRun(t, cvEdit(shelf, cvView, "+count: int\n", "+count: int\n+size: long\n"))
			sees(t, out, code, 1, cvView+":21: Entry: shop.Store.Entry at "+cvJavaDir+"Store.java:15 declares no field size")
			compact := cvEdit(cvJava(), cvJavaDir+"Store.java", "int count) {}", "int count) {\n    public Entry {}\n  }")
			out, code = cvRun(t, cvEdit(compact, cvView, "+count: int\n", "+count: int\n+Entry(String, int)\n"))
			sees(t, out, code, 0)
			out, code = cvRun(t, cvEdit(cvJava(), cvJavaDir+"DiskStore.java", "    Writer.write", "    int mask = 0xE+Writer.write"))
			sees(t, out, code, 0)
			out, code = cvRun(t, cvEdit(cvJava(), cvView, "+records: Map~String, Integer~", "+Map~String, Integer~ records"))
			sees(t, out, code, 0)
		}},
	})
}

// calloutWhy runs `fitness-install -- mermaid-callout-why <args>` on
// happyRepo with files written over it.
func calloutWhy(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"mermaid-callout-why"}, args...)...)
}

// calloutWhyDoc is a doc with a mermaid diagram followed by a table made of
// rows; the table's header sits on line 10.
func calloutWhyDoc(rows ...string) string {
	return "# Design\n\n```mermaid\nflowchart LR\n    a -->|1| b\n```\n\nThe callouts:\n\n" +
		strings.Join(rows, "\n") + "\n"
}

const (
	calloutWhyNoWhy = "| # | Description |\n| --- | --- |\n| 1 | an edge |"
	calloutWhyHead  = "| # | Description | Why |\n| --- | --- | --- |"
)

// mermaidCallouts runs `fitness-install -- mermaid-callouts <args>` on
// happyRepo with files written over it.
func mermaidCallouts(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"mermaid-callouts"}, args...)...)
}

// mermaidCalloutsDiagram fences mermaid lines as a diagram block.
func mermaidCalloutsDiagram(lines ...string) string {
	return "```mermaid\n" + strings.Join(lines, "\n") + "\n```\n"
}

// mermaidCalloutsTable builds a callout table with one row per label.
func mermaidCalloutsTable(labels ...string) string {
	table := "| # | Description |\n| --- | --- |\n"
	for _, n := range labels {
		table += "| " + n + " | step " + n + " |\n"
	}
	return table
}

// mermaidCalloutsDoc joins diagram and table blocks into one doc.
func mermaidCalloutsDoc(blocks ...string) map[string]string {
	return map[string]string{"docs/flow.md": strings.Join(blocks, "\n")}
}

// mermaidCalloutsFlow reads docs/flow.md from repo.
func mermaidCalloutsFlow(t *testing.T, repo string) string {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(repo, "docs", "flow.md"))
	mustDo(t, err)
	return string(got)
}

// mermaidCalloutsNumbered is a doc whose flowchart and C4 diagrams label
// connectors with numbers beside numbered nodes and a letter in use.
func mermaidCalloutsNumbered() map[string]string {
	return mermaidCalloutsDoc(
		mermaidCalloutsDiagram(
			"flowchart LR",
			`app[1 App] -->|2| api["3 API"]`,
			"api -->|A| db((4 Store))",
			"db -- 6 --> app",
		),
		mermaidCalloutsTable("1", "2", "3", "4", "A", "6"),
		mermaidCalloutsDiagram(`Container(web, "5", "UI")`, `Rel(web, api, "12", "HTTPS")`),
		"| # | Description |\n| --- | --- |\n| 5  | step 5 |\n| 12 | step 12 |\n",
	)
}

// mdpTable is a callout table that pairs with the diagram above it.
const mdpTable = "\n| # | Description | Why |\n| --- | --- | --- |\n| 6 | uses | because |\n"

// mdpDiagram wraps body in a mermaid fence under a heading.
func mdpDiagram(body string) string {
	return "# Architecture\n\n```mermaid\n" + body + "\n```\n"
}

// mdpPaired is a doc whose diagram is paired with a callout table.
func mdpPaired(body string) string {
	return mdpDiagram(body) + mdpTable
}

// mdpCheck runs `fitness-install -- mermaid-diagram-prose <args>` on
// happyRepo with files written over it.
func mdpCheck(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"mermaid-diagram-prose"}, args...)...)
}

// mdpError is the message a user sees for one prose label in docs/arch.md.
func mdpError(label string) string {
	return `docs/arch.md: diagram at line 3 has a relationship label with prose ("` + label + `") — put descriptions in the callout table`
}

// tableGap runs `fitness-install -- mermaid-diagram-table-gap <args>` on
// happyRepo with files written over it.
func tableGap(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"mermaid-diagram-table-gap"}, args...)...)
}

// tableGapDoc joins parts into a doc, one blank line between each.
func tableGapDoc(parts ...string) string {
	return "# Doc\n\n" + strings.Join(parts, "\n\n") + "\n"
}

// The numbered diagram, its legend, the caption, and a matching callout
// table that each Given arranges around the diagram-to-table gap.
const (
	tableGapDiagram = "```mermaid\nContainer(app, \"1\", \"UI\")\nRel(app, api, \"2\")\n```"
	tableGapLegend  = "```mermaid\nPerson(p, \"Person\")\nSystem(s, \"System\")\n```"
	tableGapCaption = "Numbers on nodes and arrows match the callout table."
	tableGapTable   = "| # | Description | Why |\n| --- | --- | --- |\n| 1 | the app | because |\n| 2 | talks to | because |"
	tableGapProse   = "Four containers sit inside one app."
)

// tableGapError is the message a reader sees for prose at line in file.
func tableGapError(file, line string) string {
	return file + ": prose between the diagram and its callout table (line " + line + ") — move detail into the callout table"
}

// mermaidLegend runs `fitness-install -- mermaid-legend <args>` on happyRepo
// with files written over it.
func mermaidLegend(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"mermaid-legend"}, args...)...)
}

// mermaidLegendDiagram is a fenced mermaid diagram with the given lines.
func mermaidLegendDiagram(lines ...string) string {
	return "```mermaid\n" + strings.Join(lines, "\n") + "\n```\n"
}

// mermaidLegendDoc is happyRepo's README followed by diagrams, whose first
// fence sits on line 9.
func mermaidLegendDoc(diagrams ...string) map[string]string {
	return map[string]string{"README.md": readme + "\n" + strings.Join(diagrams, "\n")}
}

// levelBleed runs `fitness-install -- mermaid-level-bleed` on happyRepo with
// files written over it.
func levelBleed(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "mermaid-level-bleed")
}

// levelBleedDoc is a C4 level doc with one diagram and a callout table of
// [#, Description] rows.
func levelBleedDoc(descriptions ...string) string {
	lines := []string{"# Level", "", "```mermaid", `Rel(a, b, "1")`, "```", "", "| # | Description |", "| --- | --- |"}
	for i, d := range descriptions {
		lines = append(lines, "| "+strconv.Itoa(i+1)+" | "+d+" |")
	}
	return strings.Join(lines, "\n") + "\n"
}

// levelBleedRepeat is the message for file repeating text from level.
func levelBleedRepeat(file, text string, level int) string {
	return file + `: description "` + text + `" repeats level ` + strconv.Itoa(level) +
		" verbatim — lower levels should add level-specific rationale"
}

const (
	acCheck     = "architecture-contract"
	acRoot      = "architecture/"
	acContext   = "01-system-context.md"
	acContainer = "02-containers.md"
	acComponent = "03-components.md"
	acCode      = "04-code.md"
	acJourney   = "05-user-journey.md"
	acEvidence  = "#L3)"
)

// acView writes a view as people already do, with a heading, a diagram and its
// callout table, and nothing added for the check. The diagram's Nth line is line
// N+3, and the first row sits on line count+8.
func acView(diagram string, rows ...string) string {
	body := "# View\n\n```mermaid\n" + diagram + "\n```\n\n| # | Description | Why |\n| --- | --- | --- |\n"
	for i, row := range rows {
		body += "| " + strconv.Itoa(i+1) + " | " + row + " | Shared by later checks. |\n"
	}
	return body
}

// acFixture is a small complete architecture under root. up leads from the
// root back to the repository, for links.
func acFixture(root string) map[string]string {
	up := strings.Repeat("../", strings.Count(root, "/"))
	return map[string]string{
		root + acContext: acView("C4Context\nPerson(user, \"1 User\")\nSystem(app, \"2 App\")\nRel(user, app, \"3\")",
			"A person who uses the app.", "The system.", "Uses it.") +
			"\n```mermaid\nC4Context\ntitle Legend\nPerson(p, \"Person\")\n```\n",
		root + acContainer: acView("C4Container\nPerson(user, \"1 User\")\nSystem_Boundary(app, \"2 App\") {\nContainer(service, \"3 Service\", \"Go\", \"src/cmd/service\", $tags=\"v1\")\n}\nRel(user, service, \"4\")\nUpdateElementStyle(user, $fontColor=\"red\")",
			"Same person.", "The system.", "Runs requests.", "Calls it."),
		root + acComponent: acView("C4Component\nContainer_Boundary(service, \"1 Service (src/cmd/service)\") {\nComponent(entry, \"2 Entry\", \"Go\")\nComponent(store, \"3 Store\", \"Go\", \"src/store\")\n}\nRel(entry, store, \"4\")",
			"The service.", "`src/cmd/service/main.go` starts it.", "Keeps records.", "Saves through the store, see [source]("+up+"src/cmd/service/main.go"+acEvidence+"."),
		root + acCode: acView("classDiagram\nclass Record[\"1 Record\"] {\n+ID\n}\nclass save[\"2 save()\"]\nsave --> Record : 3",
			"A saved record.", "Writes one record.", "Save returns the record.") +
			"\n```mermaid\nclassDiagram\nclass Example\nParent <|-- Child : inherits\n```\n",
		root + acJourney: "## Legend\n\n```mermaid\nflowchart TB\nperson((1 Person)):::persona\nclassDef persona fill:#fff\n```\n\n| # | Description | Why |\n| --- | --- | --- |\n| 1 | Persona | A legend row. |\n\n" +
			acView("flowchart TB\nuser((1 User)):::persona\nopen(2 Open the app):::action\ndone[3 Record saved]:::screen\nuser --> open\nopen -->|saves| done\nclassDef persona fill:#fff",
				"Same person.", "Opens it, `0001.1`.", "Shows the record, see [0001]("+up+"docs/requirements/0001-open.md)."),
		"src/cmd/service/main.go": "package main\n\nfunc main() {}\n",
		"src/store/store.go":      "package store\n",
		"src/store/store_test.go": "package store\n",
		"tools/check.sh":          "#!/bin/sh\nexit 0\n",
		acDoc:                     acDocWith("User opens the app"),
	}
}

func acRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, acCheck)
}

// acEdit replaces old with new in one view of the default fixture.
func acEdit(view, old, new string) map[string]string {
	files := acFixture(acRoot)
	files[acRoot+view] = strings.Replace(files[acRoot+view], old, new, 1)
	return files
}

const (
	ajCheck   = "architecture-journey"
	ajView    = acRoot + acJourney
	ajDoc     = "docs/requirements/0001-open.md"
	ajDoc2    = "docs/requirements/0002-save.md"
	ajOpen    = "open(2 Open the app)"
	ajDone    = "done[3 Record saved]"
	ajTagOpen = "open(2 Open the app<br/>Req 0001)"
	ajTagDone = "done[3 Record saved<br/>Req 0001]"
	ajMissing = "architecture/05-user-journey.md:15: requirement %s is missing from the journey diagram"
	ajTable   = "\n## Requirements\n\n| Requirement | Journey |\n| --- | --- |\n| [0001](../docs/requirements/0001-open.md) | `open`, `done` |\n"
)

// ajFixture is the contract fixture whose journey steps both name
// requirement 0001 in their labels. The diagram fence is line 15, the open
// step line 18, and the done step line 19.
func ajFixture() map[string]string {
	return ajEdit(ajEdit(acFixture(acRoot), ajOpen, ajTagOpen), ajDone, ajTagDone)
}

// ajSecondDoc adds requirement 0002, shaped like 0001.
func ajSecondDoc(files map[string]string) map[string]string {
	files[ajDoc2] = strings.ReplaceAll(files[ajDoc], "0001", "0002")
	return files
}

func ajRun(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, ajCheck)
}

// ajEdit replaces old with new in the journey view of a fixture.
func ajEdit(files map[string]string, old, new string) map[string]string {
	files[ajView] = strings.Replace(files[ajView], old, new, 1)
	return files
}

// ajRegion replaces the fixture's two steps with a marked region listing
// ids and holding nodes, wired by the same edges.
func ajRegion(files map[string]string, ids, nodes string) map[string]string {
	files = ajEdit(files, ajOpen+":::action\ndone", "%% journey-nodes:start "+ids+"\n"+nodes+"\n    %% journey-nodes:end\ndone")
	files = ajEdit(files, ajDone+":::screen\n", "")
	files = ajEdit(files, "user --> open\nopen -->|saves| done", "user --> req0001\nreq0001 --> req0001")
	return files
}

func ajMissingDoc(id string) string {
	return strings.Replace(ajMissing, "%s", id, 1)
}

const (
	accCheck    = "architecture-containers"
	accNoOwner  = " has no container in the containers view"
	accProject  = "project.yml"
	accPackage  = "package.json"
	accPbxproj  = "App.xcodeproj/project.pbxproj"
	accDevTools = "{\n  \"name\": \"portal\",\n  \"private\": true,\n  \"bin\": null,\n  \"devDependencies\": {\n    \"prettier\": \"3.0.0\"\n  }\n}\n"
	accService  = "{\n  \"name\": \"portal\",\n  \"scripts\": {\n    \"start\": \"node server.js\"\n  }\n}\n"
)

// accRootCommand is a Go command at the repository root, which go build
// names after the module path's last element.
var accRootCommand = map[string]string{"go.mod": "module example.com/tool\n\ngo 1.26\n", "main.go": "package main\n"}

// accNative is a native app with Node tools whose XcodeGen project puts the
// Notify target on line 17 of project.yml.
func accNative() map[string]string {
	files := acFixture(acRoot)
	files["src/cmd/service/main.go"] = "package service\n\n// Service is a library.\n"
	files[acRoot+acContainer] = acView("C4Container\nPerson(user, \"1 User\")\nSystem_Boundary(app, \"2 App\") {\nContainer(main, \"3 App\", \"Swift 6, SwiftUI\", \"App/\")\nContainer(widget, \"4 Widget\", \"Swift\")\nContainer(notify, \"5 Notify\", \"Swift\")\n}\nRel(user, main, \"6\")",
		"Same person.", "The app.", "The app target.", "The widget target.", "The notification target.", "Opens it.")
	files[accProject] = "name: App\ntargetTemplates:\n  Extension:\n    type: app-extension\n    sources: [\"${target_name}\"]\ntargets:\n  App:\n    type: application # the app\n    platform: iOS\n    sources: [App]\n  Widget:\n    templates: [Extension]\n    sources: # folders\n      - path: Widget\n        excludes:\n          - \"**/*.md\"\n  Notify:\n    templates: [Extension]\n    sources: Notify\n  AppTests:\n    type: bundle.unit-test\n    platform: iOS\n    sources: [Tests]\n"
	files["App/App.swift"] = "import SwiftUI\n"
	files["Widget/Widget.swift"] = "import WidgetKit\n"
	files["Notify/Notify.swift"] = "import UserNotifications\n"
	files["Tests/AppTests.swift"] = "import XCTest\n"
	files[accPackage] = accDevTools
	return files
}

// accXcode is the native app with an Xcode project instead of XcodeGen,
// declaring the app and the widget.
func accXcode(targets ...[2]string) map[string]string {
	files := accNative()
	delete(files, accProject)
	files[acRoot+acContainer] = strings.Replace(files[acRoot+acContainer], "Container(notify, \"5 Notify\", \"Swift\")\n", "", 1)
	files[accPbxproj] = pbxproj(targets...)
	return files
}

// pbxproj renders native targets as Xcode writes them; an empty type leaves
// the productType line out.
func pbxproj(targets ...[2]string) string {
	body := "// !$*UTF8*$!\n{\n\tobjects = {\n"
	for _, t := range targets {
		body += "\t\tA /* " + t[0] + " */ = {\n\t\t\tisa = PBXNativeTarget;\n\t\t\tname = " + t[0] + ";\n"
		if t[1] != "" {
			body += "\t\t\tproductType = \"com.apple.product-type." + t[1] + "\";\n"
		}
		body += "\t\t};\n"
	}
	return body + "\t};\n}\n"
}

func accRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, accCheck)
}

// accEdit replaces old with new in the containers view of files.
func accEdit(files map[string]string, old, new string) map[string]string {
	files[acRoot+acContainer] = strings.Replace(files[acRoot+acContainer], old, new, 1)
	return files
}

// accJavaMain is the product source of accJavaApp.
const accJavaMain = "api/src/main/java/Api.java"

// accJavaApp is a Maven project in api/ whose product source declares a
// main method below a commented-out one.
var accJavaApp = map[string]string{
	"api/pom.xml":                  javaPom(""),
	accJavaMain:                    "// static void main(String[] args) proves nothing here.\npublic class Api {\n\n  public static void main(String[] args) {}\n}\n",
	"api/src/test/java/ApiIT.java": "class ApiIT {\n  public static void main(String[] args) {}\n}\n",
}

// accLibrary is an npm library shipping only built files, with main on line
// 3 and a build script making them.
const accLibrary = "{\n  \"name\": \"@acme/kit\",\n  \"main\": \"dist/cjs/index.js\",\n  \"module\": \"dist/esm/index.js\",\n  \"exports\": {\n    \".\": {\n      \"types\": \"./dist/types/index.d.ts\",\n      \"import\": \"./dist/esm/index.js\"\n    }\n  },\n  \"files\": [\"dist\"],\n  \"scripts\": {\n    \"build\": \"rollup -c\"\n  }\n}\n"

// accKit adds a kit container with label and description to the library.
func accKit(files map[string]string, label, description string) map[string]string {
	return accEdit(with(files, nil), "}\nRel", "Container(kit, \"5 "+label+"\", \"TypeScript\", \""+description+"\")\n}\nRel")
}

const (
	aaCheck   = "architecture-actors"
	acDoc     = "docs/requirements/0001-open.md"
	aaGiven   = acDoc + ":16: Given "
	aaContext = "architecture/01-system-context.md"
	aaRel     = "Rel(user, app, \"3\")"
	aaHost    = aaRel + "\nPerson(host, \"4 Host\")"
	aaOpens   = " opens the app"
)

// acDocWith is requirement 0001 with one acceptance whose Given is given.
// The Given sits on line 16.
func acDocWith(given string) string {
	return "# 0001 Open\n\n## Why\n\nI open the app.\n\n## Measurement\n\nopens\n-\ntries\n\n## Requirements\n\n- 0001.1\n    - Given " + given + "\n        - When it starts\n            - Then it opens\n"
}

// acChained is the fixture doc whose 0001.1 chains two steps, so `0001.1.2`
// cites its second.
var acChained = strings.Replace(acDocWith("User opens the app"), "        - When it starts\n            - Then it opens\n",
	"        1. When it starts\n            - Then it opens\n        2. When it opens\n            - Then it shows the record\n", 1)

// aaFixture is the default architecture with its one Given replaced.
func aaFixture(given string) map[string]string {
	files := acFixture(acRoot)
	files[acDoc] = acDocWith(given)
	return files
}

// aaWith is aaFixture with one edit to the context diagram.
func aaWith(given, oldContext, newContext string) map[string]string {
	files := aaFixture(given)
	files[acRoot+acContext] = strings.Replace(files[acRoot+acContext], oldContext, newContext, 1)
	return files
}

func aaRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, aaCheck)
}

// aaNoActor is the failure for a Given whose text names no actor among
// the declared names.
func aaNoActor(given, where, actors string) string {
	return aaGiven + "\"" + given + "\"" + where + "; actors declared at " + aaContext + ":3 are " + actors
}

const (
	aoCheck    = "architecture-components"
	aoView     = acRoot + acComponent
	aoMain     = "src/cmd/service/main.go"
	aoStore    = "src/store/store.go"
	aoHome     = "App/Screens/HomeScreen.swift"
	aoEdge     = "Rel(entry, store, \"4\")"
	aoSaves    = "Saves through the store, see [source](../src/cmd/service/main.go#L3)."
	aoEntry    = "Component(entry, \"2 Entry\", \"Go\")\n"
	aoModel    = "Component(model, \"3 Model\", \"Swift\", \"App/Model\")\n"
	aoNoOwner  = ":1: no component in the components view names this source"
	aoNotDrawn = ", but " + aoView + " draws no "
	aoImports  = "package main\n\nimport \"example.com/app/src/store\"\n\nfunc main() { store.Save() }\n"
	aoScreen   = "import SwiftUI\n\nstruct HomeScreen {\n    let store = Store()\n}\n"
	aoDrawn    = "Rel(ui, model, \"8\")"
	aoFrames   = ".claude/skills/build-log/scripts/Frames.swift"
	aoSeed     = ".codex/scripts/Seed.swift"
)

// aoTooling adds component 9, owning the agent tooling folder, to the
// native app's components view, before its drawn edge.
func aoTooling(component, row string) map[string]string {
	files := aoEdit(aoNative(), aoDrawn, component+"\n"+aoDrawn)
	files[aoView] += "| 9 | " + row + " | Agent tooling meets the product bar. |\n"
	return files
}

// aoComponents is the components view starting on line 4, with inside added to
// the service boundary and after and rows appended.
func aoComponents(inside, after string, rows ...string) string {
	diagram := "C4Component\nContainer_Boundary(service, \"1 Service (src/cmd/service)\") {\n" + aoEntry + "Component(store, \"3 Store\", \"Go\", \"src/store\")\n" + inside + "}\n" + aoEdge + after
	return acView(diagram, append([]string{"The service.", "`src/cmd/service` starts it.", "Keeps records.", aoSaves}, rows...)...)
}

// aoFixture is the Go architecture whose entry imports the store, so the
// drawn edge and its evidence on line 3 both hold.
func aoFixture() map[string]string {
	return with(acFixture(acRoot), map[string]string{
		"go.mod": "module example.com/app\n\ngo 1.26\n",
		aoMain:   aoImports,
		aoStore:  "package store\n\n// Save writes one record.\nfunc Save() {}\n",
		aoView:   aoComponents("", ""),
	})
}

// aoNative is the native app with screens reading a model, a widget, and
// a notification extension, each under its container's boundary. The
// diagram starts on line 4; the model sits on line 7.
func aoNative() map[string]string {
	files := accNative()
	for _, gone := range []string{"App/App.swift", aoMain, aoStore, "src/store/store_test.go"} {
		files[gone] = ""
	}
	files[aoHome] = aoScreen
	files["App/Model/Store.swift"] = "import Foundation\n\nfinal class Store {\n    var records: [String] = []\n}\n"
	files[aoView] = acView("C4Component\nContainer_Boundary(main, \"1 App (App/)\") {\nComponent(ui, \"2 Screens\", \"Swift\", \"App/Screens\")\n"+aoModel+"}\nContainer_Boundary(widget, \"4 Widget\") {\nComponent(glance, \"5 Glance\", \"Swift\", \"Widget/\")\n}\nContainer_Boundary(notify, \"6 Notify\") {\nComponent(alert, \"7 Alert\", \"Swift\", \"Notify/\")\n}\nRel(ui, model, \"8\")",
		"The app.", "Screens.", "The model.", "The widget.", "Draws the glance.", "The notification extension.", "Posts alerts.", "Screens read the model.")
	return files
}

func aoRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, aoCheck)
}

// aoEdit replaces old with new in the components view of files.
func aoEdit(files map[string]string, old, new string) map[string]string {
	files[aoView] = strings.Replace(files[aoView], old, new, 1)
	return files
}

// aoWorker adds a worker command whose job reads the store, its boundary
// opening on line 10, or 11 once the store moves.
func aoWorker() map[string]string {
	files := aoFixture()
	files[acRoot+acContainer] = strings.Replace(files[acRoot+acContainer], "}\nRel", "Container(worker, \"5 Worker\", \"Go\", \"src/cmd/worker\")\n}\nRel", 1)
	files[aoView] = aoComponents("", "\nContainer_Boundary(worker, \"5 Worker (src/cmd/worker)\") {\nComponent(job, \"6 Job\", \"Go\", \"src/cmd/worker\")\n}\nRel(job, store, \"7\")", "The worker.", "Runs jobs.", "Reads records.")
	files["src/cmd/worker/main.go"] = aoImports
	return files
}

// aoResolved is the native app with a sync component beside the model, each
// declaring a load overload and save extension.
func aoResolved(line7 string) map[string]string {
	files := aoEdit(aoNative(), aoModel+"}", aoModel+"Component(sync, \"9 Sync\", \"Swift\", \"App/Sync\")\n}")
	files = aoEdit(files, "Rel(ui, model, \"8\")", "Rel(ui, model, \"8\")\nRel(sync, model, \"10\")")
	files[aoView] += "| 9 | Syncs records. | Shared by later checks. |\n| 10 | Sync fills the store. | Shared by later checks. |\n"
	files["App/Model/Load.swift"] = "func load(from path: String) -> Store {\n    Store()\n}\n"
	files["App/Model/Saving.swift"] = "#if os(iOS)\nextension Store {\n#else\nextension Store: Sendable {\n#endif\n    static func < (lhs: Store, rhs: Store) -> Bool { true }\n" +
		"    func save(_ record: String) {\n        if let kept = [record].first(where: { !$0.isEmpty }) { save(all: [kept]) }\n    }\n}\n"
	files["App/Sync/Load.swift"] = "func load(into store: Store) {\n    store.records = []\n}\n"
	files["App/Sync/Saving.swift"] = "extension Store {\n    func save(all records: [String]) {}\n}\n"
	files[aoHome] = "import SwiftUI\n\nstruct HomeScreen {\n    func refresh(save: Bool) {\n        let store = load(from: \"records\")\n        store.save(\"first\")\n" + line7 + "    }\n}\n"
	return files
}

// aoScripts is a package whose scripts run tools past runner flags, beside a
// start script, a published entry, folder and exports pattern.
const aoScripts = "{\n  \"main\": \"tools/build.mjs\",\n  \"files\": [\"lib\"],\n  \"exports\": {\n    \"./x/*\": \"./pub/*.mjs\"\n  },\n  \"scripts\": {\n" +
	"    \"analyze\": \"node scripts/analyze.mjs\",\n    \"gen\": \"npx tsx tools/gen.ts && echo done\",\n    \"seed\": \"MODE=dev node -r dotenv/config seed.mjs\",\n" +
	"    \"build\": \"node tools/build.mjs\",\n    \"start\": \"node --require dotenv/config scripts/server.mjs\",\n    \"dev\": \"node --watch scripts/server.mjs\",\n" +
	"    \"check\": \"node lib/index.js --check\",\n    \"pub\": \"node pub/a.mjs\"\n  }\n}\n"

const (
	cvCheck   = "architecture-code"
	cvView    = acRoot + acCode
	cvStore   = "src/store/store.go"
	cvModel   = "App/Model/Store.swift"
	cvHome    = "App/Screens/HomeScreen.swift"
	cvLink    = "in [store.go](../src/store/store.go)."
	cvSave    = "class save[\"2 store.Save(id string) Record\"]"
	cvEdge    = "save --> Record : 3"
	cvRecord  = "package store\n\n// Record is one saved record.\ntype Record struct {\n\tID string\n}\n\n"
	cvSaves   = "// Save writes one record.\nfunc Save(id string) Record {\n\treturn Record{ID: id}\n}\n"
	cvNoEdge  = "draws no "
	cvAllowIt = " edge that allows it"
)

// cvCode is the Go code view; the record sits on line 5, save on 8, the edge
// on 9.
func cvCode(extra string, rows ...string) string {
	diagram := "classDiagram\nclass Record[\"1 store.Record\"] {\n+ID\n}\n" + cvSave + "\n" + cvEdge + extra
	return acView(diagram, append([]string{"A saved record, " + cvLink, "Writes one record, " + cvLink, "Save returns the record."}, rows...)...)
}

// cvFixture is a Go repository whose code view matches its source.
func cvFixture() map[string]string {
	return with(acFixture(acRoot), map[string]string{
		"go.mod": "module example.com/app\n\ngo 1.26\n",
		cvStore:  cvRecord + cvSaves,
		cvView:   cvCode(""),
	})
}

// cvNative is a Swift app's code view: Saving on line 5, Store 8, HomeScreen
// 9, edges 10 and 11.
func cvNative() map[string]string {
	files := acFixture(acRoot)
	files[cvModel] = "import Foundation\n\nprotocol Saving {\n    func save(_ id: String) -> Bool\n}\n\nfinal class Store: Saving {\n    func save(_ id: String) -> Bool { true }\n}\n"
	files[cvHome] = "import SwiftUI\n\nstruct HomeScreen {\n    let store = Store()\n}\n"
	files[cvView] = acView("classDiagram\nclass Saving[\"1 Saving\"] {\n<<interface>>\n}\nclass Store[\"2 Store\"]\nclass HomeScreen[\"3 HomeScreen\"]\nStore ..|> Saving : 4\nHomeScreen --> Store : 5",
		"Saves records, in [Store.swift](../App/Model/Store.swift).", "Keeps records, in [Store.swift](../App/Model/Store.swift).",
		"Lists records, in [HomeScreen.swift](../App/Screens/HomeScreen.swift).", "The store saves.", "The screen makes a store.")
	return files
}

func cvRun(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), env, cvCheck)
}

// cvEdit replaces old with new in one file of files.
func cvEdit(files map[string]string, file, old, new string) map[string]string {
	files[file] = strings.Replace(files[file], old, new, 1)
	return files
}

// architectureChecks are the six checks that read the architecture model.
var architectureChecks = []string{"architecture-actors", cvCheck, "architecture-components", "architecture-containers", "architecture-contract", "architecture-journey"}

// cvSink adds a Go interface and a type that satisfies it, drawn as a
// realization on line 12.
func cvSink(put string) map[string]string {
	files := cvFixture()
	files[cvStore] += "\n// Sink takes records.\ntype Sink interface {\n\tPut()\n}\n\n// Disk writes records.\ntype Disk struct{}\n\n// " + put + " writes one.\nfunc (Disk) " + put + "() {}\n"
	files[cvView] = cvCode("\nclass Sink[\"4 store.Sink\"]\nclass Disk[\"5 store.Disk\"]\nDisk ..|> Sink : 6", "Takes records, "+cvLink, "Writes records, "+cvLink, "Disk is a sink.")
	return files
}

const (
	cvTSStore = "src/store/store.ts"
	cvTSDisk  = "src/store/disk.ts"
	cvTSView  = "src/store/panel.tsx"
)

// cvScript is a TypeScript repository matching its code view: store members
// on lines 10 and 11, subclass 13, edges 18 to 20.
func cvScript() map[string]string {
	files := acFixture(acRoot)
	files[cvTSStore] = "import { Sink } from './sink';\n\n// Saving writes records.\nexport interface Saving {\n  save(id: string): boolean;\n}\n\n" +
		"// Store keeps records.\nexport class Store implements Saving {\n  readonly records: Map<string, number> = new Map();\n  private label = 'Sink'; // Sink names nothing here\n\n" +
		"  save(id: string): boolean {\n    return this.records.has(id);\n  }\n}\n"
	files[cvTSDisk] = "import { Store } from './store';\n\nexport class DiskStore\n  extends Store {\n  flush(force?: boolean): Promise<void> {\n    return write(this.records);\n  }\n}\n\n" +
		"export function write(records: Map<string, number>): Promise<void> {\n  return Promise.resolve();\n}\n"
	files["src/store/sink.ts"] = "class Sink {}\n\nexport { Sink };\n"
	files[cvTSView] = cvPanel
	files[cvView] = acView("classDiagram\nclass Saving[\"1 Saving\"] {\n<<interface>>\n+save(id: string) boolean\n}\nclass Store[\"2 Store\"] {\n+records: Map~string, number~\n+save(string) boolean\n}\n"+
		"class DiskStore[\"3 DiskStore\"] {\n+flush(force: boolean) Promise~void~\n}\nclass write[\"4 write(Map~string, number~) Promise~void~\"]\nclass Sink[\"5 Sink\"]\n"+
		"Store ..|> Saving : 6\nStore <|-- DiskStore : 7\nDiskStore --> write : 8\n"+
		"class Handler[\"9 Handler\"]\nclass Panel[\"10 Panel\"] {\n+host\n+rows: Array~string~[]\n+open(number, string) void\n}\nHandler --> Store : 11",
		"Saves records, in [store.ts](../src/store/store.ts).", "Keeps records, in [store.ts](../src/store/store.ts).",
		"Writes to disk, in [disk.ts](../src/store/disk.ts).", "Writes records, in [disk.ts](../src/store/disk.ts).",
		"Takes records, in [sink.ts](../src/store/sink.ts).", "The store saves.", "A disk store is a store.", "Flushing writes.",
		"Takes a store, in [panel.tsx](../src/store/panel.tsx).", "Shows records, in [panel.tsx](../src/store/panel.tsx).", "A handler reads the store.")
	return files
}

// cvPanel is TSX without semicolons: a Store alias, braces in closing tags,
// and a class after ! or ++ with a < default.
const cvPanel = "import { Store } from './store'\n\nexport type Handler = (s: Store) => void\n\nlet renders = 0\n\n" +
	"export function Title(props: { show: boolean }) {\n  return <h1>{props.show && <b>Records</b>}</h1>\n}\n\n" +
	"const root = document.getElementById('root')!\nrenders++\nexport class Panel {\n  host = root!\n  rows: Array<string>[] = []\n\n" +
	"  open(limit: number = renders < 10 ? 1 : 2, label: string): void {\n    renders += limit\n  }\n}\n"

// cvJavaDir holds the Java sources of cvJava.
const cvJavaDir = "src/main/java/shop/"

// cvJava is a Java repository matching its code view: store members 10 and
// 11, record members 19 and 20, edges 22 to 25.
func cvJava() map[string]string {
	files := acFixture(acRoot)
	files[cvJavaDir+"Saving.java"] = "package shop;\n\n// Saving writes records.\npublic interface Saving {\n  boolean save(String id);\n}\n"
	files[cvJavaDir+"Store.java"] = "package shop;\n\nimport java.util.Map;\n\n/** Store keeps records; Sink names nothing here. */\npublic class Store implements Saving {\n  final Map<String, Integer> records = Map.of(), spare = null;\n  private String label = \"Sink\";\n\n  @Override\n  public boolean save(final String id) {\n    return records.containsKey(id);\n  }\n\n  public record Entry(String key, int count) {}\n}\n"
	files[cvJavaDir+"DiskStore.java"] = "package shop;\n\npublic class DiskStore\n    extends Store {\n  public void flush(boolean force) {\n    Writer.write(records);\n  }\n}\n"
	files[cvJavaDir+"Writer.java"] = "package shop;\n\npublic final class Writer {\n  static void write(java.util.Map<String, Integer> records) {}\n}\n"
	files[cvJavaDir+"Sink.java"] = "package shop;\n\nclass Sink {}\n"
	files["src/test/java/shop/StoreTest.java"] = "package shop;\n\npublic class StoreTest {}\n"
	files[cvView] = acView("classDiagram\nclass Saving[\"1 Saving\"] {\n<<interface>>\n+save(String id) boolean\n}\nclass Store[\"2 shop.Store\"] {\n+records: Map~String, Integer~\n+save(String) boolean\n}\n"+
		"class DiskStore[\"3 DiskStore\"] {\n+flush(boolean) void\n}\nclass Writer[\"4 Writer\"]\nclass write[\"5 Writer.write(java.util.Map~String, Integer~) void\"]\nclass Entry[\"6 Store.Entry\"] {\n+key: String\n+count: int\n}\n"+
		"Store ..|> Saving : 7\nStore <|-- DiskStore : 8\nDiskStore --> write : 9\nDiskStore --> Writer : 10",
		"Saves records, in [Saving.java](../"+cvJavaDir+"Saving.java).", "Keeps records, in [Store.java](../"+cvJavaDir+"Store.java).",
		"Writes to disk, in [DiskStore.java](../"+cvJavaDir+"DiskStore.java).", "Writes, in [Writer.java](../"+cvJavaDir+"Writer.java).",
		"Writes records, in [Writer.java](../"+cvJavaDir+"Writer.java).", "One entry, in [Store.java](../"+cvJavaDir+"Store.java).",
		"The store saves.", "A disk store is a store.", "Flushing writes.", "Flushing names the writer.")
	return files
}

// cspell:ignore Descr mindmap xychart

// mermaidSyntax runs `fitness-install -- mermaid-syntax <args>` on happyRepo
// with files written over it.
func mermaidSyntax(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"mermaid-syntax"}, args...)...)
}
