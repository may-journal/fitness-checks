package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/may-journal/fitness-checks/cmd/fitness-check-requirements/migrations"
	"github.com/may-journal/fitness-checks/internal/requirement"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

// upFlag and downFlag run one migration by name, such as
// `fitness requirements --down 0003.2-gherkin`.
const (
	upFlag   = "--up"
	downFlag = "--down"
)

// rewrite is one change a migration makes to the doc at path under root.
type rewrite func(root, path string) error

// argError is a flag mistake, which fails the check instead of stopping it.
type argError struct{ error }

// migrateResult fails the check on a flag mistake and stops it on any other
// error.
func migrateResult(err error) (checkkit.Result, error) {
	var bad argError
	if errors.As(err, &bad) {
		return checkkit.Fail(0, bad.Error()), nil
	}
	return checkkit.Result{}, err
}

// migrate rewrites each doc at paths as args ask, before the check judges.
func migrate(root string, paths, args []string) error {
	steps, err := rewrites(args)
	if err != nil {
		return argError{err}
	}
	for _, path := range paths {
		for _, step := range steps {
			err = errors.Join(err, step(root, path))
		}
	}
	return err
}

// rewrites is what args ask for: the fix, or one migration up or down.
func rewrites(args []string) ([]rewrite, error) {
	if err := flagMistake(args); err != nil {
		return nil, err
	}
	i := slices.IndexFunc(args, isMigrationFlag)
	if i < 0 {
		return fixOnly(args), nil
	}
	if err := migrationArgs(args, i); err != nil {
		return nil, err
	}
	return named(args[i], args[i+1])
}

// isMigrationFlag reports whether a is --up or --down.
func isMigrationFlag(a string) bool {
	return a == upFlag || a == downFlag
}

// flagMistake fails a migration flag joined to its name with "=", or more
// than one migration flag.
func flagMistake(args []string) error {
	for _, a := range args {
		if flag, name, ok := strings.Cut(a, "="); ok && isMigrationFlag(flag) {
			return fmt.Errorf("write %s %s, with a space", flag, name)
		}
	}
	if len(slices.DeleteFunc(slices.Clone(args), func(a string) bool { return !isMigrationFlag(a) })) > 1 {
		return fmt.Errorf("run one migration at a time")
	}
	return nil
}

// fixOnly is the fix when args ask for it.
func fixOnly(args []string) []rewrite {
	if slices.Contains(args, fixFlag) {
		return fixRewrites()
	}
	return nil
}

// migrationArgs fails the flag at index i when the fix comes with it or no
// migration name follows it.
func migrationArgs(args []string, i int) error {
	if slices.Contains(args, fixFlag) {
		return fmt.Errorf("run %s or %s alone; %s already runs every migration up", args[i], fixFlag, fixFlag)
	}
	if i+1 == len(args) {
		return fmt.Errorf("name the migration %s runs, one of %v", args[i], migrations.IDs())
	}
	return nil
}

// fixRewrites deletes Source lines, then runs every migration up.
func fixRewrites() []rewrite {
	steps := []rewrite{dropSource}
	for _, m := range migrations.All {
		steps = append(steps, m.Up)
	}
	return steps
}

// named is the up or down of the migration id, or an error listing the names.
func named(flag, id string) ([]rewrite, error) {
	m, ok := migrations.Find(id)
	if !ok {
		return nil, fmt.Errorf("no migration %q; run one of %v", id, migrations.IDs())
	}
	if flag == upFlag {
		return []rewrite{m.Up}, nil
	}
	return []rewrite{m.Down}, nil
}

// dropSource deletes the leftover Source line from the doc at path.
func dropSource(root, path string) error {
	docs, _, _ := requirement.Load(root, []string{path})
	if len(docs) == 0 {
		return nil
	}
	return requirement.DropSource(root, docs[0])
}
