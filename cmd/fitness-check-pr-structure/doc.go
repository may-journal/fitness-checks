// Command fitness-check-pr-structure validates that a pull request
// description follows the PR template. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
