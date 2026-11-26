package requirement

import (
	"cmp"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	gherkin "github.com/cucumber/gherkin/go/v36"
	messages "github.com/cucumber/messages/go/v30"
)

// The keywords a numbered acceptance writes in caps.
const (
	givenWord = "GIVEN"
	whenWord  = "WHEN"
	thenWord  = "THEN"
)

// joinWords are the Gherkin conjunctions a numbered acceptance never writes.
var joinWords = map[string]bool{"AND": true, "BUT": true}

// kind is what a numbered acceptance item states.
type kind int

const (
	other kind = iota
	given
	when
	then
	// joined is an AND or BUT item, already failed.
	joined
)

// kindOf maps each keyword to the kind of item it opens.
var kindOf = map[string]kind{givenWord: given, whenWord: when, thenWord: then}

// readWords is the Gherkin keyword each kind of item reads as, and joins
// marks the kinds whose later items read as And.
var (
	readWords = map[kind]string{given: "Given", when: "When", then: "Then"}
	joins     = map[kind]bool{given: true, then: true}
)

// joinedError is what an AND or BUT item fails with.
const joinedError = "write no AND or BUT; an extra condition is the next GIVEN and an extra result the next THEN"

// commaError is what a step holding a comma outside code fails with.
const commaError = "split this item at its comma; an extra condition is the next GIVEN and an extra result the next THEN"

// gherkinWord matches a Gherkin word between spaces in any casing. A step reads
// as one condition or result, so a word not in caps hides a second clause.
var gherkinWord = regexp.MustCompile(`(?i)(?:^|\s)(given|when|then|and|but)(?:\s|$)`)

// quoted matches text in straight or curly double quotes, which a step
// repeats as written.
var quoted = regexp.MustCompile(`"[^"]*"|“[^”]*”`)

// parserError matches one Cucumber parse error, "(line:column): message".
var parserError = regexp.MustCompile(`\((\d+):\d+\): (.*)`)

// codeSpan matches inline code, whose commas belong to the code.
var codeSpan = regexp.MustCompile("`[^`]*`")

// idLike matches a name that is an acceptance ID, which a numbered
// acceptance takes from its number instead.
var idLike = regexp.MustCompile(`^\d{4}\.\d`)

// hasGherkin reports whether a Requirements section numbers its
// acceptances, as its first top-level item shows.
func hasGherkin(s section) bool {
	for i := range s.lines {
		if top, ok := listItem(s.lines[i]); ok && top.indent == 0 {
			return opensNumbered(top, s.lines[i+1:])
		}
	}
	return false
}

// opensNumbered reports whether top is numbered, or a dash whose next item,
// opening rest, is a GIVEN in caps.
func opensNumbered(top item, rest []line) bool {
	return top.numbered || len(rest) > 0 && numberedChild(rest[0])
}

// numberedChild reports whether l is a nested GIVEN in caps.
func numberedChild(l line) bool {
	it, ok := listItem(l)
	return ok && it.indent > 0 && strings.HasPrefix(it.text, givenWord+" ")
}

// parseGherkin checks each numbered acceptance and reads it as Gherkin
// through the Cucumber parser.
func parseGherkin(p *problems, docID string, s section) []IDLine {
	checkBlockKeywords(p, s)
	groups := groupItems(p, s.lines, topLevelMessage)
	kinds := map[int]kind{}
	doc := parsed(p, s.n, gherkinText(p, s.n, groups, kinds))
	at := func(l *messages.Location) int { return s.n + int(l.Line) - 1 }
	byLine := scenarios(doc, at)
	var out []IDLine
	named := map[string]string{}
	for i, g := range groups {
		a := checkName(p, docID, g[0], i+1, named)
		if IsRemoved(a.Name) {
			checkRemoved(p, a, g)
			continue
		}
		if sc := byLine[a.Line]; checkNumbered(p, a.ID, g, kinds) && sc != nil {
			a = fromScenario(a, sc, at)
		}
		out = append(out, a)
	}
	return out
}

// RemovedName opens a deleted acceptance's item, `N. Removed in <version>:
// <reason>`, which keeps its place so later numbers never shift. Its ID
// counts as deleted, so it is no acceptance and counts toward no limit.
const RemovedName = "Removed"

// VersionLayout is a changelog entry's version, such as 2026.11.25.1230, as
// a Go time layout.
const VersionLayout = "2006.01.02.1504"

// removedLine matches a well-formed Removed item: its version and reason.
var removedLine = regexp.MustCompile(`^Removed in (\d{4}\.\d{2}\.\d{2}\.\d{4}): (\S.*)$`)

// removedStart matches a name that marks its acceptance deleted, well formed
// or not. A name such as "Removed in bulk" stays an acceptance.
var removedStart = regexp.MustCompile(`^Removed(?:$|:| in \d)`)

