package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	//	"remove_left_recursion/ll"
	"parser/util"
	"runtime"
	//	"strings"
)

type State_element struct {
	Product_id   int
	Alternate_id int
	Offset       int
}

type State_element_with_follow struct {
	state_element State_element
	follow        string
}

type State map[State_element]bool

type State_with_follow map[State_element_with_follow]bool

type State_with_follow_next struct {
	next  map[string]int
	state State_with_follow
}

type State_with_next struct {
	next  map[string]int
	state State
}

func Expand_non_terminal(state_element State_element, bnf_list []util.Bnf, current_state State) State {

	non_terminals := util.Get_nonterminal(bnf_list)
	added_state := make(State)

	if state_element.Offset < len(bnf_list[state_element.Product_id].Right[state_element.Alternate_id]) {
		right_ele := bnf_list[state_element.Product_id].Right[state_element.Alternate_id][state_element.Offset]
		if _, ok := non_terminals[right_ele]; ok {
			for prod_id, bnf := range bnf_list {
				if bnf.Left == right_ele {
					for alte_id, _ := range bnf.Right {
						new_state_element := State_element{Product_id: prod_id, Alternate_id: alte_id, Offset: 0}

						if _, ok := current_state[new_state_element]; ok {
							return State{}
						}

						added_state[new_state_element] = true
						current_state[new_state_element] = true

						new_state := Expand_non_terminal(new_state_element, bnf_list, current_state)
						for new_ele, _ := range new_state {
							added_state[new_ele] = true
						}
					}
				}
			}
		}

	}

	return added_state
}

func get_follow(state_element State_element, bnf_list []util.Bnf) (bool, string) {

	if state_element.Offset+1 < len(bnf_list[state_element.Product_id].Right[state_element.Alternate_id]) {
		return true, bnf_list[state_element.Product_id].Right[state_element.Alternate_id][state_element.Offset+1]
	}

	return false, ""
}

func Expand_non_terminal_with_follow(state_element State_element, bnf_list []util.Bnf, follow string, current_state State_with_follow) State_with_follow {

	non_terminals := util.Get_nonterminal(bnf_list)
	added_state := make(State_with_follow)

	if state_element.Offset < len(bnf_list[state_element.Product_id].Right[state_element.Alternate_id]) {
		right_ele := bnf_list[state_element.Product_id].Right[state_element.Alternate_id][state_element.Offset]
		if _, ok := non_terminals[right_ele]; ok {
			for prod_id, bnf := range bnf_list {
				if bnf.Left == right_ele {
					for alte_id, _ := range bnf.Right {
						new_state_element := State_element{Product_id: prod_id, Alternate_id: alte_id, Offset: 0}

						exist, new_follow := get_follow(state_element, bnf_list)

						follow_arg := ""
						if exist {
							follow_arg = new_follow
						} else {
							follow_arg = follow
						}

						new_state_element_with_follow := State_element_with_follow{state_element: new_state_element, follow: follow_arg}

						if _, ok := current_state[new_state_element_with_follow]; ok {
							return State_with_follow{}
						}

						added_state[new_state_element_with_follow] = true
						current_state[new_state_element_with_follow] = true

						exist, new_follow = get_follow(new_state_element, bnf_list)

						follow_arg = ""
						if exist {
							follow_arg = new_follow
						} else {
							follow_arg = new_state_element_with_follow.follow
						}
						new_state_with_follow := Expand_non_terminal_with_follow(new_state_element, bnf_list, follow_arg, current_state)

						for new_ele, _ := range new_state_with_follow {
							added_state[new_ele] = true
							current_state[new_ele] = true
						}
					}
				}
			}
		}

	}

	return added_state
}

func is_equal(state_a, state_b State) bool {
	if len(state_a) != len(state_b) {
		return false
	}

	for state_a_ele, _ := range state_a {
		if _, ok := state_b[state_a_ele]; !ok {
			return false
		}
	}

	return true
}

func is_equal_with_follow(state_a, state_b State_with_follow) bool {
	if len(state_a) != len(state_b) {
		return false
	}

	for state_a_ele, _ := range state_a {
		if _, ok := state_b[state_a_ele]; !ok {
			return false
		}
	}

	return true
}

func add_to_automaton_states_with_follow(automaton_state *[]State_with_follow_next, new_state State_with_follow) (bool, int) {

	for index, state := range *automaton_state {
		if is_equal_with_follow(state.state, new_state) {
			return false, index
		}
	}

	*automaton_state = append(*automaton_state, State_with_follow_next{state: new_state, next: make(map[string]int)})

	return true, len(*automaton_state) - 1
}

