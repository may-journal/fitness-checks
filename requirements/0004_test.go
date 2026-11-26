package requirements

import (
	"os"
	"path/filepath"
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

func Test0004_8(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
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

// hardToRead holds 104 long words in sentences within their limit.
var hardToRead = strings.Repeat("Configuration interdependencies necessitate comprehensive documentation regarding initialization sequencing, synchronization guarantees, and serialization compatibility. ", 8)

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
