package requirements

import (
	"os/exec"
	"testing"
)

func Test0003_1(t *testing.T) {
	t.Parallel()
	bin := goBuild0008(t, "./cmd/fitness-check-local-branch") + "/fitness-check-local-branch"
	out, code := execute(t, func() *exec.Cmd { return exec.Command(bin, "--describe") })
	sees(t, out, code, 0, `"name":"local-branch"`, `"scope":"github"`, `"protocol":1`)
	out, code = execute(t, func() *exec.Cmd { return exec.Command(bin, "--root", t.TempDir()) })
	sees(t, out, code, 3, "runs only inside fitness")
}
