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
	for _, c := range []struct{ name, want string }{
		{"0001.1.1.1.1 it declares the `local-branch` check", `"name":"local-branch"`},
		{"0001.1.1.1.2 it declares the scope `github`", `"scope":"github"`},
		{"0001.1.1.1.3 it declares protocol 1", `"protocol":1`},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			out, code := acceptance.Execute(t, func() *exec.Cmd { return exec.Command(bin, "--describe") })
			acceptance.Sees(t, out, code, 0, c.want)
		})
	}
	t.Run("0001.1.1.2.1 it exits 3 saying it runs only inside fitness", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Execute(t, func() *exec.Cmd { return exec.Command(bin, "--root", t.TempDir()) })
		acceptance.Sees(t, out, code, 3, "runs only inside fitness")
	})
}
