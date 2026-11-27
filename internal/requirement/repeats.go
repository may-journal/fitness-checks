package requirement

import (
	"fmt"
	"regexp"
	"strings"
)

// minSharedWords is how many leading words, not counting filler, sibling
// WHENs may share before they read as one condition said twice. A run of
// filler such as "a name is" names no condition, so it never counts.
const minSharedWords = 3

// filler are the words that carry no condition by themselves.
var filler = map[string]bool{
	"a": true, "an": true, "the": true, "it": true, "its": true, "is": true,
	"are": true, "has": true, "have": true, "no": true, "not": true,
	"in": true, "of": true, "to": true,
}

// wordPattern matches a code span or a run of non-space characters, so a code
// span counts as one word.
var wordPattern = regexp.MustCompile("`[^`]*`|\\S+")

// checkRepeats fails an acceptance that says one thing twice: WHENs that
// open with the same words, and THENs with the same text. It reads only an
// acceptance that has the numbered shape.
func checkRepeats(p *problems, g []item, kinds map[int]kind) {
	var whens []item
	var thens []sameThen
	for _, it := range g[1:] {
		switch kinds[it.n] {
		case when:
			whens = append(whens, it)
		case then:
			thens = append(thens, sameThen{it, len(whens)})
		}
	}
	checkAlikeWhens(p, whens)
	checkRepeatedThens(p, thens)
}

// words is an item's words after its keyword.
func words(it item) []string {
	_, rest := cutWord(it.text)
	return wordPattern.FindAllString(rest, -1)
}

// alike is a run of WHENs that share their leading words.
type alike struct {
	prefix []string
	whens  []item
}

// checkAlikeWhens fails each run of WHENs that share leading words worth
// naming, with the words they all share and each line.
func checkAlikeWhens(p *problems, whens []item) {
	for _, r := range alikeRuns(whens) {
		if len(r.whens) > 1 {
			p.add(r.whens[0].n, fmt.Sprintf("WHENs on lines %s open with the same words, `%s`; move those words into the GIVEN and let each WHEN say only what differs", lineList(r.whens), strings.Join(r.prefix, " ")))
		}
	}
}

// alikeRuns deals each WHEN into the first run whose shared words it keeps
// worth naming, or into a run of its own.
func alikeRuns(whens []item) []*alike {
	var runs []*alike
	for _, w := range whens {
		ws := words(w)
		r := runSharing(runs, ws)
		if r == nil {
			runs = append(runs, &alike{prefix: ws, whens: []item{w}})
			continue
		}
		r.prefix = commonPrefix(r.prefix, ws)
		r.whens = append(r.whens, w)
	}
	return runs
}

// runSharing is the first run that ws shares words worth naming with, or nil.
func runSharing(runs []*alike, ws []string) *alike {
	for _, r := range runs {
		if worthNaming(commonPrefix(r.prefix, ws)) {
			return r
		}
	}
	return nil
}

// commonPrefix is the leading words a and b share.
func commonPrefix(a, b []string) []string {
	n := 0
	for n < len(a) && n < len(b) && a[n] == b[n] {
		n++
	}
	return a[:n]
}

// worthNaming reports whether shared holds minSharedWords words that are not
// filler.
func worthNaming(shared []string) bool {
	n := 0
	for _, w := range shared {
		if !filler[strings.ToLower(w)] {
			n++
		}
	}
	return n >= minSharedWords
}

// lineList is the line numbers of g, in order, joined with commas and "and".
func lineList(g []item) string {
	var ns []string
	for _, it := range g {
		ns = append(ns, fmt.Sprint(it.n))
	}
	if len(ns) == 1 {
		return ns[0]
	}
	return strings.Join(ns[:len(ns)-1], ", ") + " and " + ns[len(ns)-1]
}

// sameThen is a THEN with the WHEN it sits under.
type sameThen struct {
	item
	parent int
}

// thenKey is a THEN's words in lower case with whitespace collapsed and
// trailing punctuation dropped, so changing case or punctuation hides no
// repeat.
func thenKey(t item) string {
	text := strings.Join(strings.Fields(strings.TrimPrefix(t.text, thenWord)), " ")
	return strings.ToLower(strings.TrimRight(text, ".,;:!?"))
}

// repeatMessage tells the author to drop a THEN repeated under its own WHEN.
// A THEN repeated under another WHEN means those WHENs belong together.
func repeatMessage(first, again sameThen) string {
	if first.parent == again.parent {
		return fmt.Sprintf("this THEN repeats the THEN on line %d; drop the repeated THEN", first.n)
	}
	return fmt.Sprintf("this THEN repeats the THEN on line %d; group the WHENs that give it under one step and state it once", first.n)
}

// checkRepeatedThens fails each THEN whose text, ignoring case, whitespace
// and trailing punctuation, matches an earlier THEN of the acceptance,
// naming both lines.
func checkRepeatedThens(p *problems, thens []sameThen) {
	seen := map[string]sameThen{}
	for _, t := range thens {
		key := thenKey(t.item)
		if first, ok := seen[key]; ok {
			p.add(t.n, repeatMessage(first, t))
			continue
		}
		seen[key] = t
	}
}
