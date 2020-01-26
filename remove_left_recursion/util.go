package main

import "strings"

func parse_bnf_right(right string) [][]string {
	rights := strings.Split(right, "|")
	parsed_list := [][]string{}
	for _, right := range rights {
		parsed_list = append(parsed_list, strings.Split(right, " "))
	}

	trimed_list := [][]string{}
	for _, right := range parsed_list {
		tmp := []string{}
		for _, r := range right {

			trimed := strings.TrimSpace(r)
			if len(trimed) > 0 {
				tmp = append(tmp, trimed)
			}
		}

		trimed_list = append(trimed_list, tmp)
	}

	return trimed_list
}

func parse_bnf_file(bnf_string string) []Bnf {

	bnf_lines := strings.Split(bnf_string, "\n")

	bnf_parsed := []Bnf{}

	for _, line := range bnf_lines {
		left_right := strings.Split(line, "::=")

		if len(left_right) > 1 {
			trimed_list := parse_bnf_right(left_right[1])
			bnf_parsed = append(bnf_parsed, Bnf{left: strings.TrimSpace(left_right[0]), right: trimed_list})
		}
	}

	return bnf_parsed

}
