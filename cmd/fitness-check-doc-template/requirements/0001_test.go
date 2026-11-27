// Package requirements proves the doc-template check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo with the docs and their templates written over it.
package requirements

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 section the template lacks or has", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.1.1", "8807486b", func(_ req.A) {
			out, code := docTemplate(t, map[string]string{
				"docs/adr/template.md":   docTemplateADR,
				"docs/adr/0001-extra.md": "# 0001 Extra\n\n## Context\n\nc\n\n## Status\n\ns\n\n## Decision\n\nd\n\n## Consequences\n\ne\n",
				"docs/adr/0002-short.md": "# 0002 Short\n\n## Context\n\nc\n\n## Consequences\n\ne\n",
			})
			acceptance.Sees(t, out, code, 1,
				"docs/adr/0001-extra.md: Unexpected `## Status` section",
				"docs/adr/0002-short.md: Missing `## Decision` section")
		})
	})
	t.Run("0001.1.2 nested and named templates nearest wins", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.2.1", "1439da03", func(_ req.A) {
			out, code := docTemplate(t, map[string]string{
				"docs/templates/adr.md":           docTemplateADR,
				"docs/deep/note.template.md":      "# Note\n\n## Summary\n\n<!-- gist -->\n",
				"docs/deep/sub/0001-note.md":      "# 0001 Note\n\n## Context\n\nc\n",
				"docs/views/templates/04-code.md": docTemplateNote,
				"docs/views/04-code.md":           "# Code\n\n## Summary\n\ns\n",
				"docs/views/notes.md":             "# Notes\n\n## Open threads\n\no\n",
			})
			acceptance.Sees(t, out, code, 1,
				"docs/deep/sub/0001-note.md: Missing `## Summary` section",
				"docs/views/notes.md: Missing `## Context` section",
				"docs/views/notes.md: Unexpected `## Open threads` section")
		})
	})
	t.Run("0001.1.3 doc matching none of several templates", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.3.1", "0fac6c6f", func(_ req.A) {
			out, code := docTemplate(t, map[string]string{
				"docs/adr/adr.template.md":  docTemplateADR,
				"docs/adr/note.template.md": docTemplateNote,
				"docs/adr/0001-bad.md":      "# 0001 Bad\n\n## Other\n\no\n",
			})
			acceptance.Sees(t, out, code, 1,
				"docs/adr/0001-bad.md: sections match none of its templates: docs/adr/adr.template.md, docs/adr/note.template.md")
		})
	})
	t.Run("0001.1.4 github templates do not govern files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.1.4.1", "3faa3bbd", func(_ req.A) {
			out, code := docTemplate(t, map[string]string{
				"docs/arch/templates/adr.md":       docTemplateADR,
				"docs/arch/templates/note.md":      docTemplateNote,
				"docs/arch/04-code.md":             "# Code\n\n## Summary\n\ns\n",
				".github/PULL_REQUEST_TEMPLATE.md": "## Background\n\n<!-- context -->\n",
				".github/other.md":                 "# Other\n\n## Totally\n\ndifferent\n",
			})
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	out, code := docTemplate(t, map[string]string{
		"docs/adr/template.md":             "# NNNN Title\n\n## Context\n\n<!-- why -->\n\n## Decision\n\nno comment here\n",
		".github/PULL_REQUEST_TEMPLATE.md": "## Background\n\nno comment\n",
		".github/ISSUE_TEMPLATE/plan.md":   "## Goal\n\nno comment\n",
	})
	req.Then(t, "0001.2.1.1.1", "d52ac9aa", func(a req.A) {
		acceptance.Sees(t, out, code, 1, "docs/adr/template.md:")
	})
	req.Then(t, "0001.2.1.1.2", "a722a8e7", func(a req.A) {
		acceptance.Sees(t, out, code, 1, "docs/adr/template.md: section `## Decision` needs a guiding <!-- comment -->")
	})
	req.Then(t, "0001.2.1.1.3", "9c8b3bad", func(a req.A) {
		acceptance.Sees(t, out, code, 1,
			".github/PULL_REQUEST_TEMPLATE.md: section `## Background` needs a guiding <!-- comment -->",
			".github/ISSUE_TEMPLATE/plan.md: section `## Goal` needs a guiding <!-- comment -->")
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 view lacks legend or callout structure", func(t *testing.T) {
		t.Parallel()
		files := dtCopied(t)
		dtEdits(t, files,
			[3]string{acCode, "```mermaid\nclassDiagram\n    class Example", "```mermaid\nflowchart TB\n    class Example"},
			[3]string{acJourney, "## Journey", "## Steps"},
		)
		out1, code1 := docTemplate(t, files)

		files = dtCopied(t)
		dtEdits(t, files,
			[3]string{acJourney, "| Type    |", "| Kind    |"},
			[3]string{acCode, "Numbers on classes", "Labels on classes"},
		)
		files[dtViews+acContainer] = files[dtViews+acContainer][:strings.Index(files[dtViews+acContainer], "| #")]
		journey := dtLine(t, files, acJourney, "| #   | Description          | Kind")
		table := dtLine(t, files, acCode, "| #")
		out2, code2 := docTemplate(t, files)
		journeyLine, tableLine := journey+":", table+":"

		files = dtShipped(t)
		dtEdits(t, files,
			[3]string{acContainer, "C4Container\n    title Product", "C4Context\n    title Product"},
			[3]string{acComponent, "C4Component\n    title Legend", "C4Container\n    title Legend"},
			[3]string{acCode, "```mermaid\nclassDiagram\n    class Record", "```mermaid\nflowchart TB\n    class Record"},
			[3]string{acJourney, "flowchart TB", "flowchart LR"},
		)
		containerLine := dtLine(t, files, acContainer, "```mermaid") + ":"
		codeLine := dtLine(t, files, acCode, "```mermaid") + ":"
		journeyType := dtLine(t, files, acJourney, "```mermaid") + ":"
		out3, code3 := docTemplate(t, files)

		req.Then(t, "0001.3.1.1.1", "aabd0f94", func(a req.A) {
			acceptance.Sees(t, out1, code1, 1, dtViews+acCode+":", dtViews+acJourney+":")
			acceptance.Sees(t, out2, code2, 1, journeyLine, tableLine, dtViews+acContainer+":")
			acceptance.Sees(t, out3, code3, 1, containerLine, dtViews+acComponent+":", codeLine, journeyType)
		})
		req.Then(t, "0001.3.1.1.2", "d7300cb5", func(a req.A) {
			acceptance.Sees(t, out1, code1, 1,
				": the opening needs a classDiagram legend as "+dtTemplates+acCode+" has",
				": sections must open with `## Legend`, `## Journey` as "+dtTemplates+acJourney+" does")
			acceptance.Sees(t, out2, code2, 1,
				": callout table columns #, Description, Kind, Why must be #, Description, Type, Why as "+dtTemplates+acJourney+" has",
				`: callout table needs the caption "Numbers on classes and letters on relationships match the callout table" above it as `+dtTemplates+acCode+" has",
				": the opening needs a callout table with columns #, Description, Why as "+dtTemplates+acContainer+" has")
			acceptance.Sees(t, out3, code3, 1,
				": diagram type C4Context must be C4Container as "+dtShippedName(acContainer)+" has",
				": the opening needs a C4Component legend as "+dtShippedName(acComponent)+" has",
				": legend type C4Container must be C4Component as "+dtShippedName(acComponent)+" has",
				": diagram type flowchart TB must be classDiagram as "+dtShippedName(acCode)+" has",
				": legend type flowchart LR must be flowchart TB as "+dtShippedName(acJourney)+" has")
		})
	})
	t.Run("0001.3.2 views from shipped architecture templates", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.3.1.2.1", "f278274c", func(_ req.A) {
			files := dtCopied(t)
			files[dtViews+acContainer] += "\n## Setup\n\nHow a consumer runs the app.\n"
			files[dtViews+acJourney] += "\n## Notes\n\nWhat the journey leaves out.\n"
			dtEdit(t, files, acComponent, "match the callout table.", "match the callout table. <!-- kept in step with the diagram -->")
			files[dtViews+"notes.md"] = "# Notes\n\n## Open threads\n\nWhat the views leave out.\n"
			files[dtViews+"archive/"+acCode] = "# Old code\n\n## Before\n\nThe first layout.\n"
			out, code := docTemplate(t, files)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWith(t, out, "doc-template", "14")

			files = dtShipped(t)
			files[dtViews+acContainer] += "\n## Setup\n\nHow a consumer runs the app.\n"
			dtSees(t, files, 0, "All 1 checks passed")
			dtEdit(t, files, acJourney, "## Journey", "## Steps")
			dtEdit(t, files, acCode, "Numbers on classes", "Labels on classes")
			out, code = docTemplate(t, files)
			acceptance.Sees(t, out, code, 1,
				dtViews+acJourney+": sections must open with `## Legend`, `## Journey` as "+dtShippedName(acJourney)+" does",
				`: callout table needs the caption "Numbers on classes and letters on relationships match the callout table" above it as `+dtShippedName(acCode)+" has")
		})
	})
	t.Run("0001.3.3 view styles differ from template naming the element and palette", func(t *testing.T) {
		t.Parallel()
		files := dtShipped(t)
		dtEdit(t, files, acContext, "    UpdateElementStyle(user,",
			"    SystemDb_Ext(db, \"4 Records\", \"Keeps records.\")\n    UpdateElementStyle(db, $fontColor=\"#FFFFFF\", $bgColor=\"#64748b\",  $borderColor=\"#475569\")\n    UpdateLayoutConfig($c4ShapeInRow=\"3\")\n    Lay_R(user, product)\n    UpdateElementStyle(user,")
		dtEdits(t, files,
			[3]string{acCode, "    save --> Record : A", "    save --> Record : A\n    Store : +save()\n    <<interface>> Store\n    style Store " + paletteClass},
			[3]string{acContainer, `Container(app, "3 App", "Language", "src/app")`, `Container(app, "3 App", $techn="Language")`},
		)
		kept, keptCode := docTemplate(t, files)
		dtEdits(t, files,
			[3]string{acContext, "UpdateElementStyle(product, $bgColor=\"#1d4ed8\"", "UpdateElementStyle(product, $bgColor=\"#ff0000\""},
			[3]string{acContainer, "    UpdateElementStyle(app, " + paletteContainer + "\n", ""},
			[3]string{acComponent, "Container_Boundary(app,", "System_Boundary(app,"},
			[3]string{acComponent, "    Rel(user, entry, \"A\")", "    Rel(user, entry, \"A\")\n    UpdateRelStyle(user, entry, $textColor=\"red\")"},
			[3]string{acCode, "    save --> Record : A", "    save --> Record : A\n    classDef hot fill:#ff0000"},
		)
		out, code := docTemplate(t, files)
		boundary := dtLine(t, files, acComponent, "    System_Boundary(app,")

		files = dtShipped(t)
		dtEdits(t, files,
			[3]string{acContext, "C4Context\n    title Product\n", "C4Context\n"},
			[3]string{acComponent, "C4Component\n    title Legend\n", "C4Component\n"},
			[3]string{acContainer, `Container(app, "3 App", "Language", "src/app")`, `Container(app, "3 App", "", "src/app")`},
			[3]string{acComponent, `Component(store, "4 Store", "Language",`, `Component(store, "4 Store", " ",`},
			[3]string{acContext, "    System_Ext(e, \"External\", \"Outside our boundary\")\n", ""},
			[3]string{acContext, c4Styled("e", paletteExternal), ""},
			[3]string{acContainer, "    System_Ext(service, \"4 Outside service\"", "    Person(service, \"4 Outside service\""},
			[3]string{acContainer, c4Styled("service", paletteExternal), c4Styled("service", palettePerson)},
		)
		titleLine := dtLine(t, files, acContext, "```mermaid")
		containerLine := dtLine(t, files, acContainer, "        Container(app,")
		componentLine := dtLine(t, files, acComponent, "        Component(store,")
		bare, bareCode := docTemplate(t, files)

		req.Then(t, "0001.3.1.3.1", "c1924e39", func(a req.A) {
			acceptance.Sees(t, out, code, 1,
				"Container app needs `UpdateElementStyle(app, "+paletteContainer+"`",
				boundary+": System_Boundary app is a kind "+dtShippedName(acComponent)+" does not use; use Person, Container_Boundary, Component, System_Ext",
				"`UpdateRelStyle(user, entry, $textColor=\"red\")` is a style "+dtShippedName(acComponent)+" does not use",
				"`classDef hot fill:#ff0000` is a style "+dtShippedName(acCode)+" does not use")
			acceptance.Sees(t, bare, bareCode, 1,
				titleLine+": C4Context diagram needs a `title` line as "+dtShippedName(acContext)+" has",
				": C4Component legend needs a `title` line as "+dtShippedName(acComponent)+" has",
				containerLine+": Container app needs a technology as "+dtShippedName(acContainer)+" gives each Container",
				componentLine+": Component store needs a technology as "+dtShippedName(acComponent)+" gives each Component",
				dtViews+acContext+":", "legend needs a System_Ext sample, a kind the diagram uses",
				dtViews+acContainer+":", "legend shows System_Ext, a kind the diagram does not use")
		})
		req.Then(t, "0001.3.1.3.2", "d8c4422f", func(a req.A) {
			acceptance.Sees(t, kept, keptCode, 0, "All 1 checks passed")
			acceptance.Sees(t, out, code, 1,
				"`UpdateElementStyle(product, $bgColor=\"#ff0000\", $borderColor=\"#1e3a8a\", $fontColor=\"#ffffff\")` must be `UpdateElementStyle(product, "+paletteSystem+"` as "+dtShippedName(acContext)+" styles each System",
				"Container app needs `UpdateElementStyle(app, "+paletteContainer+"` as "+dtShippedName(acContainer)+" styles each Container")
			for _, styled := range []string{"Person user", "Component entry", "class save", "class Record"} {
				a.True(!strings.Contains(acceptance.Read(out), styled), "%s keeps its palette and must not fail, got:\n%s", styled, out)
			}
		})
	})
}

// docTemplate runs `fitness-install -- doc-template` on happyRepo with files
// written over it.
func docTemplate(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "doc-template")
}

// docTemplateNote guides its one section with a comment.
const docTemplateNote = "# Note\n\n## Summary\n\n<!-- gist -->\n"

// docTemplateADR guides each of its three sections with a comment.
const docTemplateADR = "# NNNN Title\n\n## Context\n\n<!-- why -->\n\n## Decision\n\n<!-- what -->\n\n## Consequences\n\n<!-- fallout -->\n"

const (
	dtViews     = "docs/architecture/"
	dtTemplates = dtViews + "templates/"
)

// dtShipped is a user's docs/architecture after following the
// architecture-contract README: each view copied from its template, templates
// left behind.
func dtShipped(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	for _, view := range []string{acContext, acContainer, acComponent, acCode, acJourney} {
		body, err := os.ReadFile(filepath.Join("..", "..", "..", "internal", "architecture", "templates", view))
		acceptance.MustDo(t, err)
		files[dtViews+view] = string(body)
	}
	return files
}

// dtCopied is dtShipped with the five templates also copied into
// docs/architecture/templates, as repos did before fitness shipped them.
func dtCopied(t *testing.T) map[string]string {
	t.Helper()
	files := dtShipped(t)
	for _, view := range []string{acContext, acContainer, acComponent, acCode, acJourney} {
		files[dtTemplates+view] = files[dtViews+view]
	}
	return files
}

