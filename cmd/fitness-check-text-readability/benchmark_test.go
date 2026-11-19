package main

import (
	"path/filepath"
	"testing"
)

// BenchmarkRepositoryScan times a full scan of this repo's Markdown and code comments.
func BenchmarkRepositoryScan(b *testing.B) {
	root, err := filepath.Abs("../../..")
	if err != nil {
		b.Fatal(err)
	}
	var th thresholds
	applyConfig(&th, defaults.TextReadability)
	for b.Loop() {
		walkFiles(root, th, false)
	}
}
