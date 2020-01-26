package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"remove_left_recursion/ll"
	"remove_left_recursion/util"
	"runtime"
	//	"strings"
)

func main() {

	_, filename, _, _ := runtime.Caller(0)
	bnf_path := filepath.Join(filepath.Dir(filename), "sample.bnf")

	bnf, err := ioutil.ReadFile(bnf_path)
	util.Check(err)

	bnf_parsed := util.Parse_bnf_file(string(bnf))

	removed := ll.Remove_direct_left_recursion(bnf_parsed)
	nonterminal_set := util.Get_nonterminal(removed)
	terminal_set := util.Get_terminal(removed, nonterminal_set)

	for _, nonterminal := range nonterminal_set {
		first := ll.Get_first_set(terminal_set, nonterminal, removed)
		fmt.Printf("%v:%v\n", nonterminal, first)

	}
}
