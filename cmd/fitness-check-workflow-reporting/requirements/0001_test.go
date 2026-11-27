// Package requirements proves the workflow-reporting check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a final always summary step passes", func(t *testing.T) {
		t.Parallel()
		out, code := reportingCheck(t, map[string]string{".github/workflows/ci.yml": reportingJob(reportStep)})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "workflow-reporting", "1")
	})
	t.Run("0001.1.2 a job with no reporting step fails naming file then line with the step to add", func(t *testing.T) {
		t.Parallel()
		out, code := reportingCheck(t, map[string]string{".github/workflows/build.yaml": reportingJob("")})
		acceptance.Sees(t, out, code, 1, ".github/workflows/build.yaml:4: job build has no step that runs with if: always() and writes $GITHUB_STEP_SUMMARY",
			"add a last step with `if: ${{ always() }}` whose run writes the job's status to \"$GITHUB_STEP_SUMMARY\"")
	})
	t.Run("0001.1.3 a summary step that skips failure fails", func(t *testing.T) {
		t.Parallel()
		onSuccess := strings.Replace(reportStep, "        if: ${{ always() }}\n", "", 1)
		out, code := reportingCheck(t, map[string]string{".github/workflows/ci.yml": reportingJob(onSuccess)})
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml:4: job build writes $GITHUB_STEP_SUMMARY only when earlier steps pass, at line 9",
			"add `if: ${{ always() }}` to the step at line 9")
	})
	t.Run("0001.1.4 a job calling a reusable workflow passes", func(t *testing.T) {
		t.Parallel()
		out, code := reportingCheck(t, map[string]string{".github/workflows/ci.yml": "name: CI\non: push\njobs:\n" +
			"  fitness:\n    uses: may-journal/fitness-runner/.github/workflows/ci-reusable.yml@v1\n    secrets: inherit\n"})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "workflow-reporting", "1")
		if strings.Contains(out, "job fitness") {
			t.Errorf("a job that only calls a reusable workflow was asked for a step:\n%s", out)
		}
	})
	t.Run("0001.1.5 a matrix job is judged and named once", func(t *testing.T) {
		t.Parallel()
		matrix := "name: CI\non: push\njobs:\n  test:\n    strategy:\n      matrix:\n        os: [ubuntu-latest, macos-latest, windows-latest]\n" +
			"    runs-on: ${{ matrix.os }}\n    steps:\n      - run: make test\n"
		out, code := reportingCheck(t, map[string]string{".github/workflows/ci.yml": matrix})
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml:4: job test has no step")
		if n := strings.Count(out, "job test has no step"); n != 1 {
			t.Errorf("a matrix job must be named once, got %d times:\n%s", n, out)
		}
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 YAML GitHub rejects fails at its line then its column", func(t *testing.T) {
		t.Parallel()
		files := runsOnWorkflows(rejectedRunsOn)
		files[".github/workflows/ci.yml"] = "name: CI\non: push\njobs:\n  build:\n\truns-on: ubuntu-latest\n"
		out, code := reportingCheck(t, files)
		const bad = ": this workflow is not valid YAML: "
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml:5"+bad+"a tab indents this line",
			".github/workflows/fork.yml:5:116"+bad+"a `: ` inside a plain value starts a mapping",
			".github/workflows/reserved.yml:5:14"+bad+"a plain value cannot start with '@'",
			".github/workflows/key-line.yml:5:14"+bad+"'-' opens a list item on its key's line",
			"move it onto its own line, indented under the key",
			".github/workflows/continued.yml:6:13"+bad+"a `: ` inside a plain value",
			".github/workflows/colon-end.yml:5:27"+bad+"a plain value ends in `:`",
			".github/workflows/flow-trailing.yml:5:27"+bad+"text follows the closing bracket",
			".github/workflows/wide.yml:5:18"+bad+"a `: ` inside a plain value")
		out, code = reportingCheck(t, runsOnWorkflows(acceptedRunsOn))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2 a repo with no workflows passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := reportingCheck(t, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "workflow-reporting")
	})
	t.Run("0001.2.3 this repo's own workflows pass", func(t *testing.T) {
		t.Parallel()
		own, err := filepath.Glob(filepath.Join("..", "..", "..", ".github", "workflows", "*.y*ml"))
		acceptance.MustDo(t, err)
		if len(own) == 0 {
			t.Fatal("no workflows found under .github/workflows")
		}
		files := map[string]string{}
		for _, f := range own {
			body, err := os.ReadFile(f)
			acceptance.MustDo(t, err)
			files[".github/workflows/"+filepath.Base(f)] = string(body)
		}
		out, code := reportingCheck(t, files)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

// reportStep is a last step that reports every outcome to the job summary.
const reportStep = `      - name: Report job outcome
        if: ${{ always() }}
        run: printf 'Job %s\n' "${{ job.status }}" >> "$GITHUB_STEP_SUMMARY"
`

// reportingJob is a workflow whose job, on line 4, builds and then
// reports with tail as its last steps.
func reportingJob(tail string) string {
	return "name: CI\non: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - run: make\n" + tail
}

// reportingCheck runs workflow-reporting on happyRepo with files written
// over it.
func reportingCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "workflow-reporting")
}

// mergedRunsOn is the runs-on line that merged into three workflows and
// failed every run, since GitHub reads its `": "` as a mapping value.
const mergedRunsOn = `${{ github.event.repository.private && !github.event.pull_request.head.repo.fork && fromJSON('{"group": "fitness", "labels": "fitness"}') || 'ubuntu-24.04-arm' }}`

// runsOnJob is a reporting workflow whose job runs on runsOn, on line 5.
func runsOnJob(runsOn string) string {
	return strings.Replace(reportingJob(reportStep), "runs-on: ubuntu-latest", "runs-on: "+runsOn, 1)
}

// rejectedRunsOn are runs-on values GitHub rejects as YAML, by file,
// each failing on line 5 or 6 at the column its message follows.
var rejectedRunsOn = map[string]string{
	"fork.yml":          mergedRunsOn,
	"reserved.yml":      "@self-hosted",
	"key-line.yml":      "- self-hosted",
	"continued.yml":     "ubuntu\n      latest: now",
	"colon-end.yml":     "ubuntu-latest:",
	"flow-trailing.yml": "[self-hosted]: x",
	"wide.yml":          "café: x",
}

// acceptedRunsOn are the quoted or reshaped forms GitHub accepts.
var acceptedRunsOn = map[string]string{
	"fork.yml":          `"` + strings.ReplaceAll(mergedRunsOn, `"`, `\"`) + `"`,
	"reserved.yml":      `"@self-hosted"`,
	"key-line.yml":      "\n      - self-hosted",
	"continued.yml":     "ubuntu\n      latest now",
	"colon-end.yml":     `"ubuntu-latest:"`,
	"flow-trailing.yml": "[self-hosted] # a comment",
	"wide.yml":          `"café: x"`,
}

// runsOnWorkflows writes each runs-on value into a reporting workflow.
func runsOnWorkflows(values map[string]string) map[string]string {
	files := map[string]string{}
	for name, runsOn := range values {
		files[".github/workflows/"+name] = runsOnJob(runsOn)
	}
	return files
}
