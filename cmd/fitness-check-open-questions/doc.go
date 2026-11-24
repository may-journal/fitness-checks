// Command fitness-check-open-questions fails a list of open questions in a
// doc, page, PR or issue body. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
