// Package requirements proves the plan-structure check's docs/requirements as a user
// meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"os"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

const planStructurePitch = "> Offline edits survive the next sync."

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1.1.1 `plan-structure` passes it", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructureOf(planStructurePitch, "- [ ] Write the tests."))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
	})
	noPitch := planStructureOf("Offline edits survive the next sync.", "- [ ] Keep queued edits.")
	t.Run("0001.1.1.2.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, noPitch)
		acceptance.Sees(t, out, code, 1)
	})
	t.Run("0001.1.1.2.2 it asks for a one-line pitch", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, noPitch)
		acceptance.Sees(t, out, code, 1, "Missing pitch: add a one-line blockquote (> …) before the first ## section")
	})
	placeholder := planStructureOf("> REPLACE-ME", "- [ ] Keep queued edits.")
	t.Run("0001.1.1.3.1 it fails", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, placeholder)
		acceptance.Sees(t, out, code, 1)
	})
	t.Run("0001.1.1.3.2 it asks for the real pitch", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, placeholder)
		acceptance.Sees(t, out, code, 1, "Pitch is still the template placeholder; write the real one-line value")
	})
	t.Run("0001.1.1.4.1 it fails naming the missing section", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructurePitch+"\n\n## What needs to happen\n\n- [ ] Keep queued edits.\n")
		acceptance.Sees(t, out, code, 1, "Missing `## Background` section")
	})
	t.Run("0001.1.1.5.1 it fails because the checklist has no item", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructureOf(planStructurePitch, "```\n- [ ] Keep queued edits.\n```"))
		acceptance.Sees(t, out, code, 1, "`## What needs to happen` has no checklist item (- [ ] …)")
	})
	t.Run("0001.1.1.6.1 it fails naming the unexpected section", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructureOf(planStructurePitch, "- [ ] Keep queued edits.\n\n## Notes\n\nSee the sync design."))
		acceptance.Sees(t, out, code, 1, "Unexpected `## Notes` section — a plan has only Background and What needs to happen")
	})
	questions := func(heading string) string {
		return planStructureOf(planStructurePitch, "- [ ] Keep queued edits.\n\n## "+heading+"\n\n- Which sync wins?")
	}
	headings := []string{"Open questions", "Unresolved questions", "Questions for the owner"}
	t.Run("0001.1.1.7.1 it fails", func(t *testing.T) {
		t.Parallel()
		for _, heading := range headings {
			out, code := planStructure(t, questions(heading))
			acceptance.Sees(t, out, code, 1)
		}
	})
	t.Run("0001.1.1.7.2 it asks to turn unknowns into checklist steps", func(t *testing.T) {
		t.Parallel()
		for _, heading := range headings {
			out, code := planStructure(t, questions(heading))
			acceptance.Sees(t, out, code, 1, "Remove the `## "+heading+"` section — turn unknowns into checklist steps")
		}
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	run := func(t *testing.T) (string, int) {
		return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "plan-structure")
	}
	t.Run("0001.2.1.1.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := run(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.1.1.2 it judges zero files", func(t *testing.T) {
		t.Parallel()
		out, code := run(t)
		acceptance.Sees(t, out, code, 0)
		acceptance.PassedWithNoFiles(t, out, "plan-structure")
	})
}

// planStructure runs `fitness-install -- plan-structure --body-file plan.md`
// on happyRepo with body as plan.md.
func planStructure(t *testing.T, body string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{"plan.md": body})
	return acceptance.Fitness(t, repo, nil, "plan-structure", "--body-file", "plan.md")
}

// planStructureOf is a plan with pitch, then Background, then the steps.
func planStructureOf(pitch, steps string) string {
	return pitch + "\n\n## Background\n\nThe sync job drops edits made offline.\n\n## What needs to happen\n\n" + steps + "\n"
}
