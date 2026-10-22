package main

import (
	"fmt"
	"time"
)

/*
Channels: Buffered vs. Unbuffered
==================================

1. Unbuffered Channel (make(chan T)):
   - Buffer size = 0.
   - Synchronous: A send operation blocks until another goroutine receives from it.
   - Guarantees immediate synchronization (a "handshake").

2. Buffered Channel (make(chan T, capacity)):
   - Buffer size = capacity.
   - Asynchronous: A send operation only blocks when the buffer is FULL.
   - A receive operation only blocks when the buffer is EMPTY.
   - Decouples producer from consumer throughput.
*/

func demonstrateUnbuffered() {
	fmt.Println("--- 1. Unbuffered Channel (Synchronous Handshake) ---")

	ch := make(chan string) // cap is 0

	go func() {
		fmt.Println("[Goroutine] Ready to send message...")
		time.Sleep(100 * time.Millisecond) // Simulating work before sending
		ch <- "Hello from unbuffered sender!"
		fmt.Println("[Goroutine] Send completed! (Receiver was waiting)")
	}()

	fmt.Println("[Main] Waiting at receive site...")
	msg := <-ch // Blocks here until goroutine sends
	fmt.Println("[Main] Received message:", msg)
	fmt.Println()
}

func demonstrateBuffered() {
	fmt.Println("--- 2. Buffered Channel (Asynchronous Queue) ---")

	// Buffered channel holding up to 3 elements
	bufCh := make(chan int, 3)

	fmt.Printf("Initial state -> len: %d, cap: %d\n", len(bufCh), cap(bufCh))

	// Sends do NOT block because buffer capacity is not yet reached
	bufCh <- 10
	bufCh <- 20
	bufCh <- 30

	fmt.Printf("After 3 sends -> len: %d, cap: %d (Buffer is full)\n", len(bufCh), cap(bufCh))

	// Reading values off the buffer (FIFO order)
	fmt.Println("Received item:", <-bufCh)
	fmt.Println("Received item:", <-bufCh)
	fmt.Printf("After 2 receives -> len: %d, cap: %d\n\n", len(bufCh), cap(bufCh))
}

// Directional channels: chan<- int (send-only), <-chan int (receive-only)
func producer(out chan<- int) {
	for i := 1; i <= 3; i++ {
		out <- i
	}
	close(out) // Sender closes channel to signal no more data
}

func consumer(in <-chan int) {
	// 1. Ranging over channel until it is closed
	for val := range in {
		fmt.Printf("Consumer read: %d\n", val)
	}

	// 2. Comma-ok idiom to test if channel is closed
	val, ok := <-in
	fmt.Printf("After channel closed -> value: %d, open: %t\n", val, ok)
}

func demonstrateClosingAndDirectional() {
	fmt.Println("--- 3. Closing Channels, Range, and Directional Channels ---")
	ch := make(chan int, 3)

	producer(ch)
	consumer(ch)
	fmt.Println()
}

func main() {
	demonstrateUnbuffered()
	demonstrateBuffered()
	demonstrateClosingAndDirectional()
}
