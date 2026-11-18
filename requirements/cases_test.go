package requirements

import "testing"

// edgeCase is one case of an acceptance's table test: an edge case the
// acceptance covers, with the test that proves it.
type edgeCase struct {
	name string
	run  func(t *testing.T)
}

// runCases runs each case under its own name, inside the acceptance test.
func runCases(t *testing.T, cases []edgeCase) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
