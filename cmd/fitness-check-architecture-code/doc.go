// Command fitness-check-architecture-code resolves each class of the code
// view to a declaration a language module's reader reads. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
