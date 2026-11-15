package requirements

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func Test0007_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"build-output-untracked rejects imports from dist", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{
				".gitignore": buildOutputIgnore,
				"src/app.test.ts": "import ok from './local.js';\n" +
					"import bad from '../dist/foo.js';\n" +
					"const z = require('pkg/dist/z.js');\n" +
					"const m = await import('../dist/mod.js');\n",
			})
			out, code := buildOutputRun(t, repo)
			sees(t, out, code, 1,
				`src/app.test.ts:2 imports build output: "../dist/foo.js"`,
				`src/app.test.ts:3 imports build output: "pkg/dist/z.js"`,
				`src/app.test.ts:4 imports build output: "../dist/mod.js"`)
		}},
		{"build-output-untracked rejects a dist folder not ignored", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			write(t, repo, map[string]string{"dist/app.js": "built\n"})
			out, code := buildOutputRun(t, repo)
			sees(t, out, code, 1, "Add dist/ to .gitignore (git check-ignore reports dist is not git-ignored)")
		}},
		{"build-output-untracked rejects tracked files in ignored dist", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"sub/dist/b.js": "built\n", "dist/a.js": "built\n"})
			commit(t, repo, map[string]string{".gitignore": buildOutputIgnore}, "chore: ignore dist")
			out, code := buildOutputRun(t, repo)
			sees(t, out, code, 1, "Remove tracked files: dist/a.js, sub/dist/b.js (git rm --cached)")
		}},
		{"build-output-untracked passes an ignored dist with local imports", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{
				".gitignore": buildOutputIgnore,
				"src/a.ts":   "import x from './b.js';\n",
				"src/b.mts":  "export const b = require('./c.cjs');\n",
			})
			out, code := buildOutputRun(t, repo)
			sees(t, out, code, 0, "All 1 checks passed · 3 files scanned")
		}},
		{"build-output-untracked passes a repo with no TypeScript", func(t *testing.T) {
			t.Parallel()
			out, code := buildOutputRun(t, example(t, "happyRepo", nil))
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, buildOutputCheck)
		}},
		{"build-output-untracked ignores untracked files importing dist", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{
				".gitignore": buildOutputIgnore,
				"src/a.ts":   "import x from './b.js';\n",
			})
			write(t, repo, map[string]string{"scratch.ts": "import x from '../dist/x.js';\n"})
			out, code := buildOutputRun(t, repo)
			sees(t, out, code, 0, "All 1 checks passed · 2 files scanned")
		}},
		{"build-output-untracked rejects tracked Maven and Gradle output", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{
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
			sees(t, out, code, 1, "Remove tracked files: classes/Old.class, lib[1]/.gradle/state, lib[1]/build/out.txt, src/app/target/app.jar, target/app.jar (git rm --cached)")
		}},
		{"build-output-untracked passes a build folder without Java", func(t *testing.T) {
			t.Parallel()
			out, code := buildOutputRun(t, example(t, "happyRepo", map[string]string{"build/notes.txt": "kept\n"}))
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, buildOutputCheck)
		}},
	})
}

func Test0007_2(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"no-plans-dir rejects a plan file in docs plans", func(t *testing.T) {
			t.Parallel()
			out, code := noPlansDir(t, map[string]string{"docs/plans/03-publish.md": "# Publish\n"})
			sees(t, out, code, 1, inPlansFolder("docs/plans/03-publish.md"))
		}},
		{"no-plans-dir rejects nested plan files", func(t *testing.T) {
			t.Parallel()
			out, code := noPlansDir(t, map[string]string{"docs/plans/archive/old.md": "# Old\n"})
			sees(t, out, code, 1, inPlansFolder("docs/plans/archive/old.md"))
		}},
		{"no-plans-dir passes maintained docs and shipped images", func(t *testing.T) {
			t.Parallel()
			files := map[string]string{
				"docs/notes.md":         "# Notes\n",
				"assets/screen.png":     "png",
				"docs/planner-guide.md": "# Planner guide\n\nA Plan lists its steps under a heading.\n",
				"go/plan_test.go":       "package plan\n\nconst body = \"## What needs to happen\\n\\n- [ ] Ship it.\\n\"\n",
			}
			repo := example(t, "happyRepo", files)
			out, code := fitness(t, repo, nil, "no-plans-dir")
			sees(t, out, code, 0, "All 1 checks passed")
			for name := range files {
				if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(name))); err != nil {
					t.Errorf("no-plans-dir must leave %s in place: %v", name, err)
				}
			}
		}},
		{"no-plans-dir rejects files in design artifact folders", func(t *testing.T) {
			t.Parallel()
			out, code := noPlansDir(t, map[string]string{
				"src/Plans/next.md":               "# Next\n",
				"design/mockups/home.png":         "png",
				"reports/validation/run.txt":      "ok",
				"review-artifacts/diff.txt":       "diff",
				"assets/icons/validation-ok.png":  "png",
				"docs/validation-checklist.md":    "# Validation\n",
				"tests/fixtures/review/input.txt": "in",
			})
			sees(t, out, code, 1,
				"design/mockups/home.png: a `mockups` folder holds temporary work; "+toPlanDesign,
				"reports/validation/run.txt: a `validation` folder holds temporary work; "+toPRComment,
				"review-artifacts/diff.txt: a `review-artifacts` folder holds temporary work; "+toPRComment,
				"src/Plans/next.md: a `plans` folder holds temporary work; "+toPlanIssue,
			)
			inOrder(t, out, "design/mockups/home.png:", "reports/validation/run.txt:", "review-artifacts/diff.txt:", "src/Plans/next.md:")
			for _, kept := range []string{"assets/icons/validation-ok.png", "docs/validation-checklist.md", "tests/fixtures/review/input.txt"} {
				if strings.Contains(out, kept+":") {
					t.Errorf("%s names no artifact folder, yet it was reported:\n%s", kept, out)
				}
			}
		}},
		{"no-plans-dir rejects plan and mockup file names", func(t *testing.T) {
			t.Parallel()
			out, code := noPlansDir(t, map[string]string{
				"docs/launch-plan.md":  "# Launch\n",
				"img/home-mockup.png":  "png",
				"pr-42-before.png":     "png",
				"docs/planner.md":      "# Planner\n",
				"img/print-before.png": "png",
			})
			sees(t, out, code, 1,
				"docs/launch-plan.md: a `*-plan.md` file is a plan; "+toPlanIssue,
				"img/home-mockup.png: a `*-mockup.*` file is a design mockup; "+toPlanDesign,
				"pr-42-before.png: a `pr-<number>-*` file belongs to one pull request; "+toPRComment,
			)
			for _, kept := range []string{"docs/planner.md", "img/print-before.png"} {
				if strings.Contains(out, kept+":") {
					t.Errorf("%s matches no artifact name, yet it was reported:\n%s", kept, out)
				}
			}
		}},
		{"no-plans-dir rejects a needs checklist outside github", func(t *testing.T) {
			t.Parallel()
			out, code := noPlansDir(t, map[string]string{
				"docs/ship.md":                   planChecklist,
				".github/ISSUE_TEMPLATE/plan.md": planChecklist,
				"docs/what-needs-to-happen.md":   "# Notes\n\n## What needs to happen\n\nNothing yet.\n",
			})
			sees(t, out, code, 1, "docs/ship.md: a `## What needs to happen` checklist is a plan; "+toPlanIssue)
			for _, kept := range []string{".github/ISSUE_TEMPLATE/plan.md", "docs/what-needs-to-happen.md"} {
				if strings.Contains(out, kept+":") {
					t.Errorf("%s holds no Plan checklist outside .github/, yet it was reported:\n%s", kept, out)
				}
			}
		}},
		{"no-plans-dir rejects staged and renamed plan files", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"docs/home.png": "png"})
			write(t, repo, map[string]string{"docs/plans/next.md": "# Next\n"})
			git(t, repo, "add", "docs/plans/next.md")
			mustDo(t, os.MkdirAll(filepath.Join(repo, "docs", "mockups"), 0o755))
			git(t, repo, "mv", "docs/home.png", "docs/mockups/home.png")
			out, code := fitness(t, repo, nil, "no-plans-dir")
			sees(t, out, code, 1,
				"docs/mockups/home.png: a `mockups` folder holds temporary work; "+toPlanDesign,
				inPlansFolder("docs/plans/next.md"),
			)
		}},
		{"no-plans-dir fails when the index is unreadable", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			mustDo(t, os.WriteFile(filepath.Join(repo, ".git", "index"), []byte("not an index"), 0o644))
			out, code := fitness(t, repo, nil, "no-plans-dir")
			sees(t, out, code, 1, "cannot list the repo's files")
		}},
	})
}

