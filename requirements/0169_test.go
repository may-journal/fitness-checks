package requirements

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

// described0169 holds each check program's --describe metadata, built once.
var described0169 = sync.OnceValues(func() (map[string]checkkit.Describe, error) {
	dir, err := os.MkdirTemp("", "checks-0169-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	build := exec.Command("go", "build", "-o", dir+string(filepath.Separator), "./cmd/...")
	build.Dir = ".."
	if out, err := build.CombinedOutput(); err != nil {
		return nil, errors.New("go build ./cmd/...: " + string(out))
	}
	return describeAll0169(dir)
})

// describeAll0169 runs --describe on each check program in dir, by name.
func describeAll0169(dir string) (map[string]checkkit.Describe, error) {
	entries, err := os.ReadDir(dir)
	all := map[string]checkkit.Describe{}
	for _, e := range entries {
		if name, ok := strings.CutPrefix(e.Name(), checkkit.BinPrefix); ok && err == nil {
			all[name], err = describe0169(filepath.Join(dir, e.Name()))
		}
	}
	return all, err
}

// describe0169 decodes the --describe line program prints.
func describe0169(program string) (checkkit.Describe, error) {
	var d checkkit.Describe
	out, err := exec.Command(program, "--describe").Output()
	if err == nil {
		err = json.Unmarshal(out, &d)
	}
	return d, err
}

// patterns0169 returns check's compiled patterns, failing the test when it
// is missing or its patterns do not compile.
func patterns0169(t *testing.T, check string) checkkit.Patterns {
	t.Helper()
	all, err := described0169()
	mustDo(t, err)
	d, ok := all[check]
	if !ok {
		t.Fatalf("no check program describes %s", check)
	}
	p, err := checkkit.PatternsOf(d)
	mustDo(t, err)
	return p
}

// reaches0169 asserts whether check's patterns reach path.
func reaches0169(t *testing.T, check, path string, want bool) {
	t.Helper()
	if got := patterns0169(t, check).Matches(path); got != want {
		t.Errorf("%s reaches %s = %v, want %v", check, path, got, want)
	}
}

func Test0169_1(t *testing.T) {
	t.Parallel()
	all, err := described0169()
	mustDo(t, err)
	if len(all) == 0 {
		t.Fatal("the build made no check programs")
	}
	for name, d := range all {
		if _, err := checkkit.PatternsOf(d); err != nil {
			t.Errorf("%s: %v", name, err)
		}
		if r := checkkit.UnknownRequirement(d); r != "" {
			t.Errorf("%s requires %q, which fitness-runner does not define", name, r)
		}
	}
}

func Test0169_2(t *testing.T) {
	t.Parallel()
	reaches0169(t, "gofmt", "pkg/run.go", true)
	reaches0169(t, "gofmt", "vendor/x/run.go", false)
	reaches0169(t, "gofmt", "README.md", false)
	reaches0169(t, "markdown-links", "docs/guide.md", true)
	reaches0169(t, "markdown-links", "pkg/run.go", false)
	reaches0169(t, "workflow-reporting", ".github/workflows/ci.yml", true)
	reaches0169(t, "workflow-reporting", "docs/ci.yml", false)
	eslint := patterns0169(t, "eslint")
	if !eslint.Met([]string{"web/package.json"}) || eslint.Met([]string{"README.md"}) {
		t.Errorf("eslint must need a package.json and nothing else")
	}
}

func Test0169_3(t *testing.T) {
	t.Parallel()
	for check, path := range map[string]string{
		"pr-structure":       checkkit.PullBodyPath,
		"pr-closes-issue":    checkkit.PullBodyPath,
		"plan-structure":     checkkit.IssueBodyPath("Plan"),
		"epic-structure":     checkkit.IssueBodyPath("Epic"),
		"issue-checklist":    checkkit.IssueBodyPath("Plan"),
		"comment-structure":  checkkit.CommentBodyPath,
		"semantic-commit":    checkkit.CommitMessagePath,
		"plan-trailer":       checkkit.CommitMessagePath,
		"commit-attribution": checkkit.CommitMessagePath,
		"issue-link-once":    checkkit.CommitMessagePath,
	} {
		reaches0169(t, check, path, true)
		reaches0169(t, check, "pkg/run.go", false)
		reaches0169(t, check, checkkit.BlobPath, !strings.HasPrefix(path, ".git/"))
		if all, err := described0169(); err == nil && all[check].MinFitness != "1.75.0" {
			t.Errorf("%s needs fitness %q, want 1.75.0, the first with virtual paths", check, all[check].MinFitness)
		}
	}
	reaches0169(t, "issue-checklist", checkkit.IssueBodyPath(""), true)
	reaches0169(t, "pr-structure", checkkit.IssueBodyPath("Plan"), false)
	reaches0169(t, "plan-structure", checkkit.IssueBodyPath("Epic"), false)
	reaches0169(t, "semantic-commit", checkkit.HistoryPath, true)
	reaches0169(t, "semantic-commit", checkkit.PullBodyPath, false)
	reaches0169(t, "comment-structure", checkkit.PullBodyPath, false)
}

func Test0169_4(t *testing.T) {
	t.Parallel()
	for _, check := range []string{"prose-budget", "text-readability", "markdown-no-bold-italic", "mermaid-syntax", "open-questions"} {
		for _, path := range []string{"docs/guide.md", checkkit.PullBodyPath, checkkit.CommentBodyPath, checkkit.BlobPath} {
			reaches0169(t, check, path, true)
		}
		reaches0169(t, check, "web/app.css", false)
	}
	for _, check := range []string{"prose-budget", "text-readability", "no-contrastive-reframing"} {
		reaches0169(t, check, "pkg/run.go", true)
		reaches0169(t, check, "scripts/deploy", true)
		reaches0169(t, check, checkkit.CommitMessagePath, false)
	}
	reaches0169(t, "no-contrastive-reframing", checkkit.PullBodyPath, false)
}

// patternsExit0169 pipes line to the describe-patterns program and returns
// what it printed and its exit code.
func patternsExit0169(t *testing.T, program, line string) (string, int) {
	t.Helper()
	cmd := exec.Command(program, "fitness-check-x-v1.0.0")
	cmd.Stdin = strings.NewReader(line)
	out, err := cmd.CombinedOutput()
	if exit, ok := err.(*exec.ExitError); ok {
		return string(out), exit.ExitCode()
	}
	mustDo(t, err)
	return string(out), 0
}

func Test0169_5(t *testing.T) {
	t.Parallel()
	program := filepath.Join(t.TempDir(), "describe-patterns")
	build := exec.Command("go", "build", "-o", program, "./scripts/describe-patterns")
	build.Dir = ".."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	out, code := patternsExit0169(t, program, `{"name":"broken","files":"(oops"}`)
	if code != 1 || !strings.Contains(out, `fitness-check-x-v1.0.0: files pattern "(oops"`) {
		t.Errorf("a broken files pattern exited %d: %s", code, out)
	}
	all, err := described0169()
	mustDo(t, err)
	gofmt, err := json.Marshal(all["gofmt"])
	mustDo(t, err)
	if out, code := patternsExit0169(t, program, string(gofmt)); code != 0 {
		t.Errorf("gofmt's describe line exited %d: %s", code, out)
	}
}
