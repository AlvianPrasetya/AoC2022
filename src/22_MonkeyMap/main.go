package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

var (
	Dirs = []Pos{{X: 1}, {Y: 1}, {X: -1}, {Y: -1}}
)

type Instr interface {
	Execute(state *State, board [][]byte, xRanges []*Range, yRanges []*Range)
}

type InstrTurnLeft struct {
}

func (i *InstrTurnLeft) Execute(state *State, board [][]byte, xRanges []*Range, yRanges []*Range) {
	state.DirIdx = (state.DirIdx - 1 + len(Dirs)) % len(Dirs)
}

func (i InstrTurnLeft) String() string {
	return "L"
}

type InstrTurnRight struct {
}

func (i *InstrTurnRight) Execute(state *State, board [][]byte, xRanges []*Range, yRanges []*Range) {
	state.DirIdx = (state.DirIdx + 1) % len(Dirs)
}

func (i InstrTurnRight) String() string {
	return "R"
}

type InstrMove struct {
	Steps int
}

func (i *InstrMove) Execute(state *State, board [][]byte, xRanges []*Range, yRanges []*Range) {
	dir := Dirs[state.DirIdx]
	xRange := xRanges[state.Pos.Y]
	yRange := yRanges[state.Pos.X]
	for s := 0; s < i.Steps; s++ {
		nextPos := state.Pos.Add(dir)
		if state.DirIdx%2 == 0 {
			// Horizontal movement
			len := xRange.Max - xRange.Min + 1
			nextPos.X = xRange.Min + ((nextPos.X - xRange.Min + len) % len)
		} else {
			// Vertical movement
			len := yRange.Max - yRange.Min + 1
			nextPos.Y = yRange.Min + ((nextPos.Y - yRange.Min + len) % len)
		}
		if board[nextPos.Y][nextPos.X] == '#' {
			// Hit a wall
			break
		} else {
			state.Pos = &nextPos
		}
	}
}

func (i InstrMove) String() string {
	return strconv.Itoa(i.Steps)
}

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

func (pos Pos) String() string {
	return fmt.Sprintf("(%d, %d)", pos.Y, pos.X)
}

type Range struct {
	Min int
	Max int
}

type State struct {
	Pos    *Pos
	DirIdx int
}

func main() {
	board, instrs := parseInput("in.txt")
	fmt.Println(board, instrs)
	xRanges, yRanges := getRanges(board)
	fmt.Println(solveFirst(board, instrs, xRanges, yRanges))
}

func parseInput(in string) ([][]byte, []Instr) {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var board [][]byte
	var boardWidth int
	var instrs []Instr
	for s.Scan() {
		str := s.Text()
		board = append(board, []byte(str))
		boardWidth = max(boardWidth, len(str))

		if str == "" {
			s.Scan()
			instrStr := s.Text()

			for i := 0; i < len(instrStr); i++ {
				switch instrStr[i] {
				case 'L':
					instrs = append(instrs, &InstrTurnLeft{})
				case 'R':
					instrs = append(instrs, &InstrTurnRight{})
				default:
					j := i
					for (j+1 < len(instrStr)) && ('0' <= instrStr[j+1] && instrStr[j+1] <= '9') {
						j++
					}
					val, _ := strconv.Atoi(instrStr[i : j+1])
					instrs = append(instrs, &InstrMove{
						Steps: val,
					})
					i = j
				}
			}

			break
		}
	}

	// Right-pad board
	for i := range board {
		for len(board[i]) < boardWidth {
			board[i] = append(board[i], ' ')
		}
	}

	return board, instrs
}

func getRanges(board [][]byte) ([]*Range, []*Range) {
	xRanges := make([]*Range, len(board))
	for i := 0; i < len(board); i++ {
		xRanges[i] = &Range{Min: math.MaxInt, Max: math.MinInt}
	}
	yRanges := make([]*Range, len(board[0]))
	for i := 0; i < len(board[0]); i++ {
		yRanges[i] = &Range{Min: math.MaxInt, Max: math.MinInt}
	}

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			if board[i][j] != ' ' {
				xRanges[i].Min, xRanges[i].Max = min(xRanges[i].Min, j), max(xRanges[i].Max, j)
				yRanges[j].Min, yRanges[j].Max = min(yRanges[j].Min, i), max(yRanges[j].Max, i)
			}
		}
	}

	return xRanges, yRanges
}

func solveFirst(board [][]byte, instrs []Instr, xRanges []*Range, yRanges []*Range) int {
	state := &State{
		Pos:    &Pos{X: xRanges[0].Min},
		DirIdx: 0,
	}

	for _, instr := range instrs {
		instr.Execute(state, board, xRanges, yRanges)
		// fmt.Println(state.Pos, state.DirIdx)
	}

	return 1000*(state.Pos.Y+1) + 4*(state.Pos.X+1) + state.DirIdx
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
