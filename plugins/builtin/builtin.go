// Package builtin gathers the ready-made plugins that ship with
// ParsRubika so they can be registered in one call.
//
//	m := plugins.New(b)
//	builtin.RegisterAll(m, map[string]plugins.Config{
//		"access":    {"admins": []string{"b0abc"}},
//		"forcejoin": {"channels": []string{"c0xyz"}},
//	})
//	m.Install()
package builtin

import (
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/access"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/analytics"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/antilink"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/antispam"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/broadcast"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/flags"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/forcejoin"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/help"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/i18n"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/kv"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/moderation"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/ratelimit"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/requestlog"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/stats"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/users"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/welcome"
)

// All returns a fresh instance of every built-in plugin.
func All() []plugins.Plugin {
	return []plugins.Plugin{
		access.New(), analytics.New(), antilink.New(), antispam.New(), broadcast.New(),
		flags.New(), forcejoin.New(), help.New(), i18n.New(), kv.New(), moderation.New(),
		ratelimit.New(), requestlog.New(), stats.New(), users.New(), welcome.New(),
	}
}

// needsConfig lists plugins that refuse to start without configuration;
// RegisterAll keeps them registered but switched off until configured.
var needsConfig = map[string]bool{"forcejoin": true}

// RegisterAll registers every built-in plugin on m. cfg holds optional
// per-plugin configuration keyed by plugin name. Plugins that cannot
// work without configuration (forcejoin) are only enabled when their
// configuration is present.
func RegisterAll(m *plugins.Manager, cfg map[string]plugins.Config) error {
	for _, p := range All() {
		name := p.Meta().Name
		c, has := cfg[name]
		if needsConfig[name] && !has {
			continue
		}
		if err := m.Register(p, c); err != nil {
			return err
		}
	}
	return nil
}
