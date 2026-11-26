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
	trivial := func(t *testing.T) (string, int) {
		steps := "      - run: make test\n      - run: scripts/release.sh\n      - run: |\n          # a comment\n          make build\n"
		action := "name: A\nruns:\n  using: composite\n  steps:\n    - shell: bash\n      run: \"${{ github.action_path }}/run.sh\"\n"
		return scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps), "action.yml": action})
	}
	t.Run("0001.1.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := trivial(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.1.1.2 it counts the files", func(t *testing.T) {
		t.Parallel()
		out, code := trivial(t)
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWith(t, out, "workflow-scripts", "2")
	})
	lines := func(t *testing.T) (string, int) {
		steps := "      - name: Tag\n        run: |\n          git fetch\n          git tag v1\n"
		return scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
	}
	t.Run("0001.1.1.2.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := lines(t)
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml:")
	})
	t.Run("0001.1.1.2.2 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := lines(t)
		acceptance.Sees(t, out, code, 1, "ci.yml:8: this run step holds 2 lines of shell")
	})
	t.Run("0001.1.1.2.3 it names the script to write", func(t *testing.T) {
		t.Parallel()
		out, code := lines(t)
		acceptance.Sees(t, out, code, 1, "move it into a checked-in script such as scripts/<name>.sh, and call that script")
	})
	piped := func(t *testing.T) (string, int) {
		action := "name: A\nruns:\n  using: composite\n  steps:\n    - shell: bash\n      run: cat out | grep ok\n    - shell: bash\n      run: make && make test\n" +
			"    - shell: bash\n      run: if true; then make; fi\n    - shell: bash\n      run: \"make build && make test\"\n"
		return scriptsCheck(t, map[string]string{"action.yml": action})
	}
	t.Run("0001.1.1.3.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := piped(t)
		acceptance.Sees(t, out, code, 1, "action.yml:")
	})
	t.Run("0001.1.1.3.2 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := piped(t)
		acceptance.Sees(t, out, code, 1, "action.yml:6: this run step", "action.yml:8: this run step", "action.yml:10: this run step", "action.yml:12: this run step")
	})
	t.Run("0001.1.1.3.3 it fails naming what it holds", func(t *testing.T) {
		t.Parallel()
		out, code := piped(t)
		acceptance.Sees(t, out, code, 1, "this run step holds a pipe", "this run step holds an && chain", "this run step holds shell control flow")
	})
	t.Run("0001.1.1.4.1 it passes", func(t *testing.T) {
		t.Parallel()
		steps := "      - run: printf '%s | %s\\n' a b >> \"$GITHUB_STEP_SUMMARY\"\n      - run: echo \"${{ a || b }}\" \"x;y\"\n" +
			"      - run: make test # lint | format\n      - run: find . -name x -exec rm {} \\;\n      - run: grep ok <<< \"$out\"\n"
		out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	scalars := func(t *testing.T) (string, int) {
		steps := "      - run:\n          make a\n          make b\n      - run: make c\n          make d\n      - run: |\n          run: a | b\n"
		return scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
	}
	t.Run("0001.1.1.5.1 the scalar is judged", func(t *testing.T) {
		t.Parallel()
		out, code := scalars(t)
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml:7: this run step holds 2 lines of shell", ".github/workflows/ci.yml:10: this run step holds 2 lines of shell")
	})
	t.Run("0001.1.1.5.2 the block text is left alone", func(t *testing.T) {
		t.Parallel()
		out, code := scalars(t)
		acceptance.Sees(t, out, code, 1)
		if strings.Contains(out, ".github/workflows/ci.yml:13") {
			t.Errorf("want text inside a block left alone:\n%s", out)
		}
	})
	t.Run("0001.1.1.6.1 it passes with no files", func(t *testing.T) {
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
