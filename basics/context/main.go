package main

import (
	"context"
	"fmt"
	"time"
)

/*
The context Package
====================
The context package provides mechanisms to:
1. Propagate cancellation signals across goroutines.
2. Enforce timeouts and deadlines on long-running operations (HTTP, DB queries).
3. Carry request-scoped values (trace IDs, authentication tokens).

Key rules:
- Do not store Context inside a struct; pass ctx as the FIRST parameter to a function:
  func DoSomething(ctx context.Context, arg Type) error
- Always call the cancel function returned by WithCancel, WithTimeout, or WithDeadline (usually via defer cancel()).
- Even when a function returns, cancelling frees allocated timer resources immediately.
*/

// ==========================================
// 1. context.WithTimeout (Automatic Cancellation)
// ==========================================
func fetchWithTimeout(ctx context.Context, duration time.Duration) error {
	select {
	case <-time.After(duration):
		// Simulated database or API fetch finished
		fmt.Printf("Fetch completed successfully in %v\n", duration)
		return nil
	case <-ctx.Done():
		// Timeout expired or manual cancellation fired
		return ctx.Err() // Returns context.DeadlineExceeded or context.Canceled
	}
}

func demonstrateTimeout() {
	fmt.Println("--- 1. context.WithTimeout ---")

	// Create a context that automatically cancels after 100 milliseconds
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // Crucial: releases timer resources even if work completes early

	// Scenario A: Operation completes BEFORE timeout (Fast)
	fmt.Println("Scenario A: Fast query (50ms)")
	if err := fetchWithTimeout(ctx, 50*time.Millisecond); err != nil {
		fmt.Println("Error:", err)
	}

	// Scenario B: Operation exceeds timeout (Slow)
	fmt.Println("\nScenario B: Slow query (200ms)")
	ctxSlow, cancelSlow := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelSlow()

	if err := fetchWithTimeout(ctxSlow, 200*time.Millisecond); err != nil {
		fmt.Println("Query stopped with context error:", err) // context.DeadlineExceeded
	}
	fmt.Println()
}

// ==========================================
// 2. context.WithCancel (Manual Cancellation Propagation)
// ==========================================
func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d received cancellation signal. Cleaning up...\n", id)
			return
		default:
			// Simulating continuous processing
			time.Sleep(30 * time.Millisecond)
		}
	}
}

func demonstrateCancellation() {
	fmt.Println("--- 2. context.WithCancel (Cancellation Propagation) ---")

	ctx, cancel := context.WithCancel(context.Background())

	// Spawn 3 background workers sharing the same context tree
	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	time.Sleep(70 * time.Millisecond)
	fmt.Println("[Main] Triggering manual cancel() to stop all workers...")
	cancel() // Propagates cancellation to all child workers listening to ctx.Done()

	time.Sleep(50 * time.Millisecond) // Give workers a moment to output cleanup message
	fmt.Printf("Context error after cancel: %v\n\n", ctx.Err())
}

// ==========================================
// 3. context.WithValue (Request-Scoped Metadata)
// ==========================================
// Best practice: Use unexported custom type for context keys to prevent collisions
type contextKey string

const requestIDKey contextKey = "requestID"

func processRequest(ctx context.Context) {
	// Extract value from context
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Printf("Processing request with Trace ID: %s\n", reqID)
	} else {
		fmt.Println("No Request ID found in context")
	}
}

func demonstrateWithValue() {
	fmt.Println("--- 3. context.WithValue ---")

	// Wrap root context with request metadata
	ctx := context.WithValue(context.Background(), requestIDKey, "req-xyz-98765")
	processRequest(ctx)
}

func main() {
	demonstrateTimeout()
	demonstrateCancellation()
	demonstrateWithValue()
}
