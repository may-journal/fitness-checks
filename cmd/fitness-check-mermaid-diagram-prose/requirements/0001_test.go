// Package requirements proves the mermaid-diagram-prose check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

// callouts is a callout table that pairs with the diagram above it.
const callouts = "\n| # | Description | Why |\n| --- | --- | --- |\n| 6 | uses | because |\n"

// diagram wraps body in a mermaid fence under a heading.
func diagram(body string) string {
	return "# Architecture\n\n```mermaid\n" + body + "\n```\n"
}

// paired is a doc whose diagram is paired with a callout table.
func paired(body string) string {
	return diagram(body) + callouts
}

// prose runs the check on happyRepo with files written over it.
func prose(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, append([]string{"mermaid-diagram-prose"}, args...)...)
}

// finding is the message a user sees for one prose label in docs/arch.md.
func finding(label string) string {
	return `docs/arch.md: diagram at line 3 has a relationship label with prose ("` + label + `") — put descriptions in the callout table`
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a C4 relationship label with prose fails", func(t *testing.T) {
		t.Parallel()
		body := "Rel(app, api, \"6 Uses over HTTP\")\nBiRel(app, db, \"6 talks both ways\")\nRel_Back(api, app, '6 replies')\nRel(api, db,\n  \"6 wrapped prose\")\nRel(api, db, \"B reads\")"
		out, code := prose(t, map[string]string{"docs/arch.md": paired(body)})
		acceptance.Sees(t, out, code, 1, finding("6 Uses over HTTP"), finding("6 talks both ways"),
			finding("6 replies"), finding("6 wrapped prose"), finding("B reads"))
	})
	t.Run("0001.1.2 a flowchart pipe label or a colon edge label with prose fails", func(t *testing.T) {
		t.Parallel()
		out, code := prose(t, map[string]string{"docs/arch.md": paired("flowchart LR\napp -->|6 calls the api| api -->|6 reads| db -->|A calls back| app")})
		acceptance.Sees(t, out, code, 1, finding("6 calls the api"), finding("6 reads"), finding("A calls back"))
		out, code = prose(t, map[string]string{"docs/arch.md": paired("flowchart LR\napp --> api : 6 sends data   \napi --> db : C writes")})
		acceptance.Sees(t, out, code, 1, finding("6 sends data"), finding("C writes"))
	})
	t.Run("0001.1.3 bare number labels pass", func(t *testing.T) {
		t.Parallel()
		body := "flowchart LR\napp -->|6| api\napi --> db : 6\ndb --> app : uses 6\n%% note : 6 something\nRel(app, api, \"6\")\nclassDef x fill:#6366f1"
		out, code := prose(t, map[string]string{"docs/arch.md": paired(body)})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 a prose label fails without a callout table", func(t *testing.T) {
		t.Parallel()
		out, code := prose(t, map[string]string{"docs/arch.md": diagram(`Rel(app, api, "6 Uses over HTTP")`)})
		acceptance.Sees(t, out, code, 1, finding("6 Uses over HTTP"))
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 prose in node names passes", func(t *testing.T) {
		t.Parallel()
		body := "Container(app, \"6 App Shell\", \"SwiftUI\")\nRel(app, api, \"6\")\na[\"6 node\"] -->|6| b"
		out, code := prose(t, map[string]string{"docs/arch.md": paired(body)})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.2.2 a PR body prose label fails naming the label", func(t *testing.T) {
		t.Parallel()
		files := map[string]string{"body.md": paired(`Rel(app, api, "6 Uses over HTTP")`)}
		out, code := prose(t, files, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, `has a relationship label with prose ("6 Uses over HTTP")`)
	})
	t.Run("0001.2.3 a repo with no diagram passes with no files", func(t *testing.T) {
		t.Parallel()
		out, code := prose(t, nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "mermaid-diagram-prose")
	})
}
