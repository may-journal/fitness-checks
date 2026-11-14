// cspell:ignore modcacherw

package requirements

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// readmeSentence0165 is the sentence the tests add to gofmt's doc comment.
const readmeSentence0165 = "It also reads well."

// readmeClone copies this checkout, without its build output, and returns
// the copy.
func readmeClone(t *testing.T) string {
	t.Helper()
	clone := filepath.Join(t.TempDir(), "repo")
	cmd := exec.Command("rsync", "-a", "--exclude", "/out", "--exclude", "/.git", "../", clone+"/")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("rsync: %v\n%s", err, out)
	}
	return clone
}

// editDoc0165 adds readmeSentence0165 to the first paragraph of gofmt's doc
// comment in clone.
func editDoc0165(t *testing.T, clone string) {
	t.Helper()
	path := filepath.Join(clone, "cmd", "fitness-check-gofmt", "doc.go")
	src, err := os.ReadFile(path)
	mustDo(t, err)
	edited := strings.Replace(string(src), "// on by default.\n", "// on by default. "+readmeSentence0165+"\n", 1)
	if edited == string(src) {
		t.Fatal("gofmt's doc comment no longer ends its first paragraph with \"on by default.\"")
	}
	mustDo(t, os.WriteFile(path, []byte(edited), 0o644))
}

// generate0165 runs the clone's generator script with args and env added.
func generate0165(clone string, env []string, args ...string) (string, error) {
	cmd := exec.Command(filepath.Join(clone, "scripts", "generate-readmes.sh"), args...)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func Test0165_1(t *testing.T) {
	t.Parallel()
	here, err := filepath.Abs("..")
	mustDo(t, err)
	emptyCache := []string{"GOMODCACHE=" + t.TempDir(), "GOFLAGS=-modcacherw"}
	if out, err := generate0165(here, emptyCache, "-check"); err != nil {
		t.Errorf("generate-readmes.sh -check: %v\n%s", err, out)
	}
}

func Test0165_2(t *testing.T) {
	t.Parallel()
	clone := readmeClone(t)
	editDoc0165(t, clone)
	out, err := generate0165(clone, nil, "-check")
	if err == nil || !strings.Contains(out, "cmd/fitness-check-gofmt/README.md") {
		t.Errorf("generate-readmes.sh -check after an edit = %v, want a failure naming gofmt's README:\n%s", err, out)
	}
}

func Test0165_3(t *testing.T) {
	t.Parallel()
	clone := readmeClone(t)
	editDoc0165(t, clone)
	if out, err := generate0165(clone, nil); err != nil {
		t.Fatalf("generate-readmes.sh: %v\n%s", err, out)
	}
	readme, err := os.ReadFile(filepath.Join(clone, "cmd", "fitness-check-gofmt", "README.md"))
	mustDo(t, err)
	if !strings.Contains(string(readme), "# gofmt\n\n") || !strings.Contains(string(readme), readmeSentence0165) {
		t.Errorf("gofmt's README lacks its title or the new sentence:\n%s", readme)
	}
	if out, err := generate0165(clone, nil, "-check"); err != nil {
		t.Errorf("a second -check: %v\n%s", err, out)
	}
}
