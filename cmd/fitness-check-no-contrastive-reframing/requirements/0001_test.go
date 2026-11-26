// Package requirements proves the no-contrastive-reframing check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// hint is the start of every failure message.
const hint = `README.md: rewrite as a direct statement, not the "not X, it's Y" pattern: `

// reframing runs the check on happyRepo with prose appended to its README.
func reframing(t *testing.T, prose string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{"README.md": "---\ntitle: App\n---\n\n# App\n\nShort and clean.\n\n" + prose})
	return acceptance.Fitness(t, repo, nil, "no-contrastive-reframing")
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a negation then a restating sentence fails", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "It's not a workout. It's a lifestyle.\n\nThey’re not bugs.\n\nThey’re features.\n")
		acceptance.Sees(t, out, code, 1, hint+`"It's not a workout It's a lifestyle"`,
			hint+`"They’re not bugs They’re features"`)
	})
	t.Run("0001.1.2 one sentence that negates plus pivots fails", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "This isn't about speed; it's about consistency.\n")
		acceptance.Sees(t, out, code, 1, hint+`"This isn't about speed; it's about consistency"`)
	})
	t.Run("0001.1.3 ordinary negations and the not only X but also Y form pass", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "It's not ready yet. The file is not found, but the fallback loads.\n\n"+
			"We should not merge, but wait for review. Not only fast, but also correct.\n\n"+
			"It's not X. It's not Y.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 the pattern in code or front matter or a quote passes", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			"README.md": "---\ntitle: It's not X. It's Y.\n---\n\n# App\n\n```\nIt's not a workout. It's a lifestyle.\n```\n\n" +
				"## It's not a workout. It's a lifestyle.\n\n" +
				"Use `It's not a workout. It's a lifestyle.` verbatim.\n",
		})
		out, code := acceptance.Fitness(t, repo, nil, "no-contrastive-reframing")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = reframing(t, `The phrase "It's not a workout. It's a lifestyle." is banned.`+"\n\n"+
			"The phrase “It's not a workout. It's a lifestyle.” is banned.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.5 a long matching sentence is shortened", func(t *testing.T) {
		t.Parallel()
		long := "It's not " + strings.Repeat("a very long clause ", 10) + ", but simple."
		out, code := reframing(t, long+"\n")
		acceptance.Sees(t, out, code, 1, hint+`"`+long[:80]+`…"`)
	})
	t.Run("0001.1.6 a code comment holding the pattern fails naming its line", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{
			"app/main.go":   "package main\n\n// It's not a cache. It's the source of truth.\nfunc main() {}\n",
			"app/App.swift": "import Foundation\n\n/// It's not a cache. It's the source of truth.\nfunc run() {}\n",
			"app/index.ts":  "export {};\n\n// It's not a cache. It's the source of truth.\nexport const x = 1;\n",
		}
		repo := acceptance.Example(t, "happyRepo", files)
		out, code := acceptance.Fitness(t, repo, nil, "no-contrastive-reframing")
		want := `rewrite as a direct statement, not the "not X, it's Y" pattern`
		acceptance.Sees(t, out, code, 1, "app/App.swift:3: "+want, "app/index.ts:3: "+want, "app/main.go:3: "+want)
	})
}
