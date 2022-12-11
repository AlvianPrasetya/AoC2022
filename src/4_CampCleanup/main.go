package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AssignmentPair struct {
	First  AssignmentRange
	Second AssignmentRange
}

type AssignmentRange struct {
	Start int
	End   int
}

func main() {
	input := parseInput("in.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) []AssignmentPair {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []AssignmentPair
	for s.Scan() {
		tokens := strings.Split(s.Text(), ",")
		firstTokens := strings.Split(tokens[0], "-")
		secondTokens := strings.Split(tokens[1], "-")

		firstStart, _ := strconv.ParseInt(firstTokens[0], 10, 32)
		firstEnd, _ := strconv.ParseInt(firstTokens[1], 10, 32)
		secondStart, _ := strconv.ParseInt(secondTokens[0], 10, 32)
		secondEnd, _ := strconv.ParseInt(secondTokens[1], 10, 32)

		res = append(res, AssignmentPair{
			First: AssignmentRange{
				Start: int(firstStart),
				End:   int(firstEnd),
			},
			Second: AssignmentRange{
				Start: int(secondStart),
				End:   int(secondEnd),
			},
		})
	}

	return res
}

func solveFirst(input []AssignmentPair) int {
	var res int
	for _, pair := range input {
		if isContained(pair.First, pair.Second) || isContained(pair.Second, pair.First) {
			res++
		}
	}

	return res
}

func solveSecond(input []AssignmentPair) int {
	var res int
	for _, pair := range input {
		if isOverlap(pair.First, pair.Second) {
			res++
		}
	}

	return res
}

func isContained(outer AssignmentRange, inner AssignmentRange) bool {
	return outer.Start <= inner.Start && outer.End >= inner.End
}

func isOverlap(first AssignmentRange, second AssignmentRange) bool {
	if first.Start <= second.End && first.End >= second.Start {
		return true
	}

	if second.Start <= first.End && second.End >= first.Start {
		return true
	}

	return false
}
