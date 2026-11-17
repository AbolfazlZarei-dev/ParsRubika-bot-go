package files

import (
	"bytes"
	"context"
	"encoding/json"
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

// maxChunkSize caps the chunk size an Uploader/Downloader will allocate.
// A misconfigured caller passing a huge chunkSize would otherwise
// allocate a giant per-chunk buffer; the cap keeps RAM bounded.
const maxChunkSize = 8 << 20 // 8 MiB

// MaxUploadSize is the client-side sanity cap on a single uploaded file.
// The Rubika server enforces its own limit; this guard only rejects
// absurd sizes before any byte is read.
const MaxUploadSize = 2 << 30 // 2 GiB

// maxChunkResponse caps how many bytes of a single chunk HTTP response
// are read into memory, so a malicious server cannot exhaust RAM by
// streaming an oversized chunk body.
const maxChunkResponse = maxChunkSize + 1<<20 // 9 MiB

// maxGatherInMemory bounds the file size for which the gather (parallel,
// in-RAM) download path is allowed. Larger files must use `saveAs`
// (streamed to disk) so a big download cannot buffer the whole file in
// memory.
const maxGatherInMemory = 50 << 20 // 50 MiB

// Source is the union of an on-disk path or an in-memory blob for
// upload. Exactly one of Path / Bytes must be set.
type Source struct {
	Path  string
	Bytes []byte
	Size  int64 // populated internally when reading from Path
}

// Uploader uploads a single file using the chunked protocol. It is
// safe for concurrent use after construction (each upload runs on its
// own context).
type Uploader struct {
	httpClient *http.Client
	chunkSize  int64
	maxRetries int
	backoff    time.Duration
	headers    map[string]string
}

// NewUploader constructs an Uploader.
func NewUploader(client *http.Client, chunkSize int64, maxRetries int, backoff time.Duration, headers map[string]string) *Uploader {
	if chunkSize <= 0 {
		chunkSize = 1 << 20
	}
	if chunkSize > maxChunkSize {
		chunkSize = maxChunkSize
	}
	if maxRetries <= 0 {
		maxRetries = 3
	}
	if backoff <= 0 {
		backoff = time.Second
	}
	return &Uploader{
		httpClient: client,
		chunkSize:  chunkSize,
		maxRetries: maxRetries,
		backoff:    backoff,
		headers:    headers,
	}
}

// UploadResult is the metadata returned by the server after a
// successful upload.
type UploadResult struct {
	MIME          string `json:"mime"`
	Size          int64  `json:"size"`
	DCID          int    `json:"dc_id"`
	FileID        string `json:"file_id"`
	FileName      string `json:"file_name"`
	AccessHashRec string `json:"access_hash_rec"`
}

// RequestSendFile is the signature of the function used to obtain the
// initial (file_id, dc_id, upload_url, access_hash_send) tuple.
type RequestSendFile func(ctx context.Context, fileName string, size int64, mime string) (fileID string, dcID int, uploadURL string, accessHashSend string, err error)

// Upload uploads `src` to the server. It calls `requestSendFile` to
// obtain the upload metadata, then chunks the file and uploads each
// chunk to `uploadURL`. When the server returns `ERROR_TRY_AGAIN`, the
// upload is restarted.
//
// `cb` (if non-nil) is called after every chunk with the total size
// and the uploaded-so-far byte count.
func (u *Uploader) Upload(ctx context.Context, src Source, mime, fileName string, requestSendFile RequestSendFile, cb ProgressCallback) (*UploadResult, error) {
	// 1) Resolve size + open the source. Exactly one of Path / Bytes
	// must be set; accepting both (and silently preferring Path) would
	// mask caller bugs, so it is an error.
	if src.Path != "" && src.Bytes != nil {
		return nil, fmt.Errorf("files: source must be a path OR bytes, not both")
	}
	if src.Path != "" {
		stat, err := os.Stat(src.Path)
		if err != nil {
			return nil, fmt.Errorf("files: stat %s: %w", src.Path, err)
		}
		if stat.Size() > MaxUploadSize {
			return nil, fmt.Errorf("files: file size %d exceeds MaxUploadSize (%d)", stat.Size(), MaxUploadSize)
		}
		if fileName == "" {
			// Use only the base name so the local filesystem path is
			// never leaked to the server as file_name.
			fileName = filepath.Base(src.Path)
		} else {
			fileName = filepath.Base(fileName)
		}
		src.Size = stat.Size()
	} else if src.Bytes != nil {
		if int64(len(src.Bytes)) > MaxUploadSize {
			return nil, fmt.Errorf("files: bytes size %d exceeds MaxUploadSize (%d)", len(src.Bytes), MaxUploadSize)
		}
		if fileName == "" {
			return nil, fmt.Errorf("files: fileName is required when uploading from bytes")
		}
		fileName = filepath.Base(fileName)
		src.Size = int64(len(src.Bytes))
	} else {
		return nil, fmt.Errorf("files: source must be a path or bytes")
	}
	if mime == "" {
		mime = mimeFromName(fileName)
	}

	// 2) Request the upload metadata.
	fileID, dcID, uploadURL, accessHashSend, err := requestSendFile(ctx, fileName, src.Size, mime)
	if err != nil {
		return nil, fmt.Errorf("files: requestSendFile: %w", err)
	}

	totalParts := (src.Size + u.chunkSize - 1) / u.chunkSize
	index := int64(0)
	var lastResult map[string]any

	// Open the file ONCE for the whole upload and seek per chunk, so a
	// 1 GB upload doesn't trigger 1024 open/close syscalls.
	var fh *os.File
	if src.Path != "" {
		fh, err = os.Open(src.Path)
		if err != nil {
			return nil, err
		}
		defer fh.Close()
	}

	// maxRestarts bounds the number of ERROR_TRY_AGAIN recoveries so a
	// persistently-failing server can't trap the client in an infinite
	// re-upload loop.
	const maxRestarts = 5
	restarts := 0

	// Reusable chunk buffer (avoids a per-chunk allocation). Only one
	// upload goroutine touches it at a time, so no pool is needed.
	chunkBuf := make([]byte, u.chunkSize)

	for index < totalParts {
		// Read the chunk.
		var data []byte
		if src.Path != "" {
			if _, err := fh.Seek(index*u.chunkSize, io.SeekStart); err != nil {
				return nil, err
			}
			nr, err := io.ReadFull(fh, chunkBuf)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return nil, err
			}
			data = chunkBuf[:nr]
		} else {
			start := index * u.chunkSize
			end := start + u.chunkSize
			if end > int64(len(src.Bytes)) {
				end = int64(len(src.Bytes))
			}
			data = src.Bytes[start:end]
		}

		ok := false
		for attempt := 0; attempt < u.maxRetries; attempt++ {
			hdr := map[string]string{
				"auth":             u.headers["auth"],
				"file-id":          fileID,
				"total-part":       strconv.FormatInt(totalParts, 10),
				"part-number":      strconv.FormatInt(index+1, 10),
				"chunk-size":       strconv.Itoa(len(data)),
				"access-hash-send": accessHashSend,
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			for k, v := range hdr {
				req.Header.Set(k, v)
			}
			resp, err := u.httpClient.Do(req)
			if err != nil {
				if attempt < u.maxRetries-1 {
					sleepBackoff(ctx, u.backoff, attempt)
					continue
				}
				return nil, fmt.Errorf("files: upload chunk %d: %w", index+1, err)
			}
			// Cap the response so a malicious server cannot exhaust RAM
			// by streaming an oversized chunk-status body.
			body, _ := io.ReadAll(io.LimitReader(resp.Body, maxChunkResponse))
			_ = resp.Body.Close()
			var parsed map[string]any
			if err := json.Unmarshal(body, &parsed); err != nil {
				if attempt < u.maxRetries-1 {
					sleepBackoff(ctx, u.backoff, attempt)
					continue
				}
				return nil, fmt.Errorf("files: upload chunk %d invalid json: %w", index+1, err)
			}
			lastResult = parsed
			if s, _ := parsed["status"].(string); s == "ERROR_TRY_AGAIN" {
				// Restart the upload: re-request the file slot and
				// continue from chunk 0. Bound the number of
				// restarts so a persistently-failing server can't
				// trap us in an infinite loop.
				restarts++
				if restarts > maxRestarts {
					return nil, fmt.Errorf("files: upload aborted after %d ERROR_TRY_AGAIN restarts", maxRestarts)
				}
				fileID, dcID, uploadURL, accessHashSend, err = requestSendFile(ctx, fileName, src.Size, mime)
				if err != nil {
					return nil, fmt.Errorf("files: restart: %w", err)
				}
				index = 0
				ok = true
				goto nextIter
			}
			ok = true
			break
		}
		if !ok {
			continue
		}
		if cb != nil {
			// Report progress based on the chunk just uploaded; do
			// not regress on restart (index was already reset to 0
			// above, so (index+1) is correct here).
			cb(src.Size, minInt64((index+1)*u.chunkSize, src.Size))
		}
		index++
	nextIter:
	}

	status, _ := lastResult["status"].(string)
	statusDet, _ := lastResult["status_det"].(string)
	// The upload service has been observed answering with the hash but
	// no string `status`/`status_det` fields (or numeric ones); any
	// final response that yields an access_hash_rec is a success.
	accessHashRec := findHashRec(lastResult)
	if (status == "OK" && statusDet == "OK") || accessHashRec != "" {
		return &UploadResult{
			MIME:          mime,
			Size:          src.Size,
			DCID:          dcID,
			FileID:        fileID,
			FileName:      fileName,
			AccessHashRec: accessHashRec,
		}, nil
	}
	raw, _ := json.Marshal(lastResult)
	return nil, fmt.Errorf("files: upload failed: status=%q status_det=%q (raw response: %.400s)", status, statusDet, raw)
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// anyToString coerces a JSON-decoded value to its string form.
// Returns "" for nil / unsupported types. Used for fields like
// `access_hash_rec` that the server may return as either a string
// or a number.
func anyToString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		if x {
			return "true"
		}
		return "false"
	}
	return ""
}

