// Package requirements proves the plan-trailer check's docs/requirements as a user
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

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a message with no plan reference passes", func(t *testing.T) {
		t.Parallel()
		out, code := planTrailer(t, "feat(app): add a feature\n\nSome body.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2 a correct trailer passes", func(t *testing.T) {
		t.Parallel()
		out, code := planTrailer(t, "feat(app): add a feature\n\nSome body.\n\nPlan #63\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.3 each misspelled plan line fails with the fix", func(t *testing.T) {
		t.Parallel()
		out, code := planTrailer(t, "feat(app): add a feature\n\nPlan 63\nplan #63\nPlan: 63\n")
		acceptance.Sees(t, out, code, 1,
			`plan trailer "Plan 63" `+planTrailerFix,
			`plan trailer "plan #63" `+planTrailerFix,
			`plan trailer "Plan: 63" `+planTrailerFix)
	})
	t.Run("0001.1.4 body prose starting with Plan passes", func(t *testing.T) {
		t.Parallel()
		out, code := planTrailer(t, "feat(app): add a feature\n\nPlan the rollout in 3 steps.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.5 no message supplied reads HEAD", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", nil)
		acceptance.Commit(t, repo, nil, "feat(app): add a feature\n\nPlan: 63")
		out, code := acceptance.Fitness(t, repo, nil, "plan-trailer")
		acceptance.Sees(t, out, code, 1, `plan trailer "Plan: 63" `+planTrailerFix)
	})
	t.Run("0001.1.6 a plan of another repo passes", func(t *testing.T) {
		t.Parallel()
		out, code := planTrailer(t, "feat(app): add a feature\n\nPlan may-journal/.github#144\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = planTrailer(t, "feat(app): add a feature\n\nPlan https://github.com/may-journal/.github/issues/144\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.7 a misspelled plan line of another repo fails", func(t *testing.T) {
		t.Parallel()
		out, code := planTrailer(t, "feat(app): add a feature\n\nPlan may-journal/.github 144\nPlan: may-journal/.github#144\n")
		acceptance.Sees(t, out, code, 1,
			`plan trailer "Plan may-journal/.github 144" `+planTrailerFix,
			`plan trailer "Plan: may-journal/.github#144" `+planTrailerFix)
	})
}

// planTrailer runs `fitness-install -- plan-trailer --message <msg>` on
// happyRepo, as the commit-msg hook passes a proposed message.
func planTrailer(t *testing.T, msg string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "plan-trailer", "--message", msg)
}

// planTrailerFix is the hint every misspelled plan reference carries.
const planTrailerFix = `must read exactly "Plan #<number>", "Plan <owner>/<repo>#<number>" or "Plan https://github.com/<owner>/<repo>/issues/<number>"`
