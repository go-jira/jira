package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	q := NewGenericQueue()
	assert.NotNil(t, q)

}

func TestEnqueueAndDequeue(t *testing.T) {

	item1 := new(Generic)
	item2 := new(Generic)
	q := NewGenericQueue()

	assert.Equal(t, q, q.Enq(item1), "Enq should return the queue")
	assert.Equal(t, 1, q.Len())
	assert.Equal(t, q, q.Enq(item2), "Enq should return the queue")
	assert.Equal(t, 2, q.Len())

	assert.Equal(t, item1, q.Deq())
	assert.Equal(t, 1, q.Len())
	assert.Equal(t, item2, q.Deq())
	assert.Equal(t, 0, q.Len())

}
