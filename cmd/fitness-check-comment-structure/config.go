package main

import _ "embed"

// fileConfig is the check's key in .fitnessrc.json, which --describe lists.
type fileConfig struct {
	// Comment holds this check's budget for an issue or PR comment. Each
	// limit falls back to its default when zero.
	Comment options `json:"comment"`
}

// options holds the budget a repo sets under comment for an issue or PR
// comment. Each limit falls back to its default when zero.
type options struct {
	// MaxWords is the most words one comment may hold.
	MaxWords int `json:"maxWords"`
	// MaxTables is the most tables one comment may hold.
	MaxTables int `json:"maxTables"`
}

// defaults are the limits .fitnessrc.json's comment key may replace.
var defaults = fileConfig{Comment: options{MaxWords: 100, MaxTables: 1}}

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte
