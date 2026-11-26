// Package requirements proves the node-version check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 `node-version` passes with zero files", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, nil, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "node-version")
	})
	t.Run("0001.1.2.2.1 it fails with `missing .nvmrc`", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, map[string]string{"package.json": "{}\n"}, nil)
		acceptance.Sees(t, out, code, 1, "missing .nvmrc")
	})
	t.Run("0001.1.2.3.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, nodeVersionRepo("1\n"), nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.3.2 it reads a `v` prefix as no prefix", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, nodeVersionRepo("v1.2.3\n"), nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.4.1 it fails with `does not satisfy .nvmrc`", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, nodeVersionRepo("999\n"), nil)
		acceptance.Sees(t, out, code, 1, "does not satisfy .nvmrc (requires 999.x)")
	})
	t.Run("0001.1.2.4.2 it gives the fix `nvm use`", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, nodeVersionRepo("999\n"), nil)
		acceptance.Sees(t, out, code, 1, "Run: nvm use")
	})
	t.Run("0001.1.2.5.1 it fails as an unsatisfiable `NaN.x`", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, nodeVersionRepo("latest\n"), nil)
		acceptance.Sees(t, out, code, 1, "does not satisfy .nvmrc (requires NaN.x). Run: nvm use")
	})
	t.Run("0001.1.2.6.1 it reports info with an install hint", func(t *testing.T) {
		t.Parallel()
		out, code := nodeVersionCheck(t, nodeVersionRepo("1\n"), []string{acceptance.PathWithout(t, "node")})
		acceptance.ShowsInfo(t, out, code, "node-version", "node not found on PATH. Install Node.js")
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
