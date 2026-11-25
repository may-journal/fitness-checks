package requirement

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	gherkin "github.com/cucumber/gherkin/go/v36"
	messages "github.com/cucumber/messages/go/v30"
)

// gherkinFence opens a requirement doc's Gherkin block, and fence closes it.
const (
	gherkinFence = "```gherkin"
	fence        = "```"
)

// stepWords are the step keywords a doc writes in caps, mapped to the case the
// Cucumber parser reads. Only letter case changes, so every column holds.
var stepWords = map[string]string{"GIVEN": "Given", "WHEN": "When", "THEN": "Then", "AND": "And", "BUT": "But"}

// scenarioWord opens each acceptance, in caps, and the parser reads it as
// scenarioRead.
const (
	scenarioWord = "SCENARIO:"
	scenarioRead = "Scenario:"
)

// parserError matches one Cucumber parse error, "(line:column): message".
var parserError = regexp.MustCompile(`\((\d+):\d+\): (.*)`)

// codeSpan matches inline code, whose commas belong to the code.
var codeSpan = regexp.MustCompile("`[^`]*`")

// hasGherkin reports whether a Requirements section holds a Gherkin block.
func hasGherkin(s section) bool {
	for _, l := range s.lines {
		if strings.TrimSpace(l.text) == gherkinFence {
			return true
		}
	}
	return false
}

// parseGherkin checks the caps, indents and commas of a section's one Gherkin
// block, then has the Cucumber parser read it.
func parseGherkin(p *problems, docID string, s section) []IDLine {
	block, open, ok := gherkinLines(p, s)
	if !ok {
		return nil
	}
	text := checkWords(p, open, block)
	doc, err := gherkin.ParseGherkinDocument(strings.NewReader(text), (&gherkinIDs{}).next)
	if err != nil {
		for _, m := range parserError.FindAllStringSubmatch(err.Error(), -1) {
			n, _ := strconv.Atoi(m[1])
			p.add(open+n-1, "Gherkin: "+m[2])
		}
		return nil
	}
	return scenarios(p, docID, open, doc)
}

// gherkinIDs numbers the nodes the parser builds.
type gherkinIDs struct{ n int }

func (g *gherkinIDs) next() string {
	g.n++
	return strconv.Itoa(g.n)
}

// gherkinLines returns the lines between the section's Gherkin fences and the
// fence's line, failing any line outside the block.
func gherkinLines(p *problems, s section) ([]line, int, bool) {
	open := fenceAt(s.lines, 0, gherkinFence)
	closed := fenceAt(s.lines, open+1, fence)
	if closed < 0 {
		p.add(s.lines[open].n, "close the Gherkin block with ```")
		return nil, s.lines[open].n, false
	}
	for i, l := range s.lines {
		if i < open || i > closed {
			p.add(l.n, "Requirements holds only its one Gherkin block")
		}
	}
	return s.lines[open+1 : closed], s.lines[open].n, true
}

// fenceAt is the index of the first line from index from that is want, or -1.
func fenceAt(lines []line, from int, want string) int {
	for i := from; i < len(lines); i++ {
		if strings.TrimSpace(lines[i].text) == want {
			return i
		}
	}
	return -1
}

// checkWords checks each line of the block opened on line open. It returns the
// text the parser reads, each doc line on its own line under a Feature line.
func checkWords(p *problems, open int, block []line) string {
	out := make([]string, block[len(block)-1].n-open+1)
	out[0] = "Feature: requirements"
	stepIndent := -1
	for _, l := range block {
		out[l.n-open], stepIndent = checkLine(p, l, stepIndent)
	}
	return strings.Join(out, "\n") + "\n"
}

// checkLine checks one block line and returns it as the parser reads it, with
// the indent of the step now open.
func checkLine(p *problems, l line, stepIndent int) (string, int) {
	word, rest, indent := splitWord(l.text)
	read, isStep := checkWord(p, l.n, word)
	text := strings.Repeat(" ", indent) + read + rest
	if !isStep {
		return text, stepIndent
	}
	checkComma(p, l.n, rest)
	if read != "And" && read != "But" {
		return text, indent
	}
	if indent <= stepIndent {
		p.add(l.n, fmt.Sprintf("indent %s one level deeper than the step it extends", word))
	}
	return text, stepIndent
}

// checkComma fails a step whose text holds a comma outside code.
func checkComma(p *problems, n int, rest string) {
	if strings.Contains(codeSpan.ReplaceAllString(rest, ""), ",") {
		p.add(n, "split this step at its comma into an AND line")
	}
}

