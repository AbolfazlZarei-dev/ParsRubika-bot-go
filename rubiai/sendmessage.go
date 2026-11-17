package rubiai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	coreerrors "github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core/errors"
)

// SendMessageOptions configures an AI sendMessage call.
type SendMessageOptions struct {
	// Message is the text to send to the assistant. Required.
	Message string
	// SessionID continues an existing conversation; empty starts a new one.
	SessionID string
	// IsSupport marks the message as a support request.
	IsSupport bool
	// GenerateImage requests image generation.
	GenerateImage bool
}

// SendMessageResult is the assembled outcome of an SSE-streamed
// sendMessage call. The TextContent field holds the concatenated
// assistant reply (incrementally delivered via updateMessage events).
type SendMessageResult struct {
	SessionID   string
	MessageID   string
	Title       string
	TextContent string
	// Extra captures any other updateMessage parameters (set_value /
	// append_value) so callers do not lose data the library does not
	// model explicitly.
	Extra map[string]any
	// Error is the server-reported status_det when an `error` event
	// arrived, else empty.
	Error       string
	ErrorStatus string
}

// SendMessage sends a message to the AI service and returns the fully
// assembled result. It blocks until the SSE stream terminates.
//
// For incremental/token-by-token output use SendMessageStream, which
// invokes a callback for each parsed event as it arrives.
func (r *RubiaAI) SendMessage(ctx context.Context, opts SendMessageOptions) (*SendMessageResult, error) {
	return r.SendMessageStream(ctx, opts, nil)
}

// SendMessageStream is like SendMessage but invokes `onEvent` for each
// parsed SSE event as soon as it arrives (event name + parsed JSON data).
// The callback is invoked from the goroutine reading the response body,
// so it must not block on the same context. Pass nil to behave like
// SendMessage. The returned result is the assembled outcome.
func (r *RubiaAI) SendMessageStream(ctx context.Context, opts SendMessageOptions, onEvent func(event string, data map[string]any)) (*SendMessageResult, error) {
	if r.auth == "" {
		return nil, coreerrors.New(coreerrors.KindAuth, "rubiai: auth is required")
	}
	if strings.TrimSpace(opts.Message) == "" {
		return nil, coreerrors.New(coreerrors.KindValidation, "rubiai: message is required")
	}
	payload := map[string]any{
		"api_version": "1",
		"auth":        r.auth,
		"data": map[string]any{
			"method": "sendMessage",
			"client": ClientInfo,
			"input": map[string]any{
				"message":        opts.Message,
				"session_id":     opts.SessionID,
				"generate_image": opts.GenerateImage,
				"is_support":     opts.IsSupport,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, coreerrors.Wrap(coreerrors.KindAPI, "rubiai: marshal sendMessage payload", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("User-Agent", r.userAgent)

	resp, err := r.transport.HTTPClient().Do(req)
	if err != nil {
		return nil, coreerrors.Wrap(coreerrors.KindNetwork, "rubiai: sendMessage request", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, coreerrors.New(coreerrors.KindAPI, fmt.Sprintf("rubiai: sendMessage HTTP %d", resp.StatusCode))
	}

	result := &SendMessageResult{Extra: map[string]any{}}
	// SSE frames are delimited by blank lines. A scanner with a generous
	// buffer handles large text_content deltas; 1 MiB is well above any
	// realistic assistant chunk.
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)

	var (
		eventName string
		dataLines []string
		mu        sync.Mutex
	)
	flush := func() {
		if len(dataLines) == 0 {
			eventName = ""
			return
		}
		for _, d := range dataLines {
			parsed := parseSSEData(d)
			if onEvent != nil {
				onEvent(eventName, parsed)
			}
			assembleEvent(result, parsed)
		}
		eventName = ""
		dataLines = dataLines[:0]
	}
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		switch {
		case strings.HasPrefix(line, "event:"):
			eventName = strings.TrimSpace(line[len("event:"):])
		case strings.HasPrefix(line, "data:"):
			dataLines = append(dataLines, strings.TrimSpace(line[len("data:"):]))
		case line == "":
			mu.Lock()
			flush()
			mu.Unlock()
		default:
			// Unknown SSE field (id:, retry:, comment) — ignore.
		}
	}
	// Flush any trailing events not terminated by a blank line.
	mu.Lock()
	flush()
	mu.Unlock()
	if err := scanner.Err(); err != nil {
		return result, coreerrors.Wrap(coreerrors.KindNetwork, "rubiai: read SSE stream", err)
	}
	if result.Error == "SERVER_ERROR" {
		return result, coreerrors.New(coreerrors.KindAPI, "rubiai: AI server reported SERVER_ERROR")
	}
	return result, nil
}

// parseSSEData decodes a JSON data: line into a map. Non-JSON lines (the
// server occasionally sends a stray keep-alive) yield an empty map.
func parseSSEData(raw string) map[string]any {
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return map[string]any{}
	}
	return m
}

// assembleEvent folds a single SSE event into the running result,
// mirroring the Python _assemble_result logic:
//
//   - createSession   → session_id, title
//   - createMessage   → message_id, session_id
//   - updateSession   → title (when parameter_name=="title")
//   - updateMessage   → concatenate text_content append_value; other
//     parameters land in Extra
//   - error           → Error / ErrorStatus
func assembleEvent(result *SendMessageResult, data map[string]any) {
	if len(data) == 0 {
		return
	}
	// The Python code branches on the *event name*, but the event name is
	// not passed here. In practice the server also tags the data with a
	// "type"/"event" field; we additionally infer from the payload shape
	// so assembly is robust even when the server omits the event tag.
	switch {
	case hasKey(data, "session"):
		if sess, ok := data["session"].(map[string]any); ok {
			if v, ok := sess["session_id"].(string); ok && v != "" {
				result.SessionID = v
			}
			if v, ok := sess["title"].(string); ok && v != "" {
				result.Title = v
			}
		}
	case hasKey(data, "message"):
		if msg, ok := data["message"].(map[string]any); ok {
			if v, ok := msg["message_id"].(string); ok && v != "" {
				result.MessageID = v
			}
			if v, ok := msg["session_id"].(string); ok && v != "" {
				result.SessionID = v
			}
		}
	case hasKey(data, "parameter_name"):
		param, _ := data["parameter_name"].(string)
		switch param {
		case "title":
			if v, ok := data["set_value"].(string); ok {
				result.Title = v
			}
		case "text_content":
			if v, ok := data["append_value"].(string); ok {
				result.TextContent += v
			}
		default:
			if v, ok := data["set_value"]; ok && v != nil {
				result.Extra[param] = v
			} else if v, ok := data["append_value"]; ok && v != nil {
				result.Extra[param] = v
			}
		}
	case hasKey(data, "status_det") || hasKey(data, "data"):
		// error event: {"data":{"status_det":"SERVER_ERROR","status":"..."}}
		if d, ok := data["data"].(map[string]any); ok {
			if v, ok := d["status_det"].(string); ok {
				result.Error = v
			}
			if v, ok := d["status"].(string); ok {
				result.ErrorStatus = v
			}
		} else {
			if v, ok := data["status_det"].(string); ok {
				result.Error = v
			}
			if v, ok := data["status"].(string); ok {
				result.ErrorStatus = v
			}
		}
	}
}

func hasKey(m map[string]any, k string) bool {
	_, ok := m[k]
	return ok
}
