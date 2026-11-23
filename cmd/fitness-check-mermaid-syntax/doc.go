// Command fitness-check-mermaid-syntax fails a mermaid diagram that would not
// render, so a broken diagram never reaches a reader. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