// splitWord splits a block line into its first word, the rest, and its indent.
func splitWord(text string) (word, rest string, indent int) {
	trimmed := strings.TrimLeft(text, " ")
	indent = len(text) - len(trimmed)
	word, after, found := strings.Cut(trimmed, " ")
	if found {
		rest = " " + after
	}
	return word, rest, indent
}

// checkWord fails a keyword written in any case but caps. It returns the word
// as the parser reads it, and whether it opens a step.
func checkWord(p *problems, n int, word string) (string, bool) {
	upper := strings.ToUpper(word)
	if upper == scenarioWord {
		if word != scenarioWord {
			p.add(n, "write "+scenarioWord+" in caps")
		}
		return scenarioRead, false
	}
	read, ok := stepWords[upper]
	if !ok {
		return word, false
	}
	if word != upper {
		p.add(n, "write the keyword "+upper+" in caps")
	}
	return read, true
}

// scenarios turns each parsed SCENARIO into an acceptance, mapping parser
// lines back to doc lines.
func scenarios(p *problems, docID string, open int, doc *messages.GherkinDocument) []IDLine {
	if doc.Feature == nil {
		return nil
	}
	at := func(l *messages.Location) int { return open + int(l.Line) - 1 }
	var out []IDLine
	seen := map[string]bool{}
	for _, c := range doc.Feature.Children {
		if c.Scenario == nil {
			p.add(open, "the Gherkin block holds only SCENARIO acceptances")
			continue
		}
		a := scenario(p, docID, c.Scenario, at)
		if seen[a.ID] {
			p.add(a.Line, a.ID+" appears twice; each acceptance needs its own ID")
		}
		seen[a.ID] = true
		out = append(out, a)
	}
	return out
}

// scenario checks one SCENARIO: named by its acceptance ID, then one GIVEN,
// then WHEN and THEN pairs.
func scenario(p *problems, docID string, s *messages.Scenario, at func(*messages.Location) int) IDLine {
	a := IDLine{ID: strings.TrimSpace(s.Name), Line: at(s.Location)}
	if m := acceptID.FindStringSubmatch(a.ID); m == nil || m[1] != docID {
		p.add(a.Line, "name each SCENARIO by its acceptance ID `"+docID+".N`")
	}
	steps, ok := gherkinSteps(s.Steps)
	switch {
	case !ok:
		p.add(a.Line, a.ID+" needs one GIVEN, then WHEN and THEN pairs; AND lines may follow any of them")
	case len(steps.pairs) > p.maxSteps:
		p.add(a.Line, fmt.Sprintf("%s chains %d steps, at most %d; move the rest to a new ID", a.ID, len(steps.pairs), p.maxSteps))
	default:
		a.Given, a.GivenLine = steps.given.Text, at(steps.given.Location)
		a.Steps = pairSteps(a.ID, steps.pairs, at)
	}
	return a
}

// shape is a scenario's GIVEN and the WHEN that opens each step.
type shape struct {
	given *messages.Step
	pairs []*messages.Step
}

// gherkinSteps reads a GIVEN, then WHEN and THEN pairs, skipping AND and BUT
// lines, and reports whether the steps come in that order.
func gherkinSteps(steps []*messages.Step) (shape, bool) {
	var sh shape
	var kinds []messages.StepKeywordType
	for _, st := range steps {
		if st.KeywordType == messages.StepKeywordType_CONJUNCTION {
			continue
		}
		kinds = append(kinds, st.KeywordType)
		switch st.KeywordType {
		case messages.StepKeywordType_CONTEXT:
			sh.given = st
		case messages.StepKeywordType_ACTION:
			sh.pairs = append(sh.pairs, st)
		}
	}
	return sh, wellOrdered(kinds)
}

// wellOrdered reports whether kinds are one Context, then Action and Outcome
// pairs.
func wellOrdered(kinds []messages.StepKeywordType) bool {
	if len(kinds) < 3 || len(kinds)%2 == 0 {
		return false
	}
	want := []messages.StepKeywordType{messages.StepKeywordType_CONTEXT}
	for range len(kinds) / 2 {
		want = append(want, messages.StepKeywordType_ACTION, messages.StepKeywordType_OUTCOME)
	}
	return slices.Equal(kinds, want)
}

// pairSteps names each WHEN a step. One pair is the acceptance itself, and a
// chain numbers its steps, so step 2 of 0003.6 is 0003.6.2.
func pairSteps(id string, whens []*messages.Step, at func(*messages.Location) int) []Step {
	var steps []Step
	for i, w := range whens {
		st := Step{Line: at(w.Location), When: w.Text}
		if len(whens) > 1 {
			st.ID = fmt.Sprintf("%s.%d", id, i+1)
		}
		steps = append(steps, st)
	}
	return steps
}
