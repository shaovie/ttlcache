package ttlcache

// options provides all optional parameters
type options struct {
	bucketsCount           int
	bucketsMapPreAllocSize int
	cleanInterval          int // seconds
}

// Option function
type Option func(*options)

func setOptions(optL ...Option) *options {
	opts := &options{
		bucketsCount:           128,
		bucketsMapPreAllocSize: 128,
		cleanInterval:          10, // seconds
	}

	for _, opt := range optL {
		opt(opts)
	}
	return opts
}

// BucketsCount can effectively reduce the number of competing occurrences in concurrent access to ttlcache.
func BucketsCount(v int) Option {
	return func(o *options) {
		if v < 1 {
			panic("BucketsCount: param is illegal")
		}
		o.bucketsCount = v
	}
}

// BucketsMapPreAllocSize map prealloc size
func BucketsMapPreAllocSize(v int) Option {
	return func(o *options) {
		if v < 1 {
			panic("BucketsMapPreAllocSize: param is illegal")
		}
		o.bucketsMapPreAllocSize = v
	}
}

// CleanInterval cleans up expired object cycles.
func CleanInterval(v int) Option {
	return func(o *options) {
		if v < 1 {
			panic("CleanInterval: param is illegal")
		}
		o.cleanInterval = v
	}
}
