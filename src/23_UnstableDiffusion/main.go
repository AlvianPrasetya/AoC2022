package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

var (
	North = Pos{Y: -1}
	South = Pos{Y: 1}
	West  = Pos{X: -1}
	East  = Pos{X: 1}

	Dirs = [4]Pos{North, South, West, East}
	Adjs = [8]Pos{North, North.Add(East), East, East.Add(South), South, South.Add(West), West, West.Add(North)}
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

type Elf struct {
	DirIdx int
}

func main() {
	in := "in.txt"
	elves := parseInput(in)
	fmt.Println(solveFirst(elves, 10))
	elves = parseInput(in)
	fmt.Println(solveSecond(elves))
}

func parseInput(in string) map[Pos]*Elf {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	elves := make(map[Pos]*Elf)
	var y int
	for s.Scan() {
		row := []byte(s.Text())
		for x := 0; x < len(row); x++ {
			if row[x] == '#' {
				pos := Pos{Y: y, X: x}
				elves[pos] = &Elf{}
			}
		}
		y++
	}

	return elves
}

func solveFirst(elves map[Pos]*Elf, rounds int) int {
	for round := 0; round < rounds; round++ {
		simulateRound(elves)
	}

	minY, maxY := math.MaxInt, math.MinInt
	minX, maxX := math.MaxInt, math.MinInt
	for elfPos := range elves {
		minY, maxY = min(minY, elfPos.Y), max(maxY, elfPos.Y)
		minX, maxX = min(minX, elfPos.X), max(maxX, elfPos.X)
	}

	var res int
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if elves[Pos{Y: y, X: x}] == nil {
				res++
			}
		}
	}

	return res
}

func solveSecond(elves map[Pos]*Elf) int {
	after := make([]Pos, 0, len(elves))
	for pos := range elves {
		after = append(after, pos)
	}
	for round := 1; ; round++ {
		before := after
		simulateRound(elves)
		after = make([]Pos, 0, len(elves))
		for pos := range elves {
			after = append(after, pos)
		}

		if isEqualPosList(before, after) {
			return round
		}
	}
}

func simulateRound(elves map[Pos]*Elf) {
	proposedPosMap := make(map[Pos]Pos)
	proposedPosCount := make(map[Pos]int)
	for pos, elf := range elves {
		isSolitary := true
		for _, adj := range Adjs {
			adjPos := pos.Add(adj)
			if _, ok := elves[adjPos]; ok {
				isSolitary = false
				break
			}
		}
		if isSolitary {
			continue
		}

		for dirOffset := 0; dirOffset < len(Dirs); dirOffset++ {
			dir := Dirs[(elf.DirIdx+dirOffset)%len(Dirs)]
			proposedPos := pos.Add(dir)

			toCheckPos := make([]Pos, 0, 3)
			if dir.Y == 0 {
				toCheckPos = append(toCheckPos, proposedPos, proposedPos.Add(North), proposedPos.Add(South))
			} else {
				toCheckPos = append(toCheckPos, proposedPos, proposedPos.Add(West), proposedPos.Add(East))
			}

			isClear := true
			for _, pos := range toCheckPos {
				if _, ok := elves[pos]; ok {
					isClear = false
					break
				}
			}

			if isClear {
				proposedPosMap[pos] = proposedPos
				proposedPosCount[proposedPos]++
				break
			}
		}
	}

	toRemove := make(map[Pos]*Elf)
	toAdd := make(map[Pos]*Elf)
	for pos, elf := range elves {
		if proposedPos, ok := proposedPosMap[pos]; ok && proposedPosCount[proposedPos] == 1 {
			// Move to proposed pos
			toRemove[pos] = elf
			toAdd[proposedPos] = elf
		}
		// Cycle dir
		elf.DirIdx = (elf.DirIdx + 1) % len(Dirs)
	}

	for pos := range toRemove {
		delete(elves, pos)
	}
	for pos, elf := range toAdd {
		elves[pos] = elf
	}

	// fmt.Println(printBoard(elves))
}

func isEqualPosList(before []Pos, after []Pos) bool {
	if len(before) != len(after) {
		return false
	}

	sort.Slice(before, func(i, j int) bool {
		if before[i].Y == before[j].Y {
			return before[i].X < before[j].X
		}
		return before[i].Y < before[j].Y
	})
	sort.Slice(after, func(i, j int) bool {
		if after[i].Y == after[j].Y {
			return after[i].X < after[j].X
		}
		return after[i].Y < after[j].Y
	})

	// fmt.Println(before)
	// fmt.Println(after)
	// fmt.Println()
	for i := 0; i < len(before); i++ {
		if before[i] != after[i] {
			return false
		}
	}

	return true
}

func printBoard(elves map[Pos]*Elf) string {
	minY, maxY := math.MaxInt, math.MinInt
	minX, maxX := math.MaxInt, math.MinInt
	for elfPos := range elves {
		minY, maxY = min(minY, elfPos.Y), max(maxY, elfPos.Y)
		minX, maxX = min(minX, elfPos.X), max(maxX, elfPos.X)
	}

	var res strings.Builder
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if elves[Pos{Y: y, X: x}] != nil {
				res.WriteByte('#')
			} else {
				res.WriteByte('.')
			}
		}
		res.WriteByte('\n')
	}

	return res.String()
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
