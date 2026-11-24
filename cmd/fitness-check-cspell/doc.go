// Command fitness-check-cspell fails unknown words in the checked files. Its
// rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
