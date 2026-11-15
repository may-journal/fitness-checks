package requirements

import (
	"strings"
	"testing"
)

// scriptsCheck runs workflow-scripts on happyRepo with files written over it.
func scriptsCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "workflow-scripts")
}

// stepsWorkflow is a workflow whose job runs steps, which start on line 7.
func stepsWorkflow(steps string) string {
	return "name: CI\non: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n" + steps
}

func Test0167_1(t *testing.T) {
	t.Parallel()
	steps := "      - run: make test\n      - run: scripts/release.sh\n      - run: |\n          # a comment\n          make build\n"
	action := "name: A\nruns:\n  using: composite\n  steps:\n    - shell: bash\n      run: \"${{ github.action_path }}/run.sh\"\n"
	out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps), "action.yml": action})
	sees(t, out, code, 0, "All 1 checks passed")
	passedWith(t, out, "workflow-scripts", "2")
}

func Test0167_2(t *testing.T) {
	t.Parallel()
	steps := "      - name: Tag\n        run: |\n          git fetch\n          git tag v1\n"
	out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
	sees(t, out, code, 1, ".github/workflows/ci.yml:8: this run step holds 2 lines of shell",
		"move it into a checked-in script such as scripts/<name>.sh, and call that script")
}

func Test0167_3(t *testing.T) {
	t.Parallel()
	action := "name: A\nruns:\n  using: composite\n  steps:\n    - shell: bash\n      run: cat out | grep ok\n    - shell: bash\n      run: make && make test\n" +
		"    - shell: bash\n      run: if true; then make; fi\n    - shell: bash\n      run: \"make build && make test\"\n"
	out, code := scriptsCheck(t, map[string]string{"action.yml": action})
	sees(t, out, code, 1, "action.yml:6: this run step holds a pipe", "action.yml:8: this run step holds an && chain",
		"action.yml:10: this run step holds shell control flow", "action.yml:12: this run step holds an && chain")
}

func Test0167_4(t *testing.T) {
	t.Parallel()
	steps := "      - run: printf '%s | %s\\n' a b >> \"$GITHUB_STEP_SUMMARY\"\n      - run: echo \"${{ a || b }}\" \"x;y\"\n" +
		"      - run: make test # lint | format\n      - run: find . -name x -exec rm {} \\;\n      - run: grep ok <<< \"$out\"\n"
	out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
	sees(t, out, code, 0, "All 1 checks passed")
}

func Test0167_5(t *testing.T) {
	t.Parallel()
	steps := "      - run:\n          make a\n          make b\n      - run: make c\n          make d\n      - run: |\n          run: a | b\n"
	out, code := scriptsCheck(t, map[string]string{".github/workflows/ci.yml": stepsWorkflow(steps)})
	sees(t, out, code, 1, ".github/workflows/ci.yml:7: this run step holds 2 lines of shell", ".github/workflows/ci.yml:10: this run step holds 2 lines of shell")
	if strings.Contains(out, ".github/workflows/ci.yml:13") {
		t.Errorf("want text inside a block left alone:\n%s", out)
	}
}
