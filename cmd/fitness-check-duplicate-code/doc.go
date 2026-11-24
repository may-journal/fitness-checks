// Command fitness-check-duplicate-code fails a repo whose duplicated lines
// pass 1% of its scanned lines. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
