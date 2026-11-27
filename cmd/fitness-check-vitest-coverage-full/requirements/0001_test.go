// Package requirements proves the vitest-coverage-full check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo, with Vitest answering saved
// responses.
package requirements

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m, "vitest")
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a repo with no JavaScript source passes with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.1.1", "857dd667", func(_ req.A) {
			out, code := vcfRun(t, acceptance.Example(t, "happyRepo", nil), nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "vitest-coverage-full")
			tooling := acceptance.Example(t, "happyRepo", map[string]string{"package.json": "{\"private\": true}\n", "vitest.config.mjs": vcfConfig("50"), "vite.config.ts": "export default {};\n"})
			out, code = vcfRun(t, tooling, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "vitest-coverage-full")
		})
	})
	t.Run("0001.1.2 a project threshold below full fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "24f22fc5", func(_ req.A) {
			rp := acceptance.StandIn(t, map[string][]acceptance.Response{"vitest": {vcfFullRun}})
			out, code := vcfRun(t, vcfRepo(t, map[string]string{"vitest.config.mjs": vcfConfig("90")}), rp.Env)
			acceptance.Sees(t, out, code, 1, "Vitest coverage thresholds must be 100 for branches, functions, lines, and statements.")
		})
	})
	t.Run("0001.1.3 a shared config threshold below full fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.3.1", "b6412d72", func(_ req.A) {
			shared := "node_modules/@mayjournal/fitness-shared/config/vitest.config.mjs"
			out, code := vcfRun(t, vcfRepo(t, map[string]string{shared: vcfConfig("90")}), nil)
			acceptance.Sees(t, out, code, 1, "Fitness-runner package must have Vitest coverage thresholds set to 100 for branches, functions, lines, and statements.")
		})
	})
	t.Run("0001.1.4 full thresholds with all code covered pass with the local config", func(t *testing.T) {
		t.Parallel()
		rp := acceptance.StandIn(t, map[string][]acceptance.Response{"vitest": {vcfFullRun}})
		out, code := vcfRun(t, vcfRepo(t, nil), rp.Env)
		req.Then(t, "0001.1.2.4.1", "07b8e753", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.1.2.4.2", "a745db12", func(_ req.A) {
			vcfRanWithLocalConfig(t, rp)
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 untested code shows the coverage failure", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.1.1", "d446ed37", func(_ req.A) {
			rp := acceptance.StandIn(t, map[string][]acceptance.Response{"vitest": {vcfShortRun}})
			repo := vcfRepo(t, map[string]string{"src/sub.ts": "export const sub = (a: number, b: number): number => a - b;\n"})
			out, code := vcfRun(t, repo, rp.Env)
			acceptance.Sees(t, out, code, 1, "Coverage for statements (50%) does not meet global threshold (100%)")
			vcfRanWithLocalConfig(t, rp)
		})
	})
	t.Run("0001.2.2 a missing Vitest shows an install hint", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.2.1", "c4203556", func(_ req.A) {
			out, code := vcfRun(t, vcfRepo(t, nil), nil)
			acceptance.ShowsInfo(t, out, code, "vitest-coverage-full", "vitest not found in node_modules/.bin (walking up from root) or on PATH. Install it (npm install -D vitest)")
		})
	})
	t.Run("0001.2.3 Vitest on PATH stands in while the binary is broken", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.3.1", "294d5301", func(_ req.A) {
			rp := acceptance.StandIn(t, map[string][]acceptance.Response{"vitest": {vcfFullRun}})
			repo := vcfRepo(t, map[string]string{"node_modules/.bin/vitest": "#!/bin/sh\nexit 1\n"})
			acceptance.MustDo(t, os.Chmod(filepath.Join(repo, "node_modules", ".bin", "vitest"), 0o644))
			out, code := vcfRun(t, repo, rp.Env)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			vcfRanWithLocalConfig(t, rp)
		})
	})
	t.Run("0001.2.4 an exit with no finding shows the last lines", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.4.1", "2167bd40", func(_ req.A) {
			rp := acceptance.StandIn(t, map[string][]acceptance.Response{"vitest": {{
				Match:  []string{"run --coverage"},
				Stderr: " MISSING DEPENDENCY  Cannot find dependency '@vitest/coverage-v8'\n",
				Exit:   1,
			}}})
			out, code := vcfRun(t, vcfRepo(t, nil), rp.Env)
			acceptance.ShowsInfo(t, out, code, "vitest-coverage-full", "vitest run --coverage exited 1 without a finding", "Cannot find dependency '@vitest/coverage-v8'")
			vcfRanWithLocalConfig(t, rp)
		})
	})
}

