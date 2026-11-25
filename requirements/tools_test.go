package requirements

import (
	"os"
	"os/exec"
	"testing"
)

// killedVet returns a PATH whose go kills `go vet` as an out-of-memory kill
// would, and runs the real go otherwise.
func killedVet(t *testing.T) string {
	t.Helper()
	real, err := exec.LookPath("go")
	mustDo(t, err)
	dir := t.TempDir()
	write(t, dir, map[string]string{"go": "#!/bin/sh\n[ \"$1\" = vet ] && kill -9 $PPID\nexec \"" + real + "\" \"$@\"\n"})
	return "PATH=" + dir + string(os.PathListSeparator) + os.Getenv("PATH")
}

// brokenTool returns a PATH whose name runs the shell script ahead of the real
// tools, so one tool misbehaves like a broken install.
func brokenTool(t *testing.T, name, script string) string {
	t.Helper()
	dir := t.TempDir()
	write(t, dir, map[string]string{name: "#!/bin/sh\n" + script})
	return "PATH=" + dir + string(os.PathListSeparator) + os.Getenv("PATH")
}

// crashingGo returns a PATH whose go prints message and exits 2 when its first
// argument is sub, else runs the real go.
func crashingGo(t *testing.T, sub, message string) string {
	t.Helper()
	real, err := exec.LookPath("go")
	mustDo(t, err)
	return brokenTool(t, "go", "[ \"$1\" = "+sub+" ] && { echo '"+message+"' >&2; exit 2; }\nexec \""+real+"\" \"$@\"\n")
}
