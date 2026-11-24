// Command fitness-check-spotless checks Java formatting with the Spotless
// plugin the repo's own build declares. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