func Test0007_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"duplicate-code rejects a share past the limit", func(t *testing.T) {
			t.Parallel()
			out, code := dupRun(t, dupPair("src/a.js", "src/b.js", dupClone))
			sees(t, out, code, 1, "ERROR: duplicate code is 7.4% of scanned lines, over the 1.0% threshold", "src/b.js: 6 duplicated lines")
			if strings.Contains(read(out), "src/a.js:") {
				t.Errorf("the first copy holds no duplicated line, so only src/b.js is named:\n%s", out)
			}
		}},
		// The case commits a copied block in a script, which duplicate-code never
		// compares with a doc, then in a second doc.
		{"duplicate-code judges staged files against tracked files", func(t *testing.T) {
			t.Parallel()
			passedWith(t, dupCommit(t, "src/copy.js"), "duplicate-code", "1")
			out := dupCommit(t, "docs/copy.md")
			if !regexp.MustCompile(`│ failed\s+│ 4\s+│`).MatchString(row(out, "duplicate-code")) {
				t.Errorf("a staged doc copying a tracked doc must fail duplicate-code over the 4 docs, got row %q\n%s", row(out, "duplicate-code"), out)
			}
		}},
		{"duplicate-code compares every tracked text file", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{".gitignore": "conf/\n"})
			write(t, repo, dupPair("docs/a.md", "conf/b.json", dupClone))
			git(t, repo, "add", "-f", "docs", "conf")
			out, code := fitness(t, repo, nil, "duplicate-code")
			sees(t, out, code, 1, "ERROR: duplicate code is ")
		}},
		{"duplicate-code skips lockfiles and generated files", func(t *testing.T) {
			t.Parallel()
			files := dupPair("deps/a.lock", "deps/b.lock", dupClone)
			files[".gitattributes"] = "*.lock linguist-generated\n"
			for name, body := range dupPair("app/gradle.lockfile", "lib/gradle.lockfile", dupClone) {
				files[name] = body
			}
			for name, body := range dupPair("mvnw", "app/gradlew.bat", dupClone) {
				files[name] = body
			}
			out, code := dupRun(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		// The case tracks a binary pair and leaves a text pair untracked. Both
		// pairs share the block and their tails, so only skipping both passes.
		{"duplicate-code skips binary and untracked files", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", dupPair("blob/a.dat", "blob/b.dat", "\x00\n"+dupClone))
			write(t, repo, dupPair("src/a.js", "src/b.js", dupClone))
			out, code := fitness(t, repo, nil, "duplicate-code")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		// The case wraps the block in current markers in one pair and deprecated
		// jscpd ones in the other, so older repos keep passing.
		{"duplicate-code skips blocks inside ignore markers", func(t *testing.T) {
			t.Parallel()
			current := "// duplicate-code:ignore-start\n" + dupClone + "// duplicate-code:ignore-end\n"
			deprecated := "/* jscpd:ignore-start */\n" + dupBlock("j", 6, 10) + "/* jscpd:ignore-end */\n"
			out, code := dupRun(t, map[string]string{
				"src/a.js": current + dupBlock("fa", 24, 3),
				"src/b.js": current + dupBlock("fb", 24, 3),
				"src/c.js": deprecated + dupBlock("fc", 24, 3),
				"src/d.js": deprecated + dupBlock("fd", 24, 3),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		// The case adds a template, matching legends and a copied untitled
		// diagram to the two docs of the share past the limit case. A
		// committed third doc scopes the run to Markdown, and the exact share
		// proves the template is not scanned.
		{"duplicate-code excuses docs repeating a template", func(t *testing.T) {
			t.Parallel()
			files := dupPair("docs/guide.md", "docs/notes.md", dupClone)
			files["docs/templates/guide.md"] = dupClone + dupBlock("ft", 24, 3)
			files["docs/architecture/02-containers.md"] = dupLevelDoc("ca")
			files["docs/architecture/03-components.md"] = dupLevelDoc("cb")
			files["docs/flow-a.md"] = "# Flow\n\n" + dupCopiedDiagram + dupBlock("la", 24, 3)
			files["docs/flow-b.md"] = "# Flow\n\n" + dupLevelDiagram("lb") + dupCopiedDiagram + dupBlock("lb", 24, 3)
			repo, env := hooked(t, files)
			write(t, repo, map[string]string{"docs/new.md": "# New\n\n" + dupBlock("fn", 24, 3)})
			user(t, repo, env, "add", "docs/new.md")
			out, code := user(t, repo, env, "commit", "-m", "docs(new): add a doc")
			if code == 0 {
				t.Errorf("the hook must reject the commit:\n%s", out)
			}
			dupNames(t, out, []string{"ERROR: duplicate code is 5.6% of scanned lines", "docs/notes.md: 6 duplicated lines", "docs/flow-b.md: 11 duplicated lines"},
				[]string{"docs/architecture/02-containers.md:", "docs/architecture/03-components.md:", "docs/flow-a.md:", "docs/guide.md:", "docs/templates/guide.md:"})
		}},
		// The case puts a 3-line block in one pair and an 18-token block in
		// another, each under one minimum. Two docs share a 5-line run only because
		// it ends in a 120-dash separator row.
		{"duplicate-code passes short shared blocks", func(t *testing.T) {
			t.Parallel()
			files := dupPair("src/a.js", "src/b.js", dupBlock("w", 3, 20))
			files["src/c.js"] = dupBlock("s", 6, 3) + dupBlock("fc", 24, 3)
			files["src/d.js"] = dupBlock("s", 6, 3) + dupBlock("fd", 24, 3)
			files["docs/a.md"] = "# A\n\n" + dupSeparatorTable + "a | b\n" + dupBlock("da", 24, 3)
			files["docs/b.md"] = "# B\n\n" + dupSeparatorTable + "c | d\n" + dupBlock("db", 24, 3)
			out, code := dupRun(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
	})
}

func Test0007_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"repeated-string-literals rejects three uses of a value", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{
				"src/a.ts": rslLines("active", 2),
				"src/b.ts": rslLines("active", 1),
			})
			sees(t, out, code, 1, `"active" appears 3 times (src/a.ts:1, src/a.ts:2, src/b.ts:1) — extract a shared constant`)
		}},
		{"repeated-string-literals passes a value used twice", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.ts": rslLines("active", 2)})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"repeated-string-literals skips comments regex and templates", func(t *testing.T) {
			t.Parallel()
			src := "import x from 'active';\nconst y = require('active');\nexport { z } from 'active';\n" +
				"// 'active'\n/* 'active' */\nconst re = /'active'/g;\nconst s = `'active'`;\n"
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.ts": src})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"repeated-string-literals skips short and idiomatic values", func(t *testing.T) {
			t.Parallel()
			src := rslLines("ok", 3) + rslLines("utf8", 3) + rslLines("inherit", 3) + rslLines("object", 3) + rslLines("use strict", 3)
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.js": src})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"repeated-string-literals lists five places and a count", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.ts": rslLines("active", 7)})
			sees(t, out, code, 1, `"active" appears 7 times (src/a.ts:1, src/a.ts:2, src/a.ts:3, src/a.ts:4, src/a.ts:5, +2 more)`)
		}},
		{"repeated-string-literals reports the most repeated first", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.ts": rslLines("twice", 3) + rslLines("thrice", 4)})
			sees(t, out, code, 1, `"thrice" appears 4 times`, `"twice" appears 3 times`)
			if strings.Index(out, `"thrice"`) > strings.Index(out, `"twice"`) {
				t.Errorf("want thrice (4) reported before twice (3):\n%s", out)
			}
		}},
		{"repeated-string-literals counts test and bench files", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{
				"src/a.test.ts":  rslLines("active", 2),
				"src/b.bench.ts": rslLines("active", 1),
			})
			sees(t, out, code, 1, `"active" appears 3 times (src/a.test.ts:1, src/a.test.ts:2, src/b.bench.ts:1)`)
		}},
		{"repeated-string-literals passes an allowed value", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{
				".fitnessrc.json": `{"repeatedStringLiterals": {"allow": ["active"]}}` + "\n",
				"src/a.ts":        rslLines("active", 3),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
	})
}

func Test0007_5(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"dependency-currency passes dependencies at the latest version", func(t *testing.T) {
			t.Parallel()
			out, code := dependencyCurrency(t, dcRepo(dcRegistry(t), `{"dependencies": {"left-pad": "^1.3.0"}}`, "", dcLeftPad, "1.3.0"))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"dependency-currency rejects each kind behind latest", func(t *testing.T) {
			t.Parallel()
			files := dcRepo(dcRegistry(t), `{"dependencies": {"left-pad": "^1.0.0"}, "devDependencies": {"is-number": "^6.0.0"}, "peerDependencies": {"@types/left-pad": "^1.0.0"}}`,
				"", dcLeftPad, "1.0.0", "", dcIsNumber, "6.0.0", "", dcTypes, "1.0.0")
			out, code := dependencyCurrency(t, files)
			sees(t, out, code, 1, dcLead, "left-pad: 1.0.0 → 1.3.0", "is-number: 6.0.0 → 7.0.0", "@types/left-pad: 1.0.0 → 1.2.0")
		}},
		{"dependency-currency rejects a dependency that is missing", func(t *testing.T) {
			t.Parallel()
			out, code := dependencyCurrency(t, dcRepo(dcRegistry(t), `{"dependencies": {"left-pad": "^1.3.0"}}`))
			sees(t, out, code, 1, dcLead, "left-pad: missing → 1.3.0")
		}},
		{"dependency-currency skips internal packages", func(t *testing.T) {
			t.Parallel()
			out, code := dependencyCurrency(t, dcRepo(dcRegistry(t), `{"dependencies": {"`+dcInternal+`": "*"}}`))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"dependency-currency rejects stale workspace dependencies", func(t *testing.T) {
			t.Parallel()
			for _, workspaces := range []string{`["packages/*"]`, `{"packages": ["packages/a"]}`} {
				files := dcRepo(dcRegistry(t), `{"workspaces": `+workspaces+`}`, "", dcLeftPad, "1.0.0")
				files["packages/a/package.json"] = `{"dependencies": {"left-pad": "^1.0.0"}}` + "\n"
				out, code := dependencyCurrency(t, files)
				sees(t, out, code, 1, dcLead, "left-pad: 1.0.0 → 1.3.0")
			}
		}},
		{"dependency-currency lists shared installs once", func(t *testing.T) {
			t.Parallel()
			files := dcRepo(dcRegistry(t), `{"workspaces": ["packages/*"]}`, "", dcLeftPad, "1.0.0", "packages/a/", dcLeftPad, "1.1.0")
			for _, ws := range []string{"a", "b", "c"} {
				files["packages/"+ws+"/package.json"] = `{"devDependencies": {"left-pad": "^1.0.0"}}` + "\n"
			}
			out, code := dependencyCurrency(t, files)
			sees(t, out, code, 1, "left-pad: 1.1.0 → 1.3.0", "left-pad: 1.0.0 → 1.3.0")
			if n := strings.Count(read(out), "left-pad: 1.0.0 → 1.3.0"); n != 1 {
				t.Errorf("the shared install is listed %d times, want once:\n%s", n, out)
			}
		}},
		{"dependency-currency reports info for an unreachable registry", func(t *testing.T) {
			t.Parallel()
			for _, broken := range dcBrokenRegistries(t) {
				start := time.Now()
				out, code := dcAgainst(t, broken)
				showsInfo(t, out, code, "dependency-currency", cmp.Or(broken.want, "1 of 1 packages could not be queried from "+broken.url+", so their currency is unknown"))
				if broken.hits != nil && broken.hits.Load() == 0 {
					t.Errorf("%s: the check never asked the registry", broken.name)
				}
				if took := time.Since(start); took > dcSilentWithin {
					t.Errorf("%s: the run took %v, past the %v a silent registry may cost", broken.name, took.Round(time.Second), dcSilentWithin)
				}
			}
		}},
		{"dependency-currency passes a repo without package.json", func(t *testing.T) {
			t.Parallel()
			out, code := dependencyCurrency(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "dependency-currency")
		}},
	})
}

