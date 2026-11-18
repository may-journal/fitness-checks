// Package match builds the RE2 patterns of a check's Files, Exclude and
// Needs fields. It takes the constants the check reads files by, so no check
// hand-types an extension list.
package match

import (
	"regexp"
	"strings"
)

// Fold makes pattern match without regard to case, for a check that
// compares paths lowercased.
func Fold(pattern string) string {
	return "(?i)" + pattern
}

// Exact matches only the given repo-relative paths, whole.
func Exact(paths ...string) string {
	return "^(" + alternatives(paths, nil) + ")$"
}

// Exts matches a repo-relative path that ends in one of exts, each given
// with its leading dot as the language modules spell them.
func Exts(exts ...string) string {
	return "\\.(" + alternatives(exts, func(e string) string { return strings.TrimPrefix(e, ".") }) + ")$"
}

// VirtualSince is the oldest fitness release whose runs hold the virtual
// paths checkkit names, so a check matching only those needs it.
const VirtualSince = "1.75.0"

// NoExt matches any file name with no dot, such as LICENSE or a shell
// script the comments package knows by its shebang.
func NoExt() string {
	return "(^|/)[^/.]+$"
}

// Base matches a path whose file name is one of names, in any folder.
func Base(names ...string) string {
	return "(^|/)(" + alternatives(names, nil) + ")$"
}

// Dir matches a path under a folder named one of dirs, at any depth.
func Dir(dirs ...string) string {
	return "(^|/)(" + alternatives(dirs, nil) + ")/"
}

// Under matches a path that starts with prefix, a folder or file at the
// repo root such as `.github/workflows/`.
func Under(prefix string) string {
	return "^" + regexp.QuoteMeta(prefix)
}

// Any matches a path any of patterns matches, leaving out empty ones.
func Any(patterns ...string) string {
	var kept []string
	for _, p := range patterns {
		if p != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, "|")
}

// alternatives quotes each value, after trim when set, and joins them with |.
func alternatives(values []string, trim func(string) string) string {
	quoted := make([]string, len(values))
	for i, v := range values {
		if trim != nil {
			v = trim(v)
		}
		quoted[i] = regexp.QuoteMeta(v)
	}
	return strings.Join(quoted, "|")
}
