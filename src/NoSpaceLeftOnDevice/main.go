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
	InputPrefix = "$ "

	InputTypeCd = InputType("cd")
	InputTypeLs = InputType("ls")

	OutputTypeDir  = OutputType("dir")
	OutputTypeFile = OutputType("file")
)

type Dialogue struct {
	Input   Input
	Outputs []Output // lines of output
}

func (d *Dialogue) String() string {
	return fmt.Sprintf("Input: %v, Outputs: %v\n", d.Input, d.Outputs)
}

type InputType string

type Input struct {
	Type  InputType
	Param string
}

func (i *Input) String() string {
	return fmt.Sprintf("%s %s", i.Type, i.Param)
}

type OutputType string

type Output interface {
	GetType() OutputType
	GetName() string
}

type OutputDir struct {
	Name string
}

func (d *OutputDir) GetType() OutputType {
	return OutputTypeDir
}

func (d *OutputDir) GetName() string {
	return d.Name
}

func (d *OutputDir) String() string {
	return fmt.Sprintf("%s %s", d.GetType(), d.Name)
}

type OutputFile struct {
	Size int
	Name string
}

func (f *OutputFile) GetType() OutputType {
	return OutputTypeFile
}

func (f *OutputFile) GetName() string {
	return f.Name
}

func (f *OutputFile) String() string {
	return fmt.Sprintf("%d %s\n", f.Size, f.Name)
}

func main() {
	input := parseInput("in.txt")
	fmt.Println(solveFirst(input))
	fmt.Println(solveSecond(input))
}

func parseInput(in string) []*Dialogue {
	f, _ := os.Open(in)
	defer f.Close()

	s := bufio.NewScanner(f)

	var res []*Dialogue
	var cur *Dialogue
	for s.Scan() {
		str := s.Text()

		if strings.HasPrefix(str, InputPrefix) {
			// Input
			if cur != nil {
				// Flush current Dialogue
				res = append(res, cur)
				cur = nil
			}

			str = strings.TrimPrefix(str, InputPrefix)
			tokens := strings.Split(str, " ")

			switch tokens[0] {
			case string(InputTypeCd):
				cur = &Dialogue{
					Input: Input{
						Type:  InputTypeCd,
						Param: tokens[1],
					},
				}
			case string(InputTypeLs):
				cur = &Dialogue{
					Input: Input{
						Type: InputTypeLs,
					},
				}
			}
		} else {
			// Output
			tokens := strings.Split(str, " ")

			switch tokens[0] {
			case string(OutputTypeDir):
				cur.Outputs = append(cur.Outputs, &OutputDir{
					Name: tokens[1],
				})
			default:
				size, _ := strconv.ParseInt(tokens[0], 10, 32)
				cur.Outputs = append(cur.Outputs, &OutputFile{
					Size: int(size),
					Name: tokens[1],
				})
			}
		}
	}

	if cur != nil {
		// Flush current Dialogue
		res = append(res, cur)
		cur = nil
	}

	return res
}

// Total size of directories with size <= 100,000
func solveFirst(input []*Dialogue) int {
	dirSizes := evalDirSizes(input)

	var res int
	for _, size := range dirSizes {
		if size <= 100000 {
			res += size
		}
	}

	return res
}

// Size of smallest dir to delete such that total unused space >= 30000000
func solveSecond(input []*Dialogue) int {
	dirSizes := evalDirSizes(input)
	spaceToFree := dirSizes["."] - 40000000

	res := math.MaxInt
	for _, size := range dirSizes {
		if size >= spaceToFree && size < res {
			res = size
		}
	}

	return res
}

func evalDirSizes(input []*Dialogue) map[string]int {
	dirStack := []string{"."}
	dirSizes := make(map[string]int)
	listedDirs := make(map[string]bool)
	for _, d := range input {
		switch d.Input.Type {
		case InputTypeCd:
			switch d.Input.Param {
			case "/":
				// Go back to root
				dirStack = []string{"."}
			case "..":
				// Go up 1 level
				dirStack = dirStack[:len(dirStack)-1]
			default:
				// Go into dir
				dirStack = append(dirStack, d.Input.Param)
			}
		case InputTypeLs:
			if listedDirs[getDirPath(dirStack)] {
				// Prevent double counting
				continue
			}

			// Evaluate direct size (files immediately under this dir)
			var directSize int
			for _, o := range d.Outputs {
				if o.GetType() == OutputTypeFile {
					directSize += o.(*OutputFile).Size
				}
			}

			for i := range dirStack {
				dirSizes[getDirPath(dirStack[:i+1])] += directSize
			}
		}
	}

	return dirSizes
}

func getDirPath(dirStack []string) string {
	var res string
	for i, dir := range dirStack {
		if i != 0 {
			res += "/"
		}
		res += dir
	}

	return res
}
