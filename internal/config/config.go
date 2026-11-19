// Package config reads each check's own keys from a repo's .fitnessrc.json.
// Each check declares its keys, defaults and docs in its own package.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
	"github.com/may-journal/fitness-runner/go/pkg/conf"
)

// FileName is the config file at a repo's root.
const FileName = ".fitnessrc.json"

// Keys lists the keys a check reads for its --describe. Defaults is the
// struct Read fills, whose fields are top-level keys holding their defaults,
// and source is the embedded file that documents them.
func Keys(defaults any, source []byte) func() ([]checkkit.ConfigKey, error) {
	return func() ([]checkkit.ConfigKey, error) {
		return conf.KeysOf(defaults, source)
	}
}

// Read fills into, a pointer to a struct whose fields are the check's
// top-level keys, from root's config file. A missing file or key changes
// nothing, and keys into does not name are left alone.
func Read(root string, into any) error {
	raw, err := os.ReadFile(filepath.Join(root, FileName))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, into); err != nil {
		return fmt.Errorf("%s: %w", FileName, err)
	}
	return nil
}
