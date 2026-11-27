package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// requirementDocs are the repo's requirement docs, the root's and each
// package's.
func requirementDocs(t *testing.T) []string {
	t.Helper()
	var paths []string
	for _, pattern := range []string{"../docs/requirements/*.md", "../cmd/*/docs/requirements/*.md"} {
		found, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatal(err)
		}
		paths = append(paths, found...)
	}
	if len(paths) == 0 {
		t.Fatal("no requirement docs found")
	}
	return paths
}

// repeats are the errors that name an alike WHEN or a repeated THEN in the
// doc at path.
func repeats(t *testing.T, path string) []string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, errs := requirement.Parse(path, string(body))
	var out []string
	for _, e := range errs {
		if strings.Contains(e, "open with the same words") || strings.Contains(e, "repeats the THEN") {
			out = append(out, e)
		}
	}
	return out
}

func Test0010_1(t *testing.T) {
	t.Parallel()
	for _, path := range requirementDocs(t) {
		for _, e := range repeats(t, path) {
			t.Error(e)
		}
	}
}
