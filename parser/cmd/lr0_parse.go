package main

import (
	"fmt"
	//	"remove_left_recursion/ll"
	"parser/util"
	//	"strings"
)

type state_stack struct {
	data []int
}

func (s *state_stack) pop() int {
	top := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return top
}

func (s *state_stack) push(d int) {
	s.data = append(s.data, d)
}

func (s *state_stack) top() int {
	return s.data[len(s.data)-1]
}

type symbol_stack struct {
	data []node
}

func (s *symbol_stack) pop() node {
	top := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return top
}

func (s *symbol_stack) push(d node) {
	s.data = append(s.data, d)
}

func (s *symbol_stack) top() node {
	return s.data[len(s.data)-1]
}

type node struct {
	node_type string
	children  []node
}

func parse_lr0(automaton_states []State_with_next, input_tokens []string, bnf_list []util.Bnf) symbol_stack {

	state_stack := state_stack{}
	symbol_stack := symbol_stack{}

	state_stack.push(0)

	for _, token := range input_tokens {
		symbol_stack.push(node{node_type: token})
		next_state_id := automaton_states[state_stack.top()].next[symbol_stack.top().node_type]
		state_stack.push(next_state_id)
		handle_reduction(automaton_states, next_state_id, bnf_list, &symbol_stack, &state_stack)

	}

	return symbol_stack

}

func handle_reduction(automaton_states []State_with_next, next_state_id int, bnf_list []util.Bnf, symbol_stack *symbol_stack, state_stack *state_stack) {

	next_state := automaton_states[next_state_id]
	handlers := get_handlers(next_state, bnf_list)
	if len(handlers) > 0 {
		right := bnf_list[handlers[0].Product_id].Right[handlers[0].Alternate_id]
		left := bnf_list[handlers[0].Product_id].Left
		root_node := node{node_type: left, children: []node{}}

		for i := len(right) - 1; i >= 0; i-- {
			poped := symbol_stack.pop()
			if poped.node_type == right[i] {
				root_node.children = append([]node{poped}, root_node.children...)
				state_stack.pop()
			} else {
				fmt.Printf("parse error")
			}
		}
		symbol_stack.push(root_node)
		next_state_id = automaton_states[state_stack.top()].next[symbol_stack.top().node_type]
		state_stack.push(next_state_id)
		handle_reduction(automaton_states, next_state_id, bnf_list, symbol_stack, state_stack)
	}

}
func get_handlers(state_with_next State_with_next, bnf_list []util.Bnf) []State_element {

	state_elements := []State_element{}

	for state_element, _ := range state_with_next.state {
		if state_element.Offset >= len(bnf_list[state_element.Product_id].Right[state_element.Alternate_id]) {
			state_elements = append(state_elements, state_element)
		}
	}

	return state_elements

}