// IsRemoved reports whether an acceptance's name marks it deleted.
func IsRemoved(name string) bool {
	return removedStart.MatchString(name)
}

// checkRemoved fails a Removed item without a readable version and a
// reason, or with items nested under it.
func checkRemoved(p *problems, a IDLine, g []item) {
	if !removedWell(a.Name) {
		_, n, _ := strings.Cut(a.ID, ".")
		p.add(a.Line, "write "+a.ID+" as `"+n+". Removed in YYYY.MM.DD.HHMM: <reason>`, with the changelog version that removes it")
	}
	if len(g) > 1 {
		p.add(g[1].n, a.ID+" is Removed, so it holds no GIVEN, WHEN or THEN")
	}
}

// removedWell reports whether a Removed name holds a version that reads as a
// date and a reason.
func removedWell(name string) bool {
	m := removedLine.FindStringSubmatch(name)
	if m == nil {
		return false
	}
	_, err := time.Parse(VersionLayout, m[1])
	return err == nil
}

// scenarios maps each parsed Scenario to its doc line, or none when the
// parser rejected the doc.
func scenarios(doc *messages.GherkinDocument, at func(*messages.Location) int) map[int]*messages.Scenario {
	byLine := map[int]*messages.Scenario{}
	if doc == nil {
		return byLine
	}
	for _, c := range doc.Feature.Children {
		if c.Scenario != nil {
			byLine[at(c.Scenario.Location)] = c.Scenario
		}
	}
	return byLine
}

// parsed is text as the Cucumber parser reads it, or nil. It fails each
// line the parser rejects at its doc line, counting from heading line n.
func parsed(p *problems, n int, text string) *messages.GherkinDocument {
	doc, err := gherkin.ParseGherkinDocument(strings.NewReader(text), (&gherkinIDs{}).next)
	if err == nil {
		return doc
	}
	for _, m := range parserError.FindAllStringSubmatch(err.Error(), -1) {
		at, _ := strconv.Atoi(m[1])
		p.add(n+at-1, "Gherkin: "+m[2])
	}
	return nil
}

// gherkinIDs numbers the nodes the parser builds.
type gherkinIDs struct{ n int }

func (g *gherkinIDs) next() string {
	g.n++
	return strconv.Itoa(g.n)
}

// checkName requires the top-level item at place pos to be numbered pos and
// named in words of its own. Its ID is the doc's ID and pos, so item 2 of
// 0003 is 0003.2.
func checkName(p *problems, docID string, top item, pos int, named map[string]string) IDLine {
	a := IDLine{ID: fmt.Sprintf("%s.%d", docID, pos), Name: top.text, Line: top.n}
	if msg := nameProblem(top, pos, a.ID); msg != "" {
		p.add(top.n, msg)
	}
	if was, ok := named[a.Name]; ok && !IsRemoved(a.Name) {
		p.add(a.Line, a.ID+" has the name of "+was+"; give each acceptance its own name")
	}
	named[a.Name] = a.ID
	return a
}

// nameProblem is what is wrong with the number or name of acceptance id, at
// place pos, or empty.
func nameProblem(top item, pos int, id string) string {
	switch {
	case !top.numbered:
		return fmt.Sprintf("number each acceptance and name it, `%d. Name`; the numbered form holds no dashes", pos)
	case top.num != pos:
		return fmt.Sprintf("number this acceptance `%d.`; acceptances number from 1 in order and %s takes its ID from it", pos, id)
	}
	return wordingProblem(top.text, id)
}

// wordingProblem is what is wrong with the name of acceptance id, or empty.
func wordingProblem(name, id string) string {
	word, _ := cutWord(name)
	switch {
	case idLike.MatchString(name):
		return "name acceptance " + id + " in words; its ID comes from its number"
	case isKeyword(word):
		return "name acceptance " + id + " in words; its GIVEN, WHEN and THEN items nest under the name"
	case hasComma(name):
		return "name acceptance " + id + " without a comma"
	}
	return ""
}

// isKeyword reports whether word is a Gherkin step keyword in any case.
func isKeyword(word string) bool {
	upper := strings.ToUpper(word)
	return joinWords[upper] || kindOf[upper] != other
}

// hasComma reports whether text holds a comma outside code.
func hasComma(text string) bool {
	return strings.Contains(codeSpan.ReplaceAllString(text, ""), ",")
}

// joinWord is the first Gherkin word of text, outside code and quotes, that
// is not in caps, or "" when it holds none.
func joinWord(text string) string {
	for _, m := range gherkinWord.FindAllStringSubmatch(quoted.ReplaceAllString(codeSpan.ReplaceAllString(text, ""), ""), -1) {
		if m[1] != strings.ToUpper(m[1]) {
			return m[1]
		}
	}
	return ""
}

