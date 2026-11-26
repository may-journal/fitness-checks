package gherkin

import (
	"fmt"
	"math"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// Down rewrites the doc's numbered acceptances as bullets and drops Removed
// lines. Extra GIVENs and THENs join theirs with "and", as a bullet holds one.
func Down(root, path string) error {
	return rewrite(root, path, false, bullets)
}

// Bullets is doc, found at path, with its numbered acceptances read back as
// bullets, and whether it held them. A doc moved to the numbered form then
// compares with its bullet form word for word.
func Bullets(path, doc string) (string, bool) {
	d, _ := requirement.ParseWith(path, doc, math.MaxInt)
	if !d.Gherkin {
		return doc, false
	}
	return replaced(doc, bullets(d)), true
}

// bullets renders d's acceptances as one Given with a single step or a
// numbered chain.
func bullets(d requirement.Document) string {
	var out string
	for _, a := range d.Acceptances {
		out += "- " + a.ID + "\n  - Given " + a.Given + "\n"
		if len(a.Steps) == 1 {
			out += "    - When " + a.Steps[0].When + "\n      - Then " + a.Steps[0].Then + "\n"
			continue
		}
		for i, st := range a.Steps {
			out += fmt.Sprintf("    %d. When %s\n       - Then %s\n", i+1, st.When, st.Then)
		}
	}
	return out
}
