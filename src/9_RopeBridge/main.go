package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	DirDown  = Pos{Y: 1}
	DirUp    = Pos{Y: -1}
	DirRight = Pos{X: 1}
	DirLeft  = Pos{X: -1}
)

type Pos struct {
	Y int
	X int
}

func (pos Pos) Add(other Pos) Pos {
	return Pos{
		Y: pos.Y + other.Y,
		X: pos.X + other.X,
	}
}

func (pos Pos) Follow(other Pos) Pos {
	dY := other.Y - pos.Y
	dX := other.X - pos.X

	if dX == 0 {
		// Follow vertically
		if dY > 1 {
			return pos.Add(DirDown)
		} else if dY < -1 {
			return pos.Add(DirUp)
		}
		return pos
	} else if dY == 0 {
		// Follow horizontally
		if dX > 1 {
			return pos.Add(DirRight)
		} else if dX < -1 {
			return pos.Add(DirLeft)
		}
		return pos
	}

	// Follow diagonally
	if abs(dX)+abs(dY) == 2 {
		return pos
	}
	if dY > 0 {
		pos = pos.Add(DirDown)
	} else {
		pos = pos.Add(DirUp)
	}
	if dX > 0 {
		pos = pos.Add(DirRight)
	} else {
		pos = pos.Add(DirLeft)
	}
	return pos
}

type Move struct {
	Dir  Pos
	Dist int
}

func main() {
	input := parseInput("input.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) []Move {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []Move
	for s.Scan() {
		str := s.Text()
		tokens := strings.Split(str, " ")

		var cur Move
		switch tokens[0] {
		case "D":
			cur.Dir = DirDown
		case "U":
			cur.Dir = DirUp
		case "R":
			cur.Dir = DirRight
		case "L":
			cur.Dir = DirLeft
		}

		dist, _ := strconv.ParseInt(tokens[1], 10, 32)
		cur.Dist = int(dist)

		res = append(res, cur)
	}

	return res
}

func solveFirst(input []Move) int {
	var head, tail Pos
	tailVisited := make(map[Pos]bool)
	tailVisited[tail] = true
	for _, move := range input {
		for i := 0; i < move.Dist; i++ {
			head = head.Add(move.Dir)
			tail = tail.Follow(head)
			tailVisited[tail] = true
		}
	}

	return len(tailVisited)
}

func solveSecond(input []Move) int {
	ropes := make([]Pos, 10)
	tailVisited := make(map[Pos]bool)
	tailVisited[ropes[9]] = true
	for _, move := range input {
		for i := 0; i < move.Dist; i++ {
			ropes[0] = ropes[0].Add(move.Dir)
			for j := 1; j < 10; j++ {
				ropes[j] = ropes[j].Follow(ropes[j-1])
			}
			tailVisited[ropes[9]] = true
		}
	}

	return len(tailVisited)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
