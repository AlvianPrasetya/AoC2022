package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	surfaceDirs = [6]Pos{
		{X: -1},
		{X: 1},
		{Y: -1},
		{Y: 1},
		{Z: -1},
		{Z: 1},
	}
)

type Pos struct {
	X int
	Y int
	Z int
}

func (this Pos) Add(other Pos) Pos {
	return Pos{
		X: this.X + other.X,
		Y: this.Y + other.Y,
		Z: this.Z + other.Z,
	}
}

func main() {
	input := parseInput("in.txt")
	fmt.Println(solveFirst(input))
}

func parseInput(in string) map[Pos]bool {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	res := make(map[Pos]bool)
	for s.Scan() {
		tokens := strings.Split(s.Text(), ",")
		x, _ := strconv.ParseInt(tokens[0], 10, 32)
		y, _ := strconv.ParseInt(tokens[1], 10, 32)
		z, _ := strconv.ParseInt(tokens[2], 10, 32)

		res[Pos{int(x), int(y), int(z)}] = true
	}

	return res
}

func solveFirst(lavas map[Pos]bool) int {
	var res int
	for lava := range lavas {
		for _, dir := range surfaceDirs {
			if !lavas[lava.Add(dir)] {
				res++
			}
		}
	}

	return res
}