func Test0007_6(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"gitignore-why passes patterns that each have a comment", func(t *testing.T) {
			t.Parallel()
			out, code := gitignoreWhy(t, "# dependencies\nnode_modules/\n\n# build output\ndist\n\n# keep this tracked\n!dist/keep.js\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"gitignore-why rejects a pattern with no comment", func(t *testing.T) {
			t.Parallel()
			out, code := gitignoreWhy(t, "node_modules/\n")
			sees(t, out, code, 1, `.gitignore:1: pattern "node_modules/" `+gitignoreWhyNoComment)
		}},
		{"gitignore-why rejects a second pattern under one comment", func(t *testing.T) {
			t.Parallel()
			out, code := gitignoreWhy(t, "# deps\nnode_modules/\ndist\n")
			sees(t, out, code, 1, `.gitignore:3: pattern "dist" `+gitignoreWhyNoComment)
		}},
		{"gitignore-why rejects a blank line before the pattern", func(t *testing.T) {
			t.Parallel()
			out, code := gitignoreWhy(t, "# deps\n\nnode_modules/\n")
			sees(t, out, code, 1, `.gitignore:3: pattern "node_modules/" `+gitignoreWhyNoComment)
		}},
		{"gitignore-why rejects an empty comment", func(t *testing.T) {
			t.Parallel()
			out, code := gitignoreWhy(t, "#\nnode_modules/\n")
			sees(t, out, code, 1, `.gitignore:2: pattern "node_modules/" `+gitignoreWhyNoComment)
		}},
		{"gitignore-why passes a repo with no gitignore", func(t *testing.T) {
			t.Parallel()
			out, code := gitignoreWhy(t, "")
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "gitignore-why")
		}},
	})
}

func Test0007_7(t *testing.T) {
	runCases(t, []edgeCase{
		{"workflow-ref-exists passes paths present at their refs", func(t *testing.T) {
			t.Parallel()
			out, code, _ := refCheck(t, refWorkflow(
				"may-journal/fitness-runner/.github/workflows/ci-reusable.yml@v1",
				"may-journal/fitness-runner@"+refSha,
			), []response{
				found("repos/may-journal/fitness-runner/git/ref/tags/v1", `{"ref":"refs/tags/v1"}`),
				found("repos/may-journal/fitness-runner/contents/.github/workflows/ci-reusable.yml?ref=v1", `{"name":"ci-reusable.yml"}`),
				found("repos/may-journal/fitness-runner/contents?ref="+refSha, `[{"name":"README.md"},{"name":"action.yml"}]`),
			})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "workflow-ref-exists", "1")
		}},
		{"workflow-ref-exists rejects a workflow missing at its tag", func(t *testing.T) {
			t.Parallel()
			out, code, _ := refCheck(t, refWorkflow(
				"may-journal/fitness-runner/.github/workflows/comment-check-reusable.yml@v1",
			), []response{
				found("repos/may-journal/fitness-runner/git/ref/tags/v1", `{"ref":"refs/tags/v1"}`),
				notFound("repos/may-journal/fitness-runner/contents/.github/workflows/comment-check-reusable.yml?ref=v1"),
				found("repos/may-journal/fitness-runner", `{"full_name":"may-journal/fitness-runner"}`),
			})
			sees(t, out, code, 1, ".github/workflows/ci.yml: line 5: may-journal/fitness-runner/.github/workflows/comment-check-reusable.yml@v1: "+
				".github/workflows/comment-check-reusable.yml does not exist at v1")
		}},
		{"workflow-ref-exists rejects a missing action ref or sha", func(t *testing.T) {
			t.Parallel()
			badSha := strings.Repeat("0", 40)
			out, code, _ := refCheck(t, refWorkflow("may-journal/fitness-runner@v9", "may-journal/fitness-runner@"+badSha), []response{
				notFound("repos/may-journal/fitness-runner/git/ref/tags/v9"),
				notFound("repos/may-journal/fitness-runner/git/ref/heads/v9"),
				{Match: []string{"contents?ref=" + badSha}, Stderr: "gh: No commit found for the ref " + badSha + " (HTTP 404)\n", Exit: 1},
				found("repos/may-journal/fitness-runner", `{"full_name":"may-journal/fitness-runner"}`),
			})
			sees(t, out, code, 1, ".github/workflows/ci.yml: line 5: may-journal/fitness-runner@v9: v9 is not a tag, a branch, or a full commit sha",
				".github/workflows/ci.yml: line 7: may-journal/fitness-runner@"+badSha+": no commit "+badSha+" exists")
		}},
		{"workflow-ref-exists skips local and other owner calls", func(t *testing.T) {
			t.Parallel()
			out, code, rp := refCheck(t, refWorkflow(
				"./.github/workflows/local.yml",
				"actions/checkout@v4",
				"may-journal/fitness-runner/.github/workflows/ci-reusable.yml@${{ inputs.ref }}",
				"may-journal/fitness-runner/.github/workflows/ci-reusable.yml@main",
			)+"  script:\n    runs-on: ubuntu-latest\n    steps:\n      - run: |\n          echo 'uses: may-journal/fitness-runner/nope.yml@v1' > caller.yml\n", []response{
				notFound("repos/may-journal/fitness-runner/git/ref/tags/main"),
				found("repos/may-journal/fitness-runner/git/ref/heads/main", `{"ref":"refs/heads/main"}`),
			})
			sees(t, out, code, 0, "All 1 checks passed")
			if rp.called("gh", "contents") || rp.called("gh", "actions/checkout") {
				t.Errorf("an ignored line read GitHub: %q", rp.calls("gh"))
			}
		}},
		{"workflow-ref-exists notes a repo the token cannot read", func(t *testing.T) {
			t.Parallel()
			out, code, _ := refCheck(t, refWorkflow("may-journal/hub/.github/workflows/sync.yml@v1"), []response{
				notFound("repos/may-journal/hub"),
			})
			sees(t, out, code, 0, ".github/workflows/ci.yml: line 5: may-journal/hub/.github/workflows/sync.yml@v1: skipped, may-journal/hub cannot be read with this token")
		}},
		{"workflow-ref-exists passes a repo with no workflows", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", nil), nil, "workflow-ref-exists")
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "workflow-ref-exists")
		}},
		{"workflow-ref-exists notes calls it could not verify", func(t *testing.T) {
			t.Parallel()
			out, code, _ := refCheck(t, refWorkflow(
				"may-journal/fitness-runner/.github/workflows/nope.yml@v1",
				"may-journal/hub/.github/workflows/sync.yml@v1",
				"may-journal/journey/.github/workflows/build.yml@v1",
			), []response{
				found("repos/may-journal/fitness-runner/git/ref/tags/v1", `{"ref":"refs/tags/v1"}`),
				notFound("repos/may-journal/fitness-runner/contents/.github/workflows/nope.yml?ref=v1"),
				found("repos/may-journal/fitness-runner", `{"full_name":"may-journal/fitness-runner"}`),
				{Match: []string{"repos/may-journal/hub/"}, Stderr: "gh: Resource not accessible by integration (HTTP 403)\n", Exit: 1},
				{Match: []string{"repos/may-journal/journey/"}, Stderr: "error connecting to api.github.com\n", Exit: 1},
			})
			sees(t, out, code, 1, ".github/workflows/ci.yml: line 5: may-journal/fitness-runner/.github/workflows/nope.yml@v1: .github/workflows/nope.yml does not exist at v1",
				".github/workflows/ci.yml: line 7: may-journal/hub/.github/workflows/sync.yml@v1: skipped, GitHub answered gh: Resource not accessible by integration (HTTP 403)",
				"workflow refs from .github/workflows/ci.yml line 9 on were not verified: gh api repos/may-journal/journey/git/ref/tags/v1: error connecting to api.github.com")
		}},
		{"workflow-ref-exists reads pinned calls in parallel", func(t *testing.T) {
			// Not parallel: a wall-clock limit beside other tests would time their load.
			var uses []string
			for i := range slowRefs {
				uses = append(uses, fmt.Sprintf("may-journal/fitness-runner/.github/workflows/w%d.yml@%s", i, refSha))
			}
			slow := found("repos/may-journal/fitness-runner/contents/.github/workflows/", `{"name":"w.yml"}`)
			slow.Wait = 1
			out, code, rp := refCheck(t, refWorkflow(uses...), []response{slow})
			sees(t, out, code, 0, "All 1 checks passed")
			if reads := len(rp.calls("gh")); reads != slowRefs {
				t.Errorf("workflow-ref-exists read GitHub %d times, want %d", reads, slowRefs)
			}
			took := rowDuration.FindStringSubmatch(row(out, "workflow-ref-exists"))
			if took == nil {
				t.Fatalf("workflow-ref-exists shows no time in %q", row(out, "workflow-ref-exists"))
			}
			if d, _ := time.ParseDuration(took[1] + took[2]); d >= slowRefs*time.Second/2 {
				t.Errorf("workflow-ref-exists took %v on %d one-second reads, want under %v", d, slowRefs, slowRefs*time.Second/2)
			}
		}},
	})
}

