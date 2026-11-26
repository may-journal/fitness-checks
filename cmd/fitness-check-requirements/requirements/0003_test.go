package requirements

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// twoAcceptances is a Requirements section holding 0001.1, with two GIVENs,
// and 0001.2, a chain of two steps whose first has two THENs, each named.
const twoAcceptances = "1. Adds on a calculator\n" +
	"   1. GIVEN two numbers\n   2. GIVEN a calculator\n      1. WHEN I add them\n         1. THEN I get their sum\n" +
	"2. Doubles a sum\n" +
	"   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n         2. THEN it is a number\n" +
	"      2. WHEN I double it\n         1. THEN it is even\n"

// oneAcceptance is 0001.1, named Adds, with given, when and then as its
// numbered items.
func oneAcceptance(given, when, then string) string {
	return "1. Adds\n   1. " + given + "\n      1. " + when + "\n         1. " + then + "\n"
}

// firstThen and secondThens name the case of each THEN of twoAcceptances,
// which holds 0001.1 beside 0001.2.
const firstThen = "0001.1.2.1.1 I get their sum"

var secondThens = []string{"0001.2.1.1.1 I get their sum", "0001.2.1.1.2 it is a number", "0001.2.1.2.1 it is even"}

// firstOnly is adder_test.go proving 0001.1 alone, with its one case.
var firstOnly = acceptance.CaseFiles("adder", acceptance.CaseFunc("Test0001_1", firstThen))

// provingTests is adder_test.go proving 0001.1, and 0001.2 with the cases
// given, or one for each of its THENs when none are.
func provingTests(cases ...string) string {
	if len(cases) == 0 {
		cases = secondThens
	}
	return acceptance.CaseFiles("adder", acceptance.CaseFunc("Test0001_1", firstThen), acceptance.CaseFunc("Test0001_2", cases...))
}

// gherkinRepo is a Go repo whose root requirement doc holds the Requirements
// section reqs, proven by tests.
func gherkinRepo(t *testing.T, reqs, tests string) string {
	t.Helper()
	return acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
		".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
		rootDoc:           acceptance.DocWith(reqs),
		"adder_test.go":   tests,
	}))
}

// judged runs requirements on repo with args, returning what it printed as a
// person reads it.
func judged(t *testing.T, repo string, args ...string) (string, int) {
	t.Helper()
	out, code := acceptance.Fitness(t, repo, nil, append([]string{"requirements"}, args...)...)
	return acceptance.Read(out), code
}

// mustFail fails unless requirements fails reqs naming want.
func mustFail(t *testing.T, reqs, want string) {
	t.Helper()
	out, code := judged(t, gherkinRepo(t, reqs, acceptance.TestFile("adder", "Test0001_1")))
	acceptance.Sees(t, out, code, 1, rootDoc+":"+want)
}

// sum is 0001.1's usual numbered acceptance.
var sum = oneAcceptance("GIVEN two numbers", "WHEN I add them", "THEN I get their sum")

