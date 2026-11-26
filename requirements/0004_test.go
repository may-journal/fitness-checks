package requirements

import (
	"encoding/base64"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func Test0004_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"long words in one endless sentence", func(t *testing.T) {
			t.Parallel()
			out, code := textReadability(t, map[string]string{"dense.md": textReadabilityDense})
			sees(t, out, code, 1,
				"dense.md: readability alarm on 141 prose words",
				"(grade bands, alarm at 18)", "(alarm at 60)",
				"to fix: split sentences over ~25 words",
				"methodology: go/cmd/fitness-check-text-readability/README.md")
		}},
		{"plain short sentences pass", func(t *testing.T) {
			t.Parallel()
			out, code := textReadability(t, map[string]string{"plain.md": textReadabilityPlain})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"dense doc under hundred prose words", func(t *testing.T) {
			t.Parallel()
			short := strings.Repeat("extraordinary bureaucratic infrastructure considerations ", 20) + "conclude.\n"
			out, code := textReadability(t, map[string]string{"short.md": short})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"only one readability score in alarm", func(t *testing.T) {
			t.Parallel()
			// Half seven-letter words in 21-word sentences: only the long-word score is high.
			lix := strings.Repeat(strings.Repeat("journey to ", 10)+"end. ", 6)
			out, code := textReadability(t, map[string]string{"lix.md": lix}, "--report")
			sees(t, out, code, 0, "LIX 68.6", "All 1 checks passed")
		}},
		{"dense notation inside fenced code", func(t *testing.T) {
			t.Parallel()
			doc := textReadabilityPlain + "\n\n```text\n" + textReadabilityDense + "```\n"
			out, code := textReadability(t, map[string]string{"notation.md": doc})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"dense issue or pr body", func(t *testing.T) {
			t.Parallel()
			out, code := textReadability(t, map[string]string{"body.md": textReadabilityDense}, "--body-file", "body.md")
			sees(t, out, code, 1, "(description): readability alarm on 141 prose words")
		}},
		{"repo lowers its readability bands", func(t *testing.T) {
			t.Parallel()
			out, code := textReadability(t, map[string]string{
				".fitnessrc.json": `{"textReadability": {"maxGrade": 5, "minWords": 10}}` + "\n",
				"plain.md":        textReadabilityPlain,
			})
			sees(t, out, code, 1, "plain.md: readability alarm on 140 prose words", "(grade bands, alarm at 5)")
		}},
		{"report flag prints word counts", func(t *testing.T) {
			t.Parallel()
			out, code := textReadability(t, map[string]string{"plain.md": textReadabilityPlain}, "--report")
			sees(t, out, code, 0, "README.md 3 words (below 100 — not judged)", "plain.md 140 words", "All 1 checks passed")
		}},
	})
}

