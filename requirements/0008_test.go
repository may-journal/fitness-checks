package requirements

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/program"
)

// demoDoc is a doc.go that holds every kind of block the script converts.
const demoDoc = "// cspell:ignore zork\n\n" +
	"// Command fitness-check-demo fails a [Zork] file. It runs on every repo.\n//\n" +
	"// # Behavior\n//\n" +
	"//   - Pass: no file.\n" +
	"//   - Fail: a `.zork` file; see [the docs].\n//\n" +
	"// Run it as:\n//\n//\tfitness demo --all\n//\n" +
	"// [Zork]: https://example.com/zork\n" +
	"// [the docs]: https://example.com/docs\n" +
	"package main\n\nimport _ \"embed\"\n\n" +
	"// docSource is this file, whose comment --describe summarizes.\n//\n//go:embed doc.go\nvar docSource string\n"

// demoReadme is the README the script writes for demoDoc.
const demoReadme = "# demo\n\n" +
	"Fails a [Zork](https://example.com/zork) file. It runs on every repo.\n\n" +
	"## Behavior\n\n" +
	"- Pass: no file.\n" +
	"- Fail: a `.zork` file; see [the docs](https://example.com/docs).\n\n" +
	"Run it as:\n\n```\nfitness demo --all\n```\n\n" +
	"<!-- cspell:ignore zork -->\n"

func Test0008_1(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"the comment becomes a README and doc.go embeds it", func(t *testing.T) {
			t.Parallel()
			folder := demoFolder(t, demoDoc)
			if out, err := moveDocs(t, folder); err != nil {
				t.Fatalf("doc-to-readme failed: %v\n%s", err, out)
			}
			if got := readText(t, folder, "README.md"); got != demoReadme {
				t.Errorf("README.md is\n%s\nwant\n%s", got, demoReadme)
			}
			doc := readText(t, folder, "doc.go")
			if !isPointer(doc, "demo") || strings.Contains(doc, "Behavior") {
				t.Errorf("doc.go does not just embed the README:\n%s", doc)
			}
		}},
		{"the described summary stays", func(t *testing.T) {
			t.Parallel()
			folder := demoFolder(t, demoDoc)
			const want = "Fails a Zork file."
			if _, err := moveDocs(t, folder); err != nil {
				t.Fatal(err)
			}
			if got := program.Summary(readText(t, folder, "README.md")); got != want {
				t.Errorf("README.md summary is %q, want %q", got, want)
			}
		}},
	})
}

func Test0008_2(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"a doc.go with more than its comment fails and stays", func(t *testing.T) {
			t.Parallel()
			extra := demoDoc + "\nvar leftover = 1\n"
			folder := demoFolder(t, extra)
			out, err := moveDocs(t, folder)
			if err == nil || !strings.Contains(string(out), "fitness-check-demo") {
				t.Errorf("doc-to-readme gave %v and %q, want a failure naming the folder", err, out)
			}
			if readText(t, folder, "doc.go") != extra {
				t.Error("a refused doc.go was rewritten")
			}
			if _, statErr := os.Stat(filepath.Join(folder, "README.md")); statErr == nil {
				t.Error("a refused move wrote a README.md")
			}
		}},
	})
}

// pointerTail is what doc.go holds after its one-line comment.
const pointerTail = "\nimport _ \"embed\"\n\n" +
	"// docSource is README.md, whose first paragraph --describe summarizes.\n//\n//go:embed README.md\nvar docSource string\n"

// isPointer reports whether docGo is the one-line comment that names the README
// of check name, then only the embed of that README.
func isPointer(docGo, name string) bool {
	want := "// Command fitness-check-" + name + " is documented in README.md.\npackage main\n"
	return docGo == want+pointerTail
}

