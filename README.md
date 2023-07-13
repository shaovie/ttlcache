## ttlcache

`ttlcache` is an in-process object caching library designed specifically for managing the caching and automatic release of objects with lifecycles. The cache operates on a time-based expiration mechanism, measured in seconds, and utilizes timestamp caching to avoid real-time system time retrieval. This significantly improves efficiency. Object storage employs sharding techniques to reduce concurrent access competition.

Usage:
```
import "github.com/shaovie/ttlcache"

func main() {
	cache := ttlcache.New(ttlcache.BucketsCount(512),
		ttlcache.BucketsMapPreAllocSize(256),
		ttlcache.CleanInterval(10),
	)
	cache.Set("ttlcache", "nb", 1/*second*/) // The lifecycle is 1 second
	val, found := cache.Get("ttlcache")
	if !found {
		fmt.Println("set val error")
		return
	}
}
```

Usage Reference: `ttlcache_bench_test.go`