func Test0007_8(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"workflow-reporting passes a final always summary step", func(t *testing.T) {
			t.Parallel()
			out, code := reportingCheck(t, map[string]string{".github/workflows/ci.yml": reportingJob(reportStep)})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "workflow-reporting", "1")
		}},
		{"workflow-reporting rejects a job with no reporting step", func(t *testing.T) {
			t.Parallel()
			out, code := reportingCheck(t, map[string]string{".github/workflows/build.yaml": reportingJob("")})
			sees(t, out, code, 1, ".github/workflows/build.yaml:4: job build has no step that runs with if: always() and writes $GITHUB_STEP_SUMMARY",
				"add a last step with `if: ${{ always() }}` whose run writes the job's status to \"$GITHUB_STEP_SUMMARY\"")
		}},
		{"workflow-reporting rejects a summary step that skips failure", func(t *testing.T) {
			t.Parallel()
			onSuccess := strings.Replace(reportStep, "        if: ${{ always() }}\n", "", 1)
			out, code := reportingCheck(t, map[string]string{".github/workflows/ci.yml": reportingJob(onSuccess)})
			sees(t, out, code, 1, ".github/workflows/ci.yml:4: job build writes $GITHUB_STEP_SUMMARY only when earlier steps pass, at line 9",
				"add `if: ${{ always() }}` to the step at line 9")
		}},
		{"workflow-reporting passes a job calling a reusable workflow", func(t *testing.T) {
			t.Parallel()
			out, code := reportingCheck(t, map[string]string{".github/workflows/ci.yml": "name: CI\non: push\njobs:\n" +
				"  fitness:\n    uses: may-journal/fitness-runner/.github/workflows/ci-reusable.yml@v1\n    secrets: inherit\n"})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "workflow-reporting", "1")
		}},
		{"workflow-reporting names a matrix job once", func(t *testing.T) {
			t.Parallel()
			matrix := "name: CI\non: push\njobs:\n  test:\n    strategy:\n      matrix:\n        os: [ubuntu-latest, macos-latest, windows-latest]\n" +
				"    runs-on: ${{ matrix.os }}\n    steps:\n      - run: make test\n"
			out, code := reportingCheck(t, map[string]string{".github/workflows/ci.yml": matrix})
			sees(t, out, code, 1, ".github/workflows/ci.yml:4: job test has no step")
			if n := strings.Count(out, "job test has no step"); n != 1 {
				t.Errorf("a matrix job must be named once, got %d times:\n%s", n, out)
			}
		}},
		{"workflow-reporting rejects invalid YAML with its position", func(t *testing.T) {
			t.Parallel()
			files := runsOnWorkflows(rejectedRunsOn)
			files[".github/workflows/ci.yml"] = "name: CI\non: push\njobs:\n  build:\n\truns-on: ubuntu-latest\n"
			out, code := reportingCheck(t, files)
			const bad = ": this workflow is not valid YAML: "
			sees(t, out, code, 1, ".github/workflows/ci.yml:5"+bad+"a tab indents this line",
				".github/workflows/fork.yml:5:116"+bad+"a `: ` inside a plain value starts a mapping",
				".github/workflows/reserved.yml:5:14"+bad+"a plain value cannot start with '@'",
				".github/workflows/key-line.yml:5:14"+bad+"'-' opens a list item on its key's line",
				"move it onto its own line, indented under the key",
				".github/workflows/continued.yml:6:13"+bad+"a `: ` inside a plain value",
				".github/workflows/colon-end.yml:5:27"+bad+"a plain value ends in `:`",
				".github/workflows/flow-trailing.yml:5:27"+bad+"text follows the closing bracket",
				".github/workflows/wide.yml:5:18"+bad+"a `: ` inside a plain value")
			out, code = reportingCheck(t, runsOnWorkflows(acceptedRunsOn))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"workflow-reporting passes a repo with no workflows", func(t *testing.T) {
			t.Parallel()
			out, code := reportingCheck(t, nil)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "workflow-reporting")
		}},
		{"workflow-reporting passes the shipped fitness-runner workflows", func(t *testing.T) {
			t.Parallel()
			own, err := filepath.Glob(filepath.Join("..", "..", ".github", "workflows", "*.y*ml"))
			mustDo(t, err)
			files := map[string]string{}
			for _, f := range own {
				body, err := os.ReadFile(f)
				mustDo(t, err)
				files[".github/workflows/"+filepath.Base(f)] = string(body)
			}
			out, code := reportingCheck(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
	})
}

// buildOutputCheck is the check these tests run.
const buildOutputCheck = "build-output-untracked"

