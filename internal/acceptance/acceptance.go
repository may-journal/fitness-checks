// Package acceptance holds the helpers every requirement test folder shares:
// the repo root's requirements folder and each package's own. A test runs
// the installer at FITNESS_INSTALL in a copied example repo, as a user does.
package acceptance

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/may-journal/fitness-checks/req"
)

// ANSI matches terminal color codes, which a reader does not see as text.
var ANSI = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Example copies tests/examples/<name> into a fresh directory, writes the
// Given's files over it, and commits it as a user's repo would be. Each test
// skips first, before any setup, when no candidate installer is set.
func Example(t *testing.T, name string, files map[string]string) string {
	t.Helper()
	Installer(t)
	repo := t.TempDir()
	if err := os.CopyFS(repo, os.DirFS(filepath.Join(moduleRoot(t), "tests", "examples", name))); err != nil {
		t.Fatal(err)
	}
	Commit(t, repo, files, "docs(readme): add readme")
	return repo
}

// Commit writes files into repo and commits everything with message.
func Commit(t *testing.T, repo string, files map[string]string, message string) {
	t.Helper()
	Write(t, repo, files)
	Git(t, repo, "init", "-q")
	Git(t, repo, "add", "-A")
	Git(t, repo, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "--allow-empty", "-m", message)
}

// Write puts each file into repo, creating folders as needed. An empty body
// deletes the file.
func Write(t *testing.T, repo string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		path := filepath.Join(repo, filepath.FromSlash(name))
		if body == "" {
			_ = os.Remove(path)
			continue
		}
		MustDo(t, os.MkdirAll(filepath.Dir(path), 0o755))
		MustDo(t, os.WriteFile(path, []byte(body), 0o755))
	}
}

// Git runs git in repo with args, failing the test when git fails.
func Git(t *testing.T, repo string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", repo}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// MustDo fails the test on err.
func MustDo(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// Fitness runs the documented one-liner, `fitness-install -- <args>`, in
// repo with env added, and returns what it printed and its exit code. With
// no args it runs the full suite.
func Fitness(t *testing.T, repo string, env []string, args ...string) (string, int) {
	t.Helper()
	return Piped(t, repo, "", env, args...)
}

// Piped runs the one-liner like fitness, with stdin as the text a user pipes
// in; empty stdin means none.
func Piped(t *testing.T, repo, stdin string, env []string, args ...string) (string, int) {
	t.Helper()
	if len(args) == 0 {
		args = []string{"--policy=external", "--all"}
	}
	return Execute(t, func() *exec.Cmd {
		cmd := exec.Command(Installer(t), append([]string{"--"}, args...)...)
		cmd.Dir, cmd.Env = repo, append(UserEnv(), env...)
		if stdin != "" {
			cmd.Stdin = strings.NewReader(stdin)
		}
		return cmd
	})
}

// Execute runs the command build makes and returns its output, color codes
// dropped, and exit code, retrying a busy binary briefly.
func Execute(t *testing.T, build func() *exec.Cmd) (string, int) {
	t.Helper()
	var out []byte
	cmd := Started(t, build, func(cmd *exec.Cmd) (err error) {
		out, err = cmd.CombinedOutput()
		return err
	})
	return ANSI.ReplaceAllString(string(out), ""), cmd.ProcessState.ExitCode()
}

// Started runs each command build makes with run until one starts, and
// returns it. A busy binary is retried as execute says; any other start
// error fails the test.
func Started(t *testing.T, build func() *exec.Cmd, run func(*exec.Cmd) error) *exec.Cmd {
	t.Helper()
	for attempt := 0; ; attempt++ {
		cmd := build()
		err := run(cmd)
		if cmd.ProcessState != nil {
			return cmd
		}
		if !errors.Is(err, syscall.ETXTBSY) || attempt == 40 {
			t.Fatalf("start %s: %v", cmd.Path, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// PlaceRunner hard-links the installed runner at dst, so no parallel test's fork
// inherits an open write handle. Where linking fails, it writes, closes and
// renames a temporary file into place.
func PlaceRunner(t *testing.T, installed, dst string) {
	t.Helper()
	if os.Link(installed, dst) == nil {
		return
	}
	raw, err := os.ReadFile(installed)
	MustDo(t, err)
	tmp := dst + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o755)
	MustDo(t, err)
	_, err = f.Write(raw)
	MustDo(t, errors.Join(err, f.Close()))
	MustDo(t, os.Rename(tmp, dst))
}

// UserEnv is this process's environment without the CI job's GITHUB_*
// variables, which a user's own terminal does not have.
func UserEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "GITHUB_") {
			env = append(env, kv)
		}
	}
	return env
}

// Installer returns the candidate fitness-install the environment provides.
func Installer(t *testing.T) string {
	t.Helper()
	path := os.Getenv("FITNESS_INSTALL")
	if path == "" {
		t.Skip("set FITNESS_INSTALL to the candidate fitness-install to run requirement tests")
	}
	return path
}

// Read returns text as a person reads it, table borders dropped and wrapped
// lines joined into one.
func Read(text string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(text, "│", " ")), " ")
}

