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
	"github.com/may-journal/fitness-checks/req"
)

// TestMain lets this test binary stand in for Prettier and Biome. Run as
// either, it answers from saved responses instead of running tests.
func TestMain(m *testing.M) {
	acceptance.ServeStandIns([]string{"prettier", "biome"})
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 the repo's own config judges", func(t *testing.T) {
		t.Parallel()
		out, code, rp := prettierRun(t, prettierJS(map[string]string{".prettierrc.json": `{ "singleQuote": true, "semi": false }` + "\n"}), []acceptance.Response{prettierWarnColored})
		files := prettierJS(map[string]string{".prettierrc.json": "", "package.json": prettierPackageField})
		out2, code2, rp2 := prettierRun(t, files, []acceptance.Response{prettierAnswers(prettierWarn, "--config"), prettierClean})
		req.Then(t, "0001.1.2.1.1", "eec5ee4b", func(_ req.A) {
			prettierCalledWith(t, rp, []string{"--check", "app.js"}, "--config")
			acceptance.Sees(t, out2, code2, 0, "All 1 checks passed")
			prettierCalledWith(t, rp2, []string{"--check", "package.json"}, "--config")
		})
		req.Then(t, "0001.1.2.1.2", "eebad798", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "✖ app.js")
		})
	})
	t.Run("0001.1.2 a repo with no package json passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.2.1", "47ff4d27", func(_ req.A) {
			out, code, rp := prettierRun(t, map[string]string{"app.js": prettierMessy}, []acceptance.Response{prettierWarn})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			if calls := rp.Calls("prettier"); len(calls) > 0 {
				t.Errorf("prettier was called without a package.json: %q", calls)
			}
		})
	})
	t.Run("0001.1.3 the shared config judges a repo with none", func(t *testing.T) {
		t.Parallel()
		files := prettierJS(acceptance.With(map[string]string{".prettierrc.json": "", "app.js": prettierShared}, prettierSynced))
		shared := prettierAnswers(prettierClean, "--config", "prettier.config.cjs", "--check", "app.js", ".github/ISSUE_TEMPLATE/plan.md", ".github/fitness-setup.md")
		out, code, rp := prettierRun(t, files, []acceptance.Response{shared, prettierWarn})
		calls := rp.Calls("prettier")
		if len(calls) != 1 {
			t.Fatalf("want one prettier call, got %q", calls)
		}
		config := configFrom(t, calls[0])
		req.Then(t, "0001.1.2.3.1", "79eba689", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			prettierCalledWith(t, rp, []string{"--config", "prettier.config.cjs", "--ignore-unknown", "--ignore-path", "--with-node-modules", "--check", "app.js"})
			a.Contains(config, "singleQuote: true")
			a.Contains(config, "fromRepoRoot")
		})
		req.Then(t, "0001.1.2.3.2", "affbf1dd", func(a req.A) {
			a.Contains(config, "'.github/ISSUE_TEMPLATE/**/*.md', '.github/fitness-setup.md'")
			a.Contains(config, "singleQuote: false")
		})
	})
	t.Run("0001.1.4 staged files are checked and a symbolic link is left out", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", prettierJS(nil))
		acceptance.Write(t, repo, map[string]string{"app.js": prettierMessy, "ci_scripts/run.sh": "#!/bin/sh\necho ready\n"})
		acceptance.MustDo(t, os.Symlink("run.sh", filepath.Join(repo, "ci_scripts", "ci_post_clone.sh")))
		rp := acceptance.StandIn(t, map[string][]acceptance.Response{"prettier": {prettierLinkRefused, prettierAnswers(prettierWarn, "app.js"), prettierClean}})
		acceptance.Git(t, repo, "add", "ci_scripts")
		out, code := acceptance.Fitness(t, repo, rp.Env, "prettier")
		req.Then(t, "0001.1.2.4.1", "c5e8090c", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			prettierCalledWith(t, rp, []string{"--check", "ci_scripts/run.sh"}, "app.js")
		})
		req.Then(t, "0001.1.2.4.2", "31ca9cb7", func(_ req.A) {
			prettierCalledWith(t, rp, []string{"--check"}, "ci_scripts/ci_post_clone.sh")
		})
	})
	t.Run("0001.1.5 a write argument replaces the check", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.5.1", "fb283267", func(_ req.A) {
			out, code, rp := prettierRun(t, prettierJS(map[string]string{"app.js": prettierMessy}), []acceptance.Response{prettierAnswers(prettierWritten, "--write ."), prettierWarn}, "--write", ".")
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			prettierCalledWith(t, rp, []string{"--write ."}, "--check")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a prettierignore file fails telling to delete it and fix the findings", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", map[string]string{".prettierignore": "README.md\n"}), nil, "prettier")
		req.Then(t, "0001.2.1.1.1", "354a48ef", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, ".prettierignore exists; Prettier checks every tracked file, so delete it")
		})
		req.Then(t, "0001.2.1.1.2", "4aea09df", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "so delete it and fix the findings instead")
		})
	})
	t.Run("0001.2.2 a missing Prettier reports info with an install hint", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.2.1", "c4203556", func(_ req.A) {
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", prettierJS(nil)), []string{acceptance.PathWithout(t, "prettier")}, "prettier")
			acceptance.ShowsInfo(t, out, code, "prettier", "Prettier not installed: npm install --save-dev prettier")
		})
	})
	t.Run("0001.2.3 a missing plugin reports info quoting Prettier", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.3.1", "af82b5ec", func(_ req.A) {
			noPlugin := acceptance.Response{Stderr: "[error] Cannot find package 'prettier-plugin-packagejson' imported from /shared/prettier.config.cjs\n", Exit: 2}
			out, code, _ := prettierRun(t, prettierJS(nil), []acceptance.Response{noPlugin})
			acceptance.ShowsInfo(t, out, code, "prettier", "prettier exited 2 without a finding: [error] Cannot find package 'prettier-plugin-packagejson'")
		})
	})
	t.Run("0001.2.4 a file Prettier cannot parse fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.4.1", "f7956a3e", func(_ req.A) {
			unparsable := acceptance.Response{Stderr: "Checking formatting...\n[error] app.js: SyntaxError: Unexpected token (1:7)\n", Exit: 2}
			out, code, _ := prettierRun(t, prettierJS(nil), []acceptance.Response{unparsable})
			acceptance.Sees(t, out, code, 1, "[error] app.js: SyntaxError: Unexpected token (1:7)")
		})
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 Biome checks every tracked file in a repo with a biome config and the Biome package", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{"biome.jsonc": "{ \"formatter\": { \"enabled\": true } }\n"}
		var report strings.Builder
		for i := 1; i <= biomeMany; i++ {
			name := fmt.Sprintf("app%d.js", i)
			files[name] = prettierMessy
			fmt.Fprintf(&report, biomeReformat, name)
		}
		for _, major := range biomeMajors {
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
			req.Then(t, "0001.3.2.1.1", "fa88ed9f", func(_ req.A) {
				acceptance.Sees(t, out, code, 1, "✖ app1.js", "✖ app25.js")
			})
			req.Then(t, "0001.3.2.1.2", "7f5aef03", func(a req.A) {
				a.True(rp.Called("biome", "format", "README.md") && !strings.Contains(out, "Prettier not installed"),
					"Biome must check every tracked file in place of Prettier; calls %q:\n%s", rp.Calls("biome"), out)
			})
		}
	})
	t.Run("0001.3.2 a Biome before 1.8 fails naming its version", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{"biome.jsonc": "{ \"formatter\": { \"enabled\": true } }\n", "app.js": prettierMessy}
		files["package.json"] = fmt.Sprintf(biomePackage, "latest")
		rp := acceptance.StandIn(t, map[string][]acceptance.Response{"biome": {biomeVersion(biomeTooOld)}})
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), rp.Env, "prettier")
		req.Then(t, "0001.3.2.2.1", "16b50176", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "Biome "+biomeTooOld+" found; Biome 1.8 or later is required")
		})
		req.Then(t, "0001.3.2.2.2", "9edb5630", func(a req.A) {
			a.True(!rp.Called("biome", "--reporter=github"), "Biome %s cannot take --reporter=github, yet it ran: %q", biomeTooOld, rp.Calls("biome"))
		})
	})
	t.Run("0001.3.3 a write argument formats in place", func(t *testing.T) {
		t.Parallel()
		files := prettierJS(map[string]string{".prettierrc.json": "", "biome.jsonc": "{}\n", "package.json": fmt.Sprintf(biomePackage, "2.4.15")})
		written := acceptance.Response{Match: []string{"format " + fmt.Sprintf(biomeFlags, "none") + " --write app.js"}, Stdout: "Formatted 1 file in 2ms. Fixed 1 file.\n"}
		biomeRun := acceptance.StandIn(t, map[string][]acceptance.Response{"biome": {biomeVersion("2.4.15"), written}})
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), biomeRun.Env, "prettier", "--write", "--ignore-unknown", "app.js")
		req.Then(t, "0001.3.2.3.1", "e776c062", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
		req.Then(t, "0001.3.2.3.2", "aef997fd", func(a req.A) {
			for _, call := range biomeRun.Calls("biome") {
				a.True(!strings.Contains(call, "--ignore-unknown") && !strings.Contains(call, "README.md"),
					"Biome must get only its own flags and the given path; got %q", call)
			}
		})
	})
	t.Run("0001.3.4 a Prettier only flag fails by name", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.3.2.4.1", "b633cdce", func(_ req.A) {
			files := prettierJS(map[string]string{".prettierrc.json": "", "biome.jsonc": "{}\n", "package.json": fmt.Sprintf(biomePackage, "2.4.15")})
			biomeRun := acceptance.StandIn(t, map[string][]acceptance.Response{"biome": {biomeVersion("2.4.15")}})
			out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), biomeRun.Env, "prettier", "--write", "--print-width", "100", "app.js")
			acceptance.Sees(t, out, code, 1, "Biome formats this repo and has no counterpart for Prettier's --print-width")
		})
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
	req.Asserted(t)
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
