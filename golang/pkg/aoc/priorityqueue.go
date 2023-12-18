package aoc

import "container/heap"

type Item[T comparable] struct {
	Val         T
	GetPriority func() int
	index       int
}

type PriorityQueue[T comparable] []*Item[T]

func NewPriorityQueue[T comparable]() PriorityQueue[T] {
	pq := PriorityQueue[T]{}
	heap.Init(&pq)
	return pq
}

func (pq *PriorityQueue[T]) pushItem(item *Item[T]) {
	n := len(*pq)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) popItem() *Item[T] {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest priority so we use greater than here.
	return pq[i].GetPriority() < pq[j].GetPriority()
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[T]) Update(item *Item[T]) {
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue[T]) Push(x any) {
	pq.pushItem(x.(*Item[T]))
}

func (pq *PriorityQueue[T]) Pop() any {
	return pq.popItem()
}

func (pq *PriorityQueue[T]) PushItem(item *Item[T]) {
	heap.Push(pq, item)
}

func (pq *PriorityQueue[T]) PopItem() *Item[T] {
	return heap.Pop(pq).(*Item[T])
}
