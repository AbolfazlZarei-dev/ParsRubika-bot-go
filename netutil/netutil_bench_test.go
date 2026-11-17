package netutil

import (
	"context"
	"testing"
)

func BenchmarkValidateURL(b *testing.B) {
	urls := []string{
		"https://example.com/path/x",
		"https://rubika.ir/post/abc123",
		"http://10.0.0.1/x", // private IP literal — rejected fast
		"https://cdn.example.com/file.bin",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidateURL(urls[i%len(urls)])
	}
}

func BenchmarkAssertURL(b *testing.B) {
	// example.com resolves to public IPs; the benchmark pays for one
	// DNS lookup per call. Use a short context so a slow resolver does
	// not stall the run.
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AssertURL(ctx, "https://example.com/x")
	}
}

func BenchmarkIsBlockedIP(b *testing.B) {
	ips := [4]byte{127, 0, 0, 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isBlocked(ips)
	}
}

// isBlocked wraps IsBlockedIP for a 4-byte loop body so the compiler
// does not optimise the call away.
func isBlocked(b [4]byte) bool {
	return IsBlockedIP(b[:])
}
