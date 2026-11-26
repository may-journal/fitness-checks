package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func Test0002_1(t *testing.T) {
	t.Parallel()
	t.Run("0002.1.1.1.1 `prose-budget` fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, acceptance.InComments(strings.TrimSpace(acceptance.LongSentence)))
		acceptance.Sees(t, out, code, 1, acceptance.AtEachFile("a sentence has 30 words (max 23)")...)
		out, code = proseBudget(t, acceptance.InComments("One. Two. Three. Four. Five."))
		acceptance.Sees(t, out, code, 1, acceptance.AtEachFile("a paragraph has 5 sentences (max 4)")...)
	})
	t.Run("0002.1.1.2.1 it fails with the `comment-fragment` rule", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, acceptance.InComments(split...))
		acceptance.Sees(t, out, code, 1, acceptance.AtEachFile(`comment-fragment: a sentence starts with "So"; join it to the sentence before it or drop it`)...)
		out, code = proseBudget(t, map[string]string{"app/fence.ts": "export {};\n\n// - Builds it:\n//   ```\n//   npm run build\n//   ```\n// So it ends.\nexport const y = 2;\n"})
		acceptance.Sees(t, out, code, 1, `app/fence.ts:3: comment-fragment: a sentence starts with "So"`)
	})
	t.Run("0002.1.1.3.1 it fails with the `comment-fragment` rule", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, acceptance.InComments("Retries once on a cold cache.", "Cache warms first."))
		acceptance.Sees(t, out, code, 1, acceptance.AtEachFile("comment-fragment: a sentence has 3 words (min 4)")...)
		out, code = proseBudget(t, acceptance.InComments("See [the doc](../docs/a.md)."))
		acceptance.Sees(t, out, code, 1, acceptance.AtEachFile("comment-fragment: a sentence has 3 words (min 4)")...)
	})
	t.Run("0002.1.1.4.1 it fails with the `comment-one-statement` rule", func(t *testing.T) {
		t.Parallel()
		files := inBodies(threeSentences...)
		files["app/tail.go"] = "package main\n\nvar x = 1 // Holds one on every run. The value never changes. Tests read it directly.\n"
		for path, src := range acceptance.InComments(threeSentences...) {
			files[path] = src
		}
		out, code := proseBudget(t, files)
		want := "comment-one-statement: a comment has 3 sentences (max 2); condense it to one statement"
		acceptance.Sees(t, out, code, 1, append(append(atEachBody(want), acceptance.AtEachFile(want)...), "app/tail.go:3: "+want)...)
	})
	t.Run("0002.1.1.5.1 it passes", func(t *testing.T) {
		t.Parallel()
		files := acceptance.InComments(threeSentences[:2]...)
		for path, src := range inBodies("Retries the request once, since the first call after a deploy warms the cache that every later call reads from.") {
			files[path] = src
		}
		files["app/old.go"] = "package main\n\n// Deprecated: Use main.\nfunc old() {}\n\n// `helper` runs.\nfunc helper() {}\n"
		files["app/svc.ts"] = "export class Svc {\n  // " + strings.Join(threeSentences[:2], "\n  // ") + "\n  run(x: number): number {\n    return x;\n  }\n}\n"
		files["app/list.go"] = "package main\n\n// Codes:\n//\n//   - 1: a check found problems, or the run could not start, such as a bad\n//     config.\nvar codes = 1\n\n// unwraps reads the ? or ! of a?.b, and \"- [ ]\" or \"1. [ ]\" items.\nfunc unwraps() {}\n"
		files["README.md"] = readme + "\nThe cache warms. So it works.\n"
		for path, src := range acceptance.InComments(threeSentences[0]+" "+threeSentences[1], "", "[0008 Other devices](../docs/requirements/0008-other-devices.md)") {
			files[strings.Replace(path, "app/", "app/link/", 1)] = src
		}
		for path, src := range acceptance.InComments(`"1 Oct 2025 · 8:10a".`) {
			files[strings.Replace(path, "app/", "app/quote/", 1)] = src
		}
		for path, src := range acceptance.InComments("“Go to Mon, Sep 28”", "", "[0008][doc]") {
			files[strings.Replace(path, "app/", "app/curly/", 1)] = src
		}
		files["app/run.sh"] = "#!/bin/sh\n# \"Go to Mon, Sep 28, 9:15 AM\"\necho hi\n"
		out, code := proseBudget(t, files)
		acceptance.Sees(t, out, code, 0)
	})
}

// stagedScenario is a file whose old comment breaks a rule, and whose staged
// change adds another. It also holds what each run names.
type stagedScenario struct {
	before, after string
	staged        string   // the finding on the staged comment alone
	left          string   // the old comment's line, which the hook leaves out
	full          []string // the findings of a full scan
}

func stagedScenarios() []stagedScenario {
	long := strings.TrimSpace(acceptance.LongSentence)
	return []stagedScenario{{
		before: "package main\n\n// " + long + "\nfunc old() {}\n\n// Runs.\nfunc main() {}\n",
		after:  "package main\n\n// " + long + "\nfunc old() {}\n\n// Runs " + long + "\nfunc main() {}\n",
		staged: "app/run main.go:6: a sentence has 31 words (max 23)",
		left:   "app/run main.go:3:",
		full:   []string{"app/run main.go:3: a sentence has 30 words (max 23)", "app/run main.go:6: a sentence has 31 words (max 23)"},
	}, {
		before: "package main\n\n// Runs old code.\nfunc old() {}\n\n// Runs the main loop until it ends.\nfunc main() {}\n",
		after:  "package main\n\n// Runs old code.\nfunc old() {}\n\n// Runs the main loop. So it ends.\nfunc main() {}\n",
		staged: `app/run main.go:6: comment-fragment: a sentence starts with "So"`,
		left:   "app/run main.go:3:",
		full:   []string{"app/run main.go:3: comment-fragment: a sentence has 3 words (min 4)", `app/run main.go:6: comment-fragment: a sentence starts with "So"`},
	}, {
		before: "package main\n\n// Basically runs the old code.\nfunc old() {}\n\n// Runs the main loop until it ends.\nfunc main() {}\n",
		after:  "package main\n\n// Basically runs the old code.\nfunc old() {}\n\n// Runs the main loop, simply until it ends.\nfunc main() {}\n",
		staged: `app/run main.go:6: hedging: "simply" hedges`,
		left:   "app/run main.go:3:",
		full:   []string{`app/run main.go:3: hedging: "Basically" hedges`, `app/run main.go:6: hedging: "simply" hedges`},
	}}
}

// stage commits s.before as app/run main.go, then stages s.after over it.
func (s stagedScenario) stage(t *testing.T) string {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{"app/run main.go": s.before})
	acceptance.Write(t, repo, map[string]string{"app/run main.go": s.after})
	acceptance.Git(t, repo, "add", "-A")
	return repo
}

func Test0002_2(t *testing.T) {
	t.Parallel()
	t.Run("0002.2.1.1.1 only that comment fails", func(t *testing.T) {
		t.Parallel()
		for _, s := range stagedScenarios() {
			out, code := acceptance.Fitness(t, s.stage(t), nil, "prose-budget")
			acceptance.Sees(t, out, code, 1, s.staged)
			if strings.Contains(out, s.left) {
				t.Errorf("the hook judged a comment the change left alone:\n%s", out)
			}
		}
	})
	t.Run("0002.2.1.1.2 a full scan fails both comments", func(t *testing.T) {
		t.Parallel()
		for _, s := range stagedScenarios() {
			out, code := acceptance.Fitness(t, s.stage(t), nil, "--all", "prose-budget")
			acceptance.Sees(t, out, code, 1, s.full...)
		}
	})
}

