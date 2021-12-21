package valve

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	// TODO (twiny) implment case:
	// - key not found
	// - empty ip
}

func BenchmarkGet(b *testing.B) {
	limiter := NewLimiter(10, 5, 10*time.Minute)

	for n := 0; n < b.N; n++ {
		if limiter.Allow("127.0.0.1") {
			continue
		}
	}
}
