// Package requirements proves the mermaid-legend check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// legend runs the check on happyRepo with files written over it.
func legend(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, append([]string{"mermaid-legend"}, args...)...)
}

// The shared mermaid fixtures from acceptance.
var (
	diagram = acceptance.Diagram
	doc     = acceptance.DiagramDoc
)

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a numbered node with no style class fails", func(t *testing.T) {
		t.Parallel()
		out, code := legend(t, doc(diagram("flowchart TB", "persona((1 Persona))", "classDef persona fill:#eef")))
		req.Then(t, "0001.1.1.1.1", "08479b2a", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "README.md: callout 1 (persona) has no style class")
		})
		req.Then(t, "0001.1.1.1.2", "0c45cb5f", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "(diagram at line 9)")
		})
	})
	t.Run("0001.1.2 a diagram with numbered callouts has no classDef legend and fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "6e35af9d", func(_ req.A) {
			out, code := legend(t, doc(diagram("flowchart TB", "persona((1 Persona)):::persona")))
			acceptance.Sees(t, out, code, 1, "README.md: diagram at line 9 has numbered callouts but no classDef legend")
		})
	})
	t.Run("0001.1.3 inline class or style callouts pass", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.3.1", "ebffcd6f", func(_ req.A) {
			out, code := legend(t, doc(diagram(
				"flowchart TB",
				"persona((1 Persona)):::persona",
				"screen[2 Screen]",
				"db[(3 Store)]",
				"class screen, db thing",
				"dec{4 Choice}",
				"style dec fill:#eef",
				"classDef persona fill:#eef",
				"classDef thing fill:#f8f")))
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.4 a non flowchart diagram passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.4.1", "17e1456d", func(_ req.A) {
			out, code := legend(t, doc(diagram("block-beta", "a[1 Store]", "b[2 Screen]")))
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.5 labels without a separate number pass", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.5.1", "3e6d05a6", func(_ req.A) {
			out, code := legend(t, doc(diagram("graph LR", "screen[Screen 2]", "persona((1Persona))")))
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.6 a broken diagram is named by its own line", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.6.1", "3e057018", func(_ req.A) {
			out, code := legend(t, doc(
				diagram("flowchart TB", "a[1 A]:::x", "classDef x fill:#eef"),
				diagram("flowchart TB", "b[2 B]:::x", "classDef y fill:#eef", "c[3 C]")))
			acceptance.Sees(t, out, code, 1, "README.md: callout 3 (c) has no style class (diagram at line 15)")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a PR body unstyled callout fails naming the description", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.2.1.1", "41f704b4", func(_ req.A) {
			out, code := legend(t, map[string]string{"body.md": diagram(
				"flowchart TB", "persona((1 Persona))", "classDef persona fill:#eef")},
				"--body-file", "body.md")
			acceptance.Sees(t, out, code, 1, "(description): callout 1 (persona) has no style class (diagram at line 1)")
		})
	})
	t.Run("0001.2.2 docs with no diagram pass with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.2.2.1", "ebe8880b", func(_ req.A) {
			out, code := legend(t, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "mermaid-legend")
		})
	})
}
