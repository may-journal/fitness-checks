// Package requirements proves the mermaid-level-bleed check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"strconv"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// bleed runs the check on happyRepo with files written over it.
func bleed(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, "mermaid-level-bleed")
}

// doc is a C4 level doc with one diagram and a callout table of
// [#, Description] rows.
func doc(descriptions ...string) string {
	lines := []string{"# Level", "", "```mermaid", `Rel(a, b, "1")`, "```", "", "| # | Description |", "| --- | --- |"}
	for i, d := range descriptions {
		lines = append(lines, "| "+strconv.Itoa(i+1)+" | "+d+" |")
	}
	return strings.Join(lines, "\n") + "\n"
}

// repeat is the message for file repeating text from level.
func repeat(file, text string, level int) string {
	return file + `: description "` + text + `" repeats level ` + strconv.Itoa(level) +
		" verbatim — lower levels should add level-specific rationale"
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a level repeating its parent description fails", func(t *testing.T) {
		t.Parallel()
		out, code := bleed(t, map[string]string{
			"docs/architecture/01-context.md":    doc("User reads the timeline"),
			"docs/architecture/02-containers.md": doc("User reads the timeline", "New detail"),
		})
		acceptance.Sees(t, out, code, 1, repeat("docs/architecture/02-containers.md", "user reads the timeline", 1))
	})
	t.Run("0001.1.2 a case or spacing repeat is reported once", func(t *testing.T) {
		t.Parallel()
		out, code := bleed(t, map[string]string{
			"architecture/01-context.md": doc("The app"),
			"architecture/02-app.md":     doc("The  APP", "the app", "New detail"),
		})
		want := repeat("architecture/02-app.md", "the app", 1)
		acceptance.Sees(t, out, code, 1, want)
		if n := strings.Count(acceptance.Read(out), acceptance.Read(want)); n != 1 {
			t.Errorf("the repeat must be reported once, got %d:\n%s", n, out)
		}
	})
	t.Run("0001.1.3 distinct or blank descriptions pass", func(t *testing.T) {
		t.Parallel()
		out, code := bleed(t, map[string]string{
			"architecture/01-context.md":    doc("The whole system", ""),
			"architecture/02-containers.md": doc("The macOS app", ""),
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 levels one and three sharing text pass", func(t *testing.T) {
		t.Parallel()
		out, code := bleed(t, map[string]string{
			"architecture/01-context.md":    doc("Shared line"),
			"architecture/02-containers.md": doc("Something else"),
			"architecture/03-components.md": doc("Shared line"),
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.5 levels two and ten compare as adjacent", func(t *testing.T) {
		t.Parallel()
		out, code := bleed(t, map[string]string{
			"architecture/2-app.md":   doc("The engine"),
			"architecture/10-deep.md": doc("The engine"),
		})
		acceptance.Sees(t, out, code, 1, repeat("architecture/10-deep.md", "the engine", 2))
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 Markdown outside numbered level files passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := bleed(t, map[string]string{
			"architecture/overview.md":    doc("Anything"),
			"not-architecture/01-x.md":    doc("Anything"),
			"not-architecture/02-x.md":    doc("Anything"),
			"docs/01-context-notes.md":    doc("Anything"),
			"docs/02-container-notes.md":  doc("Anything"),
			"architecture/sub/01-deep.md": doc("Anything"),
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "mermaid-level-bleed")
	})
	t.Run("0001.2.2 the Description header column is compared", func(t *testing.T) {
		t.Parallel()
		out, code := bleed(t, map[string]string{
			"architecture/01-context.md": "| # | Label | DESCRIPTION |\n| --- | --- | --- |\n| 1 | Same label | The system |\n",
			"architecture/02-app.md":     "| # | Label | Description |\n| --- | --- | --- |\n| 1 | Same label | The system |\n",
		})
		acceptance.Sees(t, out, code, 1, repeat("architecture/02-app.md", "the system", 1))
		if strings.Contains(out, "same label") {
			t.Errorf("only the Description column may be compared:\n%s", out)
		}
	})
	t.Run("0001.2.3 a missing Description header uses the second column", func(t *testing.T) {
		t.Parallel()
		out, code := bleed(t, map[string]string{
			"architecture/01-context.md": "| # | Text | Notes |\n| --- | --- | --- |\n| 1 | Same | x |\n",
			"architecture/02-app.md":     "| # | Text | Notes |\n| --- | --- | --- |\n| 1 | same | x |\n",
		})
		acceptance.Sees(t, out, code, 1, repeat("architecture/02-app.md", "same", 1))
		if strings.Contains(out, `"x"`) {
			t.Errorf("only the second column may be compared:\n%s", out)
		}
	})
}
