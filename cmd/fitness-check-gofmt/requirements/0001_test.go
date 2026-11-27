// Package requirements proves the gofmt check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 an unformatted file fails naming the fix", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.1.1", "e31beeb1", func(_ req.A) {
			out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "bad.go": gofmtMessy})
			acceptance.Sees(t, out, code, 1, "bad.go: not gofmt-formatted (run: gofmt -w bad.go)")
		})
	})
	t.Run("0001.1.2 formatted files pass then count the files checked", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy})
		req.Then(t, "0001.1.2.2.1", "07b8e753", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.2.2", "3a68e3b4", func(a req.A) {
			acceptance.PassedWith(t, out, "gofmt", "1")
		})
	})
	t.Run("0001.1.3 testdata or vendor or untracked files are not judged", func(t *testing.T) {
		t.Parallel()
		testdata, testdataCode := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "x/testdata/fixture.go": gofmtMessy})
		vendor, vendorCode := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "vendor/lib/lib.go": gofmtMessy})
		repo := acceptance.Example(t, "happyRepo", map[string]string{"ok.go": gofmtTidy})
		acceptance.Write(t, repo, map[string]string{"scratch.go": gofmtMessy})
		untracked, untrackedCode := acceptance.Fitness(t, repo, nil, "gofmt")
		req.Then(t, "0001.1.2.3.1", "c92dad27", func(a req.A) {
			acceptance.Sees(t, testdata, testdataCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.3.2", "347111ca", func(a req.A) {
			acceptance.Sees(t, vendor, vendorCode, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.3.3", "5895109f", func(a req.A) {
			acceptance.Sees(t, untracked, untrackedCode, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.4 a repo with no Go files passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.4.1", "857dd667", func(_ req.A) {
			out, code := gofmtCheck(t, nil, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "gofmt")
		})
	})
	t.Run("0001.1.5 a file that does not parse fails as itself", func(t *testing.T) {
		t.Parallel()
		out, code := gofmtCheck(t, nil, map[string]string{"ok.go": gofmtTidy, "broken.go": "package m\n\nfunc {\n"})
		req.Then(t, "0001.1.2.5.1", "cf6d3551", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "broken.go:3:6: expected 'IDENT', found '{'")
		})
		req.Then(t, "0001.1.2.5.2", "16a12f8c", func(a req.A) {
			a.True(!strings.Contains(out, "not gofmt-formatted"), "a syntax error must fail as itself, not as an unformatted file:\n%s", out)
		})
	})
	t.Run("0001.1.6 a missing Go shows the install hint", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "188437da", func(_ req.A) {
			out, code := gofmtCheck(t, []string{acceptance.PathWithout(t, "go", "gofmt")}, map[string]string{"ok.go": gofmtTidy})
			acceptance.ShowsInfo(t, out, code, "gofmt", "Go not installed: see https://go.dev/dl")
		})
	})
	t.Run("0001.1.7 a crash without a finding shows why", func(t *testing.T) {
		t.Parallel()
		crashed := acceptance.BrokenTool(t, "gofmt", "echo 'fatal error: runtime: out of memory' >&2\nexit 2\n")
		out, code := gofmtCheck(t, []string{crashed}, map[string]string{"ok.go": gofmtTidy})
		req.Then(t, "0001.1.2.7.1", "4ca0535d", func(a req.A) {
			acceptance.ShowsInfo(t, out, code, "gofmt", "gofmt -l exited 2 without a finding: fatal error: runtime: out of memory")
		})
		req.Then(t, "0001.1.2.7.2", "d8d8dcd1", func(a req.A) {
			a.True(!strings.Contains(out, "not gofmt-formatted"), "a crash must not read as an unformatted file:\n%s", out)
		})
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
