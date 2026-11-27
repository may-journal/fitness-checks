// Package requirements proves the go-complexity check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo with Go files written over it.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a function over the branch ceiling fails naming file, line and score", func(t *testing.T) {
		t.Parallel()
		out, code := goComplexity(t, map[string]string{"a.go": complexityFile("func f()", 6)})
		req.Then(t, "0001.1.1.1.1", "f10ceae9", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "a.go:5:")
		})
		req.Then(t, "0001.1.1.1.2", "17c2cd47", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "func f has a complexity of 7")
		})
		req.Then(t, "0001.1.1.1.3", "8a50eea6", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "maximum allowed is 5")
		})
	})
	t.Run("0001.1.2 each of if or for or range or switch or select or && or || adds one point", func(t *testing.T) {
		t.Parallel()
		body := "package p\n\nfunc f(n int, ch chan int, a, b bool) bool {\n" +
			"\tswitch n {\n\tcase 1:\n\tcase 2:\n\tdefault:\n\t}\n" +
			"\tselect {\n\tcase <-ch:\n\tdefault:\n\t}\n" +
			"\tfor n > 0 {\n\t\tn--\n\t}\n" +
			"\tfor range ch {\n\t}\n" +
			"\treturn a && b || a\n}\n"
		out, code := goComplexity(t, map[string]string{"a.go": body})
		req.Then(t, "0001.1.1.2.1", "0e26e13f", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "a.go:3: func f has a complexity of 8; maximum allowed is 5")
		})
		req.Then(t, "0001.1.1.2.2", "f327bd81", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "func f has a complexity of 8")
		})
	})
	t.Run("0001.1.3 a function literal is scored alone", func(t *testing.T) {
		t.Parallel()
		body := "package p\n\nvar ok bool\n\nfunc outer() {\n\tg := func() {\n" + complexityIfs(5) + "\t}\n\tg()\n}\n"
		out, code := goComplexity(t, map[string]string{"a.go": body})
		req.Then(t, "0001.1.1.3.1", "44440599", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "a.go:6: function literal has a complexity of 6; maximum allowed is 5")
		})
		req.Then(t, "0001.1.1.3.2", "41950786", func(a req.A) {
			a.True(!strings.Contains(out, "func outer"), "outer must not carry the literal's branches:\n%s", out)
		})
	})
	t.Run("0001.1.4 a method with too many branches fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.4.1", "d0815e91", func(_ req.A) {
			file := "package p\n\ntype T struct{}\n\n" + strings.TrimPrefix(complexityFile("func (T) m()", 5), "package p\n\n")
			complexitySees(t, map[string]string{"a.go": file}, 1, "a.go:7: method m has a complexity of 6; maximum allowed is 5")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a test file is judged like any file", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.1.1", "73408155", func(_ req.A) {
			complexitySees(t, map[string]string{"a_test.go": complexityFile("func TestF()", 5)}, 1, "a_test.go:5: func TestF has a complexity of 6; maximum allowed is 5")
		})
	})
	t.Run("0001.2.2 the repo config sets the ceiling", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.2.1", "e7d3fadb", func(_ req.A) {
			out, code := goComplexity(t, map[string]string{
				".fitnessrc.json": `{"goComplexity": {"max": 2}}` + "\n",
				"a.go":            complexityFile("func f()", 2),
			})
			acceptance.Sees(t, out, code, 1, "a.go:5: func f has a complexity of 3; maximum allowed is 2")
		})
	})
	t.Run("0001.2.3 a Go file that does not parse fails with the parser's error", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.3.1", "562f823b", func(_ req.A) {
			out, code := goComplexity(t, map[string]string{"bad.go": "package p\nfunc {"})
			acceptance.Sees(t, out, code, 1, "✖ bad.go:")
			if !strings.Contains(complexityJoined(out), complexityJoined("bad.go:2:6: expected 'IDENT', found '{'")) {
				t.Errorf("output is missing the parser's error:\n%s", out)
			}
		})
	})
	t.Run("0001.2.4 a repo with no Go files passes without files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.2.1.4.1", "857dd667", func(_ req.A) {
			out, code := goComplexity(t, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "go-complexity")
		})
	})
}

// goComplexity runs `fitness-install -- go-complexity` on happyRepo with
// files written over it.
func goComplexity(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "go-complexity")
}

// complexityIfs returns n `if` branches, one point each.
func complexityIfs(n int) string {
	return strings.Repeat("\tif ok {\n\t}\n", n)
}

// complexityFile is a Go file holding decl and its body of n branches.
func complexityFile(decl string, n int) string {
	return "package p\n\nvar ok bool\n\n" + decl + " {\n" + complexityIfs(n) + "}\n"
}

// complexityJoined drops the table borders and every space, so a long
// message the table wrapped mid-word reads whole.
func complexityJoined(out string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(out, "│", " ")), "")
}

// complexitySees runs go-complexity on files and asserts its exit code and message.
func complexitySees(t *testing.T, files map[string]string, want int, msg string) {
	t.Helper()
	out, code := goComplexity(t, files)
	acceptance.Sees(t, out, code, want, msg)
}
