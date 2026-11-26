package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func Test0008_1(t *testing.T) {
	t.Parallel()
	t.Run("0008.1.1 a Gherkin word not in caps asks for caps or a split", func(t *testing.T) {
		t.Parallel()
		for word, then := range map[string]string{
			"and":   "THEN each part stands alone and code spans drop out",
			"but":   "THEN it is a sum but not a product",
			"when":  "THEN it passes when no doc names it",
			"then":  "THEN it adds them then doubles the sum",
			"given": "THEN a Given starts it",
			"AND":   "THEN an And starts it",
		} {
			want := strings.ToUpper(strings.TrimSpace(strings.Fields(word)[0]))
			gherkinFails(t, oneAcceptance("GIVEN two numbers", "WHEN I add them", then), "20: write the Gherkin word "+want+" in caps, or split it into its own GIVEN, WHEN or THEN")
		}
	})
	t.Run("0008.1.2 the word in code or in quotes or in caps passes", func(t *testing.T) {
		t.Parallel()
		reqs := oneAcceptance("GIVEN two numbers such as `1 and 2`", "WHEN I add a GIVEN or an AND", `THEN I get the word "and" or “but” from read-and-write input`)
		out, code := judged(t, gherkinRepo(t, reqs, acceptance.TestFile("adder", "Test0001_1")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0008_2(t *testing.T) {
	t.Parallel()
	t.Run("0008.2.1 a block keyword not in caps asks for caps", func(t *testing.T) {
		t.Parallel()
		for _, word := range []string{"Feature", "Rule", "Background", "Scenario", "Scenario Outline", "Example", "Examples"} {
			gherkinFails(t, strings.Replace(sum, "1. Adds", "1. "+word+": Adds", 1), "17: write the keyword "+strings.ToUpper(word)+" in caps")
		}
	})
	t.Run("0008.2.2 the keyword in caps or without a colon passes", func(t *testing.T) {
		t.Parallel()
		for _, name := range []string{"SCENARIO: Adds", "Examples add"} {
			repo := gherkinRepo(t, strings.Replace(sum, "1. Adds", "1. "+name, 1), acceptance.TestFile("adder", "Test0001_1"))
			out, code := judged(t, repo)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		}
	})
}
