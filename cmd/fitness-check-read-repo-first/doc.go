// Command fitness-check-read-repo-first reminds agents and users to read the
// repo's decisions and enabled checks first. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
