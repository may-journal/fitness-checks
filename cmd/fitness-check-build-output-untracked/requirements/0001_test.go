// Package requirements proves the build-output-untracked check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a dist folder git does not ignore fails with the hint", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", nil)
		acceptance.Write(t, repo, map[string]string{"dist/app.js": "built\n"})
		out, code := buildOutputRun(t, repo)
		acceptance.Sees(t, out, code, 1, "Add dist/ to .gitignore (git check-ignore reports dist is not git-ignored)")
	})
	t.Run("0001.1.2 tracked files in an ignored dist fail in one list", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{"sub/dist/b.js": "built\n", "dist/a.js": "built\n"})
		acceptance.Commit(t, repo, map[string]string{".gitignore": buildOutputIgnore}, "chore: ignore dist")
		out, code := buildOutputRun(t, repo)
		acceptance.Sees(t, out, code, 1, "Remove tracked files: dist/a.js, sub/dist/b.js (git rm --cached)")
	})
	t.Run("0001.1.3 tracked Maven and Gradle output joins the list", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			"pom.xml":                          "<project/>\n",
			"target/app.jar":                   "built\n",
			"lib[1]/build.gradle":              "plugins { id 'java' }\n",
			"lib[1]/build/out.txt":             "built\n",
			"lib[1]/.gradle/state":             "cache\n",
			"lib1/build/kept.txt":              "not a build's output\n",
			"src/app/pom.xml":                  "<project/>\n",
			"src/app/target/app.jar":           "built\n",
			"classes/Old.class":                "compiled\n",
			"src/test/resources/Fixture.class": "a test input\n",
			"src/main/java/A.txt":              "source\n",
		})
		out, code := buildOutputRun(t, repo)
		acceptance.Sees(t, out, code, 1, "Remove tracked files: classes/Old.class, lib[1]/.gradle/state, lib[1]/build/out.txt, src/app/target/app.jar, target/app.jar (git rm --cached)")
		// A build folder in a repo without Java is no build's output.
		out, code = buildOutputRun(t, acceptance.Example(t, "happyRepo", map[string]string{"build/notes.txt": "kept\n"}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, buildOutputCheck)
	})
	t.Run("0001.1.4 an import from dist fails at its line", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			".gitignore": buildOutputIgnore,
			"src/app.test.ts": "import ok from './local.js';\n" +
				"import bad from '../dist/foo.js';\n" +
				"const z = require('pkg/dist/z.js');\n" +
				"const m = await import('../dist/mod.js');\n",
			"src/load.cts": "const c = require('../dist/c.cjs');\n",
		})
		out, code := buildOutputRun(t, repo)
		acceptance.Sees(t, out, code, 1,
			`src/app.test.ts:2 imports build output: "../dist/foo.js"`,
			`src/app.test.ts:3 imports build output: "pkg/dist/z.js"`,
			`src/app.test.ts:4 imports build output: "../dist/mod.js"`,
			`src/load.cts:1 imports build output: "../dist/c.cjs"`)
	})
	t.Run("0001.1.5 a commented-out or untracked import passes", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			".gitignore": buildOutputIgnore,
			"src/a.ts":   "import x from './b.js';\n// import old from '../dist/old.js';\n/* const z = require('../dist/z.js'); */\n",
		})
		acceptance.Write(t, repo, map[string]string{"scratch.ts": "import x from '../dist/x.js';\n"})
		out, code := buildOutputRun(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 2 files scanned")
	})
	t.Run("0001.1.6 an ignored dist with local imports passes, counting the dist rules once", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			".gitignore": buildOutputIgnore,
			"src/a.ts":   "import x from './b.js';\n",
			"src/b.mts":  "export const b = require('./c.cjs');\n",
		})
		out, code := buildOutputRun(t, repo)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 3 files scanned")
	})
	t.Run("0001.1.7 a repo with no output or TypeScript passes", func(t *testing.T) {
		t.Parallel()
		out, code := buildOutputRun(t, acceptance.Example(t, "happyRepo", nil))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, buildOutputCheck)
	})
}

// buildOutputCheck is the check these tests run.
const buildOutputCheck = "build-output-untracked"

// buildOutputIgnore is a .gitignore that ignores dist.
const buildOutputIgnore = "dist\n"

// buildOutputRun runs `fitness-install -- build-output-untracked` in repo.
func buildOutputRun(t *testing.T, repo string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, repo, nil, buildOutputCheck)
}
