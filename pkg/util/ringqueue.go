package util

import (
	. "container/ring"
)

// Queue An unbounded queue
type Queue interface {
	// Offer a new value to the end of this queue
	Offer(interface{})

	// Poll from the head of this queue. if queue is empty, return nil
	Poll() interface{}

	// Len is the number of elements contained
	Len() int

	// Do iterate and apply the function on val
	Do(func(val interface{}))
}

func NewRingQueue(initSize int) Queue {
	r := New(initSize)
	return &ringQueue{
		head: r,
		tail: r,
		size: 0,
	}
}

type ringQueue struct {
	// head 指向空位，用于放置元素
	head *Ring

	// tail 若队列有元素，则指向最早加入的元素
	tail *Ring

	size int
}

func (q *ringQueue) Do(f func(val interface{})) {
	for p := q.tail; p != q.head; p = p.Next() {
		f(p.Value)
	}
}

func (q *ringQueue) Offer(val interface{}) {
	if q.head == q.tail && q.size > 0 {
		// 队列已满，扩容
		q.head = q.head.Prev()
		q.head.Link(New(q.size))
		q.head = q.head.Next()
	}
	q.head.Value = val
	q.head = q.head.Next()
	q.size++
}

func (q *ringQueue) Poll() interface{} {
	if q.size == 0 {
		return nil
	}
	val := q.tail.Value
	q.tail.Value = nil
	q.tail = q.tail.Next()
	q.size--
	return val
}

func (q *ringQueue) Len() int {
	return q.size
}
