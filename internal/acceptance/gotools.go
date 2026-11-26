package acceptance

import (
	"os"
	"os/exec"
	"testing"
)

// BrokenTool returns a PATH whose name runs the shell script ahead of the real
// tools, so one tool misbehaves like a broken install.
func BrokenTool(t *testing.T, name, script string) string {
	t.Helper()
	dir := t.TempDir()
	Write(t, dir, map[string]string{name: "#!/bin/sh\n" + script})
	return "PATH=" + dir + string(os.PathListSeparator) + os.Getenv("PATH")
}

// CrashingGo returns a PATH whose go prints message and exits 2 when its first
// argument is sub, else runs the real go.
func CrashingGo(t *testing.T, sub, message string) string {
	t.Helper()
	real, err := exec.LookPath("go")
	MustDo(t, err)
	return BrokenTool(t, "go", "[ \"$1\" = "+sub+" ] && { echo '"+message+"' >&2; exit 2; }\nexec \""+real+"\" \"$@\"\n")
}

// KilledGo returns a PATH whose go kills itself when its first argument is
// sub, as an out-of-memory kill would. Otherwise it runs the real go.
func KilledGo(t *testing.T, sub string) string {
	t.Helper()
	real, err := exec.LookPath("go")
	MustDo(t, err)
	return BrokenTool(t, "go", "[ \"$1\" = "+sub+" ] && kill -9 $$\nexec \""+real+"\" \"$@\"\n")
}
