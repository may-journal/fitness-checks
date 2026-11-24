// Command fitness-check-mermaid-callouts requires every mermaid diagram but a
// legend to carry callouts, and a callout table whose callouts match it one
// to one. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
