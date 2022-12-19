package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type ResourceType string

type Blueprint struct {
	OreRobotCost      [3]int
	ClayRobotCost     [3]int
	ObsidianRobotCost [3]int
	GeodeRobotCost    [3]int
	MaxRobots         [3]int
}

func main() {
	input := parseInput("in.txt")
	//fmt.Println(input)
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) []Blueprint {
	re := regexp.MustCompile("Blueprint \\d+: Each ore robot costs (\\d+) ore. Each clay robot costs (\\d+) ore. Each obsidian robot costs (\\d+) ore and (\\d+) clay. Each geode robot costs (\\d+) ore and (\\d+) obsidian.")

	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []Blueprint
	for s.Scan() {
		tokens := re.FindStringSubmatch(s.Text())

		oreRobotOreCost, _ := strconv.ParseInt(tokens[1], 10, 32)
		clayRobotOreCost, _ := strconv.ParseInt(tokens[2], 10, 32)
		obsidianRobotOreCost, _ := strconv.ParseInt(tokens[3], 10, 32)
		obsidianRobotClayCost, _ := strconv.ParseInt(tokens[4], 10, 32)
		geodeRobotOreCost, _ := strconv.ParseInt(tokens[5], 10, 32)
		geodeRobotObsidianCost, _ := strconv.ParseInt(tokens[6], 10, 32)

		bp := Blueprint{
			OreRobotCost:      [3]int{int(oreRobotOreCost), 0, 0},
			ClayRobotCost:     [3]int{int(clayRobotOreCost), 0, 0},
			ObsidianRobotCost: [3]int{int(obsidianRobotOreCost), int(obsidianRobotClayCost), 0},
			GeodeRobotCost:    [3]int{int(geodeRobotOreCost), 0, int(geodeRobotObsidianCost)},
		}

		// Optimization: we never need to make more than max cost robots of each resource type
		bp.MaxRobots = int3Max(bp.OreRobotCost, bp.ClayRobotCost, bp.ObsidianRobotCost, bp.GeodeRobotCost)

		res = append(res, bp)
	}

	return res
}

func solveFirst(input []Blueprint) int {
	var res int
	for i, bp := range input {
		res += (i + 1) * dp(bp, State{Time: 24, Robots: [3]int{1, 0, 0}, Resources: [3]int{0, 0, 0}}, make(map[State]int))
	}

	return res
}

func solveSecond(input []Blueprint) int {
	res := 1
	for _, bp := range input[:3] {
		res *= dp(bp, State{Time: 32, Robots: [3]int{1, 0, 0}, Resources: [3]int{0, 0, 0}}, make(map[State]int))
	}

	return res
}

type State struct {
	Time      int
	Robots    [3]int
	Resources [3]int
}

// Maximize number of produced geodes.
func dp(bp Blueprint, state State, memo map[State]int) int {
	if state.Time == 0 {
		return 0
	}

	if res, ok := memo[state]; ok {
		return res
	}

	var res int
	if !int3GreaterEquals(state.Resources, bp.MaxRobots) {
		// Do nothing
		res = max(res, dp(bp, State{
			Time:      state.Time - 1,
			Robots:    state.Robots,
			Resources: int3Add(state.Resources, state.Robots),
		}, memo))
	}
	if int3GreaterEquals(state.Resources, bp.OreRobotCost) && state.Robots[0] < bp.MaxRobots[0] {
		// Create ore robot
		res = max(res, dp(bp, State{
			Time:      state.Time - 1,
			Robots:    int3Add(state.Robots, [3]int{1, 0, 0}),
			Resources: int3Add(int3Sub(state.Resources, bp.OreRobotCost), state.Robots),
		}, memo))
	}
	if int3GreaterEquals(state.Resources, bp.ClayRobotCost) && state.Robots[1] < bp.MaxRobots[1] {
		// Create clay robot
		res = max(res, dp(bp, State{
			Time:      state.Time - 1,
			Robots:    int3Add(state.Robots, [3]int{0, 1, 0}),
			Resources: int3Add(int3Sub(state.Resources, bp.ClayRobotCost), state.Robots),
		}, memo))
	}
	if int3GreaterEquals(state.Resources, bp.ObsidianRobotCost) && state.Robots[2] < bp.MaxRobots[2] {
		// Create obsidian robot
		res = max(res, dp(bp, State{
			Time:      state.Time - 1,
			Robots:    int3Add(state.Robots, [3]int{0, 0, 1}),
			Resources: int3Add(int3Sub(state.Resources, bp.ObsidianRobotCost), state.Robots),
		}, memo))
	}
	if int3GreaterEquals(state.Resources, bp.GeodeRobotCost) {
		// Create geode robot
		res = max(res, dp(bp, State{
			Time:      state.Time - 1,
			Robots:    state.Robots,
			Resources: int3Add(int3Sub(state.Resources, bp.GeodeRobotCost), state.Robots),
		}, memo)+(state.Time-1)) // express directly as total number of geodes produced
	}

	memo[state] = res
	return res
}

func int3Add(first [3]int, second [3]int) [3]int {
	return [3]int{first[0] + second[0], first[1] + second[1], first[2] + second[2]}
}

func int3Sub(first [3]int, second [3]int) [3]int {
	return [3]int{first[0] - second[0], first[1] - second[1], first[2] - second[2]}
}

func int3GreaterEquals(first [3]int, second [3]int) bool {
	return first[0] >= second[0] && first[1] >= second[1] && first[2] >= second[2]
}

func int3Max(eles ...[3]int) [3]int {
	res := [3]int{}
	for _, ele := range eles {
		res[0], res[1], res[2] = max(res[0], ele[0]), max(res[1], ele[1]), max(res[2], ele[2])
	}

	return res
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
