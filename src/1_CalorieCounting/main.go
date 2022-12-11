package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	input := parseInput("in.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) [][]int {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res [][]int
	var cur []int
	for s.Scan() {
		if s.Text() == "" {
			res = append(res, cur)
			cur = []int{}
		} else {
			food, _ := strconv.ParseInt(s.Text(), 10, 32)
			cur = append(cur, int(food))
		}
	}
	if len(cur) != 0 {
		res = append(res, cur)
	}

	return res
}

func solveFirst(input [][]int) int {
	var res int
	for _, foods := range input {
		var sum int
		for _, food := range foods {
			sum += food
		}

		res = maxInt(res, sum)
	}

	return res
}

func solveSecond(input [][]int) int {
	h := &MinHeap{}
	for _, foods := range input {
		var sum int
		for _, food := range foods {
			sum += food
		}

		heap.Push(h, sum)
		if h.Len() > 3 {
			heap.Pop(h)
		}
	}

	var res int
	for _, sum := range *h {
		res += sum
	}

	return res
}

func maxInt(i, j int) int {
	if i > j {
		return i
	}
	return j
}
