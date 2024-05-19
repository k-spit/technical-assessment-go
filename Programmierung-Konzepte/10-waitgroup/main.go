package main

import (
	"fmt"
	"sync"
	"time"
)

func performTask(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Task %d starting\n", id)
	if id == 2 { // Simuliere einen Fehler in Task 2
		fmt.Printf("Task %d encountered an error!\n", id)
		return // Frühes Beenden führt zu einem sofortigen Aufruf von wg.Done() durch defer
	}
	time.Sleep(time.Second)
	fmt.Printf("Task %d completed\n", id)
}

func main() {
	var wg sync.WaitGroup
	numTasks := 5
	wg.Add(numTasks)

	for i := 0; i < numTasks; i++ {
		go performTask(i, &wg)
	}

	wg.Wait()
	fmt.Println("All tasks attempted.")
}
