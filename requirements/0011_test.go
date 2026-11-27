package requirements

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func Test0011_1(t *testing.T) {
	t.Parallel()
	const words = "it fails naming the file"
	req.Then(t, "0011.1.1.1.1", "bbaff244", func(a req.A) {
		sum := sha256.Sum256([]byte(words))
		a.Equal(hex.EncodeToString(sum[:])[:8], req.Hash(words))
	})
	req.Then(t, "0011.1.1.1.2", "41534bdb", func(a req.A) {
		a.Equal(req.Hash(words), req.Hash("  it fails\n\tnaming  the file "))
	})
}

// proof runs a Go test in a module holding a copy of the req package, and
// returns the output and whether it passed.
func proof(t *testing.T, body string) (string, bool) {
	t.Helper()
	source, err := os.ReadFile("../req/req.go")
	acceptance.MustDo(t, err)
	digest, err := os.ReadFile("../req/digest/digest.go")
	acceptance.MustDo(t, err)
	dir := t.TempDir()
	acceptance.Write(t, dir, map[string]string{
		"go.mod":               "module example.com/proof\n\ngo 1.26\n",
		"req/req.go":           strings.ReplaceAll(string(source), "github.com/may-journal/fitness-checks/req/digest", "example.com/proof/req/digest"),
		"req/digest/digest.go": string(digest),
		"proof_test.go":        "package proof\n\nimport (\n\t\"testing\"\n\n\t\"example.com/proof/req\"\n)\n\nfunc TestProof(t *testing.T) {\n" + body + "}\n",
	})
	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err == nil
}

func Test0011_2(t *testing.T) {
	t.Parallel()
	t.Run("0011.2.1 a body with no assertion fails", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0011.2.1.1.1", "30e33461", func(a req.A) {
			out, passed := proof(t, "\treq.Then(t, \"0001.1.1.1.1\", \"00000000\", func(a req.A) {})\n")
			a.True(!passed, "a body with no assertion passed:\n%s", out)
			a.Contains(out, "THEN 0001.1.1.1.1#00000000 makes no assertion")
		})
	})
	t.Run("0011.2.2 a body that asserts passes", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0011.2.1.2.1", "19dc2572", func(a req.A) {
			out, passed := proof(t, "\treq.Then(t, \"0001.1.1.1.1\", \"00000000\", func(a req.A) { a.Equal(1, 1) })\n")
			a.True(passed, "a body with an assertion failed:\n%s", out)
		})
	})
}
