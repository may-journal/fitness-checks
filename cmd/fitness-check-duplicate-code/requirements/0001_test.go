// cspell:ignore mvnw gradlew

// Package requirements proves the duplicate-code check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a
// copied example repo with the files written over it.
package requirements

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 rejects a share past the limit", func(t *testing.T) {
		t.Parallel()
		out, code := dupRun(t, dupPair("src/a.js", "src/b.js", dupClone))
		acceptance.Sees(t, out, code, 1, "ERROR: duplicate code is 7.4% of scanned lines, over the 1.0% threshold", "src/b.js: 6 duplicated lines")
		if strings.Contains(acceptance.Read(out), "src/a.js:") {
			t.Errorf("the first copy holds no duplicated line, so only src/b.js is named:\n%s", out)
		}
	})
	// The case commits a copied block in a script, which duplicate-code never
	// compares with a doc, then in a second doc.
	t.Run("0001.1.2 judges staged files against tracked files", func(t *testing.T) {
		t.Parallel()
		acceptance.PassedWith(t, dupCommit(t, "src/copy.js"), "duplicate-code", "1")
		out := dupCommit(t, "docs/copy.md")
		if !regexp.MustCompile(`│ failed\s+│ 4\s+│`).MatchString(acceptance.Row(out, "duplicate-code")) {
			t.Errorf("a staged doc copying a tracked doc must fail duplicate-code over the 4 docs, got row %q\n%s", acceptance.Row(out, "duplicate-code"), out)
		}
	})
	t.Run("0001.1.3 compares every tracked text file including Markdown or JSON", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", map[string]string{".gitignore": "conf/\n"})
		acceptance.Write(t, repo, dupPair("docs/a.md", "conf/b.json", dupClone))
		acceptance.Git(t, repo, "add", "-f", "docs", "conf")
		out, code := acceptance.Fitness(t, repo, nil, "duplicate-code")
		acceptance.Sees(t, out, code, 1, "ERROR: duplicate code is ")
	})
	// The case puts a 3-line block in one pair and an 18-token block in
	// another, each under one minimum. Two docs share a 5-line run only because
	// it ends in a 120-dash separator row.
	t.Run("0001.1.4 passes short shared blocks", func(t *testing.T) {
		t.Parallel()
		files := dupPair("src/a.js", "src/b.js", dupBlock("w", 3, 20))
		files["src/c.js"] = dupBlock("s", 6, 3) + dupBlock("fc", 24, 3)
		files["src/d.js"] = dupBlock("s", 6, 3) + dupBlock("fd", 24, 3)
		files["docs/a.md"] = "# A\n\n" + dupSeparatorTable + "a | b\n" + dupBlock("da", 24, 3)
		files["docs/b.md"] = "# B\n\n" + dupSeparatorTable + "c | d\n" + dupBlock("db", 24, 3)
		out, code := dupRun(t, files)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 skips lockfiles and generated files", func(t *testing.T) {
		t.Parallel()
		files := dupPair("deps/a.lock", "deps/b.lock", dupClone)
		files[".gitattributes"] = "*.lock linguist-generated\n"
		for name, body := range dupPair("app/gradle.lockfile", "lib/gradle.lockfile", dupClone) {
			files[name] = body
		}
		for name, body := range dupPair("mvnw", "app/gradlew.bat", dupClone) {
			files[name] = body
		}
		out, code := dupRun(t, files)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	// The case tracks a binary pair and leaves a text pair untracked. Both
	// pairs share the block and their tails, so only skipping both passes.
	t.Run("0001.2.2 skips binary and untracked files", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.Example(t, "happyRepo", dupPair("blob/a.dat", "blob/b.dat", "\x00\n"+dupClone))
		acceptance.Write(t, repo, dupPair("src/a.js", "src/b.js", dupClone))
		out, code := acceptance.Fitness(t, repo, nil, "duplicate-code")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	// The case wraps the block in current markers in one pair and deprecated
	// jscpd ones in the other, so older repos keep passing.
	t.Run("0001.2.3 skips blocks between ignore markers including the deprecated ones", func(t *testing.T) {
		t.Parallel()
		current := "// duplicate-code:ignore-start\n" + dupClone + "// duplicate-code:ignore-end\n"
		deprecated := "/* jscpd:ignore-start */\n" + dupBlock("j", 6, 10) + "/* jscpd:ignore-end */\n"
		out, code := dupRun(t, map[string]string{
			"src/a.js": current + dupBlock("fa", 24, 3),
			"src/b.js": current + dupBlock("fb", 24, 3),
			"src/c.js": deprecated + dupBlock("fc", 24, 3),
			"src/d.js": deprecated + dupBlock("fd", 24, 3),
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	// The case adds a template, matching legends and a copied untitled
	// diagram to the two docs of the share past the limit case. A
	// committed third doc scopes the run to Markdown, and the exact share
	// proves the template is not scanned.
	t.Run("0001.2.4 excuses docs repeating a template", func(t *testing.T) {
		t.Parallel()
		files := dupPair("docs/guide.md", "docs/notes.md", dupClone)
		files["docs/templates/guide.md"] = dupClone + dupBlock("ft", 24, 3)
		files["docs/architecture/02-containers.md"] = dupLevelDoc("ca")
		files["docs/architecture/03-components.md"] = dupLevelDoc("cb")
		files["docs/flow-a.md"] = "# Flow\n\n" + dupCopiedDiagram + dupBlock("la", 24, 3)
		files["docs/flow-b.md"] = "# Flow\n\n" + dupLevelDiagram("lb") + dupCopiedDiagram + dupBlock("lb", 24, 3)
		repo := acceptance.Example(t, "happyRepo", files)
		acceptance.Write(t, repo, map[string]string{"docs/new.md": "# New\n\n" + dupBlock("fn", 24, 3)})
		acceptance.Git(t, repo, "add", "docs/new.md")
		out, code := acceptance.Fitness(t, repo, nil, "duplicate-code")
		if code == 0 {
			t.Errorf("the staged doc must fail:\n%s", out)
		}
		dupNames(t, out, []string{"ERROR: duplicate code is 5.6% of scanned lines", "docs/notes.md: 6 duplicated lines", "docs/flow-b.md: 11 duplicated lines"},
			[]string{"docs/architecture/02-containers.md:", "docs/architecture/03-components.md:", "docs/flow-a.md:", "docs/guide.md:", "docs/templates/guide.md:"})
	})
}

// dupBlock builds lines of distinct words, so a shared prefix yields an
// exact copy and different prefixes never match. Six lines of ten words
// plus ";" is 66 tokens, past the 50-token minimum.
func dupBlock(prefix string, lines, perLine int) string {
	var b strings.Builder
	for i := 0; i < lines; i++ {
		for j := 0; j < perLine; j++ {
			fmt.Fprintf(&b, "%s_%d_%d ", prefix, i, j)
		}
		b.WriteString(";\n")
	}
	return b.String()
}

// dupClone is a copied block long enough to count as a clone.
var dupClone = dupBlock("c", 6, 10)

// dupPair puts block at the top of two files that differ after it.
func dupPair(a, b, block string) map[string]string {
	return map[string]string{
		a: block + dupBlock("fa", 24, 3),
		b: block + dupBlock("fb", 24, 3),
	}
}

// dupRun runs `fitness-install -- duplicate-code` on happyRepo with files over it.
func dupRun(t *testing.T, files map[string]string) (string, int) {
	t.Helper()
	return acceptance.Fitness(t, acceptance.Example(t, "happyRepo", files), nil, "duplicate-code")
}

// dupSeparatorTable is two lines of prose and a table's header and separator
// rows. Its 120 dashes carry it past 50 tokens and its last line ends a
// 5-line run.
var dupSeparatorTable = "Both docs list names.\nEach row says why.\n\n| Name | Why |\n| " +
	strings.Repeat("-", 60) + " | " + strings.Repeat("-", 60) + " |\n"

// dupCommit tracks a doc holding dupClone, stages file holding it too, and
// returns what the run over the staged file printed.
func dupCommit(t *testing.T, file string) string {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{"docs/guide.md": dupClone + dupBlock("fg", 24, 3)})
	acceptance.Write(t, repo, map[string]string{file: dupClone + dupBlock("fc", 24, 3)})
	acceptance.Git(t, repo, "add", file)
	out, _ := acceptance.Fitness(t, repo, nil, "duplicate-code")
	return out
}

// dupCopiedDiagram is an untitled diagram with no callouts, then a 48-token
// line. The lexer reads the fenced block as one token, so only a copy whose
// body stays reaches both minimums.
var dupCopiedDiagram = "```mermaid\nflowchart TB\n" + dupBlock("m", 6, 3) + "```\n\n" + dupBlock("t", 1, 47)

// dupNames asserts out holds each of wants and none of the clean paths, which
// hold only lines a doc must repeat.
func dupNames(t *testing.T, out string, wants, clean []string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(acceptance.Read(out), want) {
			t.Errorf("the commit output is missing %q:\n%s", want, out)
		}
	}
	for _, path := range clean {
		if regexp.MustCompile(regexp.QuoteMeta(path) + ` \d+ duplicated lines`).MatchString(acceptance.Read(out)) {
			t.Errorf("%s holds no duplicated line:\n%s", path, out)
		}
	}
}

// dupLevelDoc is one C4 level doc: a numbered diagram, the shared legend
// and table, then prose of its own.
func dupLevelDoc(prefix string) string {
	return "# Level\n\n" + dupLevelDiagram(prefix) + dupLegend + dupTable(prefix) + "\n" + dupBlock(prefix, 24, 3)
}

// dupTable is a level's callout rows, which name the level.
func dupTable(prefix string) string {
	return "| 1 | The " + prefix + " user | Starts " + prefix + " |\n| 2 | The " + prefix + " app | Serves " + prefix + " |\n| A | Opens " + prefix + " | Daily |\n"
}

// dupLevelDiagram is a numbered C4 diagram whose only repeated lines are
// its style statements.
func dupLevelDiagram(prefix string) string {
	return "```mermaid\nC4Container\n    Person(" + prefix + ", \"1 " + prefix + " user\")\n    Container(" + prefix + "2, \"2 " + prefix + " app\", \"Go\", \"Runs " + prefix + "\")\n" +
		"    Rel(" + prefix + ", " + prefix + "2, \"A\")\n" +
		"    UpdateElementStyle(" + prefix + ", $bgColor=\"#334155\", $borderColor=\"#1e293b\", $fontColor=\"#ffffff\")\n" +
		"    UpdateElementStyle(" + prefix + "2, $bgColor=\"#0f766e\", $borderColor=\"#115e59\", $fontColor=\"#ffffff\")\n```\n\n"
}

// dupLegend is the legend diagram every C4 level repeats, with its table.
const dupLegend = "```mermaid\nC4Container\n    title Legend\n\n    Person(p, \"Person\", \"Human or agent actor\")\n    System_Boundary(b, \"System boundary\") {\n        Container(c, \"Container\", \"Tech\", \"Runnable or deployable unit\")\n    }\n    System_Ext(e, \"External\", \"Outside our boundary\")\n\n    UpdateElementStyle(p, $bgColor=\"#334155\", $borderColor=\"#1e293b\", $fontColor=\"#ffffff\")\n    UpdateElementStyle(c, $bgColor=\"#0f766e\", $borderColor=\"#115e59\", $fontColor=\"#ffffff\")\n    UpdateElementStyle(e, $bgColor=\"#64748b\", $borderColor=\"#475569\", $fontColor=\"#ffffff\")\n```\n\nNumbers on nodes and letters on arrows match the callout table.\n\n| # | Description | Why |\n| --- | --- | --- |\n"
