package main

import (
	"fmt"
	"sync"
)

var (
	mutex   sync.Mutex
	balance int
)

func deposit(value int, wg *sync.WaitGroup) {
	mutex.Lock()
	fmt.Printf("Hinzufügen von %d zum Kontostand\n", value)
	// time.Sleep(time.Millisecond * 500) // Füge eine Verzögerung hinzu
	balance += value
	mutex.Unlock()
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	transactions := []int{50, -20, 70, -30, 100, -50, 30, -10, 60, 10, -5, 15, -25, 40, -10}
	for _, value := range transactions {
		wg.Add(1)
		go deposit(value, &wg)
	}
	wg.Wait()
	fmt.Printf("Endgültiger Kontostand: %d\n", balance)
}
