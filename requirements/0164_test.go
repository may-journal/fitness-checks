package requirements

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// dirtyOnly0164 writes a sample list holding only the dirty sample, so a
// run needs no network.
func dirtyOnly0164(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("..")
	mustDo(t, err)
	list := filepath.Join(t.TempDir(), "samples.txt")
	mustDo(t, os.WriteFile(list, []byte("dirty script "+filepath.Join(root, "tests", "parity", "make-dirty.sh")+"\n"), 0o644))
	return list
}

// programs0164 is the one build of the engine and every check that the tests copy.
var programs0164 struct {
	once sync.Once
	dir  string
	err  error
}

// builtPrograms0164 returns a folder of links to the pinned engine and every check,
// built once for the whole run.
func builtPrograms0164(t *testing.T) string {
	t.Helper()
	programs0164.once.Do(func() {
		programs0164.dir, programs0164.err = buildPrograms0164()
	})
	mustDo(t, programs0164.err)
	copied := t.TempDir()
	entries, err := os.ReadDir(programs0164.dir)
	mustDo(t, err)
	for _, entry := range entries {
		mustDo(t, os.Symlink(filepath.Join(programs0164.dir, entry.Name()), filepath.Join(copied, entry.Name())))
	}
	return copied
}

// buildPrograms0164 builds the pinned engine and every check into a new folder.
func buildPrograms0164() (string, error) {
	bin, err := os.MkdirTemp("", "parity-programs")
	if err != nil {
		return "", err
	}
	for _, target := range []string{"github.com/may-journal/fitness-runner/go/cmd/fitness", "./cmd/..."} {
		build := exec.Command("go", "build", "-o", bin+string(os.PathSeparator), target)
		build.Dir = ".."
		if out, err := build.CombinedOutput(); err != nil {
			return "", fmt.Errorf("go build %s: %v\n%s", target, err, out)
		}
	}
	return bin, nil
}

// parity0164 runs scripts/check-parity.sh with programs in bin, and returns
// its output, its exit code and its out folder.
func parity0164(t *testing.T, bin string) (string, int, string) {
	t.Helper()
	out := t.TempDir()
	cmd := exec.Command("scripts/check-parity.sh", out)
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), "CHECK_PARITY_BIN="+bin, "CHECK_PARITY_SAMPLES="+dirtyOnly0164(t))
	got, err := cmd.CombinedOutput()
	if exit, ok := err.(*exec.ExitError); ok {
		return string(got), exit.ExitCode(), out
	}
	mustDo(t, err)
	return string(got), 0, out
}

// withFindings0164 counts the results.tsv rows whose findings are above zero.
func withFindings0164(t *testing.T, out string) (rows, found int) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(out, "results.tsv"))
	mustDo(t, err)
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n")[1:] {
		rows++
		if cols := strings.Split(line, "\t"); len(cols) == 4 && cols[3] != "" && cols[3] != "0" {
			found++
		}
	}
	return rows, found
}

func Test0164_1(t *testing.T) {
	t.Parallel()
	bin := builtPrograms0164(t)
	listed, err := exec.Command(filepath.Join(bin, "fitness"), "check", "--list").Output()
	mustDo(t, err)
	got, code, out := parity0164(t, bin)
	rows, found := withFindings0164(t, out)
	tsv, err := os.ReadFile(filepath.Join(out, "results.tsv"))
	mustDo(t, err)
	if !strings.Contains(string(tsv), "\tlocal-branch\tskipped\t") {
		t.Errorf("results.tsv does not skip local-branch, which the engine judges:\n%s", tsv)
	}
	if want := len(strings.Fields(string(listed))); code != 0 || rows != want || found < 10 {
		t.Errorf("exit %d with %d rows and %d with findings, want 0, %d rows, one per built-in check, and at least 10 with findings:\n%s", code, rows, found, want, got)
	}
}

func Test0164_2(t *testing.T) {
	t.Parallel()
	bin := builtPrograms0164(t)
	tampered := "#!/bin/sh\nprintf '{\"ok\":false,\"errors\":[\"tampered\"],\"filesChecked\":0}\\n'\n"
	mustDo(t, os.Remove(filepath.Join(bin, "fitness-check-gofmt")))
	mustDo(t, os.WriteFile(filepath.Join(bin, "fitness-check-gofmt"), []byte(tampered), 0o755))
	got, code, _ := parity0164(t, bin)
	if code != 1 || !strings.Contains(got, "MISMATCH: check gofmt on sample dirty") {
		t.Errorf("exit %d with %q, want 1 naming check gofmt on sample dirty", code, got)
	}
}

func Test0164_3(t *testing.T) {
	t.Parallel()
	bin := t.TempDir()
	engine := "#!/bin/sh\nif [ \"$2\" = --list ]; then echo gofmt; exit 0; fi\n" +
		"n=$(cat \"$(dirname \"$0\")/count\" 2> /dev/null || echo 0)\necho $((n + 1)) > \"$(dirname \"$0\")/count\"\nprintf '{\"n\":%s}\\n' \"$n\"\n"
	mustDo(t, os.WriteFile(filepath.Join(bin, "fitness"), []byte(engine), 0o755))
	mustDo(t, os.WriteFile(filepath.Join(bin, "fitness-check-gofmt"), []byte("#!/bin/sh\nprintf '{\"n\":0}\\n'\n"), 0o755))
	got, code, out := parity0164(t, bin)
	b, err := os.ReadFile(filepath.Join(out, "results.tsv"))
	mustDo(t, err)
	if code != 0 || !strings.Contains(string(b), "dirty\tgofmt\tunstable") {
		t.Errorf("exit %d with results %q and output %q, want 0 and gofmt unstable", code, b, got)
	}
}

func Test0164_4(t *testing.T) {
	t.Parallel()
	bin := t.TempDir()
	engine := "#!/bin/sh\nif [ \"$2\" = --list ]; then echo gofmt; exit 0; fi\n" +
		"n=$(cat \"$(dirname \"$0\")/count\" 2> /dev/null || echo 0)\necho $((n + 1)) > \"$(dirname \"$0\")/count\"\nprintf '{\"ok\":true,\"n\":%s}\\n' \"$n\"\n"
	mustDo(t, os.WriteFile(filepath.Join(bin, "fitness"), []byte(engine), 0o755))
	mustDo(t, os.WriteFile(filepath.Join(bin, "fitness-check-gofmt"), []byte("#!/bin/sh\nprintf '{\"ok\":false,\"n\":0}\\n'\n"), 0o755))
	got, code, _ := parity0164(t, bin)
	if code != 1 || !strings.Contains(got, "MISMATCH: check gofmt on sample dirty") {
		t.Errorf("exit %d with %q, want 1 naming check gofmt on sample dirty", code, got)
	}
}
