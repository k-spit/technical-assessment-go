package main

import "fmt"

func main() {
	var name string = "Max"
	fmt.Println("Hallo,", name)
}

func greet(name string) string {
	return "Hallo " + name
}
