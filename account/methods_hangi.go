package account

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// The "hangi" (hanging / dual) profile: پروفایل هنگی. Two different
// pictures live on ONE object at the same time because each uploadAvatar
// bind carries a different main_file_id and thumbnail_file_id, and the
// two binds land at (almost) the same instant so the server keeps both
// avatars instead of settling on one.
//
// Everything the example used to do by hand lives here now: the
// tolerant upload (the upload service answers with shapes the stock
// files.Uploader rejects), the simultaneous dual bind, and the clean-up
// so the object ends up with EXACTLY the two new photos — one
// application, not one avatar entry per run.

// errTryAgain signals the upload service's ERROR_TRY_AGAIN verdict:
// the slot must be re-requested and the chunks resent from zero.
var errTryAgain = errors.New("files: ERROR_TRY_AGAIN")

// AvatarFileRef is an image that is already in Rubika's file storage
// and is ready to be bound as an avatar.
type AvatarFileRef struct {
	Label         string // human-readable name for logs ("photo1"/"photo2")
	FileID        string
	DCID          int
	AccessHashRec string
}

// HangiOptions tunes HangiProfile. The zero value is what you want in
// almost every case: ONE application that leaves exactly the two photos.
type HangiOptions struct {
	// Rounds is how many times the simultaneous dual bind is fired.
	// The default (and the point of the helper) is 1: the server keeps
	// the pair from a single round, so more rounds only pile up avatar
	// entries instead of improving the effect.
	Rounds int
	// Offset delays the SECOND bind so the two writes collide inside
	// the same server processing window. Default 300ms.
	Offset time.Duration
	// Gap is the pause between rounds (only used when Rounds > 1).
	// Default 700ms.
	Gap time.Duration
	// Clean deletes every existing avatar before binding the pair.
	Clean bool
	// KeepOld disables the clean-up pass. By default every avatar that
	// was NOT created by this run is deleted afterwards, so rerunning
	// the helper never accumulates more photos ("only once").
	KeepOld bool
	// Progress receives one human-readable line per step. Optional.
	Progress func(msg string)
}

// HangiOption mutates HangiOptions.
type HangiOption func(*HangiOptions)

// WithHangiRounds overrides the number of dual-bind rounds (default 1).
func WithHangiRounds(n int) HangiOption {
	return func(o *HangiOptions) {
		if n > 0 {
			o.Rounds = n
		}
	}
}

// WithHangiOffset overrides the stagger before the second bind.
func WithHangiOffset(d time.Duration) HangiOption {
	return func(o *HangiOptions) { o.Offset = d }
}

// WithHangiGap overrides the pause between rounds.
func WithHangiGap(d time.Duration) HangiOption {
	return func(o *HangiOptions) { o.Gap = d }
}

// WithHangiClean deletes all existing avatars before binding the pair.
func WithHangiClean() HangiOption {
	return func(o *HangiOptions) { o.Clean = true }
}

// WithHangiKeepOld keeps avatars that this run did not create (skips
// the clean-up pass).
func WithHangiKeepOld() HangiOption {
	return func(o *HangiOptions) { o.KeepOld = true }
}

// WithHangiProgress installs a progress sink.
func WithHangiProgress(fn func(string)) HangiOption {
	return func(o *HangiOptions) {
		if fn != nil {
			o.Progress = fn
		}
	}
}

// HangiProfileResult reports what the server ended up with.
type HangiProfileResult struct {
	ObjectGUID string
	Photo1     AvatarFileRef
	Photo2     AvatarFileRef
	// AvatarIDs are the avatar ids created by this run.
	AvatarIDs []string
	// Deleted is how many leftover avatars the clean-up pass removed.
	Deleted int
	// FinalCount is the avatar count the server reports at the end.
	FinalCount int
}

func (o *HangiOptions) say(format string, args ...any) {
	if o.Progress != nil {
		o.Progress(fmt.Sprintf(format, args...))
	}
}

