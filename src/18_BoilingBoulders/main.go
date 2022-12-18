package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	dirs = [6]Pos{
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

func (this Pos) IsInRange(rangeX, rangeY, rangeZ Range) bool {
	if this.X < rangeX.Min || this.X > rangeX.Max {
		return false
	}
	if this.Y < rangeY.Min || this.Y > rangeY.Max {
		return false
	}
	if this.Z < rangeZ.Min || this.Z > rangeZ.Max {
		return false
	}
	return true
}

type Range struct {
	Min int
	Max int
}

func NewRange() Range {
	return Range{Min: math.MaxInt, Max: math.MinInt}
}

func (r Range) Update(x int) Range {
	return Range{
		Min: min(r.Min, x),
		Max: max(r.Max, x),
	}
}

func main() {
	input := parseInput("in.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
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
		for _, dir := range dirs {
			if !lavas[lava.Add(dir)] {
				res++
			}
		}
	}

	return res
}

func solveSecond(lavas map[Pos]bool) int {
	rangeX, rangeY, rangeZ := NewRange(), NewRange(), NewRange()
	for lava := range lavas {
		rangeX = rangeX.Update(lava.X)
		rangeY = rangeY.Update(lava.Y)
		rangeZ = rangeZ.Update(lava.Z)
	}
	rangeX = rangeX.Update(rangeX.Min - 1).Update(rangeX.Max + 1)
	rangeY = rangeX.Update(rangeY.Min - 1).Update(rangeY.Max + 1)
	rangeZ = rangeX.Update(rangeZ.Min - 1).Update(rangeZ.Max + 1)

	replacements := []Pos{
		{rangeX.Min, math.MaxInt, math.MaxInt},
		{rangeX.Max, math.MaxInt, math.MaxInt},
		{math.MaxInt, rangeY.Min, math.MaxInt},
		{math.MaxInt, rangeY.Max, math.MaxInt},
		{math.MaxInt, math.MaxInt, rangeZ.Min},
		{math.MaxInt, math.MaxInt, rangeZ.Max},
	}

	outer := make(map[Pos]bool)
	for _, replacement := range replacements {
		surfaceX := rangeX
		if replacement.X != math.MaxInt {
			surfaceX = Range{replacement.X, replacement.X}
		}
		surfaceY := rangeY
		if replacement.Y != math.MaxInt {
			surfaceY = Range{replacement.Y, replacement.Y}
		}
		surfaceZ := rangeZ
		if replacement.Z != math.MaxInt {
			surfaceZ = Range{replacement.Z, replacement.Z}
		}

		for x := surfaceX.Min; x <= surfaceX.Max; x++ {
			for y := surfaceY.Min; y <= surfaceY.Max; y++ {
				for z := surfaceZ.Min; z <= surfaceZ.Max; z++ {
					queue := list.New()
					queue.PushBack(Pos{x, y, z})
					for queue.Len() != 0 {
						cur := queue.Remove(queue.Front()).(Pos)
						if !cur.IsInRange(rangeX, rangeY, rangeZ) || outer[cur] || lavas[cur] {
							continue
						}
						outer[cur] = true

						for _, dir := range dirs {
							queue.PushBack(cur.Add(dir))
						}
					}
				}
			}
		}
	}

	var res int
	for lava := range lavas {
		for _, dir := range dirs {
			if outer[lava.Add(dir)] {
				res++
			}
		}
	}

	return res
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
