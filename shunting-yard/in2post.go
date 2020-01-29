package main

import "fmt"

func isNumber(ch byte) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}

	return false
}

func isAlphabet(ch byte) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
		return true
	}

	return false

}

type stack struct {
	stack []byte
}

func (s *stack) pop() byte {
	last := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return last
}

func (s *stack) push(ch byte) {
	s.stack = append(s.stack, ch)
}

func (s *stack) top() byte {
	return s.stack[len(s.stack)-1]
}

func (s *stack) empty() bool {
	if len(s.stack) == 0 {
		return true
	}

	return false
}

type In2Post struct {
	precedence       map[byte]int
	rightAssociative []byte
}

func NewIn2Post(operators []Operator) In2Post {
	precedence := map[byte]int{}
	rightAssociative := []byte{}
	for _, operator := range operators {
		precedence[operator.value] = operator.precedence
		if !operator.isLeftAssociative {
			rightAssociative = append(rightAssociative, operator.value)
		}
	}

	in2post := In2Post{precedence: precedence, rightAssociative: rightAssociative}
	return in2post

}

type Operator struct {
	value             byte
	precedence        int
	isLeftAssociative bool
}

func (in2post *In2Post) isLeftAssociative(ch byte) bool {
	for _, right_op := range in2post.rightAssociative {
		if right_op == ch {
			return false
		}
	}
	return true
}

func (in2post *In2Post) isOperator(ch byte) bool {
	for op, _ := range in2post.precedence {
		if op == ch {
			return true
		}

	}
	return false
}
func (in2post *In2Post) parse(input string) []byte {

	stack := stack{}
	output := []byte{}

	for len(input) > 0 {
		if isNumber(input[0]) || isAlphabet(input[0]) {
			output = append(output, input[0])
			input = input[1:]
		} else if in2post.isOperator(input[0]) {
			for !stack.empty() && (in2post.precedence[input[0]] < in2post.precedence[stack.top()] || (in2post.precedence[input[0]] == in2post.precedence[stack.top()] && in2post.isLeftAssociative(input[0]))) {
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

	return output
}

func main() {

	operators := []Operator{}
	operators = append(operators, Operator{value: '.', precedence: 1, isLeftAssociative: true})
	operators = append(operators, Operator{value: '+', precedence: 2, isLeftAssociative: true})
	operators = append(operators, Operator{value: '-', precedence: 2, isLeftAssociative: true})
	operators = append(operators, Operator{value: '*', precedence: 3, isLeftAssociative: true})
	operators = append(operators, Operator{value: '/', precedence: 3, isLeftAssociative: true})
	operators = append(operators, Operator{value: '^', precedence: 4, isLeftAssociative: false})

	i2p := NewIn2Post(operators)
	output := i2p.parse("s+.a+")

	fmt.Printf("%s", output)
}