func Test0003_1(t *testing.T) {
	t.Parallel()
	t.Run("0003.1.1.1.1 it pairs steps with tests", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, gherkinRepo(t, twoAcceptances, provingTests()))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// A THEN without its case fails.
		out, code = judged(t, gherkinRepo(t, twoAcceptances, provingTests(secondThens[:2]...)))
		acceptance.Sees(t, out, code, 1, "THEN 0001.2.1.2.1 has no test case")
	})
	t.Run("0003.1.1.2.1 it names the keyword", func(t *testing.T) {
		t.Parallel()
		mustFail(t, oneAcceptance("Given two numbers", "WHEN I add them", "THEN I get their sum"), "18: write the keyword GIVEN in caps")
		mustFail(t, oneAcceptance("GIVEN two numbers", "WHEN I add them", "then I get their sum"), "20: write the keyword THEN in caps")
	})
	t.Run("0003.1.1.3.1 it asks for the next GIVEN or THEN", func(t *testing.T) {
		t.Parallel()
		mustFail(t, strings.Replace(sum, "      1. WHEN", "   2. AND a calculator\n      1. WHEN", 1), "19: write no AND or BUT; an extra condition is the next GIVEN and an extra result the next THEN")
		mustFail(t, sum+"         2. BUT not their product\n", "21: write no AND or BUT")
	})
	t.Run("0003.1.1.4.1 it asks to split", func(t *testing.T) {
		t.Parallel()
		mustFail(t, oneAcceptance("GIVEN two numbers, a calculator", "WHEN I add them", "THEN I get their sum"), "18: split this item at its comma")
		// A comma inside code belongs to the code.
		out, code := judged(t, gherkinRepo(t, oneAcceptance("GIVEN two numbers such as `1, 2`", "WHEN I add them", "THEN I get their sum"), acceptance.TestFile("adder", "Test0001_1")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.1.1.5.1 it names the need", func(t *testing.T) {
		t.Parallel()
		mustFail(t, strings.Replace(sum, "1. Adds", "- Adds", 1), "17: number each acceptance and name it, `1. Name`; the numbered form holds no dashes")
		mustFail(t, strings.Replace(sum, "1. Adds", "2. Adds", 1), "17: number this acceptance `1.`; acceptances number from 1 in order and 0001.1 takes its ID from it")
		mustFail(t, strings.Replace(sum, "Adds", "0001.1", 1), "17: name acceptance 0001.1 in words; its ID comes from its number")
		mustFail(t, strings.Replace(sum, "Adds", "GIVEN it adds", 1), "17: name acceptance 0001.1 in words; its GIVEN, WHEN and THEN items nest under the name")
		mustFail(t, strings.Replace(sum, "Adds", "Adds, sums", 1), "17: name acceptance 0001.1 without a comma")
		mustFail(t, sum+strings.Replace(sum, "1. Adds", "2. Adds", 1), "21: 0001.2 has the name of 0001.1; give each acceptance its own name")
		mustFail(t, strings.Replace(sum, "1. WHEN", "- WHEN", 1), "19: number each GIVEN, WHEN and THEN item; the numbered form holds no dashes")
		mustFail(t, strings.Replace(sum, "1. THEN", "2. THEN", 1), "20: number this THEN `1.`; each level numbers from 1 in order")
		mustFail(t, sum+"   2. GIVEN a calculator\n", "21: nest each GIVEN under the acceptance name, before its first WHEN")
		mustFail(t, strings.Replace(sum, "      1. WHEN", "   2. WHEN", 1), "19: nest each WHEN under the last GIVEN")
		mustFail(t, sum+"      1. THEN it is a number\n", "21: nest each THEN under its WHEN")
		mustFail(t, "1. Adds\n   1. GIVEN two numbers\n      1. WHEN I add them\n", "19: this WHEN needs a THEN nested under it")
		mustFail(t, "1. Adds\n   1. GIVEN two numbers\n", "17: 0001.1 needs a WHEN nested under its last GIVEN")
	})
	t.Run("0003.1.1.6.1 it fails at that line", func(t *testing.T) {
		t.Parallel()
		mustFail(t, strings.Replace(sum, "      1. WHEN", "   2. a calculator\n      1. WHEN", 1), "19: Gherkin: expected:")
	})
}

func docHolds(t *testing.T, repo, want string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(repo, rootDoc))
	acceptance.MustDo(t, err)
	if !strings.Contains(string(b), want) {
		t.Errorf("%s is missing:\n%s\nit holds:\n%s", rootDoc, want, b)
	}
}

// movedRun commits a bullet doc, then stages the
// doc with the Requirements section reqs and runs requirements on the staged
// change.
func movedRun(t *testing.T, reqs string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
		".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
		rootDoc:           acceptance.DocWith("- 0001.1\n  - Given two numbers\n    - When I add them\n      - Then I get their sum\n"),
	}))
	acceptance.Write(t, repo, map[string]string{rootDoc: acceptance.DocWith(reqs)})
	acceptance.Git(t, repo, "add", "-A")
	return judged(t, repo)
}

