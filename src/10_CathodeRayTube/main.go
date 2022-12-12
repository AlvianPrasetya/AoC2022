package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	InstrTypeNoop = InstrType("noop")
	InstrTypeAddx = InstrType("addx")
)

type InstrType string

type Instr struct {
	Type  InstrType
	Value int
}

func main() {
	input := parseInput("input.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) []Instr {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []Instr
	for s.Scan() {
		str := s.Text()

		tokens := strings.Split(str, " ")
		switch tokens[0] {
		case string(InstrTypeNoop):
			res = append(res, Instr{
				Type: InstrTypeNoop,
			})
		case string(InstrTypeAddx):
			value, _ := strconv.ParseInt(tokens[1], 10, 32)
			res = append(res, Instr{
				Type:  InstrTypeAddx,
				Value: int(value),
			})
		}
	}

	return res
}

func solveFirst(input []Instr) int {
	cycles := []int{20, 60, 100, 140, 180, 220}

	strs := eval(input)
	//fmt.Println(strs[1:])

	var res int
	for _, cycle := range cycles {
		//fmt.Printf("%d %d\n", cycle, strs[cycle])
		res += cycle * strs[cycle]
	}

	return res
}

func solveSecond(input []Instr) string {
	strs := eval(input)

	mat := make([][]byte, 6)
	for i := 0; i < 6; i++ {
		mat[i] = make([]byte, 40)
		for j := 0; j < 40; j++ {
			cur := strs[i*40+j+1]
			//fmt.Println(i*40+j+1, strs[i*40+j+1])
			if cur-1 <= j && j <= cur+1 {
				mat[i][j] = '#'
			} else {
				mat[i][j] = '.'
			}
		}
	}

	var res string
	for i := 0; i < 6; i++ {
		for j := 0; j < 40; j++ {
			res += string(mat[i][j])
		}
		res += "\n"
	}

	return res
}

func eval(input []Instr) []int {
	res := []int{1}

	var tmp int
	for _, in := range input {
		res = append(res, res[len(res)-1]+tmp)
		tmp = 0

		switch in.Type {
		case InstrTypeNoop:
		case InstrTypeAddx:
			res = append(res, res[len(res)-1])
			tmp = in.Value
		}
	}

	return res
}
