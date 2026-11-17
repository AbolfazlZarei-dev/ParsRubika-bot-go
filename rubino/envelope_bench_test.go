package rubino

import (
	"encoding/json"
	"testing"
)

// BenchmarkRequestPayloadMarshal measures the cost of encoding one
// Rubino request envelope — the struct + cached client-block path used
// by every API call. It isolates the encode step from the network.
func BenchmarkRequestPayloadMarshal(b *testing.B) {
	input := map[string]any{
		"profile_id":  "0123456789abcdef",
		"limit":       50,
		"sort":        "FromMax",
		"max_id":      "",
		"next_target": nil,
	}
	p := requestPayload{
		APIVersion: "0",
		Auth:       "abcdefgh0123456789ABCDEFGHIJKLMNOPQRSTUVWX",
		Client:     clientBlockJSON(),
		Data:       input,
		Method:     "getMyPosts",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(p); err != nil {
			b.Fatal(err)
		}
	}
}
