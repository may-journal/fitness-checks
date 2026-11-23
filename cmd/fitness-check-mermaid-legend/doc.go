// Command fitness-check-mermaid-legend requires numbered flowchart callout
// nodes to carry a style class, and diagrams with callouts to declare a
// `classDef` legend. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