func sleepBackoff(ctx context.Context, base time.Duration, attempt int) {
	delay := base * (1 << uint(attempt))
	select {
	case <-time.After(delay):
	case <-ctx.Done():
	}
}

// fileExtension returns the file extension (without the dot) from a
// filename. Used as a fallback when the caller does not supply a MIME
// type.
func fileExtension(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i+1:]
		}
		if name[i] == '/' || name[i] == '\\' {
			break
		}
	}
	return ""
}

// mimeFromExt returns the standard MIME type for common file extensions.
// This is used by Upload when the caller does not supply a MIME type —
// the server requires a full MIME (e.g. "image/png") to render media
// correctly; passing just the extension ("png") causes images to be
// shown as generic files.
var mimeMap = map[string]string{
	// Images
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"bmp":  "image/bmp",
	"svg":  "image/svg+xml",
	// Videos
	"mp4":  "video/mp4",
	"webm": "video/webm",
	"mov":  "video/quicktime",
	"avi":  "video/x-msvideo",
	"mkv":  "video/x-matroska",
	"m4v":  "video/x-m4v",
	// Audio
	"mp3":  "audio/mpeg",
	"ogg":  "audio/ogg",
	"wav":  "audio/wav",
	"m4a":  "audio/mp4",
	"flac": "audio/flac",
	"aac":  "audio/aac",
	"opus": "audio/opus",
	// Documents
	"pdf":  "application/pdf",
	"zip":  "application/zip",
	"rar":  "application/vnd.rar",
	"7z":   "application/x-7z-compressed",
	"tar":  "application/x-tar",
	"gz":   "application/gzip",
	"txt":  "text/plain",
	"html": "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"json": "application/json",
	"xml":  "application/xml",
	"doc":  "application/msword",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"xls":  "application/vnd.ms-excel",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"ppt":  "application/vnd.ms-powerpoint",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
}

