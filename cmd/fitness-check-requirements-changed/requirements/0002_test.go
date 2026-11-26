package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func Test0002_1(t *testing.T) {
	t.Parallel()
	edit := map[string]string{"pkgs/a/a.go": acceptance.Monorepo(true)["pkgs/a/a.go"] + "\n// Two is two.\nconst Two = 2\n"}
	pair := func(doc, test string) map[string]string {
		return acceptance.With(edit, map[string]string{doc: acceptance.RequirementDoc("0001.1", "0001.2"), test: acceptance.TestFile("requirements", "Test0001_1", "Test0001_2")})
	}
	t.Run("0002.1.1 a package's code beside another package's test plus doc fails", func(t *testing.T) {
		t.Parallel()
		out, code := staged(t, acceptance.Monorepo(true), pair(pairedDoc, "adder_test.go"))
		acceptance.Sees(t, out, code, 1, "pkgs/a/a.go changed without a test and a requirement doc")
	})
	t.Run("0002.1.2 a package's code beside its own test plus doc passes", func(t *testing.T) {
		t.Parallel()
		out, code := staged(t, acceptance.Monorepo(true), pair(acceptance.PackageDoc("a"), "pkgs/a/requirements/0001_test.go"))
		acceptance.Sees(t, out, code, 0, passed)
	})
	t.Run("0002.1.3 a package file deleted while a doc names it fails naming the file with the doc", func(t *testing.T) {
		t.Parallel()
		naming := acceptance.With(acceptance.Monorepo(true), map[string]string{pairedDoc: strings.Replace(acceptance.RequirementDoc("0001.1"), "by hand.", "by hand in pkgs/a/a.go.", 1)})
		out, code := staged(t, naming, map[string]string{"pkgs/a/a.go": ""})
		acceptance.Sees(t, out, code, 1, "pkgs/a/a.go was deleted, leaving requirement doc "+pairedDoc)
	})
}
