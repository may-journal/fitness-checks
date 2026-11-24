// Command doc-to-readme moves the doc comment of each check folder named on
// its command line from doc.go into README.md, as Markdown. doc.go keeps one
// line that names the README and embeds it, so --describe prints the README's
// first paragraph.
//
//	go run ./scripts/doc-to-readme cmd/fitness-check-<name>...
//
// A doc.go that holds more than the comment, the package clause and its own
// embed is refused. Its folder then stays as it was.
package main

import (
	"fmt"
	"go/doc/comment"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// docTail is what a check's doc.go holds after its comment.
const docTail = "package main\n\nimport _ \"embed\"\n\n" +
	"// docSource is this file, whose comment --describe summarizes.\n//\n//go:embed doc.go\nvar docSource string\n"

// newTail is what doc.go holds after its comment once the README is the docs.
const newTail = "package main\n\nimport _ \"embed\"\n\n" +
	"// docSource is README.md, whose first paragraph --describe summarizes.\n//\n//go:embed README.md\nvar docSource string\n"

// cspellLine matches a comment line that adds words to the spell checker.
var cspellLine = regexp.MustCompile(`(?m)^cspell:ignore (.*)\n?`)

func main() {
	for _, dir := range os.Args[1:] {
		if err := move(dir); err != nil {
			fmt.Fprintf(os.Stderr, "doc-to-readme: %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
}

// move rewrites the README.md and doc.go in dir.
func move(dir string) error {
	src, err := os.ReadFile(filepath.Join(dir, "doc.go"))
	if err != nil {
		return err
	}
	head, ok := strings.CutSuffix(string(src), docTail)
	if !ok {
		return fmt.Errorf("doc.go holds more than its comment, or is already moved")
	}
	text, words, err := commentOf(head)
	if err != nil {
		return err
	}
	name := strings.TrimPrefix(filepath.Base(filepath.Clean(dir)), "fitness-check-")
	readme, err := render(name, text, words)
	if err != nil {
		return err
	}
	return write(dir, readme, pointer(name)+newTail)
}

// write saves the README.md and doc.go contents in dir.
func write(dir, readme, docGo string) error {
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "doc.go"), []byte(docGo), 0o644)
}

// commentOf is the doc comment in the head of doc.go without its spell
// checker lines, and the words those lines add.
func commentOf(head string) (text, words string, err error) {
	file, err := parser.ParseFile(token.NewFileSet(), "doc.go", head+"package main\n", parser.ParseComments)
	if err != nil {
		return "", "", err
	}
	var all []string
	for _, group := range file.Comments {
		all = append(all, group.Text())
	}
	text = strings.Join(all, "\n")
	var found []string
	for _, m := range cspellLine.FindAllStringSubmatch(text, -1) {
		found = append(found, m[1])
	}
	return strings.TrimSpace(cspellLine.ReplaceAllString(text, "")), strings.Join(found, " "), nil
}

// render is the README.md for the check name. It holds the title, the comment
// as Markdown and the spell checker words.
func render(name, text, words string) (string, error) {
	rest, ok := strings.CutPrefix(text, "Command fitness-check-"+name+" ")
	if !ok {
		return "", fmt.Errorf("the comment does not open with %q", "Command fitness-check-"+name+" ")
	}
	first, size := utf8.DecodeRuneInString(rest)
	d := new(comment.Parser).Parse(string(unicode.ToUpper(first)) + rest[size:])
	parts := []string{"# " + name}
	for _, b := range d.Content {
		parts = append(parts, block(b))
	}
	out := strings.Join(parts, "\n\n") + "\n"
	if words != "" {
		out += "\n<!-- cspell:ignore " + words + " -->\n"
	}
	return out, nil
}

// block is one comment block as Markdown.
func block(b comment.Block) string {
	switch b := b.(type) {
	case *comment.Heading:
		return "## " + inline(b.Text)
	case *comment.Code:
		return "```\n" + b.Text + "```"
	case *comment.List:
		return list(b)
	case *comment.Paragraph:
		return inline(b.Text)
	}
	return ""
}

// list is a comment list as Markdown, one item to a line.
func list(l *comment.List) string {
	var lines []string
	for _, item := range l.Items {
		marker := "-"
		if item.Number != "" {
			marker = item.Number + "."
		}
		var paragraphs []string
		for _, b := range item.Content {
			paragraphs = append(paragraphs, block(b))
		}
		lines = append(lines, marker+" "+strings.Join(paragraphs, " "))
	}
	return strings.Join(lines, "\n")
}

// inline is the text of a paragraph on one line, with each link as Markdown.
func inline(texts []comment.Text) string {
	var b strings.Builder
	for _, t := range texts {
		b.WriteString(span(t))
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// span is one piece of paragraph text as Markdown.
func span(t comment.Text) string {
	switch t := t.(type) {
	case comment.Plain:
		return string(t)
	case comment.Italic:
		return string(t)
	case *comment.Link:
		return "[" + inline(t.Text) + "](" + t.URL + ")"
	case *comment.DocLink:
		return "[" + inline(t.Text) + "](" + t.DefaultURL("https://pkg.go.dev") + ")"
	}
	return ""
}

// pointer is the doc.go comment that names the README as the check's docs.
func pointer(name string) string {
	return "// Command fitness-check-" + name + " is documented in README.md.\n"
}
