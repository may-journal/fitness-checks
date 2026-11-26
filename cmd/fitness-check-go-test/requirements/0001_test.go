// Package requirements proves the go-test check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	failing := func(t *testing.T) (string, int) { return goTest(t, goTestModule("", goTestBroken)) }
	t.Run("0001.1.2.1.1 it fails naming the test", func(t *testing.T) {
		t.Parallel()
		out, code := failing(t)
		acceptance.Sees(t, out, code, 1, "--- FAIL: TestBroken")
	})
	t.Run("0001.1.2.1.2 it fails naming its `file.go:line` message", func(t *testing.T) {
		t.Parallel()
		out, code := failing(t)
		acceptance.Sees(t, out, code, 1, "m_test.go:5: boom")
	})
	t.Run("0001.1.2.2.1 each line is prefixed with that module's folder", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, goTestModule("svc/", goTestBroken))
		acceptance.Sees(t, out, code, 1, "svc: --- FAIL: TestBroken", "svc: m_test.go:5: boom")
	})
	noCompile := func(t *testing.T) (string, int) {
		return goTest(t, goTestModule("", "func TestBroken(t *testing.T) { undefinedThing() }\n"))
	}
	t.Run("0001.1.2.3.1 it fails with the build error", func(t *testing.T) {
		t.Parallel()
		out, code := noCompile(t)
		acceptance.Sees(t, out, code, 1, "undefined: undefinedThing")
	})
	t.Run("0001.1.2.3.2 it fails with its `file.go:line:col`", func(t *testing.T) {
		t.Parallel()
		out, code := noCompile(t)
		acceptance.Sees(t, out, code, 1, "m_test.go:5:33:")
	})
	passing := func(t *testing.T) (string, int) {
		return goTest(t, goTestModule("", "func TestFine(t *testing.T) {}\n"))
	}
	t.Run("0001.1.2.4.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := passing(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
	})
	t.Run("0001.1.2.4.2 it counts the modules tested", func(t *testing.T) {
		t.Parallel()
		out, _ := passing(t)
		acceptance.PassedWith(t, out, "go-test", "1")
	})
	t.Run("0001.1.2.5.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-test")
	})
	t.Run("0001.1.2.5.2 it never runs `go test`", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, nil, acceptance.CrashingGo(t, "test", "go test must not run"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-test")
	})
	t.Run("0001.1.2.5.3 a `testdata` or `vendor` module counts as none", func(t *testing.T) {
		t.Parallel()
		files := acceptance.With(goTestModule("testdata/a/", goTestBroken), goTestModule("vendor/b/", goTestBroken))
		out, code := goTest(t, files)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-test")
	})
	t.Run("0001.1.2.5.4 a change touching no Go file counts as none", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", goTestModule("", goTestBroken))
		acceptance.Write(t, repo, map[string]string{"README.md": "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n\nOne more line.\n"})
		acceptance.Git(t, repo, "add", "README.md")
		out, code := acceptance.Fitness(t, repo, nil, "go-test")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-test")
	})
	t.Run("0001.1.2.6.1 it reports info with the hint `see https://go.dev/dl`", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, goTestModule("", "func TestFine(t *testing.T) {}\n"), acceptance.PathWithout(t, "go"))
		acceptance.ShowsInfo(t, out, code, "go-test", "Go not installed: see https://go.dev/dl")
	})
	t.Run("0001.1.2.7.1 it reports info quoting what went wrong", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, goTestModule("", goTestBroken), acceptance.CrashingGo(t, "test", "go: unexpected internal error"))
		acceptance.ShowsInfo(t, out, code, "go-test", "go test ./... exited 2 without a finding: go: unexpected internal error")
	})
	t.Run("0001.1.2.7.2 a failure to start counts the same", func(t *testing.T) {
		t.Parallel()
		unstartable := t.TempDir()
		acceptance.Write(t, unstartable, map[string]string{"go": "#!/nonexistent/interpreter\n"})
		out, code := goTest(t, goTestModule("", goTestBroken), "PATH="+unstartable+string(os.PathListSeparator)+os.Getenv("PATH"))
		acceptance.ShowsInfo(t, out, code, "go-test", "go test")
	})
	t.Run("0001.1.2.7.3 a kill counts the same", func(t *testing.T) {
		t.Parallel()
		out, code := goTest(t, goTestModule("", goTestBroken), acceptance.KilledGo(t, "test"))
		acceptance.ShowsInfo(t, out, code, "go-test", "go test")
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
