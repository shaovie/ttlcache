package ttlcache

import (
	"sync"
)

type cacheItem struct {
	Object     any
	Expiration int64
}

type bucket struct {
	mtx   sync.RWMutex
	items map[string]cacheItem
}

func (b *bucket) init(bucketsMapPreAllocSize int) {
	b.items = make(map[string]cacheItem, bucketsMapPreAllocSize)
}
func (b *bucket) set(k string, x any, expiration int64) {
	b.mtx.Lock()
	b.items[k] = cacheItem{
		Object:     x,
		Expiration: expiration,
	}
	b.mtx.Unlock()
}
func (b *bucket) get(k string, now int64) (any, bool) {
	b.mtx.RLock()
	item, found := b.items[k]
	if !found || now > item.Expiration {
		b.mtx.RUnlock()
		return nil, false
	}
	b.mtx.RUnlock()
	return item.Object, true
}
func (b *bucket) exists(k string, now int64) bool {
	b.mtx.RLock()
	item, found := b.items[k]
	if !found || now > item.Expiration {
		b.mtx.RUnlock()
		return false
	}
	b.mtx.RUnlock()
	return true
}
func (b *bucket) expire(k string, now, expiration int64) bool {
	b.mtx.Lock()
	item, found := b.items[k]
	if !found || now > item.Expiration {
		b.mtx.Unlock()
		return false
	}
	item.Expiration = expiration
	b.items[k] = item
	b.mtx.Unlock()
	return true
}
func (b *bucket) add(k string, x any, now, expiration int64) bool {
	b.mtx.Lock()
	item, found := b.items[k]
	if found && now < item.Expiration {
		b.mtx.Unlock()
		return false
	}
	b.items[k] = cacheItem{
		Object:     x,
		Expiration: expiration,
	}
	b.mtx.Unlock()
	return true
}
func (b *bucket) replace(k string, x any, now, expiration int64) bool {
	b.mtx.Lock()
	item, found := b.items[k]
	if !found || now > item.Expiration {
		b.mtx.Unlock()
		return false
	}
	b.items[k] = cacheItem{
		Object:     x,
		Expiration: expiration,
	}
	b.mtx.Unlock()
	return true
}
func (b *bucket) pop(k string, now int64) (any, bool) {
	b.mtx.Lock()
	item, found := b.items[k]
	if !found || now > item.Expiration {
		b.mtx.Unlock()
		return nil, false
	}
	delete(b.items, k)
	b.mtx.Unlock()
	return item.Object, true
}
func (b *bucket) delete(k string) {
	b.mtx.Lock()
	_, found := b.items[k]
	if !found {
		b.mtx.Unlock()
		return
	}
	delete(b.items, k)
	b.mtx.Unlock()
}
func (b *bucket) incrementInt(k string, n int, now, expiration int64) (int, bool) {
	b.mtx.Lock()
	item, found := b.items[k]
	if !found || now > item.Expiration {
		b.items[k] = cacheItem{
			Object:     n,
			Expiration: expiration,
		}
		b.mtx.Unlock()
		return n, true
	}
	rv, ok := item.Object.(int)
	if !ok {
		b.mtx.Unlock()
		return 0, false
	}
	nv := rv + n
	item.Object = nv
	b.items[k] = item
	b.mtx.Unlock()
	return nv, true
}
func (b *bucket) incrementInt64(k string, n int64, now, expiration int64) (int64, bool) {
	b.mtx.Lock()
	item, found := b.items[k]
	if !found || now > item.Expiration {
		b.items[k] = cacheItem{
			Object:     n,
			Expiration: expiration,
		}
		b.mtx.Unlock()
		return n, true
	}
	rv, ok := item.Object.(int64)
	if !ok {
		b.mtx.Unlock()
		return 0, false
	}
	nv := rv + n
	item.Object = nv
	b.items[k] = item
	b.mtx.Unlock()
	return nv, true
}
func (b *bucket) incrementFloat64(k string, n float64, now, expiration int64) (float64, bool) {
	b.mtx.Lock()
	item, found := b.items[k]
	if !found || now > item.Expiration {
		b.items[k] = cacheItem{
			Object:     n,
			Expiration: expiration,
		}
		b.mtx.Unlock()
		return n, true
	}
	rv, ok := item.Object.(float64)
	if !ok {
		b.mtx.Unlock()
		return 0, false
	}
	nv := rv + n
	item.Object = nv
	b.items[k] = item
	b.mtx.Unlock()
	return nv, true
}
func (b *bucket) deleteExpired(now int64) {
	b.mtx.Lock()
	for k, v := range b.items {
		if now > v.Expiration {
			delete(b.items, k)
		}
	}
	b.mtx.Unlock()
}
func (b *bucket) allItems() int {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	return len(b.items)
}
