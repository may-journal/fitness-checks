package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func Test0008_1(t *testing.T) {
	t.Parallel()
	t.Run("0008.1.1 a lowercase and or but outside code or quotes asks to split", func(t *testing.T) {
		t.Parallel()
		gherkinFails(t, oneAcceptance("GIVEN two numbers", "WHEN I add them", "THEN each part stands alone and code spans drop out"), `20: split this item at its "and"`)
		gherkinFails(t, oneAcceptance("GIVEN two numbers", "WHEN I add them", "THEN it is a sum but not a product"), `20: split this item at its "but"`)
	})
	t.Run("0008.1.2 the word in code or in quotes passes", func(t *testing.T) {
		t.Parallel()
		reqs := oneAcceptance("GIVEN two numbers such as `1 and 2`", "WHEN I add them", `THEN I get the word "and" or “but”`)
		out, code := judged(t, gherkinRepo(t, reqs, acceptance.TestFile("adder", "Test0001_1")))
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}
