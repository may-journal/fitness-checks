// Command fitness-check-markdown-links requires every relative link in every
// markdown file to resolve to an existing file or directory. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
