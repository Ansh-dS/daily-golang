package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Goroutines:
===========
A goroutine is a lightweight thread of execution managed by the Go runtime, not the OS kernel.
- OS thread: typically allocated with 1MB-8MB of stack space; context switching requires kernel transitions.
- Goroutine: starts with just ~2KB stack size, dynamically expanding/contracting on demand; context switching happens in user space.
*/

func downloadFile(name string, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[Start] Downloading %s (takes %v)...\n", name, duration)
	time.Sleep(duration)
	fmt.Printf("[Done]  %s downloaded!\n", name)
}

func main() {
	// ==========================================
	// 1. Sequential vs Concurrent Execution
	// ==========================================
	fmt.Println("--- 1. Comparing Sequential vs Concurrent Execution ---")

	files := []string{"image.png", "video.mp4", "report.pdf"}

	// Sequential execution
	startSeq := time.Now()
	for _, file := range files {
		time.Sleep(50 * time.Millisecond) // Simulated download
		_ = file
	}
	fmt.Printf("Sequential execution time: ~%v\n", time.Since(startSeq).Round(time.Millisecond))

	// Concurrent execution with Goroutines
	startConc := time.Now()
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		// Launch each download in its own goroutine
		go downloadFile(file, 50*time.Millisecond, &wg)
	}
	wg.Wait()
	fmt.Printf("Concurrent execution time: ~%v (all 3 ran in parallel!)\n\n", time.Since(startConc).Round(time.Millisecond))

	// ==========================================
	// 2. Anonymous Functions and Closures
	// ==========================================
	fmt.Println("--- 2. Anonymous Functions & Passing Loop Variables ---")

	var wgLoop sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wgLoop.Add(1)

		// BEST PRACTICE: Pass loop variable `i` as parameter `val` to avoid closure capture issues
		go func(val int) {
			defer wgLoop.Done()
			fmt.Printf("Goroutine received value: %d\n", val)
		}(i)
	}
	wgLoop.Wait()
	fmt.Println()

	// ==========================================
	// 3. Goroutine Lifecycle: Main terminates Goroutines
	// ==========================================
	fmt.Println("--- 3. Main Goroutine Lifecycle ---")
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("This might not print if main() doesn't wait!")
	}()
	fmt.Println("Main goroutine is the boss: when main() returns, all child goroutines terminate immediately.")
	fmt.Println("That's why tools like sync.WaitGroup or channels are essential to wait for completion.")
	fmt.Println()
}
