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
	th := thresholds{maxGrade: defaults.MaxGrade, maxLix: defaults.MaxLix, minWords: defaults.MinWords}
	for b.Loop() {
		walkFiles(root, th, false)
	}
}