func Test0002_3(t *testing.T) {
	t.Parallel()
	t.Run("0002.3.1.1.1 it fails with the `commented-out-code` rule", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, inFunctions(
			"\t// x := compute(a)\n\t// if x > 0 {\n\t// \treturn x\n\t// }\n",
			"    // let x = compute(a)\n    // print(x)\n",
			"  // const x = compute(a);\n  // console.log(x);\n"))
		acceptance.Sees(t, out, code, 1,
			"app/run.go:4: commented-out-code: 3 of 3 lines are code; delete it, since git keeps old code",
			"app/Run.swift:4: commented-out-code: 2 of 2 lines are code; delete it",
			"app/run.ts:4: commented-out-code: 2 of 2 lines are code; delete it")
	})
	t.Run("0002.3.1.2.1 it fails with the `restates-code` rule", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, inFunctions(
			"\ti := 0\n\t// increment i\n\ti++\n",
			"    var i = 0\n    // increment i\n    i += 1\n",
			"  let i = 0;\n  // increment i\n  i++;\n"))
		acceptance.Sees(t, out, code, 1,
			"app/run.go:5: restates-code: every word repeats the next code line; delete it or say why",
			"app/Run.swift:5: restates-code: every word repeats the next code line",
			"app/run.ts:5: restates-code: every word repeats the next code line")
	})
	t.Run("0002.3.1.3.1 it fails with the `decoration` rule", func(t *testing.T) {
		t.Parallel()
		files := acceptance.InComments("🚀 Launches the rocket once on start.")
		files["app/banner.go"] = "package main\n\n// ==========\nvar x = 1\n"
		out, code := proseBudget(t, files)
		acceptance.Sees(t, out, code, 1, append(acceptance.AtEachFile(`decoration: a comment holds the emoji "🚀"; delete it`),
			`app/banner.go:3: decoration: a comment holds the banner "=========="; delete it`)...)
	})
	t.Run("0002.3.1.4.1 it fails with the `change-narration` rule", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, acceptance.InComments("Previously read the cache on every run."))
		acceptance.Sees(t, out, code, 1, acceptance.AtEachFile(`change-narration: "Previously" narrates a change; say what holds today, since git keeps history`)...)
	})
	t.Run("0002.3.1.5.1 it fails with the `hedging` rule", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, acceptance.InComments("Basically retries once on a cold cache."))
		acceptance.Sees(t, out, code, 1, acceptance.AtEachFile(`hedging: "Basically" hedges; drop the word or state the fact`)...)
	})
	t.Run("0002.3.1.6.1 it fails with the `tutorial-tone` rule", func(t *testing.T) {
		t.Parallel()
		files := acceptance.InComments("Note that the cache warms once on start.")
		files["app/plan.go"] = "package main\n\n// Retries once on a miss. Let's keep the cache warm.\nfunc plan() {}\n"
		files["app/Seen.swift"] = "import Foundation\n\n/// As you can see, the cache warms on start.\nfunc seen() {}\n"
		files["app/steps.ts"] = "export {};\n\n// First we read the cache, then the disk.\nexport const steps = 2;\n"
		out, code := proseBudget(t, files)
		acceptance.Sees(t, out, code, 1, append(acceptance.AtEachFile(`tutorial-tone: "Note that" opens a sentence like a tutorial; state the rule plainly`),
			`app/plan.go:3: tutorial-tone: "Let's" opens a sentence`,
			`app/Seen.swift:3: tutorial-tone: "As you can see" opens a sentence`,
			`app/steps.ts:3: tutorial-tone: "First we" opens a sentence`)...)
	})
}

func Test0002_4(t *testing.T) {
	t.Parallel()
	t.Run("0002.4.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		files := acceptance.InComments("Retries just once, since the first call warms the cache.")
		files["app/state.go"] = "package main\n\n// Maps each state → its next state on every tick.\nvar next = map[int]int{}\n"
		files["app/expiry.ts"] = "export {};\n\n// The cache no longer expires once it holds a value.\nexport const keep = true;\n"
		files["app/store.go"] = "package main\n\n// Returns the previously stored value.\nfunc stored() int { return 0 }\n"
		files["README.md"] = readme + "\nIt basically worked previously. 🚀\n\n-----\n"
		out, code := proseBudget(t, files)
		acceptance.Sees(t, out, code, 0)
	})
	t.Run("0002.4.1.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		files := acceptance.InComments("Writes a note that names the cache.")
		files["app/plan.go"] = "package main\n\n// Lets the caller retry once on a miss; note that empty keys fail.\nfunc plan() {}\n"
		files["app/greet.ts"] = "export {};\n\n// Prints \"Let's go\" when the run starts.\nexport const greet = 1;\n"
		files["README.md"] = readme + "\nNote that it works. Let's see.\n"
		out, code := proseBudget(t, files)
		acceptance.Sees(t, out, code, 0)
	})
}

