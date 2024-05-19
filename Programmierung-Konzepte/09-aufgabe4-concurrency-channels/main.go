package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func generateMessages(msgChan chan<- string, quitChan <-chan bool) {
	for {
		select {
		case <-quitChan:
			log.Println("Nachrichtengenerator wird beendet...")
			return
		default:
			msgChan <- fmt.Sprintf("Nachricht: %d", rand.Intn(100))
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}
}

func consumer(id int, msgChan <-chan string, quitChan <-chan bool) {
	for {
		select {
		case msg := <-msgChan:
			fmt.Printf("Consumer %d erhielt %s\n", id, msg)
		case <-quitChan:
			fmt.Printf("Consumer %d wird beendet...\n", id)
			return
		}
	}
}

func main() {
	msgChan := make(chan string)
	quitChan := make(chan bool)

	go generateMessages(msgChan, quitChan)

	for i := 0; i < 5; i++ {
		go consumer(i, msgChan, quitChan)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	for i := 0; i < 6; i++ {
		quitChan <- true
	}

	fmt.Println("Programm wird beendet...")
}
