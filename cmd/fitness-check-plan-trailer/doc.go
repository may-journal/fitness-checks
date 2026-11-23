// Command fitness-check-plan-trailer validates an optional `Plan #NN` commit
// trailer. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
