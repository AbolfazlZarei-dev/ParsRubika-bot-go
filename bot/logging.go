package bot

import "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"

func (b *Bot) debugf(format string, a ...any) { core.Logf(b.logger, core.LogDebug, format, a...) }
func (b *Bot) infof(format string, a ...any)  { core.Logf(b.logger, core.LogInfo, format, a...) }
func (b *Bot) warnf(format string, a ...any)  { core.Logf(b.logger, core.LogWarn, format, a...) }
func (b *Bot) errorf(format string, a ...any) { core.Logf(b.logger, core.LogError, format, a...) }
