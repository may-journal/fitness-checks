// Package requirements proves the read-repo-first check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1.1.1 it passes with zero files", func(t *testing.T) {
		t.Parallel()
		out, code := readRepoFirst(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "read-repo-first")
	})
	t.Run("0001.1.1.1.2 it asks whether the decisions were read", func(t *testing.T) {
		t.Parallel()
		out, code := readRepoFirst(t)
		acceptance.Sees(t, out, code, 0, "Did you familiarize yourself with the decisions logged in the repo")
	})
	t.Run("0001.1.1.1.3 it asks whether the enabled checks were read", func(t *testing.T) {
		t.Parallel()
		out, code := readRepoFirst(t)
		acceptance.Sees(t, out, code, 0, `specifically all "Fitness Checks" that are enabled via @mayjournal/fitness?`)
	})
	t.Run("0001.1.1.1.4 it warns never to use `--no-verify`", func(t *testing.T) {
		t.Parallel()
		out, code := readRepoFirst(t)
		acceptance.Sees(t, out, code, 0, "NOTE: Do not under any circumstance use `--no-verify`")
	})
	t.Run("0001.1.1.2.1 the table lists only the enabled checks", func(t *testing.T) {
		t.Parallel()
		out, code := readRepoFirst(t, "--policy=external", "--checks=read-repo-first,duplicate-code")
		acceptance.Sees(t, out, code, 0, "Check │ Src", "duplicate-code │ go/cmd/fitness-check-duplicate-code")
		if strings.Contains(out, "fitness-check-gofmt") {
			t.Errorf("table lists gofmt, which this run did not enable:\n%s", out)
		}
	})
	t.Run("0001.1.1.3.1 the cell is shortened with an ellipsis", func(t *testing.T) {
		t.Parallel()
		out, code := readRepoFirst(t, "--policy=external", "--checks=read-repo-first,markdown-filename-kebab-case")
		acceptance.Sees(t, out, code, 0, "markdown-filename-kebab-c… │ go/cmd/fitness-check-markdown-file"+"n…")
	})
	t.Run("0001.1.1.4.1 the table is colored", func(t *testing.T) {
		t.Parallel()
		if colored := readRepoFirstRaw(t, "NO_COLOR="); !acceptance.ANSI.MatchString(colored) {
			t.Errorf("without NO_COLOR the table should be colored:\n%q", colored)
		}
	})
	t.Run("0001.1.1.5.1 the table prints plain", func(t *testing.T) {
		t.Parallel()
		plain := readRepoFirstRaw(t, "NO_COLOR=1")
		if acceptance.ANSI.MatchString(plain) || !strings.Contains(plain, "│ Check") {
			t.Errorf("with NO_COLOR the table should print plain:\n%q", plain)
		}
	})
}

// readRepoFirst runs `fitness-install -- <args>` on happyRepo; with no args
// it runs read-repo-first alone.
func readRepoFirst(t *testing.T, args ...string) (string, int) {
	t.Helper()
	if len(args) == 0 {
		args = []string{"read-repo-first"}
	}
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, args...)
}

// readRepoFirstRaw runs read-repo-first on happyRepo with env added and
// returns its banner, up to the note, exactly as printed with color codes.
func readRepoFirstRaw(t *testing.T, env ...string) string {
	t.Helper()
	cmd := exec.Command(acceptance.Installer(t), "--", "read-repo-first")
	cmd.Dir, cmd.Env = acceptance.Example(t, "happyRepo", nil), append(acceptance.UserEnv(), env...)
	out, err := cmd.CombinedOutput()
	acceptance.MustDo(t, err)
	banner, _, _ := strings.Cut(string(out), "NOTE:")
	return banner
}
