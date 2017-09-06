package example

import "github.com/cheekybits/genny/generic"

type Generic generic.Type

// GenericQueue represents a queue of Generic types.
type GenericQueue struct {
	items []Generic
}

// NewGenericQueue makes a new empty Generic queue.
func NewGenericQueue() *GenericQueue {
	return &GenericQueue{items: make([]Generic, 0)}
}

// Enq adds an item to the queue.
func (q *GenericQueue) Enq(obj Generic) *GenericQueue {
	q.items = append(q.items, obj)
	return q
}

// Deq removes and returns the next item in the queue.
func (q *GenericQueue) Deq() Generic {
	obj := q.items[0]
	q.items = q.items[1:]
	return obj
}

// Len gets the current number of Generic items in the queue.
func (q *GenericQueue) Len() int {
	return len(q.items)
}
