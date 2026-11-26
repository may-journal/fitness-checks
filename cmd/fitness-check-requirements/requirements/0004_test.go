package requirements

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// suite runs the full suite on happyRepo plus files and returns the output.
func suite(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil)
}

// cappedRepo is pairedRepo whose 0001 holds count tested acceptances, under
// a positive limit as requirements.maxAcceptances or else the default.
func cappedRepo(limit, count int) map[string]string {
	var ids, tests []string
	for n := 1; n <= count; n++ {
		ids = append(ids, fmt.Sprintf("0001.%d", n))
		tests = append(tests, fmt.Sprintf("Test0001_%d", n))
	}
	rc := `{"policy": "external", "checks": ["requirements"]}`
	if limit > 0 {
		rc = fmt.Sprintf(`{"policy": "external", "checks": ["requirements"], "requirements": {"maxAcceptances": %d}}`, limit)
	}
	return acceptance.With(acceptance.PairedRepo(), map[string]string{
		".fitnessrc.json":    rc + "\n",
		"reqs/0001_test.go":  acceptance.TestFile("reqs", tests...),
		acceptance.PairedDoc: acceptance.RequirementDoc(ids...),
	})
}

func Test0004_1(t *testing.T) {
	t.Parallel()
	t.Run("0004.1.1 a file name and title ID that differ fail", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.DocOf("docs/requirements/0002-adds-numbers.md", acceptance.RequirementDoc("0001.1")))
		acceptance.Sees(t, out, code, 1, "name the file NNNN-kebab-title.md and title it `# NNNN Title` with the same ID")
	})
	t.Run("0004.1.2 two docs sharing an ID fail", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"docs/requirements/0001-other.md": acceptance.RequirementDoc("0001.1")}))
		acceptance.Sees(t, out, code, 1, "requirement 0001 is used by")
	})
	t.Run("0004.1.3 a repeated acceptance ID fails", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.Reqs(acceptance.Acceptances("0001.1", "0001.1")))
		acceptance.Sees(t, out, code, 1, "0001.1 appears twice")
	})
	t.Run("0004.1.4 a doc past the acceptance cap fails", func(t *testing.T) {
		t.Parallel()
		for _, row := range []struct {
			max, count int
			want       []string
		}{
			{0, 8, nil},
			{0, 9, []string{acceptance.PairedDoc + ":49: 0001 holds 9 acceptances, at most 8; fold cases into table rows of one acceptance, or move a separate goal to its own doc"}},
			{10, 9, nil},
			{10, 11, []string{acceptance.PairedDoc + ":57: 0001 holds 11 acceptances, at most 10;"}},
			{3, 4, []string{acceptance.PairedDoc + ":29: 0001 holds 4 acceptances, at most 3;"}},
		} {
			t.Run(fmt.Sprintf("cap %d holds %d", row.max, row.count), func(t *testing.T) {
				t.Parallel()
				out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", cappedRepo(row.max, row.count)), nil, "requirements")
				acceptance.Sees(t, out, code, min(len(row.want), 1), row.want...)
			})
		}
	})
}

