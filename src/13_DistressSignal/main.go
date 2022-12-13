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

func (this List) Compare(other List) int {
	for i, j := 0, 0; i < len(this) && j < len(other); i, j = i+1, j+1 {
		cmp := compare(this[i], other[j])
		if cmp != 0 {
			return cmp
		}
	}

	if len(this) < len(other) {
		return -1
	}
	if len(this) > len(other) {
		return 1
	}

	return 0
}

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
	// Both are Lists
	return first.(List).Compare(second.(List))
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
	/*str = strings.TrimSuffix(strings.TrimPrefix(str, "["), "]")
	var bal int
	var buf string
	var tokens []string
	for i := 0; i < len(str); i++ {
		buf += string(str[i])
		if str[i] == '[' {
			bal++
		} else if str[i] == ']' {
			bal--
		} else if str[i] == ',' {
			if bal == 0 && buf != "" {
				// Flush buf
				tokens = append(tokens, buf[:len(buf)-1])
				buf = ""
			}
		}
	}
	if buf != "" {
		// Flush buf
		tokens = append(tokens, buf)
		buf = ""
	}

	var res List
	for _, token := range tokens {
		if token[0] == '[' {
			// List
			res = append(res, parseList(token))
		} else {
			// Int
			val, _ := strconv.ParseInt(token, 10, 32)
			res = append(res, int(val))
		}
	}

	return res*/

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
		cmp := pair.First.Compare(pair.Second)
		if cmp < 0 {
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
		return lists[i].Compare(lists[j]) < 0
	})

	res := 1
	for i, list := range lists {
		if list.Compare(dividerA) == 0 || list.Compare(dividerB) == 0 {
			res *= i + 1
		}
	}

	return res
}
