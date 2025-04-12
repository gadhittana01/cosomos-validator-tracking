package constant

const (
	// Format of date & time
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)

const (
	// RetryCount is the number of times to retry the request
	RetryCount = 3
)

const (
	// DefaultLimit & DefaultPage is the default limit & page for pagination
	DefaultLimit = 10
	DefaultPage  = 1
)

const (
	// CacheKey is the key for the cache
	ValidatorHourlySnapshotCacheKey   = "validator_hourly_snapshot"
	ValidatorDailySnapshotCacheKey    = "validator_daily_snapshot"
	ValidatorDelegatorHistoryCacheKey = "validator_delegator_history"
)
