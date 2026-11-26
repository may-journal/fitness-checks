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
	t.Run("0001.1.1 an unformatted file fails naming the fix", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "bad.go": gofmtMessy})
		acceptance.Sees(t, out, code, 1, "bad.go: not gofmt-formatted (run: gofmt -w bad.go)")
	})
	t.Run("0001.1.2 formatted files pass then count the files checked", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "gofmt", "1")
	})
	t.Run("0001.1.3 testdata or vendor or untracked files are not judged", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{
			"ok.go":                 gofmtTidy,
			"x/testdata/fixture.go": gofmtMessy,
			"vendor/lib/lib.go":     gofmtMessy,
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		repo := acceptance.Example(t, "happyRepo", map[string]string{"ok.go": gofmtTidy})
		acceptance.Write(t, repo, map[string]string{"scratch.go": gofmtMessy})
		out, code = acceptance.Fitness(t, repo, nil, "gofmt")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 a repo with no Go files passes", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "gofmt")
	})
	t.Run("0001.1.5 a file that does not parse fails as itself", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "broken.go": "package m\n\nfunc {\n"})
		acceptance.Sees(t, out, code, 1, "broken.go:3:6: expected 'IDENT', found '{'")
		if strings.Contains(out, "not gofmt-formatted") {
			t.Errorf("a syntax error must fail as itself, not as an unformatted file:\n%s", out)
		}
	})
	t.Run("0001.1.6 a missing Go shows the install hint", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, []string{acceptance.PathWithout(t, "go", "gofmt")}, map[string]string{"ok.go": gofmtTidy})
		acceptance.ShowsInfo(t, out, code, "gofmt", "Go not installed: see https://go.dev/dl")
	})
	t.Run("0001.1.7 a crash without a finding shows why", func(t *testing.T) {
		t.Parallel()
		crashed := acceptance.BrokenTool(t, "gofmt", "echo 'fatal error: runtime: out of memory' >&2\nexit 2\n")
		out, code := gofmtCheck(t, []string{crashed}, map[string]string{"ok.go": gofmtTidy})
		acceptance.ShowsInfo(t, out, code, "gofmt", "gofmt -l exited 2 without a finding: fatal error: runtime: out of memory")
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
