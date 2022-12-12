package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type State struct {
	Stacks [][]byte
}

type Move struct {
	N    int
	From int // 0-indexed
	To   int // 0-indexed
}

func main() {
	in := "in.txt"

	state, moves := parseInput(in)
	fmt.Println(solveFirst(state, moves))

	state, moves = parseInput(in)
	fmt.Println(solveSecond(state, moves))
}

func parseInput(in string) (*State, []Move) {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	// Read state
	state := &State{}
	for s.Scan() {
		str := s.Text()

		if str == "" {
			break
		}

		for i := 0; i*4+1 < len(str); i++ {
			if i == len(state.Stacks) {
				// Add a new stack
				state.Stacks = append(state.Stacks, nil)
			}

			if str[i*4+1] >= 'A' && str[i*4+1] <= 'Z' {
				state.Stacks[i] = append([]byte{str[i*4+1]}, state.Stacks[i]...)
			}
		}
	}

	// Read moves
	var moves []Move
	for s.Scan() {
		str := s.Text()
		str = strings.Replace(str, "move ", "", 1)
		str = strings.Replace(str, "from ", "", 1)
		str = strings.Replace(str, "to ", "", 1)

		tokens := strings.Split(str, " ")

		n, _ := strconv.ParseInt(tokens[0], 10, 32)
		from, _ := strconv.ParseInt(tokens[1], 10, 32)
		to, _ := strconv.ParseInt(tokens[2], 10, 32)

		moves = append(moves, Move{
			N:    int(n),
			From: int(from - 1),
			To:   int(to - 1),
		})
	}

	return state, moves
}

func solveFirst(state *State, moves []Move) string {
	for _, move := range moves {
		toMove := state.Stacks[move.From][len(state.Stacks[move.From])-move.N:]
		state.Stacks[move.From] = state.Stacks[move.From][:len(state.Stacks[move.From])-move.N]
		state.Stacks[move.To] = append(state.Stacks[move.To], reverse(toMove)...)
	}

	var res string
	for _, stack := range state.Stacks {
		if len(stack) != 0 {
			res += string(stack[len(stack)-1])
		}
	}

	return res
}

func solveSecond(state *State, moves []Move) string {
	for _, move := range moves {
		toMove := state.Stacks[move.From][len(state.Stacks[move.From])-move.N:]
		state.Stacks[move.From] = state.Stacks[move.From][:len(state.Stacks[move.From])-move.N]
		state.Stacks[move.To] = append(state.Stacks[move.To], toMove...)
	}

	var res string
	for _, stack := range state.Stacks {
		if len(stack) != 0 {
			res += string(stack[len(stack)-1])
		}
	}

	return res
}

func reverse(arr []byte) []byte {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-i-1] = arr[len(arr)-i-1], arr[i]
	}

	return arr
}
