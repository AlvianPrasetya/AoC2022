package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := parseInput("in.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) [][]int {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res [][]int
	for s.Scan() {
		str := s.Text()

		var cur []int
		for i := 0; i < len(str); i++ {
			cur = append(cur, int(str[i]-'0'))
		}

		res = append(res, cur)
	}

	return res
}

func solveFirst(input [][]int) int {
	isVisible := make([][]bool, len(input))
	for i := 0; i < len(input); i++ {
		isVisible[i] = make([]bool, len(input[0]))
	}

	// Downwards sweep
	var cur []int
	for i := 0; i < len(input[0]); i++ {
		isVisible[0][i] = true
		cur = append(cur, input[0][i])
	}
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[0]); j++ {
			if input[i][j] > cur[j] {
				isVisible[i][j] = true
				cur[j] = input[i][j]
			}
		}
	}

	// Upwards sweep
	cur = nil
	for i := 0; i < len(input[0]); i++ {
		isVisible[len(input)-1][i] = true
		cur = append(cur, input[len(input)-1][i])
	}
	for i := len(input) - 1; i >= 0; i-- {
		for j := 0; j < len(input[0]); j++ {
			if input[i][j] > cur[j] {
				isVisible[i][j] = true
				cur[j] = input[i][j]
			}
		}
	}

	// Rightwards sweep
	cur = nil
	for i := 0; i < len(input); i++ {
		isVisible[i][0] = true
		cur = append(cur, input[i][0])
	}
	for j := 0; j < len(input[0]); j++ {
		for i := 0; i < len(input); i++ {
			if input[i][j] > cur[i] {
				isVisible[i][j] = true
				cur[i] = input[i][j]
			}
		}
	}

	// Leftwards sweep
	cur = nil
	for i := 0; i < len(input); i++ {
		isVisible[i][len(input[0])-1] = true
		cur = append(cur, input[i][len(input[0])-1])
	}
	for j := len(input[0]) - 1; j >= 0; j-- {
		for i := 0; i < len(input); i++ {
			if input[i][j] > cur[i] {
				isVisible[i][j] = true
				cur[i] = input[i][j]
			}
		}
	}

	var res int
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[0]); j++ {
			if isVisible[i][j] {
				res++
			}
		}
	}

	return res
}

func solveSecond(input [][]int) int {
	var res int
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[0]); j++ {
			// Downwards
			var down int
			for k := i + 1; k < len(input); k++ {
				down++
				if input[k][j] >= input[i][j] {
					break
				}
			}
			// Upwards
			var up int
			for k := i - 1; k >= 0; k-- {
				up++
				if input[k][j] >= input[i][j] {
					break
				}
			}
			// Rightwards
			var right int
			for k := j + 1; k < len(input[0]); k++ {
				right++
				if input[i][k] >= input[i][j] {
					break
				}
			}
			// Leftwards
			var left int
			for k := j - 1; k >= 0; k-- {
				left++
				if input[i][k] >= input[i][j] {
					break
				}
			}

			res = max(res, down*up*right*left)
		}
	}

	return res
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
