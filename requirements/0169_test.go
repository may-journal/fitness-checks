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
