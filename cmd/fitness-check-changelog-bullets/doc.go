// Command fitness-check-changelog-bullets keeps changelog entries tight:
// every entry file's `###` section must have 3 to 5 bullets, each under 365
// characters. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
