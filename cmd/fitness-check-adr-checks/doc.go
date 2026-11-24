// Command fitness-check-adr-checks pairs each architecture decision, an `ADR`
// issue, with the check in `fitness/adr-NNNN/` that proves it. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
