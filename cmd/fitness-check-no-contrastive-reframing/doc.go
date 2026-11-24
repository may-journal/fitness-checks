// Command fitness-check-no-contrastive-reframing flags the "not X, it's Y"
// reframing in markdown prose. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
