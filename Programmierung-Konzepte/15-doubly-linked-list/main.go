package main

import (
	"fmt"
)

// Node represents a node in the doubly linked list.
type Node struct {
	data int
	prev *Node
	next *Node
}

// DoublyLinkedList represents a doubly linked list.
type DoublyLinkedList struct {
	head *Node
	tail *Node
}

// Insert adds a new node at the end of the list.
func (dll *DoublyLinkedList) Insert(data int) {
	newNode := &Node{data: data}
	if dll.head == nil {
		dll.head = newNode
		dll.tail = newNode
		return
	}
	dll.tail.next = newNode
	newNode.prev = dll.tail
	dll.tail = newNode
}

// Delete removes a node with the given data.
func (dll *DoublyLinkedList) Delete(data int) {
	if dll.head == nil {
		return
	}
	if dll.head.data == data {
		dll.head = dll.head.next
		if dll.head != nil {
			dll.head.prev = nil
		} else {
			dll.tail = nil
		}
		return
	}
	current := dll.head
	for current != nil && current.data != data {
		current = current.next
	}
	if current != nil {
		if current.next != nil {
			current.next.prev = current.prev
		} else {
			dll.tail = current.prev
		}
		if current.prev != nil {
			current.prev.next = current.next
		}
	}
}

// DisplayForward prints the contents of the list from head to tail.
func (dll *DoublyLinkedList) DisplayForward() {
	current := dll.head
	for current != nil {
		fmt.Printf("%d -> ", current.data)
		current = current.next
	}
	fmt.Println("nil")
}

// DisplayBackward prints the contents of the list from tail to head.
func (dll *DoublyLinkedList) DisplayBackward() {
	current := dll.tail
	for current != nil {
		fmt.Printf("%d -> ", current.data)
		current = current.prev
	}
	fmt.Println("nil")
}

func main() {
	dll := &DoublyLinkedList{}
	dll.Insert(1)
	dll.Insert(2)
	dll.Insert(3)
	dll.DisplayForward()  // Output: 1 -> 2 -> 3 -> nil
	dll.DisplayBackward() // Output: 3 -> 2 -> 1 -> nil

	dll.Delete(2)
	dll.DisplayForward()  // Output: 1 -> 3 -> nil
	dll.DisplayBackward() // Output: 3 -> 1 -> nil
}
