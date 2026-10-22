package main

import (
	"fmt"
	"sync"
	"time"
)

/*
sync.WaitGroup
===============
A WaitGroup waits for a collection of goroutines to finish.

Golden Rules:
1. Always call wg.Add(1) BEFORE launching the goroutine (not inside it).
   Calling it inside causes a race condition where wg.Wait() might execute before Add().
2. Always use defer wg.Done() as the first line of the goroutine to guarantee it runs even on panics/returns.
3. NEVER copy a sync.WaitGroup by value! Always pass it by pointer (*sync.WaitGroup).
*/

// Task represents a unit of work that can succeed or fail
type Result struct {
	ID    int
	Value string
	Err   error
}

// performTask receives wg as a POINTER (*sync.WaitGroup)
func performTask(id int, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done() // Decrements counter when function finishes

	// Simulate work
	time.Sleep(40 * time.Millisecond)

	if id == 2 {
		results <- Result{ID: id, Err: fmt.Errorf("task %d failed due to network error", id)}
		return
	}

	results <- Result{ID: id, Value: fmt.Sprintf("Success data from task %d", id)}
}

func main() {
	fmt.Println("--- 1. Basic WaitGroup Usage ---")

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		// Rule 1: Add BEFORE starting the goroutine
		wg.Add(1)

		go func(taskID int) {
			// Rule 2: defer wg.Done() inside the goroutine
			defer wg.Done()
			fmt.Printf("Worker %d executed\n", taskID)
		}(i)
	}

	// Blocks until counter becomes 0
	wg.Wait()
	fmt.Println("All workers finished.")
	fmt.Println()

	// ==========================================
	// 2. Collecting Results & Errors Concurrently
	// ==========================================
	fmt.Println("--- 2. Aggregating Concurrent Results Safely ---")

	const totalTasks = 4
	resultsCh := make(chan Result, totalTasks)

	var taskWg sync.WaitGroup

	for i := 1; i <= totalTasks; i++ {
		taskWg.Add(1)
		// Rule 3: Pass by pointer (&taskWg)
		go performTask(i, &taskWg, resultsCh)
	}

	// Wait in background and close the channel when done
	go func() {
		taskWg.Wait()
		close(resultsCh)
	}()

	// Read results as they complete
	for res := range resultsCh {
		if res.Err != nil {
			fmt.Printf("[Error] Task %d: %v\n", res.ID, res.Err)
		} else {
			fmt.Printf("[Success] Task %d: %s\n", res.ID, res.Value)
		}
	}
}
