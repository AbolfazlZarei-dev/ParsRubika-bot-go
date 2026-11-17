# `plugins` — Discovery

_Source: [`discovery.go`](../../ParsRubika/plugins/discovery.go)_

## Symbols

### `DiscoverFromDir`

```go
func (m *Manager) DiscoverFromDir(dir string) error
```

DiscoverFromDir loads every shared-library plugin (*.so, built with
`go build -buildmode=plugin`) found in dir. Each library must export

func Register(m *plugins.Manager) error

which is called once. Go's plugin mechanism only works on Linux and
macOS with cgo enabled and a matching toolchain; elsewhere the call
returns an error. A missing directory is not an error.

### `Manifest`

```go
type Manifest struct
```

Manifest is a small JSON file that configures an already registered
plugin, so operators can tune a deployment without recompiling:

{"name": "antispam", "enabled": true, "config": {"max_messages": 8}}

### `ApplyManifests`

```go
func (m *Manager) ApplyManifests(dir string) ([]string, error)
```

ApplyManifests reads every *.plugin.json file in dir and applies the
configuration to the registered plugin it names. It returns the names
of the plugins the manifests asked to enable (enabled omitted counts
as true). Unknown plugins are reported in the returned error but do
not stop the others. A missing directory is not an error.
