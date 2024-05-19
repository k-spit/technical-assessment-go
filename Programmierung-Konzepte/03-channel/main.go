package main

import (
	"fmt"
	"time"
)

func count(thing string, c chan string) {
	for i := 1; i <= 5; i++ {
		c <- thing // sende 'thing' an den Channel
		time.Sleep(time.Millisecond * 500)
	}
	close(c) // schließe den Channel nach Fertigstellung
}

func main() {
	c := make(chan string)
	go count("Schaf", c)
	for msg := range c {
		fmt.Println(msg)
	}
}
