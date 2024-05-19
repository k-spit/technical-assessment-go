package main

import (
	"fmt"
)

// Max returns the maximum value in a slice of any ordered type.
func Max[T any](slice []T, less func(a, b T) bool) T {
	if len(slice) == 0 {
		panic("slice is empty")
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if less(max, v) {
			max = v
		}
	}
	return max
}

// LessInt compares two integers and returns true if the first is less than the second.
func LessInt(a, b int) bool {
	return a < b
}

// LessFloat compares two float64 values and returns true if the first is less than the second.
func LessFloat(a, b float64) bool {
	return a < b
}

// LessString compares two strings and returns true if the first is less than the second.
func LessString(a, b string) bool {
	return a < b
}

func main() {
	ints := []int{1, 2, 3, 4, 5}
	fmt.Printf("Max int: %d\n", Max(ints, LessInt))

	floats := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	fmt.Printf("Max float: %f\n", Max(floats, LessFloat))

	strings := []string{"apple", "banana", "cherry"}
	fmt.Printf("Max string: %s\n", Max(strings, LessString))
}
