// Command fitness-check-markdown-no-bold-italic fails bold and italic
// emphasis in markdown files: `**bold**`, `__bold__`, `*italic*`, or
// `_italic_`. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
