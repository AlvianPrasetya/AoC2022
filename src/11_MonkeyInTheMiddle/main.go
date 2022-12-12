package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	OperatorTypeAdd  = OperatorType("+")
	OperatorTypeMult = OperatorType("*")

	OperandTypeSelf  = OperandType("old")
	OperandTypeConst = OperandType("const")

	TestTypeDivBy = TestType("divisible by ")
)

type OperatorType string

type Operation struct {
	Type   OperatorType
	First  Operand
	Second Operand
}

type OperandType string

type Operand struct {
	Type  OperandType
	Value int
}

type TestType string

type Test struct {
	Type  TestType
	Value int
}

type Monkey struct {
	Items       []int
	Operation   Operation
	Test        Test
	TrueTarget  int
	FalseTarget int
}

func main() {
	in := "in.txt"

	input := parseInput(in)
	fmt.Println(solveFirst(input))

	input = parseInput(in)
	fmt.Println(solveSecond(input))

	input = parseInput(in)
	fmt.Println(solveSecondAlt(input))
}

func parseInput(in string) []*Monkey {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []*Monkey
	for s.Scan() {
		str := s.Text()

		if strings.HasPrefix(str, "Monkey") {
			var cur Monkey

			s.Scan()
			cur.Items = parseItems(s.Text())

			s.Scan()
			cur.Operation = parseOperation(s.Text())

			s.Scan()
			cur.Test = parseTest(s.Text())

			s.Scan()
			cur.TrueTarget = parseTrueTarget(s.Text())

			s.Scan()
			cur.FalseTarget = parseFalseTarget(s.Text())

			//fmt.Println(cur)
			res = append(res, &cur)
		}
	}

	return res
}

func parseItems(str string) []int {
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "Starting items: ")
	tokens := strings.Split(str, ", ")

	var res []int
	for _, token := range tokens {
		item, _ := strconv.ParseInt(token, 10, 32)
		res = append(res, int(item))
	}

	return res
}

func parseOperation(str string) Operation {
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "Operation: new = ")
	tokens := strings.Split(str, " ")

	first := parseOperand(tokens[0])
	operator := OperatorType(tokens[1])
	second := parseOperand(tokens[2])

	return Operation{
		Type:   operator,
		First:  first,
		Second: second,
	}
}

func parseOperand(str string) Operand {
	if str == "old" {
		return Operand{
			Type: OperandTypeSelf,
		}
	}

	value, _ := strconv.ParseInt(str, 10, 32)

	return Operand{
		Type:  OperandTypeConst,
		Value: int(value),
	}
}

func parseTest(str string) Test {
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "Test: ")

	var test Test
	if strings.HasPrefix(str, string(TestTypeDivBy)) {
		str = strings.TrimPrefix(str, string(TestTypeDivBy))
		value, _ := strconv.ParseInt(str, 10, 32)

		test.Type = TestTypeDivBy
		test.Value = int(value)
	}

	return test
}

func parseTrueTarget(str string) int {
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "If true: throw to monkey ")

	target, _ := strconv.ParseInt(str, 10, 32)

	return int(target)
}

func parseFalseTarget(str string) int {
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "If false: throw to monkey ")

	target, _ := strconv.ParseInt(str, 10, 32)

	return int(target)
}

func solveFirst(monkeys []*Monkey) int {
	roundCount := 20

	inspectCount := make([]int, len(monkeys))
	for r := 0; r < roundCount; r++ {
		for i, m := range monkeys {
			for _, item := range m.Items {
				newItem := eval(item, m.Operation) / 3

				var target int
				if test(newItem, m.Test) {
					// True
					target = m.TrueTarget
				} else {
					// False
					target = m.FalseTarget
				}

				monkeys[target].Items = append(monkeys[target].Items, newItem)
			}

			inspectCount[i] += len(m.Items)
			m.Items = nil
		}

		//fmt.Println(inspectCount)
	}

	sort.Slice(inspectCount, func(i, j int) bool {
		return inspectCount[i] > inspectCount[j]
	})

	return inspectCount[0] * inspectCount[1]
}

func solveSecond(monkeys []*Monkey) int {
	roundCount := 10000

	// Evaluate mod factor (LCM of all mod values)
	mod := monkeys[0].Test.Value
	for _, m := range monkeys {
		mod = lcm(mod, m.Test.Value)
	}

	inspectCount := make([]int, len(monkeys))
	for r := 0; r < roundCount; r++ {
		for i, m := range monkeys {
			for _, item := range m.Items {
				newItem := eval(item, m.Operation)

				var target int
				if test(newItem, m.Test) {
					// True
					target = m.TrueTarget
				} else {
					// False
					target = m.FalseTarget
				}

				newItem %= mod
				monkeys[target].Items = append(monkeys[target].Items, newItem)
			}

			inspectCount[i] += len(m.Items)
			m.Items = nil
		}

		//fmt.Println(inspectCount)
	}

	sort.Slice(inspectCount, func(i, j int) bool {
		return inspectCount[i] > inspectCount[j]
	})

	return inspectCount[0] * inspectCount[1]
}

func solveSecondAlt(monkeys []*Monkey) int {
	roundCount := 10000

	type ItemByModulo map[int]int

	monkeyItems := make([][]ItemByModulo, len(monkeys))
	for i, m := range monkeys {
		monkeyItems[i] = make([]ItemByModulo, len(m.Items))
		for j, item := range m.Items {
			monkeyItems[i][j] = make(map[int]int)
			for _, n := range monkeys {
				monkeyItems[i][j][n.Test.Value] = item % n.Test.Value
			}
		}
	}

	// O(round * items * monkeys)
	inspectCount := make([]int, len(monkeys))
	for r := 0; r < roundCount; r++ {
		for i, m := range monkeys {
			for _, itemByModulo := range monkeyItems[i] {
				for mod, item := range itemByModulo {
					itemByModulo[mod] = eval(item, m.Operation) % mod
				}

				var target int
				if test(itemByModulo[m.Test.Value], m.Test) {
					// True
					target = m.TrueTarget
				} else {
					// False
					target = m.FalseTarget
				}

				monkeyItems[target] = append(monkeyItems[target], itemByModulo)
			}

			inspectCount[i] += len(monkeyItems[i])
			monkeyItems[i] = nil
		}

		//fmt.Println(inspectCount)
	}

	sort.Slice(inspectCount, func(i, j int) bool {
		return inspectCount[i] > inspectCount[j]
	})

	return inspectCount[0] * inspectCount[1]
}

func eval(value int, operation Operation) int {
	var first int
	if operation.First.Type == OperandTypeSelf {
		first = value
	} else {
		first = operation.First.Value
	}

	var second int
	if operation.Second.Type == OperandTypeSelf {
		second = value
	} else {
		second = operation.Second.Value
	}

	var res int
	switch operation.Type {
	case OperatorTypeAdd:
		res = first + second
	case OperatorTypeMult:
		res = first * second
	}

	return res
}

func test(value int, test Test) bool {
	var res bool
	switch test.Type {
	case TestTypeDivBy:
		res = (value % test.Value) == 0
	}

	return res
}

func gcd(a int, b int) int {
	if b == 0 {
		return a
	}

	return gcd(b, a%b)
}

func lcm(a int, b int) int {
	return a * b / gcd(a, b)
}
