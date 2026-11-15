package requirements

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// stubLink matches a Markdown link's target.
var stubLink = regexp.MustCompile(`\]\(([^)]+)\)`)

// stub0166 is the README every check folder holds, given its name and its
// --describe summary.
func stub0166(name, summary string) string {
	return "# " + name + "\n\n" + summary + "\n\n" +
		"- [Architecture](../../docs/architecture/03-components.md)\n" +
		"- [Source](doc.go)\n" +
		"- Full docs: [pkg.go.dev](https://pkg.go.dev/github.com/may-journal/fitness-checks/cmd/fitness-check-" + name +
		") or `fitness-check-" + name + " --describe`\n"
}

// summary0166 is what the check built in bin prints as its summary.
func summary0166(t *testing.T, bin, name string) string {
	t.Helper()
	out, err := exec.Command(filepath.Join(bin, "fitness-check-"+name), "--describe").Output()
	mustDo(t, err)
	var described struct {
		Summary string `json:"summary"`
	}
	mustDo(t, json.Unmarshal(out, &described))
	return described.Summary
}

func Test0166_1(t *testing.T) {
	t.Parallel()
	bin := goBuild0160(t, "./cmd/...")
	for _, folder := range checkFolders0160(t, "..") {
		name := strings.TrimPrefix(filepath.Base(folder), "fitness-check-")
		readme, err := os.ReadFile(filepath.Join(folder, "README.md"))
		mustDo(t, err)
		if want := stub0166(name, summary0166(t, bin, name)); string(readme) != want {
			t.Errorf("%s/README.md is not the stub:\n%s\nwant:\n%s", folder, readme, want)
		}
	}
}

func Test0166_2(t *testing.T) {
	t.Parallel()
	for _, folder := range checkFolders0160(t, "..") {
		readme, err := os.ReadFile(filepath.Join(folder, "README.md"))
		mustDo(t, err)
		for _, link := range stubLink.FindAllStringSubmatch(string(readme), -1) {
			if strings.Contains(link[1], "://") {
				continue
			}
			if _, err := os.Stat(filepath.Join(folder, link[1])); err != nil {
				t.Errorf("%s/README.md links %s, which does not resolve: %v", folder, link[1], err)
			}
		}
	}
}
