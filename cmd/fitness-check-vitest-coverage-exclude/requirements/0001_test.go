// Package requirements proves the vitest-coverage-exclude check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
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
	t.Run("0001.1.1.1.1 it passes with no files judged", func(t *testing.T) {
		t.Parallel()
		out, code := coverageExclude(t, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "vitest-coverage-exclude")
	})
	t.Run("0001.1.1.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := coverageExclude(t, map[string]string{
			"package.json":     "{}\n",
			"vitest.config.js": coverageExcludeConfig("[]"),
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.1.3.1 it fails once per quoted entry telling Developer to remove it", func(t *testing.T) {
		t.Parallel()
		out, code := coverageExclude(t, map[string]string{
			"package.json":     "{}\n",
			"vitest.config.js": coverageExcludeConfig(`["**/*.test.ts", "dist"]`),
		})
		acceptance.Sees(t, out, code, 1, coverageExcludeRemove(`"**/*.test.ts"`), coverageExcludeRemove(`"dist"`))
		// The vitest field of package.json counts as a config.
		out, code = coverageExclude(t, map[string]string{
			"package.json": `{"vitest": {"test": {"coverage": {"exclude": ["src/foo.ts"]}}}}` + "\n",
		})
		acceptance.Sees(t, out, code, 1, coverageExcludeRemove(`"src/foo.ts"`))
	})
	t.Run("0001.1.1.3.2 it never reports an imported name", func(t *testing.T) {
		t.Parallel()
		// An imported name is invisible, so only the quoted path is named.
		out, code := coverageExclude(t, map[string]string{
			"package.json": "{}\n",
			"vitest.config.mjs": `import { DTS_GLOB } from "./globs.js";` + "\n" +
				`export default { test: { coverage: { exclude: [DTS_GLOB, "src/foo.ts"] } } };` + "\n",
		})
		acceptance.Sees(t, out, code, 1, coverageExcludeRemove(`"src/foo.ts"`))
		if strings.Contains(out, "DTS_GLOB") {
			t.Errorf("an imported name must not be reported:\n%s", out)
		}
	})
	t.Run("0001.1.1.4.1 the config file is the one judged", func(t *testing.T) {
		t.Parallel()
		out, code := coverageExclude(t, map[string]string{
			"package.json":     `{"vitest": {"coverage": {"exclude": ["src/foo.ts"]}}}` + "\n",
			"vitest.config.js": coverageExcludeConfig("[]"),
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.1.5.1 the shared config is judged", func(t *testing.T) {
		t.Parallel()
		out, code := coverageExclude(t, map[string]string{"package.json": "{}\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "vitest-coverage-exclude", "1")
	})
	t.Run("0001.1.1.5.2 one with an exclude fails", func(t *testing.T) {
		t.Parallel()
		out, code := coverageExclude(t, map[string]string{
			"package.json": "{}\n",
			"node_modules/@mayjournal/fitness-shared/config/vitest.config.js": coverageExcludeConfig(`["src/shared.ts"]`),
		})
		acceptance.Sees(t, out, code, 1, coverageExcludeRemove(`"src/shared.ts"`))
	})
	t.Run("0001.1.1.5.3 the built in copy passes", func(t *testing.T) {
		t.Parallel()
		out, code := coverageExclude(t, map[string]string{"package.json": "{}\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
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
