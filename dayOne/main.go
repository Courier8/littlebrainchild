package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Simulator Config
const concurrentOrders = 50000
const useMemoryProtection = true // 🚨 TOGGLE THIS TO FIX THE RACE CONDITION

type SalesTracker struct {
	mu           sync.RWMutex
	menuItemsOld map[string]int // Unprotected map

	// Protected metrics
	globalOrderCount atomic.Int64
	menuItemsSafe    map[string]int
}

func main() {
	tracker := &SalesTracker{
		menuItemsOld:  make(map[string]int),
		menuItemsSafe: make(map[string]int),
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	fmt.Printf("🚀 Starting peak-hour simulation with %d concurrent requests...\n", concurrentOrders)

	for i := 0; i < concurrentOrders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Simulate a customer ordering a "Midnight Burger"
			if useMemoryProtection {
				// 🛡️ Safe Approach
				tracker.globalOrderCount.Add(1) // Lock-free atomic increment

				tracker.mu.Lock()
				tracker.menuItemsSafe["Midnight Burger"]++
				tracker.mu.Unlock()
			} else {
				// 💥 Dangerous Approach (Data Race)
				// tracker.globalOrderCount.Add(1) // (pretend this is line below)
				// If we used a regular integer here, it would race too.
				tracker.menuItemsOld["Midnight Burger"]++
			}
		}()
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Println("\n--- Simulation Results ---")
	fmt.Printf("Execution Time: %v\n", duration)

	if useMemoryProtection {
		fmt.Printf("Safe Global Counter: %d\n", tracker.globalOrderCount.Load())
		fmt.Printf("Safe Map Counter:   %d\n", tracker.menuItemsSafe["Midnight Burger"])
	} else {
		fmt.Printf("Unsafe Map Counter: %d\n", tracker.menuItemsOld["Midnight Burger"])
	}
}
