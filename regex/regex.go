package main

import (
	"fmt"
)

type StateType int

const (
	Single StateType = 0
	Split  StateType = 1
	Match  StateType = 2
)

type State struct {
	ch         byte
	out1       *State
	out2       *State
	state_type StateType
}

type Frag struct {
	start *State
	out   []**State
}

type Stack struct {
	frags []*Frag
}

func (s *Stack) push(frag *Frag) {
	s.frags = append(s.frags, frag)
}

func (s *Stack) pop() *Frag {
	last := s.frags[len(s.frags)-1]
	s.frags = s.frags[:len(s.frags)-1]
	return last
}

func patch(out_list []**State, start *State) {
	for i := 0; i < len(out_list); i++ {
		*out_list[i] = start
	}
}

type StateList map[*State]bool

func simulateNFA(start *State, str string) StateList {
	clist := StateList{}
	nlist := StateList{}

	clist[start] = true

	for _, ch := range str {

		for current, _ := range clist {
			if current != nil {
				if current.ch == byte(ch) {
					addstate(&nlist, current.out1)
				}
			}
		}
		/*for state, _ := range nlist {*/
		//fmt.Printf("%d:%c:%v\n", i, ch, state)
		/*}*/
		clist = nlist
		nlist = StateList{}

	}

	return clist
}

func addstate(nlist *StateList, state *State) {
	if state == nil {
		return
	}
	if state.state_type == Single || state.state_type == Match {
		(*nlist)[state] = true
	} else if state.state_type == Split {
		addstate(nlist, state.out1)
		addstate(nlist, state.out2)
	}
}
func isMatch(list StateList) bool {

	for state, _ := range list {
		if state.state_type == Match {
			return true
		}
	}
	return false

}

func main() {

	stack := Stack{frags: []*Frag{}}

	postfix := []byte("s+a+.")

	for _, regex_ch := range postfix {
		switch regex_ch {
		case '.':
			{

				e2 := stack.pop()
				e1 := stack.pop()

				patch(e1.out, e2.start)
				frag := Frag{start: e1.start, out: e2.out}
				stack.push(&frag)
				break

			}
		case '|':
			{

				e2 := stack.pop()
				e1 := stack.pop()

				state := State{out1: e1.start, out2: e2.start, state_type: Split}
				frag := Frag{start: &state, out: append(e1.out, e2.out...)}
				stack.push(&frag)
				break

			}
		case '?':
			{

				e1 := stack.pop()

				state := State{out1: e1.start, out2: nil, state_type: Split}
				frag := Frag{start: &state, out: append(e1.out, &state.out2)}
				stack.push(&frag)
				break

			}
		case '*':
			{

				e1 := stack.pop()

				state := State{out1: e1.start, out2: nil, state_type: Split}
				patch(e1.out, &state)
				frag := Frag{start: &state, out: []**State{&state.out2}}
				stack.push(&frag)
				break

			}
		case '+':
			{

				e1 := stack.pop()

				state := State{out1: e1.start, out2: nil, state_type: Split}
				patch(e1.out, &state)
				frag := Frag{start: e1.start, out: []**State{&state.out2}}
				stack.push(&frag)
				break

			}

		default:
			{

				state := State{ch: regex_ch, out1: nil, out2: nil, state_type: Single}
				frag := Frag{start: &state, out: []**State{&state.out1}}
				stack.push(&frag)
				break
			}
		}
	}

	frag := stack.pop()
	matchstate := State{state_type: Match}
	patch(frag.out, &matchstate)

	//	fmt.Printf("%v", frag.start.out1.out2.out1)

	state_list := simulateNFA(frag.start, "ssaa")

	//for state, _ := range state_list {
	//fmt.Printf("%v\n", state)
	//}

	if isMatch(state_list) {
		fmt.Printf("matched!!")
	}

}
