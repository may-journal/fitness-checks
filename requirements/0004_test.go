package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test0004_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
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
