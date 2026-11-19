package main

import _ "embed"

// fileConfig is the check's key in .fitnessrc.json, which --describe lists.
type fileConfig struct {
	// RepeatedStringLiterals holds options for this check.
	RepeatedStringLiterals options `json:"repeatedStringLiterals"`
}

// options holds what a repo sets under repeatedStringLiterals.
type options struct {
	// Allow lists exact values the check never flags, for the rare repeat
	// no shared constant can fix. Keep it a shrinking list.
	Allow []string `json:"allow"`
}

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte
