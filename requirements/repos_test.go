package requirements

import (
	"encoding/json"
	"strconv"
	"testing"
)

// The Given files a user's repo adds to happyRepo. They live here, not as
// tracked example repos, so this repo's own checks never judge them.

// hungTool is a JavaScript repo whose own prettier runs prelude, then hangs
// after starting a child process. The child names prettier.pid on its command
// line until it has moved its PID into that file, so processStopped finds it
// either way.
func hungTool(prelude string) map[string]string {
	child := `/bin/sh -c 'echo $$ > "$1.tmp" && mv "$1.tmp" "$1" && exec /bin/sleep 30' child "$PWD/prettier.pid"`
	return map[string]string{
		"package.json":               "{\"name\": \"app\", \"private\": true}\n",
		"app.js":                     "const a = 1;\n",
		"node_modules/.bin/prettier": "#!/bin/sh\n" + prelude + child + " >/dev/null 2>&1 &\nwait\n",
	}
}

// fastTotal is the environment entry that lowers the fast total to ms. A test
// then sees the total stop a check without waiting a minute.
func fastTotal(ms int) string {
	return "FITNESS_INTERNAL_FAST_TOTAL_MS=" + strconv.Itoa(ms)
}

// testsTotal is the environment entry that lowers the test stage's total to
// ms, as fastTotal does the fast total.
func testsTotal(ms int) string {
	return "FITNESS_INTERNAL_TESTS_TOTAL_MS=" + strconv.Itoa(ms)
}

// planBody is a Plan that follows the template, with one task.
func planBody(task string) string {
	return "> One clear pitch for this Plan.\n\n## Background\n\nThis Plan is a test fixture.\n\n## What needs to happen\n\n" + task + "\n"
}

// planView is the saved gh answer to the pre-push gate's read of issue n in
// exampleRepo: URL, body, labels and the owner's comments.
func planView(t *testing.T, n int, body string, labels []string, comments ...string) response {
	t.Helper()
	by := make([]comment, len(comments))
	for i, c := range comments {
		by[i] = comment{Body: c, Author: "example-owner", Association: "OWNER"}
	}
	return planViewBy(t, n, body, labels, by...)
}

// comment is one Plan comment as the gate reads it: its text, its author's
// login and the author's association with the repo.
type comment struct {
	Body        string `json:"body"`
	Author      string `json:"author"`
	Association string `json:"association"`
}

// planViewBy is planView with each comment's author given.
func planViewBy(t *testing.T, n int, body string, labels []string, comments ...comment) response {
	t.Helper()
	out, err := json.Marshal(map[string]any{
		"url":      "https://github.com/" + exampleRepo + "/issues/" + itoa(n),
		"body":     body,
		"labels":   append([]string{}, labels...),
		"comments": append([]comment{}, comments...),
	})
	mustDo(t, err)
	return response{Match: []string{"issue view " + itoa(n) + " ", "--json url,body,labels,comments"}, Stdout: string(out) + "\n"}
}

// architectedRepo is a Go service whose requirement, owning test, and five
// architecture views keep every check green.
func architectedRepo() map[string]string {
	views := acFixture("docs/architecture/")
	files := map[string]string{
		"go.mod":                  "module example.com/service\n\ngo 1.22\n",
		"src/cmd/service/main.go": "// Command service saves one record.\npackage main\n\nimport \"example.com/service/src/store\"\n\nfunc main() { store.Save() }\n",
		"src/store/store.go":      "// Package store keeps records.\npackage store\n\n// Record is one saved record.\ntype Record struct{ ID int }\n\n// Save keeps one record and returns it.\nfunc Save() Record { return Record{} }\n",
		"src/store/store_test.go": testFile("store", "Test0001_1"),
		acDoc:                     acDocWith("User opens the app"),
	}
	for view, body := range architectedViews() {
		files["docs/architecture/"+view] = body
	}
	return with(views, files)
}

// The palette of each kind, as the shipped view templates style it.
const (
	palettePerson    = `$bgColor="#334155", $borderColor="#1e293b", $fontColor="#ffffff")`
	paletteExternal  = `$bgColor="#64748b", $borderColor="#475569", $fontColor="#ffffff")`
	paletteSystem    = `$bgColor="#1d4ed8", $borderColor="#1e3a8a", $fontColor="#ffffff")`
	paletteContainer = `$bgColor="#0f766e", $borderColor="#115e59", $fontColor="#ffffff")`
	paletteComponent = `$bgColor="#b45309", $borderColor="#92400e", $fontColor="#ffffff")`
	paletteClass     = "fill:#6d28d9,stroke:#4c1d95,color:#ffffff"
)

// c4Styled is the UpdateElementStyle line that gives alias its palette.
func c4Styled(alias, palette string) string {
	return "    UpdateElementStyle(" + alias + ", " + palette + "\n"
}

// archView is a view as its shipped template shapes it: the diagram, its
// legend, then the caption and the callout table of rows.
func archView(title, diagram, legend, caption string, rows ...string) string {
	body := "# " + title + "\n\n```mermaid\n" + diagram + "```\n\n```mermaid\n" + legend + "```\n\n" + caption + "\n\n| # | Description | Why |\n| --- | --- | --- |\n"
	for _, row := range rows {
		body += "| " + row + " | Shared by later checks. |\n"
	}
	return body
}

const archCaption = "Numbers on nodes and letters on arrows match the callout table."

