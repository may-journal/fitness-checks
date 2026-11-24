// Command fitness-check-markdown-filename-kebab-case requires every markdown
// file name to be kebab-case, standard root docs exempt. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
