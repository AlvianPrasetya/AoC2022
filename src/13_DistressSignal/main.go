package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Pair struct {
	First  List
	Second List
}

type List []interface{} // Each element can be int or List

func compare(first interface{}, second interface{}) int {
	aInt, aIsInt := first.(int)
	bInt, bIsInt := second.(int)
	if aIsInt && bIsInt {
		if aInt < bInt {
			return -1
		}
		if aInt > bInt {
			return 1
		}
		return 0
	}
	if aIsInt {
		return compare(List{aInt}, second)
	}
	if bIsInt {
		return compare(first, List{bInt})
	}

	aList := first.(List)
	bList := second.(List)
	for i, j := 0, 0; i < len(aList) && j < len(bList); i, j = i+1, j+1 {
		cmp := compare(aList[i], bList[j])
		if cmp != 0 {
			return cmp
		}
	}

	return compare(len(aList), len(bList))
}

func main() {
	input := parseInput("in.txt")
	//fmt.Println(input)
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) []Pair {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []Pair
	for s.Scan() {
		firstStr := s.Text()
		s.Scan()
		secondStr := s.Text()
		s.Scan()

		res = append(res, Pair{
			First:  parseList(firstStr),
			Second: parseList(secondStr),
		})
	}

	return res
}

func parseList(str string) List {
	var stack []List
	var buf string
	for i := 0; i < len(str); i++ {
		if str[i] >= '0' && str[i] <= '9' {
			buf += string(str[i])
		} else if buf != "" {
			val, _ := strconv.ParseInt(buf, 10, 32)
			stack[len(stack)-1] = append(stack[len(stack)-1], int(val))
			buf = ""
		}

		if str[i] == '[' {
			stack = append(stack, nil)
		} else if str[i] == ']' && len(stack) > 1 {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack[len(stack)-1] = append(stack[len(stack)-1], top)
		}
	}

	return stack[0]
}

func solveFirst(input []Pair) int {
	var res int
	for i, pair := range input {
		if compare(pair.First, pair.Second) < 0 {
			res += i + 1
		}
	}

	return res
}

func solveSecond(input []Pair) int {
	dividerA := List{List{2}}
	dividerB := List{List{6}}

	lists := make([]List, 0, 2*len(input))
	for _, pair := range input {
		lists = append(lists, pair.First, pair.Second)
	}
	lists = append(lists, dividerA, dividerB)

	sort.Slice(lists, func(i, j int) bool {
		return compare(lists[i], lists[j]) < 0
	})

	res := 1
	for i, list := range lists {
		if compare(list, dividerA) == 0 || compare(list, dividerB) == 0 {
			res *= i + 1
		}
	}

	return res
}
