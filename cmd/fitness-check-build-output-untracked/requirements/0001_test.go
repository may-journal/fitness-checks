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
	t.Run("0001.1.4.1.1 it fails with the hint to add `dist/` to `.gitignore`", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", nil)
		acceptance.Write(t, repo, map[string]string{"dist/app.js": "built\n"})
		out, code := buildOutputRun(t, repo)
		acceptance.Sees(t, out, code, 1, "Add dist/ to .gitignore (git check-ignore reports dist is not git-ignored)")
	})
	tracked := func(t *testing.T) (string, int) {
		repo := acceptance.Example(t, "happyRepo", map[string]string{"sub/dist/b.js": "built\n", "dist/a.js": "built\n"})
		acceptance.Commit(t, repo, map[string]string{".gitignore": buildOutputIgnore}, "chore: ignore dist")
		return buildOutputRun(t, repo)
	}
	t.Run("0001.1.4.2.1 one error lists them", func(t *testing.T) {
		t.Parallel()
		out, code := tracked(t)
		acceptance.Sees(t, out, code, 1, "Remove tracked files: dist/a.js, sub/dist/b.js")
	})
	t.Run("0001.1.4.2.2 it hints `git rm --cached`", func(t *testing.T) {
		t.Parallel()
		out, code := tracked(t)
		acceptance.Sees(t, out, code, 1, "(git rm --cached)")
	})
	t.Run("0001.1.4.3.1 those files join the same list", func(t *testing.T) {
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
	})
	t.Run("0001.1.4.3.2 a build folder without Java passes", func(t *testing.T) {
		t.Parallel()
		out, code := buildOutputRun(t, acceptance.Example(t, "happyRepo", map[string]string{"build/notes.txt": "kept\n"}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, buildOutputCheck)
	})
	imports := func(t *testing.T) (string, int) {
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			".gitignore": buildOutputIgnore,
			"src/app.test.ts": "import ok from './local.js';\n" +
				"import bad from '../dist/foo.js';\n" +
				"const z = require('pkg/dist/z.js');\n" +
				"const m = await import('../dist/mod.js');\n",
			"src/load.cts": "const c = require('../dist/c.cjs');\n",
		})
		return buildOutputRun(t, repo)
	}
	t.Run("0001.1.4.4.1 it fails as `path:line` for each `from`", func(t *testing.T) {
		t.Parallel()
		out, code := imports(t)
		acceptance.Sees(t, out, code, 1, `src/app.test.ts:2 imports build output: "../dist/foo.js"`)
	})
	t.Run("0001.1.4.4.2 it fails the same way for each `import()`", func(t *testing.T) {
		t.Parallel()
		out, code := imports(t)
		acceptance.Sees(t, out, code, 1, `src/app.test.ts:4 imports build output: "../dist/mod.js"`)
	})
	t.Run("0001.1.4.4.3 it fails the same way for each `require()`", func(t *testing.T) {
		t.Parallel()
		out, code := imports(t)
		acceptance.Sees(t, out, code, 1,
			`src/app.test.ts:3 imports build output: "pkg/dist/z.js"`,
			`src/load.cts:1 imports build output: "../dist/c.cjs"`)
	})
	t.Run("0001.1.4.4.4 tests are judged like any other file", func(t *testing.T) {
		t.Parallel()
		out, code := imports(t)
		acceptance.Sees(t, out, code, 1, `src/app.test.ts:2 imports build output`)
	})
	hidden := func(t *testing.T) (string, int) {
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			".gitignore": buildOutputIgnore,
			"src/a.ts":   "import x from './b.js';\n// import old from '../dist/old.js';\n/* const z = require('../dist/z.js'); */\n",
		})
		acceptance.Write(t, repo, map[string]string{"scratch.ts": "import x from '../dist/x.js';\n"})
		return buildOutputRun(t, repo)
	}
	t.Run("0001.1.4.5.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := hidden(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4.5.2 an untracked file passes too", func(t *testing.T) {
		t.Parallel()
		out, code := hidden(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 2 files scanned")
	})
	local := func(t *testing.T) (string, int) {
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			".gitignore": buildOutputIgnore,
			"src/a.ts":   "import x from './b.js';\n",
			"src/b.mts":  "export const b = require('./c.cjs');\n",
		})
		return buildOutputRun(t, repo)
	}
	t.Run("0001.1.4.6.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := local(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4.6.2 it counts the scanned sources plus one for the `dist` rules", func(t *testing.T) {
		t.Parallel()
		out, code := local(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 3 files scanned")
	})
	t.Run("0001.1.4.7.1 it passes with no files", func(t *testing.T) {
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
