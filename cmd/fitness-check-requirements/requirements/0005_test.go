package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
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
		req.Then(t, "0005.1.1.1.1", "4d5ed9db", func(_ req.A) {
			out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": ""}))
			acceptance.Sees(t, out, code, 1, "0001.1 has no test; name exactly one test Test0001_1")
		})
	})
	t.Run("0005.1.2 an acceptance with two tests fails naming the count", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0005.1.1.2.1", "a1be1a2a", func(_ req.A) {
			out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"more/more_test.go": acceptance.TestFile("more", "Test0001_1")}))
			acceptance.Sees(t, out, code, 1, "0001.1 is owned by 2 tests")
		})
	})
	t.Run("0005.1.3 a test named for an undefined acceptance fails naming the test then the ID", func(t *testing.T) {
		t.Parallel()
		out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "Test0001_2")}))
		req.Then(t, "0005.1.1.3.1", "cadf814f", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "Test0001_2 names")
		})
		req.Then(t, "0005.1.1.3.2", "aec3d74f", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "names 0001.2, which no requirement defines")
		})
	})
}

func Test0005_2(t *testing.T) {
	t.Parallel()
	t.Run("0005.2.1 a test named for no acceptance fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0005.2.1.1.1", "fb7c440e", func(_ req.A) {
			out, code := suite(t, acceptance.With(acceptance.GoRepo(), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "TestAddsAgain")}))
			acceptance.Sees(t, out, code, 1, "TestAddsAgain proves no requirement")
			out, code = suite(t, acceptance.With(acceptance.Reqs(acceptance.Chain("0001.1", 2)), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "TestAdds0001_1_2")}))
			acceptance.Sees(t, out, code, 1, "TestAdds0001_1_2 proves no requirement")
		})
	})
	t.Run("0005.2.2 a test named for a step fails telling which test to name", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0005.2.1.2.1", "9f6a2aad", func(_ req.A) {
			out, code := suite(t, acceptance.With(acceptance.Reqs(acceptance.Chain("0001.1", 2)), map[string]string{"adder_test.go": acceptance.TestFile("adder", "Test0001_1", "Test0001_1_2")}))
			acceptance.Sees(t, out, code, 1, "Test0001_1_2 names step 0001.1.2; one test proves every step of its acceptance, so name it Test0001_1")
		})
	})
}
