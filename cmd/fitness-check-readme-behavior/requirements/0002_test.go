package requirements

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// installedLayout builds the named checks, each in a folder of its own, and
// returns the folder of links that FITNESS_CHECKS_DIR names.
func installedLayout(t *testing.T, names ...string) string {
	t.Helper()
	checks, builds := t.TempDir(), t.TempDir()
	for _, name := range names {
		program := "fitness-check-" + name
		home := filepath.Join(builds, program+"-v1.0.0-"+runtime.GOOS+"-"+runtime.GOARCH+".tar")
		build := exec.Command("go", "build", "-o", filepath.Join(home, program), "./cmd/"+program)
		build.Dir = "../../.."
		if out, err := build.CombinedOutput(); err != nil {
			t.Fatalf("go build %s: %v\n%s", program, err, out)
		}
		acceptance.MustDo(t, os.Symlink(filepath.Join(home, program), filepath.Join(checks, program)))
	}
	return checks
}

func Test0002_1(t *testing.T) {
	t.Parallel()
	folder := "cmd/fitness-check-markdown-no-bold-italic/"
	repo := acceptance.Example(t, "happyRepo", map[string]string{
		folder + "README.md":                          readmeOf(link, command("markdown-no-bold-italic")),
		folder + "docs/requirements/0001-it-works.md": "# 0001 It Works\n",
		"other/README.md":                             "# other\n",
	})
	// The checks folder holds either one check or both. Each is linked to a
	// folder of its own.
	for _, names := range [][]string{{"readme-behavior"}, {"readme-behavior", "markdown-no-bold-italic"}} {
		out, code := acceptance.Fitness(t, repo, []string{"FITNESS_CHECKS_DIR=" + installedLayout(t, names...)}, "readme-behavior")
		acceptance.Sees(t, out, code, 1, "the Example misses the flag `--body-file`")
	}
}
