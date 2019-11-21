package libp2p

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	cache := NewTimeCache(time.Minute)

	cache.Add("test")

	if !cache.Has("test") {
		t.Fatal("should have 'test' key")
	}
}

func TestConcurrentAdd(t *testing.T) {
	cache := NewTimeCache(time.Minute)

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(item int) {
			cache.Add(strconv.Itoa(item))
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < 10; i++ {
		if !cache.Has(strconv.Itoa(i)) {
			t.Fatalf("should have '%v' key", i)
		}
	}
}

func TestExpiration(t *testing.T) {
	cache := NewTimeCache(500 * time.Millisecond)
	for i := 0; i < 5; i++ {
		cache.Add(strconv.Itoa(i))
		time.Sleep(100 * time.Millisecond)
	}

	if cache.Has(strconv.Itoa(0)) {
		t.Fatal("should have dropped '0' key from the cache already")
	}
}

func BenchmarkAdd(b *testing.B) {
	cache := NewTimeCache(time.Minute)

	for i := 0; i < b.N; i++ {
		cache.Add(strconv.Itoa(i))
	}
}

func BenchmarkConcurrentAdd(b *testing.B) {
	cache := NewTimeCache(time.Minute)

	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func(item int) {
			cache.Add(strconv.Itoa(item))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func BenchmarkHas(b *testing.B) {
	cache := NewTimeCache(time.Minute)

	for i := 0; i < b.N; i++ {
		cache.Add(strconv.Itoa(i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Has(strconv.Itoa(i))
	}
}

func BenchmarkConcurrentHas(b *testing.B) {
	cache := NewTimeCache(time.Minute)

	for i := 0; i < b.N; i++ {
		cache.Add(strconv.Itoa(i))
	}

	b.ResetTimer()

	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func(item int) {
			cache.Has(strconv.Itoa(item))
			wg.Done()
		}(i)
	}

	wg.Wait()
}
