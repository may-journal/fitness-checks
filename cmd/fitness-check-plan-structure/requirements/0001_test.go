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
	t.Run("0001.1.1 a pitch plus background plus checklist passes", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructureOf(planStructurePitch, "- [ ] Write the tests."))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
	})
	t.Run("0001.1.2 no blockquote pitch fails", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructureOf("Offline edits survive the next sync.", "- [ ] Keep queued edits."))
		acceptance.Sees(t, out, code, 1, "Missing pitch: add a one-line blockquote (> …) before the first ## section")
	})
	t.Run("0001.1.3 a placeholder pitch fails", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructureOf("> REPLACE-ME", "- [ ] Keep queued edits."))
		acceptance.Sees(t, out, code, 1, "Pitch is still the template placeholder; write the real one-line value")
	})
	t.Run("0001.1.4 no Background section fails", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructurePitch+"\n\n## What needs to happen\n\n- [ ] Keep queued edits.\n")
		acceptance.Sees(t, out, code, 1, "Missing `## Background` section")
	})
	t.Run("0001.1.5 a checkbox only inside code fails", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructureOf(planStructurePitch, "```\n- [ ] Keep queued edits.\n```"))
		acceptance.Sees(t, out, code, 1, "`## What needs to happen` has no checklist item (- [ ] …)")
	})
	t.Run("0001.1.6 an extra section fails", func(t *testing.T) {
		t.Parallel()
		out, code := planStructure(t, planStructureOf(planStructurePitch, "- [ ] Keep queued edits.\n\n## Notes\n\nSee the sync design."))
		acceptance.Sees(t, out, code, 1, "Unexpected `## Notes` section — a plan has only Background and What needs to happen")
	})
	t.Run("0001.1.7 an open questions section fails asking for checklist steps", func(t *testing.T) {
		t.Parallel()
		for _, heading := range []string{"Open questions", "Unresolved questions", "Questions for the owner"} {
			out, code := planStructure(t, planStructureOf(planStructurePitch, "- [ ] Keep queued edits.\n\n## "+heading+"\n\n- Which sync wins?"))
			acceptance.Sees(t, out, code, 1, "Remove the `## "+heading+"` section — turn unknowns into checklist steps")
		}
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "plan-structure")
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	acceptance.PassedWithNoFiles(t, out, "plan-structure")
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
