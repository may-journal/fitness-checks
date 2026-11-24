package main

import _ "embed"

// fileConfig is the check's key in .fitnessrc.json, which --describe lists.
type fileConfig struct {
	// Requirements holds options for this check.
	Requirements options `json:"requirements"`
}

// options holds what a repo sets under requirements.
type options struct {
	// MaxAcceptances is the most acceptances one requirement doc may hold.
	MaxAcceptances int `json:"maxAcceptances"`
}

// defaults hold the cap when the config sets none.
var defaults = fileConfig{Requirements: options{MaxAcceptances: 8}}

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte
