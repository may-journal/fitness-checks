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
)

// FileName is the config file at a repo's root.
const FileName = ".fitnessrc.json"

// Read fills into, a pointer holding the check's defaults, from the value
// at key in root's config file. A missing file or key changes nothing.
func Read(root, key string, into any) error {
	value, err := valueAt(root, key)
	if err != nil || value == nil {
		return err
	}
	if err := json.Unmarshal(value, into); err != nil {
		return fmt.Errorf("%s: %s: %w", FileName, key, err)
	}
	return nil
}

// valueAt is the raw value at key in root's config file, nil when the file
// or the key is missing.
func valueAt(root, key string) (json.RawMessage, error) {
	raw, err := os.ReadFile(filepath.Join(root, FileName))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return nil, fmt.Errorf("%s: %w", FileName, err)
	}
	return top[key], nil
}
