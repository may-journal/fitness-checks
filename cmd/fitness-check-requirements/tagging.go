package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/may-journal/fitness-checks/internal/marks"
	"github.com/may-journal/fitness-checks/internal/packages"
	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-runner/go/pkg/lang"
)

// reqImport is the package that holds the Go helper.
const reqImport = "github.com/may-journal/fitness-checks/req"

// tagPlan is one tag `--tag` writes into a Go case matched by name. The case
// body, after its Parallel calls, runs from offset open to offset end.
type tagPlan struct {
	tag       marks.Tag
	open, end int
	param     string
}

// caseNode is the body of one case and the name of its *testing.T.
type caseNode struct {
	body  *ast.BlockStmt
	param string
}

// parsedTest is a Go test file with its token set.
type parsedTest struct {
	fset *token.FileSet
	file *ast.File
}

// planTags lists the tags of the Go cases of each acceptance, by test file.
func planTags(root string, docs []requirement.Document, tests []reqTest, tagged map[string]bool) (map[string][]tagPlan, error) {
	plans := map[string][]tagPlan{}
	for _, a := range ownedAcceptances(docs, tests) {
		t, ok := goTest(a)
		if !ok {
			continue
		}
		found, err := planTest(root, a, t, tagged)
		if err != nil {
			return nil, err
		}
		plans[t.File] = append(plans[t.File], found...)
	}
	return plans, nil
}

// goTest is the one Go test of a, when it has just that.
func goTest(a *owned) (reqTest, bool) {
	if len(a.tests) != 1 {
		return reqTest{}, false
	}
	return a.tests[0], a.tests[0].module.Name == "go"
}

// parseTest reads the Go file, or returns nil when it does not parse.
func parseTest(root, file string) (*parsedTest, error) {
	src, err := os.ReadFile(filepath.Join(root, file))
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.SkipObjectResolution)
	if err != nil {
		return nil, nil
	}
	return &parsedTest{fset: fset, file: f}, nil
}

// planTest plans the tags of one acceptance's Go test.
func planTest(root string, a *owned, t reqTest, tagged map[string]bool) ([]tagPlan, error) {
	parsed, fn, err := testFunc(root, t)
	if fn == nil || err != nil {
		return nil, err
	}
	steps := stepsByID(a)
	var out []tagPlan
	for id, nodes := range caseNodes(fn, steps, a, testingNames(parsed.file)) {
		if p, ok := planCase(parsed.fset, nodes, steps[id], tagged, t.File); ok {
			out = append(out, p)
		}
	}
	return out, nil
}

// testFunc is the parsed file of t and its function, or none when either is unreadable.
func testFunc(root string, t reqTest) (*parsedTest, *ast.FuncDecl, error) {
	parsed, err := parseTest(root, t.File)
	if parsed == nil || err != nil {
		return nil, nil, err
	}
	fn := findFunc(parsed.file, t.Name)
	if fn == nil || fn.Body == nil {
		return nil, nil, nil
	}
	return parsed, fn, nil
}

// stepsByID maps each step of a to its ID. A single step has the ID "".
func stepsByID(a *owned) map[string]requirement.Step {
	steps := map[string]requirement.Step{}
	for _, s := range a.Steps {
		steps[s.ID] = s
	}
	return steps
}

// caseNodes maps each step ID to the cases that name it. A single step has
// the whole test as its case.
func caseNodes(fn *ast.FuncDecl, steps map[string]requirement.Step, a *owned, names map[string]bool) map[string][]caseNode {
	if _, single := steps[""]; single && len(a.Steps) == 1 {
		return map[string][]caseNode{"": {{body: fn.Body, param: testParam(fn.Type)}}}
	}
	out := map[string][]caseNode{}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if id, c, ok := namedCase(n, steps, names); ok {
			out[id] = append(out[id], c)
		}
		return true
	})
	return out
}

// namedCase reads a t.Run call whose literal name opens with a step ID of steps.
func namedCase(n ast.Node, steps map[string]requirement.Step, names map[string]bool) (string, caseNode, bool) {
	call, ok := n.(*ast.CallExpr)
	if !ok || !isRun(call, names) {
		return "", caseNode{}, false
	}
	name, body, ok := runCase(call)
	id := lang.NewCase(name, "", 0).Step
	_, held := steps[id]
	return id, body, ok && id != "" && held
}

// runCase is the name and body of the t.Run call, when its case is a function literal.
func runCase(call *ast.CallExpr) (string, caseNode, bool) {
	name, err := strconv.Unquote(call.Args[0].(*ast.BasicLit).Value)
	fn, isFunc := call.Args[1].(*ast.FuncLit)
	if err != nil || !isFunc {
		return "", caseNode{}, false
	}
	return name, caseNode{body: fn.Body, param: testParam(fn.Type)}, true
}

