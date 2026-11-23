// Command fitness-check-go-vet runs `go vet ./...` in every Go module of the
// repo, catching suspicious code the compiler accepts. Its rule is in
// README.md.
package main

import _ "embed"

// docSource is README.md, whose first paragraph --describe summarizes.
//
//go:embed README.md
var docSource string
