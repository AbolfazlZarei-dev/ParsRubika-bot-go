package transport

import (
	"bytes"
	"sync"
)

// bufPool is a sync.Pool of *bytes.Buffer used to avoid allocating a new
// buffer for every request/response serialization. The pool is process-
// wide and shared by all Transport instances.
//
// Usage:
//
//	buf := bufPool.Get().(*bytes.Buffer)
//	defer bufPool.Put(buf)
//	buf.Reset()
//	// ... use buf ...
var bufPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

// AcquireBuffer returns a reset *bytes.Buffer from the pool. The caller
// must return it with ReleaseBuffer.
func AcquireBuffer() *bytes.Buffer {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

// ReleaseBuffer returns a buffer to the pool. The buffer is reset before
// being put back so its memory can be reused.
func ReleaseBuffer(b *bytes.Buffer) {
	if b == nil {
		return
	}
	b.Reset()
	// Cap the retained buffer at 1 MiB to avoid pinning large
	// allocations in the pool indefinitely.
	if b.Cap() > 1<<20 {
		b = bytes.NewBuffer(make([]byte, 0, 4096))
	}
	bufPool.Put(b)
}

// mapPool is a sync.Pool of map[string]any used to avoid allocating a
// fresh map for every API request. The maps are small (typically <32
// keys) so the reuse benefit is significant under high QPS.
var mapPool = sync.Pool{
	New: func() any {
		return make(map[string]any, 8)
	},
}

// AcquireMap returns a cleared map[string]any from the pool.
func AcquireMap() map[string]any {
	m := mapPool.Get().(map[string]any)
	for k := range m {
		delete(m, k)
	}
	return m
}

// ReleaseMap returns a map to the pool. The map is cleared first so its
// memory can be reused.
func ReleaseMap(m map[string]any) {
	if m == nil {
		return
	}
	for k := range m {
		delete(m, k)
	}
	mapPool.Put(m)
}
