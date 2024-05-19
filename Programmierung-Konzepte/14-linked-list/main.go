package main

import (
	"fmt"
)

// Node represents a node in the singly linked list.
type Node struct {
	data int
	next *Node
}

// LinkedList represents a singly linked list.
type LinkedList struct {
	head *Node
}

// Insert adds a new node at the end of the list.
func (ll *LinkedList) Insert(data int) {
	newNode := &Node{data: data}
	if ll.head == nil {
		ll.head = newNode
		return
	}
	current := ll.head
	for current.next != nil {
		current = current.next
	}
	current.next = newNode
}

// Delete removes a node with the given data.
func (ll *LinkedList) Delete(data int) {
	if ll.head == nil {
		return
	}
	if ll.head.data == data {
		ll.head = ll.head.next
		return
	}
	current := ll.head
	for current.next != nil && current.next.data != data {
		current = current.next
	}
	if current.next != nil {
		current.next = current.next.next
	}
}

// Display prints the contents of the list.
func (ll *LinkedList) Display() {
	current := ll.head
	for current != nil {
		fmt.Printf("%d -> ", current.data)
		current = current.next
	}
	fmt.Println("nil")
}

func main() {
	ll := &LinkedList{}
	ll.Insert(1)
	ll.Insert(2)
	ll.Insert(3)
	ll.Display() // Output: 1 -> 2 -> 3 -> nil

	ll.Delete(2)
	ll.Display() // Output: 1 -> 3 -> nil
}
