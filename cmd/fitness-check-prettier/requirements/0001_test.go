// Package requirements proves the prettier check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo, with Prettier and Biome answering saved responses.
package requirements

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain lets this test binary stand in for Prettier and Biome. Run as
// either, it answers from saved responses instead of running tests.
func TestMain(m *testing.M) {
	acceptance.ServeStandIns([]string{"prettier", "biome"})
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

// sharedConfigRun runs a repo with no Prettier config of its own, which the
// shared config judges, and returns what Prettier was given.
func sharedConfigRun(t *testing.T) (string, int, acceptance.Replays) {
	t.Helper()
	files := prettierJS(acceptance.With(map[string]string{".prettierrc.json": "", "app.js": prettierShared}, prettierSynced))
	shared := prettierAnswers(prettierClean, "--config", "prettier.config.cjs", "--check", "app.js", ".github/ISSUE_TEMPLATE/plan.md", ".github/fitness-setup.md")
	return prettierRun(t, files, []acceptance.Response{shared, prettierWarn})
}

// sharedConfig is the shared config file text fitness handed Prettier.
func sharedConfig(t *testing.T, rp acceptance.Replays) string {
	t.Helper()
	calls := rp.Calls("prettier")
	if len(calls) != 1 {
		t.Fatalf("want one prettier call, got %q", calls)
	}
	return configFrom(t, calls[0])
}

// stagedRun stages a messy file and a symbolic link beside it, leaves app.js
// unstaged, and runs the check.
func stagedRun(t *testing.T) (string, int, acceptance.Replays) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", prettierJS(nil))
	acceptance.Write(t, repo, map[string]string{"app.js": prettierMessy, "ci_scripts/run.sh": "#!/bin/sh\necho ready\n"})
	acceptance.MustDo(t, os.Symlink("run.sh", filepath.Join(repo, "ci_scripts", "ci_post_clone.sh")))
	rp := acceptance.StandIn(t, map[string][]acceptance.Response{"prettier": {prettierLinkRefused, prettierAnswers(prettierWarn, "app.js"), prettierClean}})
	acceptance.Git(t, repo, "add", "ci_scripts")
	out, code := acceptance.Fitness(t, repo, rp.Env, "prettier")
	return out, code, rp
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.2.1.1 `prettier` judges every tracked file with it", func(t *testing.T) {
		t.Parallel()
		files := prettierJS(map[string]string{".prettierrc.json": "", "package.json": prettierPackageField})
		out, code, rp := prettierRun(t, files, []acceptance.Response{prettierAnswers(prettierWarn, "--config"), prettierClean})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		prettierCalledWith(t, rp, []string{"--check", "package.json"}, "--config")
		_, _, rp = prettierRun(t, prettierJS(map[string]string{".prettierrc.json": `{ "singleQuote": true, "semi": false }` + "\n"}), []acceptance.Response{prettierWarnColored})
		prettierCalledWith(t, rp, []string{"--check", "app.js"}, "--config")
	})
	t.Run("0001.1.2.1.2 a file it would reformat fails by name", func(t *testing.T) {
		t.Parallel()
		out, code, _ := prettierRun(t, prettierJS(map[string]string{".prettierrc.json": `{ "singleQuote": true, "semi": false }` + "\n"}), []acceptance.Response{prettierWarnColored})
		acceptance.Sees(t, out, code, 1, "✖ app.js")
	})
	t.Run("0001.1.2.2.1 it passes without running Prettier", func(t *testing.T) {
		t.Parallel()
		out, code, rp := prettierRun(t, map[string]string{"app.js": prettierMessy}, []acceptance.Response{prettierWarn})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if calls := rp.Calls("prettier"); len(calls) > 0 {
			t.Errorf("prettier was called without a package.json: %q", calls)
		}
	})
	t.Run("0001.1.2.3.1 the shared config judges with single quotes", func(t *testing.T) {
		t.Parallel()
		out, code, rp := sharedConfigRun(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		prettierCalledWith(t, rp, []string{"--config", "prettier.config.cjs", "--ignore-unknown", "--ignore-path", "--with-node-modules", "--check", "app.js"})
		if config := sharedConfig(t, rp); !strings.Contains(config, "singleQuote: true") {
			t.Errorf("the shared config must use single quotes; missing %q in:\n%s", "singleQuote: true", config)
		}
	})
	t.Run("0001.1.2.3.2 the synced `.github` templates keep Prettier defaults", func(t *testing.T) {
		t.Parallel()
		out, code, rp := sharedConfigRun(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		config := sharedConfig(t, rp)
		for _, want := range []string{"'.github/ISSUE_TEMPLATE/**/*.md', '.github/fitness-setup.md'", "singleQuote: false", "fromRepoRoot"} {
			if !strings.Contains(config, want) {
				t.Errorf("the shared config must keep the synced files at Prettier defaults; missing %q in:\n%s", want, config)
			}
		}
	})
	t.Run("0001.1.2.4.1 only the staged files are checked", func(t *testing.T) {
		t.Parallel()
		out, code, rp := stagedRun(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		prettierCalledWith(t, rp, []string{"--check", "ci_scripts/run.sh"}, "app.js")
	})
	t.Run("0001.1.2.4.2 a symbolic link is left out", func(t *testing.T) {
		t.Parallel()
		out, code, rp := stagedRun(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		prettierCalledWith(t, rp, []string{"--check", "ci_scripts/run.sh"}, "ci_scripts/ci_post_clone.sh")
	})
	t.Run("0001.1.2.5.1 Prettier writes in place instead of checking", func(t *testing.T) {
		t.Parallel()
		out, code, rp := prettierRun(t, prettierJS(map[string]string{"app.js": prettierMessy}), []acceptance.Response{prettierAnswers(prettierWritten, "--write ."), prettierWarn}, "--write", ".")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		prettierCalledWith(t, rp, []string{"--write ."}, "--check")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	ignored := func(t *testing.T) (string, int) {
		return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", map[string]string{".prettierignore": "README.md\n"}), nil, "prettier")
	}
	t.Run("0001.2.1.1.1 it fails telling to delete it", func(t *testing.T) {
		t.Parallel()
		out, code := ignored(t)
		acceptance.Sees(t, out, code, 1, ".prettierignore exists; Prettier checks every tracked file, so delete it")
	})
	t.Run("0001.2.1.1.2 it tells to fix the findings", func(t *testing.T) {
		t.Parallel()
		out, code := ignored(t)
		acceptance.Sees(t, out, code, 1, "and fix the findings instead")
	})
	t.Run("0001.2.1.2.1 it reports info with an install hint", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", prettierJS(nil)), []string{acceptance.PathWithout(t, "prettier")}, "prettier")
		acceptance.ShowsInfo(t, out, code, "prettier", "Prettier not installed: npm install --save-dev prettier")
	})
	t.Run("0001.2.1.3.1 it reports info quoting Prettier's error", func(t *testing.T) {
		t.Parallel()
		noPlugin := acceptance.Response{Stderr: "[error] Cannot find package 'prettier-plugin-packagejson' imported from /shared/prettier.config.cjs\n", Exit: 2}
		out, code, _ := prettierRun(t, prettierJS(nil), []acceptance.Response{noPlugin})
		acceptance.ShowsInfo(t, out, code, "prettier", "prettier exited 2 without a finding: [error] Cannot find package 'prettier-plugin-packagejson'")
	})
	t.Run("0001.2.1.4.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		unparsable := acceptance.Response{Stderr: "Checking formatting...\n[error] app.js: SyntaxError: Unexpected token (1:7)\n", Exit: 2}
		out, code, _ := prettierRun(t, prettierJS(nil), []acceptance.Response{unparsable})
		acceptance.Sees(t, out, code, 1, "[error] app.js: SyntaxError: Unexpected token (1:7)")
	})
}

// biomeFormat runs the check in a repo with a Biome config at declared version
// major[0]. It returns what Biome was given.
func biomeFormat(t *testing.T, major [3]string) (string, int, acceptance.Replays) {
	t.Helper()
	files := map[string]string{"biome.jsonc": "{ \"formatter\": { \"enabled\": true } }\n"}
	var report strings.Builder
	for i := 1; i <= biomeMany; i++ {
		name := fmt.Sprintf("app%d.js", i)
		files[name] = prettierMessy
		fmt.Fprintf(&report, biomeReformat, name)
	}
	files["package.json"] = fmt.Sprintf(biomePackage, major[0])
	repo := acceptance.Example(t, "happyRepo", files)
	reformat := acceptance.Response{
		Match:  []string{"format " + fmt.Sprintf(biomeFlags, major[2]), "app1.js"},
		Stdout: report.String(),
		Stderr: "format ━━━━━━━━━━\n\n  × Some errors were emitted while running checks.\n",
		Exit:   1,
	}
	rp := acceptance.StandIn(t, map[string][]acceptance.Response{"biome": {biomeVersion(major[1]), reformat}})
	out, code := acceptance.Fitness(t, repo, rp.Env, "prettier")
	return out, code, rp
}

// biomeTooOldRun runs the check with Biome older than 1.8 installed.
func biomeTooOldRun(t *testing.T) (string, int, acceptance.Replays) {
	t.Helper()
	files := map[string]string{"biome.jsonc": "{ \"formatter\": { \"enabled\": true } }\n", "app.js": prettierMessy}
	files["package.json"] = fmt.Sprintf(biomePackage, "latest")
	rp := acceptance.StandIn(t, map[string][]acceptance.Response{"biome": {biomeVersion(biomeTooOld)}})
	out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), rp.Env, "prettier")
	return out, code, rp
}

// biomeWriteRun runs the check with Biome and args, as a --write argument
// with or without a Prettier-only flag.
func biomeWriteRun(t *testing.T, args ...string) (string, int, acceptance.Replays) {
	t.Helper()
	files := prettierJS(map[string]string{".prettierrc.json": "", "biome.jsonc": "{}\n", "package.json": fmt.Sprintf(biomePackage, "2.4.15")})
	written := acceptance.Response{Match: []string{"format " + fmt.Sprintf(biomeFlags, "none") + " --write app.js"}, Stdout: "Formatted 1 file in 2ms. Fixed 1 file.\n"}
	rp := acceptance.StandIn(t, map[string][]acceptance.Response{"biome": {biomeVersion("2.4.15"), written}})
	out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), rp.Env, append([]string{"prettier"}, args...)...)
	return out, code, rp
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.2.1.1 `prettier` fails naming the file", func(t *testing.T) {
		t.Parallel()
		for _, major := range biomeMajors {
			out, code, _ := biomeFormat(t, major)
			acceptance.Sees(t, out, code, 1, "✖ app1.js", "✖ app25.js")
		}
	})
	t.Run("0001.3.2.1.2 it needs no Prettier", func(t *testing.T) {
		t.Parallel()
		for _, major := range biomeMajors {
			out, code, rp := biomeFormat(t, major)
			acceptance.Sees(t, out, code, 1)
			if !rp.Called("biome", "format", "README.md") || strings.Contains(out, "Prettier not installed") {
				t.Errorf("Biome must check every tracked file in place of Prettier; calls %q:\n%s", rp.Calls("biome"), out)
			}
		}
	})
	t.Run("0001.3.2.2.1 it fails naming the version", func(t *testing.T) {
		t.Parallel()
		out, code, _ := biomeTooOldRun(t)
		acceptance.Sees(t, out, code, 1, "Biome "+biomeTooOld+" found; Biome 1.8 or later is required")
	})
	t.Run("0001.3.2.2.2 it never runs Biome", func(t *testing.T) {
		t.Parallel()
		out, code, rp := biomeTooOldRun(t)
		acceptance.Sees(t, out, code, 1)
		if rp.Called("biome", "--reporter=github") {
			t.Errorf("Biome %s cannot take --reporter=github, yet it ran: %q", biomeTooOld, rp.Calls("biome"))
		}
	})
	t.Run("0001.3.2.3.1 Biome writes in place", func(t *testing.T) {
		t.Parallel()
		out, code, rp := biomeWriteRun(t, "--write", "--ignore-unknown", "app.js")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if !rp.Called("biome", "--write", "app.js") {
			t.Errorf("Biome must format in place; calls %q", rp.Calls("biome"))
		}
	})
	t.Run("0001.3.2.3.2 it gets only its own flags", func(t *testing.T) {
		t.Parallel()
		out, code, rp := biomeWriteRun(t, "--write", "--ignore-unknown", "app.js")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		for _, call := range rp.Calls("biome") {
			if strings.Contains(call, "--ignore-unknown") || strings.Contains(call, "README.md") {
				t.Errorf("Biome must get only its own flags and the given path; got %q", call)
			}
		}
	})
	t.Run("0001.3.2.4.1 it fails naming the flag", func(t *testing.T) {
		t.Parallel()
		out, code, _ := biomeWriteRun(t, "--write", "--print-width", "100", "app.js")
		acceptance.Sees(t, out, code, 1, "Biome formats this repo and has no counterpart for Prettier's --print-width")
	})
}

