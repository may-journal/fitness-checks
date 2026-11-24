// Command fitness-check-workflow-reporting fails a workflow job that can fail
// without a report on its run page (#181). Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
