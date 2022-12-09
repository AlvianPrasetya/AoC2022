package main

import (
	"bufio"
	"fmt"
	"os"
)

type Rucksack struct {
	Left  []byte
	Right []byte
}

func main() {
	input := parseInput("in.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) []*Rucksack {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []*Rucksack
	for s.Scan() {
		str := s.Text()
		n := len(str)

		cur := &Rucksack{}
		for i := 0; i < n/2; i++ {
			cur.Left = append(cur.Left, str[i])
		}
		for i := n / 2; i < n; i++ {
			cur.Right = append(cur.Right, str[i])
		}

		res = append(res, cur)
	}

	return res
}

func solveFirst(input []*Rucksack) int {
	var res int
	for _, rucksack := range input {
		leftMap := make(map[byte]bool)
		for _, ch := range rucksack.Left {
			leftMap[ch] = true
		}
		rightMap := make(map[byte]bool)
		for _, ch := range rucksack.Right {
			rightMap[ch] = true
		}

		for ch := range leftMap {
			if rightMap[ch] {
				res += getPriority(ch)
			}
		}
	}

	return res
}

func solveSecond(input []*Rucksack) int {
	var res int
	for i := 0; i < len(input); i += 3 {
		countMap := make(map[byte]int)
		for j := i; j < i+3; j++ {
			combined := append(input[j].Left, input[j].Right...)
			chMap := make(map[byte]bool)
			for _, ch := range combined {
				chMap[ch] = true
			}

			for ch := range chMap {
				countMap[ch]++
			}
		}

		for ch, count := range countMap {
			if count == 3 {
				res += getPriority(ch)
			}
		}
	}

	return res
}

func getPriority(ch byte) int {
	if ch >= 'a' && ch <= 'z' {
		return 1 + int(ch-'a')
	}
	return 27 + int(ch-'A')
}
