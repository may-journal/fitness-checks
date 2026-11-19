// Package waiting reads adr-as-issue's adrAsIssue key, the ADR files still
// waiting to move into issues. adr-as-issue owns the key and reports a bad
// entry; adr-checks reads the same list so a waiting file is not judged twice.
package waiting

import (
	"fmt"
	"path"
	"strings"

	"github.com/may-journal/fitness-checks/internal/config"
)

// Key is adr-as-issue's key in .fitnessrc.json.
const Key = "adrAsIssue"

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

// Allow is the loaded allow globs, each with a reason.
type Allow struct {
	// Globs are the entries' paths, in the config's order.
	Globs []string
}

// Load reads adrAsIssue.allow from root's config file. Each entry that
// lacks a reason or is no glob is a problem the caller fails with.
func Load(root string) (Allow, []string, error) {
	var opts options
	if err := config.Read(root, Key, &opts); err != nil {
		return Allow{}, nil, err
	}
	var allow Allow
	var problems []string
	for _, a := range opts.Allow {
		if problem := allowProblem(a.Path, a.Reason); problem != "" {
			problems = append(problems, problem)
			continue
		}
		allow.Globs = append(allow.Globs, a.Path)
	}
	return allow, problems, nil
}

// allowProblem is what is wrong with one allow entry, or "".
func allowProblem(glob, reason string) string {
	if _, err := path.Match(glob, ""); err != nil {
		return fmt.Sprintf("adrAsIssue.allow path %q is not a glob: %v", glob, err)
	}
	if strings.TrimSpace(reason) == "" {
		return fmt.Sprintf("adrAsIssue.allow path %q needs a reason", glob)
	}
	return ""
}

// Covering is the first allow glob that matches rel, or "".
func (a Allow) Covering(rel string) string {
	for _, glob := range a.Globs {
		if ok, _ := path.Match(glob, rel); ok {
			return glob
		}
	}
	return ""
}
