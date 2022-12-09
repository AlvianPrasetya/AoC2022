package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	OpInp = OperatorType("inp")
	OpAdd = OperatorType("add")
	OpMul = OperatorType("mul")
	OpDiv = OperatorType("div")
	OpMod = OperatorType("mod")
	OpEql = OperatorType("eql")

	OpInput = OperandType("input")
	OpVar   = OperandType("var")
	OpConst = OperandType("const")

	MinDigit = 1
	MaxDigit = 9
)

var (
	inputLength int
	exprID      int
	evalCount   int32
)

type Instruction struct {
	Type   OperatorType
	First  *Operand
	Second *Operand
}

func (i *Instruction) String() string {
	if i.Second == nil {
		return fmt.Sprintf("(%s %s)", i.Type, i.First)
	}

	return fmt.Sprintf("(%s %s %s)", i.First, i.Type, i.Second)
}

type OperatorType string

func (o OperatorType) String() string {
	switch o {
	case OpInp:
		return "inp"
	case OpAdd:
		return "+"
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpMod:
		return "%"
	case OpEql:
		return "=="
	}

	return ""
}

type Operand struct {
	Type  OperandType
	Value interface{} // int for Type const, string for Type var
}

type OperandType string

func (o *Operand) GetType() string {
	return string(o.Type)
}

func (o *Operand) GetValue() int {
	if o.Type != OpConst {
		panic(fmt.Sprintf("invalid GetValue of type %s", o.Type))
	}

	return o.Value.(int)
}

func (o *Operand) GetOperatorCount() int {
	return 0
}

func (o *Operand) String() string {
	switch o.Type {
	case OpInput:
		return fmt.Sprintf("in%d", o.Value.(int))
	case OpVar:
		return string(o.Value.(byte))
	case OpConst:
		return fmt.Sprintf("%d", o.Value.(int))
	}

	return ""
}

type Operation struct {
	ID     int
	Type   OperatorType
	First  Expr
	Second Expr
}

func (o *Operation) GetType() string {
	return string(o.Type)
}

func (o *Operation) GetValue() int {
	panic("invalid GetValue of type Operation")
}

func (o *Operation) GetOperatorCount() int {
	res := 1
	res += o.First.GetOperatorCount()
	res += o.Second.GetOperatorCount()

	return res
}

func (o *Operation) String() string {
	return fmt.Sprintf("(%s %s %s)", o.First, o.Type, o.Second)
}

// Expr consists of Operand and Operation
type Expr interface {
	GetType() string
	GetValue() int
	GetOperatorCount() int
}

func main() {
	input := parseInput("in.txt")
	//fmt.Println(input)
	fmt.Println(solveFirst(input))
}

func parseInput(in string) []*Instruction {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []*Instruction
	for s.Scan() {
		str := s.Text()
		tokens := strings.Split(str, " ")

		var cur Instruction
		cur.Type = OperatorType(tokens[0])

		if cur.Type == OpInp {
			inputLength++
		}

		// Parse first argument
		constVal, err := strconv.ParseInt(tokens[1], 10, 32)
		if err == nil {
			// Constant
			cur.First = &Operand{
				Type:  OpConst,
				Value: int(constVal),
			}
		} else {
			// Variable
			cur.First = &Operand{
				Type:  OpVar,
				Value: tokens[1][0],
			}
		}

		if len(tokens) == 3 {
			// Parse second argument
			constVal, err = strconv.ParseInt(tokens[2], 10, 32)
			if err == nil {
				// Constant
				cur.Second = &Operand{
					Type:  OpConst,
					Value: int(constVal),
				}
			} else {
				// Variable
				cur.Second = &Operand{
					Type:  OpVar,
					Value: tokens[2][0],
				}
			}
		}

		res = append(res, &cur)
	}

	return res
}

func solveFirst(input []*Instruction) string {
	exprMap := make(map[byte]Expr)
	for i := 0; i < 4; i++ {
		exprMap[byte('w'+i)] = &Operand{
			Type:  OpConst,
			Value: 0,
		}
	}

	var inputIdx int
	for _, in := range input {
		switch in.Type {
		case OpInp:
			exprMap[in.First.Value.(byte)] = &Operand{
				Type:  OpInput,
				Value: inputIdx,
			}
			inputIdx++
		case OpAdd, OpMul, OpDiv, OpMod, OpEql:
			evalInstr(exprMap, in.Type, in.First, in.Second)
		}

		//fmt.Printf("input=%s w=%s x=%s y=%s z=%s\n", in, exprMap["w"], exprMap["x"], exprMap["y"], exprMap["z"])
	}

	/*resCh := make(chan []int, 81)
	var wg sync.WaitGroup
	for i := MinDigit; i <= MaxDigit; i++ {
		for j := MinDigit; j <= MaxDigit; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				resCh <- searchMaxValidInput(expr, []int{i, j})
			}(i, j)
		}
	}
	wg.Wait()
	close(resCh)

	var finalRes []int
	for res := range resCh {
		if isArrLess(finalRes, res) {
			finalRes = res
		}
	}*/
	res := searchMaxValidInput(&Operation{
		Type:  OpEql,
		First: exprMap['z'],
		Second: &Operand{
			Type:  OpConst,
			Value: 0,
		},
	}, nil)

	return toString(res)
}

