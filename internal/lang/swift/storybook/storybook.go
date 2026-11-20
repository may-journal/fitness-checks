// Package storybook judges the swift-storybook-requirements check. Every
// Storybook page a Swift app declares is named by the acceptance it shows,
// and every acceptance has a page. A page is a `#Preview` or a `BookPage`.
// It judges only a repo that uses Storybook: one whose Swift code imports a
// Storybook module.
package storybook

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/comments"
	"github.com/may-journal/fitness-runner/go/pkg/lang/swift"
	"github.com/may-journal/fitness-runner/go/pkg/walkfs"
)

// Name is the check's name.
const Name = "swift-storybook-requirements"

// NoStorybook is the note a repo that imports no Storybook module passes
// with, since its acceptances need no pages.
const NoStorybook = "no Storybook in this repo: no Swift file imports a Storybook module or declares a BookPage, so acceptances need no pages"

var (
	// pageStart matches where a page is declared in a code view: the
	// `#Preview` macro, or a `BookPage(` call.
	pageStart = regexp.MustCompile(`#Preview\b|\bBookPage\s*\(`)
	// pageName is a name that starts with one acceptance ID, alone or
	// followed by a space and a label.
	pageName = regexp.MustCompile(`^([0-9]{4}\.[0-9]+)(?: \S.*)?$`)
	// idLike is a leading run that looks like an ID, which a suggested name
	// drops before its label.
	idLike = regexp.MustCompile(`^[0-9.:\s-]+`)
	// titleLabel is a `title:` argument label at the start of args.
	titleLabel = regexp.MustCompile(`^\s*title\s*:`)
	// bookPage matches a `BookPage(` call, which only Storybook declares.
	bookPage = regexp.MustCompile(`\bBookPage\s*\(`)
	// storybookImport matches an import of a Storybook module in a code
	// view, with any attributes or import kind. The modules are Storybook,
	// StorybookSandbox, StorybookKit and SwiftUIStorybook.
	storybookImport = regexp.MustCompile(`(?m)(?:^|;)[ \t]*(?:@\w+(?:\([^)\n]*\))?[ \t]+)*import[ \t]+(?:(?:struct|class|enum|protocol|typealias|func|var|let)[ \t]+)?(?:Storybook|StorybookSandbox|StorybookKit|SwiftUIStorybook)\b`)
)

// page is one declared Storybook page.
type page struct {
	file string
	line int
	// name is the page's name, or the expression it is built from.
	name string
	// kind says how the page is named.
	kind nameKind
}

// nameKind says how a page is named.
type nameKind int

const (
	// unnamed is a page with no name at all.
	unnamed nameKind = iota
	// literal is a page named by a plain string literal.
	literal
	// computed is a page named by a variable or an interpolated literal.
	computed
)

func (p page) at() string { return fmt.Sprintf("%s:%d", p.file, p.line) }

// Run judges every Swift file in the repo at root, even in a scoped run,
// since names repeat across files. A repo without docs or Storybook passes.
func Run(root string, _ []string) (checkkit.Result, error) {
	docs := requirement.DocPaths(root)
	files := walkfs.FilesByExt(root, swift.Ext)
	if len(docs) == 0 || len(files) == 0 {
		return checkkit.Pass(0), nil
	}
	pages, uses, err := readPages(root, files)
	if err != nil {
		return checkkit.Pass(0), err
	}
	if !uses {
		return checkkit.Info(len(files), NoStorybook), nil
	}
	return verdict(len(files), judgeAgainst(root, docs, pages)), nil
}

// judgeAgainst judges pages against the acceptances the docs define, then
// each acceptance no page names. A doc it cannot read is a finding too.
func judgeAgainst(root string, docs []string, pages []page) []string {
	loaded, _, skipped := requirement.Load(root, docs)
	findings := append(skipped, judge(pages, acceptances(loaded))...)
	return append(findings, missingPages(loaded, pages)...)
}

// missingPages returns a finding at each acceptance that no page's literal
// name starts with, where its doc declares it.
func missingPages(docs []requirement.Document, pages []page) []string {
	shown := shownIDs(pages)
	var out []string
	for _, d := range docs {
		for _, a := range d.Acceptances {
			if !shown[a.ID] {
				out = append(out, fmt.Sprintf("%s:%d - acceptance %s has no Storybook page; add one named %q", d.Path, a.Line, a.ID, a.ID+" <label>"))
			}
		}
	}
	return out
}

// shownIDs is the set of acceptance IDs that page literal names start with.
func shownIDs(pages []page) map[string]bool {
	shown := map[string]bool{}
	for _, p := range pages {
		if m := pageName.FindStringSubmatch(p.name); p.kind == literal && m != nil {
			shown[m[1]] = true
		}
	}
	return shown
}

// verdict passes with no findings, else fails with them.
func verdict(files int, findings []string) checkkit.Result {
	if len(findings) > 0 {
		return checkkit.Fail(files, findings...)
	}
	return checkkit.Pass(files)
}

// acceptances lists the acceptance IDs the docs define, sorted.
func acceptances(docs []requirement.Document) []string {
	var ids []string
	for _, d := range docs {
		for _, a := range d.Acceptances {
			ids = append(ids, a.ID)
		}
	}
	slices.Sort(ids)
	return ids
}

// readPages reads the pages each Swift file declares, in file order. It also
// says whether any file imports Storybook, as an app with no pages does.
func readPages(root string, files []string) ([]page, bool, error) {
	var pages []page
	uses := false
	for _, file := range files {
		src, err := comments.CodeFile(root, file)
		if err != nil {
			return nil, false, err
		}
		filed, imports := filePages(file, string(src))
		pages = append(pages, filed...)
		uses = uses || imports
	}
	return pages, uses, nil
}

