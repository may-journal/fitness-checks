// Package marks reads the tags that tie test cases to requirement THENs.
// A tag holds a THEN's ID and the hash of its text, written in the idiom of
// the test's language.
package marks

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Tag is one tag found in a test file.
type Tag struct {
	ID   string
	Hash string
	File string
	Line int
	// Func is the top-level Go function that holds the tag, or "" for another language.
	Func string
}

// Loc is where the tag is written, as file:line.
func (t Tag) Loc() string { return t.File + ":" + strconv.Itoa(t.Line) }

// Value is the tag as swift and java write it, the ID then # then the hash.
func (t Tag) Value() string { return t.ID + "#" + t.Hash }

// idPattern and hashPattern are the shapes of a THEN's ID and hash.
const (
	idPattern   = `\d{4}\.\d+\.\d+\.\d+\.\d+`
	hashPattern = `[0-9a-f]{8}`
)

// textTags match the tags written as text: a JUnit tag, a Swift activity and a Node wrapper call.
var textTags = map[string]*regexp.Regexp{
	".java":  regexp.MustCompile(`@Tag\(\s*"req:(` + idPattern + `)#(` + hashPattern + `)"\s*\)`),
	".swift": regexp.MustCompile(`named:\s*"req (` + idPattern + `)#(` + hashPattern + `)"`),
	".js":    nodeTag,
	".jsx":   nodeTag,
	".mjs":   nodeTag,
	".ts":    nodeTag,
	".tsx":   nodeTag,
}

// nodeTag matches the Node wrapper call then("ID", "hash", ...).
var nodeTag = regexp.MustCompile(`\bthen\(\s*["'](` + idPattern + `)["']\s*,\s*["'](` + hashPattern + `)["']`)

// Tags lists the tags in files, given as paths relative to root. A file it
// cannot read or parse holds none, since the language checks report it.
func Tags(root string, files []string) []Tag {
	var out []Tag
	for _, f := range files {
		raw, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			continue
		}
		if filepath.Ext(f) == ".go" {
			out = append(out, goTags(f, raw)...)
			continue
		}
		out = append(out, textual(f, string(raw))...)
	}
	return out
}

// textual lists the tags the file's language writes as text.
func textual(file, text string) []Tag {
	re := textTags[filepath.Ext(file)]
	if re == nil {
		return nil
	}
	var out []Tag
	for _, m := range re.FindAllStringSubmatchIndex(text, -1) {
		line := 1 + strings.Count(text[:m[0]], "\n")
		out = append(out, Tag{ID: text[m[2]:m[3]], Hash: text[m[4]:m[5]], File: file, Line: line})
	}
	return out
}

// goTags lists the req.Then(t, "ID", "hash", ...) calls of a Go file.
func goTags(file string, src []byte) []Tag {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.SkipObjectResolution)
	if err != nil {
		return nil
	}
	var out []Tag
	for _, d := range f.Decls {
		out = append(out, declTags(fset, file, d)...)
	}
	return out
}

// declTags lists the tags inside one declaration, noting its function name.
func declTags(fset *token.FileSet, file string, d ast.Decl) []Tag {
	name := ""
	if fn, ok := d.(*ast.FuncDecl); ok {
		name = fn.Name.Name
	}
	var out []Tag
	ast.Inspect(d, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if id, hash, ok := ThenArgs(call); ok {
				out = append(out, Tag{ID: id, Hash: hash, File: file, Line: fset.Position(call.Pos()).Line, Func: name})
			}
		}
		return true
	})
	return out
}

// ThenArgs reports the ID and hash of a call pkg.Then(t, "ID", "hash", ...)
// whose two tag arguments are string literals.
func ThenArgs(call *ast.CallExpr) (id, hash string, ok bool) {
	if !isThen(call) {
		return "", "", false
	}
	id, idOK := literal(call.Args[1])
	hash, hashOK := literal(call.Args[2])
	return id, hash, idOK && hashOK && idShape.MatchString(id)
}

// isThen reports whether call is req.Then with at least four arguments.
func isThen(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || len(call.Args) < 4 {
		return false
	}
	pkg, isPkg := sel.X.(*ast.Ident)
	return isPkg && pkg.Name == "req" && sel.Sel.Name == "Then"
}

// idShape matches a THEN ID.
var idShape = regexp.MustCompile(`^` + idPattern + `$`)

// literal is the value of a string literal expression.
func literal(e ast.Expr) (string, bool) {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(lit.Value)
	return s, err == nil
}

// Placed reports whether the tag sits where its THEN's acceptance has its test.
// A Go tag belongs in Test<doc>_<n> of its acceptance, and a tag in another
// language is placed wherever it is.
func (t Tag) Placed() bool {
	parts := strings.SplitN(t.ID, ".", 3)
	return t.Func == "" || len(parts) < 3 || t.Func == "Test"+parts[0]+"_"+parts[1]
}
