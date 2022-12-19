package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Pos struct {
	X int
	Y int
}

func (this Pos) Add(other Pos) Pos {
	return Pos{
		X: this.X + other.X,
		Y: this.Y + other.Y,
	}
}

var (
	rockTemplates = [][]Pos{
		{
			{X: 2, Y: 4},
			{X: 3, Y: 4},
			{X: 4, Y: 4},
			{X: 5, Y: 4},
		},
		{
			{X: 3, Y: 4},
			{X: 2, Y: 5},
			{X: 3, Y: 5},
			{X: 4, Y: 5},
			{X: 3, Y: 6},
		},
		{
			{X: 2, Y: 4},
			{X: 3, Y: 4},
			{X: 4, Y: 4},
			{X: 4, Y: 5},
			{X: 4, Y: 6},
		},
		{
			{X: 2, Y: 4},
			{X: 2, Y: 5},
			{X: 2, Y: 6},
			{X: 2, Y: 7},
		},
		{
			{X: 2, Y: 4},
			{X: 3, Y: 4},
			{X: 2, Y: 5},
			{X: 3, Y: 5},
		},
	}

	dirMap = map[string]Pos{
		"<": {X: -1},
		">": {X: 1},
	}
)

func main() {
	input := parseInput("in.txt")
	//fmt.Println(input)
	fmt.Println(solveFirst(input, 2022))
	fmt.Println(solveSecond(input, 1000000000000))
}

func parseInput(in string) []Pos {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []Pos
	for s.Scan() {
		tokens := strings.Split(s.Text(), "")
		for _, token := range tokens {
			res = append(res, dirMap[token])
		}
	}

	return res
}

func solveFirst(jetDirs []Pos, rockCount int) int {
	// debugStep := 10

	var board [][7]byte
	jetIdx := 0
	for r := 0; r < rockCount; r++ {
		rock := emplaceRock(&board, rockTemplates[r%len(rockTemplates)])
		// if r <= debugStep {
		// 	printBoard(board)
		// }

		for i := 0; ; i++ {
			if i%2 == 0 {
				moveRock(&board, rock, jetDirs[jetIdx])
				jetIdx = (jetIdx + 1) % len(jetDirs)
			} else {
				if !moveRock(&board, rock, Pos{Y: -1}) {
					// Can't fall down anymore
					settleRock(board, rock)
					break
				}
			}

			// if r <= debugStep {
			// 	printBoard(board)
			// }
		}

		// if r <= debugStep {
		// 	printBoard(board)
		// }
	}

	return len(board)
}

func solveSecond(jetDirs []Pos, rockCount uint64) uint64 {
	// Find periodic cycle length and height
	var cycleOffset int
	var cycleOffsetHeight int
	var cycleLength int
	var cycleHeight int
	var cycleJetIdx int
	var board [][7]byte
	jetIdx := 0
	jetIdxOccurence := make(map[int][2]int)
	for r := 0; uint64(r) < rockCount; r++ {
		heightBefore := len(board)
		jetIdxBefore := jetIdx

		rock := emplaceRock(&board, rockTemplates[r%len(rockTemplates)])

		for i := 0; ; i++ {
			if i%2 == 0 {
				moveRock(&board, rock, jetDirs[jetIdx])
				jetIdx = (jetIdx + 1) % len(jetDirs)
			} else {
				if !moveRock(&board, rock, Pos{Y: -1}) {
					// Can't fall down anymore
					settleRock(board, rock)
					break
				}
			}
		}

		if r != 0 && r%len(rockTemplates) == 0 && len(board) == heightBefore+1 {
			if prev, ok := jetIdxOccurence[jetIdxBefore]; ok {
				cycleOffset = prev[0]
				cycleOffsetHeight = prev[1]
				cycleLength = r - cycleOffset
				cycleHeight = heightBefore - prev[1]
				cycleJetIdx = jetIdxBefore
				break
			}
			jetIdxOccurence[jetIdxBefore] = [2]int{r, heightBefore}
		}
	}
	// fmt.Println(cycleOffset, cycleOffsetHeight, cycleLength, cycleHeight, cycleJetIdx)

	// Simulate (rockCount % cycleLength) rocks
	cycleCount := (rockCount - uint64(cycleOffset)) / uint64(cycleLength)
	rockCount = (rockCount - uint64(cycleOffset)) % uint64(cycleLength)
	board = nil
	jetIdx = cycleJetIdx
	for r := 0; uint64(r) < rockCount; r++ {
		rock := emplaceRock(&board, rockTemplates[r%len(rockTemplates)])

		for i := 0; ; i++ {
			if i%2 == 0 {
				moveRock(&board, rock, jetDirs[jetIdx])
				jetIdx = (jetIdx + 1) % len(jetDirs)
			} else {
				if !moveRock(&board, rock, Pos{Y: -1}) {
					// Can't fall down anymore
					settleRock(board, rock)
					break
				}
			}
		}
	}

	return uint64(cycleOffsetHeight) + cycleCount*uint64(cycleHeight) + uint64(len(board))
}

func emplaceRock(board *[][7]byte, template []Pos) []Pos {
	height := len(*board)

	rock := make([]Pos, len(template))
	for i, shard := range template {
		rock[i] = shard.Add(Pos{Y: height - 1})
		// Extend board vertically
		for rock[i].Y >= len(*board) {
			*board = append(*board, [7]byte{})
		}
		(*board)[rock[i].Y][rock[i].X] = '@'
	}

	return rock
}

func moveRock(board *[][7]byte, rock []Pos, dir Pos) bool {
	canMove := true
	for _, shard := range rock {
		movedShard := shard.Add(dir)
		if movedShard.Y < 0 || movedShard.Y == len(*board) || movedShard.X < 0 || movedShard.X == 7 {
			// Out-of-bounds
			canMove = false
			break
		}
		if (*board)[movedShard.Y][movedShard.X] == '#' {
			// Hit another settled rock
			canMove = false
			break
		}
	}

	if !canMove {
		return false
	}

	for _, shard := range rock {
		(*board)[shard.Y][shard.X] = 0
	}

	for i := range rock {
		rock[i] = rock[i].Add(dir)
		(*board)[rock[i].Y][rock[i].X] = '@'
	}

	// Shrink board vertically
	for isOccupied := false; !isOccupied; {
		for i := 0; i < 7; i++ {
			if (*board)[len(*board)-1][i] != 0 {
				isOccupied = true
				break
			}
		}

		if !isOccupied {
			*board = (*board)[:len(*board)-1]
		}
	}

	return true
}

func settleRock(board [][7]byte, rock []Pos) {
	for _, shard := range rock {
		board[shard.Y][shard.X] = '#'
	}
}

func printBoard(board [][7]byte) {
	var boardStrs []string
	for _, row := range board {
		var rowStr strings.Builder
		for _, cell := range row {
			if cell == 0 {
				cell = '.'
			}
			rowStr.WriteByte(cell)
		}
		boardStrs = append(boardStrs, rowStr.String())
	}

	var res strings.Builder
	for i := len(boardStrs) - 1; i >= 0; i-- {
		res.WriteString(boardStrs[i])
		res.WriteByte('\n')
	}

	fmt.Println(res.String())
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
