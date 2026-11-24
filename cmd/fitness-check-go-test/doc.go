// Command fitness-check-go-test runs `go test ./...` in every Go module of
// the repo and reports the failing tests. Its rule is in README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
