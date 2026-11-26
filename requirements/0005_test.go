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

func Test0005_9(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
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

// cspell:ignore Descr mindmap xychart

// mermaidSyntax runs `fitness-install -- mermaid-syntax <args>` on happyRepo
// with files written over it.
func mermaidSyntax(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"mermaid-syntax"}, args...)...)
}
