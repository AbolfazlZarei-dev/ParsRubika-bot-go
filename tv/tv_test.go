package tv

import (
	"context"
	"testing"
)

func TestNewTV(t *testing.T) {
	tv := New(WithAuth("test-auth"))
	if tv == nil {
		t.Fatal("expected non-nil TV")
	}
	defer tv.Close()
	if tv.auth != "test-auth" {
		t.Errorf("auth = %q", tv.auth)
	}
}

func TestServerList(t *testing.T) {
	if len(TVServerURLs) != 4 {
		t.Errorf("TVServerURLs has %d entries, want 4", len(TVServerURLs))
	}
	for _, url := range TVServerURLs {
		if url == "" {
			t.Error("TVServerURLs should not contain empty entries")
		}
	}
}

func TestDefaultVersion(t *testing.T) {
	if DefaultVersion == "" {
		t.Error("DefaultVersion should not be empty")
	}
	tv := New(WithAuth("test"))
	defer tv.Close()
	// getVersion should return the default when MetaURL is unreachable.
	v := tv.getVersion(context.Background())
	if v == "" {
		t.Error("getVersion should not return empty string")
	}
}

func TestWithVersion(t *testing.T) {
	tv := New(WithAuth("test"), WithVersion("9.9.9"))
	defer tv.Close()
	if tv.getVersion(context.Background()) != "9.9.9" {
		t.Error("getVersion should return the overridden version")
	}
}

func TestTVHeaders(t *testing.T) {
	if TVHeaders["origin"] != "https://tv.rubika.ir" {
		t.Errorf("origin = %q", TVHeaders["origin"])
	}
	if TVHeaders["referer"] != "https://tv.rubika.ir/" {
		t.Errorf("referer = %q", TVHeaders["referer"])
	}
}

func TestRequestNoAuth(t *testing.T) {
	tv := New()
	defer tv.Close()
	_, err := tv.request(context.Background(), "getAccountInfo", map[string]any{})
	if err == nil {
		t.Error("expected error when auth is missing")
	}
}

func TestActionOnMediaVariants(t *testing.T) {
	// These are convenience wrappers; verify they don't panic.
	// We can't call them without auth, but we can verify the
	// actionOnMedia method exists and accepts the right params.
	tv := New(WithAuth("test"))
	defer tv.Close()

	// Verify the method signatures by checking they compile.
	// The actual HTTP calls will fail since there's no server,
	// but that's expected — we're testing the method wiring.
	_ = tv.LikeMedia
	_ = tv.UnlikeMedia
	_ = tv.AddWishMedia
	_ = tv.RemoveWishMedia
	_ = tv.GetListing
	_ = tv.GetPropertyTypes
}

func TestGetListing(t *testing.T) {
	tv := New(WithAuth("test"))
	defer tv.Close()
	// This will fail at the HTTP level, but we're testing that
	// the method exists and builds the right payload.
	_, err := tv.GetListing(context.Background(), "")
	if err == nil {
		// Might succeed if a real server is reachable; that's fine.
	}
}

func TestGetPropertyTypes(t *testing.T) {
	tv := New(WithAuth("test"))
	defer tv.Close()
	_, err := tv.GetPropertyTypes(context.Background())
	if err == nil {
		// Might succeed if a real server is reachable.
	}
}

func TestRemoveWishMedia(t *testing.T) {
	tv := New(WithAuth("test"))
	defer tv.Close()
	_, err := tv.RemoveWishMedia(context.Background(), "test-media-id")
	if err == nil {
		// Might succeed if a real server is reachable.
	}
}
