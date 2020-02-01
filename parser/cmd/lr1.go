package main

import (
	"fmt"
	"io/ioutil"
	//	"remove_left_recursion/ll"
	"parser/util"
	//	"strings"
)

type State_element_with_follow struct {
	state_element State_element
	follow        string
}
type State_with_follow map[State_element_with_follow]bool

type State_with_follow_next struct {
	next  map[string]int
	state State_with_follow
}

func print_automata(automaton_states []State_with_follow_next, bnf_list []util.Bnf) {

	for _, state := range automaton_states {

		fmt.Printf("--------\n")
		for element, _ := range state.state {
			fmt.Printf("%v->%v:%v:%s\n", bnf_list[element.state_element.Product_id].Left, bnf_list[element.state_element.Product_id].Right[element.state_element.Alternate_id], element.state_element.Offset, element.follow)
		}

	}
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
