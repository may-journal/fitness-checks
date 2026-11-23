// Command fitness-check-no-plans-dir keeps plans and temporary design work
// out of the repo. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
