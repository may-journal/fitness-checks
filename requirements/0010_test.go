package requirements

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-checks/req"
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

// repeats are the errors of the doc at path that hold phrase.
func repeats(t *testing.T, path, phrase string) []string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, errs := requirement.Parse(path, string(body))
	return slices.DeleteFunc(errs, func(e string) bool { return !strings.Contains(e, phrase) })
}

// allRepeats are the errors with phrase of every requirement doc.
func allRepeats(t *testing.T, phrase string) []string {
	t.Helper()
	var out []string
	for _, path := range requirementDocs(t) {
		out = append(out, repeats(t, path, phrase)...)
	}
	return out
}

func Test0010_1(t *testing.T) {
	t.Parallel()
	req.Then(t, "0010.1.1.1.1", "b26e7f0f", func(a req.A) {
		a.Equal([]string(nil), allRepeats(t, "open with the same words"))
	})
	req.Then(t, "0010.1.1.1.2", "1e7141f9", func(a req.A) {
		a.Equal([]string(nil), allRepeats(t, "repeats the THEN"))
	})
}
