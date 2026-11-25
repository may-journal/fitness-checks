package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/may-journal/fitness-checks/internal/requirement"
)

// upFlag and downFlag run one migration by name, such as
// `fitness requirements --down 0003.6-gherkin`.
const (
	upFlag   = "--up"
	downFlag = "--down"
)

// rewrite is one change a migration makes to the doc at path under root.
type rewrite func(root, path string) error

// migrate rewrites each doc at paths as args ask, before the check judges.
func migrate(root string, paths, args []string) error {
	steps, err := rewrites(args)
	if err != nil || len(steps) == 0 {
		return err
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
	if slices.Contains(args, fixFlag) {
		return fixRewrites(), nil
	}
	for i := 0; i+1 < len(args); i++ {
		if args[i] == upFlag || args[i] == downFlag {
			return named(args[i], args[i+1])
		}
	}
	return nil, nil
}

// fixRewrites deletes Source lines, then runs every migration up.
func fixRewrites() []rewrite {
	steps := []rewrite{dropSource}
	for _, m := range requirement.Migrations {
		steps = append(steps, m.Up)
	}
	return steps
}

// named is the up or down of the migration id, or an error listing the names.
func named(flag, id string) ([]rewrite, error) {
	m, ok := requirement.FindMigration(id)
	if !ok {
		return nil, fmt.Errorf("no migration %q; run one of %v", id, migrationIDs())
	}
	if flag == upFlag {
		return []rewrite{m.Up}, nil
	}
	return []rewrite{m.Down}, nil
}

// migrationIDs lists every migration's name.
func migrationIDs() []string {
	var ids []string
	for _, m := range requirement.Migrations {
		ids = append(ids, m.ID)
	}
	return ids
}

// dropSource deletes the leftover Source line from the doc at path.
func dropSource(root, path string) error {
	docs, _, _ := requirement.Load(root, []string{path})
	if len(docs) == 0 {
		return nil
	}
	return requirement.DropSource(root, docs[0])
}