// dtShippedName is how a message names a view's shipped template.
func dtShippedName(view string) string {
	return "the shipped " + view
}

// dtEdit replaces old with new in one copied view, failing when old is not
// there to replace.
func dtEdit(t *testing.T, files map[string]string, view, old, new string) {
	t.Helper()
	if !strings.Contains(files[dtViews+view], old) {
		t.Fatalf("%s has no %q", view, old)
	}
	files[dtViews+view] = strings.Replace(files[dtViews+view], old, new, 1)
}

// dtLine is the 1-based line of the first line in a view starting with
// prefix, as an error names it.
func dtLine(t *testing.T, files map[string]string, view, prefix string) string {
	t.Helper()
	for i, line := range strings.Split(files[dtViews+view], "\n") {
		if strings.HasPrefix(line, prefix) {
			return dtViews + view + ":" + strconv.Itoa(i+1)
		}
	}
	t.Fatalf("%s has no line starting %q", view, prefix)
	return ""
}

// The views the architecture contract ships.
const (
	acContext   = "01-system-context.md"
	acContainer = "02-containers.md"
	acComponent = "03-components.md"
	acCode      = "04-code.md"
	acJourney   = "05-user-journey.md"
)

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

// dtEdits applies each {view, old, new} edit to files in order.
func dtEdits(t *testing.T, files map[string]string, edits ...[3]string) {
	t.Helper()
	for _, e := range edits {
		dtEdit(t, files, e[0], e[1], e[2])
	}
}

// dtSees runs doc-template on files and asserts its exit code and that its output holds each message.
func dtSees(t *testing.T, files map[string]string, want int, msgs ...string) {
	t.Helper()
	out, code := docTemplate(t, files)
	acceptance.Sees(t, out, code, want, msgs...)
}
