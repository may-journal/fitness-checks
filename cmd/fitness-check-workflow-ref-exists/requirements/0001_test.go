// Package requirements proves the workflow-ref-exists check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo, with GitHub answering saved responses.
package requirements

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m, "gh")
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	present := func(t *testing.T) (string, int) {
		out, code, _ := refCheck(t, refWorkflow(
			"may-journal/fitness-runner/.github/workflows/ci-reusable.yml@v1",
			"may-journal/fitness-runner@"+refSha,
		), []acceptance.Response{
			found("repos/may-journal/fitness-runner/git/ref/tags/v1", `{"ref":"refs/tags/v1"}`),
			found("repos/may-journal/fitness-runner/contents/.github/workflows/ci-reusable.yml?ref=v1", `{"name":"ci-reusable.yml"}`),
			found("repos/may-journal/fitness-runner/contents?ref="+refSha, `[{"name":"README.md"},{"name":"action.yml"}]`),
		})
		return out, code
	}
	t.Run("0001.1.2.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := present(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.1.2 it counts the workflow files", func(t *testing.T) {
		t.Parallel()
		out, code := present(t)
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWith(t, out, "workflow-ref-exists", "1")
	})
	missing := func(t *testing.T) (string, int) {
		out, code, _ := refCheck(t, refWorkflow(
			"may-journal/fitness-runner/.github/workflows/comment-check-reusable.yml@v1",
		), []acceptance.Response{
			found("repos/may-journal/fitness-runner/git/ref/tags/v1", `{"ref":"refs/tags/v1"}`),
			notFound("repos/may-journal/fitness-runner/contents/.github/workflows/comment-check-reusable.yml?ref=v1"),
			found("repos/may-journal/fitness-runner", `{"full_name":"may-journal/fitness-runner"}`),
		})
		return out, code
	}
	t.Run("0001.1.2.2.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := missing(t)
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml: line 5:")
	})
	t.Run("0001.1.2.2.2 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := missing(t)
		acceptance.Sees(t, out, code, 1, "ci.yml: line 5: may-journal/fitness-runner/.github/workflows/comment-check-reusable.yml@v1")
	})
	t.Run("0001.1.2.2.3 it fails naming the path", func(t *testing.T) {
		t.Parallel()
		out, code := missing(t)
		acceptance.Sees(t, out, code, 1, ".github/workflows/comment-check-reusable.yml does not exist")
	})
	t.Run("0001.1.2.2.4 it fails naming the ref", func(t *testing.T) {
		t.Parallel()
		out, code := missing(t)
		acceptance.Sees(t, out, code, 1, "does not exist at v1")
	})
	badSha := strings.Repeat("0", 40)
	unknown := func(t *testing.T) (string, int) {
		out, code, _ := refCheck(t, refWorkflow("may-journal/fitness-runner@v9", "may-journal/fitness-runner@"+badSha), []acceptance.Response{
			notFound("repos/may-journal/fitness-runner/git/ref/tags/v9"),
			notFound("repos/may-journal/fitness-runner/git/ref/heads/v9"),
			{Match: []string{"contents?ref=" + badSha}, Stderr: "gh: No commit found for the ref " + badSha + " (HTTP 404)\n", Exit: 1},
			found("repos/may-journal/fitness-runner", `{"full_name":"may-journal/fitness-runner"}`),
		})
		return out, code
	}
	t.Run("0001.1.2.3.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := unknown(t)
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml: line 5:", ".github/workflows/ci.yml: line 7:")
	})
	t.Run("0001.1.2.3.2 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := unknown(t)
		acceptance.Sees(t, out, code, 1, "ci.yml: line 5: may-journal/fitness-runner@v9:", "ci.yml: line 7: may-journal/fitness-runner@"+badSha+":")
	})
	t.Run("0001.1.2.3.3 it fails naming the ref", func(t *testing.T) {
		t.Parallel()
		out, code := unknown(t)
		acceptance.Sees(t, out, code, 1, "v9 is not a tag, a branch, or a full commit sha", "no commit "+badSha+" exists")
	})
	ignored := func(t *testing.T) (string, int, acceptance.Replays) {
		return refCheck(t, refWorkflow(
			"./.github/workflows/local.yml",
			"actions/checkout@v4",
			"may-journal/fitness-runner/.github/workflows/ci-reusable.yml@${{ inputs.ref }}",
			"may-journal/fitness-runner/.github/workflows/ci-reusable.yml@main",
		)+"  script:\n    runs-on: ubuntu-latest\n    steps:\n      - run: |\n          echo 'uses: may-journal/fitness-runner/nope.yml@v1' > caller.yml\n", []acceptance.Response{
			notFound("repos/may-journal/fitness-runner/git/ref/tags/main"),
			found("repos/may-journal/fitness-runner/git/ref/heads/main", `{"ref":"refs/heads/main"}`),
		})
	}
	t.Run("0001.1.2.4.1 it is not read", func(t *testing.T) {
		t.Parallel()
		out, code, rp := ignored(t)
		acceptance.Sees(t, out, code, 0)
		if rp.Called("gh", "contents") || rp.Called("gh", "actions/checkout") {
			t.Errorf("an ignored line read GitHub: %q", rp.Calls("gh"))
		}
	})
	t.Run("0001.1.2.4.2 it passes", func(t *testing.T) {
		t.Parallel()
		out, code, _ := ignored(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.5.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "workflow-ref-exists")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "workflow-ref-exists")
	})
}

