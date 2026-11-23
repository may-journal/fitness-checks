// Command fitness-check-adr-as-issue keeps architecture decisions in `ADR`
// issues, never in files. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
