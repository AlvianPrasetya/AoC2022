package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	MoveRock     = Move(0)
	MovePaper    = Move(1)
	MoveScissors = Move(2)

	OutcomeLose = Move(0)
	OutcomeDraw = Move(1)
	OutcomeWin  = Move(2)
)

var (
	moveScore = map[Move]int{
		MoveRock:     1,
		MovePaper:    2,
		MoveScissors: 3,
	}
	outcomeScore = map[Turn]int{
		{MoveRock, MoveRock}:         3,
		{MoveRock, MovePaper}:        6,
		{MoveRock, MoveScissors}:     0,
		{MovePaper, MoveRock}:        0,
		{MovePaper, MovePaper}:       3,
		{MovePaper, MoveScissors}:    6,
		{MoveScissors, MoveRock}:     6,
		{MoveScissors, MovePaper}:    0,
		{MoveScissors, MoveScissors}: 3,
	}
	myMoveMap = map[Turn]Move{
		{MoveRock, OutcomeLose}:     MoveScissors,
		{MoveRock, OutcomeDraw}:     MoveRock,
		{MoveRock, OutcomeWin}:      MovePaper,
		{MovePaper, OutcomeLose}:    MoveRock,
		{MovePaper, OutcomeDraw}:    MovePaper,
		{MovePaper, OutcomeWin}:     MoveScissors,
		{MoveScissors, OutcomeLose}: MovePaper,
		{MoveScissors, OutcomeDraw}: MoveScissors,
		{MoveScissors, OutcomeWin}:  MoveRock,
	}
)

type Move int

type Turn struct {
	OppMove Move
	MyMove  Move
}

func main() {
	input := parseInput("in.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) []Turn {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []Turn
	for s.Scan() {
		tokens := strings.Split(s.Text(), " ")

		var cur Turn
		switch tokens[0] {
		case "A":
			cur.OppMove = MoveRock
		case "B":
			cur.OppMove = MovePaper
		case "C":
			cur.OppMove = MoveScissors
		}
		switch tokens[1] {
		case "X":
			cur.MyMove = MoveRock
		case "Y":
			cur.MyMove = MovePaper
		case "Z":
			cur.MyMove = MoveScissors
		}

		res = append(res, cur)
	}

	return res
}

func solveFirst(input []Turn) int {
	var res int
	for _, turn := range input {
		res += moveScore[turn.MyMove] + outcomeScore[turn]
	}

	return res
}

func solveSecond(input []Turn) int {
	var res int
	for _, turn := range input {
		turn = Turn{
			OppMove: turn.OppMove,
			MyMove:  myMoveMap[turn],
		}

		res += moveScore[turn.MyMove] + outcomeScore[turn]
	}

	return res
}
