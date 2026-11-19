// The shared parser lives in pkg/mermaid, which keeps its fence helpers unexported, so this package duplicates them.

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/may-journal/fitness-checks/internal/diagrams"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/mdx"
	"github.com/may-journal/fitness-runner/go/pkg/mermaid"
)

func main() {
	program.Main(docSource, diagrams.Check("mermaid-diagram-table-gap", mdx.WithMDX, mermaid.RunDocCheck, validateDoc))
}

// fenceMarker captures a mermaid fence's opening backticks or tildes, so the closing fence matches as the parser does.
var fenceMarker = regexp.MustCompile("^[" + diagrams.JSWhitespace + "]*(`{3,}|~{3,})")

// captionLine matches the one caption the C4 docs place between diagram and table, opening "Numbers" and containing "match the callout table". Everything else in the gap is prose that belongs in the table.
var captionLine = regexp.MustCompile(`(?i)^numbers\b.*\bmatch the callout table\b`)

// validateDoc reports each non-blank, non-caption markdown line sitting
// between a numbered diagram (or its trailing legend block) and its callout
// table. Prose before the diagram and after the table is left alone.
func validateDoc(file, content string) []string {
	doc := mermaid.ReadDoc(content)
	lines, blocks := doc.Lines, doc.Blocks
	paired := pairedTables(blocks)
	var errors []string
	for i := range blocks {
		prev, table := precedingDiagram(blocks, i, paired)
		if prev == nil || table == nil {
			continue
		}
		errors = append(errors, gapErrors(file, lines, prev, table)...)
	}
	return errors
}

// pairedTables is the set of callout tables paired with a numbered diagram, since an orphan table has no gap to police.
func pairedTables(blocks []mermaid.DocBlock) map[*mermaid.CalloutTableBlock]bool {
	paired := map[*mermaid.CalloutTableBlock]bool{}
	for _, pair := range mermaid.PairDiagramsWithTables(blocks).Pairs {
		if pair.Table != nil {
			paired[pair.Table] = true
		}
	}
	return paired
}

// precedingDiagram returns the callout table at blocks[i] and the mermaid block before it, the legend or else the diagram. Both are nil when blocks[i] is not a diagram-preceded table.
func precedingDiagram(blocks []mermaid.DocBlock, i int, paired map[*mermaid.CalloutTableBlock]bool) (*mermaid.DiagramBlock, *mermaid.CalloutTableBlock) {
	table, ok := blocks[i].(*mermaid.CalloutTableBlock)
	if !ok || i == 0 || !paired[table] {
		return nil, nil
	}
	prev, _ := blocks[i-1].(*mermaid.DiagramBlock)
	return prev, table
}

// gapErrors reports each non-blank, non-caption line between prev's closing
// fence and the table header — the prose that belongs in the table instead.
func gapErrors(file string, lines []string, prev *mermaid.DiagramBlock, table *mermaid.CalloutTableBlock) []string {
	var errors []string
	for ln := closingFenceLine(lines, prev.Line) + 1; ln < table.Line; ln++ {
		text := strings.TrimSpace(lines[ln-1])
		if text == "" || captionLine.MatchString(text) {
			continue
		}
		errors = append(errors, fmt.Sprintf(
			"%s: prose between the diagram and its callout table (line %d) — move detail into the callout table",
			file, ln))
	}
	return errors
}

// closingFenceLine returns the 1-based line of the fence closing the diagram opened at openLine, as pkg/mermaid does. An unterminated block closes at end of file.
func closingFenceLine(lines []string, openLine int) int {
	open := openLine - 1
	marker := "```"
	if m := fenceMarker.FindStringSubmatch(lines[open]); m != nil {
		marker = m[1]
	}
	for i := open + 1; i < len(lines); i++ {
		if strings.HasPrefix(diagrams.TrimStartJS(lines[i]), marker) {
			return i + 1
		}
	}
	return len(lines)
}
