// Package requirements proves the local-branch check's docs/requirements as a
// user meets them. Each test runs the program built from this check's folder.
package requirements

import (
	"os/exec"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
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
		req.Then(t, "0001.1.1.1.1", "db35c93d", func(a req.A) {
			acceptance.Sees(t, out, code, 0, `"name":"local-branch"`)
		})
		req.Then(t, "0001.1.1.1.2", "1874452a", func(a req.A) {
			acceptance.Sees(t, out, code, 0, `"scope":"github"`)
		})
		req.Then(t, "0001.1.1.1.3", "ef391279", func(a req.A) {
			acceptance.Sees(t, out, code, 0, `"protocol":1`)
		})
	})
	t.Run("0001.1.2 run alone it exits saying it runs only inside fitness", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "d4ef3dfd", func(_ req.A) {
			out, code := acceptance.Execute(t, func() *exec.Cmd { return exec.Command(bin, "--root", t.TempDir()) })
			acceptance.Sees(t, out, code, 3, "runs only inside fitness")
		})
	})
}
