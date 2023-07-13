## ttlcache

`ttlcache` is an in-process object caching library designed specifically for managing the caching and automatic release of objects with lifecycles. The cache operates on a time-based expiration mechanism, measured in seconds, and utilizes timestamp caching to avoid real-time system time retrieval. This significantly improves efficiency. Object storage employs sharding techniques to reduce concurrent access competition.

Usage:
To use ttlcache, execute the following command: `go get github.com/shaovie/ttlcache`

Usage Reference: `ttlcache_bench_test.go`
