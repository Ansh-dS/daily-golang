package main

import (
	"fmt"
	"sync"
	"time"
)

// printMessage runs concurrently as a goroutine
func printMessage(msg string, id int) {
	fmt.Printf("[Goroutine %d]: %s\n", id, msg)
}

// worker simulates processing a job and sending the result back via a channel
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job)
		time.Sleep(50 * time.Millisecond) // Simulate work
		fmt.Printf("Worker %d finished job %d\n", id, job)
		results <- job * 2
	}
}

func main() {
	// ==========================================
	// 1. Basic Goroutine (go keyword)
	// ==========================================
	fmt.Println("--- 1. Basic Goroutine with WaitGroup ---")

	var wg sync.WaitGroup

	// Launching goroutines using named and anonymous functions
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		// Best practice: Pass loop variables as arguments into the goroutine closure
		go func(id int) {
			defer wg.Done()
			printMessage("Hello from concurrent task", id)
		}(i)
	}

	// Wait until all 3 goroutines have called wg.Done()
	wg.Wait()
	fmt.Println("All basic goroutines completed.\n")

	// ==========================================
	// 2. Communication via Channels
	// ==========================================
	fmt.Println("--- 2. Channel Communication ---")

	// Unbuffered channel: sender blocks until receiver is ready
	ch := make(chan string)

	go func() {
		// Send data to channel
		ch <- "Data sent from background goroutine!"
	}()

	// Receive data from channel (blocks until data arrives)
	received := <-ch
	fmt.Println("Received:", received, "\n")

	// ==========================================
	// 3. Worker Pool Pattern (Goroutines + Channels + WaitGroup)
	// ==========================================
	fmt.Println("--- 3. Worker Pool Pattern ---")

	const numJobs = 5
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var workerWg sync.WaitGroup

	// Start worker goroutines
	for w := 1; w <= numWorkers; w++ {
		workerWg.Add(1)
		go worker(w, jobs, results, &workerWg)
	}

	// Send jobs into the channel
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Closing signals to workers that no more jobs will be sent

	// Wait for workers in the background, then close results
	go func() {
		workerWg.Wait()
		close(results)
	}()

	// Collect results
	for res := range results {
		fmt.Println("Collected result:", res)
	}

	fmt.Println("\nWorker pool processing complete.")
}
