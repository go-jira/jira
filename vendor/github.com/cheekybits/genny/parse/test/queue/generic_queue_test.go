package queue

import "testing"

func TestSomethingQueue(t *testing.T) {

	var item1 *Something = new(Something)
	var item2 *Something = new(Something)
	var item3 *Something = new(Something)

	q := NewSomethingQueue()

	q.Push(item1)
	if len(q.items) != 1 {
		t.Error("Push should add the item")
	}
	q.Push(item2)
	if len(q.items) != 2 {
		t.Error("Push should add the item")
	}
	q.Push(item3)
	if len(q.items) != 3 {
		t.Error("Push should add the item")
	}

	if q.Pop() != item1 {
		t.Error("Pop should return items in the order in which they were added")
	}
	if len(q.items) != 2 {
		t.Error("Pop should remove the item")
	}
	if q.Pop() != item2 {
		t.Error("Pop should return items in the order in which they were added")
	}
	if len(q.items) != 1 {
		t.Error("Pop should remove the item")
	}
	if q.Pop() != item3 {
		t.Error("Pop should return items in the order in which they were added")
	}
	if len(q.items) != 0 {
		t.Error("Pop should remove the item")
	}

}

func TestIntQueue(t *testing.T) {

	var item1 int = 1
	var item2 int = 2
	var item3 int = 3

	q := NewIntQueue()

	q.Push(item1)
	if len(q.items) != 1 {
		t.Error("Push should add the item")
	}
	q.Push(item2)
	if len(q.items) != 2 {
		t.Error("Push should add the item")
	}
	q.Push(item3)
	if len(q.items) != 3 {
		t.Error("Push should add the item")
	}

	if q.Pop() != item1 {
		t.Error("Pop should return items in the order in which they were added")
	}
	if len(q.items) != 2 {
		t.Error("Pop should remove the item")
	}
	if q.Pop() != item2 {
		t.Error("Pop should return items in the order in which they were added")
	}
	if len(q.items) != 1 {
		t.Error("Pop should remove the item")
	}
	if q.Pop() != item3 {
		t.Error("Pop should return items in the order in which they were added")
	}
	if len(q.items) != 0 {
		t.Error("Pop should remove the item")
	}

}
