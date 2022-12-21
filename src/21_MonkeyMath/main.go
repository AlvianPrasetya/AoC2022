package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	OperatorTypeAdd  = OperatorType("+")
	OperatorTypeSub  = OperatorType("-")
	OperatorTypeMult = OperatorType("*")
	OperatorTypeDiv  = OperatorType("/")
	OperatorTypeEq   = OperatorType("=")

	OperandTypeInput = OperandType("input")
	OperandTypeVar   = OperandType("var")
	OperandTypeConst = OperandType("const")
)

type Expr interface {
	GetType() string
}

type OperatorType string

type Operation struct {
	Type   OperatorType
	First  Expr
	Second Expr
}

func (op Operation) GetType() string {
	return string(op.Type)
}

func (op Operation) String() string {
	return fmt.Sprintf("(%s %s %s)", op.First, op.Type, op.Second)
}

type OperandType string

type Operand struct {
	Type OperandType
	Val  interface{}
}

func (op Operand) GetType() string {
	return string(op.Type)
}

func (op Operand) String() string {
	switch op.Type {
	case OperandTypeInput:
		return fmt.Sprintf(op.Val.(string))
	case OperandTypeVar:
		return string(op.Val.(string))
	case OperandTypeConst:
		return strconv.FormatInt(op.Val.(int64), 10)
	}
	return ""
}

func main() {
	instrs := parseInput("in.txt")
	// fmt.Println(exprs)
	fmt.Println(solveFirst(instrs))
	fmt.Println(solveSecond(instrs))
}

func parseInput(in string) map[string]Expr {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	res := make(map[string]Expr)
	for s.Scan() {
		tokens := strings.Split(s.Text(), ": ")
		exprName := tokens[0]
		exprTokens := strings.Split(tokens[1], " ")
		if len(exprTokens) == 1 {
			// Const
			val, _ := strconv.ParseInt(exprTokens[0], 10, 64)
			res[exprName] = Operand{
				Type: OperandTypeConst,
				Val:  val,
			}
		} else {
			// Operation
			res[exprName] = Operation{
				Type: OperatorType(exprTokens[1]),
				First: Operand{
					Type: OperandTypeVar,
					Val:  exprTokens[0],
				},
				Second: Operand{
					Type: OperandTypeVar,
					Val:  exprTokens[2],
				},
			}
		}
	}

	return res
}

func solveFirst(instrs map[string]Expr) int64 {
	return eval(instrs, "root", make(map[string]Expr)).(Operand).Val.(int64)
}

func solveSecond(exprs map[string]Expr) int64 {
	exprs["root"] = Operation{
		Type:   OperatorTypeEq,
		First:  exprs["root"].(Operation).First,
		Second: exprs["root"].(Operation).Second,
	}
	exprs["humn"] = Operand{
		Type: OperandTypeInput,
		Val:  "input",
	}

	rootExpr := eval(exprs, "root", make(map[string]Expr))
	// fmt.Println(rootExpr)

	toEvalName := exprs["root"].(Operation).First.(Operand).Val.(string)
	// leftExpr := rootExpr.(Operation).First
	rightVal := rootExpr.(Operation).Second.(Operand).Val.(int64)
	// fmt.Println(toEvalName, rightVal)

	l := int64(math.MinInt64)
	r := int64(math.MaxInt64)
	for l <= r {
		mid := (l + r) / 2

		exprs["humn"] = Operand{
			Type: OperandTypeConst,
			Val:  mid,
		}
		evalRes := eval(exprs, toEvalName, make(map[string]Expr)).(Operand).Val.(int64)
		// fmt.Println(l, r, evalRes, rightVal)

		if evalRes < rightVal {
			r = mid - 1
		} else if evalRes > rightVal {
			l = mid + 1
		} else {
			return mid
		}
	}

	return 0
}

func eval(instrs map[string]Expr, exprName string, memo map[string]Expr) Expr {
	if expr, ok := memo[exprName]; ok {
		return expr
	}

	expr := instrs[exprName]

	var res Expr
	switch expr.(type) {
	case Operation:
		op := expr.(Operation)
		firstOperand := op.First.(Operand)
		secondOperand := op.Second.(Operand)

		firstExpr := eval(instrs, firstOperand.Val.(string), memo)
		secondExpr := eval(instrs, secondOperand.Val.(string), memo)

		res = evalOperation(op.Type, firstExpr, secondExpr)
	case Operand:
		res = expr.(Operand)
	}

	memo[exprName] = res
	return res
}

func evalOperation(op OperatorType, firstExpr Expr, secondExpr Expr) Expr {
	if firstExpr.GetType() == string(OperandTypeConst) && secondExpr.GetType() == string(OperandTypeConst) {
		first := firstExpr.(Operand).Val.(int64)
		second := secondExpr.(Operand).Val.(int64)

		var val int64
		switch op {
		case OperatorTypeAdd:
			val = first + second
		case OperatorTypeSub:
			val = first - second
		case OperatorTypeMult:
			val = first * second
		case OperatorTypeDiv:
			val = first / second
		case OperatorTypeEq:
			if first == second {
				val = 1
			}
		}

		return Operand{
			Type: OperandTypeConst,
			Val:  val,
		}
	}
	return Operation{
		Type:   op,
		First:  firstExpr,
		Second: secondExpr,
	}
}