func Test0004_2(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"misspellings in short and long files", func(t *testing.T) {
			t.Parallel()
			// 2,500 filler lines are 65,000 bytes, so the block opens just before
			// cspell's 64 KB chunk boundary and closes after it.
			filler := strings.Repeat("good words in a long file\n", 2500)
			long := "the " + cspellTypo + " fox\n" + filler +
				"cspell:disable\n" + filler[:2600] + "the " + cspellBodyTypo + " fox\ncspell:enable\n" +
				filler + "the " + cspellRepoWord + " fox\n"
			out, code := cspellRun(t, map[string]string{
				"notes.md": "good words\nthe " + cspellTypo + " fox\n",
				"long.md":  long,
			})
			first := "long.md:1:5 - Unknown word (" + cspellTypo + ")"
			last := "long.md:5105:5 - Unknown word (" + cspellRepoWord + ")"
			sees(t, out, code, 1, "notes.md:2:5 - Unknown word ("+cspellTypo+")", first, last)
			if strings.Contains(out, cspellBodyTypo) || strings.Index(out, first) > strings.Index(out, last) {
				t.Errorf("only the words outside the block may be named, in line order:\n%s", out)
			}
		}},
		{"repo word list beside shared one", func(t *testing.T) {
			t.Parallel()
			out, code := cspellRun(t, map[string]string{
				"cspell.json":   `{"words":["fitnessrc","` + cspellRepoWord + `"]}` + "\n",
				cspellInstalled: `{"words":["fitnessrc","` + cspellSharedWord + `"]}` + "\n",
				"notes.md":      cspellRepoWord + " and " + cspellSharedWord + "\n@jasonvanderslice from Vanderbot\n",
			})
			sees(t, out, code, 1, "notes.md:1:11 - Unknown word ("+cspellSharedWord+")")
			if strings.Contains(out, "Unknown word ("+cspellRepoWord+")") {
				t.Errorf("the repo's own word must be accepted:\n%s", out)
			}
			if strings.Contains(out, "Unknown word (jasonvanderslice)") || strings.Contains(out, "Unknown word (Vanderbot)") {
				t.Errorf("the owner and agent names must be accepted:\n%s", out)
			}
		}},
		{"ignore paths setting is refused", func(t *testing.T) {
			t.Parallel()
			out, code := cspellRun(t, map[string]string{"cspell.json": `{"ignorePaths":["docs"]}` + "\n"})
			sees(t, out, code, 1, "cspell.json sets ignorePaths; cspell checks every tracked file, so remove ignorePaths and fix the findings instead")
		}},
		{"unknown words in exempt files", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{
				".gitattributes":       "*.gen linguist-generated\n",
				"deps/pkg.gen":         cspellSharedWord + "\n",
				"blob.bin":             cspellBinaryWord + "\x00" + cspellBinaryOther,
				"package-lock.json":    `{"packages":{"node_modules/` + cspellTypo + `":{}}}` + "\n",
				"web/yarn.lock":        cspellRepoWord + "@1.0.0\n",
				"pnpm-lock.yaml":       cspellBodyTypo + ": 1.0.0\n",
				"go.sum":               "example.com/" + cspellBinaryOther + " v1.0.0 h1:abc=\n",
				"App/Package.resolved": `{"identity":"` + cspellBinaryWord + `"}` + "\n",
				"lib/gradle.lockfile":  "com.example:" + cspellTypo + ":1.0=compileClasspath\n",
				"mvnw":                 "url=" + cspellTypo + "\n",
				"app/gradlew.bat":      "set " + cspellRepoWord + "=1\n",
				".mvn/wrapper/MavenWrapperDownloader.java": "class Downloader { String " + cspellRepoWord + "; }\n",
				"gradle/wrapper/gradle-wrapper.properties": "distributionUrl=" + cspellBodyTypo + "\n",
			})
			write(t, repo, map[string]string{"untracked.md": "the " + cspellRepoWord + " fox\n"})
			out, code := fitness(t, repo, nil, "cspell")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"staged misspelling versus committed one", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{"old.md": "the " + cspellTypo + " fox\n", "gone.md": "clean words\n"})
			write(t, repo, map[string]string{"new.md": "the " + cspellRepoWord + " fox\n"})
			git(t, repo, "add", "new.md")
			git(t, repo, "rm", "-q", "gone.md")
			out, code := fitness(t, repo, nil, "cspell")
			sees(t, out, code, 1, "new.md:1:5 - Unknown word ("+cspellRepoWord+")")
			if strings.Contains(out, cspellTypo) || strings.Contains(out, "gone.md") {
				t.Errorf("only the staged file on disk may be checked:\n%s", out)
			}
		}},
		{"issue or pr body misspelling", func(t *testing.T) {
			t.Parallel()
			out, code := cspellRun(t, map[string]string{"body.md": "a " + cspellBodyTypo + " word\n"}, "--body-file", "body.md")
			sees(t, out, code, 1, "(description):1:3 - Unknown word ("+cspellBodyTypo+")")
		}},
		{"no repo word list falls back", func(t *testing.T) {
			t.Parallel()
			out, code := cspellRun(t, map[string]string{
				cspellInstalled: `{"words":["fitnessrc","` + cspellSharedWord + `"]}` + "\n",
				"notes.md":      cspellSharedWord + " is a shared word\n",
			})
			sees(t, out, code, 0, "All 1 checks passed")
			out, code = cspellRun(t, map[string]string{"notes.md": "mayjournal is a shared word, and addrinfo a language one\n"})
			sees(t, out, code, 0, "All 1 checks passed")
			out, code = cspellRun(t, map[string]string{"body.md": "sync `fastlane/Appfile`, `fastlane/Fastfile`, `fastlane/Gymfile`\n"}, "--body-file", "body.md")
			sees(t, out, code, 0, "All 1 checks passed")
			out, code = cspellRun(t, map[string]string{"static-preview.yml": "run: touch \"$PAGES/.nojekyll\"\n"})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"forbidden flag word phrases", func(t *testing.T) {
			t.Parallel()
			files := map[string]string{
				"cspell.json": `{"words":["fitnessrc"],"flagWords":["green apple"]}` + "\n",
				"notes.md":    "a Green  Apple and green" + "apple\n",
				"body.md":     "no green\napple here\n",
			}
			out, code := cspellRun(t, files)
			sees(t, out, code, 1, "notes.md:1:3 - Forbidden word (green apple)", "body.md:1:4 - Forbidden word (green apple)")
			if strings.Contains(out, "cspell.json:") {
				t.Errorf("the cspell.json listing flagWords must not be flagged:\n%s", out)
			}
			out, code = cspellRun(t, files, "--body-file", "body.md")
			sees(t, out, code, 1, "(description):1:4 - Forbidden word (green apple)")
		}},
		{"renovate hidden debug block passes", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			out, code := piped(t, repo, renovateBody, nil, "--blob", "-", "--checks=cspell")
			sees(t, out, code, 0, "All 1 checks passed")
			out, code = piped(t, repo, misspelledNote+renovateBody, nil, "--blob", "-", "--checks=cspell")
			sees(t, out, code, 1, "Unknown word (mispeled)")
		}},
		{"asterisk bold and italic", func(t *testing.T) {
			t.Parallel()
			out, code := noBoldItalicReadme(t, "Some **strong** and *soft* words.\n")
			sees(t, out, code, 1,
				`README.md: disallowed **bold**`+noBoldItalicRule+`"**strong**"`,
				`README.md: disallowed *italic*`+noBoldItalicRule+`"*soft*"`)
		}},
		{"underscore bold and italic", func(t *testing.T) {
			t.Parallel()
			out, code := noBoldItalicReadme(t, "Some __strong__ and _soft_ words.\n")
			sees(t, out, code, 1,
				`README.md: disallowed __bold__`+noBoldItalicRule+`"__strong__"`,
				`README.md: disallowed _italic_`+noBoldItalicRule+`"_soft_"`)
		}},
		{"emphasis inside code passes", func(t *testing.T) {
			t.Parallel()
			out, code := noBoldItalicReadme(t, "Use `**strong**` here.\n\n```\n_soft_ and **strong**\n```\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"underscores inside a markdown link", func(t *testing.T) {
			t.Parallel()
			out, code := noBoldItalicReadme(t, "See [the _dev_ path](https://example.com/_dev_/to_page_).\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"asterisk list item markers pass", func(t *testing.T) {
			t.Parallel()
			out, code := noBoldItalicReadme(t, "* first item\n* second item\n* third item\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"italic text in a body", func(t *testing.T) {
			t.Parallel()
			out, code := noBoldItalic(t, map[string]string{"body.md": "Hello _there_.\n"}, "--body-file", "body.md")
			sees(t, out, code, 1, `(description): disallowed _italic_`+noBoldItalicRule+`"_there_"`)
		}},
		{"changelog bold beside release scopes", func(t *testing.T) {
			t.Parallel()
			release := strings.Replace(notesRelease(notesTwo), "change it", "change **everything**", 1)
			out, code := notesChecks(t, notesHead+release+notesEntry, "markdown-no-bold-italic")
			sees(t, out, code, 1, `CHANGELOG.md: disallowed **bold**`+noBoldItalicRule+`"**everything**"`)
			if n := strings.Count(out, "disallowed"); n != 1 {
				t.Errorf("only the bold outside a scope must fail, got %d findings:\n%s", n, out)
			}
		}},
		{"negated sentence then restating sentence", func(t *testing.T) {
			t.Parallel()
			out, code := reframing(t, "It's not a workout. It's a lifestyle.\n\nThey’re not bugs.\n\nThey’re features.\n")
			sees(t, out, code, 1, reframingHint+`"It's not a workout It's a lifestyle"`,
				reframingHint+`"They’re not bugs They’re features"`)
		}},
		{"single sentence negates then pivots", func(t *testing.T) {
			t.Parallel()
			out, code := reframing(t, "This isn't about speed; it's about consistency.\n")
			sees(t, out, code, 1, reframingHint+`"This isn't about speed; it's about consistency"`)
		}},
		{"ordinary negations pass", func(t *testing.T) {
			t.Parallel()
			out, code := reframing(t, "It's not ready yet. The file is not found, but the fallback loads.\n\n"+
				"We should not merge, but wait for review. Not only fast, but also correct.\n\n"+
				"It's not X. It's not Y.\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"pattern in front matter or code", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", map[string]string{
				"README.md": "---\ntitle: It's not X. It's Y.\n---\n\n# App\n\n```\nIt's not a workout. It's a lifestyle.\n```\n\n" +
					"## It's not a workout. It's a lifestyle.\n\n" +
					"Use `It's not a workout. It's a lifestyle.` verbatim.\n",
			})
			out, code := fitness(t, repo, nil, "no-contrastive-reframing")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"pattern inside quotes passes", func(t *testing.T) {
			t.Parallel()
			out, code := reframing(t, `The phrase "It's not a workout. It's a lifestyle." is banned.`+"\n\n"+
				"The phrase “It's not a workout. It's a lifestyle.” is banned.\n")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"long reframed sentence is shortened", func(t *testing.T) {
			t.Parallel()
			long := "It's not " + strings.Repeat("a very long clause ", 10) + ", but simple."
			out, code := reframing(t, long+"\n")
			sees(t, out, code, 1, reframingHint+`"`+long[:80]+`…"`)
		}},
	})
}

func Test0004_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"hard to read comment", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", inComments(hardToRead)), nil, "text-readability")
			sees(t, out, code, 1, atEachFile("readability alarm on 104 prose words")...)
		}},
		{"comment saying it is not x", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", inComments("It's not a cache. It's the source of truth.")), nil, "no-contrastive-reframing")
			sees(t, out, code, 1, atEachFile(`rewrite as a direct statement, not the "not X, it's Y" pattern`)...)
		}},
		{"why comments directives and code pass", func(t *testing.T) {
			t.Parallel()
			words := strings.TrimSpace(longSentence)
			files := inComments("Retries once, since the first call warms the cache.")
			files["app/main.go"] = "//go:build " + strings.Repeat("linux && ", 29) + "darwin\n\n" + files["app/main.go"]
			files["app/build.go"] = "package main\n\n// cspell:ignore " + words + "\n\n// Runs the build:\n//\n//\t" + words + "\nfunc build() {}\n"
			files["app/Build.swift"] = "// cspell:ignore " + words + "\n/// Builds the app once.\nfunc build() {}\n"
			files["app/doc.go"] = "// " + strings.Repeat("One short paragraph here.\n//\n// ", 7) + "The end is here.\npackage main\n"
			files["app/build.ts"] = "// cspell:ignore " + words + "\n// Runs:\n//\n//     " + words + "\nexport const y = 2;\n"
			out, code := fitness(t, example(t, "happyRepo", files), nil,
				"--policy=external", "--checks=prose-budget,text-readability,no-contrastive-reframing")
			sees(t, out, code, 0, "All 3 checks passed")
		}},
	})
}

