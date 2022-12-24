package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strings"
)

var (
	DirUp    = Pos{Y: -1}
	DirDown  = Pos{Y: 1}
	DirLeft  = Pos{X: -1}
	DirRight = Pos{X: 1}

	DirMap = map[byte]Pos{
		'^': DirUp,
		'v': DirDown,
		'<': DirLeft,
		'>': DirRight,
	}
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

func (pos Pos) Mul(mult int) Pos {
	return Pos{
		Y: pos.Y * mult,
		X: pos.X * mult,
	}
}

func (pos Pos) Mod(modY int, modX int) Pos {
	return Pos{
		Y: ((pos.Y % modY) + modY) % modY,
		X: ((pos.X % modX) + modX) % modX,
	}
}

type State struct {
	Pos  Pos
	Time int
}

func main() {
	board, startPos, endPos := parseInput("in.txt")
	// fmt.Println(printBoard(board))
	// fmt.Println(startPos, endPos)
	blizzardsAtTime := generateBlizzardsAtTime(board)
	fmt.Println(solveFirst(board, startPos, endPos, blizzardsAtTime))
	fmt.Println(solveSecond(board, startPos, endPos, blizzardsAtTime))
}

func parseInput(in string) ([][]byte, Pos, Pos) {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var board [][]byte
	for s.Scan() {
		row := []byte(s.Text())
		board = append(board, row)
	}

	var startPos, endPos Pos
	for i := 0; i < len(board[0]); i++ {
		if board[0][i] == '.' {
			startPos = Pos{Y: 0, X: i}
		}
		if board[len(board)-1][i] == '.' {
			endPos = Pos{Y: len(board) - 1, X: i}
		}
	}

	return board, startPos, endPos
}

func generateBlizzardsAtTime(board [][]byte) []map[Pos]bool {
	dy := len(board) - 2
	dx := len(board[0]) - 2
	maxTime := lcm(dy, dx)

	var blizzardsAtTime []map[Pos]bool
	for t := 0; t < maxTime; t++ {
		blizzards := make(map[Pos]bool)
		for y := 1; y < len(board)-1; y++ {
			for x := 1; x < len(board[0])-1; x++ {
				if board[y][x] != '.' {
					initPos := Pos{Y: y - 1, X: x - 1}
					dir := DirMap[board[y][x]]
					pos := initPos.Add(dir.Mul(t)).Mod(dy, dx).Add(Pos{Y: 1, X: 1})

					blizzards[pos] = true
				}
			}
		}

		blizzardsAtTime = append(blizzardsAtTime, blizzards)
	}

	return blizzardsAtTime
}

func solveFirst(board [][]byte, startPos Pos, endPos Pos, blizzardsAtTime []map[Pos]bool) int {
	return bfs(board, startPos, endPos, blizzardsAtTime, 0)
}

func solveSecond(board [][]byte, startPos Pos, endPos Pos, blizzardsAtTime []map[Pos]bool) int {
	x := bfs(board, startPos, endPos, blizzardsAtTime, 0)
	y := bfs(board, endPos, startPos, blizzardsAtTime, x)
	z := bfs(board, startPos, endPos, blizzardsAtTime, y)

	return z
}

func bfs(board [][]byte, startPos Pos, endPos Pos, blizzardsAtTime []map[Pos]bool, startTime int) int {
	queue := list.New()
	queue.PushBack(State{Pos: startPos, Time: startTime})
	visited := make(map[State]bool)
	for queue.Len() != 0 {
		cur := queue.Remove(queue.Front()).(State)
		relTime := cur.Time % len(blizzardsAtTime)

		if cur.Pos.Y < 0 || cur.Pos.Y == len(board) || cur.Pos.X < 0 || cur.Pos.X == len(board[0]) {
			continue
		}
		if board[cur.Pos.Y][cur.Pos.X] == '#' {
			continue
		}
		if blizzardsAtTime[relTime][cur.Pos] {
			continue
		}
		if visited[State{
			Pos:  cur.Pos,
			Time: relTime,
		}] {
			continue
		}
		visited[State{
			Pos:  cur.Pos,
			Time: relTime,
		}] = true

		// fmt.Println(cur.Pos, cur.Time)

		if cur.Pos == endPos {
			return cur.Time
		}

		for _, dir := range DirMap {
			nextPos := cur.Pos.Add(dir)
			queue.PushBack(State{
				Pos:  nextPos,
				Time: cur.Time + 1,
			})
		}
		queue.PushBack(State{
			Pos:  cur.Pos,
			Time: cur.Time + 1,
		})
	}

	return 0
}

func printBoard(board [][]byte) string {
	var res strings.Builder
	for _, row := range board {
		for _, cell := range row {
			res.WriteByte(cell)
		}
		res.WriteByte('\n')
	}

	return res.String()
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
