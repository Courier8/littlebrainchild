package main

import "fmt"

// Global scales managed by your functions
var minHeap []int
var maxHeap []int

func main() {
	// A live stream of numbers pouring in from a simulated Kafka topic
	streamInput := []int{20, 30, 10, 5, 50, 40}

	fmt.Println("--- STARTING STREAMING MEDIAN CALCULATOR ---")

	for _, val := range streamInput {
		fmt.Printf("\n📥 New Stream Value: %d\n", val)

		// ⚖️ BALANCE COORDINATION: Decision logic to route and level the scales

		// Step 1: Ingest into the correct universe
		// If maxHeap is empty or value belongs in lower half, push to maxHeap (sign=true)
		// Note: We multiply maxHeap[0] by -1 to read its true positive value!
		if len(maxHeap) == 0 || val <= (maxHeap[0]*-1) {
			pushHeap(val, true)
		} else {
			pushHeap(val, false)
		}

		// Step 2: Critical Re-balance execution if scales tilt by more than 1 item
		if len(maxHeap) > len(minHeap)+1 {
			// Left side too heavy: pop from max, convert back to positive, push to min
			poppedVal := popHeap(true)
			pushHeap(poppedVal*-1, false)
		} else if len(minHeap) > len(maxHeap)+1 {
			// Right side too heavy: pop from min, push to max
			poppedVal := popHeap(false)
			pushHeap(poppedVal, true)
		}

		// Step 3: Extract the live Median from the tips of both heaps in O(1) time
		var currentMedian float64
		if len(maxHeap) == len(minHeap) {
			// Even count: Average of both roots (convert max root back to positive)
			maxRootRealValue := maxHeap[0] * -1
			currentMedian = float64(maxRootRealValue+minHeap[0]) / 2.0
		} else if len(maxHeap) > len(minHeap) {
			// Left side heavier: Max-heap root wins
			currentMedian = float64(maxHeap[0] * -1)
		} else {
			// Right side heavier: Min-heap root wins
			currentMedian = float64(minHeap[0])
		}

		// Debug view to trace the inner structure
		fmt.Printf("   Lower Half (MaxHeap Reverted): %v\n", revertMaxHeapView())
		fmt.Printf("   Upper Half (MinHeap):          %v\n", minHeap)
		fmt.Printf("   📊 CURRENT LIVE MEDIAN:       %.1f\n", currentMedian)
	}
}

// Helper function just to print maxHeap values as positive numbers for humans
func revertMaxHeapView() []int {
	view := make([]int, len(maxHeap))
	for i, v := range maxHeap {
		view[i] = v * -1
	}
	return view
}

func pushHeap(val int, sign bool) {
	var heapInUse []int

	if sign {
		heapInUse = maxHeap
		val = val * -1
	} else {
		heapInUse = minHeap
	}
	safetyCounter := 0

	heapInUse = append(heapInUse, val)
	newValIdx := len(heapInUse) - 1

	for newValIdx > 0 {
		parentIdx := (newValIdx - 1) / 2

		safetyCounter++
		if safetyCounter > 500 {
			fmt.Println("💥 SIFT-UP HARD-BRAKE TRIGGERED! Breaking loop to protect CPU.")
			break
		}

		if heapInUse[newValIdx] < heapInUse[parentIdx] {
			heapInUse[newValIdx], heapInUse[parentIdx] = heapInUse[parentIdx], heapInUse[newValIdx]
			newValIdx = parentIdx
		} else {
			break
		}
	}

	if sign {
		maxHeap = heapInUse
	} else {
		minHeap = heapInUse
	}
}

func popHeap(sign bool) (root int) {
	var heapInUse []int
	if sign {
		heapInUse = maxHeap
	} else {
		heapInUse = minHeap
	}

	if len(heapInUse) == 0 {
		return
	}

	root = heapInUse[0]
	heapInUse[0], heapInUse[len(heapInUse)-1] = heapInUse[len(heapInUse)-1], heapInUse[0]
	heapInUse = heapInUse[:len(heapInUse)-1]

	// 💡 THE BUG FIX: Persist slice updates immediately, even if length is now 0
	if sign {
		maxHeap = heapInUse
	} else {
		minHeap = heapInUse
	}

	if len(heapInUse) > 0 {
		reOrderHeap(sign)
	}

	return root
}

func reOrderHeap(sign bool) {
	var heapInUse []int
	if sign {
		heapInUse = maxHeap
	} else {
		heapInUse = minHeap
	}

	siftIdx := 0
	safetyCounter := 0

	for leftChildIdx := 2*siftIdx + 1; leftChildIdx < len(heapInUse); leftChildIdx = 2*siftIdx + 1 {
		safetyCounter++
		if safetyCounter > 500 {
			fmt.Println("💥 SIFT-DOWN HARD-BRAKE TRIGGERED! Breaking loop to protect CPU.")
			break
		}

		smallest := leftChildIdx
		rightChildIdx := leftChildIdx + 1

		if rightChildIdx < len(heapInUse) && heapInUse[rightChildIdx] < heapInUse[leftChildIdx] {
			smallest = rightChildIdx
		}

		if heapInUse[siftIdx] > heapInUse[smallest] {
			heapInUse[siftIdx], heapInUse[smallest] = heapInUse[smallest], heapInUse[siftIdx]
			siftIdx = smallest
		} else {
			break
		}
	}
}