func Test0004_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"hyphenated name in camelcase repo", func(t *testing.T) {
			t.Parallel()
			out, code := camelDocs(t, camelDocsRepo(map[string]string{"docs/api-design.md": "# API\n"}))
			sees(t, out, code, 1, "docs/api-design.md: filename must be camelCase")
		}},
		{"snake and pascal names in camelcase repo", func(t *testing.T) {
			t.Parallel()
			out, code := camelDocs(t, camelDocsRepo(map[string]string{
				"docs/api_design.md": "# API\n",
				"docs/ApiDesign.md":  "# API\n",
			}))
			sees(t, out, code, 1, "docs/api_design.md: filename must be camelCase", "docs/ApiDesign.md: filename must be camelCase")
		}},
		{"all camelcase names pass", func(t *testing.T) {
			t.Parallel()
			out, code := camelDocs(t, camelDocsRepo(nil))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"camelcase doc in capitalized folder", func(t *testing.T) {
			t.Parallel()
			out, code := camelDocs(t, camelDocsRepo(map[string]string{"Guides/Api_Docs/apiDesign.md": "# API\n"}))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"mostly hyphenated repo skips camelcase", func(t *testing.T) {
			t.Parallel()
			out, code := camelDocs(t, map[string]string{
				"docs/api-design.md":   "# API\n",
				"docs/user-guide.md":   "# Guide\n",
				"docs/releaseNotes.md": "# Release notes\n",
			})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, camelDocsCheck)
		}},
		{"tie goes to kebab case", func(t *testing.T) {
			t.Parallel()
			out, code := camelDocs(t, map[string]string{
				"docs/api-design.md":   "# API\n",
				"docs/releaseNotes.md": "# Release notes\n",
			})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, camelDocsCheck)
		}},
		{"staged hyphenated doc in camelcase repo", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", camelDocsRepo(map[string]string{
				"docs/old-notes.md": "# Old\n",
				"docs/userGuide.md": "# Guide\n",
			}))
			write(t, repo, map[string]string{"docs/new-notes.md": "# New\n"})
			git(t, repo, "add", "-A")
			out, code := fitness(t, repo, nil, camelDocsCheck)
			sees(t, out, code, 1, "docs/new-notes.md: filename must be camelCase")
			if strings.Contains(out, "old-notes.md") {
				t.Errorf("only the staged file may be judged:\n%s", out)
			}
		}},
		{"standard and github named docs exempt", func(t *testing.T) {
			t.Parallel()
			out, code := camelDocs(t, camelDocsRepo(map[string]string{
				"CODE_OF_CONDUCT.md":               "# Code of conduct\n",
				"docs/AGENTS.md":                   "# Agents\n",
				".github/pull_request_template.md": "# Pull request\n",
				"docs/issue_template.md":           "# Issue\n",
			}))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"camelcase name in kebab repo", func(t *testing.T) {
			t.Parallel()
			out, code := kebabDocs(t, kebabDocsRepo(map[string]string{"docs/releaseNotes.md": "# Notes\n"}))
			sees(t, out, code, 1, "docs/releaseNotes.md: filename must be kebab-case")
		}},
		{"snake, pascal and spaced names", func(t *testing.T) {
			t.Parallel()
			out, code := kebabDocs(t, kebabDocsRepo(map[string]string{
				"docs/api_notes.md": "# Notes\n",
				"docs/ApiNotes.md":  "# Notes\n",
				"docs/api notes.md": "# Notes\n",
			}))
			sees(t, out, code, 1,
				"docs/api_notes.md: filename must be kebab-case",
				"docs/ApiNotes.md: filename must be kebab-case",
				"docs/api notes.md: filename must be kebab-case")
		}},
		{"all kebab names pass", func(t *testing.T) {
			t.Parallel()
			out, code := kebabDocs(t, kebabDocsRepo(map[string]string{
				"architecture/02-containers.md": "# Containers\n",
				"docs/adr001.md":                "# ADR\n",
			}))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"kebab doc in capitalized folder", func(t *testing.T) {
			t.Parallel()
			out, code := kebabDocs(t, kebabDocsRepo(map[string]string{"Guides/Api_Docs/api-design.md": "# API\n"}))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"mostly camelcase repo skips kebab", func(t *testing.T) {
			t.Parallel()
			out, code := kebabDocs(t, map[string]string{
				"docs/releaseNotes.md": "# Notes\n",
				"docs/installGuide.md": "# Guide\n",
				"docs/api-design.md":   "# API\n",
			})
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, kebabDocsCheck)
		}},
		{"tie fails camelcase file", func(t *testing.T) {
			t.Parallel()
			out, code := kebabDocs(t, map[string]string{
				"docs/api-design.md":   "# API\n",
				"docs/releaseNotes.md": "# Notes\n",
			})
			sees(t, out, code, 1, "docs/releaseNotes.md: filename must be kebab-case")
		}},
		{"staged camelcase doc in kebab repo", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", kebabDocsRepo(map[string]string{
				"docs/oldNotes.md":    "# Old\n",
				"docs/setup-steps.md": "# Setup\n",
			}))
			write(t, repo, map[string]string{"docs/newNotes.md": "# New\n"})
			git(t, repo, "add", "-A")
			out, code := fitness(t, repo, nil, kebabDocsCheck)
			sees(t, out, code, 1, "docs/newNotes.md: filename must be kebab-case")
			if strings.Contains(out, "oldNotes.md") {
				t.Errorf("only the staged file may be judged:\n%s", out)
			}
		}},
		{"standard and github named docs exempt in kebab", func(t *testing.T) {
			t.Parallel()
			out, code := kebabDocs(t, kebabDocsRepo(map[string]string{
				"CODE_OF_CONDUCT.md":               "# Code of conduct\n",
				"docs/AGENTS.md":                   "# Agents\n",
				".github/pull_request_template.md": "# Pull request\n",
				"docs/issue_template.md":           "# Issue\n",
			}))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"link to a missing file", func(t *testing.T) {
			t.Parallel()
			out, code := markdownLinks(t, markdownLinksReadme(`See [the guide](gone.md "Guide") and [home](README.md).`+"\n"))
			sees(t, out, code, 1, "README.md:9: broken relative link: gone.md")
		}},
		{"links to file folder and heading", func(t *testing.T) {
			t.Parallel()
			files := markdownLinksReadme("See [guide](./docs/guide.md#setup), [docs](docs/), and [home](docs/../README.md).\n")
			files["docs/guide.md"] = "# Guide\n\n## Setup\n\nBack to [the readme](../README.md).\n"
			out, code := markdownLinks(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"missing sibling from a subfolder", func(t *testing.T) {
			t.Parallel()
			out, code := markdownLinks(t, map[string]string{"docs/guide.md": "# Guide\n\nSee [setup](setup.md#install).\n"})
			sees(t, out, code, 1, "docs/guide.md:3: broken relative link: setup.md#install")
		}},
		{"image or reference to missing file", func(t *testing.T) {
			t.Parallel()
			out, code := markdownLinks(t, markdownLinksReadme("![Logo](logo.png)\n\n[spec]: spec.md\n"))
			sees(t, out, code, 1, "README.md:9: broken relative link: logo.png", "README.md:11: broken relative link: spec.md")
		}},
		{"web mail and same page links", func(t *testing.T) {
			t.Parallel()
			out, code := markdownLinks(t, markdownLinksReadme(
				"See [site](https://example.com/gone), [mail](mailto:a@example.com), [cdn](//example.com/x), and [top](#app).\n"))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"broken link inside code", func(t *testing.T) {
			t.Parallel()
			out, code := markdownLinks(t, markdownLinksReadme("Write `[text](gone.md)` like this.\n\n```\n[text](gone.md)\n```\n"))
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"percent encoded link with space", func(t *testing.T) {
			t.Parallel()
			files := markdownLinksReadme("See [notes](docs/release%20notes.md).\n")
			files["docs/release notes.md"] = "# Release Notes\n"
			out, code := markdownLinks(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"markdown and mdx endings count", func(t *testing.T) {
			t.Parallel()
			out, code := markdownLinks(t, map[string]string{
				"docs/guide.markdown": "# Guide\n\nSee [setup](setup.md).\n",
				"docs/intro.mdx":      "# Intro\n\nSee [the guide](gone.mdx).\n",
			})
			sees(t, out, code, 1, "docs/guide.markdown:3: broken relative link: setup.md", "docs/intro.mdx:3: broken relative link: gone.mdx")
		}},
	})
}

func Test0004_5(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"github templates do not govern files", func(t *testing.T) {
			t.Parallel()
			out, code := docTemplate(t, map[string]string{
				"docs/arch/templates/adr.md":       docTemplateADR,
				"docs/arch/templates/note.md":      docTemplateNote,
				"docs/arch/04-code.md":             "# Code\n\n## Summary\n\ns\n",
				".github/PULL_REQUEST_TEMPLATE.md": "## Background\n\n<!-- context -->\n",
				".github/other.md":                 "# Other\n\n## Totally\n\ndifferent\n",
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"section the template lacks or has", func(t *testing.T) {
			t.Parallel()
			out, code := docTemplate(t, map[string]string{
				"docs/adr/template.md":   docTemplateADR,
				"docs/adr/0001-extra.md": "# 0001 Extra\n\n## Context\n\nc\n\n## Status\n\ns\n\n## Decision\n\nd\n\n## Consequences\n\ne\n",
				"docs/adr/0002-short.md": "# 0002 Short\n\n## Context\n\nc\n\n## Consequences\n\ne\n",
			})
			sees(t, out, code, 1,
				"docs/adr/0001-extra.md: Unexpected `## Status` section",
				"docs/adr/0002-short.md: Missing `## Decision` section")
		}},
		{"nested and named templates nearest wins", func(t *testing.T) {
			t.Parallel()
			out, code := docTemplate(t, map[string]string{
				"docs/templates/adr.md":           docTemplateADR,
				"docs/deep/note.template.md":      "# Note\n\n## Summary\n\n<!-- gist -->\n",
				"docs/deep/sub/0001-note.md":      "# 0001 Note\n\n## Context\n\nc\n",
				"docs/views/templates/04-code.md": docTemplateNote,
				"docs/views/04-code.md":           "# Code\n\n## Summary\n\ns\n",
				"docs/views/notes.md":             "# Notes\n\n## Open threads\n\no\n",
			})
			sees(t, out, code, 1,
				"docs/deep/sub/0001-note.md: Missing `## Summary` section",
				"docs/views/notes.md: Missing `## Context` section",
				"docs/views/notes.md: Unexpected `## Open threads` section")
		}},
		{"template section lacks guiding comment", func(t *testing.T) {
			t.Parallel()
			out, code := docTemplate(t, map[string]string{
				"docs/adr/template.md":             "# NNNN Title\n\n## Context\n\n<!-- why -->\n\n## Decision\n\nno comment here\n",
				".github/PULL_REQUEST_TEMPLATE.md": "## Background\n\nno comment\n",
				".github/ISSUE_TEMPLATE/plan.md":   "## Goal\n\nno comment\n",
			})
			sees(t, out, code, 1,
				"docs/adr/template.md: section `## Decision` needs a guiding <!-- comment -->",
				".github/PULL_REQUEST_TEMPLATE.md: section `## Background` needs a guiding <!-- comment -->",
				".github/ISSUE_TEMPLATE/plan.md: section `## Goal` needs a guiding <!-- comment -->")
		}},
		{"doc matching neither of two templates", func(t *testing.T) {
			t.Parallel()
			out, code := docTemplate(t, map[string]string{
				"docs/adr/adr.template.md":  docTemplateADR,
				"docs/adr/note.template.md": docTemplateNote,
				"docs/adr/0001-bad.md":      "# 0001 Bad\n\n## Other\n\no\n",
			})
			sees(t, out, code, 1,
				"docs/adr/0001-bad.md: sections match none of its templates: docs/adr/adr.template.md, docs/adr/note.template.md")
		}},
		{"view lacks legend or callout structure", func(t *testing.T) {
			t.Parallel()
			files := dtCopied(t)
			dtEdit(t, files, acCode, "```mermaid\nclassDiagram\n    class Example", "```mermaid\nflowchart TB\n    class Example")
			dtEdit(t, files, acJourney, "## Journey", "## Steps")
			out, code := docTemplate(t, files)
			sees(t, out, code, 1,
				dtViews+acCode+": the opening needs a classDiagram legend as "+dtTemplates+acCode+" has",
				dtViews+acJourney+": sections must open with `## Legend`, `## Journey` as "+dtTemplates+acJourney+" does")

			files = dtCopied(t)
			dtEdit(t, files, acJourney, "| Type    |", "| Kind    |")
			dtEdit(t, files, acCode, "Numbers on classes", "Labels on classes")
			files[dtViews+acContainer] = files[dtViews+acContainer][:strings.Index(files[dtViews+acContainer], "| #")]
			journey := dtLine(t, files, acJourney, "| #   | Description          | Kind")
			table := dtLine(t, files, acCode, "| #")
			out, code = docTemplate(t, files)
			sees(t, out, code, 1,
				journey+": callout table columns #, Description, Kind, Why must be #, Description, Type, Why as "+dtTemplates+acJourney+" has",
				table+`: callout table needs the caption "Numbers on classes and letters on relationships match the callout table" above it as `+dtTemplates+acCode+" has",
				dtViews+acContainer+": the opening needs a callout table with columns #, Description, Why as "+dtTemplates+acContainer+" has")

			files = dtShipped(t)
			dtEdit(t, files, acContainer, "C4Container\n    title Product", "C4Context\n    title Product")
			dtEdit(t, files, acComponent, "C4Component\n    title Legend", "C4Container\n    title Legend")
			dtEdit(t, files, acCode, "```mermaid\nclassDiagram\n    class Record", "```mermaid\nflowchart TB\n    class Record")
			dtEdit(t, files, acJourney, "flowchart TB", "flowchart LR")
			out, code = docTemplate(t, files)
			sees(t, out, code, 1,
				dtLine(t, files, acContainer, "```mermaid")+": diagram type C4Context must be C4Container as "+dtShippedName(acContainer)+" has",
				dtViews+acComponent+": the opening needs a C4Component legend as "+dtShippedName(acComponent)+" has",
				": legend type C4Container must be C4Component as "+dtShippedName(acComponent)+" has",
				dtLine(t, files, acCode, "```mermaid")+": diagram type flowchart TB must be classDiagram as "+dtShippedName(acCode)+" has",
				dtLine(t, files, acJourney, "```mermaid")+": legend type flowchart LR must be flowchart TB as "+dtShippedName(acJourney)+" has")
		}},
		{"views from shipped architecture templates", func(t *testing.T) {
			t.Parallel()
			files := dtCopied(t)
			files[dtViews+acContainer] += "\n## Setup\n\nHow a consumer runs the app.\n"
			files[dtViews+acJourney] += "\n## Notes\n\nWhat the journey leaves out.\n"
			dtEdit(t, files, acComponent, "match the callout table.", "match the callout table. <!-- kept in step with the diagram -->")
			files[dtViews+"notes.md"] = "# Notes\n\n## Open threads\n\nWhat the views leave out.\n"
			files[dtViews+"archive/"+acCode] = "# Old code\n\n## Before\n\nThe first layout.\n"
			out, code := docTemplate(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
			passedWith(t, out, "doc-template", "14")

			files = dtShipped(t)
			files[dtViews+acContainer] += "\n## Setup\n\nHow a consumer runs the app.\n"
			out, code = docTemplate(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
			dtEdit(t, files, acJourney, "## Journey", "## Steps")
			dtEdit(t, files, acCode, "Numbers on classes", "Labels on classes")
			out, code = docTemplate(t, files)
			sees(t, out, code, 1,
				dtViews+acJourney+": sections must open with `## Legend`, `## Journey` as "+dtShippedName(acJourney)+" does",
				`: callout table needs the caption "Numbers on classes and letters on relationships match the callout table" above it as `+dtShippedName(acCode)+" has")
		}},
		{"view styles differ from template", func(t *testing.T) {
			t.Parallel()
			files := dtShipped(t)
			dtEdit(t, files, acContext, "    UpdateElementStyle(user,",
				"    SystemDb_Ext(db, \"4 Records\", \"Keeps records.\")\n    UpdateElementStyle(db, $fontColor=\"#FFFFFF\", $bgColor=\"#64748b\",  $borderColor=\"#475569\")\n    UpdateLayoutConfig($c4ShapeInRow=\"3\")\n    Lay_R(user, product)\n    UpdateElementStyle(user,")
			dtEdit(t, files, acCode, "    save --> Record : A", "    save --> Record : A\n    Store : +save()\n    <<interface>> Store\n    style Store "+paletteClass)
			dtEdit(t, files, acContainer, `Container(app, "3 App", "Language", "src/app")`, `Container(app, "3 App", $techn="Language")`)
			out, code := docTemplate(t, files)
			sees(t, out, code, 0, "All 1 checks passed")
			dtEdit(t, files, acContext, "UpdateElementStyle(product, $bgColor=\"#1d4ed8\"", "UpdateElementStyle(product, $bgColor=\"#ff0000\"")
			dtEdit(t, files, acContainer, "    UpdateElementStyle(app, "+paletteContainer+"\n", "")
			dtEdit(t, files, acComponent, "Container_Boundary(app,", "System_Boundary(app,")
			dtEdit(t, files, acComponent, "    Rel(user, entry, \"A\")", "    Rel(user, entry, \"A\")\n    UpdateRelStyle(user, entry, $textColor=\"red\")")
			dtEdit(t, files, acCode, "    save --> Record : A", "    save --> Record : A\n    classDef hot fill:#ff0000")
			out, code = docTemplate(t, files)
			sees(t, out, code, 1,
				"`UpdateElementStyle(product, $bgColor=\"#ff0000\", $borderColor=\"#1e3a8a\", $fontColor=\"#ffffff\")` must be `UpdateElementStyle(product, "+paletteSystem+"` as "+dtShippedName(acContext)+" styles each System",
				"Container app needs `UpdateElementStyle(app, "+paletteContainer+"` as "+dtShippedName(acContainer)+" styles each Container",
				dtLine(t, files, acComponent, "    System_Boundary(app,")+": System_Boundary app is a kind "+dtShippedName(acComponent)+" does not use; use Person, Container_Boundary, Component, System_Ext",
				"`UpdateRelStyle(user, entry, $textColor=\"red\")` is a style "+dtShippedName(acComponent)+" does not use",
				"`classDef hot fill:#ff0000` is a style "+dtShippedName(acCode)+" does not use")
			for _, styled := range []string{"Person user", "Component entry", "class save", "class Record"} {
				if strings.Contains(read(out), styled) {
					t.Errorf("%s keeps its palette and must not fail, got:\n%s", styled, out)
				}
			}

			files = dtShipped(t)
			dtEdit(t, files, acContext, "C4Context\n    title Product\n", "C4Context\n")
			dtEdit(t, files, acComponent, "C4Component\n    title Legend\n", "C4Component\n")
			dtEdit(t, files, acContainer, `Container(app, "3 App", "Language", "src/app")`, `Container(app, "3 App", "", "src/app")`)
			dtEdit(t, files, acComponent, `Component(store, "4 Store", "Language",`, `Component(store, "4 Store", " ",`)
			dtEdit(t, files, acContext, "    System_Ext(e, \"External\", \"Outside our boundary\")\n", "")
			dtEdit(t, files, acContext, c4Styled("e", paletteExternal), "")
			dtEdit(t, files, acContainer, "    System_Ext(service, \"4 Outside service\"", "    Person(service, \"4 Outside service\"")
			dtEdit(t, files, acContainer, c4Styled("service", paletteExternal), c4Styled("service", palettePerson))
			out, code = docTemplate(t, files)
			sees(t, out, code, 1,
				dtLine(t, files, acContext, "```mermaid")+": C4Context diagram needs a `title` line as "+dtShippedName(acContext)+" has",
				": C4Component legend needs a `title` line as "+dtShippedName(acComponent)+" has",
				dtLine(t, files, acContainer, "        Container(app,")+": Container app needs a technology as "+dtShippedName(acContainer)+" gives each Container",
				dtLine(t, files, acComponent, "        Component(store,")+": Component store needs a technology as "+dtShippedName(acComponent)+" gives each Component",
				dtViews+acContext+":", "legend needs a System_Ext sample, a kind the diagram uses",
				dtViews+acContainer+":", "legend shows System_Ext, a kind the diagram does not use")
		}},
		{"adr files in adr folder", func(t *testing.T) {
			t.Parallel()
			out, code := adrAsIssue(t, map[string]string{
				"docs/architecture/adr/0001-use-sqlite.md": adrFile,
				"docs/decisions/ADR/0002-sync.md":          adrFile,
				"docs/architecture/adr/README.md":          "# ADRs\n\n- [Use SQLite](0001-use-sqlite.md)\n",
			})
			sees(t, out, code, 1, inADRFolder("docs/decisions/ADR/0002-sync.md"), inADRFolder("docs/architecture/adr/0001-use-sqlite.md"))
			if strings.Contains(out, "adr/README.md:") {
				t.Errorf("the adr folder's README lists ADRs, yet it was reported:\n%s", out)
			}
		}},
		{"markdown file titled as adr", func(t *testing.T) {
			t.Parallel()
			out, code := adrAsIssue(t, map[string]string{
				"docs/storage.md":    "# ADR 12: Use SQLite\n\nWe sync offline.\n",
				"docs/adr-checks.md": "# adr-checks\n\nPairs each decision with a check.\n",
				"docs/guide.md":      "# Guide\n\n## ADR 3\n\nA later heading.\n",
			})
			sees(t, out, code, 1, "docs/storage.md: a title starting with ADR makes the file an ADR; "+toADRIssue)
			for _, kept := range []string{"docs/adr-checks.md", "docs/guide.md"} {
				if strings.Contains(out, kept+":") {
					t.Errorf("%s has no title starting with the word ADR, yet it was reported:\n%s", kept, out)
				}
			}
		}},
		{"check readme holds adr sections", func(t *testing.T) {
			t.Parallel()
			out, code := adrAsIssue(t, with(heldCheck(), map[string]string{
				"docs/writing-checks.md": "# Writing checks\n\n```markdown\n" + heldDecision + "```\n",
			}))
			sees(t, out, code, 1, "fitness/adr-0001/README.md: `## Context`, `## Decision` and `## Consequences` sections make the file an ADR; "+toADRIssue)
			if strings.Contains(out, "docs/writing-checks.md:") {
				t.Errorf("an ADR in a code block is an example, yet it was reported:\n%s", out)
			}
		}},
		{"every adr kept as linked issue", func(t *testing.T) {
			t.Parallel()
			out, code := adrAsIssue(t, map[string]string{
				"docs/architecture/adr/0001-use-sqlite.md": "https://github.com/may-journal/example/issues/41\n",
				"fitness/adr-0001/README.md":               "# 0001 — Small files\n\nSmall files read fast.\n\nhttps://github.com/may-journal/example/issues/5\n\n## Check\n\nFails a long Markdown file.\n",
				"fitness/adr-0001/check":                   verdictScript(passing),
				".github/ISSUE_TEMPLATE/adr.md":            "---\nname: ADR\nlabels: ADR\n---\n\n## Status\n\n## Context\n\n## Decision\n\n## Consequences\n",
				"docs/notes.md":                            "# Notes\n",
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"allowed adr files with reasons", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", with(heldCheck(), map[string]string{
				".fitnessrc.json": allowADRs(`{"path": "docs/architecture/adr/*.md", "reason": "the migration waits on repo access"}, ` +
					`{"path": "fitness/adr-*/README.md", "reason": "the checks move with them"}`),
				"docs/architecture/adr/0001-use-sqlite.md": adrFile,
				"docs/architecture/adr/0002-sync.md":       adrFile,
				"docs/decisions/adr/0003-theme.md":         adrFile,
			}))
			out, code := fitness(t, repo, nil, "adr-as-issue")
			sees(t, out, code, 1, inADRFolder("docs/decisions/adr/0003-theme.md"))
			for _, kept := range []string{"0001-use-sqlite.md", "0002-sync.md", "fitness/adr-0001/README.md"} {
				if strings.Contains(out, kept+":") {
					t.Errorf("%s is allowed, yet it was reported:\n%s", kept, out)
				}
			}
			out, code = fitness(t, repo, nil, "adr-checks")
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"allow entry without reason or match", func(t *testing.T) {
			t.Parallel()
			out, code := adrAsIssue(t, map[string]string{
				".fitnessrc.json":                          allowADRs(`{"path": "docs/architecture/adr/*.md"}, {"path": "docs/old/adr/*.md", "reason": "moved"}`),
				"docs/architecture/adr/0001-use-sqlite.md": adrFile,
			})
			sees(t, out, code, 1,
				`adrAsIssue.allow path "docs/architecture/adr/*.md" needs a reason`,
				`adrAsIssue.allow path "docs/old/adr/*.md" matches no ADR file; remove it`,
				inADRFolder("docs/architecture/adr/0001-use-sqlite.md"),
			)
		}},
	})
}

func Test0004_6(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"issue body with unchecked item", func(t *testing.T) {
			t.Parallel()
			out, code := issueChecklist(t, "- [x] Write the doc.\n- [ ] Write the tests.\n")
			sees(t, out, code, 1, `unchecked item: "Write the tests."`)
		}},
		{"issue body with every item ticked", func(t *testing.T) {
			t.Parallel()
			out, code := issueChecklist(t, "- [x] Write the doc.\n- [x] Write the tests.\n")
			sees(t, out, code, 0, "issue-checklist")
		}},
	})
}

func Test0004_8(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"comment over word budget", func(t *testing.T) {
			t.Parallel()
			out, code := commentStructure(t, overBudget, nil)
			sees(t, out, code, 1, "(comment): 120 words outside tables, code, and diagrams; the budget is 100")
		}},
		{"repo raises the word budget", func(t *testing.T) {
			t.Parallel()
			out, code := commentStructure(t, overBudget, map[string]string{".fitnessrc.json": `{"comment": {"maxWords": 150}}`})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"two plain tables in comment", func(t *testing.T) {
			t.Parallel()
			table := "| A | B |\n| - | - |\n| 1 | 2 |\n"
			out, code := commentStructure(t, table+"\nAnd more.\n\n"+table+"\n"+numbered, nil)
			sees(t, out, code, 1, "(comment): 2 tables besides callout tables; the budget is 1")
		}},
		{"unnumbered diagram versus numbered legend", func(t *testing.T) {
			t.Parallel()
			out, code := commentStructure(t, "Flow.\n\n```mermaid\nflowchart TD\n  A --> B\n```\n", nil)
			sees(t, out, code, 1, "(comment): line 3: mermaid diagram has no numbered callouts")
			legend := strings.Replace(numbered, "\n\n|", "\n\n```mermaid\nflowchart LR\n  k[Step]:::s\n  classDef s fill:#eef\n```\n\n|", 1)
			out, code = commentStructure(t, legend, nil)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"change log note in comment", func(t *testing.T) {
			t.Parallel()
			out, code := commentStructure(t, "The parser works.\n\nEdit 2: it handles tabs too.\n", nil)
			sees(t, out, code, 1, `(comment): line 3: "Edit 2:": change-log note; edit the text in place and drop the note`)
			out, code = commentStructure(t, "The cache is updated after every push.\n", nil)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"agent frame model id and marker", func(t *testing.T) {
			t.Parallel()
			out, code := commentStructure(t, "Done.\n\n<system-reminder>\n\nRan on claude-opus-4-1 today.\n\n<!-- fitness:comment-check:31:abc -->\n", nil)
			sees(t, out, code, 1, `(comment): line 3: "<system-reminder": agent frame; post only the finished text`,
				`(comment): line 5: "claude-opus-4-1": model identifier; name the agent, not the model`,
				`(comment): line 7: "<!-- fitness:": carries a fitness marker; only fitness writes those, so remove it`)
			out, code = commentStructure(t, "Pinned to `claude-opus-4-1` in config.\n", nil)
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"plan or epic body exempt", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			for _, body := range []string{planBody("- [ ] Write the tests."), epicBody("1. #22\n2. #23, after #22")} {
				out, code := piped(t, repo, body, nil, "--blob", "-")
				sees(t, out, code, 0, "Skipped (does not apply to this text): comment-structure")
			}
		}},
		{"no comment to check", func(t *testing.T) {
			t.Parallel()
			out, code := fitness(t, example(t, "happyRepo", nil), nil, "comment-structure")
			sees(t, out, code, 0, "All 1 checks passed")
			passedWithNoFiles(t, out, "comment-structure")
		}},
		{"pr and plan body questions line", func(t *testing.T) {
			t.Parallel()
			plan := planBody("- [x] Write the tests.")
			out, code := judgeBody(t, "test(app): add one file", prClosing(11), ghIssue(t, 11, plan, "Plan"))
			sees(t, out, code, 0, "Passed.")
			out, code = judgeBody(t, "test(app): add one file", questionsLabel+"\n"+prClosing(11), ghIssue(t, 11, plan, "Plan"))
			sees(t, out, code, 1, "(description):", questionsFinding)
			out, code, _ = judgeKind(t, planTemplate("- [ ] Write the tests."), nil, kindGitHub("plan", "Plan"))
			sees(t, out, code, 0, "Passed.")
			out, code, _ = judgeKind(t, planTemplate("- [ ] Write the tests.\n\n"+questionsLabel), nil, kindGitHub("plan", "Plan"))
			sees(t, out, code, 1, "(description):", questionsFinding)
		}},
		{"artifact page publish denied", func(t *testing.T) {
			t.Parallel()
			repo := example(t, "happyRepo", nil)
			page := filepath.Join(t.TempDir(), "plan.html")
			mustDo(t, os.WriteFile(page, []byte("<h1>Sync plan</h1>\n<h2>Questions for the owner</h2>\n<ul><li>Which wins?</li></ul>\n"), 0o644))
			out, code := hookEvent(t, repo, artifactEvent(page))
			sees(t, out, code, 0, `"permissionDecision":"deny"`, "failed on the artifact page plan.html", "[open-questions]",
				`(description):2: a heading opens a list of open questions: \"Questions for the owner\"`)
			mustDo(t, os.WriteFile(page, []byte("<h1>Sync plan</h1>\n<p>The server wins.</p>\n"), 0o644))
			out, code = hookEvent(t, repo, artifactEvent(page))
			silentPass(t, out, code)
		}},
	})
}

// proseBudget runs `fitness-install -- prose-budget <args>` on happyRepo
// with files written over it.
func proseBudget(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"prose-budget"}, args...)...)
}

