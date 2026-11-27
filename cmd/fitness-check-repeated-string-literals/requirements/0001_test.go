// Package requirements proves the repeated-string-literals check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

// repeated runs `fitness-install -- repeated-string-literals` on happyRepo
// with files written over it.
func repeated(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "repeated-string-literals")
}

// uses is a source file that holds each value on its own line.
func uses(values ...string) string {
	var b strings.Builder
	for i, v := range values {
		fmt.Fprintf(&b, "export const v%d = %s;\n", i, v)
	}
	return b.String()
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a value used three times fails listing each location and asking for a shared constant", func(t *testing.T) {
		t.Parallel()
		out, code := repeated(t, map[string]string{
			"src/a.ts": "if (order.status === 'active') setStatus('active');\n",
			"src/b.ts": "export const rows = all.filter((r) => r.status === 'active');\n",
		})
		req.Then(t, "0001.1.1.1.1", "915c530c", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, `"active" appears 3 times`)
		})
		req.Then(t, "0001.1.1.1.2", "d51852bd", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, `(src/a.ts:1, src/a.ts:1, src/b.ts:1)`)
		})
		req.Then(t, "0001.1.1.1.3", "43bc551f", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, `— extract a shared constant`)
		})
	})
	t.Run("0001.1.2 a value used twice passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "07b8e753", func(_ req.A) {
			out, code := repeated(t, map[string]string{"src/a.ts": uses("'active'", "'active'")})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.3 a short value is ignored", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.3.1", "d958af26", func(_ req.A) {
			out, code := repeated(t, map[string]string{"src/a.ts": uses("'ab'", "'ab'", "'ab'", "'ab'")})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.1.4 a test file counts like any source", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.4.1", "d605117c", func(_ req.A) {
			out, code := repeated(t, map[string]string{
				"src/a.ts":      uses("'pending'"),
				"src/a.test.ts": uses("'pending'", "'pending'"),
			})
			acceptance.Sees(t, out, code, 1, `"pending" appears 3 times (src/a.test.ts:1, src/a.test.ts:2, src/a.ts:1)`)
		})
	})
	t.Run("0001.1.5 the most repeated value comes first", func(t *testing.T) {
		t.Parallel()
		out, code := repeated(t, map[string]string{"src/a.ts": uses("'alpha'", "'alpha'", "'alpha'", "'beta'", "'beta'", "'beta'", "'beta'")})
		req.Then(t, "0001.1.1.5.1", "6e5e9f66", func(a req.A) {
			acceptance.Sees(t, out, code, 1, `"beta" appears 4 times`, `"alpha" appears 3 times`)
			a.True(strings.Index(acceptance.Read(out), `"beta" appears`) < strings.Index(acceptance.Read(out), `"alpha" appears`), "the more repeated value must come first:\n%s", out)
		})
	})
	t.Run("0001.1.6 locations stop at five then the rest are counted", func(t *testing.T) {
		t.Parallel()
		out, code := repeated(t, map[string]string{"src/a.ts": uses("'many'", "'many'", "'many'", "'many'", "'many'", "'many'", "'many'")})
		req.Then(t, "0001.1.1.6.1", "32952591", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, `"many" appears 7 times (src/a.ts:1, src/a.ts:2, src/a.ts:3, src/a.ts:4, src/a.ts:5,`)
		})
		req.Then(t, "0001.1.1.6.2", "e83cc100", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, `, +2 more)`)
		})
	})
	t.Run("0001.1.7 a repo with no sources passes with zero files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.7.1", "91ecdad2", func(_ req.A) {
			out, code := repeated(t, nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "repeated-string-literals")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	req.Then(t, "0001.2.1.1.1", "07b8e753", func(_ req.A) {
		src := "import a from './shared';\nimport b from './shared';\nconst c = require('./shared');\nexport { a, b, c };\n"
		out, code := repeated(t, map[string]string{"src/a.ts": src})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		src = "export const a = typeof x === 'object';\nexport const b = typeof y === 'object';\nexport const c = typeof z === 'object';\n" +
			"export const d = buf.toString('utf8');\nexport const e = buf.toString('utf8');\nexport const f = buf.toString('utf8');\n"
		out, code = repeated(t, map[string]string{"src/a.ts": src})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		src = "export const a = `status`;\nexport const b = `status`;\nexport const c = `status`;\n" +
			"export const d = /status/;\nexport const e = /status/;\nexport const f = /status/;\n" +
			"// 'status' 'status' 'status'\n"
		out, code = repeated(t, map[string]string{"src/a.ts": src})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		".fitnessrc.json": `{"repeatedStringLiterals": {"allow": ["kindValue"]}}` + "\n",
		"src/a.ts":        uses("'kindValue'", "'kindValue'", "'kindValue'", "'otherValue'", "'otherValue'", "'otherValue'"),
	}
	t.Run("0001.3.1 a listed value passes", func(t *testing.T) {
		t.Parallel()
		out, code := repeated(t, files)
		req.Then(t, "0001.3.1.1.1", "07b8e753", func(a req.A) {
			acceptance.Sees(t, out, code, 1)
			a.True(!strings.Contains(out, `"kindValue" appears`), "an allowed value must pass:\n%s", out)
		})
	})
	t.Run("0001.3.2 another value still fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.3.1.2.1", "fb54b835", func(_ req.A) {
			out, code := repeated(t, files)
			acceptance.Sees(t, out, code, 1, `"otherValue" appears 3 times`)
		})
	})
}
