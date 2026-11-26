package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// assertionWords are the words an assertion call's name opens with, such as
// t.Fatalf, require.Equal, Sees or expectEqual. A word counts whole: "check"
// opens checkEqual and not checked.
const assertionWords = "assert|require|expect|verify|check|sees|must|ensure|fail|fatal|equal|contains|passed|shows"

// assertionName matches a call name, its first letter lowered, that opens
// with an assertion word and ends there or at the next word.
var assertionName = regexp.MustCompile(`^(?:` + assertionWords + `)(?:f|now)?(?:$|[A-Z0-9_])`)

// assertionText matches an assertion call in code other than Go: a call
// named like assertEquals, XCTAssertTrue or expect, or a Swift Testing macro.
var assertionText = regexp.MustCompile(`\b(?:(?i:assert)|XCTAssert|XCTFail|(?i:require|expect|verify|fail))\w*\s*\(|#expect|#require`)

// plainPackages are the packages whose calls only read or shape values, so
// strings.Contains is no assertion.
var plainPackages = map[string]bool{"strings": true, "slices": true, "bytes": true, "errors": true, "fmt": true, "os": true, "filepath": true, "maps": true, "regexp": true, "template": true}

// assertionErrors fails each case named for a THEN whose body holds no
// assertion call. A case row of a table is judged by its test's whole body.
func assertionErrors(root string, tests []reqTest) []string {
	files := map[string]*sourceFile{}
	var errs []string
	for _, t := range tests {
		if casesReported[t.module.Name] {
			errs = append(errs, silentCases(files, root, t)...)
		}
	}
	return errs
}

// silentCases fails each case of t named for a THEN that holds no assertion.
func silentCases(files map[string]*sourceFile, root string, t reqTest) []string {
	var errs []string
	for i, c := range t.Cases {
		if !caseID.MatchString(c.Name) {
			continue
		}
		src := sourceOf(files, root, c.Loc)
		switch {
		case src == nil:
			errs = append(errs, fmt.Sprintf("%s: case %q is in a file this check cannot read, so its assertion call is unchecked", c.Loc, c.Name))
		case !src.asserts(t, i):
			errs = append(errs, fmt.Sprintf("%s: case %q has no assertion call; assert the result its THEN states", c.Loc, c.Name))
		}
	}
	return errs
}

// sourceFile is one test file's text, and its Go syntax tree when it is Go.
type sourceFile struct {
	text string
	fset *token.FileSet
	tree *ast.File
}

// sourceOf reads the file a case location names once. A file it cannot read
// is judged by the language's own scan, so it returns nil.
func sourceOf(files map[string]*sourceFile, root, loc string) *sourceFile {
	name, _ := splitLoc(loc)
	if src, ok := files[name]; ok {
		return src
	}
	raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		files[name] = nil
		return nil
	}
	src := &sourceFile{text: string(raw), fset: token.NewFileSet()}
	if strings.HasSuffix(name, ".go") {
		src.tree, _ = parser.ParseFile(src.fset, name, raw, parser.SkipObjectResolution)
	}
	files[name] = src
	return src
}

// splitLoc splits a file:line location.
func splitLoc(loc string) (string, int) {
	i := strings.LastIndex(loc, ":")
	if i < 0 {
		return loc, 0
	}
	n, _ := strconv.Atoi(loc[i+1:])
	return loc[:i], n
}

// asserts reports whether case i of t holds an assertion call.
func (s *sourceFile) asserts(t reqTest, i int) bool {
	_, line := splitLoc(t.Cases[i].Loc)
	if s.tree != nil {
		return s.goAsserts(line)
	}
	return assertionText.MatchString(s.region(t, line, i))
}

// goAsserts reports whether the t.Run literal declared on line holds an
// assertion call. A case with no such literal is a table row, which the
// function around it runs, so that function's body is judged.
func (s *sourceFile) goAsserts(line int) bool {
	var body ast.Node
	ast.Inspect(s.tree, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok && body == nil {
			if lit := runLiteral(call); lit != nil && s.fset.Position(call.Pos()).Line == line {
				body = lit.Body
			}
		}
		return body == nil
	})
	if body == nil {
		body = s.enclosing(line)
	}
	return body == nil || holdsAssertion(body)
}

