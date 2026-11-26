package gherkin

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// Up numbers the doc's bullet acceptances, naming each for its first WHEN,
// and leaves commas and names to a person. Each ID a bullet doc skips becomes
// a Removed line, so every acceptance keeps its ID.
func Up(root, path string) error {
	version := removedVersion(root)
	return rewrite(root, path, true, func(d requirement.Document) string { return numbered(d, version) })
}

// numbered renders d's acceptances in ID order, each a number and name over
// one GIVEN, its WHENs and a THEN under each.
func numbered(d requirement.Document, version string) string {
	byNumber := map[int]requirement.IDLine{}
	last := 0
	for _, a := range d.Acceptances {
		_, n, _ := strings.Cut(a.ID, ".")
		i, _ := strconv.Atoi(n)
		byNumber[i], last = a, max(last, i)
	}
	var out string
	taken := map[string]bool{requirement.RemovedName: true}
	for i := 1; i <= last; i++ {
		out += acceptance(i, byNumber, taken, version)
	}
	return out
}

// gapReason is why up writes a Removed line for an ID the bullets skip.
const gapReason = "absent before the move to numbered steps"

// acceptance renders acceptance number i, or a Removed line for a skipped
// number.
func acceptance(i int, byNumber map[int]requirement.IDLine, taken map[string]bool, version string) string {
	a, ok := byNumber[i]
	if !ok {
		return fmt.Sprintf("%d. %s in %s: %s\n", i, requirement.RemovedName, version, gapReason)
	}
	pad := strings.Repeat(" ", len(strconv.Itoa(i))+2)
	out := fmt.Sprintf("%d. %s\n%s1. GIVEN %s\n", i, name(a, taken), pad, a.Given)
	for n, st := range a.Steps {
		out += fmt.Sprintf("%s   %d. WHEN %s\n%s      1. THEN %s\n", pad, n+1, st.When, pad, st.Then)
	}
	return out
}

// name is a's first WHEN without a leading keyword, opening with a capital,
// for a person to reword. A name taken before adds the first THEN, then a
// count, and taken records it.
func name(a requirement.IDLine, taken map[string]bool) string {
	first := a.Steps[0]
	base := capital(withoutKeyword(first.When))
	n := base
	if taken[n] {
		n = base + " then " + first.Then
	}
	for i := 2; taken[n]; i++ {
		n = fmt.Sprintf("%s then %s %d", base, first.Then, i)
	}
	taken[n] = true
	return n
}

// withoutKeyword is text without a leading Gherkin keyword in any case, which
// a name may not open with.
func withoutKeyword(text string) string {
	word, rest, _ := strings.Cut(text, " ")
	if keywords[strings.ToUpper(word)] && rest != "" {
		return rest
	}
	return text
}

// keywords are the Gherkin step keywords.
var keywords = map[string]bool{"GIVEN": true, "WHEN": true, "THEN": true, "AND": true, "BUT": true}

// capital is text opening with a capital.
func capital(text string) string {
	r, size := utf8.DecodeRuneInString(text)
	return string(unicode.ToUpper(r)) + text[size:]
}
