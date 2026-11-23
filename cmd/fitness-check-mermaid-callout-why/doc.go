// Command fitness-check-mermaid-callout-why requires each numbered callout
// table to carry a "Why" column, and flags numbered rows whose "Why" cell is
// empty. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
