package account

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/enums"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/files"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/text"
	"path/filepath"
	"strings"
	"time"
)

// Sending messages and media.

// SendOptions carries the optional arguments for SendMessage.
type SendOptions struct {
	ReplyToMessageID string
	Sticker          map[string]any
	FileInline       *FileInlineInput
	IsSpoil          bool
	AutoDelete       time.Duration
	ParseMode        enums.ParseMode
	Metadata         *text.Result
}

// FileInlineInput bundles everything SendMessage needs to know about
// attaching a file to the outgoing message.
type FileInlineInput struct {
	Update    *events.Event // upload metadata
	Type      string        // "File", "Image", "Voice", "Music", "Gif", "Video", "VideoMessage"
	Width     int
	Height    int
	Time      int
	Performer string
}

// SendMessage sends a message to the given object GUID.
func (a *Account) SendMessage(ctx context.Context, objectGUID, msg string, opts ...SendOptions) (*events.Event, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}

	// Allow magic aliases "me" / "cloud" / "self" → current user.
	if g := normalizeSelf(objectGUID, a.guid); g != objectGUID {
		objectGUID = g
	}

	// Build the base input. IMPORTANT: `reply_to_message_id` MUST be
	// omitted when empty — the server rejects the literal empty string
	// with `INVALID_INPUT: sendMessage failed`. The wire payload must
	// use JSON null (or simply omit the key) to mean "no reply";
	// sending `""` instead fails the request.
	input := map[string]any{
		"object_guid": objectGUID,
		"rnd":         randInt63(),
	}
	if o.ReplyToMessageID != "" {
		input["reply_to_message_id"] = o.ReplyToMessageID
	}
	if msg != "" {
		if o.Metadata != nil {
			mergeMetadata(input, o.Metadata)
		}
		input["text"] = strings.TrimSpace(msg)
		mode := o.ParseMode
		if mode == "" {
			mode = a.parseMode
		}
		if mode == enums.ParseModeHTML {
			md := a.formatter.Format(a.formatter.ToMarkdown(msg))
			mergeMetadata(input, &md)
		} else if mode == enums.ParseModeMarkdown {
			md := a.formatter.Format(msg)
			mergeMetadata(input, &md)
		}
	}
	if o.Sticker != nil {
		input["sticker"] = o.Sticker
	}
	if o.FileInline != nil && o.FileInline.Update != nil {
		upd := o.FileInline.Update
		// IMPORTANT: `file_id` and `access_hash_rec` may come back as
		// JSON numbers depending on the server; use AnyStr to coerce.
		// Also include `file_name` — the upload returns it and the
		// wire payload includes it.
		file := map[string]any{
			"file_id":         upd.AnyStr("file_id"),
			"dc_id":           upd.Int("dc_id"),
			"size":            upd.Int("size"),
			"mime":            upd.Str("mime"),
			"file_name":       upd.Str("file_name"),
			"access_hash_rec": upd.AnyStr("access_hash_rec"),
		}
		// When a thumbnail was generated, attach its base64 payload
		// together with the real image dimensions.
		if thumb := upd.Str("thumb_inline"); thumb != "" {
			file["thumb_inline"] = thumb
		}

		t := o.FileInline.Type
		if t == "" {
			t = "File"
		}
		// Wire type: 'VideoMessage' → 'Video' (with is_round), otherwise
		// the type as-is ('Image', 'Video', 'Voice', 'Music', 'Gif',
		// 'File'). Hard-coding 'Video' for every type caused the server
		// to reject /photo with INVALID_INPUT: sendMessage failed.
		wireType := t
		if t == "VideoMessage" {
			wireType = "Video"
			file["is_round"] = true
		}
		file["type"] = wireType
		// Sensible defaults are applied below: sending 0 here makes the
		// server reject the request.
		if o.FileInline.Time > 0 {
			file["time"] = o.FileInline.Time
		} else {
			file["time"] = 1
		}
		// Width/height: prefer the explicit user-supplied value, then
		// the actual image dimensions (from the thumbnail), then 200.
		imgWidth := upd.Int("width")
		imgHeight := upd.Int("height")
		if o.FileInline.Width > 0 {
			file["width"] = o.FileInline.Width
		} else if imgWidth > 0 {
			file["width"] = imgWidth
		} else {
			file["width"] = 200
		}
		if o.FileInline.Height > 0 {
			file["height"] = o.FileInline.Height
		} else if imgHeight > 0 {
			file["height"] = imgHeight
		} else {
			file["height"] = 200
		}
		file["music_performer"] = o.FileInline.Performer
		file["is_spoil"] = o.IsSpoil
		input["file_inline"] = file

		if a.debug {
			if jb, err := json.Marshal(file); err == nil {
				a.debugf("sendMessage file_inline=%s", string(jb))
			}
		}
	}

	// Chunk long messages into 4200-char pieces.
	if msg != "" && len(msg) > core.MaxMessageTextLength && o.FileInline == nil {
		var last *events.Event
		var err error
		for start := 0; start < len(msg); start += core.MaxMessageTextLength {
			end := start + core.MaxMessageTextLength
			if end > len(msg) {
				end = len(msg)
			}
			input["text"] = strings.TrimSpace(msg[start:end])
			last, err = a.Request(ctx, RequestOptions{Method: "sendMessage", Input: input, Encrypt: true})
			if err != nil {
				return last, err
			}
		}
		return last, nil
	}

	result, err := a.Request(ctx, RequestOptions{Method: "sendMessage", Input: input, Encrypt: true})
	if err != nil {
		// If the server rejected the message with INVALID_INPUT and
		// the payload contained metadata, retry once WITHOUT metadata.
		// The Rubika server is picky about metadata: a single bad
		// part (e.g. a Pre block whose language it doesn't like, or
		// overlapping parts) causes it to reject the whole message.
		// Retrying without metadata guarantees the text still goes
		// through — it just renders without formatting.
		if md, hasMD := input["metadata"]; hasMD && md != nil {
			if s := err.Error(); strings.Contains(s, "INVALID_INPUT") {
				plainInput := map[string]any{
					"object_guid": input["object_guid"],
					"rnd":         randInt63(),
					"text":        input["text"],
				}
				if rt, ok := input["reply_to_message_id"]; ok && rt != "" {
					plainInput["reply_to_message_id"] = rt
				}
				if a.debug {
					a.debugf("sendMessage retrying without metadata (INVALID_INPUT)")
				}
				result, err = a.Request(ctx, RequestOptions{Method: "sendMessage", Input: plainInput, Encrypt: true})
				if err != nil {
					return nil, err
				}
				return result, nil
			}
		}
		return nil, err
	}
	if o.AutoDelete > 0 {
		go func(obj, msgID string, d time.Duration) {
			ctx2, cancel := context.WithTimeout(context.Background(), d)
			defer cancel()
			_ = a.AutoDeleteMessage(ctx2, obj, msgID, d)
		}(result.Str("object_guid"), result.Str("message_id"), o.AutoDelete)
	}
	return result, nil
}