// findFunc is the top-level function named name.
func findFunc(f *ast.File, name string) *ast.FuncDecl {
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}

// testParam is the name of the first parameter of ft.
func testParam(ft *ast.FuncType) string {
	if ft.Params == nil || len(ft.Params.List) == 0 || len(ft.Params.List[0].Names) == 0 {
		return ""
	}
	return ft.Params.List[0].Names[0].Name
}

// planCase plans the tag of the one case of a step with one THEN that no tag
// names yet.
func planCase(fset *token.FileSet, nodes []caseNode, s requirement.Step, tagged map[string]bool, file string) (tagPlan, bool) {
	if len(nodes) != 1 || len(s.Thens) != 1 || tagged[s.Thens[0].ID] {
		return tagPlan{}, false
	}
	return planBody(fset, nodes[0], s.Thens[0], file)
}

// planBody plans the statements of the case after its Parallel calls, which
// move into a req.Then. A body with none, or a call of req.Then already, is
// left alone.
func planBody(fset *token.FileSet, c caseNode, th requirement.Then, file string) (tagPlan, bool) {
	stmts := afterParallel(c.body.List)
	if len(stmts) == 0 || !usable(c) {
		return tagPlan{}, false
	}
	start := fset.Position(stmts[0].Pos())
	return tagPlan{
		tag:   marks.Tag{ID: th.ID, Hash: th.Hash(), File: file, Line: start.Line},
		open:  start.Offset,
		end:   fset.Position(stmts[len(stmts)-1].End()).Offset,
		param: c.param,
	}, true
}

// afterParallel is stmts without its leading calls of Parallel.
func afterParallel(stmts []ast.Stmt) []ast.Stmt {
	for len(stmts) > 0 && isParallel(stmts[0]) {
		stmts = stmts[1:]
	}
	return stmts
}

// usable reports whether the case names its *testing.T and its body is fit to wrap.
func usable(c caseNode) bool {
	return c.param != "" && c.param != "_" && fit(c)
}

// fit reports whether a req.Then would count the body's assertions and compile.
func fit(c caseNode) bool {
	return callsCounted(c.body) && !thenInside(c.body) && !rawAssertion(c) && !usesName(c.body, "req")
}

// countedHelpers are the acceptance helpers that call req.Asserted, so a req.Then counts their assertion.
var countedHelpers = []string{"Sees", "PassedWith", "PassedWithNoFiles", "ShowsInfo", "JavaSeesSnippet", "JavaCalledOnce"}

// callsCounted reports whether body calls one of the acceptance helpers that count an assertion.
func callsCounted(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok && calls(call, "acceptance") {
			found = found || slices.Contains(countedHelpers, call.Fun.(*ast.SelectorExpr).Sel.Name)
		}
		return !found
	})
	return found
}

// usesName reports whether body mentions an identifier called name.
func usesName(body *ast.BlockStmt, name string) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name == name {
			found = true
		}
		return !found
	})
	return found
}

// rawAssertion reports whether the case calls Error, Fatal or Fail on its test.
func rawAssertion(c caseNode) bool {
	found := false
	ast.Inspect(c.body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			found = found || failsTest(call, c.param)
		}
		return !found
	})
	return found
}

// failsTest reports whether call is param.Error, Errorf, Fatal, Fatalf, Fail or FailNow.
func failsTest(call *ast.CallExpr, param string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || !slices.ContainsFunc([]string{"Fail", "Error", "Fatal"}, func(p string) bool { return strings.HasPrefix(sel.Sel.Name, p) }) {
		return false
	}
	id, isID := sel.X.(*ast.Ident)
	return isID && id.Name == param
}

// isParallel reports whether s is a call of Parallel.
func isParallel(s ast.Stmt) bool {
	e, ok := s.(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := e.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "Parallel"
}

// thenInside reports whether body already calls a req.Then.
func thenInside(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if _, _, ok := marks.ThenArgs(call); ok {
				found = true
			}
		}
		return !found
	})
	return found
}

// writeTags applies each file's plans unless dry, and returns every tag planned.
func writeTags(root string, plans map[string][]tagPlan, dry bool) ([]marks.Tag, error) {
	var done []marks.Tag
	for file, list := range plans {
		for _, p := range list {
			done = append(done, p.tag)
		}
		if err := rewriteFile(root, file, list, dry); err != nil {
			return nil, err
		}
	}
	return done, nil
}

