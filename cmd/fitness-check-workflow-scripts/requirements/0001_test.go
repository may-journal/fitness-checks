// Package requirements proves the workflow-scripts check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 trivial commands and script calls pass", func(t *testing.T) {
		t.Parallel()
		steps := "      - run: make test\n      - run: scripts/release.sh\n      - run: |\n          # a comment\n          make build\n"
		action := "name: A\nruns:\n  using: composite\n  steps:\n    - shell: bash\n      run: \"${{ github.action_path }}/run.sh\"\n"
		out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps), "action.yml": action})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "workflow-scripts", "2")
	})
	t.Run("0001.1.2 several lines of shell fail naming file then line with the fix", func(t *testing.T) {
		t.Parallel()
		steps := "      - name: Tag\n        run: |\n          git fetch\n          git tag v1\n"
		out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml:8: this run step holds 2 lines of shell",
			"move it into a checked-in script such as scripts/<name>.sh, and call that script")
	})
	t.Run("0001.1.3 a pipe or a chain or control flow fails naming file then line then what it holds", func(t *testing.T) {
		t.Parallel()
		action := "name: A\nruns:\n  using: composite\n  steps:\n    - shell: bash\n      run: cat out | grep ok\n    - shell: bash\n      run: make && make test\n" +
			"    - shell: bash\n      run: if true; then make; fi\n    - shell: bash\n      run: \"make build && make test\"\n"
		out, code := scriptsCheck(t, map[string]string{"action.yml": action})
		acceptance.Sees(t, out, code, 1, "action.yml:6: this run step holds a pipe", "action.yml:8: this run step holds an && chain",
			"action.yml:10: this run step holds shell control flow", "action.yml:12: this run step holds an && chain")
		if n := strings.Count(out, "action.yml:"); n != 4 {
			t.Errorf("want the file and line named once per step, got %d:\n%s", n, out)
		}
	})
	t.Run("0001.1.4 operators in quotes and expressions pass", func(t *testing.T) {
		t.Parallel()
		steps := "      - run: printf '%s | %s\\n' a b >> \"$GITHUB_STEP_SUMMARY\"\n      - run: echo \"${{ a || b }}\" \"x;y\"\n" +
			"      - run: make test # lint | format\n      - run: find . -name x -exec rm {} \\;\n      - run: grep ok <<< \"$out\"\n"
		out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "this run step") {
			t.Errorf("operators in quotes or expressions were read as plain text, want no finding:\n%s", out)
		}
	})
	t.Run("0001.1.5 a plain scalar is judged and block text is left alone", func(t *testing.T) {
		t.Parallel()
		steps := "      - run:\n          make a\n          make b\n      - run: make c\n          make d\n      - run: |\n          run: a | b\n"
		out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml:7: this run step holds 2 lines of shell", ".github/workflows/ci.yml:10: this run step holds 2 lines of shell")
		if strings.Contains(out, ".github/workflows/ci.yml:13") {
			t.Errorf("want text inside a block left alone:\n%s", out)
		}
	})
	t.Run("0001.1.6 a repo with no workflow or action files passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := scriptsCheck(t, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "workflow-scripts")
	})
}

// scriptsCheck runs workflow-scripts on happyRepo with files written over it.
func scriptsCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "workflow-scripts")
}

// stepsWorkflow is a workflow whose job runs steps, which start on line 7.
func stepsWorkflow(steps string) string {
	return "name: CI\non: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n" + steps
}
