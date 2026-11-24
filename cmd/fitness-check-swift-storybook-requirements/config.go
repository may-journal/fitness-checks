package main

import _ "embed"

// fileConfig is the check's key in .fitnessrc.json, which --describe lists.
type fileConfig struct {
	// SwiftStorybookRequirements holds options for this check.
	SwiftStorybookRequirements options `json:"swiftStorybookRequirements"`
}

// options holds what a repo sets under swiftStorybookRequirements.
type options struct {
	// MaxPagesPerAcceptance is the most Storybook pages one acceptance ID may
	// name. A positive value replaces the default; zero keeps it.
	MaxPagesPerAcceptance int `json:"maxPagesPerAcceptance"`
}

// defaults give each acceptance one page when the config sets none.
var defaults = fileConfig{SwiftStorybookRequirements: options{MaxPagesPerAcceptance: 1}}

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte
