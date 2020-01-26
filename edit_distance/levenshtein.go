package main

import "fmt"

func main() {

	Levenshtein("abce", "a")

}

func Levenshtein(str1, str2 string) {

	matrix := make([][]int, len(str1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(str2)+1)
	}

	for i, row := range matrix {
		for j, _ := range row {
			matrix[i][j] = 0
		}
	}

	for i := 1; i <= len(str1); i++ {
		matrix[i][0] = i
	}

	for i := 1; i <= len(str2); i++ {
		matrix[0][i] = i
	}

	for i := 1; i <= len(str2); i++ {
		for j := 1; j <= len(str1); j++ {
			substitute := 0
			if str1[j-1] == str2[i-1] {
				substitute = 0
			} else {
				substitute = 1
			}

			del_cost := matrix[j-1][i] + 1
			insert_cost := matrix[j][i-1] + 1
			sub_cost := matrix[j-1][i-1] + substitute
			matrix[j][i] = min_int(del_cost, insert_cost, sub_cost)
		}
	}

	fmt.Printf("%v", matrix)

}

func min_int(int1, int2, int3 int) int {
	if int1 > int2 {
		if int2 > int3 {
			return int3
		} else {
			return int2
		}
	} else {
		if int1 > int3 {
			return int3
		} else {
			return int1
		}
	}
}
