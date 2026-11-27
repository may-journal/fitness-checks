package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// agedReqs is Adds, a Removed line at 2 and Doubles at 3.
var agedReqs = strings.Replace(twoAcceptances, "2. Doubles", "2. Removed in 2026.11.25.1230: folded into 0001.1\n3. Doubles", 1)

// droppedReqs is agedReqs without its Removed line, Doubles moved up to 2.
var droppedReqs = twoAcceptances

// releases are two released versions newer than the Removed line.
var releases = []string{"2026.11.26.1000", "2026.11.27.1000"}

// agedTests is adder_test.go proving 0001.1 and Doubles at number, a case per step.
func agedTests(number string) string {
	return strings.NewReplacer("0001_2", "0001_"+number, "0001.2.", "0001."+number+".").Replace(provingTests("0001.2.1 add", "0001.2.2 double"))
}

// agedRun commits agedReqs with its tests and the releases, then changes the doc and tests.
// The change is committed, or only staged as a hook sees it.
func agedRun(t *testing.T, versions []string, reqs, tests string, staged bool) (string, int) {
	t.Helper()
	files := map[string]string{
		".fitnessrc.json": `{"policy": "external", "checks": ["requirements"]}` + "\n",
		rootDoc:           acceptance.DocWith(agedReqs),
		"adder_test.go":   agedTests("3"),
	}
	for _, v := range versions {
		files["changelog/"+strings.NewReplacer(".", "-").Replace(v[:15])+"-change.md"] = acceptance.BareEntryText(v, "Feat: change it.")
	}
	repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), files))
	change := map[string]string{rootDoc: acceptance.DocWith(reqs), "adder_test.go": tests}
	if !staged {
		acceptance.Commit(t, repo, change, "docs(requirements): drop a Removed line")
		return judged(t, repo)
	}
	acceptance.Write(t, repo, change)
	acceptance.Git(t, repo, "add", "-A")
	return judged(t, repo)
}

func Test0010_1(t *testing.T) {
	t.Parallel()
	for _, staged := range []bool{false, true} {
		out, code := agedRun(t, releases, droppedReqs, agedTests("2"), staged)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	}
}

func Test0010_2(t *testing.T) {
	t.Parallel()
	t.Run("0010.2.1 a moved acceptance keeping its old test name fails naming the test", func(t *testing.T) {
		t.Parallel()
		for _, staged := range []bool{false, true} {
			out, code := agedRun(t, releases, droppedReqs, agedTests("3"), staged)
			acceptance.Sees(t, out, code, 1, "0001.2")
		}
	})
	t.Run("0010.2.2 a moved acceptance keeping its old step case name fails naming the case", func(t *testing.T) {
		t.Parallel()
		out, code := agedRun(t, releases, droppedReqs, strings.NewReplacer("0001.2.", "0001.3.").Replace(agedTests("2")), false)
		acceptance.Sees(t, out, code, 1, "0001.2.1")
	})
}

func Test0010_3(t *testing.T) {
	t.Parallel()
	out, code := agedRun(t, releases[1:], droppedReqs, agedTests("2"), false)
	acceptance.Sees(t, out, code, 1, rootDoc+":1: 0001.2 was Removed in \"2026.11.25.1230\" and 1 released versions are newer; keep its `2. Removed` line until 2 are",
		"To drop old markers by hand, delete each `Removed` line with 2 or more newer released versions", "rename its test and step cases to match in the same change")
}
