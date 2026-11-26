// Package migrations lists the requirement doc format migrations, oldest
// first. Each lives in a folder named by its migration, such as
// 0003.2-gherkin, which holds its up.go and down.go.
package migrations

import gherkin "github.com/may-journal/fitness-checks/cmd/fitness-check-requirements/migrations/0003.2-gherkin"

// Migration is one reversible rewrite of a requirement doc's format. Its ID
// is the acceptance that states it, so the name holds as formats change.
type Migration struct {
	ID   string
	Up   func(root, path string) error
	Down func(root, path string) error
}

// All lists every migration oldest first. The fix runs each one up.
var All = []Migration{
	{ID: gherkin.ID, Up: gherkin.Up, Down: gherkin.Down},
}

// Find is the migration named id.
func Find(id string) (Migration, bool) {
	for _, m := range All {
		if m.ID == id {
			return m, true
		}
	}
	return Migration{}, false
}

// IDs lists every migration's name.
func IDs() []string {
	var ids []string
	for _, m := range All {
		ids = append(ids, m.ID)
	}
	return ids
}
