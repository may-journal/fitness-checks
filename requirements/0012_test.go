package requirements

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// describing is a check program whose --describe answer says summary.
func describing(summary string) string {
	return "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(`{\"name\":\"demo\",\"summary\":\"" + summary + "\"}`)\n}\n"
}

// compared runs scripts/describe-diff.sh in a repo that tagged the demo check
// with the first summary and now holds the second. It returns the output and
// exit code.
func compared(t *testing.T, tagged, now string) (string, int) {
	t.Helper()
	repo := t.TempDir()
	script, err := os.ReadFile(filepath.Join("..", "scripts", "describe-diff.sh"))
	mustDo(t, err)
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "tag.gpgsign=false"}, args...)...)
		cmd.Dir = repo
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	write := func(name, body string, mode os.FileMode) {
		t.Helper()
		mustDo(t, os.MkdirAll(filepath.Dir(filepath.Join(repo, name)), 0o755))
		mustDo(t, os.WriteFile(filepath.Join(repo, name), []byte(body), mode))
	}
	git("init", "-q")
	write("go.mod", "module demo\n\ngo 1.21\n", 0o644)
	write("scripts/describe-diff.sh", string(script), 0o755)
	write("cmd/fitness-check-demo/main.go", describing(tagged), 0o644)
	git("add", "-A")
	git("commit", "-q", "-m", "feat(demo): first release")
	git("tag", "fitness-check-demo-v1.0.0")
	write("cmd/fitness-check-demo/main.go", describing(now), 0o644)
	git("add", "-A")
	git("commit", "-q", "--allow-empty", "-m", "feat(demo): the batch")
	cmd := exec.Command("scripts/describe-diff.sh", "HEAD", "demo")
	cmd.Dir = repo
	out, _ := cmd.CombinedOutput()
	return string(out), cmd.ProcessState.ExitCode()
}

func Test0012_1(t *testing.T) {
	t.Parallel()
	t.Run("0012.1.1 an unchanged answer prints identical and exits 0", func(t *testing.T) {
		t.Parallel()
		out, code := compared(t, "Says one.", "Says one.")
		if code != 0 || !strings.Contains(out, "demo: identical to fitness-check-demo-v1.0.0") {
			t.Errorf("want identical and exit 0, got %d:\n%s", code, out)
		}
	})
	t.Run("0012.1.2 a changed answer prints the diff and exits 1", func(t *testing.T) {
		t.Parallel()
		out, code := compared(t, "Says one.", "Says two.")
		if code != 1 || !strings.Contains(out, "-{\"name\":\"demo\",\"summary\":\"Says one.\"}") || !strings.Contains(out, "+{\"name\":\"demo\",\"summary\":\"Says two.\"}") {
			t.Errorf("want the diff and exit 1, got %d:\n%s", code, out)
		}
	})
}
