// Command fitness-check-mermaid-diagram-table-gap flags loose markdown prose
// wedged between a numbered mermaid diagram (or its legend block) and the
// callout table that follows it. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
