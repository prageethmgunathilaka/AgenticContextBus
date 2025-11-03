package constants

import "testing"

func TestConstantsHaveExpectedRanges(t *testing.T) {
	if DefaultChunkSize <= 0 || MaxChunkSize < DefaultChunkSize {
		t.Fatal("invalid chunk size constants")
	}
	if DefaultRateLimitRequests <= 0 || DefaultMaxConcurrentStreams <= 0 {
		t.Fatal("invalid rate limit constants")
	}
	if DefaultAccessTokenTTL <= 0 || DefaultRefreshTokenTTL <= 0 {
		t.Fatal("invalid token TTLs")
	}
	if DefaultContextTTL <= 0 || IdempotencyKeyTTL <= 0 {
		t.Fatal("invalid storage TTLs")
	}
	if DefaultPageLimit <= 0 || MaxPageLimit < DefaultPageLimit {
		t.Fatal("invalid pagination limits")
	}
}


