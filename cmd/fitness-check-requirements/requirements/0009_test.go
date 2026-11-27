package requirements

import (
	"strconv"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-checks/req"
)

// chain is 0001.1, named Adds, with two WHENs that each give their THEN.
func chain(when1, then1, when2, then2 string) string {
	return "1. Adds\n   1. GIVEN two numbers\n" +
		"      1. " + when1 + "\n         1. " + then1 + "\n" +
		"      2. " + when2 + "\n         1. " + then2 + "\n"
}

// tagsFor is adder_test.go with a tagged test for each acceptance of the doc that reqs makes.
func tagsFor(reqs string) string {
	doc, _ := requirement.Parse(rootDoc, acceptance.DocWith(reqs))
	body := "package adder\n\nimport (\n\t\"testing\"\n\n\t\"github.com/may-journal/fitness-checks/req\"\n)\n"
	for _, a := range doc.Acceptances {
		body += "\nfunc Test" + strings.Replace(a.ID, ".", "_", 1) + "(t *testing.T) {\n"
		for _, s := range a.Steps {
			for _, th := range s.Thens {
				body += "\treq.Then(t, " + strconv.Quote(th.ID) + ", " + strconv.Quote(th.Hash()) + ", func(a req.A) { a.Equal(1, 1) })\n"
			}
		}
		body += "}\n"
	}
	return body
}

// passes fails unless requirements passes reqs, whose second acceptance chains
// two steps.
func passes(t *testing.T, reqs string) {
	t.Helper()
	out, code := judged(t, gherkinRepo(t, reqs, tagsFor(reqs)))
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
}

func Test0009_1(t *testing.T) {
	t.Parallel()
	t.Run("0009.1.1 sibling WHENs sharing three meaningful words name them for the GIVEN", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0009.1.1.1.1", "85699bcd", func(_ req.A) {
			gherkinFails(t, chain("WHEN I add two small numbers", "THEN I get their sum", "WHEN I add two large numbers", "THEN I get a total"),
				"19: WHENs on lines 19 and 21 open with the same words, `I add two`; move those words into the GIVEN")
		})
	})
	t.Run("0009.1.2 shared filler or fewer than three words pass", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0009.1.1.2.1", "07b8e753", func(_ req.A) {
			passes(t, sum+"2. Doubles a sum\n   1. GIVEN two numbers\n"+
				"      1. WHEN a number is small\n         1. THEN it adds fast\n"+
				"      2. WHEN a number is large\n         1. THEN it adds slowly\n")
			passes(t, sum+"2. Doubles a sum\n   1. GIVEN two numbers\n"+
				"      1. WHEN I add them\n         1. THEN it adds fast\n"+
				"      2. WHEN I add three\n         1. THEN it adds slowly\n")
		})
	})
}

func Test0009_2(t *testing.T) {
	t.Parallel()
	t.Run("0009.2.1 a THEN that repeats names both lines and tells the author what to do", func(t *testing.T) {
		t.Parallel()
		group := "22: this THEN repeats the THEN on line 20; group the WHENs that give it under one step and state it once"
		req.Then(t, "0009.2.1.1.1", "5703ca02", func(_ req.A) {
			gherkinFails(t, chain("WHEN I add them", "THEN I get their sum", "WHEN a sum is doubled", "THEN I  get their sum"), group)
			gherkinFails(t, chain("WHEN I add them", "THEN I get their sum", "WHEN a sum is doubled", "THEN i get their Sum."), group)
		})
		req.Then(t, "0009.2.1.1.2", "e14e5488", func(_ req.A) {
			gherkinFails(t, sum+"         2. THEN i get their sum.\n", "21: this THEN repeats the THEN on line 20; drop the repeated THEN")
		})
	})
	t.Run("0009.2.2 the same THEN in two acceptances passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0009.2.1.2.1", "07b8e753", func(_ req.A) {
			passes(t, twoAcceptances)
		})
	})
}