// mimeFromName returns the MIME type for the given filename based on
// its extension. Returns "application/octet-stream" for unknown
// extensions.
func mimeFromName(name string) string {
	ext := strings.ToLower(fileExtension(name))
	if m, ok := mimeMap[ext]; ok {
		return m
	}
	return "application/octet-stream"
}

// MaxParallelDownloads bounds the number of concurrent chunk
// downloads when `gather` is enabled.
const MaxParallelDownloads = 8

// Result is the union type returned by Downloader.Download. Exactly
// one of Path / Bytes will be populated depending on `saveAs`.
type Result struct {
	Path  string
	Bytes []byte
}

// Downloader fetches a file by (dc_id, file_id, access_hash) from the
// Rubika file servers using chunked HTTP requests.
type Downloader struct {
	httpClient *http.Client
	chunkSize  int64
	maxRetries int
	backoff    time.Duration
	headers    map[string]string
}

// NewDownloader constructs a Downloader.
func NewDownloader(client *http.Client, chunkSize int64, maxRetries int, backoff time.Duration, headers map[string]string) *Downloader {
	if chunkSize <= 0 {
		chunkSize = 128 * 1024
	}
	if chunkSize > maxChunkSize {
		chunkSize = maxChunkSize
	}
	if maxRetries <= 0 {
		maxRetries = 3
	}
	if backoff <= 0 {
		backoff = time.Second
	}
	return &Downloader{
		httpClient: client,
		chunkSize:  chunkSize,
		maxRetries: maxRetries,
		backoff:    backoff,
		headers:    headers,
	}
}

