// Command fitness-check-build-output-untracked keeps build output out of git
// and out of source imports. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
