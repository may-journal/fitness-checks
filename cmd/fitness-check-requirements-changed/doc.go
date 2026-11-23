// Command fitness-check-requirements-changed fails a change that adds or
// changes product code without also changing a test and a requirement doc.
// Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
