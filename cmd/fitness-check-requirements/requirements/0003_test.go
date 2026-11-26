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

// provingTests is adder_test.go proving 0001.1, and 0001.2 with one case per
// step named.
func provingTests(steps ...string) string {
	body := "package adder\n\nimport \"testing\"\n\nfunc Test0001_1(t *testing.T) {}\n\nfunc Test0001_2(t *testing.T) {\n"
	for _, s := range steps {
		body += "\tt.Run(\"" + s + "\", func(t *testing.T) {})\n"
	}
	return body + "}\n"
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

// gherkinFails fails unless requirements fails reqs naming want.
func gherkinFails(t *testing.T, reqs, want string) {
	t.Helper()
	out, code := judged(t, gherkinRepo(t, reqs, acceptance.TestFile("adder", "Test0001_1")))
	acceptance.Sees(t, out, code, 1, rootDoc+":"+want)
}

// sum is 0001.1's usual numbered acceptance.
var sum = oneAcceptance("GIVEN two numbers", "WHEN I add them", "THEN I get their sum")

func Test0003_1(t *testing.T) {
	t.Parallel()
	t.Run("0003.1.1 every keyword in caps pairs like bullets", func(t *testing.T) {
		t.Parallel()
		out, code := judged(t, gherkinRepo(t, twoAcceptances, provingTests("0001.2.1 add", "0001.2.2 double")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// A step without its test case fails as it would for bullets.
		out, code = judged(t, gherkinRepo(t, twoAcceptances, provingTests("0001.2.1 add")))
		acceptance.Sees(t, out, code, 1, "0001.2.2")
	})
	t.Run("0003.1.2 a keyword not in caps fails", func(t *testing.T) {
		t.Parallel()
		gherkinFails(t, oneAcceptance("Given two numbers", "WHEN I add them", "THEN I get their sum"), "18: write the keyword GIVEN in caps")
		gherkinFails(t, oneAcceptance("GIVEN two numbers", "WHEN I add them", "then I get their sum"), "20: write the keyword THEN in caps")
	})
	t.Run("0003.1.3 an AND or BUT asks for the next GIVEN or THEN", func(t *testing.T) {
		t.Parallel()
		gherkinFails(t, strings.Replace(sum, "      1. WHEN", "   2. AND a calculator\n      1. WHEN", 1), "19: write no AND or BUT; an extra condition is the next GIVEN and an extra result the next THEN")
		gherkinFails(t, sum+"         2. BUT not their product\n", "21: write no AND or BUT")
	})
	t.Run("0003.1.4 a comma outside code asks to split the item", func(t *testing.T) {
		t.Parallel()
		gherkinFails(t, oneAcceptance("GIVEN two numbers, a calculator", "WHEN I add them", "THEN I get their sum"), "18: split this item at its comma")
		// A comma inside code belongs to the code.
		out, code := judged(t, gherkinRepo(t, oneAcceptance("GIVEN two numbers such as `1, 2`", "WHEN I add them", "THEN I get their sum"), acceptance.TestFile("adder", "Test0001_1")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.1.5 a dash or misnumbered or misplaced or misnamed item names what it needs", func(t *testing.T) {
		t.Parallel()
		gherkinFails(t, strings.Replace(sum, "1. Adds", "- Adds", 1), "17: number each acceptance and name it, `1. Name`; the numbered form holds no dashes")
		gherkinFails(t, strings.Replace(sum, "1. Adds", "2. Adds", 1), "17: number this acceptance `1.`; acceptances number from 1 in order and 0001.1 takes its ID from it")
		gherkinFails(t, strings.Replace(sum, "Adds", "0001.1", 1), "17: name acceptance 0001.1 in words; its ID comes from its number")
		gherkinFails(t, strings.Replace(sum, "Adds", "GIVEN it adds", 1), "17: name acceptance 0001.1 in words; its GIVEN, WHEN and THEN items nest under the name")
		gherkinFails(t, strings.Replace(sum, "Adds", "Adds, sums", 1), "17: name acceptance 0001.1 without a comma")
		gherkinFails(t, sum+strings.Replace(sum, "1. Adds", "2. Adds", 1), "21: 0001.2 has the name of 0001.1; give each acceptance its own name")
		gherkinFails(t, strings.Replace(sum, "1. WHEN", "- WHEN", 1), "19: number each GIVEN, WHEN and THEN item; the numbered form holds no dashes")
		gherkinFails(t, strings.Replace(sum, "1. THEN", "2. THEN", 1), "20: number this THEN `1.`; each level numbers from 1 in order")
		gherkinFails(t, sum+"   2. GIVEN a calculator\n", "21: nest each GIVEN under the acceptance name, before its first WHEN")
		gherkinFails(t, strings.Replace(sum, "      1. WHEN", "   2. WHEN", 1), "19: nest each WHEN under the last GIVEN")
		gherkinFails(t, sum+"      1. THEN it is a number\n", "21: nest each THEN under its WHEN")
		gherkinFails(t, "1. Adds\n   1. GIVEN two numbers\n      1. WHEN I add them\n", "19: this WHEN needs a THEN nested under it")
		gherkinFails(t, "1. Adds\n   1. GIVEN two numbers\n", "17: 0001.1 needs a WHEN nested under its last GIVEN")
	})
	t.Run("0003.1.6 an item the parser rejects fails at its line", func(t *testing.T) {
		t.Parallel()
		gherkinFails(t, strings.Replace(sum, "      1. WHEN", "   2. a calculator\n      1. WHEN", 1), "19: Gherkin: expected:")
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

// movedRun commits a bullet doc whose Given holds a comma. It stages the
// doc with the Requirements section reqs and runs requirements on the staged
// change.
func movedRun(t *testing.T, reqs string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
		".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
		rootDoc:           acceptance.DocWith("- 0001.1\n  - Given two numbers, a calculator\n    - When I add them\n      - Then I get their sum\n"),
	}))
	acceptance.Write(t, repo, map[string]string{rootDoc: acceptance.DocWith(reqs)})
	acceptance.Git(t, repo, "add", "-A")
	return judged(t, repo)
}

func Test0003_3(t *testing.T) {
	t.Parallel()
	moved := "1. Adds on a calculator\n   1. GIVEN two numbers\n   2. GIVEN a calculator\n      1. WHEN I add them\n         1. THEN I get their sum\n"
	out, code := movedRun(t, moved)
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	// A move that also changes what a step says still needs its test changed.
	out, code = movedRun(t, strings.Replace(moved, "their sum", "their product", 1))
	acceptance.Sees(t, out, code, 1, rootDoc)
}

// committed commits twoAcceptances with their tests, then commits the doc
// with each Requirements section of edits in turn, and runs requirements.
func committed(t *testing.T, edits ...string) (string, int) {
	t.Helper()
	repo := gherkinRepo(t, twoAcceptances, provingTests("0001.2.1 add", "0001.2.2 double"))
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
	t.Run("0003.4.1 a rename in place passes", func(t *testing.T) {
		t.Parallel()
		out, code := committed(t, strings.Replace(twoAcceptances, "Adds on a calculator", "Adds with a calculator", 1))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.4.2 a move to another number fails naming both IDs", func(t *testing.T) {
		t.Parallel()
		out, code := committed(t, swapped())
		acceptance.Sees(t, out, code, 1,
			rootDoc+":17: 0001.1 holds the acceptance 0001.2 held before; a numbered acceptance keeps its place, so number a new one after the last",
			"0001.2 holds the acceptance 0001.1 held before")
	})
	t.Run("0003.4.3 a move back passes again", func(t *testing.T) {
		t.Parallel()
		out, code := committed(t, swapped(), twoAcceptances)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

// firstRemoved is twoAcceptances with 0001.1 deleted as a Removed line.
var firstRemoved = "1. Removed in 2026.11.25.1230: folded into 0001.2\n" + twoAcceptances[strings.Index(twoAcceptances, "2. Doubles"):]

// onlySecond is adder_test.go proving 0001.2 alone, a case per step.
var onlySecond = "package adder\n\nimport \"testing\"\n\nfunc Test0001_2(t *testing.T) {\n" +
	"\tt.Run(\"0001.2.1 add\", func(t *testing.T) {})\n\tt.Run(\"0001.2.2 double\", func(t *testing.T) {})\n}\n"

// removedRun commits twoAcceptances with their tests, then commits the doc
// with each Requirements section of edits in turn beside tests, and runs
// requirements.
func removedRun(t *testing.T, tests string, edits ...string) (string, int) {
	t.Helper()
	repo := gherkinRepo(t, twoAcceptances, provingTests("0001.2.1 add", "0001.2.2 double"))
	for _, reqs := range edits {
		acceptance.Commit(t, repo, map[string]string{rootDoc: acceptance.DocWith(reqs), "adder_test.go": tests}, "docs(requirements): edit acceptances")
	}
	return judged(t, repo)
}

func Test0003_5(t *testing.T) {
	t.Parallel()
	t.Run("0003.5.1 a Removed line keeps later numbers", func(t *testing.T) {
		t.Parallel()
		out, code := removedRun(t, onlySecond, firstRemoved)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.5.2 a Removed line holding items or no version and reason fails", func(t *testing.T) {
		t.Parallel()
		gherkinFails(t, "1. Removed in 2026.11.25.1230: folded\n   1. GIVEN two numbers\n", "18: 0001.1 is Removed, so it holds no GIVEN, WHEN or THEN")
		// A Removed line needs a version that reads as a date and a reason.
		for _, bad := range []string{"Removed", "Removed: folded", "Removed in 2026.13.40.9999: folded", "Removed in 2026.11.25.1230:", "Removed in 2026.11.25: folded"} {
			gherkinFails(t, "1. "+bad+"\n", "17: write 0001.1 as `1. Removed in YYYY.MM.DD.HHMM: <reason>`, with the changelog version that removes it")
		}
		gherkinFails(t, "1. Removed in 2026.11.25.1230: folded, merged\n", "17: name acceptance 0001.1 without a comma")
		// A name that only opens with the word stays an acceptance.
		out, code := judged(t, gherkinRepo(t, strings.Replace(sum, "Adds", "Removed in bulk", 1), acceptance.TestFile("adder", "Test0001_1")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.5.3 a test naming a Removed acceptance fails", func(t *testing.T) {
		t.Parallel()
		out, code := removedRun(t, provingTests("0001.2.1 add", "0001.2.2 double"), firstRemoved)
		acceptance.Sees(t, out, code, 1, "Test0001_1")
	})
	t.Run("0003.5.4 a Removed line turned back fails as reuse", func(t *testing.T) {
		t.Parallel()
		back := strings.Replace(twoAcceptances, "Adds on a calculator", "Adds again", 1)
		out, code := removedRun(t, provingTests("0001.2.1 add", "0001.2.2 double"), firstRemoved, back)
		acceptance.Sees(t, out, code, 1, rootDoc+":17: 0001.1 was deleted before; IDs are never reused")
	})
	t.Run("0003.5.5 a Removed acceptance removed then appended in one commit is no move", func(t *testing.T) {
		t.Parallel()
		appended := firstRemoved + "3. Adds on a calculator\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n"
		tests := strings.Replace(onlySecond, "func Test0001_2", "func Test0001_3(t *testing.T) {}\n\nfunc Test0001_2", 1)
		out, code := removedRun(t, tests, appended)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

// atCap is a Requirements section of ten numbered items: eight acceptances,
// the default cap, and Removed lines at 3 and 7.
func atCap() (string, string) {
	reqs, tests := "", "package adder\n\nimport \"testing\"\n"
	for n := 1; n <= 10; n++ {
		if n == 3 || n == 7 {
			reqs += fmt.Sprintf("%d. Removed in 2026.11.25.1230: folded into 0001.%d\n", n, n+1)
			continue
		}
		reqs += fmt.Sprintf("%d. Adds pair %d\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n", n, n)
		tests += fmt.Sprintf("\nfunc Test0001_%d(t *testing.T) {}\n", n)
	}
	return reqs, tests
}

func Test0003_6(t *testing.T) {
	t.Parallel()
	reqs, tests := atCap()
	repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
		".fitnessrc.json": `{"policy": "external", "checks": ["requirements", "prose-budget"]}` + "\n",
		rootDoc:           acceptance.DocWith(reqs),
		"adder_test.go":   tests,
	}))
	for i, check := range []string{"requirements", "prose-budget"} {
		t.Run(fmt.Sprintf("0003.6.%d %s passes at the cap", i+1, check), func(t *testing.T) {
			t.Parallel()
			out, code := acceptance.Fitness(t, repo, nil, check)
			acceptance.Sees(t, acceptance.Read(out), code, 0, check+" passed")
		})
	}
}

func Test0003_7(t *testing.T) {
	t.Parallel()
	repo := gherkinRepo(t, twoAcceptances, provingTests("0001.2.1 add", "0001.2.2 double"))
	out, code := judged(t, repo, "--up", "0003.2-gherkin")
	acceptance.Sees(t, out, code, 1, "unknown flag --up; the only flag is --fix")
}

// onlyFirst is twoAcceptances without 0001.2, the last.
var onlyFirst = twoAcceptances[:strings.Index(twoAcceptances, "2. Doubles")]

func Test0003_8(t *testing.T) {
	t.Parallel()
	t.Run("0003.8.1 the last acceptance as a Removed line passes", func(t *testing.T) {
		t.Parallel()
		out, code := removedRun(t, acceptance.TestFile("adder", "Test0001_1"), onlyFirst+"2. Removed in 2026.11.25.1230: folded into 0001.1\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0003.8.2 the last acceptance deleted outright fails naming the Removed line to add", func(t *testing.T) {
		t.Parallel()
		out, code := removedRun(t, acceptance.TestFile("adder", "Test0001_1"), onlyFirst)
		acceptance.Sees(t, out, code, 1, rootDoc+":1: 0001.2 was an acceptance and this doc no longer holds it; a removed acceptance keeps its number, so add `2. Removed in YYYY.MM.DD.HHMM: <reason>` to Requirements")
	})
}

func Test0003_9(t *testing.T) {
	t.Parallel()
	t.Run("0003.9.1 a lowercase and or but outside code or quotes asks to split", func(t *testing.T) {
		t.Parallel()
		gherkinFails(t, oneAcceptance("GIVEN two numbers", "WHEN I add them", "THEN each part stands alone and code spans drop out"), `20: split this item at its "and"`)
		gherkinFails(t, oneAcceptance("GIVEN two numbers", "WHEN I add them", "THEN it is a sum but not a product"), `20: split this item at its "but"`)
	})
	t.Run("0003.9.2 the word in code or in quotes passes", func(t *testing.T) {
		t.Parallel()
		reqs := oneAcceptance("GIVEN two numbers such as `1 and 2`", "WHEN I add them", `THEN I get the word "and" or “but”`)
		out, code := judged(t, gherkinRepo(t, reqs, acceptance.TestFile("adder", "Test0001_1")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}
