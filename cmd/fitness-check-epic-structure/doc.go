// Command fitness-check-epic-structure validates that an epic follows the
// Epic issue template, or the repo's own template labeled `Epic`. Its rule is
// in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
