package ttlcache

type Options struct {
	bucketsCount           int
	bucketsMapPreAllocSize int
	cleanInterval          int // seconds
}

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

func BucketsCount(v int) Option {
	return func(o *Options) {
		if v > 0 {
			o.bucketsCount = v
		}
	}
}
func BucketsMapPreAllocSize(v int) Option {
	return func(o *Options) {
		if v > 0 {
			o.bucketsMapPreAllocSize = v
		}
	}
}
func CleanInterval(v int) Option {
	return func(o *Options) {
		if v > 0 {
			o.cleanInterval = v
		}
	}
}
