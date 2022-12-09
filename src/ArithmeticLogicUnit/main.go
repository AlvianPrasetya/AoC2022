package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
)

type Instruction struct {
	Type     OperatorType
	Operands []Operand
}

type OperatorType string

type Operand struct {
	Type  OperandType
	Value interface{} // int for Type const, string for Type var
}

type OperandType string

func (o Operand) String() string {
	switch o.Type {
	case OpVar:
		return o.Value.(string)
	case OpConst:
		return fmt.Sprintf("%d", o.Value.(int))
	}

	return ""
}

type Expr interface {
	String() string
}

type Operation struct {
	Type     OperatorType
	Operands []Expr
}

func (o Operation) String() string {
	var res string
	for i, op := range o.Operands {
		if i != 0 {
			res += ", "
		}
		res += op.String()
	}

	res = fmt.Sprintf("%s(%s)", o.Type, res)
	return res
}

func main() {
	input := parseInput("in.txt")
	//fmt.Println(input)
	fmt.Println(solveFirst(input))
}

func parseInput(in string) []Instruction {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []Instruction
	for s.Scan() {
		str := s.Text()
		tokens := strings.Split(str, " ")

		var cur Instruction
		cur.Type = OperatorType(tokens[0])

		for i := 1; i < len(tokens); i++ {
			constVal, err := strconv.ParseInt(tokens[i], 10, 32)
			if err == nil {
				// Constant
				cur.Operands = append(cur.Operands, Operand{
					Type:  OpConst,
					Value: int(constVal),
				})
			} else {
				// Variable
				cur.Operands = append(cur.Operands, Operand{
					Type:  OpVar,
					Value: tokens[i],
				})
			}
		}

		res = append(res, cur)
	}

	return res
}

func solveFirst(input []Instruction) Expr {
	exprMap := make(map[string]Expr)
	for i := 0; i < 4; i++ {
		exprMap[string(byte('w')+byte(i))] = Operand{
			Type:  OpConst,
			Value: 0,
		}
	}

	var inputIdx int
	for i, in := range input {
		if i == 100 {
			break
		}
		switch in.Type {
		case OpInp:
			exprMap[in.Operands[0].Value.(string)] = Operand{
				Type:  OpVar,
				Value: fmt.Sprintf("in[%d]", inputIdx),
			}
			inputIdx++
		case OpAdd:
			first := exprMap[in.Operands[0].Value.(string)]

			var second Expr
			switch in.Operands[1].Type {
			case OpVar:
				second = exprMap[in.Operands[1].Value.(string)]
			case OpConst:
				second = in.Operands[1]
			}
			exprMap[in.Operands[0].Value.(string)] = Operation{
				Type:     OpAdd,
				Operands: []Expr{first, second},
			}
		case OpMul:
			first := exprMap[in.Operands[0].Value.(string)]

			var second Expr
			switch in.Operands[1].Type {
			case OpVar:
				second = exprMap[in.Operands[1].Value.(string)]
			case OpConst:
				second = in.Operands[1]
			}
			exprMap[in.Operands[0].Value.(string)] = Operation{
				Type:     OpMul,
				Operands: []Expr{first, second},
			}
		case OpDiv:
			first := exprMap[in.Operands[0].Value.(string)]

			var second Expr
			switch in.Operands[1].Type {
			case OpVar:
				second = exprMap[in.Operands[1].Value.(string)]
			case OpConst:
				second = in.Operands[1]
			}
			exprMap[in.Operands[0].Value.(string)] = Operation{
				Type:     OpDiv,
				Operands: []Expr{first, second},
			}
		case OpMod:
			first := exprMap[in.Operands[0].Value.(string)]

			var second Expr
			switch in.Operands[1].Type {
			case OpVar:
				second = exprMap[in.Operands[1].Value.(string)]
			case OpConst:
				second = in.Operands[1]
			}
			exprMap[in.Operands[0].Value.(string)] = Operation{
				Type:     OpMod,
				Operands: []Expr{first, second},
			}
		case OpEql:
			first := exprMap[in.Operands[0].Value.(string)]

			var second Expr
			switch in.Operands[1].Type {
			case OpVar:
				second = exprMap[in.Operands[1].Value.(string)]
			case OpConst:
				second = in.Operands[1]
			}
			exprMap[in.Operands[0].Value.(string)] = Operation{
				Type:     OpEql,
				Operands: []Expr{first, second},
			}
		}
	}

	return exprMap["z"]
}