// Not parallel: a wall-clock limit beside other tests would time their load.
func Test0001_2(t *testing.T) {
	t.Run("0001.2.1.1.1 its calls are skipped with a note", func(t *testing.T) {
		t.Parallel()
		out, code, _ := refCheck(t, refWorkflow("may-journal/hub/.github/workflows/sync.yml@v1"), []acceptance.Response{
			notFound("repos/may-journal/hub"),
		})
		acceptance.Sees(t, out, code, 0, ".github/workflows/ci.yml: line 5: may-journal/hub/.github/workflows/sync.yml@v1: skipped, may-journal/hub cannot be read with this token")
	})
	refused := func(t *testing.T) (string, int) {
		out, code, _ := refCheck(t, refWorkflow(
			"may-journal/fitness-runner/.github/workflows/nope.yml@v1",
			"may-journal/hub/.github/workflows/sync.yml@v1",
			"may-journal/journey/.github/workflows/build.yml@v1",
		), []acceptance.Response{
			found("repos/may-journal/fitness-runner/git/ref/tags/v1", `{"ref":"refs/tags/v1"}`),
			notFound("repos/may-journal/fitness-runner/contents/.github/workflows/nope.yml?ref=v1"),
			found("repos/may-journal/fitness-runner", `{"full_name":"may-journal/fitness-runner"}`),
			{Match: []string{"repos/may-journal/hub/"}, Stderr: "gh: Resource not accessible by integration (HTTP 403)\n", Exit: 1},
			{Match: []string{"repos/may-journal/journey/"}, Stderr: "error connecting to api.github.com\n", Exit: 1},
		})
		return out, code
	}
	t.Run("0001.2.1.2.1 the findings so far are kept", func(t *testing.T) {
		t.Parallel()
		out, code := refused(t)
		acceptance.Sees(t, out, code, 1, ".github/workflows/ci.yml: line 5: may-journal/fitness-runner/.github/workflows/nope.yml@v1: .github/workflows/nope.yml does not exist at v1")
	})
	t.Run("0001.2.1.2.2 the calls left unverified are noted", func(t *testing.T) {
		t.Parallel()
		out, code := refused(t)
		acceptance.Sees(t, out, code, 1,
			".github/workflows/ci.yml: line 7: may-journal/hub/.github/workflows/sync.yml@v1: skipped, GitHub answered gh: Resource not accessible by integration (HTTP 403)",
			"workflow refs from .github/workflows/ci.yml line 9 on were not verified: gh api repos/may-journal/journey/git/ref/tags/v1: error connecting to api.github.com")
	})
	down := func(t *testing.T) (string, int) {
		out, code, _ := refCheck(t, refWorkflow("may-journal/fitness-runner/.github/workflows/ci-reusable.yml@v1"), []acceptance.Response{
			{Stderr: "error connecting to api.github.com\n", Exit: 1},
		})
		return out, code
	}
	t.Run("0001.2.1.3.1 it passes with info", func(t *testing.T) {
		t.Parallel()
		out, code := down(t)
		acceptance.Sees(t, out, code, 0, "workflow refs were not verified: gh api repos/may-journal/fitness-runner/git/ref/tags/v1: error connecting to api.github.com")
	})
	t.Run("0001.2.1.3.2 it never fails", func(t *testing.T) {
		t.Parallel()
		out, code := down(t)
		acceptance.Sees(t, out, code, 0)
		if strings.Contains(acceptance.Row(out, "workflow-ref-exists"), "failed") {
			t.Errorf("a check whose service cannot answer must not fail, got row %q", acceptance.Row(out, "workflow-ref-exists"))
		}
	})
	t.Run("0001.2.1.4.1 they are read in parallel", func(t *testing.T) {
		// Not parallel: a wall-clock limit beside other tests would time their load.
		var uses []string
		for i := range slowRefs {
			uses = append(uses, fmt.Sprintf("may-journal/fitness-runner/.github/workflows/w%d.yml@%s", i, refSha))
		}
		slow := found("repos/may-journal/fitness-runner/contents/.github/workflows/", `{"name":"w.yml"}`)
		slow.Wait = 1
		out, code, rp := refCheck(t, refWorkflow(uses...), []acceptance.Response{slow})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if reads := len(rp.Calls("gh")); reads != slowRefs {
			t.Errorf("workflow-ref-exists read GitHub %d times, want %d", reads, slowRefs)
		}
		took := rowDuration.FindStringSubmatch(acceptance.Row(out, "workflow-ref-exists"))
		if took == nil {
			t.Fatalf("workflow-ref-exists shows no time in %q", acceptance.Row(out, "workflow-ref-exists"))
		}
		if d, _ := time.ParseDuration(took[1] + took[2]); d >= slowRefs*time.Second/2 {
			t.Errorf("workflow-ref-exists took %v on %d one-second reads, want under %v", d, slowRefs, slowRefs*time.Second/2)
		}
	})
}

