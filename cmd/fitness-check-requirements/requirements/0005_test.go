package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// checked runs only the requirements check on happyRepo plus files.
func checked(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "requirements")
}

func Test0005_1(t *testing.T) {
	t.Parallel()
	t.Run("0005.1.1 an acceptance with no test fails naming the test to write", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": ""}))
		acceptance.Sees(t, out, code, 1, "0001.1 has no test; name exactly one test Test0001_1")
	})
	t.Run("0005.1.2 an acceptance with two tests fails naming the count", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"more/more_test.go": acceptance.TestFile("more", "Test0001_1")}))
		acceptance.Sees(t, out, code, 1, "0001.1 is owned by 2 tests")
	})
	t.Run("0005.1.3 a test named for an undefined acceptance fails naming the test then the ID", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "Test0001_2")}))
		acceptance.Sees(t, out, code, 1, "Test0001_2 names 0001.2, which no requirement defines")
	})
}

func Test0005_2(t *testing.T) {
	t.Parallel()
	t.Run("0005.2.1 a test named for no acceptance fails", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "TestAddsAgain")}))
		acceptance.Sees(t, out, code, 1, "TestAddsAgain proves no requirement")
		out, code = suite(t, acceptance.With(acceptance.Reqs(acceptance.Chain("0001.1", 2)), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "TestAdds0001_1_2")}))
		acceptance.Sees(t, out, code, 1, "TestAdds0001_1_2 proves no requirement")
	})
	t.Run("0005.2.2 a test named for a step fails telling which test to name", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.With(acceptance.Reqs(acceptance.Chain("0001.1", 2)), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "Test0001_1_2")}))
		acceptance.Sees(t, out, code, 1, "Test0001_1_2 names step 0001.1.2; one test proves every step of its acceptance, so name it Test0001_1")
	})
}

func Test0005_3(t *testing.T) {
	t.Parallel()
	chained := acceptance.Reqs(acceptance.Acceptances("0001.1") + acceptance.Chain("0001.2", 2))
	named := func(cases ...string) map[string]string {
		return acceptance.With(chained, map[string]string{"chain_test.go": acceptance.CaseFile("adder", "Test0001_2", cases...)})
	}
	t.Run("0005.3.1 a case named for each step passes beside other cases", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named("0001.2.1 adds", "0001.2.2 doubles", "edge"))
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0005.3.2 a step with no case fails naming the step then the test", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named("0001.2.1 adds"))
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+":25: step 0001.2.2 has no test case; add a case to Test0001_2 whose name opens with 0001.2.2")
	})
	t.Run("0005.3.3 a case naming a step its acceptance does not chain fails", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named("0001.2.1 adds", "0001.2.2 doubles", "0001.2.3 halves", "0001.1.1 sums"))
		acceptance.Sees(t, out, code, 1, `chain_test.go:8: case "0001.2.3 halves" names step 0001.2.3, which 0001.2 does not chain`,
			`chain_test.go:9: case "0001.1.1 sums" names step 0001.1.1, which 0001.2 does not chain`)
		out, code = checked(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": acceptance.CaseFile("adder", "Test0001_1", "0001.1.1 adds")}))
		acceptance.Sees(t, out, code, 1, `case "0001.1.1 adds" names step 0001.1.1, which 0001.1 does not chain`)
	})
	t.Run("0005.3.4 a malformed chain never fails its good cases", func(t *testing.T) {
		t.Parallel()
		out, _ := checked(t, acceptance.With(acceptance.Reqs(acceptance.Chain("0001.1", 8)), map[string]string{"adder_test.go": acceptance.CaseFile("adder", "Test0001_1", acceptance.StepCases("0001.1", 7)...)}))
		if strings.Contains(out, "does not chain") {
			t.Errorf("a malformed chain fails its good cases:\n%s", out)
		}
	})
}
