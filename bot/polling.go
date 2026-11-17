package bot

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// GetUpdates fetches a single batch of updates without touching the
// stored offset or the duplicate filter.
func (b *Bot) GetUpdates(ctx context.Context, limit int, offsetID string) ([]*models.Update, error) {
	payload := map[string]any{"limit": limit}
	if offsetID != "" {
		payload["offset_id"] = offsetID
	}
	data, err := b.makeRequest(ctx, "getUpdates", payload)
	if err != nil {
		return nil, err
	}
	raw, _ := data["updates"].([]any)
	out := make([]*models.Update, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		upd, err := parseUpdate(m, b)
		if err != nil || upd == nil {
			continue
		}
		out = append(out, upd)
	}
	return out, nil
}

// Updater performs one long-poll round: it uses the stored offset,
// advances it, filters duplicates and returns the fresh updates.
//
// When the bot starts without a stored offset and SkipBacklog is on,
// the first successful batch is discarded so old messages are not
// replayed.
func (b *Bot) Updater(ctx context.Context, limit int, pollTimeout time.Duration) ([]*models.Update, error) {
	payload := map[string]any{"limit": limit}
	b.offsetMu.Lock()
	offset := b.nextOffsetID
	b.offsetMu.Unlock()
	if offset != "" {
		payload["offset_id"] = offset
	}
	if pollTimeout > 0 {
		payload["timeout"] = int(pollTimeout / time.Second)
	}
	data, err := b.call(ctx, "getUpdates", payload, false)
	if err != nil {
		return nil, err
	}
	if newOff := parsedString(data, "next_offset_id"); newOff != "" {
		b.offsetMu.Lock()
		changed := newOff != b.nextOffsetID
		b.nextOffsetID = newOff
		b.offsetMu.Unlock()
		if changed {
			b.persistOffset(newOff)
		}
	}
	if b.skipPending.CompareAndSwap(true, false) {
		return nil, nil
	}
	raw, _ := data["updates"].([]any)
	out := make([]*models.Update, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if key := buildRawUpdateKey(m); key != "" && !b.seen.add("raw:"+key) {
			b.stats.Duplicate.Add(1)
			continue
		}
		upd, err := parseUpdate(m, b)
		if err != nil || upd == nil {
			continue
		}
		out = append(out, upd)
	}
	return out, nil
}

// ResetOffset forgets the stored polling offset.
func (b *Bot) ResetOffset() {
	b.offsetMu.Lock()
	b.nextOffsetID = ""
	b.offsetMu.Unlock()
	b.persistOffset("")
}

// persistOffset writes the offset atomically so a crash never leaves a
// truncated file behind.
func (b *Bot) persistOffset(v string) {
	if !b.opts.PersistOffset || b.offsetFile == "" {
		return
	}
	if v == "" {
		_ = os.Remove(b.offsetFile)
		return
	}
	tmp := b.offsetFile + ".tmp"
	if err := os.WriteFile(tmp, []byte(v), 0o600); err != nil {
		b.warnf("cannot persist offset: %v", err)
		return
	}
	if err := os.Rename(tmp, b.offsetFile); err != nil {
		_ = os.Remove(tmp)
		b.warnf("cannot persist offset: %v", err)
	}
}

// OffsetFile returns the path of the on-disk offset file (empty when
// persistence is disabled).
func (b *Bot) OffsetFile() string {
	if b.offsetFile == "" {
		return ""
	}
	if abs, err := filepath.Abs(b.offsetFile); err == nil {
		return abs
	}
	return b.offsetFile
}

// runPolling is the long-poll loop.
func (b *Bot) runPolling(ctx context.Context) error {
	if err := b.onStart(ctx); err != nil {
		return err
	}
	defer b.onStop(context.Background())
	defer b.drain()

	failures := 0
	for {
		if err := ctx.Err(); err != nil {
			return nil
		}
		updates, err := b.Updater(ctx, 100, b.opts.LongPollTimeout)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			failures++
			b.warnf("polling error (%d in a row): %v", failures, err)
			if !b.sleepBackoff(ctx, failures-1, retryAfterOf(err)) {
				return nil
			}
			continue
		}
		failures = 0
		for _, u := range updates {
			b.submit(ctx, u)
		}
		if len(updates) == 0 && b.opts.PollInterval > 0 {
			t := time.NewTimer(b.opts.PollInterval)
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return nil
			}
		}
	}
}

// retryAfterOf extracts the server suggested delay from err.
func retryAfterOf(err error) time.Duration {
	if ae, ok := err.(*APIError); ok {
		return ae.RetryAfter
	}
	return 0
}

// containsString reports whether `in` contains `want`.
func containsString(in []string, want string) bool {
	for _, s := range in {
		if s == want {
			return true
		}
	}
	return false
}