// biomePackage is a JavaScript project that installs Biome at version %s.
const biomePackage = "{\"name\": \"app\", \"private\": true, \"devDependencies\": {\"@biomejs/biome\": \"%s\"}}\n"

// biomeFlags are the flags fitness passes Biome before the paths, with %s
// the cap that lifts its default of 20 diagnostics.
const biomeFlags = "--reporter=github --max-diagnostics=%s --no-errors-on-unmatched --files-ignore-unknown=true"

// biomeMajors pairs declared version, installed version and cap: Biome 1.8
// takes at most 65535, and `latest` hides the version.
var biomeMajors = [][3]string{{"2.4.15", "2.4.15", "none"}, {"^1.9.4", "1.9.4", "none"}, {"latest", "1.8.3", "65535"}}

// biomeTooOld is a Biome before 1.8, whose CLI rejects `--reporter=github`,
// declared under a spec that hides its version.
const biomeTooOld = "1.7.3"

// biomeVersion is real `biome --version` output for version v.
func biomeVersion(v string) acceptance.Response {
	return acceptance.Response{Match: []string{"--version"}, Stdout: "Version: " + v + "\n"}
}

// biomeMany is more findings than Biome's default cap of 20 diagnostics.
const biomeMany = 25

// prettierJS is a formatted project with its own Prettier config, which
// keeps single quotes.
func prettierJS(files map[string]string) map[string]string {
	return acceptance.With(map[string]string{
		"package.json":     "{}\n",
		".prettierrc.json": `{ "singleQuote": true }` + "\n",
		"app.js":           prettierFormatted,
	}, files)
}

