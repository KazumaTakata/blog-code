package main

import "fmt"

func isNumber(ch byte) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}

	return false
}

func isOperator(ch byte) bool {
	if ch == '+' || ch == '-' || ch == '/' || ch == '*' || ch == '^' {
		return true
	}
	return false
}

type Stack struct {
	stack []byte
}

func (s *Stack) pop() byte {
	last := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return last
}

func (s *Stack) push(ch byte) {
	s.stack = append(s.stack, ch)
}

func (s *Stack) top() byte {
	return s.stack[len(s.stack)-1]
}

func (s *Stack) empty() bool {
	if len(s.stack) == 0 {
		return true
	}

	return false
}

func isLeft(ch byte) bool {
	if ch != '^' {
		return true
	}

	return false
}

func main() {

	precedence := map[byte]int{}
	precedence['+'] = 2
	precedence['-'] = 2
	precedence['*'] = 3
	precedence['/'] = 3
	precedence['^'] = 4

	stack := Stack{}
	output := []byte{}

	input := "3+4*2/(1-5)^2^3"

	for len(input) > 0 {
		fmt.Printf("%s\n", input)
		if isNumber(input[0]) {
			output = append(output, input[0])
			input = input[1:]
		} else if isOperator(input[0]) {
			for !stack.empty() && (precedence[input[0]] < precedence[stack.top()] || (precedence[input[0]] == precedence[stack.top()] && isLeft(input[0]))) {
				output = append(output, stack.pop())
			}

			stack.push(input[0])
			input = input[1:]
		} else if input[0] == '(' {
			stack.push(input[0])
			input = input[1:]
		} else if input[0] == ')' {
			token := stack.pop()
			for token != '(' {
				output = append(output, token)
				token = stack.pop()
			}
			input = input[1:]
		}
	}

	for !stack.empty() {
		output = append(output, stack.pop())
	}

	fmt.Printf("%s", output)
}
