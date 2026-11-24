// Command fitness-check-dependency-currency fails a declared npm dependency
// behind its latest published version — an installable update means the
// project isn't at peak fitness. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
