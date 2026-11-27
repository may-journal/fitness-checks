// Package requirements proves the go-test check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a failing test is named with its message", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, goTestModule("", goTestBroken))
		req.Then(t, "0001.1.2.1.1", "cadf814f", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "--- FAIL: TestBroken")
		})
		req.Then(t, "0001.1.2.1.2", "06a86cc7", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "m_test.go:5: boom")
		})
	})
	t.Run("0001.1.2 a failure in a subfolder module shows its folder", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "56647b75", func(_ req.A) {
			out, code := goTest(t, goTestModule("svc/", goTestBroken))
			acceptance.Sees(t, out, code, 1, "svc: --- FAIL: TestBroken", "svc: m_test.go:5: boom")
		})
	})
	t.Run("0001.1.3 a test that does not compile fails with the build error then its location", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, goTestModule("", "func TestBroken(t *testing.T) { undefinedThing() }\n"))
		req.Then(t, "0001.1.2.3.1", "95206fff", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "undefined: undefinedThing")
		})
		req.Then(t, "0001.1.2.3.2", "8508af08", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "m_test.go:5:33:")
		})
	})
	t.Run("0001.1.4 passing tests pass then count the modules", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, goTestModule("", "func TestFine(t *testing.T) {}\n"))
		req.Then(t, "0001.1.2.4.1", "07b8e753", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.4.2", "a03182df", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "1 files scanned")
			acceptance.PassedWith(t, out, "go-test", "1")
		})
	})
	t.Run("0001.1.5 a repo with nothing to test passes without running go test", func(t *testing.T) {
		t.Parallel()
		bare, bareCode := goTest(t, nil)
		// Without Go on PATH a run of go test would show info, not a pass.
		noGo, noGoCode := goTest(t, nil, acceptance.PathWithout(t, "go"))
		// Fixtures under testdata and vendor are not modules to test.
		files := acceptance.With(goTestModule("testdata/a/", goTestBroken), goTestModule("vendor/b/", goTestBroken))
		fixtures, fixturesCode := goTest(t, files)
		// A change touching no Go file leaves a failing module alone.
		repo := acceptance.Example(t, "happyRepo", goTestModule("", goTestBroken))
		acceptance.Write(t, repo, map[string]string{"README.md": "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n\nOne more line.\n"})
		acceptance.Git(t, repo, "add", "README.md")
		changed, changedCode := acceptance.Fitness(t, repo, nil, "go-test")
		req.Then(t, "0001.1.2.5.1", "857dd667", func(a req.A) {
			acceptance.Sees(t, bare, bareCode, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, bare, "go-test")
		})
		req.Then(t, "0001.1.2.5.2", "56cf814f", func(a req.A) {
			acceptance.Sees(t, noGo, noGoCode, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, noGo, "go-test")
		})
		req.Then(t, "0001.1.2.5.3", "54749f6e", func(a req.A) {
			acceptance.Sees(t, fixtures, fixturesCode, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, fixtures, "go-test")
		})
		req.Then(t, "0001.1.2.5.4", "3b8d9bc6", func(a req.A) {
			acceptance.Sees(t, changed, changedCode, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, changed, "go-test")
		})
	})
	t.Run("0001.1.6 a missing Go shows the install hint", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "188437da", func(_ req.A) {
			out, code := goTest(t, goTestModule("", "func TestFine(t *testing.T) {}\n"), acceptance.PathWithout(t, "go"))
			acceptance.ShowsInfo(t, out, code, "go-test", "Go not installed: see https://go.dev/dl")
		})
	})
	t.Run("0001.1.7 a run that fails without a finding shows why", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.7.1", "4ca0535d", func(a req.A) {
			out, code := goTest(t, goTestModule("", goTestBroken), acceptance.CrashingGo(t, "test", "go: unexpected internal error"))
			acceptance.ShowsInfo(t, out, code, "go-test", "go test ./... exited 2 without a finding: go: unexpected internal error")
		})
		req.Then(t, "0001.1.2.7.2", "0b9ea8db", func(a req.A) {
			// A go whose interpreter does not exist cannot start.
			dir := t.TempDir()
			acceptance.Write(t, dir, map[string]string{"go": "#!/nonexistent/interpreter\n"})
			out, code := goTest(t, goTestModule("", goTestBroken), "PATH="+dir+string(os.PathListSeparator)+os.Getenv("PATH"))
			acceptance.ShowsInfo(t, out, code, "go-test")
		})
		req.Then(t, "0001.1.2.7.3", "68ea2962", func(a req.A) {
			out, code := goTest(t, goTestModule("", goTestBroken), acceptance.KilledGo(t, "test"))
			acceptance.ShowsInfo(t, out, code, "go-test", "go test")
		})
	})
}

// goTest runs `fitness-install -- go-test` on happyRepo with files written
// over it and env added.
func goTest(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), env, "go-test")
}

// goTestModule is a Go module in dir whose test file holds body.
func goTestModule(dir, body string) map[string]string {
	return map[string]string{
		dir + "go.mod":    "module example.com/m\n\ngo 1.22\n",
		dir + "m_test.go": "package m\n\nimport \"testing\"\n\n" + body,
	}
}

// goTestBroken is a test that fails with the message "boom" on line 5.
const goTestBroken = "func TestBroken(t *testing.T) { t.Error(\"boom\") }\n"
