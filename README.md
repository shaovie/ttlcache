## ttlcache

`ttlcache` is an in-process object caching library designed specifically for managing the caching and automatic release of objects with lifecycles. The cache operates on a time-based expiration mechanism, measured in seconds, and utilizes timestamp caching to avoid real-time system time retrieval. This significantly improves efficiency. Object storage employs sharding techniques to reduce concurrent access competition.

Usage:
```
import "github.com/shaovie/ttlcache"

func main() {
	cache := ttlcache.New(ttlcache.SetBucketsCount(512),
		ttlcache.SetBucketsMapPreAllocSize(256),
		ttlcache.SetCleanInterval(10),
	)
	cache.Set("ttlcache", "nb", 1)
	val, found := cache.Get("ttlcache")
	if !found {
		fmt.Println("set val error")
		return
	}
}
```

Usage Reference: `ttlcache_bench_test.go`
