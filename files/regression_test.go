package files

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// chunkRT is a RoundTripper that serves slices of a fixed in-memory file
// based on the start-index / last-index headers the Downloader sends.
// It lets us exercise the chunked download logic end-to-end without
// contacting the real Rubika file servers (the base URL is hardcoded
// inside Downloader.Download, so we intercept at the transport layer).
type chunkRT struct {
	data    []byte
	shortBy int // if >0, serve a body that many bytes shorter than requested (to test the length check)
	status  int // override status code (0 => 200)
}

func (rt chunkRT) RoundTrip(req *http.Request) (*http.Response, error) {
	start, _ := strconv.ParseInt(req.Header.Get("start-index"), 10, 64)
	end, _ := strconv.ParseInt(req.Header.Get("last-index"), 10, 64)
	if start < 0 {
		start = 0
	}
	if end >= int64(len(rt.data)) {
		end = int64(len(rt.data)) - 1
	}
	if end < start {
		end = start - 1
	}
	chunk := rt.data[start : end+1]
	if rt.shortBy > 0 && len(chunk) > rt.shortBy {
		chunk = chunk[:len(chunk)-rt.shortBy]
	}
	status := rt.status
	if status == 0 {
		status = http.StatusOK
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader(chunk)),
		Header:     http.Header{},
	}, nil
}

func newDownloaderWithRT(rt http.RoundTripper, chunkSize int64) *Downloader {
	// Tiny backoff so the short-chunk retry test stays fast.
	return NewDownloader(&http.Client{Transport: rt}, chunkSize, 2, time.Millisecond, map[string]string{"auth": "x"})
}

func TestDownloadSaveAsRoundTripAndMode0600(t *testing.T) {
	data := bytes.Repeat([]byte("ParsRubika-"), 1000) // ~10000 bytes
	dir := t.TempDir()
	saveAs := filepath.Join(dir, "out.bin")
	d := newDownloaderWithRT(chunkRT{data: data}, 1024)

	res, err := d.Download(context.Background(), 1, "fid", "hash", int64(len(data)), saveAs, false, nil)
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	if res.Path != saveAs {
		t.Fatalf("Path = %q, want %q", res.Path, saveAs)
	}
	got, err := os.ReadFile(saveAs)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("downloaded content mismatch (len %d vs %d)", len(got), len(data))
	}
	fi, _ := os.Stat(saveAs)
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("downloaded file mode = %v, want 0600", fi.Mode().Perm())
	}
	// No .download temp leftover.
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.Contains(e.Name(), ".download.") {
			t.Fatalf("temp leftover present: %s", e.Name())
		}
	}
}

func TestDownloadShortChunkRejectedAndCleanedUp(t *testing.T) {
	data := bytes.Repeat([]byte{0xAB}, 4096)
	dir := t.TempDir()
	saveAs := filepath.Join(dir, "bad.bin")
	// Serve each chunk 10 bytes short -> the length check must reject it.
	d := newDownloaderWithRT(chunkRT{data: data, shortBy: 10}, 512)

	_, err := d.Download(context.Background(), 1, "fid", "hash", int64(len(data)), saveAs, false, nil)
	if err == nil {
		t.Fatal("expected an error for a truncated chunk")
	}
	// The final file must NOT exist (cleanup of the temp on error).
	if _, err := os.Stat(saveAs); err == nil {
		t.Fatal("partial file left at saveAs after a failed download")
	}
	// And no temp leftover.
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.Contains(e.Name(), ".download.") {
			t.Fatalf("temp leftover present after failure: %s", e.Name())
		}
	}
}

func TestDownloadGatherRefusesLargeFile(t *testing.T) {
	d := newDownloaderWithRT(chunkRT{}, 1024)
	_, err := d.Download(context.Background(), 1, "fid", "hash", maxGatherInMemory+1, "", true, nil)
	if err == nil || !strings.Contains(err.Error(), "saveAs") {
		t.Fatalf("expected gather refusal for large file, got %v", err)
	}
}

func TestUploadRejectsBothSources(t *testing.T) {
	u := NewUploader(&http.Client{}, 1024, 1, 0, nil)
	_, err := u.Upload(context.Background(), Source{Path: "x", Bytes: []byte("y")}, "", "", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "not both") {
		t.Fatalf("expected 'not both' error, got %v", err)
	}
}

func TestUploadFileNameBaseNotFullPath(t *testing.T) {
	u := NewUploader(&http.Client{}, 1024, 1, 0, nil)
	var seenName string
	// requestSendFile records the fileName the uploader passes along.
	rsf := func(ctx context.Context, name string, size int64, mime string) (string, int, string, string, error) {
		seenName = name
		return "fid", 1, "http://invalid.invalid/u", "hash", nil
	}
	// Use a Bytes source with an empty fileName so it falls through... no,
	// we want to test the Path branch's filepath.Base. Use a Path source
	// with a nested path and an empty fileName so Upload derives it.
	dir := t.TempDir()
	full := filepath.Join(dir, "sub", "secret.txt")
	_ = os.MkdirAll(filepath.Dir(full), 0o755)
	_ = os.WriteFile(full, []byte("hi"), 0o600)
	_, err := u.Upload(context.Background(), Source{Path: full}, "", "", rsf, nil)
	// The upload itself will fail (invalid host), but requestSendFile has
	// already been called by then, so seenName is populated.
	if err == nil {
		t.Fatal("expected upload to fail against an invalid host")
	}
	if seenName != "secret.txt" {
		t.Fatalf("fileName sent to server = %q, want base name %q", seenName, "secret.txt")
	}
}

func TestNewUploaderCapsChunkSize(t *testing.T) {
	u := NewUploader(&http.Client{}, 1<<30, 1, 0, nil)
	if u.chunkSize > maxChunkSize {
		t.Fatalf("chunkSize = %d, want <= %d (capped)", u.chunkSize, maxChunkSize)
	}
}
