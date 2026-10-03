package main

import "fmt"

func main() {
	// 1. Initializing a slice with make([]T, len, cap)
	// Here len = 2, cap = 4
	s := make([]int, 2, 4)
	s[0] = 10
	s[1] = 20

	fmt.Println("--- Initial Slice ---")
	fmt.Printf("s: %v | len: %d | cap: %d\n\n", s, len(s), cap(s))

	// 2. Appending when len < cap
	// The underlying array has room, so cap does NOT change.
	s = append(s, 30)
	fmt.Println("--- After appending 30 (within capacity) ---")
	fmt.Printf("s: %v | len: %d | cap: %d\n", s, len(s), cap(s))

	s = append(s, 40)
	fmt.Println("--- After appending 40 (reached capacity limit) ---")
	fmt.Printf("s: %v | len: %d | cap: %d\n\n", s, len(s), cap(s))

	// 3. Appending when len == cap
	// The backing array is full!
	// Go allocates a new, larger backing array (often doubles capacity for smaller slices),
	// copies existing elements over, and appends the new value.
	s = append(s, 50)
	fmt.Println("--- After appending 50 (capacity exceeded -> reallocation) ---")
	fmt.Printf("s: %v | len: %d | cap: %d\n\n", s, len(s), cap(s))

	// 4. Growth demonstration in a loop starting from an empty slice
	fmt.Println("--- Growth Demonstration (Dynamic Resizing) ---")
	var dynamic []int
	prevCap := cap(dynamic)

	for i := 1; i <= 10; i++ {
		dynamic = append(dynamic, i)
		if cap(dynamic) != prevCap {
			fmt.Printf("Appended %2d -> len: %2d, cap changed: %2d -> %2d\n", i, len(dynamic), prevCap, cap(dynamic))
			prevCap = cap(dynamic)
		} else {
			fmt.Printf("Appended %2d -> len: %2d, cap remains: %2d\n", i, len(dynamic), cap(dynamic))
		}
	}
}
