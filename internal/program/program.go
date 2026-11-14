// Package program runs one check as its own program. It answers
// --describe with the check's summary, taken from its doc comment, and the
// protocol number, then hands every other run to checkkit.
package program

import (
	"cmp"
	"encoding/json"
	"go/doc"
	"go/doc/comment"
	"go/parser"
	"go/token"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

// Protocol is the check protocol version these programs speak. The engine
// fails a run whose check reports another one.
const Protocol = 1

// described is what --describe prints: the check's metadata and protocol.
type described struct {
	checkkit.Describe
	Protocol int `json:"protocol"`
}

// Main runs c, given its doc.go source.
func Main(source string, c checkkit.Check) {
	if slices.Contains(os.Args[1:], "--describe") {
		os.Exit(describe(os.Stdout, source, c.Describe))
	}
	checkkit.Main(c)
}

// describe prints d with its summary and the protocol, and returns the exit
// code.
func describe(out io.Writer, source string, d checkkit.Describe) int {
	d.Summary = cmp.Or(d.Summary, Summary(source))
	if err := json.NewEncoder(out).Encode(described{d, Protocol}); err != nil {
		return checkkit.ExitInternal
	}
	return 0
}

// Summary is the first sentence of the doc comment in source, which opens
// with "Command fitness-check-<name>", as plain text.
func Summary(source string) string {
	file, err := parser.ParseFile(token.NewFileSet(), "doc.go", source, parser.PackageClauseOnly|parser.ParseComments)
	if err != nil || file.Doc == nil {
		return ""
	}
	text := afterCommand(file.Doc.Text())
	return unlinked(new(doc.Package).Synopsis(text), text)
}

// afterCommand drops the "Command fitness-check-<name>" opening from text and
// capitalizes what follows, so the sentence reads on its own.
func afterCommand(text string) string {
	rest, ok := strings.CutPrefix(text, "Command fitness-check-")
	if !ok {
		return text
	}
	if _, rest, ok = strings.Cut(rest, " "); !ok || rest == "" {
		return text
	}
	return strings.ToUpper(rest[:1]) + rest[1:]
}

// unlinked drops the brackets Go keeps around a link in summary when it cuts
// the sentence before text defines the link.
func unlinked(summary, text string) string {
	for _, def := range new(comment.Parser).Parse(text).Links {
		summary = strings.ReplaceAll(summary, "["+def.Text+"]", def.Text)
	}
	return summary
}
