package main

import "fmt"

type Entry struct {
	Name string
	Prev *Entry
	Next *Entry
}

var (
	head *Entry
	tail *Entry
)

var (
	cacheMap   map[string]*Entry
	listLength = 3
)

func main() {
	cacheMap = make(map[string]*Entry)

	testEntries := []string{"A", "B", "C", "A", "A", "B", "F", "G", "C", "I", "J"}

	head = &Entry{Name: "dummy_head"}
	tail = &Entry{Name: "dummy_tail"}
	head.Next = tail
	tail.Prev = head

	fmt.Println("--- STARTING CACHE SIMULATION ---")
	printList()
	fmt.Println("---------------------------------")

	// Run our test loop
	for i, entryName := range testEntries {
		isNew := allocateEntry(entryName)

		status := "MISS (Inserted)"
		if !isNew {
			status = "HIT (Moved to Top)"
		}

		fmt.Printf("Index %d '%s' | Status: %s\n", i, entryName, status)

		// 🚨 WARNING: The simulation will print perfectly until it hits index 3 ("A").
		// Once it hits "A", the missing detach step will cause an infinite loop here.
		printList()
		fmt.Println("---------------------------------")
	}
}

func printList() {
	curr := head.Next
	fmt.Print("Current Timeline: [Head] <-> ")

	// Safe counter to prevent your terminal from crashing during infinite loops
	safetyCounter := 0

	for curr != tail {
		fmt.Printf("%s <-> ", curr.Name)
		curr = curr.Next

		safetyCounter++
		if safetyCounter > 10 {
			fmt.Println("💥 INFINITE LOOP DETECTED! Stopping print.")
			return
		}
	}
	fmt.Println("[Tail]")
}

func allocateEntry(e string) bool {

	// check entry existence
	if entry, exists := cacheMap[e]; exists {
		currentHead := head.Next
		if entry == currentHead {
			// if exists and is already at the top of the list, do nothing and return false
			return false
		}
		// if exists, move to the top of the list while evicting the entry from its current position and return false
		prevNode := entry.Prev
		prevNode.Next = entry.Next

		nextNode := entry.Next
		nextNode.Prev = prevNode

		entry.Prev = head
		entry.Next = currentHead

		currentHead.Prev = entry
		head.Next = entry

		return false
	} else if len(cacheMap) < listLength {
		currentHead := head.Next
		// if not exists and list is not at cap, append to the top of the list and return true
		newEntry := &Entry{Name: e}

		currentHead.Prev = newEntry
		newEntry.Next = currentHead
		newEntry.Prev = head
		head.Next = newEntry

		cacheMap[e] = newEntry
		return true
	} else if len(cacheMap) == listLength { // cap reached
		currentHead := head.Next
		// get node before the last node (tail.Prev.Prev)
		nodeBeforeLast := tail.Prev.Prev
		evictingEntry := tail.Prev
		// evict said node from the map
		delete(cacheMap, evictingEntry.Name)
		// evict said node from the list
		nodeBeforeLast.Next = tail
		tail.Prev = nodeBeforeLast

		// append to the top of the list
		newEntry := &Entry{Name: e}
		// current head
		currentHead.Prev = newEntry

		newEntry.Next = currentHead
		newEntry.Prev = head

		head.Next = newEntry
		// entry to map
		cacheMap[e] = newEntry
		return true
	}
	return false
}
