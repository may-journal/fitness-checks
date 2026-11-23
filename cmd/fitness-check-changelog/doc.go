// Command fitness-check-changelog judges the `changelog/` entry files, the
// only changelog source. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
