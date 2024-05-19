package main

import (
	"context"
	"fmt"
	"time"
)

func operation(ctx context.Context, sleepTime int) {
	completed := make(chan bool)

	go func() {
		time.Sleep(time.Duration(sleepTime) * time.Second) // Simuliert eine lange Operation
		completed <- true
	}()

	select {
	case <-ctx.Done(): // Wird ausgeführt, wenn der Kontext abgebrochen wird
		fmt.Println("Operation aborted")
	case <-completed:
		fmt.Println("Operation completed successfully")
	}
}

type contextKey string // define a type to avoid key collisions

func greet(ctx context.Context) {
	// Retrieve the value from the context
	user := ctx.Value(contextKey("user")).(string) // type assertion is necessary
	fmt.Printf("Hello, %s!\n", user)
}

func main() {
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	// _time := time.Now().Add(time.Duration(time.Second * 5))
	// ctx, cancel := context.WithDeadline(context.Background(), _time)

	var userKey contextKey = "user"
	ctx := context.WithValue(context.Background(), userKey, "Alice123")

	// defer cancel() // Stellt sicher, dass alle Ressourcen freigegeben werden, wenn sie nicht mehr benötigt werden

	operation(ctx, 1)
	operation(ctx, 2)
	operation(ctx, 3)

	greet(ctx)
}
