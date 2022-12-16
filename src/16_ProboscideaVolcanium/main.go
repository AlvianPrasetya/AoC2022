package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Valve struct {
	Idx  int
	Flow int
	Ins  map[int]bool
	Outs map[int]bool
}

func main() {
	input := parseInput("in.txt")
	fmt.Printf("%+v\n", input)
	//fmt.Println(solveFirst(valveMap, valves, 30))
	fmt.Println(solveSecond(input, 26))
}

func parseInput(in string) []Valve {
	re := regexp.MustCompile("Valve (?P<valve_id>[A-Z]{2}) has flow rate=(?P<flow_rate>\\d+); tunnel[s]? lead[s]? to valve[s]? (?P<connected_valves>[A-Z, ]+)")

	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	valveIdx := make(map[string]int)
	var edges [][2]string
	var valves []Valve
	for s.Scan() {
		tokens := re.FindStringSubmatch(s.Text())
		name := tokens[1]
		flow, _ := strconv.ParseInt(tokens[2], 10, 32)
		outs := strings.Split(tokens[3], ", ")

		valveIdx[name] = len(valves)
		valves = append(valves, Valve{
			Idx:  len(valves),
			Flow: int(flow),
		})

		for _, out := range outs {
			edges = append(edges, [2]string{name, out})
		}
	}

	for _, edge := range edges {
		valves[valveIdx[edge[0]]].Outs[valveIdx[edge[1]]] = true
		valves[valveIdx[edge[1]]].Ins[valveIdx[edge[0]]] = true
	}

	return valves
}

type State struct {
	Cur           int
	Prev          int
	TimeLeft      int
	OpenedBitmask uint64
}

type PairState struct {
	FirstCur      int
	FirstPrev     int
	SecondCur     int
	SecondPrev    int
	TimeLeft      int
	OpenedBitmask uint64
}

func solveFirst(valves []Valve, time int) int {
	return dfs(valves, State{
		Cur:           0,
		Prev:          -1,
		TimeLeft:      time,
		OpenedBitmask: 0,
	}, make(map[State]int))
}

func dfs(valves []Valve, cur State, memo map[State]int) int {
	if cur.TimeLeft == 0 {
		return 0
	}

	if res, ok := memo[cur]; ok {
		return res
	}

	states, values := evalStatesValues(valves, cur)

	var res int
	for i := 0; i < len(states); i++ {
		res = max(res, dfs(valves, states[i], memo)+values[i])
	}

	memo[cur] = res
	return res
}

func solveSecond(valves []Valve, time int) int {
	return 0
}

func dfsPair(valveMap map[string]int, valves []Valve, cur PairState, memo map[PairState]int) int {
	if cur.TimeLeft == 0 {
		return 0
	}

	if res, ok := memo[cur]; ok {
		return res
	}

	firstStates, firstValues := evalStatesValues(valves, State{
		Cur:           cur.FirstCur,
		Prev:          cur.FirstPrev,
		TimeLeft:      cur.TimeLeft,
		OpenedBitmask: cur.OpenedBitmask,
	})
	secondStates, secondValues := evalStatesValues(valves, State{
		Cur:           cur.SecondCur,
		Prev:          cur.SecondPrev,
		TimeLeft:      cur.TimeLeft,
		OpenedBitmask: cur.OpenedBitmask,
	})

	var res int
	for i := 0; i < len(firstStates); i++ {
		for j := 0; j < len(secondStates); j++ {
			firstState := firstStates[i]
			secondState := secondStates[j]

			// Invalid pruning
			/*if firstState.Idx == secondState.Idx && firstState.PrevIdx == secondState.PrevIdx {
				continue
			}*/

			if firstState.Cur > secondState.Cur {
				firstState, secondState = secondState, firstState
			}

			res = max(res, dfsPair(valveMap, valves, PairState{
				FirstCur:      firstState.Cur,
				FirstPrev:     firstState.Prev,
				SecondCur:     secondState.Cur,
				SecondPrev:    secondState.Prev,
				TimeLeft:      cur.TimeLeft - 1,
				OpenedBitmask: firstState.OpenedBitmask | secondState.OpenedBitmask,
			}, memo)+firstValues[i]+secondValues[j])
		}
	}

	memo[cur] = res
	return res
}

func evalStatesValues(valves []Valve, state State) ([]State, []int) {
	var states []State
	var values []int
	if state.OpenedBitmask&(1<<state.Cur) == 0 && valves[state.Cur].Flow != 0 {
		states = append(states, State{
			Cur:           state.Cur,
			Prev:          state.Cur,
			TimeLeft:      state.TimeLeft - 1,
			OpenedBitmask: state.OpenedBitmask | (1 << state.Cur),
		})
		values = append(values, valves[state.Cur].Flow*(state.TimeLeft-1))
	}
	for _, out := range valves[state.Cur].Outs {
		if out == state.Prev {
			continue
		}

		states = append(states, State{
			Cur:           out,
			Prev:          state.Cur,
			TimeLeft:      state.TimeLeft - 1,
			OpenedBitmask: state.OpenedBitmask,
		})
		values = append(values, 0)
	}

	return states, values
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
