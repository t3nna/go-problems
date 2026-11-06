package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinHeapInsertAndDeleteMin(t *testing.T) {
	h := &minHeap{}
	assert.Nil(t, h.peek())
	assert.Nil(t, h.deleteMin())

	late := &entry{key: "late", expiresAt: time.Unix(30, 0)}
	early := &entry{key: "early", expiresAt: time.Unix(10, 0)}
	middle := &entry{key: "middle", expiresAt: time.Unix(20, 0)}
	equal := &entry{key: "equal", expiresAt: time.Unix(10, 0)}
	h.insert(late)
	h.insert(early)
	h.insert(middle)
	h.insert(equal)

	assert.Same(t, early, h.peek())
	assert.Len(t, h.h, 4)

	first := h.deleteMin()
	second := h.deleteMin()
	require.NotNil(t, first)
	require.NotNil(t, second)
	assert.ElementsMatch(t, []*entry{early, equal}, []*entry{first, second})
	assert.Equal(t, -1, first.index)
	assert.Equal(t, -1, second.index)
	assert.Same(t, middle, h.deleteMin())
	assert.Equal(t, -1, middle.index)
	assert.Same(t, late, h.deleteMin())
	assert.Equal(t, -1, late.index)
	assert.Nil(t, h.peek())
	assert.Nil(t, h.deleteMin())
}

func TestMinHeapRemove(t *testing.T) {
	// Removing an interior entry can move the last entry down.
	h := &minHeap{}
	one := &entry{expiresAt: time.Unix(1, 0)}
	two := &entry{expiresAt: time.Unix(2, 0)}
	three := &entry{expiresAt: time.Unix(3, 0)}
	four := &entry{expiresAt: time.Unix(4, 0)}
	five := &entry{expiresAt: time.Unix(5, 0)}
	h.insert(one)
	h.insert(two)
	h.insert(three)
	h.insert(four)
	h.insert(five)
	backing := h.h

	assert.Same(t, two, h.remove(two.index))
	assert.Equal(t, -1, two.index)
	assert.Nil(t, backing[4])
	assert.Equal(t, 1, four.index)
	assert.Equal(t, 3, five.index)
	assert.Same(t, one, h.deleteMin())
	assert.Same(t, three, h.deleteMin())
	assert.Same(t, four, h.deleteMin())
	assert.Same(t, five, h.deleteMin())

	// Removing an interior entry can also move the last entry up.
	h = &minHeap{}
	one = &entry{expiresAt: time.Unix(1, 0)}
	four = &entry{expiresAt: time.Unix(4, 0)}
	two = &entry{expiresAt: time.Unix(2, 0)}
	five = &entry{expiresAt: time.Unix(5, 0)}
	six := &entry{expiresAt: time.Unix(6, 0)}
	three = &entry{expiresAt: time.Unix(3, 0)}
	h.insert(one)
	h.insert(four)
	h.insert(two)
	h.insert(five)
	h.insert(six)
	h.insert(three)

	assert.Same(t, five, h.remove(five.index))
	assert.Equal(t, -1, five.index)
	assert.Equal(t, 1, three.index)
	assert.Equal(t, 3, four.index)
	assert.Same(t, one, h.deleteMin())
	assert.Same(t, two, h.deleteMin())
	assert.Same(t, three, h.deleteMin())
	assert.Same(t, four, h.deleteMin())
	assert.Same(t, six, h.deleteMin())
}

func TestMinHeapRemoveRootAndLast(t *testing.T) {
	h := &minHeap{}
	first := &entry{expiresAt: time.Unix(1, 0)}
	last := &entry{expiresAt: time.Unix(2, 0)}
	h.insert(first)
	h.insert(last)

	assert.Same(t, last, h.remove(last.index))
	assert.Equal(t, -1, last.index)
	assert.Same(t, first, h.remove(first.index))
	assert.Equal(t, -1, first.index)
	assert.Empty(t, h.h)

	h.insert(first)
	h.insert(last)
	assert.Same(t, first, h.remove(0))
	assert.Equal(t, -1, first.index)
	assert.Equal(t, 0, last.index)
	assert.Same(t, last, h.peek())
}