// refSha is a full commit sha a workflow pins an action to.
const refSha = "6b02ebef8736d2508d47e085806b5b3cd19efa6c"

// refWorkflow is a workflow file whose jobs call each `uses:` value, one
// job per value, so value i sits on line 4+2i.
func refWorkflow(uses ...string) string {
	var b strings.Builder
	b.WriteString("name: CI\non: push\njobs:\n")
	for i, u := range uses {
		b.WriteString("  job" + string(rune('a'+i)) + ":\n    uses: " + u + "\n")
	}
	return b.String()
}

// found is a saved GitHub answer for a path that exists.
func found(match, body string) acceptance.Response {
	return acceptance.Response{Match: []string{"api " + match}, Stdout: body}
}

// notFound is a saved GitHub answer for a path that does not exist.
func notFound(match string) acceptance.Response {
	return acceptance.Response{Match: []string{"api " + match}, Stderr: "gh: Not Found (HTTP 404)\n", Exit: 1}
}

// refCheck runs workflow-ref-exists on happyRepo with .github/workflows/ci.yml
// set to workflow, GitHub answering from saved responses.
func refCheck(t *testing.T, workflow string, saved []acceptance.Response) (string, int, acceptance.Replays) {
	t.Helper()
	rp := acceptance.StandIn(t, map[string][]acceptance.Response{"gh": saved})
	repo := acceptance.Example(t, "happyRepo", map[string]string{".github/workflows/ci.yml": workflow})
	out, code := acceptance.Fitness(t, repo, rp.Env, "workflow-ref-exists")
	return out, code, rp
}

// rowDuration matches a results-table row's time, in milliseconds or seconds.
var rowDuration = regexp.MustCompile(`│ (\d+(?:\.\d+)?)(ms|s)\s+│\s*$`)

// slowRefs is how many pinned workflow calls the parallel case makes, each answered
// after one second. Read one at a time, they take slowRefs seconds.
const slowRefs = 8
