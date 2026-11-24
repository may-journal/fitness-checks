// Command fitness-check-semantic-commit requires the commit subject to follow
// Conventional Commits, `type(scope): description`, with the scope required.
// Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
