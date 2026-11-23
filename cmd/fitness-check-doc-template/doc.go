// Command fitness-check-doc-template makes each markdown document match the
// shape of a template that governs it. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
