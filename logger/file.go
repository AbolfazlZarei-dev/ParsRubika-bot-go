package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// RotatingFile is an io.WriteCloser that appends to a file and rotates
// it once it grows past a size limit, keeping a bounded number of
// numbered backups (app.log.1, app.log.2, ...).
type RotatingFile struct {
	mu         sync.Mutex
	path       string
	maxBytes   int64
	maxBackups int
	f          *os.File
	size       int64
}

// NewRotatingFile opens path for appending, creating directories as
// needed. maxSizeMB <= 0 disables rotation.
func NewRotatingFile(path string, maxSizeMB, maxBackups int) (*RotatingFile, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	r := &RotatingFile{path: path, maxBytes: int64(maxSizeMB) << 20, maxBackups: maxBackups}
	if err := r.open(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RotatingFile) open() error {
	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	st, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return err
	}
	r.f, r.size = f, st.Size()
	return nil
}

// Write implements io.Writer.
func (r *RotatingFile) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.f == nil {
		return 0, os.ErrClosed
	}
	if r.maxBytes > 0 && r.size+int64(len(p)) > r.maxBytes {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := r.f.Write(p)
	r.size += int64(n)
	return n, err
}

func (r *RotatingFile) rotate() error {
	_ = r.f.Close()
	r.f = nil
	if r.maxBackups <= 0 {
		_ = os.Remove(r.path)
	} else {
		_ = os.Remove(fmt.Sprintf("%s.%d", r.path, r.maxBackups))
		for i := r.maxBackups - 1; i >= 1; i-- {
			_ = os.Rename(fmt.Sprintf("%s.%d", r.path, i), fmt.Sprintf("%s.%d", r.path, i+1))
		}
		if err := os.Rename(r.path, r.path+".1"); err != nil {
			_ = os.Remove(r.path)
		}
	}
	return r.open()
}

// Close closes the underlying file.
func (r *RotatingFile) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.f == nil {
		return nil
	}
	err := r.f.Close()
	r.f = nil
	return err
}

// Rotate forces an immediate rotation.
func (r *RotatingFile) Rotate() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.f == nil {
		return os.ErrClosed
	}
	return r.rotate()
}