// demoFolder writes a check folder named fitness-check-demo holding docGo, and
// returns it.
func demoFolder(t *testing.T, docGo string) string {
	t.Helper()
	folder := filepath.Join(t.TempDir(), "fitness-check-demo")
	mustDo(t, os.MkdirAll(folder, 0o755))
	mustDo(t, os.WriteFile(filepath.Join(folder, "doc.go"), []byte(docGo), 0o644))
	return folder
}

// moveDocs runs the doc move script on folder, with its combined output.
func moveDocs(t *testing.T, folder string) ([]byte, error) {
	t.Helper()
	return exec.Command(filepath.Join(goBuild(t, "./scripts/doc-to-readme"), "doc-to-readme"), folder).CombinedOutput()
}

// readText reads file in folder, failing t when it cannot.
func readText(t *testing.T, folder, file string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(folder, file))
	mustDo(t, err)
	return string(data)
}

// movedChecks are the checks whose rule is their README, which each batch of
// the move to Markdown grows.
var movedChecks = []string{
	"adr-as-issue", "adr-checks", "architecture-actors", "architecture-code",
	"architecture-components", "architecture-containers", "architecture-contract",
	"architecture-journey", "build-output-untracked",
	"changelog", "changelog-bullets", "changelog-updated", "checkstyle", "comment-structure",
	"commit-attribution", "cspell", "dependency-currency", "doc-template", "duplicate-code",
	"epic-structure", "eslint", "gitignore-why", "go-complexity", "go-test",
	"go-test-coverage", "go-vet", "gofmt", "issue-checklist", "issue-link-once", "java-test",
	"java-version", "local-branch", "localized-text", "markdown-filename-camel-case",
	"markdown-filename-kebab-case", "markdown-links", "markdown-no-bold-italic",
	"mermaid-callout-why", "mermaid-callouts", "mermaid-diagram-prose",
	"mermaid-diagram-table-gap", "mermaid-legend", "mermaid-level-bleed", "mermaid-syntax",
	"no-contrastive-reframing", "no-eslint-disable", "no-fitness-front-matter",
	"no-plans-dir", "node-version", "open-questions", "plan-structure", "plan-trailer",
	"pr-closes-issue", "pr-structure", "prettier", "prose-budget", "read-repo-first",
	"release-changelog", "repeated-string-literals", "requirements", "requirements-changed",
	"semantic-commit", "spotless", "swift-storybook-requirements", "swiftlint",
	"text-readability", "vitest-coverage-exclude", "vitest-coverage-full",
	"workflow-ref-exists", "workflow-reporting", "workflow-scripts",
}

func Test0008_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"a moved check's doc.go is one line that names its README and embeds it", func(t *testing.T) {
			t.Parallel()
			for _, name := range movedChecks {
				docGo := readText(t, filepath.Join("..", "cmd", "fitness-check-"+name), "doc.go")
				if !isPointer(docGo, name) {
					t.Errorf("%s doc.go is not the one-line pointer and the README embed:\n%s", name, docGo)
				}
			}
		}},
		{"no check README repeats the architecture sentence and the repo README links the page once", func(t *testing.T) {
			t.Parallel()
			const page = "docs/architecture/03-components.md"
			for _, name := range movedChecks {
				if strings.Contains(readText(t, filepath.Join("..", "cmd", "fitness-check-"+name), "README.md"), "03-components.md") {
					t.Errorf("%s README links the architecture page", name)
				}
			}
			if got := strings.Count(readText(t, "..", "README.md"), "]("+page+")"); got != 1 {
				t.Errorf("the repo README links %s %d times, want once", page, got)
			}
		}},
	})
}

