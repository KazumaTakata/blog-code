package ll

import "remove_left_recursion/util"

func Check_direct_left_recursion(bnf_list []util.Bnf) []util.Bnf {

	direct_left := []util.Bnf{}

	for _, bnf := range bnf_list {
		for _, right := range bnf.Right {
			if bnf.Left == right[0] {
				direct_left = append(direct_left, bnf)
			}
		}
	}

	return direct_left
}

func Check_direct_left_recursion_of_one_production(bnf util.Bnf) bool {
	for _, right := range bnf.Right {
		if bnf.Left == right[0] {
			return true
		}
	}
	return false
}

//https://en.wikipedia.org/wiki/Left_recursion

func Remove_direct_left_recursion(bnf_list []util.Bnf) []util.Bnf {

	removed_left := []util.Bnf{}

	for _, bnf := range bnf_list {
		if Check_direct_left_recursion_of_one_production(bnf) {
			left_rec := util.Bnf{Left: bnf.Left + "'", Right: [][]string{}}
			non_left_rec := util.Bnf{Left: bnf.Left, Right: [][]string{}}

			for _, right := range bnf.Right {
				if bnf.Left == right[0] {
					left_rec.Right = append(left_rec.Right, append(right[1:], bnf.Left+"'"))

				} else {
					non_left_rec.Right = append(non_left_rec.Right, append(right, bnf.Left+"'"))
				}
			}
			left_rec.Right = append(left_rec.Right, []string{"epsilon"})

			removed_left = append(removed_left, non_left_rec)
			removed_left = append(removed_left, left_rec)
		} else {
			removed_left = append(removed_left, bnf)
		}
	}
	return removed_left
}

func Get_first_set(terminal_set []string, non_terminal string, bnf_list []util.Bnf) []string {

	first := []string{}

	for _, bnf := range bnf_list {
		if bnf.Left == non_terminal {
			for _, right := range bnf.Right {
				if util.Contains(terminal_set, right[0]) {
					first = append(first, right[0])
				} else {
					first_ := Get_first_set(terminal_set, right[0], bnf_list)
					first = append(first, first_...)
				}
			}
		}
	}
	return first
}