// filePages returns the pages one file declares, and whether it imports
// Storybook. It skips comments and strings, which its code view blanks.
func filePages(file, src string) ([]page, bool) {
	code, literals := swift.Literals(src)
	starts := map[int]swift.Literal{}
	for _, lit := range literals {
		starts[lit.Start] = lit
	}
	var pages []page
	for _, m := range pageStart.FindAllStringIndex(code, -1) {
		p := page{file: file, line: strings.Count(code[:m[0]], "\n") + 1}
		if strings.HasPrefix(code[m[0]:], "#Preview") {
			p.name, p.kind = previewName(src, m[1], starts)
		} else {
			p.name, p.kind = titleName(src, code, m[1], starts)
		}
		pages = append(pages, p)
	}
	return pages, storybookImport.MatchString(code) || bookPage.MatchString(code)
}

// previewName reads the literal a `#Preview` ending at i passes first, from
// src since the code view blanks it. Without one the preview has no name.
func previewName(src string, i int, starts map[int]swift.Literal) (string, nameKind) {
	i = skipSpace(src, i)
	if i >= len(src) || src[i] != '(' {
		return "", unnamed
	}
	lit, ok := starts[skipSpace(src, i+1)]
	if !ok {
		return "", unnamed
	}
	return named(lit)
}

// titleName reads the `title:` argument of the `BookPage(` call opened at
// i. A label of a call nested inside it does not count.
func titleName(src, code string, i int, starts map[int]swift.Literal) (string, nameKind) {
	end := closing(code, i)
	for _, at := range topLevelArgs(code, i, end) {
		if m := titleLabel.FindStringIndex(code[at:end]); m != nil {
			return titleArg(src, skipSpace(src, at+m[1]), end, starts)
		}
	}
	return "", unnamed
}

// titleArg is the name a `title:` argument starting at i gives: its literal,
// or the expression it is built from.
func titleArg(src string, i, end int, starts map[int]swift.Literal) (string, nameKind) {
	if lit, ok := starts[i]; ok {
		return named(lit)
	}
	return strings.TrimSpace(firstArg(src[i:end])), computed
}

// named is a literal's name. An interpolated literal is computed, since its
// lexed text leaves the interpolations out.
func named(lit swift.Literal) (string, nameKind) {
	if strings.Contains(lit.Source, `\(`) {
		return lit.Source, computed
	}
	return lit.Text, literal
}

// judge returns a finding for each page whose name is wrong or repeats an
// earlier page's name.
func judge(pages []page, ids []string) []string {
	var out []string
	seen := map[string]page{}
	for _, p := range pages {
		finding := nameFinding(p, ids)
		if first, ok := seen[p.name]; finding == "" && ok {
			finding = fmt.Sprintf("page %q repeats the name of the page at %s; give each page its own label after the ID", p.name, first.at())
		}
		if finding != "" {
			out = append(out, p.at()+" - "+finding)
			continue
		}
		seen[p.name] = p
	}
	return out
}

// nameFinding is a page's finding about its own name, or "".
func nameFinding(p page, ids []string) string {
	switch p.kind {
	case unnamed:
		return fmt.Sprintf("page has no name; name it by the acceptance it shows, such as #Preview(%q)", suggestion(ids, ""))
	case computed:
		return fmt.Sprintf("page name %s is not a plain literal; write the acceptance name, such as %q", p.name, suggestion(ids, ""))
	}
	return idFinding(p.name, ids, suggestion(ids, p.name))
}

// idFinding is the finding for a name that does not start with a defined
// acceptance ID, or "".
func idFinding(name string, ids []string, example string) string {
	m := pageName.FindStringSubmatch(name)
	if m == nil {
		return fmt.Sprintf("page %q does not start with an acceptance ID; name it such as %q", name, example)
	}
	if _, ok := slices.BinarySearch(ids, m[1]); !ok {
		return fmt.Sprintf("page %q names %s, which no doc in %s defines; use an acceptance ID from those docs", name, m[1], strings.TrimSuffix(requirement.Dir, "/"))
	}
	return ""
}

// suggestion is a name to use: the repo's first acceptance ID, then the
// page's label with any ID-like start dropped.
func suggestion(ids []string, name string) string {
	label := strings.TrimSpace(idLike.ReplaceAllString(name, ""))
	if label == "" {
		label = "label"
	}
	if len(ids) == 0 {
		return "0001.1 " + label
	}
	return ids[0] + " " + label
}

// topLevelArgs returns where each argument of the call between i and end
// starts, at the call's own depth.
func topLevelArgs(code string, i, end int) []int {
	starts := []int{i}
	depth := 0
	for j := i; j < end; j++ {
		depth += nesting(code[j])
		if code[j] == ',' && depth == 0 {
			starts = append(starts, j+1)
		}
	}
	return starts
}

// skipSpace returns the first index at or after i that is not white space.
func skipSpace(code string, i int) int {
	for i < len(code) && strings.ContainsRune(" \t\r\n", rune(code[i])) {
		i++
	}
	return i
}

// closing returns the index of the parenthesis that closes the call whose
// arguments start at i, or the end of code.
func closing(code string, i int) int {
	depth := 1
	for ; i < len(code); i++ {
		if depth += nesting(code[i]); depth == 0 {
			return i
		}
	}
	return len(code)
}

// firstArg returns args up to the first comma at its own depth.
func firstArg(args string) string {
	if starts := topLevelArgs(args, 0, len(args)); len(starts) > 1 {
		return args[:starts[1]-1]
	}
	return args
}

// nesting is how far a bracket moves the depth: one in, one out, or none.
func nesting(c byte) int {
	switch c {
	case '(', '[', '{':
		return 1
	case ')', ']', '}':
		return -1
	}
	return 0
}
