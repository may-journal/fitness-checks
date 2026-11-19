package waiting

import (
	_ "embed"

	"github.com/may-journal/fitness-checks/internal/config"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

// fileConfig is adr-as-issue's key in .fitnessrc.json, which its
// --describe lists.
type fileConfig struct {
	// AdrAsIssue holds the adr-as-issue check's options.
	AdrAsIssue options `json:"adrAsIssue"`
}

// options holds what a repo sets under adrAsIssue.
type options struct {
	// Allow lists ADR files still waiting to move, each with its reason,
	// as a shrinking fallback list.
	Allow []struct {
		// Path is a glob, as Go's path.Match reads it, of the waiting files.
		Path string `json:"path"`
		// Reason says why the files have not moved yet.
		Reason string `json:"reason"`
	} `json:"allow"`
}

// configSource is this file, whose doc comments describe each key.
//
//go:embed config.go
var configSource []byte

// Keys lists adrAsIssue and its keys for adr-as-issue's --describe. The
// engine fails two checks that declare one key, so adr-checks, which reads
// the list too, declares none.
func Keys() func() ([]checkkit.ConfigKey, error) {
	return config.Keys(fileConfig{}, configSource)
}
