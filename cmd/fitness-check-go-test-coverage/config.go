package main

import _ "embed"

// fileConfig is the check's key in .fitnessrc.json, which --describe lists.
type fileConfig struct {
	// GoTestCoverage declares library entry packages for this check.
	GoTestCoverage settings `json:"goTestCoverage"`
}

// settings holds what a repo sets under goTestCoverage.
type settings struct {
	// Entries lists the full import paths of library entry packages,
	// which must reach full statement coverage, as command packages do.
	Entries []string `json:"entries"`
}

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte
