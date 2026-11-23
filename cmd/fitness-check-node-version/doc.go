// Command fitness-check-node-version requires the machine's active Node
// version to satisfy the repo's `.nvmrc`. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
