// cspell:ignore Liau

package main

import _ "embed"

// fileConfig is the check's key in .fitnessrc.json, which --describe lists.
type fileConfig struct {
	// TextReadability holds the alarm bands for this check.
	TextReadability options `json:"textReadability"`
}

// options holds the alarm bands a repo sets under textReadability.
type options struct {
	// MaxGrade is the Coleman-Liau and ARI grade at which a file alarms.
	MaxGrade float64 `json:"maxGrade"`
	// MaxLix is the LIX score at which a file alarms.
	MaxLix float64 `json:"maxLix"`
	// MinWords is the fewest prose words a file needs to be judged.
	MinWords int `json:"minWords"`
}

// Defaults are calibrated to this repo's baseline of grade 11-15 and LIX 40-48, so only genuine outliers fire.
var defaults = fileConfig{TextReadability: options{MaxGrade: 18, MaxLix: 60, MinWords: 100}}

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte
