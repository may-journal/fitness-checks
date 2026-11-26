// Package requirements proves the go-complexity check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo with Go files written over it.
package requirements

import (
	"os"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	over := func(t *testing.T) (string, int) {
		return goComplexity(t, map[string]string{"a.go": complexityFile("func f()", 6)})
	}
	t.Run("0001.1.1.1.1 it fails naming `file:line`", func(t *testing.T) {
		t.Parallel()
		out, code := over(t)
		acceptance.Sees(t, out, code, 1, "a.go:5: func f")
	})
	t.Run("0001.1.1.1.2 it names its complexity", func(t *testing.T) {
		t.Parallel()
		out, code := over(t)
		acceptance.Sees(t, out, code, 1, "has a complexity of 7")
	})
	t.Run("0001.1.1.1.3 it names the maximum allowed", func(t *testing.T) {
		t.Parallel()
		out, code := over(t)
		acceptance.Sees(t, out, code, 1, "maximum allowed is 5")
	})
	t.Run("0001.1.1.2.1 each adds one point to a start of 1", func(t *testing.T) {
		t.Parallel()
		body := "package p\n\nfunc f(n int, ch chan int, a, b bool) bool {\n" +
			"\tswitch n {\n\tcase 1:\n\tcase 2:\n\tdefault:\n\t}\n" +
			"\tselect {\n\tcase <-ch:\n\tdefault:\n\t}\n" +
			"\tfor n > 0 {\n\t\tn--\n\t}\n" +
			"\tfor range ch {\n\t}\n" +
			"\treturn a && b || a\n}\n"
		out, code := goComplexity(t, map[string]string{"a.go": body})
		acceptance.Sees(t, out, code, 1, "a.go:3: func f has a complexity of 8; maximum allowed is 5")
	})
	t.Run("0001.1.1.2.2 default cases add none", func(t *testing.T) {
		t.Parallel()
		body := "package p\n\nvar ok bool\n\nfunc f(n int, ch chan int) {\n" +
			"\tswitch n {\n\tdefault:\n\t}\n" +
			"\tselect {\n\tdefault:\n\t}\n" + complexityIfs(5) + "}\n"
		out, code := goComplexity(t, map[string]string{"a.go": body})
		acceptance.Sees(t, out, code, 1, "a.go:5: func f has a complexity of 6; maximum allowed is 5")
	})
	literal := func(t *testing.T) (string, int) {
		body := "package p\n\nvar ok bool\n\nfunc outer() {\n\tg := func() {\n" + complexityIfs(5) + "\t}\n\tg()\n}\n"
		return goComplexity(t, map[string]string{"a.go": body})
	}
	t.Run("0001.1.1.3.1 it scores alone", func(t *testing.T) {
		t.Parallel()
		out, code := literal(t)
		acceptance.Sees(t, out, code, 1, "a.go:6: function literal has a complexity of 6; maximum allowed is 5")
	})
	t.Run("0001.1.1.3.2 the enclosing function does not carry its branches", func(t *testing.T) {
		t.Parallel()
		out, code := literal(t)
		acceptance.Sees(t, out, code, 1)
		if strings.Contains(out, "func outer") {
			t.Errorf("outer must not carry the literal's branches:\n%s", out)
		}
	})
	t.Run("0001.1.1.4.1 it fails reported as `method <name>`", func(t *testing.T) {
		t.Parallel()
		file := "package p\n\ntype T struct{}\n\n" + strings.TrimPrefix(complexityFile("func (T) m()", 5), "package p\n\n")
		out, code := goComplexity(t, map[string]string{"a.go": file})
		acceptance.Sees(t, out, code, 1, "a.go:7: method m has a complexity of 6; maximum allowed is 5")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1.1.1 it fails like any other Go file", func(t *testing.T) {
		t.Parallel()
		out, code := goComplexity(t, map[string]string{"a_test.go": complexityFile("func TestF()", 5)})
		acceptance.Sees(t, out, code, 1, "a_test.go:5: func TestF has a complexity of 6; maximum allowed is 5")
	})
	t.Run("0001.2.1.2.1 that ceiling judges instead of 5", func(t *testing.T) {
		t.Parallel()
		out, code := goComplexity(t, map[string]string{
			".fitnessrc.json": `{"goComplexity": {"max": 2}}` + "\n",
			"a.go":            complexityFile("func f()", 2),
		})
		acceptance.Sees(t, out, code, 1, "a.go:5: func f has a complexity of 3; maximum allowed is 2")
	})
	t.Run("0001.2.1.3.1 it fails with the parser's error", func(t *testing.T) {
		t.Parallel()
		out, code := goComplexity(t, map[string]string{"bad.go": "package p\nfunc {"})
		acceptance.Sees(t, out, code, 1, "✖ bad.go:")
		if !strings.Contains(complexityJoined(out), complexityJoined("bad.go:2:6: expected 'IDENT', found '{'")) {
			t.Errorf("output is missing the parser's error:\n%s", out)
		}
	})
	t.Run("0001.2.1.4.1 it passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := goComplexity(t, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "go-complexity")
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