// UploadImage uploads a local image file (jpg/png/gif/webp) to Rubika's
// upload service and returns the reference needed to bind it as an
// avatar.
//
// It talks to the upload service through the same tolerant chunk parser
// the hangprofile example used: success is an access_hash_rec anywhere
// in the response — string or number, top level or nested — because the
// service has been seen to answer without the status/status_det fields
// the stock files.Uploader insists on.
func (a *Account) UploadImage(ctx context.Context, path string) (*AvatarFileRef, error) {
	mime := mimeFromExt(filepath.Ext(path))
	switch mime {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
	default:
		return nil, fmt.Errorf("upload %s: unsupported image type %q (use jpg/png/gif/webp)", filepath.Base(path), filepath.Ext(path))
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("upload %s: %w", filepath.Base(path), err)
	}
	label := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	fileID, dcID, hashRec, err := a.uploadImageBytes(ctx, label, filepath.Base(path), data, mime)
	if err != nil {
		return nil, err
	}
	return &AvatarFileRef{Label: label, FileID: fileID, DCID: dcID, AccessHashRec: hashRec}, nil
}

// uploadImageBytes asks the API for an upload slot (requestSendFile),
// pushes the bytes chunk by chunk, and returns the (file_id, dc_id,
// access_hash_rec) tuple needed to bind the file as an avatar. On the
// server's ERROR_TRY_AGAIN verdict it re-requests the slot and starts
// over, bounded by maxImageRestarts.
func (a *Account) uploadImageBytes(ctx context.Context, label, fileName string, data []byte, mime string) (string, int, string, error) {
	client := &http.Client{Timeout: 2 * time.Minute}

	const maxRestarts = 5
	for restart := 0; ; restart++ {
		if restart > maxRestarts {
			return "", 0, "", fmt.Errorf("%s: aborted after %d ERROR_TRY_AGAIN restarts", label, maxRestarts)
		}
		if restart > 0 {
			a.debugf("%s: requesting a fresh upload slot (restart %d/%d)", label, restart, maxRestarts)
		}

		req, err := a.RequestSendFile(ctx, fileName, int64(len(data)), mime)
		if err != nil {
			return "", 0, "", fmt.Errorf("%s: requestSendFile: %w", label, err)
		}
		fileID := req.AnyStr("id")
		dcID := req.Int("dc_id")
		uploadURL := req.AnyStr("upload_url")
		hashSend := req.AnyStr("access_hash_send")
		if fileID == "" || dcID == 0 || uploadURL == "" || hashSend == "" {
			return "", 0, "", fmt.Errorf("%s: requestSendFile incomplete", label)
		}

		hash, err := pushImageChunks(ctx, client, uploadURL, a.Auth(), a.UserAgent(), fileID, hashSend, data)
		if err == nil {
			return fileID, dcID, hash, nil
		}
		if !errors.Is(err, errTryAgain) {
			return "", 0, "", err
		}
	}
}

// pushImageChunks uploads data to uploadURL one chunk at a time.
// Success is reported by an access_hash_rec anywhere in the response —
// string or number, top-level or nested under data (the server has been
// seen to answer without status/status_det, which is what makes the
// stock files.Uploader fail). A final chunk without a hash is an error.
func pushImageChunks(ctx context.Context, client *http.Client, uploadURL, auth, ua, fileID, hashSend string, data []byte) (string, error) {
	const chunkSize = 1 << 20 // 1 MiB, same as the library default
	total := (int64(len(data)) + chunkSize - 1) / chunkSize
	if total < 1 {
		total = 1
	}
	for index := int64(0); index < total; index++ {
		start := index * chunkSize
		end := start + chunkSize
		if end > int64(len(data)) {
			end = int64(len(data))
		}
		parsed, status, err := postImageChunk(ctx, client, uploadURL, auth, ua,
			fileID, hashSend, total, index+1, data[start:end])
		if err != nil {
			return "", fmt.Errorf("chunk %d/%d: %w", index+1, total, err)
		}
		if hash := avatarHashRec(parsed); hash != "" {
			return hash, nil
		}
		if status == "ERROR_TRY_AGAIN" {
			return "", errTryAgain
		}
	}
	return "", errors.New("upload finished without an access_hash_rec")
}