// prettierFormatted and prettierMessy are one statement, formatted and not.
const (
	prettierFormatted = "const answer = { value: 42 };\n"
	prettierMessy     = "const answer={value:42}\n"
)

// prettierShared is formatted to the shared config's single quotes, which
// Prettier's own defaults would reformat.
const prettierShared = "const greeting = 'hi';\n"

// prettierPackageField is a formatted package.json holding the Prettier
// config, in place of a config file.
const prettierPackageField = "{\n  \"prettier\": {\n    \"singleQuote\": true\n  }\n}\n"

// What Prettier 3 prints: `--check` on formatted files, `--check` naming
// app.js, the same with a colored [warn] prefix, and `--write`.
var (
	prettierClean = acceptance.Response{Stdout: "Checking formatting...\nAll matched files use Prettier code style!\n"}
	prettierWarn  = acceptance.Response{
		Stdout: "Checking formatting...\n",
		Stderr: "[warn] app.js\n[warn] Code style issues found in the above file. Run Prettier with --write to fix.\n",
		Exit:   1,
	}
	prettierWarnColored = acceptance.Response{
		Stdout: "Checking formatting...\n",
		Stderr: "\x1b[33m[warn]\x1b[39m app.js\n\x1b[33m[warn]\x1b[39m Code style issues found in the above file. Run Prettier with --write to fix.\n",
		Exit:   1,
	}
	prettierWritten = acceptance.Response{Stdout: "CHANGELOG.md 9ms (unchanged)\nREADME.md 4ms (unchanged)\napp.js 12ms\npackage.json 1ms (unchanged)\n"}
)

