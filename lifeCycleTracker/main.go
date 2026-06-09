package main

import "fmt"

func simulateCells(ticks int) int {
	// Our fixed-size bucket array representing ages 0 to 4
	// Index represents the age, Value represents the count of cells at that age.
	// T = 0 starts with one cell at age 0: [1, 0, 0, 0, 0]
	current := [5]int{1, 0, 0, 0, 0}

	// Run the simulation for T ticks
	for tick := 1; tick <= ticks; tick++ {
		var next [5]int

		// 1. Spawning Logic:
		// Cells at age 1 will become age 2 -> spawns a 0
		// Cells at age 3 will become age 4 -> spawns a 0
		next[0] = current[1] + current[3]

		// 2. Shifting/Aging Logic:
		next[1] = current[0] // 0s become 1s
		next[2] = current[1] // 1s become 2s
		next[3] = current[2] // 2s become 3s
		next[4] = current[3] // 3s become 4s

		// Note: current[4] is completely ignored here,
		// which automatically satisfies the rule: "If age hits 5, eliminate it."

		// Move to the next state
		current = next
	}

	// Sum up all the active cells left in the buckets
	totalCells := 0
	for _, count := range current {
		totalCells += count
	}

	return totalCells
}

func main() {
	// Let's verify against your exact examples
	for t := 0; t <= 1000; t++ {
		fmt.Printf("T = %d | Total Cells: %d\n", t, simulateCells(t))
	}
}
