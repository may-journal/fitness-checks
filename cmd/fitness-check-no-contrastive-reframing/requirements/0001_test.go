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
	t.Run("0001.1.2.1.1 it fails once per match", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "It's not a workout. It's a lifestyle.\n\nThey’re not bugs.\n\nThey’re features.\n")
		acceptance.Sees(t, out, code, 1, "README.md: rewrite as a direct statement")
		if n := strings.Count(out, "rewrite as a direct statement"); n != 2 {
			t.Errorf("want one failure per match (2), got %d:\n%s", n, out)
		}
	})
	t.Run("0001.1.2.1.2 it quotes the snippet", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "It's not a workout. It's a lifestyle.\n\nThey’re not bugs.\n\nThey’re features.\n")
		acceptance.Sees(t, out, code, 1, hint+`"It's not a workout It's a lifestyle"`,
			hint+`"They’re not bugs They’re features"`)
	})
	t.Run("0001.1.2.2.1 it fails quoting the snippet", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "This isn't about speed; it's about consistency.\n")
		acceptance.Sees(t, out, code, 1, hint+`"This isn't about speed; it's about consistency"`)
	})
	t.Run("0001.1.2.3.1 it passes", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "It's not ready yet. The file is not found, but the fallback loads.\n\n"+
			"We should not merge, but wait for review.\n\n"+
			"It's not X. It's not Y.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.3.2 the form `not only X but also Y` passes", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "Not only fast, but also correct.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.4.1 it passes", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{
			"README.md": "---\ntitle: It's not X. It's Y.\n---\n\n# App\n\n```\nIt's not a workout. It's a lifestyle.\n```\n\n" +
				"## It's not a workout. It's a lifestyle.\n",
		})
		out, code := acceptance.Fitness(t, repo, nil, "no-contrastive-reframing")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.4.2 the same holds in inline code or a quoted span", func(t *testing.T) {
		t.Parallel()
		out, code := reframing(t, "Use `It's not a workout. It's a lifestyle.` verbatim.\n\n"+
			`The phrase "It's not a workout. It's a lifestyle." is banned.`+"\n\n"+
			"The phrase “It's not a workout. It's a lifestyle.” is banned.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.2.5.1 the quoted snippet is shortened with an ellipsis", func(t *testing.T) {
		t.Parallel()
		long := "It's not " + strings.Repeat("a very long clause ", 10) + ", but simple."
		out, code := reframing(t, long+"\n")
		acceptance.Sees(t, out, code, 1, hint+`"`+long[:80]+`…"`)
	})
	comments := func(t *testing.T) (string, int) {
		files := map[string]string{
			"app/main.go":   "package main\n\n// It's not a cache. It's the source of truth.\nfunc main() {}\n",
			"app/App.swift": "import Foundation\n\n/// It's not a cache. It's the source of truth.\nfunc run() {}\n",
			"app/index.ts":  "export {};\n\n// It's not a cache. It's the source of truth.\nexport const x = 1;\n",
		}
		return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "no-contrastive-reframing")
	}
	const want = `rewrite as a direct statement, not the "not X, it's Y" pattern`
	t.Run("0001.1.2.6.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := comments(t)
		acceptance.Sees(t, out, code, 1, "app/App.swift:", "app/index.ts:", "app/main.go:", want)
	})
	t.Run("0001.1.2.6.2 it names the line of the block", func(t *testing.T) {
		t.Parallel()
		out, code := comments(t)
		acceptance.Sees(t, out, code, 1, "app/App.swift:3: "+want, "app/index.ts:3: "+want, "app/main.go:3: "+want)
	})
}