// postImageChunk POSTs one chunk to the upload service and returns the
// parsed JSON plus a best-effort string form of the `status` field
// (string or number).
func postImageChunk(ctx context.Context, client *http.Client, uploadURL, auth, ua, fileID, hashSend string, totalParts, part int64, chunk []byte) (map[string]any, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, bytes.NewReader(chunk))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("auth", auth)
	req.Header.Set("user-agent", ua)
	req.Header.Set("file-id", fileID)
	req.Header.Set("total-part", strconv.FormatInt(totalParts, 10))
	req.Header.Set("part-number", strconv.FormatInt(part, 10))
	req.Header.Set("chunk-size", strconv.Itoa(len(chunk)))
	req.Header.Set("access-hash-send", hashSend)

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 9<<20))
	if err != nil {
		return nil, "", err
	}
	var parsed map[string]any
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	if err := dec.Decode(&parsed); err != nil {
		return nil, "", fmt.Errorf("HTTP %d, non-JSON body", resp.StatusCode)
	}
	if resp.StatusCode >= 300 {
		return parsed, "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return parsed, avatarJSONStr(parsed["status"]), nil
}

// avatarHashRec digs access_hash_rec out of a parsed chunk response:
// top level or nested under data.
func avatarHashRec(m map[string]any) string {
	if m == nil {
		return ""
	}
	if s := avatarJSONStr(m["access_hash_rec"]); s != "" {
		return s
	}
	if d, ok := m["data"].(map[string]any); ok {
		return avatarHashRec(d)
	}
	return ""
}

// avatarJSONStr coerces a decoded JSON value to a string. Numbers
// arrive as json.Number (dec.UseNumber) and are rendered exactly.
func avatarJSONStr(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case json.Number:
		return string(x)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	}
	return ""
}

