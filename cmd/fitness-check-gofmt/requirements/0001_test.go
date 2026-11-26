// Package requirements proves the gofmt check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 it fails as `path: not gofmt-formatted (run: gofmt -w path)`", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "bad.go": gofmtMessy})
		acceptance.Sees(t, out, code, 1, "bad.go: not gofmt-formatted (run: gofmt -w bad.go)")
	})
	formatted := func(t *testing.T) (string, int) {
		return gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy})
	}
	t.Run("0001.1.2.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := formatted(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.2.2 it counts the files checked", func(t *testing.T) {
		t.Parallel()
		out, _ := formatted(t)
		acceptance.PassedWith(t, out, "gofmt", "1")
	})
	t.Run("0001.1.2.3.1 a file in `testdata` is not judged", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "x/testdata/fixture.go": gofmtMessy})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.3.2 a file in `vendor` is not judged", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "vendor/lib/lib.go": gofmtMessy})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.3.3 an untracked file is not judged", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{"ok.go": gofmtTidy})
		acceptance.Write(t, repo, map[string]string{"scratch.go": gofmtMessy})
		out, code := acceptance.Fitness(t, repo, nil, "gofmt")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.4.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "gofmt")
	})
	broken := func(t *testing.T) (string, int) {
		return gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "broken.go": "package m\n\nfunc {\n"})
	}
	t.Run("0001.1.2.5.1 it fails as `path:line:column: message`", func(t *testing.T) {
		t.Parallel()
		out, code := broken(t)
		acceptance.Sees(t, out, code, 1, "broken.go:3:6: expected 'IDENT', found '{'")
	})
	t.Run("0001.1.2.5.2 it never reads as an unformatted file", func(t *testing.T) {
		t.Parallel()
		out, code := broken(t)
		acceptance.Sees(t, out, code, 1)
		if strings.Contains(out, "not gofmt-formatted") {
			t.Errorf("a syntax error must fail as itself, not as an unformatted file:\n%s", out)
		}
	})
	t.Run("0001.1.2.6.1 it reports info with the hint `see https://go.dev/dl`", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, []string{acceptance.PathWithout(t, "go", "gofmt")}, map[string]string{"ok.go": gofmtTidy})
		acceptance.ShowsInfo(t, out, code, "gofmt", "Go not installed: see https://go.dev/dl")
	})
	crash := func(t *testing.T) (string, int) {
		crashed := acceptance.BrokenTool(t, "gofmt", "echo 'fatal error: runtime: out of memory' >&2\nexit 2\n")
		return gofmtCheck(t, []string{crashed}, map[string]string{"ok.go": gofmtTidy})
	}
	t.Run("0001.1.2.7.1 it reports info quoting what went wrong", func(t *testing.T) {
		t.Parallel()
		out, code := crash(t)
		acceptance.ShowsInfo(t, out, code, "gofmt", "gofmt -l exited 2 without a finding: fatal error: runtime: out of memory")
	})
	t.Run("0001.1.2.7.2 a crash never reads as an unformatted file", func(t *testing.T) {
		t.Parallel()
		out, code := crash(t)
		acceptance.ShowsInfo(t, out, code, "gofmt", "gofmt -l exited 2")
		if strings.Contains(out, "not gofmt-formatted") {
			t.Errorf("a crash must not read as an unformatted file:\n%s", out)
		}
	})
}

// gofmtCheck runs `fitness-install -- gofmt` on happyRepo with files written
// over it and env added.
func gofmtCheck(t *testing.T, env []string, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), env, "gofmt")
}

// gofmtTidy is a Go file gofmt leaves as it is; gofmtMessy is one it rewrites.
const (
	gofmtTidy  = "package m\n\nfunc F() int { return 1 }\n"
	gofmtMessy = "package m\nfunc   G( ) int {return 2}\n"
)
