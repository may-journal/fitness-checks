package requirements

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// changedClone clones this repo, appends a line to each file, commits, and
// returns the clone with the commit before the change.
func changedClone(t *testing.T, files ...string) (string, string) {
	t.Helper()
	clone := filepath.Join(t.TempDir(), "repo")
	if out, err := exec.Command("git", "clone", "-q", "--local", "..", clone).CombinedOutput(); err != nil {
		t.Fatalf("git clone: %v\n%s", err, out)
	}
	base, err := exec.Command("git", "-C", clone, "rev-parse", "HEAD").Output()
	mustDo(t, err)
	for _, file := range files {
		f, err := os.OpenFile(filepath.Join(clone, file), os.O_APPEND|os.O_WRONLY, 0)
		mustDo(t, err)
		_, err = f.WriteString("\n// changed\n")
		mustDo(t, errors.Join(err, f.Close()))
	}
	git(t, clone, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-qam", "test(change): change files")
	return clone, string(base[:len(base)-1])
}

// changedChecks runs this checkout's scripts/changed-checks.sh in clone with
// args and decodes the checks it names.
func changedChecks(t *testing.T, clone string, args ...string) []string {
	t.Helper()
	script, err := filepath.Abs(filepath.Join("..", "scripts", "changed-checks.sh"))
	mustDo(t, err)
	cmd := exec.Command(script, args...)
	cmd.Dir = clone
	out, err := cmd.Output()
	mustDo(t, err)
	var checks []string
	mustDo(t, json.Unmarshal(out, &checks))
	return checks
}

func Test0161_1(t *testing.T) {
	t.Parallel()
	clone, base := changedClone(t, "cmd/fitness-check-gofmt/main.go")
	if got := changedChecks(t, clone, base); !slices.Equal(got, []string{"gofmt"}) {
		t.Errorf("changed-checks = %v, want only gofmt", got)
	}
}

func Test0161_2(t *testing.T) {
	t.Parallel()
	clone, base := changedClone(t, "internal/spell/mask.go")
	if got := changedChecks(t, clone, base); !slices.Equal(got, []string{"cspell"}) {
		t.Errorf("changed-checks = %v, want only cspell, the one check that imports spell", got)
	}
}

func Test0161_3(t *testing.T) {
	t.Parallel()
	clone, base := changedClone(t, "go.mod")
	every := checkFolders0160(t)
	for _, got := range [][]string{changedChecks(t, clone, base), changedChecks(t, clone)} {
		if len(got) != len(every) {
			t.Errorf("changed-checks names %d checks, want all %d", len(got), len(every))
		}
	}
}

func Test0161_4(t *testing.T) {
	t.Parallel()
	script, err := filepath.Abs(filepath.Join("..", "scripts", "test-check.sh"))
	mustDo(t, err)
	out, err := exec.Command(script, "gofmt").CombinedOutput()
	if err != nil || !strings.Contains(string(out), `{"name":"gofmt"`) {
		t.Errorf("test-check.sh gofmt: %v\n%s", err, out)
	}
}