// vcfConfig is a local Vitest config holding the given coverage thresholds
// over src/**/*.ts, the way a consumer repo writes one.
func vcfConfig(branches string) string {
	return "export default {\n  test: {\n    coverage: {\n      include: ['src/**/*.ts'],\n      provider: 'v8',\n" +
		"      thresholds: { branches: " + branches + ", functions: 100, lines: 100, statements: 100 },\n" +
		"    },\n    include: ['src/**/*.test.ts'],\n  },\n};\n"
}

// vcfProject is a small JavaScript project whose one function is fully tested.
var vcfProject = map[string]string{
	"package.json":      `{"type": "module"}` + "\n",
	"vitest.config.mjs": vcfConfig("100"),
	"src/add.ts":        "export const add = (a: number, b: number): number => a + b;\n",
	"src/add.test.ts": "import { expect, test } from 'vitest';\nimport { add } from './add';\n\n" +
		"test('adds', () => {\n  expect(add(1, 2)).toBe(3);\n});\n",
}

// vcfRepo is happyRepo with vcfProject and then files written over it.
func vcfRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	all := map[string]string{}
	for _, set := range []map[string]string{vcfProject, files} {
		for name, body := range set {
			all[name] = body
		}
	}
	return acceptance.Example(t, "happyRepo", all)
}

func vcfRun(t *testing.T, repo string, env []string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, repo, env, "vitest-coverage-full")
}

// vcfTable is the coverage table Vitest's v8 provider prints after a run.
func vcfTable(rows string) string {
	rule := "----------|---------|----------|---------|---------|-------------------\n"
	return " % Coverage report from v8\n" + rule +
		"File      | % Stmts | % Branch | % Funcs | % Lines | Uncovered Line #s \n" + rule + rows + rule
}

// vcfFullRun is what `vitest run --coverage` prints when every line ran.
var vcfFullRun = acceptance.Response{
	Match: []string{"run --coverage"},
	Stdout: "\n RUN  v3.2.4 /repo\n\n \u2713 src/add.test.ts (1 test) 2ms\n\n" +
		" Test Files  1 passed (1)\n      Tests  1 passed (1)\n   Start at  10:15:02\n   Duration  412ms\n\n" +
		vcfTable("All files |     100 |      100 |     100 |     100 |                   \n"+
			" add.ts   |     100 |      100 |     100 |     100 |                   \n"),
}

// vcfShortRun is what `vitest run --coverage` prints when sub.ts never ran:
// the table on stdout, then one threshold error per metric on stderr.
var vcfShortRun = acceptance.Response{
	Match: []string{"run --coverage"},
	Stdout: "\n RUN  v3.2.4 /repo\n\n \u2713 src/add.test.ts (1 test) 2ms\n\n" +
		" Test Files  1 passed (1)\n      Tests  1 passed (1)\n   Start at  10:15:02\n   Duration  418ms\n\n" +
		vcfTable("All files |      50 |      100 |      50 |      50 |                   \n"+
			" add.ts   |     100 |      100 |     100 |     100 |                   \n"+
			" sub.ts   |       0 |      100 |       0 |       0 | 1                 \n"),
	Stderr: "ERROR: Coverage for lines (50%) does not meet global threshold (100%)\n" +
		"ERROR: Coverage for functions (50%) does not meet global threshold (100%)\n" +
		"ERROR: Coverage for statements (50%) does not meet global threshold (100%)\n",
	Exit: 1,
}

// vcfRanWithLocalConfig fails t unless fitness ran vitest once with
// `run --coverage` alone, leaving Vitest to load the repo's own config.
func vcfRanWithLocalConfig(t *testing.T, rp acceptance.Replays) {
	t.Helper()
	req.Asserted(t)
	if calls := rp.Calls("vitest"); len(calls) != 1 || calls[0] != "run --coverage" {
		t.Errorf("vitest calls = %q, want one \"run --coverage\"", calls)
	}
}
