package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

/*
Go's M:N Scheduler (The G-M-P Model)
=====================================
Go multiplexes M goroutines onto N operating system threads using the G-M-P model:

- G (Goroutine):
    Lightweight user-space thread. Starts with a tiny 2KB stack that grows and shrinks dynamically.
- M (Machine / OS Thread):
    An actual OS thread managed by the OS kernel scheduler.
- P (Processor / Logical Context):
    A logical processor required to execute Go code. The number of P's is determined
    by runtime.GOMAXPROCS (defaults to the number of logical CPU cores).

Key Mechanisms:
1. Local Run Queue: Each P has its own local queue of runnable G's (lock-free, fast).
2. Global Run Queue: Shared queue checked periodically or when local queues are empty.
3. Work Stealing: When a P runs out of G's, it steals half of the tasks from another P's queue.
4. Syscall Handling: When a G performs a blocking syscall, M and G detach from P; P finds another M to keep working.
*/

func demonstrateRuntimeStats() {
	fmt.Println("--- 1. Scheduler Runtime Configuration ---")
	fmt.Printf("Logical CPU Cores: %d\n", runtime.NumCPU())
	fmt.Printf("Current GOMAXPROCS (P count): %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("Initial active Goroutines: %d\n\n", runtime.NumGoroutine())
}

func demonstrateCooperativeYield() {
	fmt.Println("--- 2. Cooperative Scheduling with runtime.Gosched() ---")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= 3; i++ {
			fmt.Println("Task A running step", i)
			// runtime.Gosched() voluntarily yields the processor (P),
			// allowing other Goroutines on the same P to execute.
			runtime.Gosched()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= 3; i++ {
			fmt.Println("Task B running step", i)
			runtime.Gosched()
		}
	}()

	wg.Wait()
	fmt.Println()
}

func demonstrateWorkStealingConcept() {
	fmt.Println("--- 3. Tracking Active Goroutines (M:N multiplexing) ---")

	const taskCount = 20
	var wg sync.WaitGroup
	wg.Add(taskCount)

	for i := 0; i < taskCount; i++ {
		go func(id int) {
			defer wg.Done()
			// Simulate computational task
			time.Sleep(20 * time.Millisecond)
		}(i)
	}

	// Observe active goroutines while scheduler distributes them among P's and M's
	fmt.Printf("Goroutines running concurrently: %d (M:N multiplexed across %d Ps)\n",
		runtime.NumGoroutine(), runtime.GOMAXPROCS(0))

	wg.Wait()
	fmt.Println("All distributed tasks finished successfully.")
}

func main() {
	demonstrateRuntimeStats()
	demonstrateCooperativeYield()
	demonstrateWorkStealingConcept()
}
