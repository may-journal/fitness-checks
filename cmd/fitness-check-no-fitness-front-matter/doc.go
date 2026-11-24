// Command fitness-check-no-fitness-front-matter fails a Markdown file whose
// front matter still holds `fitnessFunctions` or `relatedConfigurations`. Its
// rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
