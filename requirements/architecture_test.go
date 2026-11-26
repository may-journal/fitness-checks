package requirements

import "github.com/may-journal/fitness-checks/internal/acceptance"

// The architecture fixtures live in acceptance, so each architecture check's
// own requirement tests share them with these.
const (
	acCode      = acceptance.AcCode
	acComponent = acceptance.AcComponent
	acContainer = acceptance.AcContainer
	acContext   = acceptance.AcContext
	acDoc       = acceptance.AcDoc
	acJourney   = acceptance.AcJourney
)

var (
	acDocWith          = acceptance.AcDocWith
	acFixture          = acceptance.AcFixture
	aoFixture          = acceptance.AoFixture
	architectureChecks = acceptance.ArchitectureChecks
	javaPom            = acceptance.JavaPom
)
