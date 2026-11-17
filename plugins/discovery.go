package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

// DiscoverFromDir loads every shared-library plugin (*.so, built with
// `go build -buildmode=plugin`) found in dir. Each library must export
//
//	func Register(m *plugins.Manager) error
//
// which is called once. Go's plugin mechanism only works on Linux and
// macOS with cgo enabled and a matching toolchain; elsewhere the call
// returns an error. A missing directory is not an error.
func (m *Manager) DiscoverFromDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("plugins: read dir %q: %w", dir, err)
	}
	var errs []error
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".so") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		lib, err := plugin.Open(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("plugins: open %s: %w", path, err))
			continue
		}
		sym, err := lib.Lookup("Register")
		if err != nil {
			errs = append(errs, fmt.Errorf("plugins: %s: %w", path, err))
			continue
		}
		register, ok := sym.(func(*Manager) error)
		if !ok {
			errs = append(errs, fmt.Errorf("plugins: %s: Register must be func(*plugins.Manager) error", path))
			continue
		}
		if err := register(m); err != nil {
			errs = append(errs, fmt.Errorf("plugins: %s: %w", path, err))
		}
	}
	return errors.Join(errs...)
}

// Manifest is a small JSON file that configures an already registered
// plugin, so operators can tune a deployment without recompiling:
//
//	{"name": "antispam", "enabled": true, "config": {"max_messages": 8}}
type Manifest struct {
	Name    string `json:"name"`
	Enabled *bool  `json:"enabled,omitempty"`
	Config  Config `json:"config,omitempty"`
}

// ApplyManifests reads every *.plugin.json file in dir and applies the
// configuration to the registered plugin it names. It returns the names
// of the plugins the manifests asked to enable (enabled omitted counts
// as true). Unknown plugins are reported in the returned error but do
// not stop the others. A missing directory is not an error.
func (m *Manager) ApplyManifests(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("plugins: read dir %q: %w", dir, err)
	}
	var (
		names []string
		errs  []error
	)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".plugin.json") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		var mf Manifest
		if err := json.Unmarshal(data, &mf); err != nil {
			errs = append(errs, fmt.Errorf("plugins: %s: %w", path, err))
			continue
		}
		name := normalize(mf.Name)
		if len(mf.Config) > 0 {
			if err := m.Configure(name, mf.Config); err != nil {
				errs = append(errs, fmt.Errorf("plugins: %s: %w", path, err))
				continue
			}
		}
		if mf.Enabled == nil || *mf.Enabled {
			names = append(names, name)
		}
	}
	return names, errors.Join(errs...)
}
