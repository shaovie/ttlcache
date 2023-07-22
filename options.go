package ttlcache

// Options provides all optional parameters
type Options struct {
	bucketsCount           int
	bucketsMapPreAllocSize int
	cleanInterval          int // seconds
}

// Option function
type Option func(*Options)

func setOptions(optL ...Option) *Options {
	opts := &Options{
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
	return func(o *Options) {
		if v > 0 {
			o.bucketsCount = v
		}
	}
}

// BucketsMapPreAllocSize map prealloc size
func BucketsMapPreAllocSize(v int) Option {
	return func(o *Options) {
		if v > 0 {
			o.bucketsMapPreAllocSize = v
		}
	}
}

// CleanInterval cleans up expired object cycles.
func CleanInterval(v int) Option {
	return func(o *Options) {
		if v > 0 {
			o.cleanInterval = v
		}
	}
}
