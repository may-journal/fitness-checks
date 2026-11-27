// Package requirements proves the node-version check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a repo with no package json passes with zero files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.1.1", "7fe38931", func(_ req.A) {
			out, code := nodeVersionCheck(t, nil, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "node-version")
		})
	})
	t.Run("0001.1.2 a missing nvmrc fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "17b409d8", func(_ req.A) {
			out, code := nodeVersionCheck(t, map[string]string{"package.json": "{}\n"}, nil)
			acceptance.Sees(t, out, code, 1, "missing .nvmrc")
		})
	})
	t.Run("0001.1.3 the active Node meets an nvmrc", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.3.1", "07b8e753", func(_ req.A) {
			nvmrcSees(t, "1\n", 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.3.2", "df850705", func(_ req.A) {
			nvmrcSees(t, "v1.2.3\n", 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.4 an nvmrc outranks the active Node", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, nodeVersionRepo("999\n"), nil)
		req.Then(t, "0001.1.2.4.1", "a053d73f", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "does not satisfy .nvmrc (requires 999.x)")
		})
		req.Then(t, "0001.1.2.4.2", "30034ad3", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "Run: nvm use")
		})
	})
	t.Run("0001.1.5 an nvmrc without a number fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.5.1", "c3e4b0a3", func(_ req.A) {
			nvmrcSees(t, "latest\n", 1, "does not satisfy .nvmrc (requires NaN.x). Run: nvm use")
		})
	})
	t.Run("0001.1.6 a missing node reports info with an install hint", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "c4203556", func(_ req.A) {
			out, code := nodeVersionCheck(t, nodeVersionRepo("1\n"), []string{acceptance.PathWithout(t, "node")})
			acceptance.ShowsInfo(t, out, code, "node-version", "node not found on PATH. Install Node.js")
		})
	})
}

// nodeVersionCheck runs `fitness-install -- node-version` on happyRepo with
// files written over it and env added.
func nodeVersionCheck(t *testing.T, files map[string]string, env []string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), env, "node-version")
}

// nodeVersionRepo is a JavaScript repo whose `.nvmrc` holds nvmrc.
func nodeVersionRepo(nvmrc string) map[string]string {
	return map[string]string{"package.json": "{}\n", ".nvmrc": nvmrc}
}

// nvmrcSees runs the check on a repo whose `.nvmrc` holds nvmrc and expects code and want.
func nvmrcSees(t *testing.T, nvmrc string, wantCode int, want string) {
	t.Helper()
	out, code := nodeVersionCheck(t, nodeVersionRepo(nvmrc), nil)
	acceptance.Sees(t, out, code, wantCode, want)
}
