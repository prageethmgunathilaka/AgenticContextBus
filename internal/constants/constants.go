package constants

const (
	// DefaultAccessTokenTTL is the default access token TTL in seconds
	DefaultAccessTokenTTL = 3600 // 1 hour

	// DefaultRefreshTokenTTL is the default refresh token TTL in seconds
	DefaultRefreshTokenTTL = 604800 // 7 days

	// DefaultContextTTL is the default context TTL in seconds
	DefaultContextTTL = 86400 // 24 hours

	// DefaultChunkSize is the default chunk size for streaming (1MB)
	DefaultChunkSize = 1024 * 1024

	// MaxDirectContextSize is the maximum size for direct context (1MB)
	MaxDirectContextSize = 1024 * 1024

	// MaxStreamingContextSize is the maximum size for streaming contexts (100MB)
	MaxStreamingContextSize = 100 * 1024 * 1024
)

