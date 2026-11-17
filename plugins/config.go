package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config is the free-form configuration of one plugin. Values coming
// from JSON files are float64/[]any/map[string]any; the typed getters
// understand both those and native Go values.
type Config map[string]any

// Clone returns a deep-enough copy (nested maps and slices are copied).
func (c Config) Clone() Config {
	out := make(Config, len(c))
	for k, v := range c {
		out[k] = cloneValue(v)
	}
	return out
}

func cloneValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		m := make(map[string]any, len(t))
		for k, x := range t {
			m[k] = cloneValue(x)
		}
		return m
	case Config:
		return t.Clone()
	case []any:
		s := make([]any, len(t))
		for i, x := range t {
			s[i] = cloneValue(x)
		}
		return s
	case []string:
		return append([]string(nil), t...)
	}
	return v
}

// Merge returns defaults overlaid with override.
func Merge(defaults, override Config) Config {
	out := defaults.Clone()
	for k, v := range override {
		out[k] = cloneValue(v)
	}
	return out
}

// String returns the value of key as a string.
func (c Config) String(key, def string) string {
	switch t := c[key].(type) {
	case string:
		return t
	case nil:
		return def
	default:
		return fmt.Sprint(t)
	}
}

// Int returns the value of key as an int.
func (c Config) Int(key string, def int) int {
	switch t := c[key].(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	case string:
		if n, err := strconv.Atoi(t); err == nil {
			return n
		}
	}
	return def
}

// Int64 returns the value of key as an int64.
func (c Config) Int64(key string, def int64) int64 {
	switch t := c[key].(type) {
	case int:
		return int64(t)
	case int64:
		return t
	case float64:
		return int64(t)
	case string:
		if n, err := strconv.ParseInt(t, 10, 64); err == nil {
			return n
		}
	}
	return def
}

// Bool returns the value of key as a bool.
func (c Config) Bool(key string, def bool) bool {
	switch t := c[key].(type) {
	case bool:
		return t
	case string:
		if b, err := strconv.ParseBool(t); err == nil {
			return b
		}
	}
	return def
}

// Duration returns the value of key as a time.Duration. Strings such as
// "10s" or "2m" are parsed; plain numbers are read as seconds.
func (c Config) Duration(key string, def time.Duration) time.Duration {
	switch t := c[key].(type) {
	case time.Duration:
		return t
	case string:
		if d, err := time.ParseDuration(t); err == nil {
			return d
		}
	case int:
		return time.Duration(t) * time.Second
	case int64:
		return time.Duration(t) * time.Second
	case float64:
		return time.Duration(t * float64(time.Second))
	}
	return def
}

// Strings returns the value of key as a string slice.
func (c Config) Strings(key string) []string {
	switch t := c[key].(type) {
	case []string:
		return append([]string(nil), t...)
	case []any:
		out := make([]string, 0, len(t))
		for _, x := range t {
			out = append(out, fmt.Sprint(x))
		}
		return out
	case string:
		if t != "" {
			return []string{t}
		}
	}
	return nil
}

// Has reports whether key is present.
func (c Config) Has(key string) bool { _, ok := c[key]; return ok }

// ParseConfigJSON parses a document of the form
//
//	{"plugins": {"antispam": {"max_messages": 5}, "welcome": {...}}}
//
// into per-plugin configurations (keys are lower-cased).
func ParseConfigJSON(data []byte) (map[string]Config, error) {
	var doc struct {
		Plugins map[string]map[string]any `json:"plugins"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("plugins: parse config: %w", err)
	}
	out := make(map[string]Config, len(doc.Plugins))
	for name, c := range doc.Plugins {
		out[normalize(name)] = Config(c)
	}
	return out, nil
}

// LoadConfigFile reads a JSON configuration file (see ParseConfigJSON).
func LoadConfigFile(path string) (map[string]Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("plugins: read config: %w", err)
	}
	return ParseConfigJSON(data)
}