// SendPhoto sends a photo.
func (a *Account) SendPhoto(ctx context.Context, objectGUID string, photo any, caption string, opts ...SendOptions) (*events.Event, error) {
	return a.sendFileType(ctx, objectGUID, photo, caption, "Image", opts...)
}

// SendVideo sends a video.
func (a *Account) SendVideo(ctx context.Context, objectGUID string, video any, caption string, opts ...SendOptions) (*events.Event, error) {
	return a.sendFileType(ctx, objectGUID, video, caption, "Video", opts...)
}

// SendVoice sends a voice message.
func (a *Account) SendVoice(ctx context.Context, objectGUID string, voice any, caption string, opts ...SendOptions) (*events.Event, error) {
	return a.sendFileType(ctx, objectGUID, voice, caption, "Voice", opts...)
}

// SendMusic sends an audio file.
func (a *Account) SendMusic(ctx context.Context, objectGUID string, music any, caption string, opts ...SendOptions) (*events.Event, error) {
	return a.sendFileType(ctx, objectGUID, music, caption, "Music", opts...)
}

// SendGif sends a GIF.
func (a *Account) SendGif(ctx context.Context, objectGUID string, gif any, caption string, opts ...SendOptions) (*events.Event, error) {
	return a.sendFileType(ctx, objectGUID, gif, caption, "Gif", opts...)
}

// SendDocument sends a generic file.
func (a *Account) SendDocument(ctx context.Context, objectGUID string, file any, caption string, opts ...SendOptions) (*events.Event, error) {
	return a.sendFileType(ctx, objectGUID, file, caption, "File", opts...)
}

// SendVideoMessage sends a round video (video-message).
func (a *Account) SendVideoMessage(ctx context.Context, objectGUID string, video any, caption string, opts ...SendOptions) (*events.Event, error) {
	return a.sendFileType(ctx, objectGUID, video, caption, "VideoMessage", opts...)
}

