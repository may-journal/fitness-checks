// Command fitness-check-repeated-string-literals fails a string literal that
// appears three or more times across the repo's JS and TS sources. Its rule
// is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
