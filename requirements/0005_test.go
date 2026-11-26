package requirements

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

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