func Test0003_3(t *testing.T) {
	t.Parallel()
	moved := "1. Adds on a calculator\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n"
	t.Run("0003.3.1.1.1 its tests need no change", func(t *testing.T) {
		t.Parallel()
		out, code := movedRun(t, moved)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// A move that also changes what a step says still needs its test changed.
		out, code = movedRun(t, strings.Replace(moved, "their sum", "their product", 1))
		acceptance.Sees(t, out, code, 1, rootDoc)
	})
}

// committed commits twoAcceptances with their tests, then commits the doc
// with each Requirements section of edits in turn, and runs requirements.
func committed(t *testing.T, edits ...string) (string, int) {
	t.Helper()
	repo := gherkinRepo(t, twoAcceptances, provingTests())
	for _, reqs := range edits {
		acceptance.Commit(t, repo, map[string]string{rootDoc: acceptance.DocWith(reqs)}, "docs(requirements): edit acceptances")
	}
	return judged(t, repo)
}

// swapped is twoAcceptances with its two acceptances in each other's place.
func swapped() string {
	first, second, _ := strings.Cut(twoAcceptances, "2. Doubles a sum\n")
	return "1. Doubles a sum\n" + second + strings.Replace(first, "1. Adds", "2. Adds", 1)
}

func Test0003_4(t *testing.T) {
	t.Parallel()
	t.Run("0003.4.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := committed(t, strings.Replace(twoAcceptances, "Adds on a calculator", "Adds with a calculator", 1))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.4.1.2.1 it names both IDs", func(t *testing.T) {
		t.Parallel()
		out, code := committed(t, swapped())
		acceptance.Sees(t, out, code, 1,
			rootDoc+":17: 0001.1 holds the acceptance 0001.2 held before; a numbered acceptance keeps its place, so number a new one after the last",
			"0001.2 holds the acceptance 0001.1 held before")
	})
	t.Run("0003.4.1.3.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := committed(t, swapped(), twoAcceptances)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

// firstRemoved is twoAcceptances with 0001.1 deleted as a Removed line.
var firstRemoved = "1. Removed in 2026.11.25.1230: folded into 0001.2\n" + twoAcceptances[strings.Index(twoAcceptances, "2. Doubles"):]

// onlySecond is adder_test.go proving 0001.2 alone, a case per THEN.
var onlySecond = acceptance.CaseFiles("adder", acceptance.CaseFunc("Test0001_2", secondThens...))

// removedRun commits twoAcceptances with their tests, then commits the doc
// with each Requirements section of edits in turn beside tests, and runs
// requirements.
func removedRun(t *testing.T, tests string, edits ...string) (string, int) {
	t.Helper()
	repo := gherkinRepo(t, twoAcceptances, provingTests())
	for _, reqs := range edits {
		acceptance.Commit(t, repo, map[string]string{rootDoc: acceptance.DocWith(reqs), "adder_test.go": tests}, "docs(requirements): edit acceptances")
	}
	return judged(t, repo)
}

func Test0003_5(t *testing.T) {
	t.Parallel()
	t.Run("0003.5.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := removedRun(t, onlySecond, firstRemoved)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.5.1.2.1 it fails", func(t *testing.T) {
		t.Parallel()
		mustFail(t, "1. Removed in 2026.11.25.1230: folded\n   1. GIVEN two numbers\n", "18: 0001.1 is Removed, so it holds no GIVEN, WHEN or THEN")
		// A Removed line needs a version that reads as a date and a reason.
		for _, bad := range []string{"Removed", "Removed: folded", "Removed in 2026.13.40.9999: folded", "Removed in 2026.11.25.1230:", "Removed in 2026.11.25: folded"} {
			mustFail(t, "1. "+bad+"\n", "17: write 0001.1 as `1. Removed in YYYY.MM.DD.HHMM: <reason>`, with the changelog version that removes it")
		}
		mustFail(t, "1. Removed in 2026.11.25.1230: folded, merged\n", "17: name acceptance 0001.1 without a comma")
		// A name that only opens with the word stays an acceptance.
		out, code := judged(t, gherkinRepo(t, strings.Replace(sum, "Adds", "Removed in bulk", 1), acceptance.TestFile("adder", "Test0001_1")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.5.1.3.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := removedRun(t, provingTests(), firstRemoved)
		acceptance.Sees(t, out, code, 1, "Test0001_1")
	})
	t.Run("0003.5.1.4.1 it fails as reused", func(t *testing.T) {
		t.Parallel()
		back := strings.Replace(twoAcceptances, "Adds on a calculator", "Adds again", 1)
		out, code := removedRun(t, provingTests(), firstRemoved, back)
		acceptance.Sees(t, out, code, 1, rootDoc+":17: 0001.1 was deleted before; IDs are never reused")
	})
	t.Run("0003.5.1.5.1 it passes", func(t *testing.T) {
		t.Parallel()
		appended := firstRemoved + "3. Adds on a calculator\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n"
		tests := acceptance.CaseFiles("adder", acceptance.CaseFunc("Test0001_3", "0001.3.1.1.1 I get their sum"), acceptance.CaseFunc("Test0001_2", secondThens...))
		out, code := removedRun(t, tests, appended)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

// atCap is a Requirements section of ten numbered items: eight acceptances,
// the default cap, and Removed lines at 3 and 7.
func atCap() (string, string) {
	reqs := ""
	var funcs []string
	for n := 1; n <= 10; n++ {
		if n == 3 || n == 7 {
			reqs += fmt.Sprintf("%d. Removed in 2026.11.25.1230: folded into 0001.%d\n", n, n+1)
			continue
		}
		reqs += fmt.Sprintf("%d. Adds pair %d\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n", n, n)
		funcs = append(funcs, acceptance.CaseFunc(fmt.Sprintf("Test0001_%d", n), fmt.Sprintf("0001.%d.1.1.1 I get their sum", n)))
	}
	return reqs, acceptance.CaseFiles("adder", funcs...)
}

func Test0003_6(t *testing.T) {
	t.Parallel()
	t.Run("0003.6.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		reqs, tests := atCap()
		repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
			".fitnessrc.json": `{"policy": "external", "checks": ["requirements", "prose-budget"]}` + "\n",
			rootDoc:           acceptance.DocWith(reqs),
			"adder_test.go":   tests,
		}))
		for _, check := range []string{"requirements", "prose-budget"} {
			out, code := acceptance.Fitness(t, repo, nil, check)
			acceptance.Sees(t, acceptance.Read(out), code, 0, check+" passed")
		}
	})
}

func Test0003_7(t *testing.T) {
	t.Parallel()
	t.Run("0003.7.1.1.1 it fails naming the flag", func(t *testing.T) {
		t.Parallel()
		repo := gherkinRepo(t, twoAcceptances, provingTests())
		out, code := judged(t, repo, "--up", "0003.2-gherkin")
		acceptance.Sees(t, out, code, 1, "unknown flag --up; the only flag is --fix")
	})
}

// onlyFirst is twoAcceptances without 0001.2, the last.
var onlyFirst = twoAcceptances[:strings.Index(twoAcceptances, "2. Doubles")]

func Test0003_8(t *testing.T) {
	t.Parallel()
	t.Run("0003.8.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := removedRun(t, firstOnly, onlyFirst+"2. Removed in 2026.11.25.1230: folded into 0001.1\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.8.1.2.1 it names the `Removed` line", func(t *testing.T) {
		t.Parallel()
		out, code := removedRun(t, firstOnly, onlyFirst)
		acceptance.Sees(t, out, code, 1, rootDoc+":1: 0001.2 was an acceptance and this doc no longer holds it; a removed acceptance keeps its number, so add `2. Removed in YYYY.MM.DD.HHMM: <reason>` to Requirements")
	})
}
