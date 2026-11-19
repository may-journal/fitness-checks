package main

import _ "embed"

// fileConfig is the check's key in .fitnessrc.json, which --describe lists.
type fileConfig struct {
	// GoComplexity holds options for this check.
	GoComplexity options `json:"goComplexity"`
}

// options holds what a repo sets under goComplexity.
type options struct {
	// Max is the highest cyclomatic complexity one Go function may have.
	Max int `json:"max"`
}

// defaults hold the ceiling when the config sets none.
var defaults = fileConfig{GoComplexity: options{Max: 5}}

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte
