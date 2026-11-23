// Command fitness-check-comment-structure judges one GitHub issue or PR
// comment, so a reply stays short, readable, and free of agent noise. Its
// rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
