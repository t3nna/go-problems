package main

import "time"

type entry struct {
	key       string
	value     string
	expiresAt time.Time
	index     int // for min heap, and ttl cache
}

type minHeap struct {
	h []*entry
}

func (h *minHeap) length() int {
	return len(h.h)
}

func (h *minHeap) peek() *entry {
	if len(h.h) == 0 {
		return nil
	}
	return h.h[0]
}

func (h *minHeap) insert(e *entry) {
	e.index = len(h.h)
	h.h = append(h.h, e)
	h.heapifyUp(e.index)

}
func (h *minHeap) deleteMin() *entry {
	if len(h.h) == 0 {
		return nil
	}
	return h.remove(0)
}

func (h *minHeap) remove(idx int) *entry {
	lastIdx := len(h.h) - 1
	removed := h.h[idx]

	h.swap(idx, lastIdx)
	h.h[lastIdx] = nil
	h.h = h.h[:lastIdx]
	removed.index = -1

	if idx < len(h.h) {
		p := h.parent(idx)

		if idx > 0 && h.less(idx, p) {
			h.heapifyUp(idx)
		} else {
			h.heapifyDown(idx)
		}

	}
	return removed

}

func (h *minHeap) parent(idx int) int {
	if idx <= 0 {
		return 0
	}
	return (idx - 1) / 2

}

func (h *minHeap) rightChild(idx int) int {
	return 2*idx + 2
}

func (h *minHeap) leftChild(idx int) int {
	return 2*idx + 1
}

func (h *minHeap) heapifyUp(idx int) {
	for idx > 0 {
		p := h.parent(idx)
		if !h.less(idx, p) {
			return
		}
		h.swap(idx, p)
		idx = p
	}

}
func (h *minHeap) heapifyDown(idx int) {
	for {
		l := h.leftChild(idx)
		r := h.rightChild(idx)
		smallest := idx

		if l < len(h.h) && h.less(l, smallest) {
			smallest = l
		}
		if r < len(h.h) && h.less(r, smallest) {
			smallest = r
		}

		if smallest == idx {
			return
		}

		h.swap(idx, smallest)
		idx = smallest

	}

}

func (h *minHeap) less(i, j int) bool {
	return h.h[i].expiresAt.Before(h.h[j].expiresAt)
}

func (h *minHeap) swap(i, j int) {
	h.h[i], h.h[j] = h.h[j], h.h[i]
	h.h[i].index = i
	h.h[j].index = j
}
