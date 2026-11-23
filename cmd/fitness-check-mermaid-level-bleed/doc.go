// Command fitness-check-mermaid-level-bleed warns when a callout description
// at one C4 level repeats the previous level verbatim. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
