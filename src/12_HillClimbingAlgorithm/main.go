package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

const (
	StartHeight = 'a'
	EndHeight   = 'z'
)

var (
	moves = []Pos{
		{X: 1},
		{Y: 1},
		{X: -1},
		{Y: -1},
	}
)

type Pos struct {
	Y int
	X int
}

type State struct {
	Pos  Pos
	Dist int
}

func main() {
	heightMap, start, end := parseInput("in.txt")
	fmt.Println(solveFirst(heightMap, start, end))
	fmt.Println(solveSecond(heightMap, start, end))
}

func parseInput(in string) ([][]int, Pos, Pos) {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var heightMap [][]int
	var start Pos
	var end Pos
	for s.Scan() {
		str := s.Text()
		heightMap = append(heightMap, nil)
		for i := 0; i < len(str); i++ {
			switch str[i] {
			case 'S':
				heightMap[len(heightMap)-1] = append(heightMap[len(heightMap)-1], StartHeight-'a')
				start = Pos{
					Y: len(heightMap) - 1,
					X: i,
				}
			case 'E':
				heightMap[len(heightMap)-1] = append(heightMap[len(heightMap)-1], EndHeight-'a')
				end = Pos{
					Y: len(heightMap) - 1,
					X: i,
				}
			default:
				heightMap[len(heightMap)-1] = append(heightMap[len(heightMap)-1], int(str[i]-'a'))
			}
		}
	}

	return heightMap, start, end
}

func solveFirst(heightMap [][]int, start Pos, end Pos) int {
	return shortestPath(heightMap, start, func(pos Pos) bool {
		return pos == end
	}, func(cur, next Pos) bool {
		return heightMap[next.Y][next.X]-heightMap[cur.Y][cur.X] <= 1
	})
}

func solveSecond(heightMap [][]int, start Pos, end Pos) int {
	return shortestPath(heightMap, end, func(pos Pos) bool {
		return heightMap[pos.Y][pos.X] == 0
	}, func(cur, next Pos) bool {
		return heightMap[cur.Y][cur.X]-heightMap[next.Y][next.X] <= 1
	})
}

func shortestPath(heightMap [][]int, start Pos, hasReached func(pos Pos) bool, canMove func(cur, next Pos) bool) int {
	queue := list.New()
	queue.PushBack(State{Pos: start, Dist: 0})
	visited := make(map[Pos]bool)
	for queue.Len() != 0 {
		cur := queue.Remove(queue.Front()).(State)
		if visited[cur.Pos] {
			// Already visited
			continue
		}
		visited[cur.Pos] = true

		if hasReached(cur.Pos) {
			// Reached end
			return cur.Dist
		}

		for _, move := range moves {
			newPos := Pos{
				Y: cur.Pos.Y + move.Y,
				X: cur.Pos.X + move.X,
			}
			if newPos.Y < 0 || newPos.Y == len(heightMap) || newPos.X < 0 || newPos.X == len(heightMap[0]) {
				// Out-of-bounds
				continue
			}
			if !canMove(cur.Pos, newPos) {
				// Too high
				continue
			}

			queue.PushBack(State{Pos: newPos, Dist: cur.Dist + 1})
		}
	}

	return 0
}