// issueChecklist runs `fitness-install -- issue-checklist` on an issue body.
func issueChecklist(t *testing.T, body string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", map[string]string{"issue.md": body})
	return fitness(t, repo, nil, "issue-checklist", "--body-file", "issue.md")
}

// cspell:ignore mvnw gradlew

// Misspellings are built from pieces so this file's own spell check passes.
var (
	cspellTypo        = "qu" + "ik"
	cspellRepoWord    = "bor" + "wn"
	cspellSharedWord  = "zz" + "zq" + "qqv"
	cspellBinaryWord  = "xq" + "zzt"
	cspellBinaryOther = "vb" + "nmm"
	cspellBodyTypo    = "mis" + "pe" + "led"
)

// cspellRun runs `fitness-install -- cspell <args>` on happyRepo with files
// written over it and committed.
func cspellRun(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"cspell"}, args...)...)
}

// cspellInstalled is where an installed shared package keeps its cspell.json.
const cspellInstalled = "node_modules/@mayjournal/fitness-shared/config/cspell.json"

// docTemplate runs `fitness-install -- doc-template` on happyRepo with files
// written over it.
func docTemplate(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "doc-template")
}

// docTemplateNote guides its one section with a comment.
const docTemplateNote = "# Note\n\n## Summary\n\n<!-- gist -->\n"

