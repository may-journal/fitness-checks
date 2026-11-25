package requirement

// Migration is one reversible rewrite of a requirement doc's format. Its ID
// is the acceptance that states it, so the name holds as formats change.
type Migration struct {
	ID   string
	Up   func(root, path string) error
	Down func(root, path string) error
}

// Migrations lists every migration oldest first. The fix runs each one up.
var Migrations = []Migration{
	{ID: "0003.6-gherkin", Up: ToGherkin, Down: ToBullets},
}

// FindMigration is the migration named id.
func FindMigration(id string) (Migration, bool) {
	for _, m := range Migrations {
		if m.ID == id {
			return m, true
		}
	}
	return Migration{}, false
}
