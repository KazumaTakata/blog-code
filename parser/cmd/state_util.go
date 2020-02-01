package main

import (
	"parser/util"
	//	"strings"
)

type not_explored_queue struct {
	queue []int
}

func (q *not_explored_queue) enqueue(new_item int) {
	q.queue = append(q.queue, new_item)
}

func (q *not_explored_queue) dequeue() int {
	dequeued := q.queue[0]
	q.queue = q.queue[1:]
	return dequeued
}

func (q *not_explored_queue) empty() bool {
	if len(q.queue) == 0 {
		return true
	}

	return false
}

func is_last(state_element State_element, bnf_list []util.Bnf) bool {
	if len(bnf_list[state_element.Product_id].Right[state_element.Alternate_id]) > state_element.Offset {
		return false
	}

	return true
}

type State map[State_element]bool

type State_element struct {
	Product_id   int
	Alternate_id int
	Offset       int
}
