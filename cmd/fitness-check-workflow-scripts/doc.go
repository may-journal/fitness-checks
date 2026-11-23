// Command fitness-check-workflow-scripts fails a `run:` step in a workflow or
// action that holds more than a trivial command (#946). Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
