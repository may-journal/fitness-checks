package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// thenText is the one THEN of acceptance 0001.1 as acceptance.Acceptances writes it.
const thenText = "I get their sum"

// thenID is the ID of that THEN.
const thenID = "0001.1.1.1.1"

// tagged is a Go test file whose Test0001_1 holds one case per tag, each
// tag an ID and a hash.
func tagged(tags ...[2]string) map[string]string {
	body := "package adder\n\nimport (\n\t\"testing\"\n\n\t\"github.com/may-journal/fitness-checks/req\"\n)\n\nfunc Test0001_1(t *testing.T) {\n"
	for _, tag := range tags {
		body += "\treq.Then(t, \"" + tag[0] + "\", \"" + tag[1] + "\", func(a req.A) { a.Equal(1, 1) })\n"
	}
	return map[string]string{"adder_test.go": body + "}\n"}
}

// currentTag is the tag of the THEN with its current hash.
func currentTag() [2]string { return [2]string{thenID, req.Hash(thenText)} }

func Test0010_1(t *testing.T) {
	t.Parallel()
	t.Run("0010.1.1 a tag naming no THEN or holding another hash fails", func(t *testing.T) {
		t.Parallel()
		gone, goneCode := checked(t, acceptance.With(acceptance.GoRepo(), tagged(currentTag(), [2]string{"0001.1.1.9.9", "00000000"})))
		stale, staleCode := checked(t, acceptance.With(acceptance.GoRepo(), tagged([2]string{thenID, "00000000"})))
		req.Then(t, "0010.1.2.1.1", "9da7bd33", func(_ req.A) {
			acceptance.Sees(t, gone, goneCode, 1, "adder_test.go:11: tag 0001.1.1.9.9#00000000 names no THEN")
			acceptance.Sees(t, stale, staleCode, 1, "adder_test.go:10: tag "+thenID+"#00000000 is stale")
		})
		req.Then(t, "0010.1.2.1.2", "c11c5ad1", func(_ req.A) {
			acceptance.Sees(t, stale, staleCode, 1, "then tag the case "+thenID+"#"+req.Hash(thenText))
		})
	})
	t.Run("0010.1.2 two tags on one THEN fail", func(t *testing.T) {
		t.Parallel()
		out, code := checked(t, acceptance.With(acceptance.GoRepo(), tagged(currentTag(), currentTag())))
		req.Then(t, "0010.1.2.2.1", "f0ff953d", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "THEN "+thenID+" is tagged twice")
		})
		req.Then(t, "0010.1.2.2.2", "3ef949ee", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "also at adder_test.go:10")
		})
	})
	t.Run("0010.1.3 current tags pass beside named cases", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0010.1.2.3.1", "1ff11095", func(_ req.A) {
			out, code := checked(t, acceptance.With(acceptance.GoRepo(), tagged(currentTag())))
			acceptance.Sees(t, out, code, 0)
			acceptance.PassedWith(t, out, "requirements", "2")
		})
	})
	t.Run("0010.1.4 a tag outside the test of its acceptance fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0010.1.2.4.1", "6b31dd4c", func(_ req.A) {
			moved := strings.Replace(tagged(currentTag())["adder_test.go"], "func Test0001_1(", "func TestHelper(", 1) + "\nfunc Test0001_1(t *testing.T) {}\n"
			out, code := checked(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": moved}))
			acceptance.Sees(t, out, code, 1, "sits in TestHelper, outside the test of its acceptance; move the case into Test0001_1")
		})
	})
}

// untagged is a repo whose Test0001_1 runs a case named for step 1 of a two-step chain.
func untagged() map[string]string {
	cases := acceptance.StepCases("0001.1", 2)
	return acceptance.With(acceptance.Reqs(acceptance.Chain("0001.1", 3)), map[string]string{
		"adder_test.go": strings.ReplaceAll(acceptance.CaseFile("adder", "Test0001_1", cases...), "func(t *testing.T) {})", "func(t *testing.T) {\n\t\tt.Parallel()\n\t\tacceptance.Sees(t, out, code, 0)\n\t})"),
	})
}

func Test0010_2(t *testing.T) {
	t.Parallel()
	t.Run("0010.2.1 the tag mode tags the cases it matches and lists the rest", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", untagged())
		out, code := acceptance.Fitness(t, repo, nil, "requirements", "--tag")
		req.Then(t, "0010.2.1.1.1", "14021cef", func(a req.A) {
			raw, err := os.ReadFile(filepath.Join(repo, "adder_test.go"))
			a.NoError(err)
			a.Contains(string(raw), `req.Then(t, "0001.1.1.1.1", "`+req.Hash(thenText))
		})
		req.Then(t, "0010.2.1.1.2", "9f8e32d9", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "THEN 0001.1.1.3.1 has no tagged test case")
		})
	})
	t.Run("0010.2.2 a dry run writes no file", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0010.2.1.2.1", "fbbb4ffc", func(a req.A) {
			repo := acceptance.Example(t, "happyRepo", untagged())
			before, err := os.ReadFile(filepath.Join(repo, "adder_test.go"))
			a.NoError(err)
			acceptance.Fitness(t, repo, nil, "requirements", "--tag", "--dry-run")
			after, err := os.ReadFile(filepath.Join(repo, "adder_test.go"))
			a.NoError(err)
			a.Equal(string(before), string(after))
		})
	})
}
