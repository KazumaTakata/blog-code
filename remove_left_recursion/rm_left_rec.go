package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	//	"remove_left_recursion/ll"
	"remove_left_recursion/util"
	"runtime"
	//	"strings"
)

type State_element struct {
	Product_id   int
	Alternate_id int
	Offset       int
}

type State map[State_element]bool

type State_with_next struct {
	next  map[string]int
	state State
}

func Expand_non_terminal(state_element State_element, bnf_list []util.Bnf) State {

	non_terminals := util.Get_nonterminal(bnf_list)
	added_state := make(State)

	if state_element.Offset < len(bnf_list[state_element.Product_id].Right[state_element.Alternate_id]) {
		right_ele := bnf_list[state_element.Product_id].Right[state_element.Alternate_id][state_element.Offset]
		if _, ok := non_terminals[right_ele]; ok {
			for prod_id, bnf := range bnf_list {
				if bnf.Left == right_ele {
					for alte_id, _ := range bnf.Right {
						new_state_element := State_element{Product_id: prod_id, Alternate_id: alte_id, Offset: 0}
						added_state[new_state_element] = true
						new_state := Expand_non_terminal(new_state_element, bnf_list)
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
			new_elements := Expand_non_terminal(element, bnf_list)
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

type not_explored_queue struct {
	queue []int
}

func (q *not_explored_queue) enqueu(new_item int) {
	q.queue = append(q.queue, new_item)
}

func (q *not_explored_queue) dequeu() int {
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

func main() {

	_, filename, _, _ := runtime.Caller(0)
	bnf_path := filepath.Join(filepath.Dir(filename), "sample2.bnf")

	bnf, err := ioutil.ReadFile(bnf_path)
	util.Check(err)

	bnf_parsed := util.Parse_bnf_file(string(bnf))
	start_state := make(State)
	start_state[State_element{Product_id: 0, Alternate_id: 0, Offset: 0}] = true

	new_start_state := make(State)
	for state_ele, _ := range start_state {
		new_elements := Expand_non_terminal(state_ele, bnf_parsed)
		for new_ele, _ := range new_elements {
			new_start_state[new_ele] = true
		}
	}
	for new_elem, _ := range new_start_state {
		start_state[new_elem] = true
	}

	automaton_states := []State_with_next{}
	_, root_index := add_to_automaton_states(&automaton_states, start_state)
	not_explored_queue := not_explored_queue{queue: []int{root_index}}

	for !not_explored_queue.empty() {
		root_index := not_explored_queue.dequeu()
		root_state := automaton_states[root_index]
		new_states := create_new_states(bnf_parsed, root_state.state)
		not_explored := add_all_to_automaton_states(&automaton_states, root_index, new_states)

		for _, index := range not_explored {
			not_explored_queue.enqueu(index)
		}
	}
	for _, state := range automaton_states {
		fmt.Printf("%v\n", state)
	}

	//removed := ll.Remove_direct_left_recursion(bnf_parsed)
	//nonterminal_set := util.Get_nonterminal(removed)
	//terminal_set := util.Get_terminal(removed, nonterminal_set)

	//for _, nonterminal := range nonterminal_set {
	//first := ll.Get_first_set(terminal_set, nonterminal, removed)
	//fmt.Printf("%v:%v\n", nonterminal, first)

	/*}*/
}