func Test0002_5(t *testing.T) {
	t.Parallel()
	t.Run("0002.5.1.1.1 it fails with the `overlong` rule", func(t *testing.T) {
		t.Parallel()
		out, code := proseBudget(t, inFunctions(
			"\t// Hands the number back to the caller without any change at all.\n\n\t// Callers rely on getting their own number back from this function.\n",
			"    // Hands the number back to the caller without any change at all.\n\n    // Callers rely on getting their own number back from this function.\n",
			"  // Hands the number back to the caller without any change at all.\n\n  // Callers rely on getting their own number back from this function.\n"))
		acceptance.Sees(t, out, code, 1,
			"app/run.go:3: overlong: a function body holds 23 comment words in 2 comments for 2 code tokens (max 3 words per token); keep the comments that say why and cut the rest",
			"app/Run.swift:3: overlong: a function body holds 23 comment words in 2 comments for 2 code tokens",
			"app/run.ts:3: overlong: a function body holds 23 comment words in 2 comments for 3 code tokens")
	})
	t.Run("0002.5.1.2.1 it passes", func(t *testing.T) {
		t.Parallel()
		files := inFunctions(
			"\t// Hands the number back unchanged, since callers compare it with the value they passed in.\n",
			"    var b = a\n    // Only the copy changes below this line.\n    b += 1\n    // The copy is dropped once the caller has its number.\n    _ = b\n",
			"  let b = a;\n  // Only the copy changes below this line.\n  b++;\n  // The copy is dropped once the caller has its number.\n  void b;\n")
		files["app/swap.go"] = "package main\n\nfunc swap(a, b int) (int, int) {\n\t// Swaps the pair, as this call shows:\n" +
			"\t//     first, second := swap(the value the caller passed first, the value the caller passed second)\n" +
			"\t//     first, second = swap(the value that came back first, the value that came back second)\n" +
			"\tt := a\n\n\t// The copy keeps a safe while b moves.\n\ta = b\n\treturn a, t\n}\n"
		files["app/greet.sh"] = "#!/bin/sh\ngreet() {\n# echo hi\necho hi\n}\n"
		files["README.md"] = readme + "\nAn over-commented example:\n\n```go\nfunc run(a int) int {\n" +
			"\t// Hands the number back to the caller without any change at all.\n\n" +
			"\t// Callers rely on getting their own number back from this function.\n\treturn a\n}\n```\n"
		out, code := proseBudget(t, files)
		acceptance.Sees(t, out, code, 0)
	})
}

// inBodies writes a comment of lines on line 4, inside a Go, a Swift and a TypeScript function.
func inBodies(lines ...string) map[string]string {
	mark := func(indent, before, after string) string {
		var b strings.Builder
		b.WriteString(before)
		for _, l := range lines {
			b.WriteString(indent + "// " + l + "\n")
		}
		return b.String() + after
	}
	return map[string]string{
		"app/body.go":    mark("\t", "package main\n\nfunc body() {\n", "\tprintln()\n}\n"),
		"app/Body.swift": mark("    ", "import Foundation\n\nfunc body() {\n", "    print(1)\n}\n"),
		"app/body.ts":    mark("  ", "export {};\n\nexport function body(): void {\n", "  console.log(1);\n}\n"),
	}
}

// atEachBody names line 4 of each inBodies file, followed by message.
func atEachBody(message string) []string {
	return []string{"app/Body.swift:4: " + message, "app/body.go:4: " + message, "app/body.ts:4: " + message}
}

// split is the sentence and fragment PR #746 passed with.
var split = []string{
	"run reads the whole model and every requirement doc unless the run is architecture.Unmoved.",
	"So a removed actor fails even when no requirement changed.",
}

// threeSentences is a comment of three whole sentences.
var threeSentences = []string{
	"Reads the cache first on every run.",
	"The cache may be cold after a deploy.",
	"One retry warms it for the next call.",
}

// inFunctions writes each Go, Swift and TypeScript body, its comment from line 4, inside a function.
func inFunctions(goBody, swiftBody, tsBody string) map[string]string {
	return map[string]string{
		"app/run.go":    "package main\n\nfunc run(a int) int {\n" + goBody + "\treturn a\n}\n",
		"app/Run.swift": "import Foundation\n\nfunc run(a: Int) -> Int {\n" + swiftBody + "    return a\n}\n",
		"app/run.ts":    "export {};\n\nexport function run(a: number): number {\n" + tsBody + "  return a;\n}\n",
	}
}