// docTemplateADR guides each of its three sections with a comment.
const docTemplateADR = "# NNNN Title\n\n## Context\n\n<!-- why -->\n\n## Decision\n\n<!-- what -->\n\n## Consequences\n\n<!-- fallout -->\n"

const (
	dtViews     = "docs/architecture/"
	dtTemplates = dtViews + "templates/"
)

// dtShipped is a user's docs/architecture after following the
// architecture-contract README: each view copied from its template, templates
// left behind.
func dtShipped(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	for _, view := range []string{acContext, acContainer, acComponent, acCode, acJourney} {
		body, err := os.ReadFile(filepath.Join("..", "internal", "architecture", "templates", view))
		mustDo(t, err)
		files[dtViews+view] = string(body)
	}
	return files
}

// dtCopied is dtShipped with the five templates also copied into
// docs/architecture/templates, as repos did before fitness shipped them.
func dtCopied(t *testing.T) map[string]string {
	t.Helper()
	files := dtShipped(t)
	for _, view := range []string{acContext, acContainer, acComponent, acCode, acJourney} {
		files[dtTemplates+view] = files[dtViews+view]
	}
	return files
}

// dtShippedName is how a message names a view's shipped template.
func dtShippedName(view string) string {
	return "the shipped " + view
}

// dtEdit replaces old with new in one copied view, failing when old is not
// there to replace.
func dtEdit(t *testing.T, files map[string]string, view, old, new string) {
	t.Helper()
	if !strings.Contains(files[dtViews+view], old) {
		t.Fatalf("%s has no %q", view, old)
	}
	files[dtViews+view] = strings.Replace(files[dtViews+view], old, new, 1)
}

