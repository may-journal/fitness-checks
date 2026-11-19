// The shared parser lives in pkg/mermaid, which keeps its styling patterns unexported, so this package duplicates them.

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
	program.Main(docSource, diagrams.Check("mermaid-diagram-prose", mdx.WithMDX, mermaid.RunDocCheck, validateDoc))
}

// jsWS is the JavaScript \s character class as an RE2 class body.
const jsWS = diagrams.JSWhitespace

// trimJS trims JavaScript whitespace from both ends — String.prototype.trim.
var trimJS = diagrams.TrimJS

// calloutThenProse matches a callout, a number or capital letters, followed
// by prose, e.g. `6 Uses` or `A calls` — the duplication #20 targets.
var calloutThenProse = regexp.MustCompile(`^(\d+|[A-Z]+)[` + jsWS + `]+[A-Za-z]`)

// isCalloutProse reports a label that starts with a node number or a
// connector letter that follows the doubling, then carries prose.
func isCalloutProse(label string) bool {
	m := calloutThenProse.FindStringSubmatch(label)
	return m != nil && (m[1][0] <= '9' || mermaid.IsLetterCallout(m[1]))
}

// relCall matches C4 relationship calls whose label argument may carry prose.
var relCall = regexp.MustCompile(`\b(?:Bi)?Rel\w*\([^)]*\)`)

// quoted matches quoted argument content.
var quoted = regexp.MustCompile(`["']([^"']*)["']`)

// pipeLabel matches a flowchart pipe edge label: `-->|6 uses|`.
var pipeLabel = regexp.MustCompile(`\|[` + jsWS + `]*([^|]*?)[` + jsWS + `]*\|`)

// styleLine matches classDef/style lines, which carry non-label numbers;
// skip them for colon labels.
var styleLine = regexp.MustCompile(`^[` + jsWS + `]*(classDef|style|linkStyle|%%)`)

// colonLabel matches an edge label after a colon: `A --> B : 6 uses`.
var colonLabel = regexp.MustCompile(`:[` + jsWS + `]*((?:\d+|[A-Z]+)[` + jsWS + `]+[A-Za-z][^|]*?)[` + jsWS + `]*$`)

// relCallProse returns prose labels from C4 `Rel(...)` call arguments.
func relCallProse(body string) []string {
	var found []string
	for _, rel := range relCall.FindAllString(body, -1) {
		for _, arg := range quoted.FindAllStringSubmatch(rel, -1) {
			if label := trimJS(arg[1]); isCalloutProse(label) {
				found = append(found, label)
			}
		}
	}
	return found
}

// pipeProse returns prose labels from a line's flowchart pipe edge labels
// (`-->|…|`).
func pipeProse(line string) []string {
	var found []string
	for _, pipe := range pipeLabel.FindAllStringSubmatch(line, -1) {
		if label := trimJS(pipe[1]); isCalloutProse(label) {
			found = append(found, label)
		}
	}
	return found
}

// colonProse returns the prose label from a line's colon edge label
// (`A --> B : 6 uses`), if any.
func colonProse(line string) []string {
	colon := colonLabel.FindStringSubmatch(line)
	if colon == nil {
		return nil
	}
	if label := trimJS(colon[1]); isCalloutProse(label) {
		return []string{label}
	}
	return nil
}

// lineEdgeProse returns prose labels from a single diagram line's edge
// labels (styling lines skipped).
func lineEdgeProse(line string) []string {
	if styleLine.MatchString(line) {
		return nil
	}
	return append(pipeProse(line), colonProse(line)...)
}

// proseLabels returns relationship/edge labels in a diagram body that carry
// prose beyond a callout.
func proseLabels(diagram *mermaid.DiagramBlock) []string {
	found := relCallProse(diagram.Body)
	for _, line := range strings.Split(diagram.Body, "\n") {
		found = append(found, lineEdgeProse(line)...)
	}
	return found
}

// validateDoc flags edge labels that carry prose beyond a callout in every
// diagram but a legend.
func validateDoc(file, content string) []string {
	errors := []string{}
	for _, diagram := range mermaid.Judged(mermaid.ParseDoc(content)) {
		for _, label := range proseLabels(diagram) {
			errors = append(errors, fmt.Sprintf(
				`%s: diagram at line %d has a relationship label with prose ("%s") — put descriptions in the callout table`,
				file, diagram.Line, label))
		}
	}
	return errors
}
