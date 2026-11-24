package main

import (
	"fmt"

	"github.com/may-journal/fitness-checks/internal/config"
	"github.com/may-journal/fitness-checks/internal/requirement"
)

// maxAcceptances resolves the cap: the .fitnessrc.json
// requirements.maxAcceptances value when positive, else the default.
func maxAcceptances(root string) int {
	var c fileConfig
	if config.Read(root, &c) == nil && c.Requirements.MaxAcceptances > 0 {
		return c.Requirements.MaxAcceptances
	}
	return defaults.Requirements.MaxAcceptances
}

// capErrors fails each doc holding more than limit acceptances, at the first
// acceptance past the cap.
func capErrors(docs []requirement.Document, limit int) []string {
	var errs []string
	for _, d := range docs {
		if n := len(d.Acceptances); n > limit {
			msg := fmt.Sprintf("%s holds %d acceptances, at most %d; fold cases into table rows of one acceptance, or move a separate goal to its own doc", d.ID, n, limit)
			errs = append(errs, requirement.At(d.Path, d.Acceptances[limit].Line, msg))
		}
	}
	return errs
}
