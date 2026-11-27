// Package requirements proves the swift-storybook-requirements check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"maps"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a repo without Storybook passes with a note", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.1.1", "35310532", func(_ req.A) {
			view := storybookView("import StorybookHistoryKit\n// import Storybook\nlet note = \"\"\"\nimport StorybookKit\n\"\"\"\n\n#Preview {\n    BadgeView()\n}\n")
			out, code := storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": view})
			acceptance.ShowsInfo(t, out, code, "swift-storybook-requirements", "no Storybook in this repo: no Swift file imports a Storybook module or declares a BookPage, so acceptances need no pages")
			if strings.Contains(out, "has no Storybook page") || strings.Contains(out, "page has no name") {
				t.Errorf("want no finding in a repo without Storybook:\n%s", out)
			}
		})
	})
	t.Run("0001.1.2 a Storybook import or page makes the repo judged", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "ea6b0e46", func(_ req.A) {
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
	})
	t.Run("0001.1.3 one page per acceptance passes reading every file", func(t *testing.T) {
		t.Parallel()
		view := storybookView("#Preview(\"0001.1 day\", traits: .sizeThatFitsLayout) { BadgeView() }\n")
		pages := "let pages = [\n    BookPage(title: \"0001.2\") { BadgeView() },\n]\n"
		out, code := storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": view, "App/Pages.swift": pages})
		// A page in a comment or a string is not read.
		view = storybookView("// #Preview { BadgeView() }\n/* #Preview { BadgeView() } */\nlet note = \"#Preview { BadgeView() }\"\n")
		out2, code2 := storybookCheck(t, map[string]string{"App/BadgeView.swift": view})
		out3, code3 := storybookCheck(t, map[string]string{storybookShown: ""})
		// With no requirement docs there is nothing to name a page by.
		out4, code4 := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview { BadgeView() }\n"), "docs/requirements/0001-badge.md": ""})
		req.Then(t, "0001.1.1.3.1", "07b8e753", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.Sees(t, out2, code2, 0, "All 1 checks passed")
			acceptance.Sees(t, out3, code3, 0, "All 1 checks passed")
			acceptance.Sees(t, out4, code4, 0, "All 1 checks passed")
			a.True(!strings.Contains(out4, "page has no name"), "want no finding without requirement docs:\n%s", out4)
		})
		req.Then(t, "0001.1.1.3.2", "29f2bb4e", func(_ req.A) {
			acceptance.PassedWith(t, out, "swift-storybook-requirements", "2")
			acceptance.PassedWith(t, out2, "swift-storybook-requirements", "2")
			acceptance.PassedWithNoFiles(t, out3, "swift-storybook-requirements")
		})
	})
	t.Run("0001.1.4 a page without a literal name fails naming the file then the line", func(t *testing.T) {
		t.Parallel()
		pages := "let pages = [BookPage(\"Pages.swift\", 1, content: Card(title: \"0001.1 day\"))]\n"
		out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview {\n    BadgeView()\n}\n"), "App/Pages.swift": pages})
		pages = "let name = \"0001.1 day\"\nlet pages = [BookPage(title: name, file: #file) { BadgeView() }]\n"
		view := storybookView("#Preview(\"0001.1 \\(mode)\") { BadgeView() }\n")
		out2, code2 := storybookCheck(t, map[string]string{"App/Pages.swift": pages, "App/BadgeView.swift": view})
		req.Then(t, "0001.1.1.4.1", "f7956a3e", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:", "App/Pages.swift:")
			acceptance.Sees(t, out2, code2, 1, "App/Pages.swift:", "App/BadgeView.swift:")
		})
		req.Then(t, "0001.1.1.4.2", "aabd0f94", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:7 - page has no name", "App/Pages.swift:1 - page has no name")
			acceptance.Sees(t, out2, code2, 1, "App/Pages.swift:2 - page name name is not a plain literal",
				"App/BadgeView.swift:7 - page name \"0001.1 \\(mode)\" is not a plain literal")
		})
		req.Then(t, "0001.1.1.4.3", "d7270adf", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "name it by the acceptance it shows, such as #Preview(\"0001.1 label\")")
			acceptance.Sees(t, out2, code2, 1, "write the acceptance name, such as \"0001.1 label\"")
		})
	})
	t.Run("0001.1.5 a name without a defined acceptance ID fails naming the file and line that hold it", func(t *testing.T) {
		t.Parallel()
		pages := "let pages = [BookPage(title: \"BadgeView: empty\") { BadgeView() }]\n"
		out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"day\") { BadgeView() }\n"), "App/Pages.swift": pages})
		out2, code2 := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.9 day\") { BadgeView() }\n")})
		// An acceptance with chained steps keeps one page, so a step name fails.
		chained := strings.Replace(storybookDoc, "        - When the photo shows\n            - Then it has no badge\n",
			"        1. When the photo shows\n            - Then it has no badge\n        2. When the tester taps it\n            - Then the badge stays hidden\n", 1)
		out3, code3 := storybookCheck(t, map[string]string{"docs/requirements/0001-badge.md": chained})
		out4, code4 := storybookCheck(t, map[string]string{"docs/requirements/0001-badge.md": chained, "App/BadgeView.swift": storybookView("#Preview(\"0001.2.2 tapped\") { BadgeView() }\n")})
		out5, code5 := storybookCheck(t, map[string]string{"docs/requirements/0001-badge.md": chained, "App/BadgeView.swift": storybookView("#Preview(\"0001.9.2 tapped\") { BadgeView() }\n")})
		req.Then(t, "0001.1.1.5.1", "a3526dc6", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:", "App/Pages.swift:")
			acceptance.Sees(t, out2, code2, 1, "App/BadgeView.swift:")
			acceptance.Sees(t, out4, code4, 1, "App/BadgeView.swift:")
			acceptance.Sees(t, out5, code5, 1, "App/BadgeView.swift:")
		})
		req.Then(t, "0001.1.1.5.2", "3ff7ef1d", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:7 - page \"day\" does not start with an acceptance ID",
				"App/Pages.swift:1 - page \"BadgeView: empty\" does not start with an acceptance ID")
			acceptance.Sees(t, out2, code2, 1, "App/BadgeView.swift:7 - page \"0001.9 day\" names 0001.9, which no doc in docs/requirements defines")
			acceptance.Sees(t, out4, code4, 1, `App/BadgeView.swift:7 - page "0001.2.2 tapped" names step 0001.2.2`)
			acceptance.Sees(t, out5, code5, 1, `App/BadgeView.swift:7 - page "0001.9.2 tapped" does not start with an acceptance ID`)
		})
		req.Then(t, "0001.1.1.5.3", "21f664d4", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "name it such as \"0001.1 day\"")
			acceptance.Sees(t, out3, code3, 0, "All 1 checks passed")
			acceptance.Sees(t, out4, code4, 1, `one page shows every step of an acceptance, so name it "0001.2 tapped"`)
		})
	})
	t.Run("0001.1.6 a repeated name or acceptance ID fails", func(t *testing.T) {
		t.Parallel()
		view := storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n\n#Preview(\"0001.1 night\") { BadgeView() }\n")
		out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": view})
		// A preview and a BookPage share the count.
		out2, code2 := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.2 empty\") { BadgeView() }\n")})
		pages := "\nlet pages = [BookPage(title: \"0001.1 day\") { BadgeView() }]\n"
		out3, code3 := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n"), "App/Pages.swift": pages})
		// The repo may allow more pages per acceptance.
		three := view + "\n#Preview(\"0001.1 dusk\") { BadgeView() }\n"
		raised := `{"swiftStorybookRequirements": {"maxPagesPerAcceptance": 2}}` + "\n"
		out4, code4 := storybookCheck(t, map[string]string{".fitnessrc.json": raised, storybookShown: storybookPages, "App/BadgeView.swift": view})
		out5, code5 := storybookCheck(t, map[string]string{".fitnessrc.json": raised, storybookShown: storybookPages, "App/BadgeView.swift": three})
		zero := `{"swiftStorybookRequirements": {"maxPagesPerAcceptance": 0}}` + "\n"
		out6, code6 := storybookCheck(t, map[string]string{".fitnessrc.json": zero, storybookShown: storybookPages, "App/BadgeView.swift": view})
		req.Then(t, "0001.1.1.6.1", "25af8c41", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "App/BadgeView.swift:9 - page \"0001.1 night\" repeats acceptance 0001.1 from line 7; draw each variant as a labeled row on one page",
				"App/Shown.swift:4 - page \"0001.1 shown\" repeats acceptance 0001.1 from App/BadgeView.swift:7")
			acceptance.Sees(t, out2, code2, 1, "App/Shown.swift:5 - page \"0001.2 shown\" repeats acceptance 0001.2 from App/BadgeView.swift:7; draw each variant as a labeled row on one page")
			acceptance.Sees(t, out3, code3, 1, "App/Pages.swift:2 - page \"0001.1 day\" repeats the name of the page at App/BadgeView.swift:7")
		})
		req.Then(t, "0001.1.1.6.2", "cc72b5c3", func(_ req.A) {
			acceptance.Sees(t, out4, code4, 0, "All 1 checks passed")
			acceptance.Sees(t, out5, code5, 1, "App/BadgeView.swift:11 - page \"0001.1 dusk\" repeats acceptance 0001.1 from line 7")
			acceptance.Sees(t, out6, code6, 1, "App/BadgeView.swift:9 - page \"0001.1 night\" repeats acceptance 0001.1 from line 7")
		})
	})
	t.Run("0001.1.7 an acceptance with no page fails naming its doc line", func(t *testing.T) {
		t.Parallel()
		catalog := "import Storybook\n\nenum AppStorybook {}\n"
		out, code := storybookCheck(t, map[string]string{storybookShown: "", storybookCatalog: catalog, "App/BadgeView.swift": storybookView("")})
		out2, code2 := storybookCheck(t, map[string]string{storybookShown: "", storybookCatalog: catalog, "App/BadgeView.swift": storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n")})
		req.Then(t, "0001.1.1.7.1", "09837be2", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "docs/requirements/0001-badge.md:15 - acceptance 0001.1 has no Storybook page",
				"docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
			acceptance.Sees(t, out2, code2, 1, "docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
			a.True(!strings.Contains(out2, "acceptance 0001.1"), "want no finding for the acceptance a page names:\n%s", out2)
		})
		req.Then(t, "0001.1.1.7.2", "06d2e1b2", func(_ req.A) {
			acceptance.Sees(t, out, code, 1, "add one named \"0001.1 <label>\"")
		})
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
