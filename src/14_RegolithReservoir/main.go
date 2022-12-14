package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := "in.txt"

	board, minX, maxX, _ := parseInput(in)
	//fmt.Println(board, minX, maxX)
	fmt.Println(solveFirst(board, minX, maxX))

	board, minX, maxX, maxY := parseInput(in)
	//fmt.Println(board, minX, maxX, maxY)
	fmt.Println(solveSecond(board, minX, maxX, maxY))
}

func parseInput(in string) (map[int]map[int]bool, int, int, int) {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	board := make(map[int]map[int]bool)
	minX, maxX, maxY := math.MaxInt, math.MinInt, math.MinInt
	for s.Scan() {
		str := s.Text()
		tokens := strings.Split(str, " -> ")

		prevX, prevY := -1, -1
		for _, token := range tokens {
			xy := strings.Split(token, ",")
			xInt32, _ := strconv.ParseInt(xy[0], 10, 32)
			yInt32, _ := strconv.ParseInt(xy[1], 10, 32)
			x, y := int(xInt32), int(yInt32)

			if prevX != -1 && prevY != -1 {
				startX, endX := min(x, prevX), max(x, prevX)
				startY, endY := min(y, prevY), max(y, prevY)
				minX, maxX, maxY = min(minX, startX), max(maxX, endX), max(maxY, endY)

				for i := startX; i <= endX; i++ {
					if board[i] == nil {
						board[i] = make(map[int]bool)
					}
					for j := startY; j <= endY; j++ {
						board[i][j] = true
					}
				}
			}
			prevX, prevY = x, y
		}
	}

	return board, minX, maxX, maxY
}

func solveFirst(board map[int]map[int]bool, minX int, maxX int) int {
	var res int
	for {
		x, y := 500, 0
		for {
			minY := math.MaxInt
			for filledY := range board[x] {
				if y > filledY {
					continue
				}
				minY = min(minY, filledY)
			}

			if minY == math.MaxInt {
				// Fell into the abyss
				return res
			}

			y = minY - 1
			if !board[x-1][y+1] {
				// Down-left
				x = x - 1
				y = y + 1
				continue
			}
			if !board[x+1][y+1] {
				// Down-right
				x = x + 1
				y = y + 1
				continue
			}
			//fmt.Println(x, y)
			board[x][y] = true
			res++
			break
		}
	}
}

func solveSecond(board map[int]map[int]bool, minX int, maxX int, maxY int) int {
	floorY := maxY + 2

	var res int
	for {
		x, y := 500, 0
		for {
			minY := math.MaxInt
			for filledY := range board[x] {
				if y > filledY {
					continue
				}
				minY = min(minY, filledY)
			}

			if minY == math.MaxInt {
				minY = floorY
			}

			y = minY - 1
			if !board[x-1][y+1] && y+1 != floorY {
				// Down-left
				x = x - 1
				y = y + 1
				continue
			}
			if !board[x+1][y+1] && y+1 != floorY {
				// Down-right
				x = x + 1
				y = y + 1
				continue
			}

			//fmt.Println(x, y)
			if board[x] == nil {
				board[x] = make(map[int]bool)
			}
			board[x][y] = true
			res++

			if x == 500 && y == 0 {
				return res
			}
			break
		}
	}
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
