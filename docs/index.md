# ParsRubika documentation

Per-package reference. Each package has an `index.md` overview plus one
file per source file (method group) listing every exported symbol with
its signature and documentation.

## Quick start (everything in one import)

The root package wires all five clients together — account, bot,
rubino, tv and rubiai — and ships one-line helpers so a whole job
fits in a few lines of code:

```go
package main

import pr "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4"

func main() {
    // one line, one job: the dual (hangi) profile on a channel
    pr.Hangi("session", "+989123456789",
        "c0DBedk0525ec07fb54a35b01baf2303", "1.jpg", "2.jpg")
}
```

See [`parsrubika`](parsrubika/index.md) for the full facade API

## Packages

| Package | Path |
|---|---|
| [`parsrubika`](parsrubika/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4` |
| [`account`](account/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account` |
| [`listeners`](account/listeners/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account/listeners` |
| [`matchers`](account/matchers/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/account/matchers` |
| [`bot`](bot/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot` |
| [`flow`](bot/flow/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/flow` |
| [`models`](bot/models/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models` |
| [`router`](bot/router/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/router` |
| [`rules`](bot/rules/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules` |
| [`widgets`](bot/widgets/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/widgets` |
| [`core`](core/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core` |
| [`convert`](core/convert/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/convert` |
| [`errors`](core/errors/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors` |
| [`guard`](core/guard/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/guard` |
| [`guid`](core/guid/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/guid` |
| [`idgen`](core/idgen/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/idgen` |
| [`validate`](core/validate/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/validate` |
| [`enums`](enums/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/enums` |
| [`events`](events/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events` |
| [`ai`](extras/ai/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai` |
| [`anthropic`](extras/ai/providers/anthropic/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/anthropic` |
| [`deepseek`](extras/ai/providers/deepseek/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/deepseek` |
| [`gapgpt`](extras/ai/providers/gapgpt/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/gapgpt` |
| [`gemini`](extras/ai/providers/gemini/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/gemini` |
| [`groq`](extras/ai/providers/groq/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/groq` |
| [`httpx`](extras/ai/providers/internal/httpx/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/internal/httpx` |
| [`oai`](extras/ai/providers/internal/oai/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/internal/oai` |
| [`openai`](extras/ai/providers/openai/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/openai` |
| [`openrouter`](extras/ai/providers/openrouter/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/openrouter` |
| [`zai`](extras/ai/providers/zai/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/ai/providers/zai` |
| [`cron`](extras/cron/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/cron` |
| [`forms`](extras/forms/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/forms` |
| [`handy`](extras/handy/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/handy` |
| [`qr`](extras/qr/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/qr` |
| [`scheduler`](extras/scheduler/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/scheduler` |
| [`tasks`](extras/tasks/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/tasks` |
| [`voicekit`](extras/voicekit/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/extras/voicekit` |
| [`files`](files/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/files` |
| [`logger`](logger/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/logger` |
| [`netutil`](netutil/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/netutil` |
| [`plugins`](plugins/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins` |
| [`builtin`](plugins/builtin/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin` |
| [`access`](plugins/builtin/access/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/access` |
| [`analytics`](plugins/builtin/analytics/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/analytics` |
| [`antilink`](plugins/builtin/antilink/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/antilink` |
| [`antispam`](plugins/builtin/antispam/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/antispam` |
| [`broadcast`](plugins/builtin/broadcast/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/broadcast` |
| [`flags`](plugins/builtin/flags/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/flags` |
| [`forcejoin`](plugins/builtin/forcejoin/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/forcejoin` |
| [`help`](plugins/builtin/help/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/help` |
| [`i18n`](plugins/builtin/i18n/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/i18n` |
| [`kv`](plugins/builtin/kv/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/kv` |
| [`moderation`](plugins/builtin/moderation/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/moderation` |
| [`ratelimit`](plugins/builtin/ratelimit/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/ratelimit` |
| [`requestlog`](plugins/builtin/requestlog/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/requestlog` |
| [`stats`](plugins/builtin/stats/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/stats` |
| [`users`](plugins/builtin/users/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/users` |
| [`welcome`](plugins/builtin/welcome/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/builtin/welcome` |
| [`atomicfile`](plugins/internal/atomicfile/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/internal/atomicfile` |
| [`plugintest`](plugins/plugintest/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/plugins/plugintest` |
| [`rubiai`](rubiai/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/rubiai` |
| [`rubino`](rubino/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/rubino` |
| [`security`](security/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/security` |
| [`storage`](storage/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/storage` |
| [`text`](text/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/text` |
| [`transport`](transport/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/transport` |
| [`tv`](tv/index.md) | `github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/tv` |

## Examples

| Example | Doc |
|---|---|
| [`hangprofile`](examples/hangprofile/index.md) | the dual (hangi) profile, one call |

