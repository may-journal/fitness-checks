package requirements

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func Test0003_2(t *testing.T) {
	t.Parallel()
	t.Run("0003.2.1 fix or up turns bullets into numbered items", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
			".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
		}))
		out, code := judged(t, repo, "--fix")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		docHolds(t, repo, "## Requirements\n\n"+strings.Replace(sum, "Adds", "I add them", 1))
		// A chain keeps its step numbers, and each ID its acceptance number.
		chained := gherkinRepo(t, acceptance.Acceptances("0001.1")+acceptance.Chain("0001.2", 2), provingTests("0001.2.1 add", "0001.2.2 double"))
		out, code = judged(t, chained, "--up", "0003.2-gherkin")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// A name an acceptance before took adds the first THEN.
		docHolds(t, chained, "2. I add them then I get their sum\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN ")
		docHolds(t, chained, "      2. WHEN I double that sum\n         1. THEN ")
		// An acceptance commented out stays out of the list.
		commented := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
			".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
			rootDoc:           "<!-- kept -->\n" + acceptance.RequirementDoc("0001.1") + "<!--\n- 0001.9\n  - Given two numbers\n    - When I add them\n      - Then I get their sum\n-->\n",
		}))
		out, code = judged(t, commented, "--fix")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		docHolds(t, commented, "         1. THEN I get their sum\n")
		// A comment outside the Requirements section stays.
		docHolds(t, commented, "<!-- kept -->\n")
		if b, _ := os.ReadFile(filepath.Join(commented, rootDoc)); strings.Contains(string(b), "0001.9") {
			t.Errorf("--fix brought back a commented-out acceptance:\n%s", b)
		}
	})
	t.Run("0003.2.2 down turns numbered items back into bullets", func(t *testing.T) {
		t.Parallel()
		repo := gherkinRepo(t, twoAcceptances, provingTests("0001.2.1 add", "0001.2.2 double"))
		out, code := judged(t, repo, "--down", "0003.2-gherkin")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		docHolds(t, repo, "- 0001.1\n  - Given two numbers and a calculator\n    - When I add them\n      - Then I get their sum\n"+
			"- 0001.2\n  - Given two numbers\n    1. When I add them\n       - Then I get their sum and it is a number\n    2. When I double it\n       - Then it is even\n")
	})
	t.Run("0003.2.3 a flag naming no migration fails listing them", func(t *testing.T) {
		t.Parallel()
		repo := gherkinRepo(t, twoAcceptances, provingTests("0001.2.1 add", "0001.2.2 double"))
		out, code := judged(t, repo, "--down")
		acceptance.Sees(t, out, code, 1, "name the migration --down runs, one of [0003.2-gherkin]")
		out, code = judged(t, repo, "--up", "0003.9-gherkin")
		acceptance.Sees(t, out, code, 1, "no migration \"0003.9-gherkin\"; run one of [0003.2-gherkin]")
		out, code = judged(t, repo, "--fix", "--up", "0003.2-gherkin")
		acceptance.Sees(t, out, code, 1, "run --up or --fix alone")
		out, code = judged(t, repo, "--up=0003.2-gherkin")
		acceptance.Sees(t, out, code, 1, "write --up 0003.2-gherkin, with a space")
		out, code = judged(t, repo, "--up", "0003.2-gherkin", "--down", "0003.2-gherkin")
		acceptance.Sees(t, out, code, 1, "run one migration at a time")
	})
}

// docHolds fails unless repo's root requirement doc holds want.
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
	t.Run("0003.5.5 up writes a Removed line for each skipped ID", func(t *testing.T) {
		t.Parallel()
		repo := gherkinRepo(t, acceptance.Acceptances("0001.2"), acceptance.TestFile("adder", "Test0001_2"))
		out, code := judged(t, repo, "--up", "0003.2-gherkin")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// With no new changelog entry, the line carries the run's time.
		b, err := os.ReadFile(filepath.Join(repo, rootDoc))
		acceptance.MustDo(t, err)
		if !regexp.MustCompile(`## Requirements\n\n1\. Removed in \d{4}\.\d{2}\.\d{2}\.\d{4}: absent before the move to numbered steps\n2\. I add them\n`).Match(b) {
			t.Errorf("%s has no Removed line with the run's version:\n%s", rootDoc, b)
		}
		// A new changelog entry lends the line its version.
		entry := gherkinRepo(t, acceptance.Acceptances("0001.2"), acceptance.TestFile("adder", "Test0001_2"))
		acceptance.Write(t, entry, map[string]string{"changelog/2026-12-01-0930-numbered.md": "### 2026.12.01.0930\n\n> Numbered.\n"})
		out, code = judged(t, entry, "--up", "0003.2-gherkin")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		docHolds(t, entry, "1. Removed in 2026.12.01.0930: absent before the move to numbered steps\n")
	})
	t.Run("0003.5.6 down skips a Removed ID", func(t *testing.T) {
		t.Parallel()
		repo := gherkinRepo(t, firstRemoved, onlySecond)
		out, code := judged(t, repo, "--down", "0003.2-gherkin")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		docHolds(t, repo, "## Requirements\n\n- 0001.2\n  - Given two numbers\n    1. When I add them\n")
	})
	t.Run("0003.5.7 a Removed acceptance appended again is no move", func(t *testing.T) {
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
	for _, check := range []string{"requirements", "prose-budget"} {
		out, code := acceptance.Fitness(t, repo, nil, check)
		acceptance.Sees(t, acceptance.Read(out), code, 0, check+" passed")
	}
}
