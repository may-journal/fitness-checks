package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/program"
)

// stubTargets are the stub's relative links, as a check folder resolves them.
var stubTargets = []string{"../../docs/architecture/03-components.md", "doc.go"}

// stub0166 is the README every check folder holds, given its name and the
// summary its doc comment opens with, as --describe prints it.
func stub0166(name, summary string) string {
	return "# " + name + "\n\n" + summary + "\n\n" +
		"- [Architecture](" + stubTargets[0] + ")\n" +
		"- [Source](" + stubTargets[1] + ")\n" +
		"- Full docs: [pkg.go.dev](https://pkg.go.dev/github.com/may-journal/fitness-checks/cmd/fitness-check-" + name +
		") or `fitness-check-" + name + " --describe`\n"
}

// readFolder0166 reads file in folder, or reports why it could not.
func readFolder0166(t *testing.T, folder, file string) (string, bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(folder, file))
	if err != nil {
		t.Errorf("%s: %v", folder, err)
		return "", false
	}
	return string(data), true
}

func Test0166_1(t *testing.T) {
	t.Parallel()
	for _, folder := range checkFolders0160(t, "..") {
		name := strings.TrimPrefix(filepath.Base(folder), "fitness-check-")
		source, okDoc := readFolder0166(t, folder, "doc.go")
		readme, okReadme := readFolder0166(t, folder, "README.md")
		if want := stub0166(name, program.Summary(source)); okDoc && okReadme && readme != want {
			t.Errorf("%s/README.md is not the stub:\n%s\nwant:\n%s", folder, readme, want)
		}
	}
}

func Test0166_2(t *testing.T) {
	t.Parallel()
	for _, folder := range checkFolders0160(t, "..") {
		for _, target := range stubTargets {
			if _, err := os.Stat(filepath.Join(folder, target)); err != nil {
				t.Errorf("%s/README.md links %s, which does not resolve: %v", folder, target, err)
			}
		}
	}
}
