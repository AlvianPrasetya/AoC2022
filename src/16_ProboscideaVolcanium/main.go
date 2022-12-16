package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Valve struct {
	Idx  int
	Name string
	Flow int
	Ins  map[int]int
	Outs map[int]int
}

func (v Valve) String() string {
	return fmt.Sprintf("{ Idx: %d, Name: %s, Flow: %d, Ins: %v, Outs: %v }\n", v.Idx, v.Name, v.Flow, v.Ins, v.Outs)
}

func main() {
	input := parseInput("in.txt")
	//fmt.Printf("%+v\n", input)
	valves := pruneZeros(input)
	//fmt.Printf("%+v\n", valves)
	dists := evalDists(valves)
	//fmt.Printf("%+v\n", dists)
	fmt.Println(solveFirst(valves, dists, 30))
	fmt.Println(solveSecond(valves, dists, 26))
}

func parseInput(in string) []*Valve {
	re := regexp.MustCompile("Valve (?P<valve_id>[A-Z]{2}) has flow rate=(?P<flow_rate>\\d+); tunnel[s]? lead[s]? to valve[s]? (?P<connected_valves>[A-Z, ]+)")

	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	valveIdx := make(map[string]int)
	var edges [][2]string
	var valves []*Valve
	for s.Scan() {
		tokens := re.FindStringSubmatch(s.Text())
		name := tokens[1]
		flow, _ := strconv.ParseInt(tokens[2], 10, 32)
		outs := strings.Split(tokens[3], ", ")

		valveIdx[name] = len(valves)
		valves = append(valves, &Valve{
			Idx:  len(valves),
			Name: name,
			Flow: int(flow),
			Ins:  make(map[int]int),
			Outs: make(map[int]int),
		})

		for _, out := range outs {
			edges = append(edges, [2]string{name, out})
		}
	}

	for _, edge := range edges {
		valves[valveIdx[edge[0]]].Outs[valveIdx[edge[1]]] = 1
		valves[valveIdx[edge[1]]].Ins[valveIdx[edge[0]]] = 1
	}

	return valves
}

func pruneZeros(valves []*Valve) []*Valve {
	type IdxDist struct {
		Idx  int
		Dist int
	}

	queue := list.New()
	for i := range valves {
		if valves[i].Flow == 0 && valves[i].Name != "AA" {
			ins := make(map[int]int)
			queue.PushBack(IdxDist{Idx: i, Dist: 0})
			isVisited := make(map[int]bool)
			for queue.Len() != 0 {
				cur := queue.Remove(queue.Front()).(IdxDist)
				if isVisited[cur.Idx] {
					continue
				}
				isVisited[cur.Idx] = true

				for in := range valves[cur.Idx].Ins {
					if valves[in].Flow == 0 && valves[in].Name != "AA" {
						queue.PushBack(IdxDist{Idx: in, Dist: cur.Dist + 1})
					} else {
						if _, ok := ins[in]; !ok || ins[in] > cur.Dist+1 {
							ins[in] = cur.Dist + 1
						}
					}
				}
			}

			outs := make(map[int]int)
			queue.PushBack(IdxDist{Idx: i, Dist: 0})
			isVisited = make(map[int]bool)
			for queue.Len() != 0 {
				cur := queue.Remove(queue.Front()).(IdxDist)
				if isVisited[cur.Idx] {
					continue
				}
				isVisited[cur.Idx] = true

				for out := range valves[cur.Idx].Outs {
					if valves[out].Flow == 0 && valves[out].Name != "AA" {
						queue.PushBack(IdxDist{Idx: out, Dist: cur.Dist + 1})
					} else {
						if _, ok := outs[out]; !ok || outs[out] > cur.Dist+1 {
							outs[out] = cur.Dist + 1
						}
					}
				}
			}

			for in, inDist := range ins {
				for out, outDist := range outs {
					if in == out {
						continue
					}
					if _, ok := valves[in].Outs[out]; !ok || valves[in].Outs[out] > inDist+outDist {
						valves[in].Outs[out] = inDist + outDist
					}
					if _, ok := valves[out].Ins[in]; !ok || valves[out].Ins[in] > inDist+outDist {
						valves[out].Ins[in] = inDist + outDist
					}
				}
			}
		}
	}

	newIndexMap := make(map[int]int)
	for i := range valves {
		if valves[i].Flow == 0 && valves[i].Name != "AA" {
			continue
		}
		newIndexMap[i] = len(newIndexMap)
	}

	var newValves []*Valve
	for i := range valves {
		if valves[i].Flow == 0 && valves[i].Name != "AA" {
			continue
		}

		valves[i].Idx = newIndexMap[i]

		newIns := make(map[int]int)
		for in, dist := range valves[i].Ins {
			if _, ok := newIndexMap[in]; !ok {
				continue
			}

			newIns[newIndexMap[in]] = dist
		}
		valves[i].Ins = newIns

		newOuts := make(map[int]int)
		for out, dist := range valves[i].Outs {
			if _, ok := newIndexMap[out]; !ok {
				continue
			}

			newOuts[newIndexMap[out]] = dist
		}
		valves[i].Outs = newOuts

		newValves = append(newValves, valves[i])
	}

	return newValves
}

