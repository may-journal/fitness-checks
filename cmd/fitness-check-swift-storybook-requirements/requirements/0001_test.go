// Package requirements proves the swift-storybook-requirements check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"maps"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1.1.1 it passes with a note that no Storybook is in the repo", func(t *testing.T) {
		t.Parallel()
		view := storybookView("import StorybookHistoryKit\n// import Storybook\nlet note = \"\"\"\nimport StorybookKit\n\"\"\"\n\n#Preview {\n    BadgeView()\n}\n")
		out, code := storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": view})
		acceptance.ShowsInfo(t, out, code, "swift-storybook-requirements", "no Storybook in this repo: no Swift file imports a Storybook module or declares a BookPage, so acceptances need no pages")
		if strings.Contains(out, "has no Storybook page") || strings.Contains(out, "page has no name") {
			t.Errorf("want no finding in a repo without Storybook:\n%s", out)
		}
	})
	t.Run("0001.1.1.2.1 the repo is judged even with no pages yet", func(t *testing.T) {
		t.Parallel()
		for _, imports := range []string{
			"@testable import StorybookKit\n",
			"import SwiftUIStorybook\n",
			"#if DEBUG\n    @_spi(Storybook, Internal) import StorybookSandbox\n#endif\n",
			"import SwiftUI; import Storybook\n",
			"let pages = [BookPage(title: \"0001.1 day\") { Text(verbatim: \"\") }]\n",
		} {
			out, code := storybookCheck(t, map[string]string{storybookShown: "", storybookCatalog: imports})
			acceptance.Sees(t, out, code, 1, "docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
		}
	})
	t.Run("0001.1.1.3.1 it passes", func(t *testing.T) {
		t.Parallel()
		view := storybookView("#Preview(\"0001.1 day\", traits: .sizeThatFitsLayout) { BadgeView() }\n")
		pages := "let pages = [\n    BookPage(title: \"0001.2\") { BadgeView() },\n]\n"
		out, code := storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": view, "App/Pages.swift": pages})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// A page in a comment or a string is not read.
		view = storybookView("// #Preview { BadgeView() }\n/* #Preview { BadgeView() } */\nlet note = \"#Preview { BadgeView() }\"\n")
		out, code = storybookCheck(t, map[string]string{"App/BadgeView.swift": view})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		// With no requirement docs there is nothing to name a page by.
		out, code = storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview { BadgeView() }\n"), "docs/requirements/0001-badge.md": ""})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if strings.Contains(out, "page has no name") {
			t.Errorf("want no finding without requirement docs:\n%s", out)
		}
	})
	t.Run("0001.1.1.3.2 it reads every Swift file", func(t *testing.T) {
		t.Parallel()
		view := storybookView("#Preview(\"0001.1 day\", traits: .sizeThatFitsLayout) { BadgeView() }\n")
		pages := "let pages = [\n    BookPage(title: \"0001.2\") { BadgeView() },\n]\n"
		out, code := storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": view, "App/Pages.swift": pages})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "swift-storybook-requirements", "2")
		view = storybookView("// #Preview { BadgeView() }\n/* #Preview { BadgeView() } */\nlet note = \"#Preview { BadgeView() }\"\n")
		out, code = storybookCheck(t, map[string]string{"App/BadgeView.swift": view})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "swift-storybook-requirements", "2")
		out, code = storybookCheck(t, map[string]string{storybookShown: ""})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "swift-storybook-requirements")
	})

	// unnamed runs a page with no name in a preview and in a BookPage.
	unnamed := func(t *testing.T) (string, int) {
		pages := "let pages = [BookPage(\"Pages.swift\", 1, content: Card(title: \"0001.1 day\"))]\n"
		return storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview {\n    BadgeView()\n}\n"), "App/Pages.swift": pages})
	}
	// notLiteral runs a page name that is a variable and one with interpolation.
	notLiteral := func(t *testing.T) (string, int) {
		pages := "let name = \"0001.1 day\"\nlet pages = [BookPage(title: name, file: #file) { BadgeView() }]\n"
		view := storybookView("#Preview(\"0001.1 \\(mode)\") { BadgeView() }\n")
		return storybookCheck(t, map[string]string{"App/Pages.swift": pages, "App/BadgeView.swift": view})
	}
	t.Run("0001.1.1.4.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		out, code := unnamed(t)
		acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:7 - page has no name", "App/Pages.swift:1 - page has no name")
		out, code = notLiteral(t)
		acceptance.Sees(t, out, code, 1, "App/Pages.swift:2 - page name", "App/BadgeView.swift:7 - page name")
	})
	t.Run("0001.1.1.4.2 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		out, code := unnamed(t)
		acceptance.Sees(t, out, code, 1, ".swift:7 - page has no name", ".swift:1 - page has no name")
		out, code = notLiteral(t)
		acceptance.Sees(t, out, code, 1, ".swift:2 - page name name is not a plain literal", ".swift:7 - page name \"0001.1 \\(mode)\" is not a plain literal")
	})
	t.Run("0001.1.1.4.3 it says how to name the page", func(t *testing.T) {
		t.Parallel()
		out, code := unnamed(t)
		acceptance.Sees(t, out, code, 1, "name it by the acceptance it shows, such as #Preview(\"0001.1 label\")")
		out, code = notLiteral(t)
		acceptance.Sees(t, out, code, 1, "write the acceptance name, such as \"0001.1 label\"")
	})

	// chained is the doc with chained steps under acceptance 0001.2.
	chained := strings.Replace(storybookDoc, "        - When the photo shows\n            - Then it has no badge\n",
		"        1. When the photo shows\n            - Then it has no badge\n        2. When the tester taps it\n            - Then the badge stays hidden\n", 1)
	// undefined runs names that no acceptance ID starts, and a step name.
	undefined := func(t *testing.T) []string {
		var outs []string
		for _, files := range []map[string]string{
			{"App/BadgeView.swift": storybookView("#Preview(\"day\") { BadgeView() }\n"), "App/Pages.swift": "let pages = [BookPage(title: \"BadgeView: empty\") { BadgeView() }]\n"},
			{"App/BadgeView.swift": storybookView("#Preview(\"0001.9 day\") { BadgeView() }\n")},
			{"docs/requirements/0001-badge.md": chained, "App/BadgeView.swift": storybookView("#Preview(\"0001.2.2 tapped\") { BadgeView() }\n")},
			{"docs/requirements/0001-badge.md": chained, "App/BadgeView.swift": storybookView("#Preview(\"0001.9.2 tapped\") { BadgeView() }\n")},
		} {
			out, code := storybookCheck(t, files)
			if code != 1 {
				t.Errorf("want exit code 1, got %d:\n%s", code, out)
			}
			outs = append(outs, out)
		}
		return outs
	}
	t.Run("0001.1.1.5.1 it fails naming the file", func(t *testing.T) {
		t.Parallel()
		outs := undefined(t)
		acceptance.Sees(t, outs[0], 1, 1, "App/BadgeView.swift:7 - page \"day\" does not start with an acceptance ID", "App/Pages.swift:1 - page \"BadgeView: empty\" does not start with an acceptance ID")
		acceptance.Sees(t, outs[1], 1, 1, "App/BadgeView.swift:7 - page \"0001.9 day\" names 0001.9, which no doc in docs/requirements defines")
		acceptance.Sees(t, outs[2], 1, 1, `App/BadgeView.swift:7 - page "0001.2.2 tapped" names step 0001.2.2`)
		acceptance.Sees(t, outs[3], 1, 1, `App/BadgeView.swift:7 - page "0001.9.2 tapped" does not start with an acceptance ID`)
	})
	t.Run("0001.1.1.5.2 it fails naming the line", func(t *testing.T) {
		t.Parallel()
		outs := undefined(t)
		acceptance.Sees(t, outs[0], 1, 1, ".swift:7 - page \"day\" does not start with an acceptance ID", ".swift:1 - page \"BadgeView: empty\" does not start with an acceptance ID")
		acceptance.Sees(t, outs[2], 1, 1, `.swift:7 - page "0001.2.2 tapped" names step 0001.2.2`)
	})
	t.Run("0001.1.1.5.3 it names the page name to use", func(t *testing.T) {
		t.Parallel()
		outs := undefined(t)
		acceptance.Sees(t, outs[0], 1, 1, "name it such as \"0001.1 day\"")
		acceptance.Sees(t, outs[2], 1, 1, `one page shows every step of an acceptance, so name it "0001.2 tapped"`)
		// An acceptance with chained steps keeps one page, so its own name passes.
		out, code := storybookCheck(t, map[string]string{"docs/requirements/0001-badge.md": chained})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})

	// view holds two pages of acceptance 0001.1, and storybookShown a third.
	view := storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n\n#Preview(\"0001.1 night\") { BadgeView() }\n")
	t.Run("0001.1.1.6.1 it fails naming both places", func(t *testing.T) {
		t.Parallel()
		out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": view})
		acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:9 - page \"0001.1 night\" repeats acceptance 0001.1 from line 7; draw each variant as a labeled row on one page",
			"App/Shown.swift:4 - page \"0001.1 shown\" repeats acceptance 0001.1 from App/BadgeView.swift:7")
		// A preview and a BookPage share the count.
		out, code = storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.2 empty\") { BadgeView() }\n")})
		acceptance.Sees(t, out, code, 1, "App/Shown.swift:5 - page \"0001.2 shown\" repeats acceptance 0001.2 from App/BadgeView.swift:7; draw each variant as a labeled row on one page")
		pages := "\nlet pages = [BookPage(title: \"0001.1 day\") { BadgeView() }]\n"
		out, code = storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n"), "App/Pages.swift": pages})
		acceptance.Sees(t, out, code, 1, "App/Pages.swift:2 - page \"0001.1 day\" repeats the name of the page at App/BadgeView.swift:7")
	})
	t.Run("0001.1.1.6.2 a `maxPagesPerAcceptance` setting allows more pages per acceptance", func(t *testing.T) {
		t.Parallel()
		three := view + "\n#Preview(\"0001.1 dusk\") { BadgeView() }\n"
		raised := `{"swiftStorybookRequirements": {"maxPagesPerAcceptance": 2}}` + "\n"
		out, code := storybookCheck(t, map[string]string{".fitnessrc.json": raised, storybookShown: storybookPages, "App/BadgeView.swift": view})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		out, code = storybookCheck(t, map[string]string{".fitnessrc.json": raised, storybookShown: storybookPages, "App/BadgeView.swift": three})
		acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:11 - page \"0001.1 dusk\" repeats acceptance 0001.1 from line 7")
		zero := `{"swiftStorybookRequirements": {"maxPagesPerAcceptance": 0}}` + "\n"
		out, code = storybookCheck(t, map[string]string{".fitnessrc.json": zero, storybookShown: storybookPages, "App/BadgeView.swift": view})
		acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:9 - page \"0001.1 night\" repeats acceptance 0001.1 from line 7")
	})

	catalog := "import Storybook\n\nenum AppStorybook {}\n"
	noPage := func(t *testing.T) []string {
		var outs []string
		for _, previews := range []string{"", "#Preview(\"0001.1 day\") { BadgeView() }\n"} {
			out, code := storybookCheck(t, map[string]string{storybookShown: "", storybookCatalog: catalog, "App/BadgeView.swift": storybookView(previews)})
			if code != 1 {
				t.Errorf("want exit code 1, got %d:\n%s", code, out)
			}
			outs = append(outs, out)
		}
		return outs
	}
	t.Run("0001.1.1.7.1 it fails naming the doc line", func(t *testing.T) {
		t.Parallel()
		outs := noPage(t)
		acceptance.Sees(t, outs[0], 1, 1, "docs/requirements/0001-badge.md:15 - acceptance 0001.1 has no Storybook page",
			"docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
		acceptance.Sees(t, outs[1], 1, 1, "docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
		if strings.Contains(outs[1], "acceptance 0001.1") {
			t.Errorf("want no finding for the acceptance a page names:\n%s", outs[1])
		}
	})
	t.Run("0001.1.1.7.2 it names the page name to add", func(t *testing.T) {
		t.Parallel()
		outs := noPage(t)
		acceptance.Sees(t, outs[0], 1, 1, "has no Storybook page; add one named \"0001.1 <label>\"")
	})
}

// storybookDoc is the fixture app's requirement doc, defining 0001.1 and
// 0001.2.
const storybookDoc = `# 0001 Badge

## Why

Testers see which photos are new.

## Measurement

badges seen
-
photos opened

## Requirements

- 0001.1
    - Given Tester opens a new photo
        - When the photo shows
            - Then it has a badge
- 0001.2
    - Given Tester opens an old photo
        - When the photo shows
            - Then it has no badge
`

// storybookShown is the fixture app's page file that imports Storybook and
// shows both acceptances. A test's own pages are all that it judges.
const storybookShown = "App/Shown.swift"

// storybookPages is a storybookShown body that shows only 0001.2, so a test
// draws the pages of 0001.1 itself.
const storybookPages = "import Storybook\n\nlet shown = [BookPage(title: \"0001.2 shown\") { BadgeView() }]\n"

// storybookCatalog is a file that imports Storybook and declares no page, as
// an app with no pages yet has.
const storybookCatalog = "App/AppStorybook.swift"

// storybookCheck runs swift-storybook-requirements on happyRepo with the
// fixture requirement doc, storybookShown, and files written over them. An
// empty body deletes a fixture file.
func storybookCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	app := map[string]string{
		"docs/requirements/0001-badge.md": storybookDoc,
		storybookShown:                    "import Storybook\n\nlet shown = [\n    BookPage(title: \"0001.1 shown\") { BadgeView() },\n    BookPage(title: \"0001.2 shown\") { BadgeView() },\n]\n",
	}
	maps.Copy(app, files)
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", app), nil, "swift-storybook-requirements")
}

// storybookView wraps previews in a SwiftUI view file whose first preview
// starts on line 7.
func storybookView(previews string) string {
	return "import SwiftUI\n\nstruct BadgeView: View {\n    var body: some View { Text(verbatim: \"\") }\n}\n\n" + previews
}
