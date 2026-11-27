package requirements

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// syncedScript is a shell script the org sync writes into the example repo.
const syncedScript = "scripts/static-preview-link.sh"

func Test0003_1(t *testing.T) {
	t.Parallel()
	given := acceptance.With(acceptance.GoRepo(), map[string]string{
		syncedScript:                "echo one\n",
		".github/synced-files.json": syncRecord("echo one\n"),
	})
	t.Run("0003.1.1 a synced script matching its hash passes with no test or doc", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0003.1.2.1.1", "2701e931", func(_ req.A) {
			out, code := staged(t, given, map[string]string{syncedScript: "echo two\n", ".github/synced-files.json": syncRecord("echo two\n")})
			acceptance.Sees(t, out, code, 0, passed)
		})
	})
	t.Run("0003.1.2 a synced script edited by hand fails naming the file and its source", func(t *testing.T) {
		t.Parallel()
		out, code := staged(t, given, map[string]string{syncedScript: "echo by hand\n"})
		req.Then(t, "0003.1.2.2.1", "f7956a3e", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, syncedScript+" is synced from")
		})
		req.Then(t, "0003.1.2.2.2", "5f469d81", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "is synced from may-journal/.github and was edited by hand; change it there instead")
		})
	})
}

// syncRecord is the sync record naming syncedScript with body as its content.
func syncRecord(body string) string {
	blob := sha1.Sum(fmt.Appendf(nil, "blob %d\x00%s", len(body), body))
	return fmt.Sprintf("{\n  \"source\": \"may-journal/.github\",\n  \"files\": {\n    %q: \"%x\"\n  }\n}\n", syncedScript, blob)
}
