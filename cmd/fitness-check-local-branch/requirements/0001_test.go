// Package requirements proves the local-branch check's docs/requirements as a
// user meets them. Each test runs the program built from this check's folder.
package requirements

import (
	"os/exec"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	bin := t.TempDir() + "/fitness-check-local-branch"
	build := exec.Command("go", "build", "-o", bin, "..")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	t.Run("0001.1.1 it declares the check, its scope and its protocol", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Execute(t, func() *exec.Cmd { return exec.Command(bin, "--describe") })
		acceptance.Sees(t, out, code, 0, `"name":"local-branch"`, `"scope":"github"`, `"protocol":1`)
	})
	t.Run("0001.1.2 run alone it exits saying it runs only inside fitness", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Execute(t, func() *exec.Cmd { return exec.Command(bin, "--root", t.TempDir()) })
		acceptance.Sees(t, out, code, 3, "runs only inside fitness")
	})
}
