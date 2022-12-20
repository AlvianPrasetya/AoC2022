package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LinkedList struct {
	Head *Node
	Tail *Node
}

func (ll *LinkedList) AddNode(node *Node) {
	if ll.Head == nil {
		ll.Head = node
		ll.Tail = node
		node.Prev = node
		node.Next = node
		return
	}

	ll.Tail.Next = node
	node.Prev = ll.Tail
	node.Next = ll.Head
	ll.Head.Prev = node
	ll.Tail = node
}

func (ll LinkedList) String() string {
	if ll.Head == nil {
		return ""
	}

	var res strings.Builder
	for cur := ll.Head; ; cur = cur.Next {
		res.WriteString(cur.String())
		res.WriteString(" -> ")
		if cur.Next == ll.Head {
			break
		}
	}
	//res.WriteString(ll.Tail.String())

	return res.String()
}

type Node struct {
	Value int64
	Next  *Node
	Prev  *Node
}

func (cur *Node) SwapNext() {
	prev := cur.Prev
	next := cur.Next
	nextNext := cur.Next.Next
	// fmt.Println(prev, cur, next, nextNext)

	prev.Next = next
	cur.Prev = next
	cur.Next = nextNext
	next.Prev = prev
	next.Next = cur
	nextNext.Prev = cur
}

func (cur *Node) SwapPrev() {
	prev := cur.Prev
	prevPrev := cur.Prev.Prev
	next := cur.Next

	prevPrev.Next = cur
	prev.Prev = cur
	prev.Next = next
	cur.Prev = prevPrev
	cur.Next = prev
	next.Prev = prev
}

func (cur Node) String() string {
	return strconv.FormatInt(cur.Value, 10)
}

func main() {
	in := "in.txt"

	// input, ll, zero := parseInput(in)
	// fmt.Println(solveFirst(input, ll, zero))

	input, ll, zero := parseInput(in)
	fmt.Println(solveSecond(input, ll, zero, 811589153, 10))
}

func parseInput(in string) ([]*Node, *LinkedList, *Node) {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var input []*Node
	ll := &LinkedList{}
	var zero *Node
	for s.Scan() {
		val, _ := strconv.ParseInt(s.Text(), 10, 64)
		node := &Node{
			Value: val,
		}
		input = append(input, node)
		ll.AddNode(node)
		if val == 0 {
			zero = node
		}
	}

	return input, ll, zero
}

func solveFirst(input []*Node, ll *LinkedList, zero *Node) int64 {
	// fmt.Println(ll)
	for _, in := range input {
		for i := 0; int64(i) < in.Value; i++ {
			in.SwapNext()
		}
		for i := 0; int64(i) > in.Value; i-- {
			in.SwapPrev()
		}
		// fmt.Println(ll)
	}

	var res int64
	cur := zero
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			cur = cur.Next
		}
		res += cur.Value
	}

	return res
}

func solveSecond(input []*Node, ll *LinkedList, zero *Node, decryptionKey int64, cycleCount int) int64 {
	for _, in := range input {
		in.Value = in.Value * decryptionKey
	}

	// fmt.Println(ll)
	for i := 0; i < cycleCount; i++ {
		// fmt.Println("Cycle ", i+1)
		for _, in := range input {
			if in.Value == 0 {
				continue
			}

			var fn func()
			if in.Value < 0 {
				fn = in.SwapPrev
			} else {
				fn = in.SwapNext
			}

			swapTimes := int(abs(in.Value) % int64(len(input)-1))
			for j := 0; j < swapTimes; j++ {
				fn()
			}
			// fmt.Println(swapTimes, in.Value, ll)
		}
		// fmt.Println(ll)
	}

	var res int64
	cur := zero
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			cur = cur.Next
		}
		res += cur.Value
	}

	return res
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
