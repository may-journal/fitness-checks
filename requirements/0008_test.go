package requirements

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
	"The check's place in the system is in the [architecture](../../docs/architecture/03-components.md).\n\n" +
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
			if !strings.Contains(doc, "//go:embed README.md\nvar docSource string") || strings.Contains(doc, "Behavior") {
				t.Errorf("doc.go does not just embed the README:\n%s", doc)
			}
		}},
		{"the described summary stays", func(t *testing.T) {
			t.Parallel()
			folder := demoFolder(t, demoDoc)
			want := program.Summary(demoDoc[strings.Index(demoDoc, "// Command"):])
			if _, err := moveDocs(t, folder); err != nil {
				t.Fatal(err)
			}
			for _, file := range []string{"README.md", "doc.go"} {
				if got := program.Summary(readText(t, folder, file)); got != want || want != "Fails a Zork file." {
					t.Errorf("%s summary is %q, want %q", file, got, "Fails a Zork file.")
				}
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
	"architecture-components", "architecture-containers", "build-output-untracked",
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
		{"a moved check's doc.go is one comment that embeds its README", func(t *testing.T) {
			t.Parallel()
			for _, name := range movedChecks {
				docGo := readText(t, filepath.Join("..", "cmd", "fitness-check-"+name), "doc.go")
				if problem := singleCommentProblem(name, docGo); problem != "" {
					t.Errorf("%s doc.go %s:\n%s", name, problem, docGo)
				}
			}
		}},
	})
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

// singleCommentProblem says how docGo of the check name is not one comment
// followed by the README embed; it is "" when it is. The comment opens with
// "Command fitness-check-<name>".
func singleCommentProblem(name, docGo string) string {
	comment, tail, ok := strings.Cut(docGo, newTail)
	if !ok || tail != "" {
		return "holds more than the comment and the README embed"
	}
	lines := strings.Split(strings.TrimSuffix(comment, "\n"), "\n")
	if !strings.HasPrefix(lines[0], "// Command fitness-check-"+name+" ") {
		return "does not open with its Command comment"
	}
	return commentLinesProblem(lines)
}

// commentLinesProblem says how lines break the "// " form of 78 columns
// and a closing pointer to the README; it is "" when they hold.
func commentLinesProblem(lines []string) string {
	if slices.ContainsFunc(lines, notCommentLine) {
		return "comment is not one run of lines of at most 78 columns"
	}
	if !strings.HasSuffix(strings.ReplaceAll(strings.Join(lines, " "), " // ", " "), " Its rule is in README.md.") {
		return "comment does not end by pointing at the README"
	}
	return ""
}

// notCommentLine reports whether line is not a "// " comment line of at most
// 78 columns.
func notCommentLine(line string) bool {
	return !strings.HasPrefix(line, "// ") || len(line) > 78
}

// newTail is what a moved check's doc.go holds after its comment.
const newTail = "package main\n\nimport _ \"embed\"\n\n" +
	"// docSource is README.md, whose first paragraph --describe summarizes.\n//\n//go:embed README.md\nvar docSource string\n"
