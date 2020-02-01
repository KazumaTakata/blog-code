package main

import (
	"fmt"
	"path/filepath"
	//	"remove_left_recursion/ll"
	"parser/util"
	"runtime"
	//	"strings"
)

type Action_type int

const (
	Shift  Action_type = 0
	Accept Action_type = 1
	Reduce Action_type = 2
)

type Reduction struct {
	left  string
	right []string
}

type Action struct {
	action_type Action_type
	reduction   Reduction
}

type Table_element struct {
	goto_table map[string]int
	action     Action
}

type Table struct {
	table_elements []Table_element
}

func Construct_lr0_Table(state_with_next_list []State_with_next, bnf_list []util.Bnf) Table {

	table := Table{}

	for _, state_with_next := range state_with_next_list {
		table_element := Table_element{goto_table: state_with_next.next}

		handlers := get_handlers(state_with_next, bnf_list)
		if len(handlers) > 0 {
			if len(handlers) > 1 {
				fmt.Printf("reduction conflicts")
			}
			handler := handlers[0]
			reduction := Reduction{left: bnf_list[handler.Product_id].Left, right: bnf_list[handler.Product_id].Right[handler.Alternate_id]}
			action := Action{action_type: Reduce, reduction: reduction}

			table_element.action = action
		} else {
			action := Action{action_type: Shift}
			table_element.action = action
		}

		table.table_elements = append(table.table_elements, table_element)
	}

	return table
}

func Print_lr0_table(table Table, bnf_list []util.Bnf) {

	fmt.Printf("------------------------\n")

	nonterminal_and_terminal := util.Get_nonterminal_and_terminal(bnf_list)
	terms := []string{}
	for term, _ := range nonterminal_and_terminal {
		terms = append(terms, term)
	}

	fmt.Printf(" ")
	for _, term := range terms {
		fmt.Printf(" %v  ", term)
	}

	fmt.Printf("\n")

	for i, table_ele := range table.table_elements {

		fmt.Printf("%d", i)

		for _, term := range terms {
			if index, ok := table_ele.goto_table[term]; ok {
				fmt.Printf(" %v ", index)
			} else {
				fmt.Printf("    ")
			}

		}
		fmt.Printf("\n")
	}

}

func main() {

	_, filename, _, _ := runtime.Caller(0)
	bnf_path := filepath.Join(filepath.Dir(filepath.Dir(filename)), "sample2.bnf")

	automaton_states, bnf_list := lr0_automata(bnf_path)

	table := Construct_lr0_Table(automaton_states, bnf_list)

	Print_lr0_table(table, bnf_list)
	//input_tokens := []string{"int", "+", "(", "int", "+", "int", ";", ")", ";"}

	//symbol_stack := parse_lr0(automaton_states, input_tokens, bnf_list)
	//fmt.Printf("%+v", symbol_stack.data[0])

	//print_automata(automaton_states, bnf_list)

	//removed := ll.Remove_direct_left_recursion(bnf_parsed)
	//nonterminal_set := util.Get_nonterminal(removed)
	//terminal_set := util.Get_terminal(removed, nonterminal_set)

	//for _, nonterminal := range nonterminal_set {
	//first := ll.Get_first_set(terminal_set, nonterminal, removed)
	//fmt.Printf("%v:%v\n", nonterminal, first)

	/*}*/
}