// buildOutputIgnore is a .gitignore that ignores dist.
const buildOutputIgnore = "dist\n"

// buildOutputRun runs `fitness-install -- build-output-untracked` in repo.
func buildOutputRun(t *testing.T, repo string) (string, int) {
	t.Helper()
	return fitness(t, repo, nil, buildOutputCheck)
}

// The packages these tests declare, each with its saved latest version.
const (
	dcLeftPad  = "left-pad"        // the saved latest version is 1.3.0
	dcIsNumber = "is-number"       // the saved latest version is 7.0.0
	dcTypes    = "@types/left-pad" // the saved latest version is 1.2.0
	dcInternal = "@mayjournal/fitness-shared"
)

// dcLead is the line that opens every failure.
const dcLead = "Dependencies behind their latest published version — update them (npm install <pkg>@latest) or pin intentionally:"

// dcUnreachable is a registry address nothing listens on.
const dcUnreachable = "http://127.0.0.1:1"

// dcDoc is a package's abbreviated registry document, as npm serves it for
// Accept: application/vnd.npm.install-v1+json. The last version is latest.
func dcDoc(name string, versions ...string) string {
	all := map[string]any{}
	for _, v := range versions {
		all[v] = map[string]any{"name": name, "version": v, "dist": map[string]string{"shasum": strings.Repeat("0", 40)}}
	}
	doc, _ := json.Marshal(map[string]any{
		"name":      name,
		"modified":  "2024-01-01T00:00:00.000Z",
		"dist-tags": map[string]string{"latest": versions[len(versions)-1]},
		"versions":  all,
	})
	return string(doc)
}

// dcRegistry serves the saved documents for every package these tests
// declare. The internal package has a newer published version, so it would
// fail if the check did not skip it.
func dcRegistry(t *testing.T) string {
	t.Helper()
	return registry(t, map[string]string{
		dcLeftPad:  dcDoc(dcLeftPad, "1.0.0", "1.1.0", "1.3.0"),
		dcIsNumber: dcDoc(dcIsNumber, "6.0.0", "7.0.0"),
		dcTypes:    dcDoc(dcTypes, "1.0.0", "1.2.0"),
		dcInternal: dcDoc(dcInternal, "9.9.9"),
	})
}

// dcAnswer is a registry that answers every request with status and body,
// and counts the requests it received.
func dcAnswer(t *testing.T, status int, body string) (string, *atomic.Int32) {
	t.Helper()
	hits := new(atomic.Int32)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv.URL, hits
}

// dcNpmrc is an .npmrc pointing npm at registry.
func dcNpmrc(registry string) string {
	return "registry=" + registry + "\n"
}

// dependencyCurrency runs `fitness-install -- dependency-currency` on happyRepo
// with files written over it, ignoring node_modules as a real repo does. A local
// listener stands in for GitHub's API so the probe needs no network, and env may
// name another.
func dependencyCurrency(t *testing.T, files map[string]string, env ...string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", with(map[string]string{".gitignore": "node_modules\n"}, files))
	return fitness(t, repo, append([]string{"HOME=" + t.TempDir(), "GITHUB_API_URL=" + silentListener(t)}, env...), "dependency-currency")
}

// dcInstalled is the node_modules manifest of name installed at version,
// under dir ("" for the repo root).
func dcInstalled(dir, name, version string) (string, string) {
	return dir + "node_modules/" + name + "/package.json", `{"name": "` + name + `", "version": "` + version + `"}` + "\n"
}

// dcRepo is a package.json declaring manifest, an .npmrc pointing at
// registry, plus each install given as dir, name, version triples.
func dcRepo(registry, manifest string, installs ...string) map[string]string {
	files := map[string]string{"package.json": manifest + "\n", ".npmrc": dcNpmrc(registry)}
	for i := 0; i+2 < len(installs); i += 3 {
		path, body := dcInstalled(installs[i], installs[i+1], installs[i+2])
		files[path] = body
	}
	return files
}

// dcBroken is one registry that cannot tell the latest version, named in
// the repo or the home .npmrc. want is the info the run shows, when not
// the count of packages it could not query.
type dcBroken struct {
	name, url string
	want      string
	hits      *atomic.Int32
	home      bool
	env       []string
}

// dcBrokenRegistries are unreachable, erroring, garbage-answering and silent
// registries, set in the repo or home .npmrc or behind a proxy.
func dcBrokenRegistries(t *testing.T) []dcBroken {
	failing, failHits := dcAnswer(t, http.StatusInternalServerError, `{"error":"Internal Server Error"}`)
	garbage, garbageHits := dcAnswer(t, http.StatusOK, "<html>proxy login</html>")
	homeGarbage, homeHits := dcAnswer(t, http.StatusOK, `{"dist-tags": `)
	return []dcBroken{
		{name: "unreachable", url: dcUnreachable},
		{name: "server error", url: failing, hits: failHits},
		{name: "garbage", url: garbage, hits: garbageHits},
		{name: "silent", url: silentListener(t)},
		{name: "home garbage", url: homeGarbage, hits: homeHits, home: true},
		{name: "silent proxy", url: "https://registry.example.invalid", want: "did not answer within 2s through proxy", env: []string{"HTTPS_PROXY=" + silentListener(t), "NO_PROXY=", "GITHUB_API_URL=https://api.github.com"}},
	}
}

// dcAgainst runs dependency-currency on a repo with left-pad behind latest,
// asking the broken registry from the repo or home .npmrc.
func dcAgainst(t *testing.T, broken dcBroken) (string, int) {
	t.Helper()
	files := dcRepo(broken.url, `{"dependencies": {"left-pad": "^1.0.0"}}`, "", dcLeftPad, "1.0.0")
	home := t.TempDir()
	if broken.home {
		delete(files, ".npmrc")
		write(t, home, map[string]string{".npmrc": dcNpmrc(broken.url)})
	}
	return dependencyCurrency(t, files, append([]string{"HOME=" + home}, broken.env...)...)
}

// dcSilentWithin bounds a run whose registry never answers: the ten-second
// request timeout plus install, well under the 30-second budget.
const dcSilentWithin = 20 * time.Second

// gitignoreWhy runs `fitness-install -- gitignore-why` on happyRepo with
// the given .gitignore; an empty one means the repo has none.
func gitignoreWhy(t *testing.T, gitignore string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", map[string]string{".gitignore": gitignore})
	return fitness(t, repo, nil, "gitignore-why")
}

// gitignoreWhyNoComment is the tail of every gitignore-why failure.
const gitignoreWhyNoComment = "has no explanatory # comment on the line above"

// cspell:ignore mvnw gradlew

// dupBlock builds lines of distinct words, so a shared prefix yields an
// exact copy and different prefixes never match. Six lines of ten words
// plus ";" is 66 tokens, past the 50-token minimum.
func dupBlock(prefix string, lines, perLine int) string {
	var b strings.Builder
	for i := 0; i < lines; i++ {
		for j := 0; j < perLine; j++ {
			fmt.Fprintf(&b, "%s_%d_%d ", prefix, i, j)
		}
		b.WriteString(";\n")
	}
	return b.String()
}

// dupClone is a copied block long enough to count as a clone.
var dupClone = dupBlock("c", 6, 10)

// dupPair puts block at the top of two files that differ after it.
func dupPair(a, b, block string) map[string]string {
	return map[string]string{
		a: block + dupBlock("fa", 24, 3),
		b: block + dupBlock("fb", 24, 3),
	}
}

// dupRun runs `fitness-install -- duplicate-code` on happyRepo with files over it.
func dupRun(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "duplicate-code")
}

// dupSeparatorTable is two lines of prose and a table's header and separator
// rows. Its 120 dashes carry it past 50 tokens and its last line ends a
// 5-line run.
var dupSeparatorTable = "Both docs list names.\nEach row says why.\n\n| Name | Why |\n| " +
	strings.Repeat("-", 60) + " | " + strings.Repeat("-", 60) + " |\n"

// dupCommit tracks a doc holding dupClone in a hooked repo, stages file
// holding it too, and returns what the commit printed.
func dupCommit(t *testing.T, file string) string {
	t.Helper()
	repo, env := hooked(t, map[string]string{"docs/guide.md": dupClone + dupBlock("fg", 24, 3)})
	write(t, repo, map[string]string{file: dupClone + dupBlock("fc", 24, 3)})
	user(t, repo, env, "add", file)
	out, _ := user(t, repo, env, "commit", "-m", "docs(guide): copy the block")
	return out
}

// dupCopiedDiagram is an untitled diagram with no callouts, then a 48-token
// line. The lexer reads the fenced block as one token, so only a copy whose
// body stays reaches both minimums.
var dupCopiedDiagram = "```mermaid\nflowchart TB\n" + dupBlock("m", 6, 3) + "```\n\n" + dupBlock("t", 1, 47)

// dupNames asserts out holds each of wants and none of the clean paths, which
// hold only lines a doc must repeat.
func dupNames(t *testing.T, out string, wants, clean []string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(read(out), want) {
			t.Errorf("the commit output is missing %q:\n%s", want, out)
		}
	}
	for _, path := range clean {
		if regexp.MustCompile(regexp.QuoteMeta(path) + ` \d+ duplicated lines`).MatchString(read(out)) {
			t.Errorf("%s holds no duplicated line:\n%s", path, out)
		}
	}
}

// dupLevelDoc is one C4 level doc: a numbered diagram, the shared legend
// and table, then prose of its own.
func dupLevelDoc(prefix string) string {
	return "# Level\n\n" + dupLevelDiagram(prefix) + dupLegend + dupTable(prefix) + "\n" + dupBlock(prefix, 24, 3)
}

// dupTable is a level's callout rows, which name the level.
func dupTable(prefix string) string {
	return "| 1 | The " + prefix + " user | Starts " + prefix + " |\n| 2 | The " + prefix + " app | Serves " + prefix + " |\n| A | Opens " + prefix + " | Daily |\n"
}

// dupLevelDiagram is a numbered C4 diagram whose only repeated lines are
// its style statements.
func dupLevelDiagram(prefix string) string {
	return "```mermaid\nC4Container\n    Person(" + prefix + ", \"1 " + prefix + " user\")\n    Container(" + prefix + "2, \"2 " + prefix + " app\", \"Go\", \"Runs " + prefix + "\")\n" +
		"    Rel(" + prefix + ", " + prefix + "2, \"A\")\n" +
		"    UpdateElementStyle(" + prefix + ", $bgColor=\"#334155\", $borderColor=\"#1e293b\", $fontColor=\"#ffffff\")\n" +
		"    UpdateElementStyle(" + prefix + "2, $bgColor=\"#0f766e\", $borderColor=\"#115e59\", $fontColor=\"#ffffff\")\n```\n\n"
}

// dupLegend is the legend diagram every C4 level repeats, with its table.
const dupLegend = "```mermaid\nC4Container\n    title Legend\n\n    Person(p, \"Person\", \"Human or agent actor\")\n    System_Boundary(b, \"System boundary\") {\n        Container(c, \"Container\", \"Tech\", \"Runnable or deployable unit\")\n    }\n    System_Ext(e, \"External\", \"Outside our boundary\")\n\n    UpdateElementStyle(p, $bgColor=\"#334155\", $borderColor=\"#1e293b\", $fontColor=\"#ffffff\")\n    UpdateElementStyle(c, $bgColor=\"#0f766e\", $borderColor=\"#115e59\", $fontColor=\"#ffffff\")\n    UpdateElementStyle(e, $bgColor=\"#64748b\", $borderColor=\"#475569\", $fontColor=\"#ffffff\")\n```\n\nNumbers on nodes and letters on arrows match the callout table.\n\n| # | Description | Why |\n| --- | --- | --- |\n"

// noPlansDir runs `fitness-install -- no-plans-dir` on happyRepo with files
// written over it.
func noPlansDir(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "no-plans-dir")
}

// Where the check sends each kind of temporary work.
const (
	toPlanIssue  = "open a `Plan`-labeled GitHub Issue and put the plan there, not a file in the repo"
	toPlanDesign = "post it as a comment on the Plan Issue that it designs, not a file in the repo"
	toPRComment  = "post it as a comment on the pull request it reviews, not a file in the repo"
)

// inPlansFolder is the error line for a file in a `plans` folder.
func inPlansFolder(path string) string {
	return path + ": a `plans` folder holds temporary work; " + toPlanIssue
}

// planChecklist is a Plan's checklist as a Markdown file would hold it.
const planChecklist = "# Ship it\n\n## What needs to happen\n\n- [ ] Write the tests.\n"

// inOrder asserts each want appears in out after the one before it.
func inOrder(t *testing.T, out string, want ...string) {
	t.Helper()
	text, at := read(out), 0
	for _, w := range want {
		i := strings.Index(text[at:], w)
		if i < 0 {
			t.Errorf("output is missing %q after position %d, or lists it out of order:\n%s", w, at, out)
			return
		}
		at += i + len(w)
	}
}

// repeatedStringLiterals runs `fitness-install -- repeated-string-literals`
// on happyRepo with files written over it.
func repeatedStringLiterals(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "repeated-string-literals")
}

// rslLines returns n source lines that each use the literal value.
func rslLines(value string, n int) string {
	return strings.Repeat("use('"+value+"');\n", n)
}

// refSha is a full commit sha a workflow pins an action to.
const refSha = "6b02ebef8736d2508d47e085806b5b3cd19efa6c"

// refWorkflow is a workflow file whose jobs call each `uses:` value, one
// job per value, so value i sits on line 4+2i.
func refWorkflow(uses ...string) string {
	var b strings.Builder
	b.WriteString("name: CI\non: push\njobs:\n")
	for i, u := range uses {
		b.WriteString("  job" + string(rune('a'+i)) + ":\n    uses: " + u + "\n")
	}
	return b.String()
}

// found is a saved GitHub answer for a path that exists.
func found(match, body string) response {
	return response{Match: []string{"api " + match}, Stdout: body}
}

// notFound is a saved GitHub answer for a path that does not exist.
func notFound(match string) response {
	return response{Match: []string{"api " + match}, Stderr: "gh: Not Found (HTTP 404)\n", Exit: 1}
}

// refCheck runs workflow-ref-exists on happyRepo with .github/workflows/ci.yml
// set to workflow, GitHub answering from saved responses.
func refCheck(t *testing.T, workflow string, saved []response) (string, int, replays) {
	t.Helper()
	rp := standIn(t, map[string][]response{"gh": saved})
	repo := example(t, "happyRepo", map[string]string{".github/workflows/ci.yml": workflow})
	out, code := fitness(t, repo, rp.env, "workflow-ref-exists")
	return out, code, rp
}

// rowDuration matches a results-table row's time, in milliseconds or seconds.
var rowDuration = regexp.MustCompile(`│ (\d+(?:\.\d+)?)(ms|s)\s+│\s*$`)

// slowRefs is how many pinned workflow calls the "workflow-ref-exists reads pinned calls in parallel" case makes, each answered
// after one second. Read one at a time, they take slowRefs seconds.
const slowRefs = 8

// reportStep is a last step that reports every outcome to the job summary.
const reportStep = `      - name: Report job outcome
        if: ${{ always() }}
        run: printf 'Job %s\n' "${{ job.status }}" >> "$GITHUB_STEP_SUMMARY"
`

// reportingJob is a workflow whose job, on line 4, builds and then
// reports with tail as its last steps.
func reportingJob(tail string) string {
	return "name: CI\non: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - run: make\n" + tail
}

// reportingCheck runs workflow-reporting on happyRepo with files written
// over it.
func reportingCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "workflow-reporting")
}

// mergedRunsOn is the runs-on line that merged into three workflows and
// failed every run, since GitHub reads its `": "` as a mapping value.
const mergedRunsOn = `${{ github.event.repository.private && !github.event.pull_request.head.repo.fork && fromJSON('{"group": "fitness", "labels": "fitness"}') || 'ubuntu-24.04-arm' }}`

// runsOnJob is a reporting workflow whose job runs on runsOn, on line 5.
func runsOnJob(runsOn string) string {
	return strings.Replace(reportingJob(reportStep), "runs-on: ubuntu-latest", "runs-on: "+runsOn, 1)
}

// rejectedRunsOn are runs-on values GitHub rejects as YAML, by file,
// each failing on line 5 or 6 at the column its message follows.
var rejectedRunsOn = map[string]string{
	"fork.yml":          mergedRunsOn,
	"reserved.yml":      "@self-hosted",
	"key-line.yml":      "- self-hosted",
	"continued.yml":     "ubuntu\n      latest: now",
	"colon-end.yml":     "ubuntu-latest:",
	"flow-trailing.yml": "[self-hosted]: x",
	"wide.yml":          "café: x",
}

// acceptedRunsOn are the quoted or reshaped forms GitHub accepts.
var acceptedRunsOn = map[string]string{
	"fork.yml":          `"` + strings.ReplaceAll(mergedRunsOn, `"`, `\"`) + `"`,
	"reserved.yml":      `"@self-hosted"`,
	"key-line.yml":      "\n      - self-hosted",
	"continued.yml":     "ubuntu\n      latest now",
	"colon-end.yml":     `"ubuntu-latest:"`,
	"flow-trailing.yml": "[self-hosted] # a comment",
	"wide.yml":          `"café: x"`,
}

// runsOnWorkflows writes each runs-on value into a reporting workflow.
func runsOnWorkflows(values map[string]string) map[string]string {
	files := map[string]string{}
	for name, runsOn := range values {
		files[".github/workflows/"+name] = runsOnJob(runsOn)
	}
	return files
}

func Test0007_9(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"workflow-scripts passes trivial commands and script calls", func(t *testing.T) {
			t.Parallel()
			steps := "      - run: make test\n      - run: scripts/release.sh\n      - run: |\n          # a comment\n          make build\n"
			action := "name: A\nruns:\n  using: composite\n  steps:\n    - shell: bash\n      run: \"${{ github.action_path }}/run.sh\"\n"
			out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps), "action.yml": action})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "workflow-scripts", "2")
		}},
		{"workflow-scripts rejects several lines of shell", func(t *testing.T) {
			t.Parallel()
			steps := "      - name: Tag\n        run: |\n          git fetch\n          git tag v1\n"
			out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
			sees(t, out, code, 1, ".github/workflows/ci.yml:8: this run step holds 2 lines of shell",
				"move it into a checked-in script such as scripts/<name>.sh, and call that script")
		}},
		{"workflow-scripts rejects a pipe, a chain and control flow", func(t *testing.T) {
			t.Parallel()
			action := "name: A\nruns:\n  using: composite\n  steps:\n    - shell: bash\n      run: cat out | grep ok\n    - shell: bash\n      run: make && make test\n" +
				"    - shell: bash\n      run: if true; then make; fi\n    - shell: bash\n      run: \"make build && make test\"\n"
			out, code := scriptsCheck(t, map[string]string{"action.yml": action})
			sees(t, out, code, 1, "action.yml:6: this run step holds a pipe", "action.yml:8: this run step holds an && chain",
				"action.yml:10: this run step holds shell control flow", "action.yml:12: this run step holds an && chain")
		}},
		{"workflow-scripts passes operators in quotes and expressions", func(t *testing.T) {
			t.Parallel()
			steps := "      - run: printf '%s | %s\\n' a b >> \"$GITHUB_STEP_SUMMARY\"\n      - run: echo \"${{ a || b }}\" \"x;y\"\n" +
				"      - run: make test # lint | format\n      - run: find . -name x -exec rm {} \\;\n      - run: grep ok <<< \"$out\"\n"
			out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"workflow-scripts judges a plain scalar and leaves block text alone", func(t *testing.T) {
			t.Parallel()
			steps := "      - run:\n          make a\n          make b\n      - run: make c\n          make d\n      - run: |\n          run: a | b\n"
			out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
			sees(t, out, code, 1, ".github/workflows/ci.yml:7: this run step holds 2 lines of shell", ".github/workflows/ci.yml:10: this run step holds 2 lines of shell")
			if strings.Contains(out, ".github/workflows/ci.yml:13") {
				t.Errorf("want text inside a block left alone:\n%s", out)
			}
		}},
	})
}

// scriptsCheck runs workflow-scripts on happyRepo with files written over it.
func scriptsCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "workflow-scripts")
}

// stepsWorkflow is a workflow whose job runs steps, which start on line 7.
func stepsWorkflow(steps string) string {
	return "name: CI\non: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n" + steps
}