func evalDists(valves []*Valve) [][]int {
	dists := make([][]int, len(valves))
	for i := 0; i < len(valves); i++ {
		dists[i] = make([]int, len(valves))
		for j := 0; j < len(valves); j++ {
			if i == j {
				dists[i][j] = 0
			} else {
				dists[i][j] = 10000
			}
		}
	}

	for i := 0; i < len(valves); i++ {
		for out, dist := range valves[i].Outs {
			dists[i][out] = dist
		}
	}

	for k := 0; k < len(valves); k++ {
		for i := 0; i < len(valves); i++ {
			for j := 0; j < len(valves); j++ {
				if dists[i][k]+dists[k][j] < dists[i][j] {
					dists[i][j] = dists[i][k] + dists[k][j]
				}
			}
		}
	}

	return dists
}

type State struct {
	Cur           int
	Time          int
	OpenedBitmask uint64
}

func solveFirst(valves []*Valve, dists [][]int, time int) int {
	var startIdx int
	for _, valve := range valves {
		if valve.Name == "AA" {
			startIdx = valve.Idx
		}
	}

	return dfs(valves, dists, State{
		Cur:           startIdx,
		Time:          time,
		OpenedBitmask: 0,
	}, make(map[State]int))
}

func dfs(valves []*Valve, dists [][]int, cur State, memo map[State]int) int {
	if cur.Time == 0 {
		return 0
	}

	if res, ok := memo[cur]; ok {
		return res
	}

	states, values := evalStatesValues(valves, dists, cur)

	var res int
	for i, state := range states {
		res = max(res, dfs(valves, dists, state, memo)+values[i])
	}

	memo[cur] = res
	return res
}

type PairState struct {
	FirstCur      int
	FirstTime     int
	SecondCur     int
	SecondTime    int
	OpenedBitmask uint64
}

func solveSecond(valves []*Valve, dists [][]int, time int) int {
	var startIdx int
	for _, valve := range valves {
		if valve.Name == "AA" {
			startIdx = valve.Idx
		}
	}

	return dfsPair(valves, dists, PairState{
		FirstCur:      startIdx,
		FirstTime:     time,
		SecondCur:     startIdx,
		SecondTime:    time,
		OpenedBitmask: 0,
	}, make(map[PairState]int))
}

func dfsPair(valves []*Valve, dists [][]int, cur PairState, memo map[PairState]int) int {
	curTime := max(cur.FirstTime, cur.SecondTime)
	if curTime == 0 {
		return 0
	}

	if res, ok := memo[cur]; ok {
		return res
	}

	firstStates, firstValues := evalStatesValues(valves, dists, State{
		Cur:           cur.FirstCur,
		Time:          cur.FirstTime,
		OpenedBitmask: cur.OpenedBitmask,
	})

	secondStates, secondValues := evalStatesValues(valves, dists, State{
		Cur:           cur.SecondCur,
		Time:          cur.SecondTime,
		OpenedBitmask: cur.OpenedBitmask,
	})

	var res int
	if len(firstStates) != 0 && len(secondStates) != 0 {
		for i, firstState := range firstStates {
			for j, secondState := range secondStates {
				if firstState.Cur == secondState.Cur {
					// Can't open the same valve
					continue
				}
				res = max(res, dfsPair(valves, dists, PairState{
					FirstCur:      firstState.Cur,
					FirstTime:     firstState.Time,
					SecondCur:     secondState.Cur,
					SecondTime:    secondState.Time,
					OpenedBitmask: firstState.OpenedBitmask | secondState.OpenedBitmask,
				}, memo)+firstValues[i]+secondValues[j])
			}
		}
	} else if len(firstStates) != 0 {
		for i, firstState := range firstStates {
			res = max(res, dfsPair(valves, dists, PairState{
				FirstCur:      firstState.Cur,
				FirstTime:     firstState.Time,
				SecondCur:     cur.SecondCur,
				SecondTime:    cur.SecondTime,
				OpenedBitmask: firstState.OpenedBitmask,
			}, memo)+firstValues[i])
		}
	} else if len(secondStates) != 0 {
		for i, secondState := range secondStates {
			res = max(res, dfsPair(valves, dists, PairState{
				FirstCur:      cur.FirstCur,
				FirstTime:     cur.FirstTime,
				SecondCur:     secondState.Cur,
				SecondTime:    secondState.Time,
				OpenedBitmask: secondState.OpenedBitmask,
			}, memo)+secondValues[i])
		}
	}

	memo[cur] = res
	return res
}

func evalStatesValues(valves []*Valve, dists [][]int, cur State) ([]State, []int) {
	var states []State
	var values []int

	for i := 0; i < len(valves); i++ {
		if cur.OpenedBitmask&(1<<i) == 0 {
			// Not opened yet
			timeNeeded := dists[cur.Cur][i] + 1
			if cur.Time < timeNeeded {
				// Not enough time
				continue
			}

			states = append(states, State{
				Cur:           i,
				Time:          cur.Time - timeNeeded,
				OpenedBitmask: cur.OpenedBitmask | (1 << i),
			})
			values = append(values, valves[i].Flow*(cur.Time-timeNeeded))
		}
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