// Sees asserts the user saw every want and got exit code code.
func Sees(t *testing.T, out string, gotCode, code int, want ...string) {
	t.Helper()
	req.Asserted(t)
	for _, w := range want {
		if !strings.Contains(Read(out), Read(w)) {
			t.Errorf("output is missing %q:\n%s", w, out)
		}
	}
	if gotCode != code {
		t.Errorf("exit code = %d, want %d:\n%s", gotCode, code, out)
	}
}

// ResultRow matches a results-table line: name, status, files, time.
var ResultRow = regexp.MustCompile(`^│ (\S+)\s+│ (passed|failed|info|skipped)\s+│`)

// Row returns the results-table line for check, or "". The table shortens
// long names with "…", as a reader sees them.
func Row(out, check string) string {
	for _, line := range strings.Split(out, "\n") {
		if m := ResultRow.FindStringSubmatch(line); m != nil && Names(m[1], check) {
			return line
		}
	}
	return ""
}

// Names reports whether a table cell shows check, whole or shortened.
func Names(cell, check string) bool {
	short, cut := strings.CutSuffix(cell, "…")
	return cell == check || cut && strings.HasPrefix(check, short)
}

// PassedWithNoFiles asserts check's row shows a pass that judged 0 files, or
// a skip. The engine skips a check that declares it applies to nothing here.
func PassedWithNoFiles(t *testing.T, out, check string) {
	t.Helper()
	req.Asserted(t)
	if !regexp.MustCompile(`│ (passed\s+│ 0|skipped\s+│ -)\s+│`).MatchString(Row(out, check)) {
		t.Errorf("%s must pass with 0 files, got row %q", check, Row(out, check))
	}
}

// ShowsInfo asserts the run exited 0 printing every want, with check's row
// showing info because the check could not judge.
func ShowsInfo(t *testing.T, out string, code int, check string, want ...string) {
	t.Helper()
	Sees(t, out, code, 0, want...)
	if !strings.Contains(Row(out, check), "info") {
		t.Errorf("%s must show info, got row %q", check, Row(out, check))
	}
}

// PassedWith asserts check's row shows a pass that judged files.
func PassedWith(t *testing.T, out, check, files string) {
	t.Helper()
	req.Asserted(t)
	if !regexp.MustCompile(`│ passed\s+│ ` + files + `\s+│`).MatchString(Row(out, check)) {
		t.Errorf("%s must pass with %s files, got row %q", check, files, Row(out, check))
	}
}

// moduleRoot is the folder holding this module's go.mod, found from the
// test's folder upward. A package's tests find the examples there too.
func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	MustDo(t, err)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the test folder")
		}
		dir = parent
	}
}

// repoGitConfig is the git config every test repo gets, whatever the host's,
// since a signed tag needs a message no test gives. Background maintenance after
// a commit can still be writing a pack when the test removes the repo.
var repoGitConfig = [][2]string{{"tag.gpgSign", "false"}, {"maintenance.auto", "false"}, {"gc.auto", "0"}}

// KeepHostGitOut puts repoGitConfig in the environment, so every git the
// tests or fitness start reads it.
func KeepHostGitOut() {
	for i, kv := range repoGitConfig {
		_ = os.Setenv(fmt.Sprintf("GIT_CONFIG_KEY_%d", i), kv[0])
		_ = os.Setenv(fmt.Sprintf("GIT_CONFIG_VALUE_%d", i), kv[1])
	}
	_ = os.Setenv("GIT_CONFIG_COUNT", fmt.Sprint(len(repoGitConfig)))
}

// Remote gives repo an empty bare origin to push to, so no test needs GitHub.
func Remote(t *testing.T, repo string) {
	t.Helper()
	bare := t.TempDir()
	Git(t, bare, "init", "-q", "--bare")
	Git(t, repo, "remote", "add", "origin", bare)
}

// User runs git in repo as a person at the terminal does, with hooks,
// returning its output and exit code.
func User(t *testing.T, repo string, env []string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t"}, args...)...)
	cmd.Dir, cmd.Env = repo, append(UserEnv(), env...)
	out, _ := cmd.CombinedOutput()
	return ANSI.ReplaceAllString(string(out), ""), cmd.ProcessState.ExitCode()
}

// Main runs a requirements package's tests with the host's git config kept
// out, serving the named stand-ins when run under their names.
func Main(m *testing.M, standIns ...string) {
	if len(standIns) > 0 {
		ServeStandIns(standIns)
	}
	KeepHostGitOut()
	os.Exit(m.Run())
}

// Diagram is a fenced mermaid diagram with the given lines.
func Diagram(lines ...string) string {
	return "```mermaid\n" + strings.Join(lines, "\n") + "\n```\n"
}

// DiagramDoc is a README followed by diagrams, whose first fence sits on
// line 9.
func DiagramDoc(diagrams ...string) map[string]string {
	readme := "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n\n"
	return map[string]string{"README.md": readme + strings.Join(diagrams, "\n")}
}
