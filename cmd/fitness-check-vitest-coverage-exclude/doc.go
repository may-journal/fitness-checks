// Command fitness-check-vitest-coverage-exclude fails any Vitest coverage
// exclude entry, so coverage measures every file. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
