package valve

import "testing"

func TestGet(t *testing.T) {
	// TODO (twiny) implment case:
	// - key not found
	// - empty ip
}

func BenchmarkGet(b *testing.B) {
	limiter := NewLimiter(10, 5)

	for n := 0; n < b.N; n++ {
		limit := limiter.Get("127.0.0.1")

		if limit.Allow() {
			continue
		}
	}
}
