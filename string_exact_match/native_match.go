package main

import (
	"bytes"
	"fmt"
)

type matched_pos struct {
	start  int
	length int
}

func naive_string_match(pattern, text []byte) []matched_pos {
	searched_id := 0
	pattern_id := 0

	matched := []matched_pos{}

	for searched_id <= len(text)-len(pattern) {
		for pattern_id < len(pattern) {
			if text[searched_id+pattern_id] == pattern[pattern_id] {
				pattern_id += 1
				if pattern_id == len(pattern) {
					pattern_id = 0
					matched = append(matched, matched_pos{start: searched_id, length: len(pattern)})
					break
				}
			} else {
				pattern_id = 0
				break
			}
		}
		searched_id += 1
	}

	return matched

}

func get_l_i(suffix, pattern []byte) int {

	stlide := 1
	full_length := len(pattern)

	preceding := len(pattern) - len(suffix)

	if preceding >= 0 {
		p_i_1 := pattern[preceding]

		for len(suffix) <= len(pattern) {
			slice_index := len(pattern) - len(suffix)
			if bytes.Equal(suffix, pattern[slice_index:]) && (slice_index-1 < 0 || p_i_1 != pattern[slice_index-1]) {
				//fmt.Printf("%v\n", suffix)
				//fmt.Printf("%v\n", pattern[slice_index:])
				return full_length - stlide + 1
			} else {
				//fmt.Printf("%d", stlide)
				stlide = stlide + 1
				pattern = pattern[:len(pattern)-1]
			}
		}
	}
	return 0

}

func preprocess_good_suffix(pattern []byte) ([]int, []int) {

	table_l := make([]int, len(pattern))
	table_h := make([]int, len(pattern))

	for i, _ := range pattern {
		suffix := pattern[i:]
		index := get_l_i(suffix, pattern[:len(pattern)-1])
		table_l[i] = index
	}

	for i, _ := range pattern {
		suffix := pattern[i:]

		for len(suffix) > 0 {
			if bytes.Equal(suffix, pattern[:len(suffix)]) {
				table_h[i] = len(suffix)
				break
			} else {
				suffix = suffix[1:]
			}
		}
	}

	return table_l, table_h
}

func preprocess_bad_char(pattern []byte) []map[byte]int {
	bad_ch_table := []map[byte]int{}

	for i, p_ch := range pattern {
		new_map := make(map[byte]int)

		if i != 0 {

			for j, ch := range pattern[:i] {
				if ch != p_ch {
					skip_length := i - j - 1
					new_map[ch] = skip_length
				}
			}

		}
		bad_ch_table = append(bad_ch_table, new_map)
	}

	return bad_ch_table
}

func boyer_moore_string_match(pattern, text []byte, bad_ch_table []map[byte]int, table_l, table_h []int) []matched_pos {

	matched_list := []matched_pos{}

	for i := 0; i < len(text)-len(pattern)+1; i++ {
		fmt.Printf("%d\n", i)
		for j, _ := range pattern {
			p_index := len(pattern) - j - 1
			if text[i+p_index] == pattern[p_index] {
				if p_index == 0 {
					new_pos := matched_pos{start: i, length: len(pattern)}
					matched_list = append(matched_list, new_pos)
				}
			} else {

				slide_bad_ch := 0
				if p_index > 0 {

					bad_ch := text[i+p_index]
					if skip_i, ok := bad_ch_table[p_index][bad_ch]; ok {
						slide_bad_ch = skip_i
					} else {
						slide_bad_ch = p_index
					}
				}

				slide_l_h := 0
				if p_index != len(pattern)-1 {
					if table_l[p_index+1] != 0 {
						slide_l_h = len(pattern) - table_l[p_index+1] - 1
					} else {
						slide_l_h = len(pattern) - table_h[p_index+1] - 1
					}
				}

				slide := 0
				if slide_bad_ch >= slide_l_h {
					slide = slide_bad_ch
				} else {
					slide = slide_l_h
				}

				i += slide

				break
			}
		}
	}

	return matched_list
}

func main() {

	text := []byte("aivaie avava a kvava")
	pattern := []byte("avava")
	//pattern := []byte("hvvavaava")
	//pattern := []byte("vavavaava")

	//matched_list := naive_string_match(pattern, text)

	//fmt.Printf("%v", matched_list)

	bad_ch_table := preprocess_bad_char(pattern)

	l, h := preprocess_good_suffix(pattern)

	fmt.Printf("%v:%v", l, h)

	matched_list := boyer_moore_string_match(pattern, text, bad_ch_table, l, h)

	fmt.Printf("%v", matched_list)

}
