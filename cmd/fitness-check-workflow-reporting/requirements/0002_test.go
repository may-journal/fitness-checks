package requirements

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func Test0002_1(t *testing.T) {
	t.Parallel()
	const bad = ": this workflow is not valid YAML: "
	t.Run("0002.1.1 a colon plus space or a trailing colon fails at its place", func(t *testing.T) {
		t.Parallel()
		out, code := reportingCheck(t, runsOnWorkflows(map[string]string{
			"continued.yml": "ubuntu\n      latest: now",
			"colon-end.yml": "ubuntu-latest:",
			"wide.yml":      "café: x",
		}))
		acceptance.Sees(t, out, code, 1,
			".github/workflows/continued.yml:6:13"+bad+"a `: ` inside a plain value",
			".github/workflows/colon-end.yml:5:27"+bad+"a plain value ends in `:`",
			".github/workflows/wide.yml:5:18"+bad+"a `: ` inside a plain value")
		out, code = reportingCheck(t, runsOnWorkflows(map[string]string{
			"continued.yml": "ubuntu\n      latest now",
			"colon-end.yml": `"ubuntu-latest:"`,
			"wide.yml":      `"café: x"`,
		}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0002.1.2 a value that starts with a reserved character fails naming it", func(t *testing.T) {
		t.Parallel()
		reserved := []string{",", "]", "}", "%", "@", "`"}
		for _, char := range reserved {
			out, code := reportingCheck(t, runsOnWorkflows(map[string]string{"reserved.yml": char + "self-hosted"}))
			acceptance.Sees(t, out, code, 1, ".github/workflows/reserved.yml:5:14"+bad+"a plain value cannot start with '"+char+"'")
			out, code = reportingCheck(t, runsOnWorkflows(map[string]string{"reserved.yml": `"` + char + `self-hosted"`}))
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		}
	})
	t.Run("0002.1.3 a list or a key or a value opened on its key's line fails telling to move it", func(t *testing.T) {
		t.Parallel()
		for _, opener := range []string{"-", "?", ":"} {
			out, code := reportingCheck(t, runsOnWorkflows(map[string]string{"key-line.yml": opener + " self-hosted"}))
			acceptance.Sees(t, out, code, 1, ".github/workflows/key-line.yml:5:14"+bad+"'"+opener+"' opens", "move it onto its own line, indented under the key")
		}
		out, code := reportingCheck(t, runsOnWorkflows(map[string]string{"key-line.yml": "\n      - self-hosted"}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0002.1.4 text after a closing bracket fails saying so", func(t *testing.T) {
		t.Parallel()
		out, code := reportingCheck(t, runsOnWorkflows(map[string]string{"flow-trailing.yml": "[self-hosted]: x", "mapping.yml": "{group: fitness} x"}))
		acceptance.Sees(t, out, code, 1,
			".github/workflows/flow-trailing.yml:5:27"+bad+"text follows the closing bracket",
			".github/workflows/mapping.yml:5:31"+bad+"text follows the closing bracket")
		out, code = reportingCheck(t, runsOnWorkflows(map[string]string{"flow-trailing.yml": "[self-hosted] # a comment"}))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0002_2(t *testing.T) {
	t.Parallel()
	const ran = "ran.txt"
	t.Run("0002.2.1 a staged or pull request change judges only the workflow files it changes", func(t *testing.T) {
		t.Parallel()
		unreported := map[string]string{".github/workflows/old.yml": reportingJob("")}
		fixed := map[string]string{".github/workflows/ci.yml": reportingJob(reportStep)}
		repo := acceptance.Example(t, "happyRepo", unreported)
		acceptance.Write(t, repo, fixed)
		acceptance.Git(t, repo, "add", "-A")
		out, code := acceptance.Fitness(t, repo, nil, "workflow-reporting")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		repo = acceptance.Example(t, "happyRepo", unreported)
		acceptance.Git(t, repo, "update-ref", "refs/remotes/origin/main", "HEAD")
		acceptance.Commit(t, repo, fixed, "ci(report): add the report step")
		out, code = acceptance.Fitness(t, repo, acceptance.PullRequestEnv, "workflow-reporting")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0002.2.2 a step is read and never run", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{".github/workflows/ci.yml": reportingJob("      - run: touch " + ran + "\n" + reportStep)})
		out, code := acceptance.Fitness(t, repo, nil, "workflow-reporting")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if _, err := os.Stat(filepath.Join(repo, ran)); err == nil {
			t.Errorf("a workflow step ran and wrote %s", ran)
		}
	})
}