func add_to_automaton_states(automaton_state *[]State_with_next, new_state State) (bool, int) {

	for index, state := range *automaton_state {
		if is_equal(state.state, new_state) {
			return false, index
		}
	}

	*automaton_state = append(*automaton_state, State_with_next{state: new_state, next: make(map[string]int)})

	return true, len(*automaton_state) - 1
}

func is_last(state_element State_element, bnf_list []util.Bnf) bool {
	if len(bnf_list[state_element.Product_id].Right[state_element.Alternate_id]) > state_element.Offset {
		return false
	}

	return true
}

func create_new_states(bnf_list []util.Bnf, root_state State) map[string]State {

	nonterminal_and_terminal := util.Get_nonterminal_and_terminal(bnf_list)

	new_state_elements := map[string][]State_element{}

	for node, _ := range nonterminal_and_terminal {
		for state_ele, _ := range root_state {
			if !is_last(state_ele, bnf_list) {
				if node == bnf_list[state_ele.Product_id].Right[state_ele.Alternate_id][state_ele.Offset] {
					new_state_elements[node] = append(new_state_elements[node], State_element{Product_id: state_ele.Product_id, Alternate_id: state_ele.Alternate_id, Offset: state_ele.Offset + 1})
				}
			}
		}

	}
	new_states := map[string]State{}

	for key, elements := range new_state_elements {
		new_state := State{}
		for _, element := range elements {
			new_elements := Expand_non_terminal(element, bnf_list, State{})
			for new_ele, _ := range new_elements {
				new_state[new_ele] = true
			}
			new_state[element] = true
		}
		new_states[key] = new_state
	}

	return new_states
}

func create_new_states_with_follow(bnf_list []util.Bnf, root_state State_with_follow) map[string]State_with_follow {

	nonterminal_and_terminal := util.Get_nonterminal_and_terminal(bnf_list)

	new_state_elements := map[string][]State_element_with_follow{}

	for node, _ := range nonterminal_and_terminal {
		for state_element_with_follow, _ := range root_state {
			if !is_last(state_element_with_follow.state_element, bnf_list) {
				if node == bnf_list[state_element_with_follow.state_element.Product_id].Right[state_element_with_follow.state_element.Alternate_id][state_element_with_follow.state_element.Offset] {
					new_state_element_with_follow := State_element_with_follow{state_element: State_element{Product_id: state_element_with_follow.state_element.Product_id, Alternate_id: state_element_with_follow.state_element.Alternate_id, Offset: state_element_with_follow.state_element.Offset + 1}, follow: state_element_with_follow.follow}
					new_state_elements[node] = append(new_state_elements[node], new_state_element_with_follow)
				}
			}
		}

	}

	new_states := map[string]State_with_follow{}

	for key, elements := range new_state_elements {
		new_state := State_with_follow{}
		for _, element := range elements {
			new_elements := Expand_non_terminal_with_follow(element.state_element, bnf_list, element.follow, State_with_follow{})
			for new_ele, _ := range new_elements {
				new_state[new_ele] = true
			}
			new_state[element] = true
		}
		new_states[key] = new_state
	}

	return new_states
}

func add_all_to_automaton_states(automaton_states *[]State_with_next, root_index int, new_states map[string]State) []int {

	not_explored := []int{}
	for key, new_state := range new_states {
		is_new, index := add_to_automaton_states(automaton_states, new_state)
		if is_new {
			not_explored = append(not_explored, index)
		}
		(*automaton_states)[root_index].next[key] = index
	}

	return not_explored
}