// runLiteral is the function literal of a call such as t.Run("name", func…).
func runLiteral(call *ast.CallExpr) *ast.FuncLit {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Run" || len(call.Args) != 2 {
		return nil
	}
	lit, _ := call.Args[1].(*ast.FuncLit)
	return lit
}

// enclosing is the body of the function declaration holding line.
func (s *sourceFile) enclosing(line int) ast.Node {
	for _, d := range s.tree.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && s.spans(fn, line) {
			return fn.Body
		}
	}
	return nil
}

// spans reports whether function fn has a body that holds line.
func (s *sourceFile) spans(fn *ast.FuncDecl, line int) bool {
	return fn.Body != nil && s.fset.Position(fn.Pos()).Line <= line && line <= s.fset.Position(fn.End()).Line
}

// holdsAssertion reports whether any call under n is named like an
// assertion.
func holdsAssertion(n ast.Node) bool {
	found := false
	ast.Inspect(n, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok && assertingCall(call) {
			found = true
		}
		return !found
	})
	return found
}

// assertingCall reports whether call's name opens with an assertion word.
func assertingCall(call *ast.CallExpr) bool {
	switch f := call.Fun.(type) {
	case *ast.Ident:
		return assertingName(f.Name)
	case *ast.SelectorExpr:
		recv, _ := f.X.(*ast.Ident)
		switch {
		case recv == nil:
			return assertingName(f.Sel.Name)
		case plainPackages[recv.Name]:
			return false
		case f.Sel.Name == "Error" || f.Sel.Name == "Errorf":
			return testingReceivers[recv.Name]
		}
		return assertingName(f.Sel.Name)
	}
	return false
}

// testingReceivers are the names a test gives its testing value, so only
// t.Error counts as an assertion and err.Error() does not.
var testingReceivers = map[string]bool{"t": true, "tb": true, "b": true}

// assertingName reports whether a call name opens with an assertion word.
// Must parsers such as mustCompile only build values, so they do not count.
func assertingName(name string) bool {
	lowered := lowerFirst(name)
	return assertionName.MatchString(lowered) && !strings.HasPrefix(lowered, "mustCompile") && !strings.HasPrefix(lowered, "mustParse")
}

// lowerFirst is name with its first letter lowered.
func lowerFirst(name string) string {
	if name == "" {
		return name
	}
	r := []rune(name)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// region is the text of case i of t outside Go. It runs from the case's line
// to the next case, the next test or the end of the file.
func (s *sourceFile) region(t reqTest, from, i int) string {
	lines := strings.Split(s.text, "\n")
	end := len(lines)
	if next := nextBoundary(t, i, lines); next > from {
		end = min(end, next-1)
	}
	if from < 1 || from > end {
		return ""
	}
	return strings.Join(lines[from-1:end], "\n")
}

// nextBoundary is the line where case i of t ends: the next case's, else the
// next test's, or 0 for none.
func nextBoundary(t reqTest, i int, lines []string) int {
	if i+1 < len(t.Cases) {
		_, from := splitLoc(t.Cases[i].Loc)
		if _, n := splitLoc(t.Cases[i+1].Loc); n > from {
			return n
		}
	}
	return nextTest(t, lines)
}

// nextTest is the line of the test that follows t's last case, or 0.
func nextTest(t reqTest, lines []string) int {
	_, at := splitLoc(t.Loc)
	var after []int
	for n := at + 1; n <= len(lines); n++ {
		if testStart.MatchString(lines[n-1]) {
			after = append(after, n)
		}
	}
	sort.Ints(after)
	if len(after) == 0 {
		return 0
	}
	return after[0]
}

// testStart matches a line opening a test declaration outside Go: a JUnit
// annotation, a Swift Testing attribute or an XCTest method.
var testStart = regexp.MustCompile(`^\s*(?:@Test|@TestFactory|func\s+test)`)
