// Package requirements proves the vitest-coverage-exclude check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
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
	t.Run("0001.1.1 a repo with no package json passes with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.1.1", "4b30883c", func(_ req.A) {
			out, code := coverageExclude(t, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "vitest-coverage-exclude")
		})
	})
	t.Run("0001.1.2 an empty exclude passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "07b8e753", func(_ req.A) {
			out, code := coverageExclude(t, map[string]string{
				"package.json":     "{}\n",
				"vitest.config.js": coverageExcludeConfig("[]"),
			})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.3 each exclude entry fails by name", func(t *testing.T) {
		t.Parallel()
		out1, code1 := coverageExclude(t, map[string]string{
			"package.json":     "{}\n",
			"vitest.config.js": coverageExcludeConfig(`["**/*.test.ts", "dist"]`),
		})
		// The vitest field of package.json counts as a config.
		out2, code2 := coverageExclude(t, map[string]string{
			"package.json": `{"vitest": {"test": {"coverage": {"exclude": ["src/foo.ts"]}}}}` + "\n",
		})
		// An imported name is invisible, so only the quoted path is named.
		out3, code3 := coverageExclude(t, map[string]string{
			"package.json": "{}\n",
			"vitest.config.mjs": `import { DTS_GLOB } from "./globs.js";` + "\n" +
				`export default { test: { coverage: { exclude: [DTS_GLOB, "src/foo.ts"] } } };` + "\n",
		})
		req.Then(t, "0001.1.1.3.1", "825597c0", func(_ req.A) {
			acceptance.Sees(t, out1, code1, 1, coverageExcludeRemove(`"**/*.test.ts"`), coverageExcludeRemove(`"dist"`))
			acceptance.Sees(t, out2, code2, 1, coverageExcludeRemove(`"src/foo.ts"`))
			acceptance.Sees(t, out3, code3, 1, coverageExcludeRemove(`"src/foo.ts"`))
		})
		req.Then(t, "0001.1.1.3.2", "3d782499", func(a req.A) {
			a.True(!strings.Contains(out3, "DTS_GLOB"), "an imported name must not be reported:\n%s", out3)
		})
	})
	t.Run("0001.1.4 the config file beside a package json field wins", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.4.1", "1274a8e0", func(_ req.A) {
			out, code := coverageExclude(t, map[string]string{
				"package.json":     `{"vitest": {"coverage": {"exclude": ["src/foo.ts"]}}}` + "\n",
				"vitest.config.js": coverageExcludeConfig("[]"),
			})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.5 with no config the shared config is judged", func(t *testing.T) {
		t.Parallel()
		out, code := coverageExclude(t, map[string]string{
			"package.json": "{}\n",
			"node_modules/@mayjournal/fitness-shared/config/vitest.config.js": coverageExcludeConfig(`["src/shared.ts"]`),
		})
		out2, code2 := coverageExclude(t, map[string]string{"package.json": "{}\n"})
		req.Then(t, "0001.1.1.5.1", "01f73c4d", func(_ req.A) {
			acceptance.PassedWith(t, out2, "vitest-coverage-exclude", "1")
		})
		req.Then(t, "0001.1.1.5.2", "7a365f76", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, coverageExcludeRemove(`"src/shared.ts"`))
		})
		req.Then(t, "0001.1.1.5.3", "2b914c7e", func(_ req.A) {
			acceptance.Sees(t, out2, code2, 0, "All 1 checks passed")
		})
	})
}

// coverageExclude runs `fitness-install -- vitest-coverage-exclude` on
// happyRepo with files written over it.
func coverageExclude(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "vitest-coverage-exclude")
}

// coverageExcludeRemove is the failure line naming one quoted entry.
func coverageExcludeRemove(quoted string) string {
	return "Vitest coverage exclude must be empty, so coverage judges every file; remove: " + quoted
}

// coverageExcludeConfig is a CommonJS Vitest config with exclude as given.
func coverageExcludeConfig(exclude string) string {
	return "module.exports = { test: { coverage: { exclude: " + exclude + " } } };\n"
}
