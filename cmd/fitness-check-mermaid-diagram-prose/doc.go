// Command fitness-check-mermaid-diagram-prose fails edge labels that carry
// prose beyond a callout, number or letter, in any mermaid diagram. Its rule
// is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
