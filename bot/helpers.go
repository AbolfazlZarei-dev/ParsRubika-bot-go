package bot

import (
	"context"
	"io"
	"net/http"
	"os"
)

// newGetRequest is a small helper that builds a GET request with the
// supplied context.
func newGetRequest(ctx context.Context, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

// createFile opens (or creates) a file for writing.
func createFile(path string) (*os.File, error) {
	return os.Create(path)
}

// closeCloser is a tiny helper that closes `c` and discards the error.
func closeCloser(c io.Closer) { _ = c.Close() }
