package main

import (
	"fmt"
	"path/filepath"
	//	"remove_left_recursion/ll"
	"runtime"
	//	"strings"
)

func main() {

	_, filename, _, _ := runtime.Caller(0)
	bnf_path := filepath.Join(filepath.Dir(filepath.Dir(filename)), "sample2.bnf")

	automaton_states, bnf_list := lr0_automata(bnf_path)

	input_tokens := []string{"int", "+", "(", "int", "+", "int", ";", ")", ";"}

	symbol_stack := parse_lr0(automaton_states, input_tokens, bnf_list)
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
