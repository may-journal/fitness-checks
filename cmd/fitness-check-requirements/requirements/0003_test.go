package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// twoScenarios is a Gherkin block body holding 0001.1, whose GIVEN gains an
// AND line, and 0001.2, a chain of two steps.
const twoScenarios = "SCENARIO: 0001.1\n" +
	"  GIVEN two numbers\n    AND a calculator\n  WHEN I add them\n  THEN I get their sum\n\n" +
	"SCENARIO: 0001.2\n" +
	"  GIVEN two numbers\n  WHEN I add them\n  THEN I get their sum\n  WHEN I double it\n  THEN it is even\n"

// provingTests is adder_test.go proving 0001.1, and 0001.2 with one case per
// step named.
func provingTests(steps ...string) string {
	body := "package adder\n\nimport \"testing\"\n\nfunc Test0001_1(t *testing.T) {}\n\nfunc Test0001_2(t *testing.T) {\n"
	for _, s := range steps {
		body += "\tt.Run(\"" + s + "\", func(t *testing.T) {})\n"
	}
	return body + "}\n"
}

// gherkinRepo is a Go repo whose root requirement doc holds the Gherkin block
// body, proven by tests.
func gherkinRepo(t *testing.T, body, tests string) string {
	t.Helper()
	return acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
		".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
		rootDoc:           acceptance.DocWith("```gherkin\n" + body + "```\n"),
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

// gherkinFails fails unless requirements fails body naming want.
func gherkinFails(t *testing.T, body, want string) {
	t.Helper()
	out, code := judged(t, gherkinRepo(t, body, acceptance.TestFile("adder", "Test0001_1")))
	acceptance.Sees(t, out, code, 1, rootDoc+":"+want)
}

func Test0003_1(t *testing.T) {
	t.Parallel()
	out, code := judged(t, gherkinRepo(t, twoScenarios, provingTests("0001.2.1 add", "0001.2.2 double")))
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	// A step without its test case fails as it would for bullets.
	out, code = judged(t, gherkinRepo(t, twoScenarios, provingTests("0001.2.1 add")))
	acceptance.Sees(t, out, code, 1, "0001.2.2")
}

func Test0003_2(t *testing.T) {
	t.Parallel()
	gherkinFails(t, "SCENARIO: 0001.1\n  given two numbers\n  WHEN I add them\n  THEN I get their sum\n", "19: write the keyword GIVEN in caps")
	gherkinFails(t, "Scenario: 0001.1\n  GIVEN two numbers\n  WHEN I add them\n  THEN I get their sum\n", "18: write SCENARIO: in caps")
}

func Test0003_3(t *testing.T) {
	t.Parallel()
	gherkinFails(t, "SCENARIO: 0001.1\n  GIVEN two numbers\n  AND a calculator\n  WHEN I add them\n  THEN I get their sum\n", "20: indent AND one level deeper than the step it extends")
}

func Test0003_4(t *testing.T) {
	t.Parallel()
	gherkinFails(t, "SCENARIO: 0001.1\n  GIVEN two numbers, a calculator\n  WHEN I add them\n  THEN I get their sum\n", "19: split this step at its comma into an AND line")
	// A comma inside code belongs to the code.
	out, code := judged(t, gherkinRepo(t, "SCENARIO: 0001.1\n  GIVEN two numbers such as `1, 2`\n  WHEN I add them\n  THEN I get their sum\n", acceptance.TestFile("adder", "Test0001_1")))
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
}

func Test0003_5(t *testing.T) {
	t.Parallel()
	// A table whose rows hold different cell counts is one the parser rejects.
	gherkinFails(t, "SCENARIO: 0001.1\n  GIVEN two numbers\n    | a |\n    | b | c |\n  WHEN I add them\n  THEN I get their sum\n", "21: Gherkin: inconsistent cell count")
}

func Test0003_6(t *testing.T) {
	t.Parallel()
	t.Run("0003.6.1 up turns bullets into one Gherkin block", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
			".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
		}))
		out, code := judged(t, repo, "--fix")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		docHolds(t, repo, "```gherkin\nSCENARIO: 0001.1\n  GIVEN two numbers\n  WHEN I add them\n  THEN I get their sum\n```\n")
	})
	t.Run("0003.6.2 down turns the block back into bullets", func(t *testing.T) {
		t.Parallel()
		repo := gherkinRepo(t, twoScenarios, provingTests("0001.2.1 add", "0001.2.2 double"))
		out, code := judged(t, repo, "--down", "0003.6-gherkin")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		docHolds(t, repo, "- 0001.1\n  - Given two numbers and a calculator\n    - When I add them\n      - Then I get their sum\n"+
			"- 0001.2\n  - Given two numbers\n    1. When I add them\n       - Then I get their sum\n    2. When I double it\n       - Then it is even\n")
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

// movedRun commits a bullet doc whose GIVEN holds a comma. It stages the doc
// as the Gherkin block body and runs requirements on the staged change.
func movedRun(t *testing.T, body string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
		".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
		rootDoc:           acceptance.DocWith("- 0001.1\n  - Given two numbers, a calculator\n    - When I add them\n      - Then I get their sum\n"),
	}))
	acceptance.Write(t, repo, map[string]string{rootDoc: acceptance.DocWith("```gherkin\n" + body + "```\n")})
	acceptance.Git(t, repo, "add", "-A")
	return judged(t, repo)
}

func Test0003_7(t *testing.T) {
	t.Parallel()
	moved := "SCENARIO: 0001.1\n  GIVEN two numbers\n    AND a calculator\n  WHEN I add them\n  THEN I get their sum\n"
	out, code := movedRun(t, moved)
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	// A move that also changes what a step says still needs its test changed.
	out, code = movedRun(t, strings.Replace(moved, "their sum", "their product", 1))
	acceptance.Sees(t, out, code, 1, rootDoc)
}
