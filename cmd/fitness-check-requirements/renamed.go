package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/may-journal/fitness-runner/go/pkg/gitx"
	"github.com/may-journal/fitness-runner/go/pkg/lang"
)

// withoutTagEdits drops each changed Go test file that differs from base in case names, req tags or req wrappers alone.
// A tag ties a case to its THEN, so such an edit needs no doc change.
func (p pairs) withoutTagEdits(root, base string, changed []string) ([]string, error) {
	files := p.goTests(changed)
	if base == "" || len(files) == 0 {
		return changed, nil
	}
	var objects []string
	for _, f := range files {
		objects = append(objects, base+":"+f, gitx.JudgedObject(f))
	}
	text, err := gitx.ShowAll(root, objects)
	same := map[string]bool{}
	for _, f := range files {
		same[f] = tagEditOnly(f, text[base+":"+f], text[gitx.JudgedObject(f)])
	}
	return slices.DeleteFunc(slices.Clone(changed), func(f string) bool { return same[f] }), err
}

// goTests keeps the changed Go test files that prove a requirement.
func (p pairs) goTests(changed []string) []string {
	var out []string
	for _, f := range p.owned(changed) {
		if strings.HasSuffix(f, "_test.go") && len(p.numbers[f]) > 0 {
			out = append(out, f)
		}
	}
	return out
}

// tagEditOnly reports whether two versions of a Go file differ in case names, tags and wrappers alone.
func tagEditOnly(file, was, now string) bool {
	unwrapped := untagged(file, now)
	return was != "" && unwrapped != "" && untagged(file, was) == unwrapped
}

// untagged is a Go file's code without case names, req tags, wrappers and imports, or "" if it fails to parse.
func untagged(file, src string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path.Base(file), src, parser.SkipObjectResolution)
	if err != nil {
		return ""
	}
	f.Comments = nil
	names := testingNames(f)
	ast.Inspect(f, func(n ast.Node) bool {
		blankName(n, names)
		unwrapIn(n)
		return true
	})
	dropImports(f)
	var out bytes.Buffer
	if printer.Fprint(&out, fset, f) != nil {
		return ""
	}
	clean, err := format.Source(out.Bytes())
	if err != nil {
		return ""
	}
	return withoutBlankLines(string(clean))
}

// blankName cuts the name of a t.Run case to its step ID.
func blankName(n ast.Node, names map[string]bool) {
	if call, ok := n.(*ast.CallExpr); ok && isRun(call, names) {
		lit := call.Args[0].(*ast.BasicLit)
		name, _ := strconv.Unquote(lit.Value)
		lit.Value = strconv.Quote(lang.NewCase(name, "", 0).Step)
	}
}

// unwrapIn replaces each req.Then statement of the statement list n holds by its body.
func unwrapIn(n ast.Node) {
	switch x := n.(type) {
	case *ast.BlockStmt:
		x.List = flatten(x.List)
	case *ast.CaseClause:
		x.Body = flatten(x.Body)
	}
}

// flatten is stmts with each req.Then or req.Asserted statement replaced by
// the statements of its body, or dropped.
func flatten(stmts []ast.Stmt) []ast.Stmt {
	var out []ast.Stmt
	for _, s := range stmts {
		call := reqCall(s)
		switch {
		case call == nil:
			out = append(out, s)
		case thenBody(call) != nil:
			out = append(out, flatten(thenBody(call))...)
		}
	}
	return out
}

// thenBody is the statements of the function literal a req.Then call ends with, or nil.
func thenBody(call *ast.CallExpr) []ast.Stmt {
	sel := call.Fun.(*ast.SelectorExpr)
	if sel.Sel.Name != "Then" || len(call.Args) < 4 {
		return nil
	}
	if fn, ok := call.Args[len(call.Args)-1].(*ast.FuncLit); ok {
		return fn.Body.List
	}
	return nil
}

// reqCall is the call of a req package function that s holds, or nil.
func reqCall(s ast.Stmt) *ast.CallExpr {
	if e, ok := s.(*ast.ExprStmt); ok {
		if call, ok := e.X.(*ast.CallExpr); ok && calls(call, "req") {
			return call
		}
	}
	return nil
}

// calls reports whether call is pkg.Name(...) for some name.
func calls(call *ast.CallExpr, pkg string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == pkg
}

// withoutBlankLines is text without its empty lines, which a wrapper adds around a body.
func withoutBlankLines(text string) string {
	lines := strings.Split(text, "\n")
	return strings.Join(slices.DeleteFunc(lines, func(l string) bool { return strings.TrimSpace(l) == "" }), "\n")
}

// dropImports removes the import declarations of f, which a tag edit changes.
func dropImports(f *ast.File) {
	f.Decls = slices.DeleteFunc(f.Decls, func(d ast.Decl) bool {
		gen, ok := d.(*ast.GenDecl)
		return ok && gen.Tok == token.IMPORT
	})
}
