package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	//	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Bnf struct {
	left  string
	right [][]string
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func get_terminal(bnf_list []Bnf, non_terminal []string) []string {

	terminal := []string{}

	for _, prod := range bnf_list {
		for _, rights := range prod.right {
			for _, right := range rights {
				if !Contains(non_terminal, right) {
					terminal = append(terminal, right)
				}
			}
		}
	}

	return terminal

}

func get_nonterminal(bnf_list []Bnf) []string {

	non_terminal := []string{}

	for _, prod := range bnf_list {
		non_terminal = append(non_terminal, prod.left)
	}

	return non_terminal

}

func check_direct_left_recursion(bnf_list []Bnf) []Bnf {

	direct_left := []Bnf{}

	for _, bnf := range bnf_list {
		for _, right := range bnf.right {
			if bnf.left == right[0] {
				direct_left = append(direct_left, bnf)
			}
		}
	}

	return direct_left
}

func check_direct_left_recursion_of_one_production(bnf Bnf) bool {
	for _, right := range bnf.right {
		if bnf.left == right[0] {
			return true
		}
	}
	return false
}

//https://en.wikipedia.org/wiki/Left_recursion

func remove_direct_left_recursion(bnf_list []Bnf) []Bnf {

	removed_left := []Bnf{}

	for _, bnf := range bnf_list {
		if check_direct_left_recursion_of_one_production(bnf) {
			left_rec := Bnf{left: bnf.left + "'", right: [][]string{}}
			non_left_rec := Bnf{left: bnf.left, right: [][]string{}}

			for _, right := range bnf.right {
				if bnf.left == right[0] {
					left_rec.right = append(left_rec.right, append(right[1:], bnf.left+"'"))

				} else {
					non_left_rec.right = append(non_left_rec.right, append(right, bnf.left+"'"))
				}
			}
			left_rec.right = append(left_rec.right, []string{"epsilon"})

			removed_left = append(removed_left, non_left_rec)
			removed_left = append(removed_left, left_rec)
		} else {
			removed_left = append(removed_left, bnf)
		}
	}
	return removed_left
}

func get_first_set(terminal_set []string, non_terminal string, bnf_list []Bnf) []string {

	first := []string{}

	for _, bnf := range bnf_list {
		if bnf.left == non_terminal {
			for _, right := range bnf.right {
				if Contains(terminal_set, right[0]) {
					first = append(first, right[0])
				} else {
					first_ := get_first_set(terminal_set, right[0], bnf_list)
					first = append(first, first_...)
				}
			}
		}
	}
	return first
}

//func parse_bnf_right(right string) [][]string {
//rights := strings.Split(right, "|")
//parsed_list := [][]string{}
//for _, right := range rights {
//parsed_list = append(parsed_list, strings.Split(right, " "))
//}

//trimed_list := [][]string{}
//for _, right := range parsed_list {
//tmp := []string{}
//for _, r := range right {

//trimed := strings.TrimSpace(r)
//if len(trimed) > 0 {
//tmp = append(tmp, trimed)
//}
//}

//trimed_list = append(trimed_list, tmp)
//}

//return trimed_list
//}

//func parse_bnf_file(bnf_string string) []Bnf {

//bnf_lines := strings.Split(bnf_string, "\n")

//bnf_parsed := []Bnf{}

//for _, line := range bnf_lines {
//left_right := strings.Split(line, "::=")

//if len(left_right) > 1 {
//trimed_list := parse_bnf_right(left_right[1])
//bnf_parsed = append(bnf_parsed, Bnf{left: strings.TrimSpace(left_right[0]), right: trimed_list})
//}
//}

//return bnf_parsed

//}

func main() {

	_, filename, _, _ := runtime.Caller(0)
	bnf_path := filepath.Join(filepath.Dir(filename), "sample.bnf")

	bnf, err := ioutil.ReadFile(bnf_path)
	check(err)

	bnf_parsed := parse_bnf_file(string(bnf))

	//non_term := get_nonterminal(bnf_parsed)
	//term := get_terminal(bnf_parsed, non_term)
	//fmt.Printf("non_term is:%v\nterm is :%v", non_term, term)

	//direct_rec_bnf := check_direct_left_recursion(bnf_parsed)
	//fmt.Printf("%v", direct_rec_bnf)

	removed := remove_direct_left_recursion(bnf_parsed)
	nonterminal_set := get_nonterminal(removed)
	terminal_set := get_terminal(removed, nonterminal_set)

	//for _, r := range removed {
	//fmt.Printf("%v\n", r)

	//}

	for _, nonterminal := range nonterminal_set {
		first := get_first_set(terminal_set, nonterminal, removed)
		fmt.Printf("%v:%v\n", nonterminal, first)

	}
}
