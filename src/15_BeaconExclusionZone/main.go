package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Pos struct {
	X int
	Y int
}

type Sensor struct {
	Pos    Pos
	Beacon Pos
}

type Range struct {
	Start int
	End   int
}

func main() {
	input := parseInput("in.txt")
	//fmt.Println(input)
	fmt.Println(solveFirst(input, 2000000))
	fmt.Println(solveSecond(input, 4000000))
}

func parseInput(in string) []Sensor {
	re := regexp.MustCompile("-?[0-9]+")

	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []Sensor
	for s.Scan() {
		tokens := re.FindAllString(s.Text(), -1)
		sensorX, _ := strconv.ParseInt(tokens[0], 10, 32)
		sensorY, _ := strconv.ParseInt(tokens[1], 10, 32)
		beaconX, _ := strconv.ParseInt(tokens[2], 10, 32)
		beaconY, _ := strconv.ParseInt(tokens[3], 10, 32)

		res = append(res, Sensor{
			Pos: Pos{
				X: int(sensorX),
				Y: int(sensorY),
			},
			Beacon: Pos{
				X: int(beaconX),
				Y: int(beaconY),
			},
		})
	}

	return res
}

func solveFirst(input []Sensor, queryY int) int {
	beaconPos := make(map[Pos]bool)

	// Evaluate minX and maxX
	minX, maxX := math.MaxInt, math.MinInt
	for _, s := range input {
		beaconPos[s.Beacon] = true

		distToBeacon := dist(s.Pos, s.Beacon)
		yDistToQuery := abs(queryY - s.Pos.Y)
		if yDistToQuery > distToBeacon {
			// Further than dist to beacon
			continue
		}
		xDistToQuery := distToBeacon - yDistToQuery

		minX, maxX = min(minX, s.Pos.X-xDistToQuery), max(maxX, s.Pos.X+xDistToQuery)
	}

	//fmt.Println(minX, maxX)

	// Check each position between [minX, maxX] if it is nearer to any sensor than the corresponding beacon
	var res int
	for x := minX; x <= maxX; x++ {
		pos := Pos{X: x, Y: queryY}
		if beaconPos[pos] {
			continue
		}

		for _, s := range input {
			if dist(pos, s.Pos) <= dist(s.Pos, s.Beacon) {
				res++
				break
			}
		}
	}

	return res
}

func solveSecond(input []Sensor, searchSpace int) int {
	ranges := make([][]Range, searchSpace+1)
	for y := 0; y <= searchSpace; y++ {
		for _, s := range input {
			distToBeacon := dist(s.Pos, s.Beacon)
			yDistToQuery := abs(y - s.Pos.Y)
			if yDistToQuery > distToBeacon {
				// Further than dist to beacon
				continue
			}
			xDistToQuery := distToBeacon - yDistToQuery

			ranges[y] = append(ranges[y], Range{
				Start: s.Pos.X - xDistToQuery,
				End:   s.Pos.X + xDistToQuery,
			})
		}

		if len(ranges[y]) == 0 {
			panic(fmt.Sprintf("no range on Y = %d", y))
		}

		// Sort by range.Start
		sort.Slice(ranges[y], func(i, j int) bool {
			return ranges[y][i].Start < ranges[y][j].Start
		})

		// Merge ranges
		var lastIdx int
		for i := 1; i < len(ranges[y]); i++ {
			last := ranges[y][lastIdx]
			cur := ranges[y][i]
			if last.End+1 >= cur.Start {
				// Overlapping ranges, merge into last
				ranges[y][lastIdx] = Range{
					Start: last.Start,
					End:   max(last.End, cur.End),
				}
			} else {
				// Non-overlapping ranges, finalize last, move last pointer
				lastIdx++
				ranges[y][lastIdx] = cur
			}
		}
		ranges[y] = ranges[y][:lastIdx+1]

		//fmt.Println(ranges[queryY])

		if len(ranges[y]) == 2 {
			//fmt.Println(y, ranges[y])
			return (ranges[y][1].Start-1)*4000000 + y
		}
	}

	return 0
}

func dist(first, second Pos) int {
	return abs(second.X-first.X) + abs(second.Y-first.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
