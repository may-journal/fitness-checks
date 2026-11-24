package requirements

import (
	"os"
	"os/exec"
	"path/filepath"
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
	"architecture-components", "architecture-containers", "build-output-untracked", "changelog",
	"changelog-bullets", "changelog-updated", "checkstyle", "comment-structure", "commit-attribution",
	"cspell", "dependency-currency", "doc-template", "duplicate-code", "epic-structure", "eslint",
	"gitignore-why", "go-complexity", "go-test", "gofmt",
}

func Test0008_3(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"a moved check's doc.go is one comment that embeds its README", func(t *testing.T) {
			t.Parallel()
			for _, name := range movedChecks {
				docGo := readText(t, filepath.Join("..", "cmd", "fitness-check-"+name), "doc.go")
				comment, _, _ := strings.Cut(docGo, "package main")
				if !readmeIsDocs(docGo) || strings.Count(comment, "\n") > 3 {
					t.Errorf("%s doc.go holds more than a short comment and the README embed:\n%s", name, docGo)
				}
			}
		}},
	})
}
