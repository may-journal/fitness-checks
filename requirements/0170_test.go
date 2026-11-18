package requirements

import (
	"strings"
	"testing"
)

func Test0170_1(t *testing.T) {
	t.Parallel()
	repo := folded0170(t)
	out, code := staged(t, repo, reborn0170)
	sees(t, out, code, 1, "0002 was deleted before", "0002.1 was deleted before")
	out, code = staged(t, repo, map[string]string{renumberMap: "# 0002 Adds Again folded into 0001\n0002 -> 0001\n"})
	sees(t, out, code, 0)
}

func Test0170_2(t *testing.T) {
	t.Parallel()
	out, code := staged(t, folded0170(t), with(reborn0170, map[string]string{renumberMap: "0002.1 -> 0001.2\n"}))
	sees(t, out, code, 1, "0002 was deleted before")
	if strings.Contains(out, "0002.1 was deleted before") {
		t.Errorf("the moved acceptance 0002.1 must be free to reuse:\n%s", out)
	}
}

func Test0170_3(t *testing.T) {
	t.Parallel()
	for body, want := range map[string]string{
		"0001 => 0003\n":               renumberMap + ":1: write each entry as `old -> new`",
		"0002 -> 0001.1\n":             renumberMap + ":1: write each entry as `old -> new`",
		"0002 -> 0001\n0002 -> 0001\n": renumberMap + ":2: 0002 is listed twice",
		"\xef\xbb\xbf0001 -> 0001\n":   renumberMap + ":1: 0001 maps to itself",
		"0002 -> 0001 # folded\n":      renumberMap + ":1: write each entry as `old -> new`, both doc numbers such as 0142 -> 0005 or both acceptance IDs such as 0142.3 -> 0005.2, with comments on their own `#` line",
		"\n0002.1 -> 0001.7\n":         renumberMap + ":2: 0002.1 -> 0001.7 points at 0001.7, which no doc defines",
	} {
		out, code := fitness(t, example(t, "happyRepo", with(pairedRepo(), map[string]string{renumberMap: body})), nil, "requirements")
		sees(t, out, code, 1, want)
	}
}

func Test0170_4(t *testing.T) {
	t.Parallel()
	moved := map[string]string{
		pairedDoc:                                "",
		renumberMap:                              "0001 -> 0003\n",
		"docs/requirements/0003-adds-numbers.md": strings.Replace(docWith(acceptances("0003.1")), "# 0001 ", "# 0003 ", 1),
	}
	out, code := stagedRun(t, pairedRepo(), moved)
	sees(t, out, code, 1, "Test0001_1 names 0001.1, which no requirement defines")
	out, code = stagedRun(t, pairedRepo(), with(moved, map[string]string{pairedTest: "", "reqs/0003_test.go": testFile("reqs", "Test0003_1")}))
	sees(t, out, code, 0)
}

func Test0170_5(t *testing.T) {
	t.Parallel()
	repo := folded0170(t)
	commit(t, repo, with(reborn0170, map[string]string{renumberMap: "0002 -> 0001\n"}), "docs(requirements): add 0002 again")
	commit(t, repo, map[string]string{"docs/requirements/0002-adds-again.md": "", "reqs/0002_test.go": ""}, "docs(requirements): drop 0002")
	out, code := staged(t, repo, reborn0170)
	sees(t, out, code, 1, "0002 was deleted before", "0002.1 was deleted before")
}

// renumberMap is where a repo lists each moved requirement ID.
const renumberMap = "docs/requirements/renumber.txt"

// reborn0170 brings back requirement 0002, with its acceptance 0002.1 and
// its test, after folded0170 removed it.
var reborn0170 = map[string]string{
	"docs/requirements/0002-adds-again.md": twoRepo()["docs/requirements/0002-adds-again.md"],
	"reqs/0002_test.go":                    testFile("reqs", "Test0002_1"),
}

// folded0170 is twoRepo after a commit folds requirement 0002 into 0001 as
// 0001.2, so history shows 0002 and 0002.1 removed.
func folded0170(t *testing.T) string {
	t.Helper()
	repo := example(t, "happyRepo", twoRepo())
	commit(t, repo, map[string]string{
		"docs/requirements/0002-adds-again.md": "",
		"reqs/0001_0002_test.go":               "",
		pairedDoc:                              requirementDoc("0001.1", "0001.2"),
		pairedTest:                             testFile("reqs", "Test0001_1", "Test0001_2"),
	}, "docs(requirements): fold 0002 into 0001")
	return repo
}