// architectedViews are the five views of architectedRepo, each in its
// level's one form and palette.
func architectedViews() map[string]string {
	return map[string]string{
		acContext: archView("System Context",
			"C4Context\n    title Service\n    Person(user, \"1 User\")\n    System(app, \"2 App\")\n    Rel(user, app, \"A\")\n"+c4Styled("user", palettePerson)+c4Styled("app", paletteSystem),
			"C4Context\n    title Legend\n    Person(p, \"Person\")\n    System(s, \"System\")\n"+c4Styled("p", palettePerson)+c4Styled("s", paletteSystem),
			archCaption, "1 | A person who uses the app.", "2 | The system.", "A | Uses it."),
		acContainer: archView("Containers",
			"C4Container\n    title Service\n    Person(user, \"1 User\")\n    System_Boundary(app, \"2 App\") {\n        Container(service, \"3 Service\", \"Go\", \"src/cmd/service\")\n    }\n    Rel(user, service, \"A\")\n"+c4Styled("user", palettePerson)+c4Styled("service", paletteContainer),
			"C4Container\n    title Legend\n    Person(p, \"Person\")\n    System_Boundary(b, \"System boundary\") {\n        Container(c, \"Container\", \"Tech\")\n    }\n"+c4Styled("p", palettePerson)+c4Styled("c", paletteContainer),
			archCaption, "1 | Same person.", "2 | The app that holds the service.", "3 | Runs requests.", "A | Calls it."),
		acComponent: archView("Components",
			"C4Component\n    title Service\n    Container_Boundary(service, \"1 Service (src/cmd/service)\") {\n        Component(entry, \"2 Entry\", \"Go\")\n        Component(store, \"3 Store\", \"Go\", \"src/store\")\n    }\n    Rel(entry, store, \"A\")\n"+c4Styled("entry", paletteComponent)+c4Styled("store", paletteComponent),
			"C4Component\n    title Legend\n    Container_Boundary(cb, \"Container boundary\") {\n        Component(c, \"Component\", \"Tech\")\n    }\n"+c4Styled("c", paletteComponent),
			archCaption, "1 | The service.", "2 | `src/cmd/service/main.go` starts it.", "3 | Keeps records.", "A | Saves through the store, see [source](../../src/cmd/service/main.go#L6)."),
		acCode: archView("Code",
			"classDiagram\n    class Record[\"1 Record\"] {\n        +ID\n    }\n    class Save[\"2 Save()\"]\n    Save --> Record : A\n    style Record "+paletteClass+"\n    style Save "+paletteClass+"\n",
			"classDiagram\n    class Example\n    Parent <|-- Child : inherits\n    style Example "+paletteClass+"\n    style Parent "+paletteClass+"\n    style Child "+paletteClass+"\n",
			"Numbers on classes and letters on relationships match the callout table.",
			"1 | A saved record, in [store.go](../../src/store/store.go#L5).", "2 | Writes one record, in [store.go](../../src/store/store.go#L8).", "A | Save returns the record."),
		acJourney: "# User Journey\n\n## Legend\n\n```mermaid\nflowchart TB\n    person((1 Person)):::persona\n    classDef persona fill:#fff\n```\n\n" +
			"Numbers match the callout table.\n\n| # | Description | Why |\n| --- | --- | --- |\n| 1 | Persona | A legend row. |\n\n" +
			"## Journey\n\n```mermaid\nflowchart TB\n    user((1 User)):::persona\n    open(2 Open the app, Req 0001):::action\n    user --> open\n    classDef persona fill:#fff\n```\n\n" +
			"Numbers match the callout table.\n\n| # | Description | Type | Why |\n| --- | --- | --- | --- |\n| 1 | Same person. | Persona | Shared by later checks. |\n| 2 | Opens it, `0001.1`. | Action | Shared by later checks. |\n",
	}
}

// notesHead opens a changelog in the house layout.
const notesHead = "# Changelog\n\n## Changes\n\n"

// notesEntry is a timestamped entry that follows the house rules.
const notesEntry = "### 2026.10.09.1500\n\n- Fix: one.\n- Fix: two.\n- Fix: three.\n\n"

// notesRelease is a 1.1.1 release as Release Please writes it: a linked
// version heading, then each subheading with one bullet per scope given.
func notesRelease(subheadings map[string][]string) string {
	s := "## [1.1.1](https://github.com/o/r/compare/v1.1.0...v1.1.1) (2026-10-10)\n\n"
	for _, h := range []string{"Features", "Bug Fixes"} {
		if len(subheadings[h]) == 0 {
			continue
		}
		s += "\n### " + h + "\n\n"
		for _, scope := range subheadings[h] {
			s += "* **" + scope + ":** change it ([#239](https://github.com/o/r/issues/239)) ([c9f4729](https://github.com/o/r/commit/c9f4729))\n"
		}
		s += "\n"
	}
	return s
}

// notesChecks runs the named checks on happyRepo with changelog.
func notesChecks(t *testing.T, changelog, checks string) (string, int) {
	t.Helper()
	repo := example(t, "happyRepo", map[string]string{"CHANGELOG.md": changelog})
	return fitness(t, repo, nil, "--policy=external", "--checks="+checks, "--all")
}

// notesTwo is a release with one feature and one fix.
var notesTwo = map[string][]string{"Features": {"install"}, "Bug Fixes": {"release"}}
