// Command fitness-check-workflow-ref-exists fails a workflow calling a
// may-journal workflow or action at a tag or sha lacking that path (#485).
// Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
