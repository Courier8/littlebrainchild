package main

import "fmt"

type Bucket struct {
	Timestamp int
	Count     int
}

var (
	capCount   = 3
	windowSize = 5
	timeWindow = [5]Bucket{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}
)

func checkOrderAllowance(timestamp int) bool {

	currentCount := 0
	// calculate the position in the circular buffer, using modulo operation.
	pos := timestamp % windowSize
	// check if the current bucket expired
	if timeWindow[pos].Timestamp != timestamp {
		timeWindow[pos].Timestamp = timestamp
		timeWindow[pos].Count = 0
	}

	// loop through the time window and count the orders
	for i := 0; i < windowSize; i++ {
		if timeWindow[i].Timestamp > timestamp-windowSize {
			currentCount += timeWindow[i].Count
		}
	}

	// cap check
	if currentCount < capCount {
		timeWindow[pos].Timestamp = timestamp
		timeWindow[pos].Count++
		return true
	}

	return false
}

func main() {
	fmt.Println("🚀 Starting order processing simulation with sliding window rate limiter...")

	orderSchedule := []int{1, 3, 5, 6, 7, 8, 9, 10, 11, 14, 15, 16, 18, 19, 20}

	for _, t := range orderSchedule {
		fmt.Printf("Time %ds -> Allowed: %t\n", t, checkOrderAllowance(t))
	}

}
