package requirement

import "github.com/may-journal/fitness-checks/internal/config"

// DefaultMaxSteps is the step limit when .fitnessrc.json sets no
// requirements.maxSteps, a span of control of five plus or minus two.
const DefaultMaxSteps = 7

// stepConfig is the requirements key of .fitnessrc.json that sets the step limit.
type stepConfig struct {
	Requirements struct {
		MaxSteps int `json:"maxSteps"`
	} `json:"requirements"`
}

// MaxSteps is the step limit root's .fitnessrc.json sets, else the default. A chain holds two steps at least, so a
// smaller limit falls back to the default too.
func MaxSteps(root string) int {
	var c stepConfig
	if config.Read(root, &c) == nil && c.Requirements.MaxSteps >= 2 {
		return c.Requirements.MaxSteps
	}
	return DefaultMaxSteps
}