// dtLine is the 1-based line of the first line in a view starting with
// prefix, as an error names it.
func dtLine(t *testing.T, files map[string]string, view, prefix string) string {
	t.Helper()
	for i, line := range strings.Split(files[dtViews+view], "\n") {
		if strings.HasPrefix(line, prefix) {
			return dtViews + view + ":" + strconv.Itoa(i+1)
		}
	}
	t.Fatalf("%s has no line starting %q", view, prefix)
	return ""
}

// camelDocsCheck is the check under test.
const camelDocsCheck = "markdown-filename-camel-case"

// camelDocsRepo is a repo whose docs follow camelCase, with files added.
func camelDocsRepo(files map[string]string) map[string]string {
	return with(map[string]string{
		"docs/releaseNotes.md": "# Release notes\n",
		"docs/installGuide.md": "# Install guide\n",
	}, files)
}

// camelDocs runs `fitness-install -- markdown-filename-camel-case` on
// happyRepo with files written over it.
func camelDocs(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, camelDocsCheck)
}

// kebabDocsCheck is the check under test.
const kebabDocsCheck = "markdown-filename-kebab-case"

// kebabDocsRepo is a repo whose docs follow kebab-case, with files added.
func kebabDocsRepo(files map[string]string) map[string]string {
	repo := map[string]string{
		"docs/api-design.md": "# API\n",
		"docs/user-guide.md": "# Guide\n",
	}
	maps.Copy(repo, files)
	return repo
}

// kebabDocs runs `fitness-install -- markdown-filename-kebab-case` on
// happyRepo with files written over it.
func kebabDocs(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, kebabDocsCheck)
}

// markdownLinks runs `fitness-install -- markdown-links` on happyRepo with
// files written over it.
func markdownLinks(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "markdown-links")
}

// markdownLinksReadme is the happy README with body appended; body starts on
// line 9.
func markdownLinksReadme(body string) map[string]string {
	return map[string]string{"README.md": readme + "\n" + body}
}

// noBoldItalic runs `fitness-install -- markdown-no-bold-italic <args>` on
// happyRepo with files written over it.
func noBoldItalic(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"markdown-no-bold-italic"}, args...)...)
}

// noBoldItalicReadme runs the check with text added to happyRepo's README.
func noBoldItalicReadme(t *testing.T, text string) (string, int) {
	t.Helper()
	return noBoldItalic(t, map[string]string{"README.md": readme + "\n" + text})
}

// noBoldItalicRule is the advice every finding carries.
const noBoldItalicRule = " (use only when explicitly required): "

// reframing runs `fitness-install -- no-contrastive-reframing` on happyRepo
// with prose appended to its README.
func reframing(t *testing.T, prose string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", map[string]string{"README.md": readme + "\n" + prose})
	return fitness(t, repo, nil, "no-contrastive-reframing")
}

