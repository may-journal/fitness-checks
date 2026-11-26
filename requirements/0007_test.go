package requirements

import (
	"strings"
	"testing"
)

func Test0007_4(t *testing.T) {
	t.Parallel()
	runCases(t, []edgeCase{
		{"repeated-string-literals rejects three uses of a value", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{
				"src/a.ts": rslLines("active", 2),
				"src/b.ts": rslLines("active", 1),
			})
			sees(t, out, code, 1, `"active" appears 3 times (src/a.ts:1, src/a.ts:2, src/b.ts:1) — extract a shared constant`)
		}},
		{"repeated-string-literals passes a value used twice", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.ts": rslLines("active", 2)})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"repeated-string-literals skips comments regex and templates", func(t *testing.T) {
			t.Parallel()
			src := "import x from 'active';\nconst y = require('active');\nexport { z } from 'active';\n" +
				"// 'active'\n/* 'active' */\nconst re = /'active'/g;\nconst s = `'active'`;\n"
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.ts": src})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"repeated-string-literals skips short and idiomatic values", func(t *testing.T) {
			t.Parallel()
			src := rslLines("ok", 3) + rslLines("utf8", 3) + rslLines("inherit", 3) + rslLines("object", 3) + rslLines("use strict", 3)
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.js": src})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
		{"repeated-string-literals lists five places and a count", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.ts": rslLines("active", 7)})
			sees(t, out, code, 1, `"active" appears 7 times (src/a.ts:1, src/a.ts:2, src/a.ts:3, src/a.ts:4, src/a.ts:5, +2 more)`)
		}},
		{"repeated-string-literals reports the most repeated first", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{"src/a.ts": rslLines("twice", 3) + rslLines("thrice", 4)})
			sees(t, out, code, 1, `"thrice" appears 4 times`, `"twice" appears 3 times`)
			if strings.Index(out, `"thrice"`) > strings.Index(out, `"twice"`) {
				t.Errorf("want thrice (4) reported before twice (3):\n%s", out)
			}
		}},
		{"repeated-string-literals counts test and bench files", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{
				"src/a.test.ts":  rslLines("active", 2),
				"src/b.bench.ts": rslLines("active", 1),
			})
			sees(t, out, code, 1, `"active" appears 3 times (src/a.test.ts:1, src/a.test.ts:2, src/b.bench.ts:1)`)
		}},
		{"repeated-string-literals passes an allowed value", func(t *testing.T) {
			t.Parallel()
			out, code := repeatedStringLiterals(t, map[string]string{
				".fitnessrc.json": `{"repeatedStringLiterals": {"allow": ["active"]}}` + "\n",
				"src/a.ts":        rslLines("active", 3),
			})
			sees(t, out, code, 0, "All 1 checks passed")
		}},
	})
}

// Where the check sends each kind of temporary work.
const (
	toPlanIssue  = "open a `Plan`-labeled GitHub Issue and put the plan there, not a file in the repo"
	toPlanDesign = "post it as a comment on the Plan Issue that it designs, not a file in the repo"
	toPRComment  = "post it as a comment on the pull request it reviews, not a file in the repo"
)

// inPlansFolder is the error line for a file in a `plans` folder.
func inPlansFolder(path string) string {
	return path + ": a `plans` folder holds temporary work; " + toPlanIssue
}

// planChecklist is a Plan's checklist as a Markdown file would hold it.
const planChecklist = "# Ship it\n\n## What needs to happen\n\n- [ ] Write the tests.\n"

// inOrder asserts each want appears in out after the one before it.
func inOrder(t *testing.T, out string, want ...string) {
	t.Helper()
	text, at := read(out), 0
	for _, w := range want {
		i := strings.Index(text[at:], w)
		if i < 0 {
			t.Errorf("output is missing %q after position %d, or lists it out of order:\n%s", w, at, out)
			return
		}
		at += i + len(w)
	}
}

// repeatedStringLiterals runs `fitness-install -- repeated-string-literals`
// on happyRepo with files written over it.
func repeatedStringLiterals(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return fitness(t, example(t, "happyRepo", files), nil, "repeated-string-literals")
}

// rslLines returns n source lines that each use the literal value.
func rslLines(value string, n int) string {
	return strings.Repeat("use('"+value+"');\n", n)
}
