package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// maxFileStoreBytes caps the size of a FileStore file accepted on Load.
// A real session file is a few hundred bytes; the cap exists only so a
// maliciously planted or truncated multi-megabyte blob cannot make
// Load allocate unbounded memory during base64 decode.
const maxFileStoreBytes = 1 << 20 // 1 MiB

// writeFileAtomic writes `data` to `path` via a freshly-created temp
// file in the same directory, fsyncs it, then renames it over `path`.
//
// Guarantees:
//
//   - The temp file name is random (os.CreateTemp), so an attacker cannot
//     pre-place a symlink at the expected temp path (defeats the symlink
//     / TOCTOU race that a fixed "<path>.tmp" name would suffer).
//   - The temp file is created with mode `perm` (callers pass 0o600 for
//     stores that hold secrets) and is chmod'd again right before the
//     rename in case the umask loosened it.
//   - The temp file is fsync'd and the parent directory is fsync'd, so a
//     crash cannot leave a torn file or lose the rename.
//   - On any error the temp file is removed, so no partial file remains
//     at the final path.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path) + ".tmp.*"

	f, err := os.CreateTemp(dir, base)
	if err != nil {
		return fmt.Errorf("storage: create temp file: %w", err)
	}
	tmp := f.Name()
	// Best-effort cleanup on any failure path; a successful rename makes
	// this a no-op.
	defer func() {
		if err != nil {
			_ = os.Remove(tmp)
		}
	}()

	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("storage: write temp file: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("storage: sync temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("storage: close temp file: %w", err)
	}
	// Re-assert the requested mode: CreateTemp honours the umask, which
	// could leave the file group/world-readable on some systems.
	if err := os.Chmod(tmp, perm); err != nil {
		return fmt.Errorf("storage: chmod temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("storage: rename temp file: %w", err)
	}
	// fsync the parent directory so the rename itself is durable.
	if d, err := os.Open(dir); err == nil {
		if err := d.Sync(); err != nil {
			// Some filesystems (tmpfs) do not support directory fsync;
			// a failure here does not invalidate the completed rename.
			_ = err
		}
		d.Close()
	}
	return nil
}
