package ttlcache

import (
	"sync/atomic"
	"time"

	xxhash "github.com/cespare/xxhash/v2"
)

// TTLCache is an in-process object caching library designed specifically for managing the caching
// and automatic release of objects with lifecycles
type TTLCache struct {
	bucketsCount uint64
	timeCached   atomic.Int64 // millisecond
	buckets      []bucket
	janitor      *janitor
}

// New return an ttlcache instance
func New(opts ...Option) *TTLCache {
	opt := setOptions(opts...)
	c := &TTLCache{
		bucketsCount: uint64(opt.bucketsCount),
		buckets:      make([]bucket, opt.bucketsCount),
	}
	for i := 0; i < len(c.buckets); i++ {
		c.buckets[i].init(opt.bucketsMapPreAllocSize)
	}
	c.timeCached.Store(time.Now().UnixMilli())
	go c.syncTimeCache()
	runJanitor(c, time.Duration(opt.cleanInterval)*time.Second)
	return c
}

// Set add an item to the cache, replacing any existing item.
//
// ttl is in seconds, must > 0,
func (c *TTLCache) Set(k string, x any, ttl int64) {
	if ttl < 1 {
		panic("TTLCache:Set ttl < 1")
	}
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	c.buckets[idx].set(k, x, c.timeCached.Load()+ttl*1000)
}

// Expire set item expire, Return false if item not found or had expired
//
// ttl is in seconds, must > 0,
func (c *TTLCache) Expire(k string, ttl int64) bool {
	if ttl < 1 {
		panic("TTLCache:Expire ttl < 1")
	}
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	now := c.timeCached.Load()
	return c.buckets[idx].expire(k, now, now+ttl*1000)
}

// Add add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns an error otherwise.
//
// ttl is in seconds, must > 0,
func (c *TTLCache) Add(k string, x any, ttl int64) bool {
	if ttl < 1 {
		panic("TTLCache:Add ttl < 1")
	}
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	now := c.timeCached.Load()
	return c.buckets[idx].add(k, x, now, now+ttl*1000)
}

// Replace set a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns false otherwise.
//
// ttl is in seconds, must > 0,
func (c *TTLCache) Replace(k string, x any, ttl int64) bool {
	if ttl < 1 {
		panic("TTLCache:Replace ttl < 1")
	}
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	now := c.timeCached.Load()
	return c.buckets[idx].replace(k, x, now, now+ttl*1000)
}

// Get get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *TTLCache) Get(k string) (any, bool) {
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	return c.buckets[idx].get(k, c.timeCached.Load())
}

// Pop pop gets an item from the cache and deletes it.
//
// The bool return indicates if the item was set.
func (c *TTLCache) Pop(k string) (any, bool) {
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	return c.buckets[idx].pop(k, c.timeCached.Load())
}

// Exist return the keys exists or not
func (c *TTLCache) Exist(k string) bool {
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	return c.buckets[idx].exists(k, c.timeCached.Load())
}

// IncrementInt an item of type int by n. Returns false if the item's value is
// not an int. If there is no error, the incremented value is returned.
//
// ttl is in seconds, must > 0,
// ttl is only effective for newly added or expired.
func (c *TTLCache) IncrementInt(k string, n int, ttl int64) (int, bool) {
	if ttl < 1 {
		panic("TTLCache:IncrementInt ttl < 1")
	}
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	now := c.timeCached.Load()
	return c.buckets[idx].incrementInt(k, n, now, now+ttl*1000)
}

// IncrementInt64 an item of type int64 by n. Returns false if the item's value is
// not an int64. If there is no error, the incremented value is returned.
//
// ttl is in seconds, must > 0,
// ttl is only effective for newly added or expired.
func (c *TTLCache) IncrementInt64(k string, n int64, ttl int64) (int64, bool) {
	if ttl < 1 {
		panic("TTLCache:IncrementInt64 ttl < 1")
	}
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	now := c.timeCached.Load()
	return c.buckets[idx].incrementInt64(k, n, now, now+ttl*1000)
}

// IncrementFloat64 an item of type float64 by n. Returns false if the item's value
// is not an float64. If there is no error, the incremented value is returned.
//
// ttl is in seconds, must > 0,
// ttl is only effective for newly added or expired.
func (c *TTLCache) IncrementFloat64(k string, n float64, ttl int64) (float64, bool) {
	if ttl < 1 {
		panic("TTLCache:IncrementFloat64 ttl < 1")
	}
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	now := c.timeCached.Load()
	return c.buckets[idx].incrementFloat64(k, n, now, now+ttl*1000)
}

// Delete an item from the cache. Does nothing if the key is not in the cache.
func (c *TTLCache) Delete(k string) {
	h := xxhash.Sum64String(k)
	idx := h % c.bucketsCount
	c.buckets[idx].delete(k)
}

// Items return cached objects(include expired)
func (c *TTLCache) Items() int {
	n := 0
	for i := 0; i < len(c.buckets); i++ {
		n += c.buckets[i].allItems()
	}
	return n
}

// Delete all expired items from the cache.
func (c *TTLCache) deleteExpired() {
	now := c.timeCached.Load()
	for i := 0; i < len(c.buckets); i++ {
		c.buckets[i].deleteExpired(now)
	}
}

func (c *TTLCache) syncTimeCache() {
	ticker := time.NewTicker(time.Millisecond * 300)
	for {
		select {
		case now := <-ticker.C:
			c.timeCached.Store(now.UnixMilli())
		}
	}
}

type janitor struct {
	Interval time.Duration
}

func (j *janitor) Run(c *TTLCache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		}
	}
}

func runJanitor(c *TTLCache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
	}
	c.janitor = j
	go j.Run(c)
}
