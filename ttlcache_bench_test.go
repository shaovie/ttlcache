package ttlcache

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) string {
	if n < 1 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

type live struct {
	i int64
	n string
}

func (l *live) loop() {
	ticker := time.NewTicker(time.Millisecond * 20)
	for {
		select {
		case now := <-ticker.C:
			l.i = now.Unix()
		}
	}
}
func newLive() *live {
	l := &live{}
	go l.loop()
	return l
}
func livex(l *live, wg *sync.WaitGroup) {
	defer wg.Done()
	//runtime.LockOSThread()
	for i := int64(0); i < 10000; i++ {
		j := rand.Int63() % 100
		if i == j {
		}
	}
}
func get(cache *TTLCache, wg *sync.WaitGroup) {
	defer wg.Done()
	//runtime.LockOSThread()

	for i := 0; i < 100000; i++ {
		val, found := cache.Get("small-" + strconv.FormatInt(int64(i), 10))
		if found {
			if val.(int) != i {
				panic("set val invalid")
			}
		}
	}
}
func set(cache *TTLCache, wg *sync.WaitGroup) {
	defer wg.Done()
	//runtime.LockOSThread()

	for i := 0; i < 100000; i++ {
		cache.Set("small-"+strconv.FormatInt(int64(i), 10), i, 3 /*seconds*/)
	}
}
func BenchmarkLive(b *testing.B) {
	fmt.Println("hello boy")
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	l := &live{}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go livex(l, &wg)
	}
	wg.Wait()
}
func BenchmarkSetGet(b *testing.B) {
	fmt.Println("hello boy")
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	cache := New(SetBucketsCount(512),
		SetBucketsMapPreAllocSize(256),
		SetCleanInterval(10),
	)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		set(cache, &wg)
	}
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go set(cache, &wg)
		go get(cache, &wg)
	}
	wg.Wait()
}
func TestAll(t *testing.T) {
	cache := New(SetBucketsCount(512),
		SetBucketsMapPreAllocSize(256),
		SetCleanInterval(10),
	)
	cache.Set("ttlcache", "nb", 1)
	val, found := cache.Get("ttlcache")
	if val.(string) != "nb" {
		t.Error("get set val error")
		return
	}
	if exist := cache.Exist("ttlcache"); !exist {
		t.Error("exist error")
		return
	}
	time.Sleep(time.Millisecond * 1500)
	_, found = cache.Get("ttlcache")
	if found {
		t.Error("set expiration error")
		return
	}

	if ok := cache.Add("ttlcache", "nb", 2); !ok {
		t.Error("add error")
		return
	}
	if ok := cache.Add("ttlcache", "nb", 2); ok {
		t.Error("add an exist item error")
		return
	}
	val, found = cache.Get("ttlcache")
	if !found || val.(string) != "nb" {
		t.Error("add item error")
		return
	}

	if ok := cache.Replace("ttlcache", "no1", 1); !ok {
		t.Error("replace an exist item error")
		return
	}
	val, found = cache.Get("ttlcache")
	if !found || val.(string) != "no1" {
		t.Error("replace item error")
		return
	}
	time.Sleep(time.Millisecond * 1300)
	_, found = cache.Get("ttlcache")
	if found {
		t.Error("replace expiration error")
		return
	}

	n, ok := cache.IncrementInt("ttlcacheno", 1, 1)
	if !ok || n != 1 {
		t.Error("incrementInt error")
		return
	}
	time.Sleep(time.Millisecond * 1200)
	_, found = cache.Get("ttlcacheno")
	if found {
		t.Error("incrementInt expiration error")
		return
	}
	n1, ok := cache.IncrementInt64("ttlcacheno", 10, 1)
	if !ok || n1 != 10 {
		t.Error("incrementInt64 error")
		return
	}
	f, ok := cache.IncrementFloat64("ttlcacheno", 10, 1)
	if ok {
		t.Error("incrementFloat64 by int64 error")
		return
	}
	cache.Delete("ttlcacheno")
	f, ok = cache.IncrementFloat64("ttlcacheno", 10, 1)
	if !ok || !(math.Abs(f-10) < 0.000001) {
		t.Error("incrementFloat64 error")
		return
	}
	fmt.Printf("items: %d\n", cache.Items())
	time.Sleep(time.Second * 10)
	fmt.Printf("items: %d\n", cache.Items())
}