func Test0004_2(t *testing.T) {
	t.Parallel()
	t.Run("0004.2.1 sections that are missing or extra or empty fail", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.DocOf("docs/requirements/0001-adds-numbers.md", acceptance.RequirementDoc("0001.1")+"\n## Out of scope\n\n- Subtraction.\n"))
		acceptance.Sees(t, out, code, 1, "sections must be exactly ## Why, ## Measurement, ## Requirements")
	})
	t.Run("0004.2.2 a Measurement that is not a numerator and denominator fails", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.DocOf("docs/requirements/0001-adds-numbers.md", strings.Replace(acceptance.RequirementDoc("0001.1"), "\n-\n", "\nover\n", 1)))
		acceptance.Sees(t, out, code, 1, "Measurement needs three lines: numerator, `-`, denominator")
	})
	t.Run("0004.2.3 a Source line fails and the fix deletes it", func(t *testing.T) {
		t.Parallel()
		doc := acceptance.DocWith("1. Adds\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n")
		sourced := acceptance.WithSource(doc, "`gone.go:99`")
		repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.PairedRepo(), map[string]string{acceptance.PairedDoc: sourced}))
		out, code := acceptance.Fitness(t, repo, nil, "requirements")
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+":15: delete this `Source:` line, or run `fitness requirements --fix`")
		out, code = acceptance.Fitness(t, repo, nil, "requirements", "--fix")
		acceptance.Sees(t, out, code, 0)
		if got, err := os.ReadFile(filepath.Join(repo, acceptance.PairedDoc)); err != nil || string(got) != doc {
			t.Errorf("--fix must leave the doc without its Source line, got %v:\n%s", err, got)
		}
		acceptance.Git(t, repo, "add", "-A")
		out, code = acceptance.Fitness(t, repo, nil, "requirements")
		acceptance.Sees(t, out, code, 0)
		if strings.Contains(out, "Source") {
			t.Errorf("a doc without a Source line gets no error, got:\n%s", out)
		}
	})
}

func Test0004_3(t *testing.T) {
	t.Parallel()
	out, code := suite(t, acceptance.Reqs(acceptance.Acceptances("0001.1")+"        - When I add them again\n"))
	acceptance.Sees(t, out, code, 1, "needs exactly one Given, one When, and one Then")
	out, code = suite(t, acceptance.Reqs("- 0001.1\n    - When I add them\n        - Given two numbers\n            - Then I get their sum\n"))
	acceptance.Sees(t, out, code, 1, "each acceptance is one Given, then one When, then one Then, each nested under the line above")
}

func Test0004_4(t *testing.T) {
	t.Parallel()
	t.Run("0004.4.1 steps chain up to the limit beside a single step", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1"), "chain_test.go": acceptance.CaseFile("adder", "Test0001_2", acceptance.StepCases("0001.2", 7)...)}
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", acceptance.With(acceptance.Reqs(acceptance.Acceptances("0001.1")+acceptance.Chain("0001.2", 7)), tests)), nil, "requirements")
		acceptance.Sees(t, out, code, 0)
		out, code = suite(t, acceptance.Reqs(acceptance.Chain("0001.1", 8)))
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+":17: 0001.1 chains 8 steps, at most 7; move the rest to a new ID")
	})
	t.Run("0004.4.2 misnumbered or unnumbered or single steps fail", func(t *testing.T) {
		t.Parallel()
		for body, want := range map[string]string{
			strings.Replace(acceptance.Chain("0001.1", 3), "    3. When", "    4. When", 1): ":23: number this step `3.`; a chain numbers its steps `1.` to `7.` in order",
			strings.Replace(acceptance.Chain("0001.1", 2), "    2. When", "    - When", 1):  ":21: number this step `2.`",
			acceptance.Chain("0001.1", 1): ":19: a single step is the plain bullet `- When`; number steps only in a chain of two or more",
			strings.Replace(acceptance.Chain("0001.1", 2), "       - Then I get their sum\n", "", 1): ":19: step 1 needs exactly one Then nested under its When",
			strings.Replace(acceptance.Chain("0001.1", 2), "    2. When", "    2. Then", 1):          ":21: " + "each step is a numbered When nested under the Given",
		} {
			out, code := suite(t, acceptance.Reqs(body))
			acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+want)
		}
	})
	t.Run("0004.4.3 two docs sharing a WHEN fail", func(t *testing.T) {
		t.Parallel()
		other := "docs/requirements/0002-adds-again.md"
		adds := "- 0001.1\n    - Given Writer has two whole numbers\n        - When I add them\n            - Then I get their sum\n"
		again := "- 0002.1\n    - Given Writer has two whole numbers that carry\n        - When I add them\n            - Then I get their sum\n"
		out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{
			"adder_test.go":      acceptance.TestFile("adder", "Test0001_1", "Test0002_1"),
			acceptance.PairedDoc: acceptance.DocWith(adds),
			other:                strings.Replace(acceptance.DocWith(again), "0001 Adds Numbers", "0002 Adds Again", 1),
		}))
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+" and "+other+` both cover "When I add them"`, "blob/main/fitness/adr-0004")
	})
}

