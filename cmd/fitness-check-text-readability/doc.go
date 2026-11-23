// Command fitness-check-text-readability fails a markdown file whose prose
// scores as hard to read. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
