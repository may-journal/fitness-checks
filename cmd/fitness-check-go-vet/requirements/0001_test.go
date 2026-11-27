// Package requirements proves the go-vet check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a vet problem at the repo root fails naming the finding", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule(".", goVetBad)))
		acceptance.Sees(t, out, code, 1, "✖ "+goVetFinding)
	})
	t.Run("0001.1.2 a vet problem in a subfolder module shows its folder", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule("api", goVetBad)))
		acceptance.Sees(t, out, code, 1, "✖ api: "+goVetFinding)
	})
	t.Run("0001.1.3 clean modules pass then count the modules vetted", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule("api", "package m\n\nfunc F() int { return 1 }\n")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "go-vet", "1")
	})
	t.Run("0001.1.4 a repo with no module to vet or a run scoped to changed files with none of them Go passes", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", map[string]string{"main.go": goVetBad}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-vet")
		// Fixtures under testdata and vendor are not modules to vet.
		files := goVetModule("testdata/fixture", goVetBad)
		for name, body := range goVetModule("vendor/dep", goVetBad) {
			files[name] = body
		}
		out, code = goVet(t, acceptance.Example(t, "happyRepo", files))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-vet")
		repo := acceptance.Example(t, "happyRepo", goVetModule("api", goVetBad))
		acceptance.Write(t, repo, map[string]string{"README.md": "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n\nMore text.\n"})
		acceptance.Git(t, repo, "add", "README.md")
		out, code = goVet(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-vet")
	})
	t.Run("0001.1.5 a missing Go shows the install hint", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule("api", goVetBad)), acceptance.PathWithout(t, "go"))
		acceptance.ShowsInfo(t, out, code, "go-vet", "Go not installed: see https://go.dev/dl")
	})
	t.Run("0001.1.6 a run that fails without a finding shows why", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", goVetModule(".", goVetBad))
		out, code := goVet(t, repo, acceptance.CrashingGo(t, "vet", "go: unexpected internal error"))
		acceptance.ShowsInfo(t, out, code, "go-vet", "go vet ./... exited 2 without a finding: go: unexpected internal error")
		out, code = goVet(t, repo, acceptance.KilledGo(t, "vet"))
		acceptance.ShowsInfo(t, out, code, "go-vet", "go vet")
	})
}

// goVetMod is a minimal go.mod for a Given's module.
const goVetMod = "module m\n\ngo 1.21\n"

// goVetBad is Go source the compiler accepts but go vet flags.
const goVetBad = "package m\n\nimport \"fmt\"\n\nfunc F() { fmt.Printf(\"%d\\n\", \"x\") }\n"

// goVetFinding is what go vet reports for goVetBad.
const goVetFinding = `m.go:5:24: fmt.Printf format %d has arg "x" of wrong type string`

// goVetModule returns the files of a module in dir holding src.
func goVetModule(dir, src string) map[string]string {
	return map[string]string{dir + "/go.mod": goVetMod, dir + "/m.go": src}
}

// goVet runs `fitness-install -- go-vet` in repo with env added.
func goVet(t *testing.T, repo string, env ...string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, repo, env, "go-vet")
}
