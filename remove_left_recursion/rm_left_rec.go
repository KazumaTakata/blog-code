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
	//	fmt.Printf("%v\n", start_state)

	nonterminal_and_terminal := util.Get_nonterminal_and_terminal(bnf_parsed)

	new_state_elements := map[string][]State_element{}

	for node, _ := range nonterminal_and_terminal {
		for state_ele, _ := range start_state {
			if node == bnf_parsed[state_ele.Product_id].Right[state_ele.Alternate_id][state_ele.Offset] {
				new_state_elements[node] = append(new_state_elements[node], State_element{Product_id: state_ele.Product_id, Alternate_id: state_ele.Alternate_id, Offset: state_ele.Offset + 1})
			}
		}

	}
	//fmt.Printf("%v\n", new_state_elements)

	automaton_states := []State{}
	new_states := map[string]State{}

	for key, elements := range new_state_elements {
		new_state := State{}
		for _, element := range elements {
			new_elements := Expand_non_terminal(element, bnf_parsed)
			for new_ele, _ := range new_elements {
				new_state[new_ele] = true
			}
			new_state[element] = true
		}
		new_states[key] = new_state
	}

	fmt.Printf("%v", new_states)

	//removed := ll.Remove_direct_left_recursion(bnf_parsed)
	//nonterminal_set := util.Get_nonterminal(removed)
	//terminal_set := util.Get_terminal(removed, nonterminal_set)

	//for _, nonterminal := range nonterminal_set {
	//first := ll.Get_first_set(terminal_set, nonterminal, removed)
	//fmt.Printf("%v:%v\n", nonterminal, first)

	/*}*/
}
