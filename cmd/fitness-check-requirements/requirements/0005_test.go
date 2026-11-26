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

func Test0005_4(t *testing.T) {
	t.Parallel()
	chained := acceptance.Reqs(acceptance.Acceptances("0001.1") + acceptance.Chain("0001.2", 2))
	named := func(cases ...string) map[string]string {
		return acceptance.With(chained, map[string]string{"chain_test.go": acceptance.CaseFile("adder", "Test0001_2", cases...)})
	}
	const sum, even = "0001.2.1.1.1 I get their sum", "0001.2.1.2.1 the total is even"
	t.Run("0005.4.1.1.1 it passes beside cases named otherwise", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named(sum, even, "edge"))
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0005.4.1.2.1 it fails naming the THEN", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named(sum))
		acceptance.Sees(t, out, code, 1, acceptance.PairedDoc+":26: THEN 0001.2.1.2.1 has no test case")
	})
	t.Run("0005.4.1.2.2 it fails naming the test to add the case to", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named(sum))
		acceptance.Sees(t, out, code, 1, `add a case to Test0001_2 named "`+even+`"`)
	})
	t.Run("0005.4.1.3.1 it fails naming the case and the name it needs", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named(sum, "0001.2.1.2.1 totals"))
		acceptance.Sees(t, out, code, 1, `THEN 0001.2.1.2.1 has a case named "0001.2.1.2.1 totals" in Test0001_2; rename it to "`+even+`"`)
	})
	t.Run("0005.4.1.4.1 it fails naming the count", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named(sum, even, even))
		acceptance.Sees(t, out, code, 1, "THEN 0001.2.1.2.1 has 2 test cases in Test0001_2")
	})
	t.Run("0005.4.1.5.1 it fails naming the case", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, named(sum, even, "0001.2.1.9.1 halves", "0001.1.1.1.1 I get their sum"))
		acceptance.Sees(t, out, code, 1, `case "0001.2.1.9.1 halves" names 0001.2.1.9.1, which no requirement doc holds`,
			`case "0001.1.1.1.1 I get their sum" names 0001.1.1.1.1, which 0001.1 holds, not 0001.2`)
	})
	t.Run("0005.4.1.6.1 it fails naming the case", func(t *testing.T) {
		t.Parallel()
		silent := "package adder\n\nimport \"testing\"\n\nfunc Test0001_1(t *testing.T) {\n\tt.Run(\"0001.1.1.1.1 I get their sum\", func(t *testing.T) {})\n}\n"
		out, code := checked(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": silent}))
		acceptance.Sees(t, out, code, 1, `case "0001.1.1.1.1 I get their sum" has no assertion call`)
		table := "package adder\n\nimport \"testing\"\n\nfunc Test0001_1(t *testing.T) {\n\tfor _, c := range []struct{ name string }{{\"0001.1.1.1.1 I get their sum\"}} {\n\t\tt.Run(c.name, func(t *testing.T) { t.Parallel(); _ = c })\n\t}\n}\n"
		out, code = checked(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": table}))
		acceptance.Sees(t, out, code, 1, `case "0001.1.1.1.1 I get their sum" has no assertion call`)
		helpers := strings.Replace(silent, "func(t *testing.T) {}", "func(t *testing.T) { if t.Name() == \"\" { t.Fatalf(\"x\") } }", 1)
		out, code = checked(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": helpers}))
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0005.4.1.7.1 it needs only its own test until the Node switch is on", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, scriptRepo())
		acceptance.Sees(t, out, code, 0)
	})
}
