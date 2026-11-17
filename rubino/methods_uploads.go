package rubino

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/core"
)

// RequestUploadFile asks the Rubino API to allocate an upload slot.
// `fileType` is "Picture", "Video", "Audio" etc.
func (s *Rubino) RequestUploadFile(ctx context.Context, profileID, fileName string, fileSize int64, fileType string) (*Event, error) {
	return s.request(ctx, "requestUploadFile", map[string]any{
		"profile_id": s.profileOr(profileID),
		"file_name":  fileName,
		"file_size":  fileSize,
		"file_type":  fileType,
	})
}

// UploadFile uploads a local file to the supplied `uploadURL` using
// multipart/form-data. Returns the file_id assigned by the server.
//
// The Rubino upload protocol is similar to the messenger one: the
// client requests an upload URL via RequestUploadFile, then POSTs the
// file as `multipart/form-data` to the returned URL.
func (s *Rubino) UploadFile(ctx context.Context, uploadURL, fileName, filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	body, w := createStudioMultipartBody(fileName, f)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := s.transport.HTTPClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("rubino: upload http %d: %s", resp.StatusCode, string(raw))
	}
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	if d, ok := parsed["data"].(map[string]any); ok {
		if id, ok := d["file_id"].(string); ok {
			return id, nil
		}
	}
	return "", fmt.Errorf("rubino: missing file_id in upload response")
}

// UploadAndPost is a high-level helper that uploads a local file then
// publishes a post. `postType` is "Picture" or "Video".
func (s *Rubino) UploadAndPost(ctx context.Context, profileID, postType, filePath, caption string, width, height int) (*Event, error) {
	fileName := baseName(filePath)
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("rubino: stat %s: %w", filePath, err)
	}
	rubinoType := "Picture"
	if postType == "Video" {
		rubinoType = "Video"
	}
	res, err := s.RequestUploadFile(ctx, profileID, fileName, stat.Size(), rubinoType)
	if err != nil {
		return nil, err
	}
	uploadURL := res.Str("upload_url")
	if uploadURL == "" {
		return nil, fmt.Errorf("rubino: no upload_url in response")
	}
	fileID, err := s.UploadFile(ctx, uploadURL, fileName, filePath)
	if err != nil {
		return nil, err
	}
	return s.AddPost(ctx, profileID, postType, fileID, caption, fileName, width, height)
}

// baseName returns the trailing path segment of `p`.
func baseName(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[i+1:]
		}
	}
	return p
}

// createStudioMultipartBody builds a multipart/form-data body for a
// single file. The body is streamed via an io.Pipe so the upload is
// not buffered in memory.
func createStudioMultipartBody(fileName string, f io.Reader) (io.Reader, *studioMultipartWriter) {
	r, pw := io.Pipe()
	w := newStudioMultipartWriter(pw, fileName, f)
	go w.run()
	return r, w
}

// studioMultipartWriter is a tiny multipart writer that uses a
// hand-rolled boundary (so the studio package doesn't import
// mime/multipart).
type studioMultipartWriter struct {
	pw       *io.PipeWriter
	fileName string
	body     io.Reader
	boundary string
}

func newStudioMultipartWriter(pw *io.PipeWriter, fileName string, body io.Reader) *studioMultipartWriter {
	return &studioMultipartWriter{
		pw:       pw,
		fileName: fileName,
		body:     body,
		boundary: "parsrubika-rubino-" + strconv.Itoa(int(core.Version[0])),
	}
}

func (w *studioMultipartWriter) FormDataContentType() string {
	return "multipart/form-data; boundary=" + w.boundary
}

func (w *studioMultipartWriter) run() {
	defer w.pw.Close()
	// Header
	header := fmt.Sprintf("--%s\r\n", w.boundary) +
		fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=%q\r\n", w.fileName) +
		"Content-Type: application/octet-stream\r\n\r\n"
	if _, err := w.pw.Write([]byte(header)); err != nil {
		return
	}
	// Body
	if _, err := io.Copy(w.pw, w.body); err != nil {
		return
	}
	// Footer
	footer := fmt.Sprintf("\r\n--%s--\r\n", w.boundary)
	_, _ = w.pw.Write([]byte(footer))
}

// unused import guard