func Test0008_6(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"a README names no config key its check does not declare", func(t *testing.T) {
			t.Parallel()
			declared := declaredKeys(t)
			tops := topKeys(declared)
			for name, keys := range declared {
				readme := readText(t, filepath.Join("..", "cmd", "fitness-check-"+name), "README.md")
				for _, key := range undeclaredKeys(readme, keys, tops) {
					t.Errorf("%s README names the key %s, which its --describe does not declare", name, key)
				}
			}
		}},
		{"a README whose check declares keys points to fitness explain", func(t *testing.T) {
			t.Parallel()
			for name, keys := range declaredKeys(t) {
				readme := readText(t, filepath.Join("..", "cmd", "fitness-check-"+name), "README.md")
				if len(keys) > 0 && !strings.Contains(readme, "`fitness explain "+name+"`") {
					t.Errorf("%s README does not say to run `fitness explain %s`", name, name)
				}
			}
		}},
	})
}

// declaredKeys maps each check's name to the .fitnessrc.json keys its
// --describe declares, nested keys included.
func declaredKeys(t *testing.T) map[string]map[string]bool {
	t.Helper()
	bin := goBuild(t, "./cmd/...")
	all := map[string]map[string]bool{}
	for _, folder := range checkFolders(t, "..") {
		name := strings.TrimPrefix(filepath.Base(folder), "fitness-check-")
		keys := map[string]bool{}
		collectKeys(describeMap(t, filepath.Join(bin, filepath.Base(folder)))["config"], keys)
		all[name] = keys
	}
	return all
}

// collectKeys adds each key in the decoded --describe config list, and the
// keys under it, to into.
func collectKeys(list any, into map[string]bool) {
	items, _ := list.([]any)
	for _, item := range items {
		entry, _ := item.(map[string]any)
		if key, _ := entry["key"].(string); key != "" {
			into[key] = true
		}
		collectKeys(entry["keys"], into)
	}
}

// topKeys is each top-level key, such as proseBudget, that any check declares.
func topKeys(declared map[string]map[string]bool) map[string]bool {
	tops := map[string]bool{}
	for _, keys := range declared {
		for key := range keys {
			tops[strings.Split(key, ".")[0]] = true
		}
	}
	return tops
}

// codeSpan matches a Markdown inline code span.
var codeSpan = regexp.MustCompile("`([^`\n]+)`")

// undeclaredKeys lists the code spans in readme that name a top-level key or
// a key under one that own does not declare.
func undeclaredKeys(readme string, own, tops map[string]bool) []string {
	var bad []string
	for _, m := range codeSpan.FindAllStringSubmatch(readme, -1) {
		span := m[1]
		if tops[strings.Split(span, ".")[0]] && !own[span] {
			bad = append(bad, span)
		}
	}
	return bad
}
func Test0008_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"a link whose URL holds parentheses leaves no stray parenthesis", func(t *testing.T) {
			t.Parallel()
			got := program.Summary("# demo\n\nFails a [Zork](https://example.com/Zork_(game)) file. It runs on every repo.\n")
			if want := "Fails a Zork file."; got != want {
				t.Errorf("summary is %q, want %q", got, want)
			}
		}},
		{"a first block that is not prose gives no summary", func(t *testing.T) {
			t.Parallel()
			for name, block := range map[string]string{
				"list": "- Fails a file.", "numbered list": "1. Fails a file.", "fence": "```\nfitness demo\n```",
				"quote": "> Fails a file.", "table": "| a | b |", "comment": "<!-- Fails a file. -->",
			} {
				if got := program.Summary("# demo\n\n" + block + "\n\nFails a file.\n"); got != "" {
					t.Errorf("a README opening with a %s has summary %q, want none", name, got)
				}
			}
		}},
	})
}

func Test0008_5(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"a comment opening with a multi-byte letter keeps it whole", func(t *testing.T) {
			t.Parallel()
			folder := demoFolder(t, strings.Replace(demoDoc, "fails a [Zork] file", "élan a [Zork] file", 1))
			if out, err := moveDocs(t, folder); err != nil {
				t.Fatalf("doc-to-readme failed: %v\n%s", err, out)
			}
			if got := readText(t, folder, "README.md"); !strings.Contains(got, "\n\nÉlan a [Zork]") {
				t.Errorf("README.md does not open with a capital É:\n%s", got)
			}
		}},
	})
}
