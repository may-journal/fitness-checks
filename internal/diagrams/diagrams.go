// Package diagrams holds the code the mermaid doc checks share. It gives
// their Describe and their run over a pasted description or every document.
// It also exports the JavaScript whitespace class pkg/mermaid keeps to itself.
package diagrams

import (
	"strings"

	"github.com/may-journal/fitness-checks/internal/match"
	"github.com/may-journal/fitness-runner/go/pkg/bodycheck"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

// Validate judges one document's content, naming it file in each error.
type Validate func(file, content string) []string

// Driver walks root's documents through validate, as mermaid.RunDocCheck
// and mermaid.RunDiagramCheck do.
type Driver func(root string, validate func(file, content string) []string) ([]string, int, error)

// Describe is a blob check named name that reads the files ending in exts.
func Describe(name string, exts []string) checkkit.Describe {
	return checkkit.Describe{
		Name:  name,
		Scope: checkkit.ScopeFiles,
		Blob:  true,
		Files: match.Exts(exts...),
	}
}

// Check is the check named name: it judges a pasted description with
// validate, and otherwise every document drive walks.
func Check(name string, exts []string, drive Driver, validate Validate) checkkit.Check {
	return checkkit.Check{
		Describe: Describe(name, exts),
		Run: func(root string, args []string) (checkkit.Result, error) {
			if res, handled, err := Description(root, args, validate); err != nil || handled {
				return res, err
			}
			errors, filesChecked, err := drive(root, validate)
			if err != nil {
				return checkkit.Result{}, err
			}
			return Verdict(filesChecked, errors), nil
		},
	}
}

// Description judges a pasted description with validate, reporting whether
// the run held one.
func Description(root string, args []string, validate Validate) (checkkit.Result, bool, error) {
	res, handled, err := bodycheck.RunDoc(root, args, func(_, content string) []string {
		return validate("(description)", content)
	})
	if err != nil {
		return checkkit.Result{}, false, err
	}
	return res, handled, nil
}

// Verdict fails filesChecked files with errors, or passes them when there
// are none.
func Verdict(filesChecked int, errors []string) checkkit.Result {
	if len(errors) > 0 {
		return checkkit.Fail(filesChecked, errors...)
	}
	return checkkit.Pass(filesChecked)
}

// JSWhitespace is the JavaScript \s class as an RE2 class body. It holds
// ASCII whitespace with \v, the Unicode spaces, line separators and the BOM.
const JSWhitespace = `\t\n\v\f\r \x{00A0}\x{1680}\x{2000}-\x{200A}\x{2028}\x{2029}\x{202F}\x{205F}\x{3000}\x{FEFF}`

// IsJSSpace reports whether r is in the JavaScript \s class, the rune-level
// twin of JSWhitespace.
func IsJSSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ',
		0x00A0, 0x1680, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000, 0xFEFF:
		return true
	}
	return r >= 0x2000 && r <= 0x200A
}

// TrimJS trims JavaScript whitespace from both ends, as String.prototype.trim.
func TrimJS(s string) string { return strings.TrimFunc(s, IsJSSpace) }

// TrimStartJS trims leading JavaScript whitespace, as
// String.prototype.trimStart.
func TrimStartJS(s string) string { return strings.TrimLeftFunc(s, IsJSSpace) }