func add_all_to_automaton_states_with_follow(automaton_states *[]State_with_follow_next, root_index int, new_states map[string]State_with_follow) []int {

	not_explored := []int{}
	for key, new_state := range new_states {
		is_new, index := add_to_automaton_states_with_follow(automaton_states, new_state)
		if is_new {
			not_explored = append(not_explored, index)
		}
		(*automaton_states)[root_index].next[key] = index
	}

	return not_explored
}

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
func gen_start_state(bnf_list []util.Bnf) State {
	start_state := make(State)
	start_state[State_element{Product_id: 0, Alternate_id: 0, Offset: 0}] = true

	new_start_state := make(State)
	for state_ele, _ := range start_state {
		new_elements := Expand_non_terminal(state_ele, bnf_list, State{})
		for new_ele, _ := range new_elements {
			new_start_state[new_ele] = true
		}
	}
	for new_elem, _ := range new_start_state {
		start_state[new_elem] = true
	}

	return start_state

}
func gen_start_state_with_follow(bnf_list []util.Bnf) State_with_follow {
	start_state := make(State_with_follow)
	start_element := State_element{Product_id: 0, Alternate_id: 0, Offset: 0}
	start_element_follow := State_element_with_follow{state_element: start_element, follow: "$"}
	start_state[start_element_follow] = true

	new_start_state := make(State_with_follow)
	for state_ele, _ := range start_state {
		new_elements := Expand_non_terminal_with_follow(state_ele.state_element, bnf_list, start_element_follow.follow, State_with_follow{})
		for new_ele, _ := range new_elements {
			new_start_state[new_ele] = true
		}
	}
	for new_elem, _ := range new_start_state {
		start_state[new_elem] = true
	}

	return start_state

}
func lr0_automata(filepath string) ([]State_with_next, []util.Bnf) {

	bnf, err := ioutil.ReadFile(filepath)
	util.Check(err)

	bnf_parsed := util.Parse_bnf_file(string(bnf))

	start_state := gen_start_state(bnf_parsed)

	automaton_states := []State_with_next{}
	_, root_index := add_to_automaton_states(&automaton_states, start_state)
	not_explored_queue := not_explored_queue{queue: []int{root_index}}

	for !not_explored_queue.empty() {
		root_index := not_explored_queue.dequeue()
		root_state := automaton_states[root_index]
		new_states := create_new_states(bnf_parsed, root_state.state)
		not_explored := add_all_to_automaton_states(&automaton_states, root_index, new_states)

		for _, index := range not_explored {
			not_explored_queue.enqueue(index)
		}
	}

	return automaton_states, bnf_parsed
}

func lr1_automata(filepath string) ([]State_with_follow_next, []util.Bnf) {

	bnf, err := ioutil.ReadFile(filepath)
	util.Check(err)

	bnf_parsed := util.Parse_bnf_file(string(bnf))

	start_state := gen_start_state_with_follow(bnf_parsed)

	automaton_states := []State_with_follow_next{}
	_, root_index := add_to_automaton_states_with_follow(&automaton_states, start_state)
	not_explored_queue := not_explored_queue{queue: []int{root_index}}

	for !not_explored_queue.empty() {
		root_index := not_explored_queue.dequeue()
		root_state := automaton_states[root_index]
		new_states := create_new_states_with_follow(bnf_parsed, root_state.state)
		not_explored := add_all_to_automaton_states_with_follow(&automaton_states, root_index, new_states)

		for _, index := range not_explored {
			not_explored_queue.enqueue(index)
		}
	}

	return automaton_states, bnf_parsed
}

func print_automata(automaton_states []State_with_follow_next, bnf_list []util.Bnf) {

	for _, state := range automaton_states {

		fmt.Printf("--------\n")
		for element, _ := range state.state {
			fmt.Printf("%v->%v:%v:%s\n", bnf_list[element.state_element.Product_id].Left, bnf_list[element.state_element.Product_id].Right[element.state_element.Alternate_id], element.state_element.Offset, element.follow)
		}

	}
}

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

func handle_reduction() {

}
func parse_lr0(automaton_states []State_with_next, input_tokens []string, bnf_list []util.Bnf) symbol_stack {

	state_stack := state_stack{}
	symbol_stack := symbol_stack{}

	state_stack.push(0)

	for _, token := range input_tokens {
		symbol_stack.push(node{node_type: token})
		next_state_id := automaton_states[state_stack.top()].next[symbol_stack.top().node_type]
		state_stack.push(next_state_id)

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

		}

	}

	return symbol_stack

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

func main() {

	_, filename, _, _ := runtime.Caller(0)
	bnf_path := filepath.Join(filepath.Dir(filepath.Dir(filename)), "sample2.bnf")

	automaton_states, bnf_list := lr0_automata(bnf_path)

	input_tokens := []string{"int", "+", "E"}

	symbol_stack := parse_lr0(automaton_states, input_tokens, bnf_list)

	//for _, state := range automaton_states {
	//fmt.Printf("%v\n", state)
	/*}*/

	fmt.Printf("%+v", symbol_stack.data[0])

	//print_automata(automaton_states, bnf_list)

	//removed := ll.Remove_direct_left_recursion(bnf_parsed)
	//nonterminal_set := util.Get_nonterminal(removed)
	//terminal_set := util.Get_terminal(removed, nonterminal_set)

	//for _, nonterminal := range nonterminal_set {
	//first := ll.Get_first_set(terminal_set, nonterminal, removed)
	//fmt.Printf("%v:%v\n", nonterminal, first)

	/*}*/
}