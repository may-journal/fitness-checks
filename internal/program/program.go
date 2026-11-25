// Package program runs one check as its own program. It answers
// --describe with the check's summary, its keys and the protocol number.
// It hands every other run to checkkit.
package program

import (
	"cmp"
	"encoding/json"
	"fmt"
	"go/doc"
	"io"
	"os"
	"regexp"
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

// Main runs c, given its README.md source.
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

// Summary is the first sentence of the paragraph after the title of a check's
// README.md. A first block that is not prose, such as
// a list, fence, quote, table or HTML, gives no summary.
func Summary(source string) string {
	source = strings.ReplaceAll(source, "\r\n", "\n")
	_, body, _ := strings.Cut(strings.TrimSpace(source), "\n")
	paragraph, _, _ := strings.Cut(strings.TrimSpace(body), "\n\n")
	if !markdownProse.MatchString(paragraph) {
		return ""
	}
	return new(doc.Package).Synopsis(markdownLink.ReplaceAllString(paragraph, "$1"))
}

// markdownProse matches a block that starts as plain prose: not a heading,
// list item, code fence, quote, table, HTML or indented code.
var markdownProse = regexp.MustCompile(`^[^\s#>|<` + "`~" + `*+\-\d]|^\d+[^.)\d]|^[*+\-][^\s]`)

// markdownLink matches an inline Markdown link, whose text is group 1. Its URL
// may hold one level of balanced parentheses.
var markdownLink = regexp.MustCompile(`\[([^\]]*)\]\((?:[^()]|\([^()]*\))*\)`)
