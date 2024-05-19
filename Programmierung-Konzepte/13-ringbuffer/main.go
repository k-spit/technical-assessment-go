package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// RingBuffer represents a circular buffer.
type RingBuffer struct {
	buffer []interface{}
	size   int
	start  int
	end    int
	count  int
}

// NewRingBuffer creates a new ring buffer with the given size.
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]interface{}, size),
		size:   size,
	}
}

// Push adds an element to the ring buffer.
func (rb *RingBuffer) Push(value interface{}) {
	if rb.count == rb.size {
		// Overwrite the oldest element when the buffer is full
		rb.start = (rb.start + 1) % rb.size
	} else {
		rb.count++
	}
	rb.buffer[rb.end] = value
	rb.end = (rb.end + 1) % rb.size
}

// Pop removes and returns the oldest element from the ring buffer.
func (rb *RingBuffer) Pop() (interface{}, error) {
	if rb.count == 0 {
		return nil, errors.New("buffer is empty")
	}
	value := rb.buffer[rb.start]
	rb.start = (rb.start + 1) % rb.size
	rb.count--
	return value, nil
}

// Peek returns the oldest element without removing it.
func (rb *RingBuffer) Peek() (interface{}, error) {
	if rb.count == 0 {
		return nil, errors.New("buffer is empty")
	}
	return rb.buffer[rb.start], nil
}

// IsEmpty checks if the buffer is empty.
func (rb *RingBuffer) IsEmpty() bool {
	return rb.count == 0
}

// IsFull checks if the buffer is full.
func (rb *RingBuffer) IsFull() bool {
	return rb.count == rb.size
}

// Size returns the current number of elements in the buffer.
func (rb *RingBuffer) Size() int {
	return rb.count
}

// Capacity returns the total capacity of the buffer.
func (rb *RingBuffer) Capacity() int {
	return rb.size
}

// PrintBuffer prints the current state of the ring buffer.
func (rb *RingBuffer) PrintBuffer() {
	fmt.Print("Buffer: ")
	for i := 0; i < rb.size; i++ {
		if i >= rb.start && i < rb.start+rb.count || (rb.start+rb.count > rb.size && i < (rb.start+rb.count)%rb.size) {
			fmt.Printf("%v ", rb.buffer[i])
		} else {
			fmt.Print("nil ")
		}
	}
	fmt.Println()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the size of the ring buffer: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	size, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid size")
		return
	}

	rb := NewRingBuffer(size)

	for {
		fmt.Println("\nChoose an operation:")
		fmt.Println("1. Push")
		fmt.Println("2. Pop")
		fmt.Println("3. Peek")
		fmt.Println("4. Is Empty")
		fmt.Println("5. Is Full")
		fmt.Println("6. Size")
		fmt.Println("7. Capacity")
		fmt.Println("8. Print Buffer")
		fmt.Println("9. Exit")
		fmt.Print("Enter your choice: ")

		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid choice")
			continue
		}

		switch choice {
		case 1:
			fmt.Print("Enter a value to push: ")
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)
			rb.Push(value)
			fmt.Println("Value pushed")
		case 2:
			value, err := rb.Pop()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Popped value: %v\n", value)
			}
		case 3:
			value, err := rb.Peek()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Peeked value: %v\n", value)
			}
		case 4:
			fmt.Printf("Is Empty: %v\n", rb.IsEmpty())
		case 5:
			fmt.Printf("Is Full: %v\n", rb.IsFull())
		case 6:
			fmt.Printf("Size: %d\n", rb.Size())
		case 7:
			fmt.Printf("Capacity: %d\n", rb.Capacity())
		case 8:
			rb.PrintBuffer()
		case 9:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}