// rewriteFile wraps each planned case body of the file, unless dry or none.
func rewriteFile(root, file string, plans []tagPlan, dry bool) error {
	if dry || len(plans) == 0 {
		return nil
	}
	path := filepath.Join(root, file)
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	out, err := wrapped(src, plans)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return writeLike(path, out)
}

// writeLike writes out over path with the mode path has.
func writeLike(path string, out []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, info.Mode().Perm())
}

// wrapped is src with each planned body in a req.Then, imported and formatted.
func wrapped(src []byte, plans []tagPlan) ([]byte, error) {
	sort.Slice(plans, func(i, j int) bool { return plans[i].open > plans[j].open })
	for _, p := range plans {
		head := fmt.Sprintf("req.Then(%s, %q, %q, func(_ req.A) {\n", p.param, p.tag.ID, p.tag.Hash)
		src = slices.Concat(src[:p.open], []byte(head), src[p.open:p.end], []byte("\n})"), src[p.end:])
	}
	src, err := withImport(src)
	if err != nil {
		return nil, err
	}
	return format.Source(src)
}

// sameGroup matches an import of another may-journal package.
var sameGroup = regexp.MustCompile(`(?m)^\t"github.com/may-journal/[^\n]*\n`)

// withImport adds the req import to src unless it holds it.
func withImport(src []byte) ([]byte, error) {
	if bytes.Contains(src, []byte(strconv.Quote(reqImport))) {
		return src, nil
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	if gen := importDecl(f); gen != nil {
		return importedIn(src, fset, gen, "\t"+strconv.Quote(reqImport)+"\n"), nil
	}
	at := fset.Position(f.Name.End()).Offset
	return slices.Concat(src[:at], []byte("\n\nimport "+strconv.Quote(reqImport)), src[at:]), nil
}

// importDecl is the first import declaration of f, or nil.
func importDecl(f *ast.File) *ast.GenDecl {
	for _, d := range f.Decls {
		if gen, ok := d.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
			return gen
		}
	}
	return nil
}

// importedIn adds spec to the import declaration gen of src, beside the
// imports of other may-journal packages when it holds any.
func importedIn(src []byte, fset *token.FileSet, gen *ast.GenDecl, spec string) []byte {
	if !gen.Lparen.IsValid() {
		from, to := fset.Position(gen.Pos()).Offset, fset.Position(gen.End()).Offset
		block := "import (\n" + strings.TrimPrefix(string(src[from:to]), "import ") + "\n" + spec + ")"
		return slices.Concat(src[:from], []byte(block), src[to:])
	}
	at := fset.Position(gen.Rparen).Offset
	insert, text := at, "\n"+spec
	if m := sameGroup.FindAllIndex(src[:at], -1); len(m) > 0 {
		insert, text = m[len(m)-1][1], spec
	}
	return slices.Concat(src[:insert], []byte(text), src[insert:])
}

// inRoot keeps the items whose file lies in package root r among roots.
func inRoot[T any](roots []string, r string, items []T, file func(T) string) []T {
	var out []T
	for _, it := range items {
		if packages.Of(roots, file(it)) == r {
			out = append(out, it)
		}
	}
	return out
}

// testFile and tagFile are the files of a test and of a tag.
func testFile(t reqTest) string  { return t.File }
func tagFile(t marks.Tag) string { return t.File }

// tagMode plans, and unless dry writes, the tags of each package root. It
// returns the tags it planned.
func tagMode(root string, roots, paths []string, tests []reqTest, tags []marks.Tag, dry bool) ([]marks.Tag, error) {
	var all []marks.Tag
	for _, r := range roots {
		docs, _, _ := requirement.Load(root, packages.Owned(roots, r, paths))
		own := inRoot(roots, r, tags, tagFile)
		plans, err := planTags(root, docs, inRoot(roots, r, tests, testFile), taggedThens(own))
		if err != nil {
			return nil, err
		}
		done, err := writeTags(root, plans, dry)
		if err != nil {
			return nil, err
		}
		all = append(all, done...)
	}
	return all, nil
}

// taggedRun reads the tags of the test files. With `--tag` it first writes
// the tags it can, or counts them as written with `--dry-run`.
func taggedRun(root string, roots, paths []string, tests []reqTest, files, args []string) ([]marks.Tag, error) {
	tags := marks.Tags(root, files)
	if !slices.Contains(args, tagFlag) {
		return tags, nil
	}
	wrote, err := tagMode(root, roots, paths, tests, tags, slices.Contains(args, dryFlag))
	if slices.Contains(args, dryFlag) {
		return append(tags, wrote...), err
	}
	return marks.Tags(root, files), err
}