// prettierAnswers matches a acceptance.Response to calls holding every match string.
func prettierAnswers(r acceptance.Response, match ...string) acceptance.Response {
	r.Match = match
	return r
}

// prettierRun runs `fitness-install -- prettier <args>` on happyRepo with
// files written over it and a Prettier stand-in answering from saved.
func prettierRun(t *testing.T, files map[string]string, saved []acceptance.Response, args ...string) (string, int, acceptance.Replays) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	rp := acceptance.StandIn(t, map[string][]acceptance.Response{"prettier": saved})
	out, code := acceptance.Fitness(t, repo, rp.Env, append([]string{"prettier"}, args...)...)
	return out, code, rp
}

// prettierCalledWith fails t unless fitness called Prettier with want and
// without each of not.
func prettierCalledWith(t *testing.T, rp acceptance.Replays, want []string, not ...string) {
	t.Helper()
	if !rp.Called("prettier", want...) {
		t.Errorf("prettier was not called with %q; calls: %q", want, rp.Calls("prettier"))
	}
	for _, n := range not {
		if rp.Called("prettier", n) {
			t.Errorf("prettier was called with %q; calls: %q", n, rp.Calls("prettier"))
		}
	}
}

// biomeReformat is real `biome format --reporter=github` output naming a
// file it would reformat.
const biomeReformat = "::error title=format,file=%s,line=1,endLine=1,col=2,endColumn=2::Formatter would have printed the following content:\n"

// prettierSynced are the files may-journal/.github syncs into every repo,
// formatted with Prettier's defaults: double-quoted front matter.
var prettierSynced = map[string]string{
	".github/ISSUE_TEMPLATE/plan.md": "---\nname: Plan\ntitle: \"[Plan] Short name of the work\"\n---\n\n> REPLACE-ME\n",
	".github/fitness-setup.md":       "# Fitness setup\n\nRun \"fitness\" before each push.\n",
}

// configFrom returns the file a Prettier call names after --config.
func configFrom(t *testing.T, call string) string {
	t.Helper()
	args := strings.Fields(call)
	i := slices.Index(args, "--config")
	if i < 0 || i+1 == len(args) {
		t.Fatalf("prettier was called without --config: %q", call)
	}
	raw, err := os.ReadFile(args[i+1])
	acceptance.MustDo(t, err)
	return string(raw)
}

// prettierLinkRefused is what Prettier 3 prints when a symbolic link is
// named on its command line.
var prettierLinkRefused = acceptance.Response{
	Match:  []string{"ci_scripts/ci_post_clone.sh"},
	Stderr: "[error] Explicitly specified pattern \"ci_scripts/ci_post_clone.sh\" is a symbolic link.\n",
	Exit:   2,
}
