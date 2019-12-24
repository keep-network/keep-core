package libp2p

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	cache := newTimeCache(time.Minute)

	cache.add("test")

	if !cache.has("test") {
		t.Fatal("should have 'test' key")
	}
}

func TestConcurrentAdd(t *testing.T) {
	cache := newTimeCache(time.Minute)

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(item int) {
			cache.add(strconv.Itoa(item))
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < 10; i++ {
		if !cache.has(strconv.Itoa(i)) {
			t.Fatalf("should have '%v' key", i)
		}
	}
}

func TestSweeping(t *testing.T) {
	cache := newTimeCache(200 * time.Millisecond)

	cache.add("0")
	time.Sleep(100 * time.Millisecond)

	cache.add("1")
	time.Sleep(100 * time.Millisecond)

	cache.sweep()

	if cache.has("0") {
		t.Fatal("should have '0' dropped from the cache")
	}
	if !cache.has("1") {
		t.Fatal("should not have '1' dropped from the cache")
	}
}

func TestExpiration(t *testing.T) {
	cache := newTimeCache(400 * time.Millisecond)
	for i := 0; i < 6; i++ {
		cache.add(strconv.Itoa(i))
		time.Sleep(100 * time.Millisecond)
	}

	if cache.has("0") {
		t.Fatal("should have '0' dropped from the cache")
	}
	if cache.has("1") {
		t.Fatal("should have '1' dropped from the cache")
	}
	if !cache.has("2") {
		t.Fatal("should not have '2' dropped from the cache")
	}
	if !cache.has("3") {
		t.Fatal("should not have '3' dropped from the cache")
	}
	if !cache.has("4") {
		t.Fatal("should not have '4' dropped from the cache")
	}
	if !cache.has("5") {
		t.Fatal("should not have '5' dropped from the cache")
	}
}

func TestExpirationSameElement(t *testing.T) {
	cache := newTimeCache(500 * time.Millisecond)

	for i := 0; i < 4; i++ {
		cache.add(strconv.Itoa(i))
	}
	time.Sleep(200 * time.Millisecond)

	cache.add(strconv.Itoa(1))
	time.Sleep(400 * time.Millisecond)
	cache.add(strconv.Itoa(4))

	if cache.has("0") {
		t.Fatal("should have 'zero' dropped from the cache")
	}
	if !cache.has("1") {
		t.Fatal("should not have 'one' dropped from the cache")
	}
	if cache.has("2") {
		t.Fatal("should have 'two' dropped from the cache")
	}
	if cache.has("3") {
		t.Fatal("should have 'three' dropped from the cache")
	}
	if !cache.has("4") {
		t.Fatal("should not have 'four' dropped from the cache")
	}
}

func BenchmarkAdd(b *testing.B) {
	cache := newTimeCache(time.Minute)

	for i := 0; i < b.N; i++ {
		cache.add(strconv.Itoa(i))
	}
}

func BenchmarkConcurrentAdd(b *testing.B) {
	cache := newTimeCache(time.Minute)

	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func(item int) {
			cache.add(strconv.Itoa(item))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func BenchmarkHas(b *testing.B) {
	cache := newTimeCache(time.Minute)

	for i := 0; i < b.N; i++ {
		cache.add(strconv.Itoa(i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.has(strconv.Itoa(i))
	}
}

func BenchmarkConcurrentHas(b *testing.B) {
	cache := newTimeCache(time.Minute)

	for i := 0; i < b.N; i++ {
		cache.add(strconv.Itoa(i))
	}

	b.ResetTimer()

	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func(item int) {
			cache.has(strconv.Itoa(item))
			wg.Done()
		}(i)
	}

	wg.Wait()
}