func (a *Account) sendFileType(ctx context.Context, objectGUID string, file any, caption, fileType string, opts ...SendOptions) (*events.Event, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	var fileInline *FileInlineInput
	if upd, ok := file.(*events.Event); ok {
		fileInline = &FileInlineInput{Update: upd, Type: fileType}
	} else {
		var src files.Source
		var fileName string
		switch v := file.(type) {
		case string:
			src.Path = v
			// IMPORTANT: the server rejects `file_name` values that
			// contain path separators (e.g. "examples/mybot/hero-barber.jpg").
			// Only the basename is sent.
			fileName = filepath.Base(v)
		case []byte:
			src.Bytes = v
			fileName = "upload.bin"
		default:
			return nil, fmt.Errorf("account: unsupported file type %T", file)
		}
		uploader := files.NewUploader(a.transport.HTTPClient(), core.DefaultChunkSize, a.maxRetries, time.Second,
			map[string]string{"auth": a.auth, "user-agent": a.userAgent})
		uploadResult, err := uploader.Upload(ctx, src, "", fileName, func(ctx context.Context, name string, size int64, mime string) (string, int, string, string, error) {
			ev, err := a.RequestSendFile(ctx, name, size, mime)
			if err != nil {
				return "", 0, "", "", err
			}
			// IMPORTANT: `id` and `access_hash_send` may be returned
			// as JSON numbers (e.g. `"id": 12345`) — `ev.Str()` only
			// matches string values and would return "" for those,
			// causing the server to reject the subsequent sendMessage
			// with INVALID_INPUT. `AnyStr` coerces numbers to strings.
			fileID := ev.AnyStr("id")
			dcID := ev.Int("dc_id")
			uploadURL := ev.Str("upload_url")
			accessHashSend := ev.AnyStr("access_hash_send")
			// Debug log of the requestSendFile response.
			if a.debug {
				a.debugf("requestSendFile result: id=%q dc_id=%d upload_url=%q access_hash_send=%q",
					fileID, dcID, uploadURL, accessHashSend)
			}
			return fileID, dcID, uploadURL, accessHashSend, nil
		}, nil)
		if err != nil {
			return nil, err
		}
		// Debug log of the final upload result.
		if a.debug {
			a.debugf("upload result: file_id=%q dc_id=%d size=%d mime=%q file_name=%q access_hash_rec=%q",
				uploadResult.FileID, uploadResult.DCID, uploadResult.Size, uploadResult.MIME, uploadResult.FileName, uploadResult.AccessHashRec)
		}
		// Generate a thumbnail for image types — the server REQUIRES
		// `thumb_inline` and the REAL `width`/`height` for Image, Video,
		// Gif, VideoMessage types. Without it the server rejects the
		// request with INVALID_INPUT: sendMessage failed.
		extraFields := map[string]any{}
		needsThumb := fileType == "Image" || fileType == "Video" || fileType == "Gif" || fileType == "VideoMessage"
		if needsThumb {
			thumb, err := files.GenerateThumbnail(src)
			if err != nil {
				// Fall back to a plain File so the message can still go out.
				a.warnf("thumbnail generation failed (%v); downgrading type to File", err)
				fileType = "File"
			} else {
				extraFields["thumb_inline"] = thumb.ThumbInline
				extraFields["width"] = thumb.Width
				extraFields["height"] = thumb.Height
				a.debugf("thumbnail generated: width=%d height=%d thumb_len=%d",
					thumb.Width, thumb.Height, len(thumb.ThumbInline))
			}
		}
		// Wrap the upload result + any thumbnail fields into an *events.Event.
		upd := events.New(map[string]any{
			"mime":            uploadResult.MIME,
			"size":            uploadResult.Size,
			"dc_id":           uploadResult.DCID,
			"file_id":         uploadResult.FileID,
			"file_name":       uploadResult.FileName,
			"access_hash_rec": uploadResult.AccessHashRec,
		})
		for k, v := range extraFields {
			upd.Set(k, v)
		}
		fileInline = &FileInlineInput{Update: upd, Type: fileType}
	}
	o.FileInline = fileInline
	return a.SendMessage(ctx, objectGUID, caption, o)
}

// SendSticker sends a sticker with the full sticker payload.
func (a *Account) SendSticker(ctx context.Context, objectGUID, emojiCharacter, stickerID, stickerSetID, wHRatio string, stickerFile map[string]any, replyToMessageID string) (*events.Event, error) {
	if wHRatio == "" {
		wHRatio = "1.0"
	}
	in := map[string]any{
		"object_guid": objectGUID,
		"rnd":         randInt63(),
		"sticker": map[string]any{
			"emoji_character": emojiCharacter,
			"sticker_id":      stickerID,
			"sticker_set_id":  stickerSetID,
			"w_h_ratio":       wHRatio,
			"file":            stickerFile,
		},
	}
	if replyToMessageID != "" {
		in["reply_to_message_id"] = replyToMessageID
	}
	return a.Request(ctx, RequestOptions{Method: "sendMessage", Input: in, Encrypt: true})
}

// SendSimpleSticker sends a sticker by sticker_id only.
func (a *Account) SendSimpleSticker(ctx context.Context, objectGUID, stickerID, replyToMessageID string) (*events.Event, error) {
	in := map[string]any{
		"object_guid": objectGUID,
		"sticker_id":  stickerID,
		"rnd":         randInt63(),
	}
	if replyToMessageID != "" {
		in["reply_to_message_id"] = replyToMessageID
	}
	return a.Request(ctx, RequestOptions{Method: "sendMessage", Input: in, Encrypt: true})
}

func normalizeSelf(guid, selfGUID string) string {
	switch strings.ToLower(guid) {
	case "me", "cloud", "self":
		return selfGUID
	}
	return guid
}

func mergeMetadata(in map[string]any, md *text.Result) {
	if md == nil {
		return
	}
	if md.Text != "" {
		in["text"] = md.Text
	}
	if md.Metadata != nil {
		in["metadata"] = md.Metadata
	}
}