func Test0004_5(t *testing.T) {
	t.Parallel()
	t.Run("0004.5.1 a deleted ID brought back fails", func(t *testing.T) {
		t.Parallel()
		two := acceptance.With(acceptance.GoRepo(), map[string]string{
			"adder_test.go":                          acceptance.TestFile("adder", "Test0001_1", "Test0001_2"),
			"docs/requirements/0001-adds-numbers.md": acceptance.RequirementDoc("0001.1", "0001.2"),
		})
		repo := acceptance.Example(t, "happyRepo", two)
		acceptance.Commit(t, repo, acceptance.GoRepo(), "docs(requirements): drop 0001.2")
		acceptance.Write(t, repo, two)
		out, code := acceptance.Fitness(t, repo, nil)
		acceptance.Sees(t, out, code, 1, "0001.2 was deleted before")
	})
	t.Run("0004.5.2 a deleted step brought back fails", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", acceptance.GoRepo())
		cased := func(steps int) map[string]string {
			return acceptance.With(acceptance.Reqs(acceptance.Chain("0001.1", steps)), map[string]string{"adder_test.go": acceptance.CaseFile("adder", "Test0001_1", acceptance.StepCases("0001.1", steps)...)})
		}
		acceptance.Commit(t, repo, cased(3), "docs(requirements): chain 0001.1 in three steps")
		out, code := acceptance.Fitness(t, repo, nil, "requirements")
		acceptance.Sees(t, out, code, 0)
		acceptance.Commit(t, repo, cased(2), "docs(requirements): drop step 0001.1.3")
		acceptance.Write(t, repo, cased(3))
		out, code = acceptance.Fitness(t, repo, nil, "requirements")
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+":23: 0001.1.3 was deleted before; IDs are never reused, so chain the step under a new acceptance ID")
	})
	t.Run("0004.5.3 a commit removing every doc restarts the history", func(t *testing.T) {
		t.Parallel()
		two := acceptance.With(acceptance.GoRepo(), map[string]string{
			"adder_test.go":                          acceptance.TestFile("adder", "Test0001_1", "Test0001_2"),
			"docs/requirements/0001-adds-numbers.md": acceptance.RequirementDoc("0001.1", "0001.2"),
		})
		repo := acceptance.Example(t, "happyRepo", two)
		acceptance.Commit(t, repo, acceptance.GoRepo(), "docs(requirements): drop 0001.2")
		restructured := acceptance.With(two, map[string]string{
			"docs/requirements/0001-adds-numbers.md": "",
			"docs/requirements/0001-sums-numbers.md": acceptance.RequirementDoc("0001.1", "0001.2"),
		})
		acceptance.Write(t, repo, restructured)
		acceptance.Git(t, repo, "add", "-A")
		out, code := acceptance.Fitness(t, repo, nil, "requirements")
		acceptance.Sees(t, out, code, 0)
		acceptance.Commit(t, repo, restructured, "docs(requirements): restart the numbers")
		out, code = acceptance.Fitness(t, repo, nil, "requirements")
		acceptance.Sees(t, out, code, 0)
		dropped := map[string]string{"adder_test.go": acceptance.GoRepo()["adder_test.go"], "docs/requirements/0001-sums-numbers.md": acceptance.RequirementDoc("0001.1")}
		acceptance.Commit(t, repo, dropped, "docs(requirements): drop 0001.2 again")
		acceptance.Write(t, repo, restructured)
		out, code = acceptance.Fitness(t, repo, nil, "requirements")
		acceptance.Sees(t, out, code, 1, "0001.2 was deleted before")
	})
}
