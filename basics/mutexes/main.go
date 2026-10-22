package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Mutexes (sync.Mutex and sync.RWMutex)
=====================================
When multiple goroutines access shared memory concurrently:
- If at least one goroutine writes to that memory, you have a DATA RACE.
- Mutexes provide Mutual Exclusion to protect critical sections.

1. sync.Mutex:
   - Only ONE goroutine can hold the lock at any time (for read or write).
   - Use mu.Lock() and defer mu.Unlock().

2. sync.RWMutex (Reader/Writer Mutex):
   - Multiple concurrent readers can hold RLock() simultaneously.
   - Writers acquire Lock() exclusively (blocking both readers and other writers).
   - Ideal for read-heavy workloads (e.g., caches, lookups).
*/

// SafeCounter uses sync.Mutex to protect an integer counter
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++ // Critical Section: protected from concurrent writes
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// Cache uses sync.RWMutex for high-concurrency read operations
type Cache struct {
	rwMu sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

// Set requires an exclusive write lock
func (c *Cache) Set(key, value string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.data[key] = value
}

// Get allows multiple readers at the exact same time
func (c *Cache) Get(key string) (string, bool) {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func main() {
	// ==========================================
	// 1. sync.Mutex (Mutual Exclusion)
	// ==========================================
	fmt.Println("--- 1. sync.Mutex: Safe Counter ---")

	counter := SafeCounter{}
	var wg sync.WaitGroup

	const routines = 1000
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Expected count: %d | Actual count: %d (Zero race conditions!)\n\n", routines, counter.Value())

	// ==========================================
	// 2. sync.RWMutex (Read/Write Mutex)
	// ==========================================
	fmt.Println("--- 2. sync.RWMutex: Concurrent Cache (Many Readers, Few Writers) ---")

	cache := NewCache()
	cache.Set("app_name", "daily-golang")

	var rwWg sync.WaitGroup

	// Start 5 concurrent readers (they read simultaneously without blocking each other)
	for r := 1; r <= 5; r++ {
		rwWg.Add(1)
		go func(readerID int) {
			defer rwWg.Done()
			val, _ := cache.Get("app_name")
			fmt.Printf("Reader %d fetched: %s\n", readerID, val)
		}(r)
	}

	// Start 1 concurrent writer (acquires exclusive lock)
	rwWg.Add(1)
	go func() {
		defer rwWg.Done()
		time.Sleep(10 * time.Millisecond)
		cache.Set("version", "v1.0.0")
		fmt.Println("[Writer] Updated cache version")
	}()

	rwWg.Wait()
	version, _ := cache.Get("version")
	fmt.Printf("Final Cache 'version': %s\n", version)
}