// Download fetches the file. When `saveAs` is non-empty the bytes are
// written to disk and the Result.Path is populated. When `gather` is
// true and `saveAs` is empty, chunks are fetched in parallel.
func (d *Downloader) Download(ctx context.Context, dcID int, fileID, accessHash string, size int64, saveAs string, gather bool, cb ProgressCallback) (*Result, error) {
	baseURL := fmt.Sprintf("https://messenger%d.iranlms.ir", dcID)
	baseHeaders := map[string]string{
		"auth":            d.headers["auth"],
		"access-hash-rec": accessHash,
		"file-id":         fileID,
		"user-agent":      d.headers["user-agent"],
	}

	fetchChunk := func(start, end int64) ([]byte, error) {
		expected := end - start + 1
		var lastErr error
		for attempt := 0; attempt < d.maxRetries; attempt++ {
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/GetFile.ashx", nil)
			if err != nil {
				return nil, err
			}
			for k, v := range baseHeaders {
				req.Header.Set(k, v)
			}
			req.Header.Set("start-index", strconv.FormatInt(start, 10))
			req.Header.Set("last-index", strconv.FormatInt(end, 10))
			resp, err := d.httpClient.Do(req)
			if err != nil {
				lastErr = err
				sleepBackoff(ctx, d.backoff, attempt)
				continue
			}
			// Cap the body so a malicious server cannot exhaust RAM by
			// streaming an oversized chunk.
			body, err := io.ReadAll(io.LimitReader(resp.Body, maxChunkResponse))
			_ = resp.Body.Close()
			if err != nil {
				lastErr = err
				sleepBackoff(ctx, d.backoff, attempt)
				continue
			}
			if resp.StatusCode != 200 {
				lastErr = fmt.Errorf("download chunk %d-%d: http %d", start, end, resp.StatusCode)
				sleepBackoff(ctx, d.backoff, attempt)
				continue
			}
			// Verify the chunk is exactly the requested range. A short
			// or truncated chunk would otherwise be silently written,
			// producing a corrupt file.
			if int64(len(body)) != expected {
				lastErr = fmt.Errorf("download chunk %d-%d: short body (%d bytes, want %d)", start, end, len(body), expected)
				sleepBackoff(ctx, d.backoff, attempt)
				continue
			}
			return body, nil
		}
		return nil, lastErr
	}

	if saveAs != "" {
		// Write to a temp file in the same directory, then rename atomically
		// onto saveAs on success. On any error the temp file is removed so
		// no partial/corrupt file is left at the final path. Mode 0600
		// keeps downloaded media readable only by the owning user.
		dir := filepath.Dir(saveAs)
		if dir == "" {
			dir = "."
		}
		tmp, err := os.CreateTemp(dir, filepath.Base(saveAs)+".download.*")
		if err != nil {
			return nil, fmt.Errorf("files: create download temp: %w", err)
		}
		tmpName := tmp.Name()
		success := false
		defer func() {
			if !success {
				_ = tmp.Close()
				_ = os.Remove(tmpName)
			}
		}()
		for start := int64(0); start < size; start += d.chunkSize {
			end := start + d.chunkSize - 1
			if end >= size {
				end = size - 1
			}
			data, err := fetchChunk(start, end)
			if err != nil {
				return nil, err
			}
			if _, err := tmp.Write(data); err != nil {
				return nil, err
			}
			if cb != nil {
				cb(size, end+1)
			}
		}
		if err := tmp.Sync(); err != nil {
			return nil, fmt.Errorf("files: sync download: %w", err)
		}
		if err := tmp.Close(); err != nil {
			return nil, fmt.Errorf("files: close download temp: %w", err)
		}
		// Final integrity check: the written file must match the expected
		// size, otherwise a truncated download is rejected rather than
		// silently accepted as a corrupt file.
		if fi, err := os.Stat(tmpName); err != nil {
			return nil, fmt.Errorf("files: stat download: %w", err)
		} else if fi.Size() != size {
			return nil, fmt.Errorf("files: downloaded size %d != expected %d", fi.Size(), size)
		}
		if err := os.Chmod(tmpName, 0o600); err != nil {
			return nil, fmt.Errorf("files: chmod download: %w", err)
		}
		if err := os.Rename(tmpName, saveAs); err != nil {
			return nil, fmt.Errorf("files: finalize download: %w", err)
		}
		success = true
		return &Result{Path: saveAs}, nil
	}

	if gather {
		// The gather path buffers the whole file in memory; refuse it for
		// large files so a big download cannot exhaust RAM. Callers that
		// need a large file must pass a non-empty `saveAs` (streamed to
		// disk via the atomic temp-file path above).
		if size > maxGatherInMemory {
			return nil, fmt.Errorf("files: gather download refused for %d-byte file (max %d); use saveAs", size, maxGatherInMemory)
		}
		type chunk struct {
			start int64
			end   int64
		}
		var chunks []chunk
		for start := int64(0); start < size; start += d.chunkSize {
			end := start + d.chunkSize - 1
			if end >= size {
				end = size - 1
			}
			chunks = append(chunks, chunk{start, end})
		}
		results := make([][]byte, len(chunks))
		errs := make([]error, len(chunks))
		sem := make(chan struct{}, MaxParallelDownloads)
		var wg sync.WaitGroup
		for i, c := range chunks {
			wg.Add(1)
			sem <- struct{}{}
			go func(i int, c chunk) {
				defer wg.Done()
				defer func() { <-sem }()
				data, err := fetchChunk(c.start, c.end)
				results[i] = data
				errs[i] = err
			}(i, c)
		}
		wg.Wait()
		var buf bytes.Buffer
		for i, data := range results {
			if errs[i] != nil {
				return nil, errs[i]
			}
			buf.Write(data)
		}
		if cb != nil {
			cb(size, int64(buf.Len()))
		}
		return &Result{Bytes: buf.Bytes()}, nil
	}

	// Sequential in-memory download. Like gather, this buffers the whole
	// file in RAM, so refuse it for large files (use saveAs instead).
	if size > maxGatherInMemory {
		return nil, fmt.Errorf("files: in-memory download refused for %d-byte file (max %d); use saveAs", size, maxGatherInMemory)
	}
	var buf bytes.Buffer
	for start := int64(0); start < size; start += d.chunkSize {
		end := start + d.chunkSize - 1
		if end >= size {
			end = size - 1
		}
		data, err := fetchChunk(start, end)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		if cb != nil {
			cb(size, int64(buf.Len()))
		}
	}
	return &Result{Bytes: buf.Bytes()}, nil
}

// findHashRec digs access_hash_rec out of a parsed upload response:
// top level, under data, or under data.data — as string or number.
func findHashRec(m map[string]any) string {
	if m == nil {
		return ""
	}
	if s := anyToString(m["access_hash_rec"]); s != "" {
		return s
	}
	if d, ok := m["data"].(map[string]any); ok {
		if s := findHashRec(d); s != "" {
			return s
		}
	}
	return ""
}
