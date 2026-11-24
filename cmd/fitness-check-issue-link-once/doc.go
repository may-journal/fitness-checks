// Command fitness-check-issue-link-once keeps each Issue's timeline readable
// (#153). Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
