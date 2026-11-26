// Package requirements proves the go-vet check's docs/requirements as a user
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
	t.Run("0001.1.2.1.1 it fails naming the finding as `file.go:line:col: message`", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule(".", goVetBad)))
		acceptance.Sees(t, out, code, 1, "✖ "+goVetFinding)
	})
	t.Run("0001.1.2.2.1 each finding is prefixed with that module's folder", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule("api", goVetBad)))
		acceptance.Sees(t, out, code, 1, "✖ api: "+goVetFinding)
	})
	cleanModule := func(t *testing.T) (string, int) {
		return goVet(t, acceptance.Example(t, "happyRepo", goVetModule("api", "package m\n\nfunc F() int { return 1 }\n")))
	}
	t.Run("0001.1.2.3.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := cleanModule(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.3.2 it counts the modules vetted", func(t *testing.T) {
		t.Parallel()
		out, _ := cleanModule(t)
		acceptance.PassedWith(t, out, "go-vet", "1")
	})
	t.Run("0001.1.2.4.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", map[string]string{"main.go": goVetBad}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-vet")
	})
	t.Run("0001.1.2.4.2 it never runs `go vet`", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", map[string]string{"main.go": goVetBad}), acceptance.CrashingGo(t, "vet", "go vet must not run"))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if acceptance.Row(out, "go-vet") == "" {
			t.Errorf("no go-vet row in:\n%s", out)
		}
		acceptance.PassedWithNoFiles(t, out, "go-vet")
	})
	t.Run("0001.1.2.4.3 a `testdata` module counts as none", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule("testdata/fixture", goVetBad)))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-vet")
	})
	t.Run("0001.1.2.4.4 a `vendor` module counts as none", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule("vendor/dep", goVetBad)))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-vet")
	})
	t.Run("0001.1.2.5.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", goVetModule("api", goVetBad))
		acceptance.Write(t, repo, map[string]string{"README.md": "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n\nMore text.\n"})
		acceptance.Git(t, repo, "add", "README.md")
		out, code := goVet(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-vet")
	})
	t.Run("0001.1.2.6.1 it reports info with the hint `see https://go.dev/dl`", func(t *testing.T) {
		t.Parallel()
		out, code := goVet(t, acceptance.Example(t, "happyRepo", goVetModule("api", goVetBad)), acceptance.PathWithout(t, "go"))
		acceptance.ShowsInfo(t, out, code, "go-vet", "Go not installed: see https://go.dev/dl")
	})
	t.Run("0001.1.2.7.1 it reports info quoting what went wrong", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", goVetModule(".", goVetBad))
		out, code := goVet(t, repo, acceptance.CrashingGo(t, "vet", "go: unexpected internal error"))
		acceptance.ShowsInfo(t, out, code, "go-vet", "go vet ./... exited 2 without a finding: go: unexpected internal error")
	})
	t.Run("0001.1.2.7.2 a failure to start counts the same", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", goVetModule(".", goVetBad))
		unstartable := t.TempDir()
		acceptance.Write(t, unstartable, map[string]string{"go": "#!/nonexistent/interpreter\n"})
		out, code := goVet(t, repo, "PATH="+unstartable+string(os.PathListSeparator)+os.Getenv("PATH"))
		acceptance.ShowsInfo(t, out, code, "go-vet", "go vet")
	})
	t.Run("0001.1.2.7.3 a kill counts the same", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", goVetModule(".", goVetBad))
		out, code := goVet(t, repo, acceptance.KilledGo(t, "vet"))
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
