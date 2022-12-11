package main

import (
	"bufio"
	"fmt"
	"os"
)

type Test struct {
}

func main() {
	input := parseInput("input.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) string {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res string
	for s.Scan() {
		res = s.Text()
	}

	return res
}

func solveFirst(input string) int {
	return solve(input, 4)
}

func solveSecond(input string) int {
	return solve(input, 14)
}

func solve(s string, targetCount int) int {
	chCount := make(map[byte]int)
	var uniqCount int
	for i := 0; i < targetCount; i++ {
		chCount[s[i]]++
		if chCount[s[i]] == 1 {
			uniqCount++
		}

		if uniqCount == targetCount {
			return i + 1
		}
	}

	for i := targetCount; i < len(s); i++ {
		chCount[s[i-targetCount]]--
		if chCount[s[i-targetCount]] == 0 {
			uniqCount--
		}

		chCount[s[i]]++
		if chCount[s[i]] == 1 {
			uniqCount++
		}

		if uniqCount == targetCount {
			return i + 1
		}
	}

	return 0
}