// gherkinText is the acceptances as the Cucumber parser reads them, each on
// its doc line below a Feature line on heading line n. It fails each item's
// wording and files its kind by line.
func gherkinText(p *problems, n int, groups [][]item, kinds map[int]kind) string {
	last := groups[len(groups)-1]
	out := make([]string, last[len(last)-1].n-n+1)
	out[0] = "Feature: requirements"
	for _, g := range groups {
		out[g[0].n-n] = "Scenario: " + g[0].text
		prev := other
		for _, it := range g[1:] {
			word, rest := cutWord(it.text)
			k := checkItem(p, it.n, word, rest)
			kinds[it.n] = k
			out[it.n-n] = readAs(k, prev, word) + rest
			prev = k
		}
	}
	return strings.Join(out, "\n") + "\n"
}

// cutWord splits an item's text into its first word and the rest, the rest
// keeping its leading space.
func cutWord(text string) (string, string) {
	word, rest, found := strings.Cut(text, " ")
	if found {
		rest = " " + rest
	}
	return word, rest
}

// checkItem fails an AND or BUT, or the wording of a step. It returns the
// item's kind.
func checkItem(p *problems, n int, word, rest string) kind {
	upper := strings.ToUpper(word)
	if joinWords[upper] {
		p.add(n, joinedError)
		return joined
	}
	k := kindOf[upper]
	if k == other {
		return other
	}
	itemWording(p, n, word, rest)
	return k
}

// itemWording fails a keyword not in caps, a comma outside code, or a
// lowercase "and" or "but" outside code and quotes.
func itemWording(p *problems, n int, word, rest string) {
	if upper := strings.ToUpper(word); word != upper {
		p.add(n, "write the keyword "+upper+" in caps")
	}
	if hasComma(rest) {
		p.add(n, commaError)
	}
	if w := joinWord(rest); w != "" {
		p.add(n, "write the Gherkin word "+strings.ToUpper(w)+" in caps, or split it into its own GIVEN, WHEN or THEN")
	}
}

// readAs is the Gherkin keyword for kind k after kind prev. A failed AND or
// BUT or a repeated GIVEN or THEN reads as And; any other item keeps its word.
func readAs(k, prev kind, word string) string {
	if k == joined || k == prev && joins[k] {
		return "And"
	}
	return cmp.Or(readWords[k], word)
}

// fromScenario fills a's Given and steps from its parsed Scenario. GIVENs
// join with "and" into one Given, as each step's THENs do into one Then.
func fromScenario(a IDLine, sc *messages.Scenario, at func(*messages.Location) int) IDLine {
	var givens []string
	var thens [][]string
	var thenItems [][]Then
	for _, st := range sc.Steps {
		switch {
		case st.KeywordType == messages.StepKeywordType_ACTION:
			a.Steps = append(a.Steps, Step{Line: at(st.Location), When: st.Text})
			thens = append(thens, nil)
			thenItems = append(thenItems, nil)
		case len(a.Steps) == 0:
			givens = append(givens, st.Text)
			a.GivenLine = cmp.Or(a.GivenLine, at(st.Location))
		default:
			i := len(thens) - 1
			thens[i] = append(thens[i], st.Text)
			id := fmt.Sprintf("%s.%d.%d.%d", a.ID, len(givens), i+1, len(thens[i]))
			thenItems[i] = append(thenItems[i], Then{ID: id, Line: at(st.Location), Text: st.Text})
		}
	}
	a.Given = strings.Join(givens, " and ")
	a = numberSteps(a, thens)
	for i := range a.Steps {
		a.Steps[i].Thens = thenItems[i]
	}
	return a
}

// numberSteps sets each step's Then. One WHEN is the acceptance itself, and
// a chain numbers its steps, so step 2 of 0003.6 is 0003.6.2.
func numberSteps(a IDLine, thens [][]string) IDLine {
	for i := range a.Steps {
		a.Steps[i].Then = strings.Join(thens[i], " and ")
		if len(a.Steps) > 1 {
			a.Steps[i].ID = fmt.Sprintf("%s.%d", a.ID, i+1)
		}
	}
	return a
}

// blockKeyword matches a Gherkin block keyword and its colon opening a line of
// a Requirements section, after any list number. Group 1 is the keyword.
var blockKeyword = regexp.MustCompile(`(?i)^\s*(?:\d+\.\s+|[-*+]\s+)?(feature|rule|background|scenario outline|scenario|examples|example):`)

// checkBlockKeywords fails each Gherkin block keyword that opens a line and
// is not in caps. It is the caps rule for the keywords no step opens with.
func checkBlockKeywords(p *problems, s section) {
	for _, l := range s.lines {
		if m := blockKeyword.FindStringSubmatch(l.text); m != nil && m[1] != strings.ToUpper(m[1]) {
			p.add(l.n, "write the keyword "+strings.ToUpper(m[1])+" in caps")
		}
	}
}
