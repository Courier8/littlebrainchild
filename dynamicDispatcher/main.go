package main

import "fmt"

type Rider struct {
	ID    string
	Score uint64
}

var riderPool []*Rider

func main() {
	// A messy input list of riders arriving in random order
	mockRiders := []*Rider{
		{ID: "Rider-A", Score: 40},
		{ID: "Rider-B", Score: 12},
		{ID: "Rider-C", Score: 89},
		{ID: "Rider-D", Score: 34},
		{ID: "Rider-E", Score: 19},
		{ID: "Rider-F", Score: 9},
		{ID: "Rider-G", Score: 3},
		{ID: "Rider-H", Score: 16},
	}

	fmt.Println("--- PHASE 1: SYSTEM INGESTION (SIFT-UP) ---")
	for _, r := range mockRiders {
		addRiderSiftUp(r)
		fmt.Printf("Added %s (Score: %d) | Current Root is: %s (Score: %d)\n",
			r.ID, r.Score, riderPool[0].ID, riderPool[0].Score)
	}

	fmt.Println("\n--- PHASE 2: DISPATCH SIMULATION (SIFT-DOWN) ---")
	// Expected extraction order: 3, 9, 12, 16, 19, 34
	expectedScores := []uint64{3, 9, 12, 16, 19, 34}

	for i := 0; i < len(expectedScores); i++ {
		bestRider := topTrim()
		if bestRider == nil {
			fmt.Println("❌ Error: Pool emptied prematurely!")
			return
		}

		if bestRider.Score == expectedScores[i] {
			fmt.Printf("✅ MATCH: Successfully dispatched %s with lowest score: %d\n", bestRider.ID, bestRider.Score)
		} else {
			fmt.Printf("❌ MISMATCH: Expected score %d, but got %s with score %d\n", expectedScores[i], bestRider.ID, bestRider.Score)
		}
	}
}

func addRiderSiftUp(r *Rider) {
	// 1. Append r to riderPool
	riderPool = append(riderPool, r)
	newRiderIdx := len(riderPool) - 1
	safetyCounter := 0

	// 2. Loop and compare to parent index: (idx - 1) / 2
	for newRiderIdx > 0 {
		parentIdx := (newRiderIdx - 1) / 2

		safetyCounter++
		if safetyCounter > 500 {
			fmt.Println("💥 SIFT-UP HARD-BRAKE TRIGGERED! Breaking loop to protect CPU.")
			break
		}
		// 3. Swap if smaller, then update idx
		if riderPool[newRiderIdx].Score < riderPool[parentIdx].Score {
			riderPool[newRiderIdx], riderPool[parentIdx] = riderPool[parentIdx], riderPool[newRiderIdx]
			newRiderIdx = parentIdx
		} else {
			// if not swapping, must stop
			break
		}
	}
}

func topTrim() *Rider {
	if len(riderPool) == 0 {
		return nil
	}

	// 1. Keep a reference to the target rider at index 0
	rider := riderPool[0]
	// 2. Swap index 0 with the last index: len(riderPool) - 1
	riderPool[0], riderPool[len(riderPool)-1] = riderPool[len(riderPool)-1], riderPool[0]
	// 3. Shrink the slice:
	riderPool = riderPool[:len(riderPool)-1]

	// 4. If there are still items left in the pool, call siftDownReorderPool()
	if len(riderPool) > 0 {
		siftDownReorderPool()
	}

	// 5. Return the target rider reference
	return rider
}

func siftDownReorderPool() {
	siftIdx := 0
	safetyCounter := 0
	// Loop while left child (2*idx + 1) < len(riderPool)
	for leftChildIdx := 2*siftIdx + 1; leftChildIdx < len(riderPool); leftChildIdx = 2*siftIdx + 1 {

		// 🚨 HARD-BRAKE SAFEGUARD
		safetyCounter++
		if safetyCounter > 500 {
			fmt.Println("💥 SIFT-DOWN HARD-BRAKE TRIGGERED! Breaking loop to protect CPU.")
			break
		}

		smallest := leftChildIdx
		rightChildIdx := leftChildIdx + 1
		//   Find the smallest of the two children
		if rightChildIdx < len(riderPool) && riderPool[rightChildIdx].Score < riderPool[leftChildIdx].Score {
			smallest = rightChildIdx
		}

		if riderPool[siftIdx].Score > riderPool[smallest].Score {
			//   If parent is larger than smallest child, swap and update idx
			riderPool[siftIdx], riderPool[smallest] = riderPool[smallest], riderPool[siftIdx]
			siftIdx = smallest
		} else {
			// break because the new node has found its correct position
			break
		}
	}
}