func evalInstr(exprMap map[byte]Expr, opType OperatorType, first *Operand, second *Operand) {
	// First operand is always a var
	firstExpr := exprMap[first.Value.(byte)]

	var secondExpr Expr
	switch second.Type {
	case OpInput:
		secondExpr = second
	case OpVar:
		secondExpr = exprMap[second.Value.(byte)]
	case OpConst:
		secondExpr = second
	}

	if res, ok := evalSpecial(opType, firstExpr, secondExpr); ok {
		// There is a special optimizations possible with this operator, e.g.: multiply by 0
		exprMap[first.Value.(byte)] = res
	} else if firstExpr.GetType() == string(OpConst) && secondExpr.GetType() == string(OpConst) {
		// Both operands are constants, directly evaluate
		exprMap[first.Value.(byte)] = evalOp(opType, firstExpr.(*Operand), secondExpr.(*Operand))
	} else {
		exprMap[first.Value.(byte)] = &Operation{
			ID:     exprID,
			Type:   opType,
			First:  firstExpr,
			Second: secondExpr,
		}
		exprID++
	}
}

func evalSpecial(opType OperatorType, first Expr, second Expr) (Expr, bool) {
	switch opType {
	case OpAdd:
		if first.GetType() == string(OpConst) && first.GetValue() == 0 {
			// Addition with 0
			return second, true
		}
		if second.GetType() == string(OpConst) && second.GetValue() == 0 {
			// Addition with 0
			return first, true
		}
	case OpMul:
		if (first.GetType() == string(OpConst) && first.GetValue() == 0) ||
			(second.GetType() == string(OpConst) && second.GetValue() == 0) {
			// Multiplication by 0
			return &Operand{
				Type:  OpConst,
				Value: 0,
			}, true
		}
	case OpDiv:
		if first.GetType() == string(OpConst) && first.GetValue() == 0 {
			// Division from 0
			return first, true
		}
		if second.GetType() == string(OpConst) && second.GetValue() == 1 {
			// Division by 1
			return first, true
		}
	}

	return nil, false
}

func evalOp(op OperatorType, first *Operand, second *Operand) *Operand {
	if first.Type != OpConst || second.Type != OpConst {
		panic(fmt.Sprintf("invalid eval of type %s and %s", first.Type, second.Type))
	}

	firstVal := first.GetValue()
	secondVal := second.GetValue()

	var val int
	switch op {
	case OpAdd:
		val = firstVal + secondVal
	case OpMul:
		val = firstVal * secondVal
	case OpDiv:
		val = firstVal / secondVal
	case OpMod:
		val = firstVal % secondVal
	case OpEql:
		if firstVal == secondVal {
			val = 1
		}
	}

	return &Operand{
		Type:  OpConst,
		Value: val,
	}
}

func searchMaxValidInput(expr Expr, input []int) []int {
	if len(input) == inputLength {
		newCount := atomic.AddInt32(&evalCount, 1)
		if newCount%1000 == 0 {
			fmt.Printf("Eval count: %d\n", newCount)
		}

		if evalExpr(expr, input, make(map[int]int), &sync.RWMutex{}) == 1 {
			fmt.Printf("Valid eval %v\n", input)
			return input
		}
		return nil
	}

	var res []int
	for i := MinDigit; i <= MaxDigit; i++ {
		cur := searchMaxValidInput(expr, append(input, i))
		if isArrLess(res, cur) {
			res = cur
		}
	}

	return res
}

func evalExpr(expr Expr, input []int, exprMemo map[int]int, mutex *sync.RWMutex) int {
	if _, ok := expr.(*Operand); ok {
		// Operand
		operand := expr.(*Operand)

		switch operand.Type {
		case OpInput:
			return input[operand.Value.(int)]
		case OpConst:
			return operand.Value.(int)
		}
		// OpVar should not appear in the final expression
		panic(fmt.Sprintf("invalid operand type %s", operand.Type))
	}

	// Operation
	operation := expr.(*Operation)

	// Check if result is in memo
	mutex.RLock()
	res, ok := exprMemo[operation.ID]
	mutex.RUnlock()
	if ok {
		return res
	}

	first := &Operand{
		Type:  OpConst,
		Value: evalExpr(operation.First, input, exprMemo, mutex),
	}
	second := &Operand{
		Type:  OpConst,
		Value: evalExpr(operation.Second, input, exprMemo, mutex),
	}
	res = evalOp(operation.Type, first, second).Value.(int)

	// Put result in memo
	mutex.Lock()
	exprMemo[operation.ID] = res
	mutex.Unlock()

	return res
}

func isArrLess(first []int, second []int) bool {
	if len(second) == 0 {
		return false
	}
	if len(first) == 0 {
		return true
	}

	if len(first) != len(second) {
		panic(fmt.Sprintf("mismatched array length %d and %d", len(first), len(second)))
	}

	for i := 0; i < len(first); i++ {
		if first[i] < second[i] {
			return true
		} else if first[i] > second[i] {
			return false
		}
	}

	return false
}

func toString(input []int) string {
	var sb strings.Builder
	for _, digit := range input {
		sb.WriteByte(byte('0' + digit))
	}
	return sb.String()
}
