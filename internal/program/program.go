// Package program runs one check as its own program. It answers
// --describe with the check's summary, its keys and the protocol number.
// It hands every other run to checkkit.
package program

import (
	"cmp"
	"encoding/json"
	"fmt"
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
	Configured(source, nil, c)
}

// Configured runs c as Main does, and its --describe also lists the
// .fitnessrc.json keys that keys returns. Only --describe calls keys, so
// other runs skip the work.
func Configured(source string, keys func() ([]checkkit.ConfigKey, error), c checkkit.Check) {
	if slices.Contains(os.Args[1:], "--describe") {
		os.Exit(describe(os.Stdout, source, keys, c.Describe))
	}
	checkkit.Main(c)
}

// describe prints d with its summary, keys and the protocol, and returns
// the exit code.
func describe(out io.Writer, source string, keys func() ([]checkkit.ConfigKey, error), d checkkit.Describe) int {
	d.Summary = cmp.Or(d.Summary, Summary(source))
	if keys != nil {
		config, err := keys()
		if err != nil {
			fmt.Fprintln(os.Stderr, "--describe:", err)
			return checkkit.ExitInternal
		}
		d.Config = config
	}
	if err := json.NewEncoder(out).Encode(described{d, Protocol}); err != nil {
		return checkkit.ExitInternal
	}
	return 0
}

// Summary is the first sentence of a check's documentation as plain text.
// The source is a README.md, which opens with a title, or a doc.go, whose
// doc comment opens with "Command fitness-check-<name>".
func Summary(source string) string {
	if strings.HasPrefix(strings.TrimSpace(source), "#") {
		return markdownSummary(source)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "doc.go", source, parser.PackageClauseOnly|parser.ParseComments)
	if err != nil || file.Doc == nil {
		return ""
	}
	text := afterCommand(file.Doc.Text())
	return unlinked(new(doc.Package).Synopsis(text), text)
}

// markdownSummary is the first sentence of the paragraph after the title of
// the Markdown source. A source with no such paragraph has no summary.
func markdownSummary(source string) string {
	source = strings.ReplaceAll(source, "\r\n", "\n")
	_, body, _ := strings.Cut(strings.TrimSpace(source), "\n")
	paragraph, _, _ := strings.Cut(strings.TrimSpace(body), "\n\n")
	if strings.HasPrefix(paragraph, "#") {
		return ""
	}
	return new(doc.Package).Synopsis(paragraph)
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
