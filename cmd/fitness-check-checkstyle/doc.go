// Command fitness-check-checkstyle lints Java with the Checkstyle plugin the
// repo's own build declares, and holds its config to a complexity ceiling.
// Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