// HangiProfile sets the dual (hangi) profile — پروفایل هنگی — on
// objectGUID with two local image files. Pass an empty objectGUID to
// target the signed-in account itself; pass a group/channel GUID to
// hang the pair on that chat instead.
//
// What it does, in order:
//  1. optionally deletes every existing avatar (WithHangiClean);
//  2. uploads both images (tolerant upload, see UploadImage);
//  3. fires ONE round of two simultaneous uploadAvatar binds — photo1
//     main + photo2 thumb, and the swapped pair — so the object ends up
//     with two different pictures (offset staggers the second bind);
//  4. deletes every avatar this run did NOT create, so rerunning the
//     helper never piles up entries: the object ends with exactly the
//     two new photos. Disable with WithHangiKeepOld.
//
// The default is a single application: Rounds=1. More rounds only add
// avatar entries server-side; the hang comes from the two writes
// colliding, not from repetition.
func (a *Account) HangiProfile(ctx context.Context, objectGUID, photo1Path, photo2Path string, opts ...HangiOption) (*HangiProfileResult, error) {
	o := &HangiOptions{Rounds: 1, Offset: 300 * time.Millisecond, Gap: 700 * time.Millisecond}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}

	guid := objectGUID
	if guid == "" {
		guid = a.GUID()
	}
	if guid == "" {
		return nil, errors.New("hangi: no object GUID (pass one or sign in first)")
	}

	res := &HangiProfileResult{ObjectGUID: guid}

	// 1. Optional clean slate.
	if o.Clean {
		o.say("deleting all existing avatars of %s...", guid)
		if err := a.DeleteAllAvatars(ctx, guid); err != nil {
			return nil, fmt.Errorf("hangi: delete avatars: %w", err)
		}
	}

	// 2. Upload both photos first: uploading does not touch the avatar
	// state, so only the bind step below is raced.
	refs := make([]*AvatarFileRef, 2)
	for i, p := range []string{photo1Path, photo2Path} {
		ref, err := a.UploadImage(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("hangi: photo%d: %w", i+1, err)
		}
		ref.Label = fmt.Sprintf("photo%d", i+1)
		refs[i] = ref
		o.say("%s uploaded (file_id=%s dc=%d)", ref.Label, ref.FileID, ref.DCID)
	}
	res.Photo1, res.Photo2 = *refs[0], *refs[1]

	// 3. ONE application: fire both uploadAvatar binds at (almost) the
	// same instant — never one after the other. Bind 1 makes photo1 the
	// main photo with photo2 as its thumbnail; bind 2 fires the swapped
	// pair, so whichever avatar a viewer gets, the object shows two
	// different pictures.
	o.say("hanging the profile: %d round(s) of simultaneous dual-binds (offset=%s)", o.Rounds, o.Offset)
	type bindResult struct {
		avatarID string
		err      error
	}
	for r := 1; r <= o.Rounds; r++ {
		if r > 1 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(o.Gap):
			}
		}
		results := make([]bindResult, 2)
		var wg sync.WaitGroup
		wg.Add(2)
		for i := 0; i < 2; i++ {
			i := i
			go func() {
				defer wg.Done()
				if i == 1 && o.Offset > 0 {
					time.Sleep(o.Offset)
				}
				resp, err := a.UploadAvatar(ctx, guid, refs[i].FileID, refs[1-i].FileID)
				results[i].err = err
				if err != nil || resp == nil {
					return
				}
				av, _ := resp.Get("avatar")
				if avm, ok := av.(map[string]any); ok {
					results[i].avatarID, _ = avm["avatar_id"].(string)
				}
				if results[i].avatarID == "" {
					results[i].avatarID = resp.AnyStr("avatar_id")
				}
			}()
		}
		wg.Wait()
		for i, br := range results {
			if br.err != nil {
				o.say("  bind%d(main=%s,thumb=%s): FAILED: %v", i+1, refs[i].Label, refs[1-i].Label, br.err)
				continue
			}
			if br.avatarID != "" {
				res.AvatarIDs = append(res.AvatarIDs, br.avatarID)
			}
			o.say("  bind%d(main=%s,thumb=%s): bound", i+1, refs[i].Label, refs[1-i].Label)
		}
	}

	// 4. Clean-up pass: remove every avatar this run did not create, so
	// the object ends with EXACTLY the pair and reruns never accumulate
	// photos. Skipped when nothing was bound (never wipe an object
	// because the binds all failed) or when KeepOld is set.
	if !o.KeepOld && len(res.AvatarIDs) > 0 {
		created := make(map[string]bool, len(res.AvatarIDs))
		for _, id := range res.AvatarIDs {
			created[id] = true
		}
		resp, err := a.GetAvatars(ctx, guid)
		if err != nil {
			o.say("cleanup: getAvatars failed: %v", err)
		} else {
			for _, av := range resp.Slice("avatars") {
				avm, ok := av.(map[string]any)
				if !ok {
					continue
				}
				id, _ := avm["avatar_id"].(string)
				if id == "" || created[id] {
					continue
				}
				if _, err := a.DeleteAvatar(ctx, guid, id); err != nil {
					o.say("cleanup: delete %s: %v", id, err)
					continue
				}
				res.Deleted++
			}
			if res.Deleted > 0 {
				o.say("cleanup: removed %d leftover avatar(s)", res.Deleted)
			}
		}
	}

	// Final state as the server sees it.
	if resp, err := a.GetAvatars(ctx, guid); err != nil {
		o.say("getAvatars failed: %v", err)
	} else {
		res.FinalCount = len(resp.Slice("avatars"))
		o.say("server lists %d avatar(s)", res.FinalCount)
	}
	return res, nil
}
