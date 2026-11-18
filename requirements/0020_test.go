package requirements

import (
	"maps"
	"strings"
	"testing"
)

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

// storybookShown is the fixture app's page file that shows both
// acceptances, so a test's own pages are all that it judges.
const storybookShown = "App/Shown.swift"

// storybookCheck runs swift-storybook-requirements on happyRepo with the
// fixture requirement doc, storybookShown, and files written over them. An
// empty body deletes a fixture file.
func storybookCheck(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	app := map[string]string{
		"docs/requirements/0001-badge.md": storybookDoc,
		storybookShown:                    "let shown = [\n    BookPage(title: \"0001.1 shown\") { BadgeView() },\n    BookPage(title: \"0001.2 shown\") { BadgeView() },\n]\n",
	}
	maps.Copy(app, files)
	return fitness(t, example(t, "happyRepo", app), nil, "swift-storybook-requirements")
}

// storybookView wraps previews in a SwiftUI view file whose first preview
// starts on line 7.
func storybookView(previews string) string {
	return "import SwiftUI\n\nstruct BadgeView: View {\n    var body: some View { Text(verbatim: \"\") }\n}\n\n" + previews
}

func Test0020_1(t *testing.T) {
	t.Parallel()
	view := storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n\n#Preview(\"0001.1 night\", traits: .sizeThatFitsLayout) { BadgeView() }\n\n#Preview(\"0001.2\") {\n    BadgeView()\n}\n")
	pages := "let pages = [\n    BookPage(title: \"0001.2 empty\") { BadgeView() },\n]\n"
	out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": view, "App/Pages.swift": pages})
	sees(t, out, code, 0, "All 1 checks passed")
	passedWith(t, out, "swift-storybook-requirements", "3")
}

func Test0020_2(t *testing.T) {
	t.Parallel()
	pages := "let pages = [BookPage(\"Pages.swift\", 1, content: Card(title: \"0001.1 day\"))]\n"
	out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview {\n    BadgeView()\n}\n"), "App/Pages.swift": pages})
	sees(t, out, code, 1, "App/BadgeView.swift:7 - page has no name; name it by the acceptance it shows, such as #Preview(\"0001.1 label\")",
		"App/Pages.swift:1 - page has no name")
}

func Test0020_3(t *testing.T) {
	t.Parallel()
	pages := "let pages = [BookPage(title: \"BadgeView: empty\") { BadgeView() }]\n"
	out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"day\") { BadgeView() }\n"), "App/Pages.swift": pages})
	sees(t, out, code, 1, "App/BadgeView.swift:7 - page \"day\" does not start with an acceptance ID; name it such as \"0001.1 day\"",
		"App/Pages.swift:1 - page \"BadgeView: empty\" does not start with an acceptance ID")
}

func Test0020_4(t *testing.T) {
	t.Parallel()
	out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.9 day\") { BadgeView() }\n")})
	sees(t, out, code, 1, "App/BadgeView.swift:7 - page \"0001.9 day\" names 0001.9, which no doc in docs/requirements defines")
}

func Test0020_5(t *testing.T) {
	t.Parallel()
	pages := "\nlet pages = [BookPage(title: \"0001.1 day\") { BadgeView() }]\n"
	out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n"), "App/Pages.swift": pages})
	sees(t, out, code, 1, "App/Pages.swift:2 - page \"0001.1 day\" repeats the name of the page at App/BadgeView.swift:7")
}

func Test0020_6(t *testing.T) {
	t.Parallel()
	pages := "let name = \"0001.1 day\"\nlet pages = [BookPage(title: name, file: #file) { BadgeView() }]\n"
	view := storybookView("#Preview(\"0001.1 \\(mode)\") { BadgeView() }\n")
	out, code := storybookCheck(t, map[string]string{"App/Pages.swift": pages, "App/BadgeView.swift": view})
	sees(t, out, code, 1, "App/Pages.swift:2 - page name name is not a plain literal; write the acceptance name, such as \"0001.1 label\"",
		"App/BadgeView.swift:7 - page name \"0001.1 \\(mode)\" is not a plain literal")
}

func Test0020_7(t *testing.T) {
	t.Parallel()
	view := storybookView("// #Preview { BadgeView() }\n/* #Preview { BadgeView() } */\nlet note = \"#Preview { BadgeView() }\"\n")
	out, code := storybookCheck(t, map[string]string{"App/BadgeView.swift": view})
	sees(t, out, code, 0, "All 1 checks passed")
	passedWith(t, out, "swift-storybook-requirements", "2")
	out, code = storybookCheck(t, map[string]string{storybookShown: ""})
	sees(t, out, code, 0, "All 1 checks passed")
	passedWithNoFiles(t, out, "swift-storybook-requirements")
	out, code = storybookCheck(t, map[string]string{"App/BadgeView.swift": storybookView("#Preview { BadgeView() }\n"), "docs/requirements/0001-badge.md": ""})
	sees(t, out, code, 0, "All 1 checks passed")
	if strings.Contains(out, "page has no name") {
		t.Errorf("want no finding without requirement docs:\n%s", out)
	}
}

func Test0020_8(t *testing.T) {
	t.Parallel()
	out, code := storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": storybookView("")})
	sees(t, out, code, 1, "docs/requirements/0001-badge.md:15 - acceptance 0001.1 has no Storybook page; add one named \"0001.1 <label>\"",
		"docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
	out, code = storybookCheck(t, map[string]string{storybookShown: "", "App/BadgeView.swift": storybookView("#Preview(\"0001.1 day\") { BadgeView() }\n")})
	sees(t, out, code, 1, "docs/requirements/0001-badge.md:19 - acceptance 0001.2 has no Storybook page")
	if strings.Contains(out, "acceptance 0001.1") {
		t.Errorf("want no finding for the acceptance a page names:\n%s", out)
	}
}
