package account

import (
	"context"
	"fmt"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/files"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/netutil"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// File transfer: raw upload/download and convenience wrappers.

// UploadFromPath uploads a local file by path and returns the upload
// result event. The MIME type is auto-detected from the file extension.
func (a *Account) UploadFromPath(ctx context.Context, path string) (*UploadResult, error) {
	mime := mimeFromExt(filepath.Ext(path))
	src := files.Source{Path: path}
	ev, err := a.Upload(ctx, src, mime, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	return &UploadResult{
		Event: ev,
		MIME:  mime,
		Name:  filepath.Base(path),
	}, nil
}

// UploadFromBytes uploads raw bytes as a file.
func (a *Account) UploadFromBytes(ctx context.Context, data []byte, fileName, mime string) (*UploadResult, error) {
	if mime == "" {
		mime = mimeFromExt(filepath.Ext(fileName))
	}
	src := files.Source{Bytes: data, Size: int64(len(data))}
	ev, err := a.Upload(ctx, src, mime, fileName)
	if err != nil {
		return nil, err
	}
	return &UploadResult{
		Event: ev,
		MIME:  mime,
		Name:  fileName,
	}, nil
}

// UploadResult wraps an upload response with metadata.
type UploadResult struct {
	Event interface{ AnyStr(keys ...string) string }
	MIME  string
	Name  string
}

// DownloadToPath downloads a file inline to the given path.
func (a *Account) DownloadToPath(ctx context.Context, fileInline *events.Event, saveAs string) error {
	if fileInline == nil {
		return fmt.Errorf("no file to download")
	}
	_, err := a.Download(ctx, fileInline, saveAs, 0, nil)
	return err
}

// DownloadToBytes downloads a file inline and returns the raw bytes.
func (a *Account) DownloadToBytes(ctx context.Context, fileInline *events.Event) ([]byte, error) {
	if fileInline == nil {
		return nil, fmt.Errorf("no file to download")
	}
	result, err := a.Download(ctx, fileInline, "", 0, nil)
	if err != nil {
		return nil, err
	}
	return result.Bytes, nil
}

// UploadFromURL downloads a file from a URL and uploads it to Rubika.
// This is useful for forwarding media from other sources.
//
// Security: the URL is validated and resolved through the netutil SSRF
// policy (http/https only; private, loopback, link-local and cloud
// metadata IPs are refused both at pre-flight and at the dial layer),
// and the response body is capped at netutil.MaxFetchBytes so a
// malicious remote cannot exhaust memory.
func (a *Account) UploadFromURL(ctx context.Context, url, fileName string) (*UploadResult, error) {
	if _, err := netutil.AssertURL(ctx, url); err != nil {
		return nil, err
	}
	client := netutil.NewSafeClient(a.timeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ParsRubika/4.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("download %s: HTTP %d", url, resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, netutil.MaxFetchBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > netutil.MaxFetchBytes {
		return nil, fmt.Errorf("download %s: response exceeds %d bytes", url, netutil.MaxFetchBytes)
	}
	if fileName == "" {
		// Try to extract from Content-Disposition
		if cd := resp.Header.Get("Content-Disposition"); cd != "" {
			fileName = extractFileName(cd)
		}
		if fileName == "" {
			fileName = "download"
		}
	}
	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = mimeFromExt(filepath.Ext(fileName))
	}
	return a.UploadFromBytes(ctx, data, fileName, mime)
}

// mimeFromExt returns a MIME type for common file extensions.
func mimeFromExt(ext string) string {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	switch ext {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	case "mp4":
		return "video/mp4"
	case "webm":
		return "video/webm"
	case "mov":
		return "video/quicktime"
	case "avi":
		return "video/x-msvideo"
	case "mkv":
		return "video/x-matroska"
	case "mp3":
		return "audio/mpeg"
	case "ogg":
		return "audio/ogg"
	case "wav":
		return "audio/wav"
	case "m4a":
		return "audio/mp4"
	case "flac":
		return "audio/flac"
	case "pdf":
		return "application/pdf"
	case "zip":
		return "application/zip"
	case "rar":
		return "application/vnd.rar"
	case "7z":
		return "application/x-7z-compressed"
	case "tar":
		return "application/x-tar"
	case "gz":
		return "application/gzip"
	case "txt":
		return "text/plain"
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "js":
		return "application/javascript"
	case "json":
		return "application/json"
	case "xml":
		return "application/xml"
	case "doc":
		return "application/msword"
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "xls":
		return "application/vnd.ms-excel"
	case "xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "ppt":
		return "application/vnd.ms-powerpoint"
	case "pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	default:
		return "application/octet-stream"
	}
}

// extractFileName parses a Content-Disposition header for the filename
// and returns a sanitized base name. Path separators (both "/" and "\")
// are stripped so a header like `filename="../../etc/passwd"` cannot
// inject a path component, and control characters are dropped. The
// result is therefore safe to use as a file_name sent to the server or
// as a local file base name.
func extractFileName(cd string) string {
	// Look for filename="..." or filename=...
	idx := strings.Index(strings.ToLower(cd), "filename=")
	if idx < 0 {
		return ""
	}
	rest := cd[idx+9:]
	rest = strings.TrimSpace(rest)
	if strings.HasPrefix(rest, "\"") {
		rest = rest[1:]
		if end := strings.Index(rest, "\""); end >= 0 {
			return sanitizeFileName(rest[:end])
		}
	}
	if end := strings.IndexAny(rest, "; "); end >= 0 {
		rest = rest[:end]
	}
	return sanitizeFileName(rest)
}

// sanitizeFileName keeps only the base-name portion of `name` (after the
// last "/" or "\") and strips control characters. It returns "" for
// names that would resolve to a directory (".", "..") or be empty.
func sanitizeFileName(name string) string {
	if i := strings.LastIndexAny(name, `/\`); i >= 0 {
		name = name[i+1:]
	}
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == ".." {
		return ""
	}
	name = strings.Map(func(r rune) rune {
		if r < 0x20 {
			return -1
		}
		return r
	}, name)
	return name
}

// FileExists reports whether a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileSize returns the size of a file in bytes, or 0 on error.
func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// DownloadProfilePicture downloads the avatar of an object. If saveAs is
// empty the raw bytes are returned; otherwise the file path is returned.
func (a *Account) DownloadProfilePicture(ctx context.Context, objectGUID string, saveAs string) (any, error) {
	resp, err := a.GetAvatars(ctx, objectGUID)
	if err != nil {
		return nil, err
	}
	avatars := resp.Slice("avatars")
	if len(avatars) == 0 {
		return nil, fmt.Errorf("no avatars for %s", objectGUID)
	}
	av, _ := avatars[0].(map[string]any)
	if av == nil {
		return nil, fmt.Errorf("invalid avatar response")
	}
	fileInline, _ := av["file_inline"].(map[string]any)
	if fileInline == nil {
		return nil, fmt.Errorf("avatar has no file_inline")
	}
	return a.Download(ctx, events.New(fileInline), saveAs, 0, nil)
}

// RequestSendFile asks the server to allocate an upload slot.
func (a *Account) RequestSendFile(ctx context.Context, fileName string, size int64, mime string) (*events.Event, error) {
	return a.Request(ctx, RequestOptions{
		Method:     "requestSendFile",
		Input:      map[string]any{"file_name": fileName, "size": size, "mime": mime},
		TmpSession: a.auth == "",
		Encrypt:    true,
	})
}

// Upload uploads a file and returns the metadata wrapped in an
// *events.Event.
func (a *Account) Upload(ctx context.Context, src files.Source, mime, fileName string) (*events.Event, error) {
	uploader := files.NewUploader(a.transport.HTTPClient(), core.DefaultChunkSize, a.maxRetries, time.Second,
		map[string]string{"auth": a.auth, "user-agent": a.userAgent})
	result, err := uploader.Upload(ctx, src, mime, fileName, func(ctx context.Context, name string, size int64, mime string) (string, int, string, string, error) {
		ev, err := a.RequestSendFile(ctx, name, size, mime)
		if err != nil {
			return "", 0, "", "", err
		}
		return ev.AnyStr("id"), ev.Int("dc_id"), ev.AnyStr("upload_url"), ev.AnyStr("access_hash_send"), nil
	}, nil)
	if err != nil {
		return nil, err
	}
	return events.New(map[string]any{
		"mime":            result.MIME,
		"size":            result.Size,
		"dc_id":           result.DCID,
		"file_id":         result.FileID,
		"file_name":       result.FileName,
		"access_hash_rec": result.AccessHashRec,
	}), nil
}

// Download downloads a file by its file_inline metadata. When `saveAs`
// is empty the raw bytes are returned.
func (a *Account) Download(ctx context.Context, fileInline *events.Event, saveAs string, chunkSize int64, cb files.ProgressCallback) (*files.Result, error) {
	if fileInline == nil {
		return nil, fmt.Errorf("account: file_inline is nil")
	}
	if chunkSize <= 0 {
		chunkSize = core.DefaultChunkSize
	}
	d := files.NewDownloader(a.transport.HTTPClient(), chunkSize, a.maxRetries, time.Second,
		map[string]string{"auth": a.auth, "user-agent": a.userAgent})
	res, err := d.Download(ctx,
		fileInline.Int("dc_id"),
		fileInline.Str("file_id"),
		fileInline.Str("access_hash_rec"),
		int64(fileInline.Int("size")),
		saveAs, false, cb)
	if err != nil {
		return nil, err
	}
	if saveAs != "" && filepath.Ext(saveAs) == "" {
		// Derive a SAFE extension from the server-supplied MIME. The
		// previous code appended fileInline.Str("mime") verbatim, which
		// (a) produced names like "file.image/jpeg" and (b) let a
		// malicious server set mime="../../etc/x" for a path-traversal
		// rename. safeExtFromMIME returns only a sanitized alphanumeric
		// extension (or ""), so the rename stays inside the intended
		// directory.
		if ext := safeExtFromMIME(fileInline.Str("mime")); ext != "" {
			newName := saveAs + "." + ext
			if err := os.Rename(res.Path, newName); err == nil {
				res.Path = newName
			}
		}
	}
	return res, nil
}

// safeExtFromMIME derives a sanitized file extension from a MIME type
// string (e.g. "image/jpeg" -> "jpg"). It returns "" for anything that
// is not a clean alphanumeric subtype, so it can never introduce a path
// separator, a ".." sequence, or any other character that could escape
// the intended directory during a rename.
func safeExtFromMIME(mime string) string {
	if mime == "" {
		return ""
	}
	sub := mime
	if i := strings.IndexByte(mime, '/'); i >= 0 {
		sub = mime[i+1:]
	}
	// Strip parameters like "; charset=utf-8".
	if i := strings.IndexByte(sub, ';'); i >= 0 {
		sub = sub[:i]
	}
	sub = strings.TrimSpace(sub)
	// Map common subtypes to their conventional file extensions.
	switch strings.ToLower(sub) {
	case "jpeg":
		return "jpg"
	case "mpeg":
		return "mp3"
	case "quicktime":
		return "mov"
	case "x-msvideo":
		return "avi"
	case "x-matroska":
		return "mkv"
	case "svg+xml":
		return "svg"
	}
	// Reject anything that is not pure alphanumeric — this is the guard
	// against path traversal / weird characters from a server-controlled
	// field.
	for _, r := range sub {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return ""
		}
	}
	if sub == "" {
		return ""
	}
	return strings.ToLower(sub)
}
