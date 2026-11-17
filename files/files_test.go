package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestUploaderUploadWithServer runs the chunked upload protocol against
// a real HTTP test server that mimics the Rubika upload endpoints.
func TestUploaderUploadWithServer(t *testing.T) {
	// Test file: 3.5 chunks of 1 MiB → 4 parts.
	data := bytes.Repeat([]byte("A"), 3*1024*1024+512*1024) // 3.5 MiB
	tmpDir, err := os.MkdirTemp("", "parsrubika-upload-")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	srcPath := filepath.Join(tmpDir, "test.bin")
	if err := os.WriteFile(srcPath, data, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var chunkCount int
	var mu sync.Mutex
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		chunkCount++
		mu.Unlock()
		_, _ = io.Copy(io.Discard, r.Body)
		_, _ = fmt.Fprintf(w, `{"status":"OK","status_det":"OK","data":{"access_hash_rec":"abc"}}`)
	}))
	defer srv.Close()

	requestSendFile := func(ctx context.Context, name string, size int64, mime string) (string, int, string, string, error) {
		return "file_123", 1, srv.URL, "access_hash_send", nil
	}

	uploader := NewUploader(srv.Client(), 1024*1024, 3, time.Millisecond, map[string]string{"auth": "test"})
	src := Source{Path: srcPath}
	result, err := uploader.Upload(context.Background(), src, "bin", "test.bin", requestSendFile, func(total, current int64) {
		if current > total {
			t.Errorf("current > total: %d > %d", current, total)
		}
	})
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if result.FileID != "file_123" {
		t.Errorf("FileID = %q, want file_123", result.FileID)
	}
	if result.AccessHashRec != "abc" {
		t.Errorf("AccessHashRec = %q, want abc", result.AccessHashRec)
	}
	if result.Size != int64(len(data)) {
		t.Errorf("Size = %d, want %d", result.Size, len(data))
	}
	// Should have received 4 chunks (3.5 MiB / 1 MiB chunk = 4 parts).
	mu.Lock()
	got := chunkCount
	mu.Unlock()
	if got != 4 {
		t.Errorf("chunk count = %d, want 4", got)
	}
}

// TestUploaderUploadBytes tests the in-memory (bytes) upload path.
func TestUploaderUploadBytes(t *testing.T) {
	data := []byte("hello world from bytes")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.Copy(io.Discard, r.Body)
		_, _ = fmt.Fprintf(w, `{"status":"OK","status_det":"OK","data":{"access_hash_rec":"abc"}}`)
	}))
	defer srv.Close()

	requestSendFile := func(ctx context.Context, name string, size int64, mime string) (string, int, string, string, error) {
		return "file_456", 1, srv.URL, "ahs", nil
	}

	uploader := NewUploader(srv.Client(), 1024, 3, time.Millisecond, map[string]string{"auth": "test"})
	src := Source{Bytes: data}
	result, err := uploader.Upload(context.Background(), src, "bin", "upload.bin", requestSendFile, nil)
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if result.FileID != "file_456" {
		t.Errorf("FileID = %q, want file_456", result.FileID)
	}
	if result.Size != int64(len(data)) {
		t.Errorf("Size = %d, want %d", result.Size, len(data))
	}
}

// TestUploaderUploadRetry verifies that the uploader retries on
// ERROR_TRY_AGAIN by restarting the upload.
func TestUploaderUploadRetry(t *testing.T) {
	data := []byte("retry me")
	attempt := 0
	var attemptMu sync.Mutex
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.Copy(io.Discard, r.Body)
		attemptMu.Lock()
		attempt++
		curAttempt := attempt
		attemptMu.Unlock()
		if curAttempt == 1 {
			// First chunk: server requests restart.
			_, _ = fmt.Fprintf(w, `{"status":"ERROR_TRY_AGAIN","status_det":"ERROR_TRY_AGAIN"}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"status":"OK","status_det":"OK","data":{"access_hash_rec":"rec"}}`)
	}))
	defer srv.Close()

	var reqCount int
	var reqMu sync.Mutex
	requestSendFile := func(ctx context.Context, name string, size int64, mime string) (string, int, string, string, error) {
		reqMu.Lock()
		reqCount++
		id := reqCount
		reqMu.Unlock()
		return fmt.Sprintf("file_%d", id), 1, srv.URL, "ahs", nil
	}

	uploader := NewUploader(srv.Client(), 1024, 3, time.Millisecond, map[string]string{"auth": "test"})
	src := Source{Bytes: data}
	result, err := uploader.Upload(context.Background(), src, "bin", "retry.bin", requestSendFile, nil)
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if reqCount != 2 {
		t.Errorf("requestSendFile should have been called twice, got %d", reqCount)
	}
	if result.AccessHashRec != "rec" {
		t.Errorf("AccessHashRec = %q, want rec", result.AccessHashRec)
	}
}