// reframingHint is the start of every failure message.
const reframingHint = `README.md: rewrite as a direct statement, not the "not X, it's Y" pattern: `

// textReadability runs `fitness-install -- text-readability <args>` on
// happyRepo with files written over it.
func textReadability(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, append([]string{"text-readability"}, args...)...)
}

// textReadabilityPlain is 140 prose words of short, plain sentences.
var textReadabilityPlain = strings.Repeat("The runner builds each check and then runs the full suite on every commit. ", 10)

// textReadabilityDense is 141 long words in one sentence.
var textReadabilityDense = strings.Repeat("extraordinary bureaucratic infrastructure considerations regarding administrative complexity ", 20) + "conclude.\n"

// cspell:ignore antml

// commentStructure runs comment-structure on comment, saved over happyRepo
// with files, as --body-file.
func commentStructure(t *testing.T, comment string, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", with(files, map[string]string{"comment.md": comment})), nil,
		"comment-structure", "--body-file", "comment.md")
}

// words is n words of plain prose.
func words(n int) string {
	return strings.TrimSpace(strings.Repeat("word ", n)) + "."
}

// numbered is a numbered diagram and its matching callout table.
const numbered = "```mermaid\nflowchart TD\n  A[1 Start]:::s --> B[2 End]:::s\n  classDef s fill:#eef\n```\n\n" +
	"| # | Element | Description | Why |\n| - | - | - | - |\n| 1 | Start | First. | Begins. |\n| 2 | End | Last. | Ends. |\n"

// overBudget has 120 prose words, and more inside a table and code.
var overBudget = words(120) + "\n\n| A | B |\n| - | - |\n| " + words(30) + " | x |\n\n```\n" + words(30) + "\n```\n"

// hardToRead holds 104 long words in sentences within their limit.
var hardToRead = strings.Repeat("Configuration interdependencies necessitate comprehensive documentation regarding initialization sequencing, synchronization guarantees, and serialization compatibility. ", 8)

// adrAsIssue runs `fitness-install -- adr-as-issue` on happyRepo with files
// written over it.
func adrAsIssue(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "adr-as-issue")
}

// toADRIssue is the fix every ADR file's error names.
const toADRIssue = "open an `ADR` issue instead; `fitness adr-migrate` opens one per file and points each link at it"

// adrFile is an ADR as a file in an `adr` folder holds it.
const adrFile = "# ADR 0001 — Use SQLite\n\nStatus: accepted\n\n## Context\n\nWe sync offline.\n"

// inADRFolder is the error line for an ADR file in an `adr` folder.
func inADRFolder(path string) string {
	return path + ": a Markdown file in an `adr` folder is an ADR; " + toADRIssue
}

// heldDecision is a check README that still holds its decision.
const heldDecision = "# 0001 — Small files\n\nSmall files read fast.\n\n## Context\n\nLong files hide their point.\n\n" +
	"## Decision\n\n1. Every Markdown file stays small.\n\n## Consequences\n\n1. Readers find the point fast.\n\n" +
	"## Check\n\nFails a long Markdown file.\n"

// heldCheck is ADR 0001's check folder, its README holding the decision.
func heldCheck() map[string]string {
	return map[string]string{"fitness/adr-0001/README.md": heldDecision, "fitness/adr-0001/check": verdictScript(passing)}
}

// allowADRs is a config whose adrAsIssue.allow holds entries.
func allowADRs(entries string) string {
	return `{"adrAsIssue": {"allow": [` + entries + `]}}` + "\n"
}

// cspell:ignore mispeled

// renovateBody is the end of a Renovate PR body, with the hidden debug block
// that holds its state encoded as Renovate encodes it.
var renovateBody = "This PR updates the tooling.\n\n<!--renovate-debug:" +
	base64.StdEncoding.EncodeToString([]byte(`{"createdInVer":"44.222.2","updatedInVer":"44.222.2","targetBranch":"main","labels":[]}`)) +
	"-->\n"

// questionsLabel is a label line that opens a list of open questions, and
// the finding it earns.
const (
	questionsLabel   = "Open questions:\n\n- Which sync wins?\n"
	questionsFinding = `a label line opens a list of open questions: "Open questions:"`
)

// artifactEvent is the PreToolUse event for an Artifact call publishing page.
func artifactEvent(page string) map[string]any {
	return map[string]any{"hook_event_name": "PreToolUse", "tool_name": "Artifact", "tool_input": map[string]any{"file_path": page}}
}
