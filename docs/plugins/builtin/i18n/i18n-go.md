# `i18n` — Package i18n is a plugin that publishes a translator ("i18n" service)

_Source: [`i18n.go`](../../../../ParsRubika/plugins/builtin/i18n/i18n.go)_

## Symbols

### `Translator`

```go
type Translator interface
```

Translator is the published service.

### `Plugin`

```go
type Plugin struct
```

Plugin implements plugins.Plugin and Translator.

### `New`

```go
func New() *Plugin
```

New creates the plugin.

### `DefaultConfig`

```go
func (p *Plugin) DefaultConfig() plugins.Config
```

DefaultConfig implements plugins.Defaulter.

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.

### `AddLanguage`

```go
func (p *Plugin) AddLanguage(lang string, msgs map[string]string)
```

AddLanguage merges translations for lang (usable from code).

### `T`

```go
func (p *Plugin) T(userID, key string, args ...any) string
```

T implements Translator.

### `Lang`

```go
func (p *Plugin) Lang(userID string) string
```

Lang implements Translator.

### `SetLang`

```go
func (p *Plugin) SetLang(userID, lang string) bool
```

SetLang implements Translator; it refuses unknown languages.

### `Languages`

```go
func (p *Plugin) Languages() []string
```

Languages implements Translator.
