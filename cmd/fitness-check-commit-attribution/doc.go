// Command fitness-check-commit-attribution requires a commit message to
// disclose AI use in "AI-Tools:" and "AI-Models:" git trailers. Its rule is
// in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
