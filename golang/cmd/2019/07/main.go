package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	program := parse(in)
	var max int
	for _, phases := range permutations([]int{0, 1, 2, 3, 4}) {
		if out := run(program, phases, false); out > max {
			max = out
		}
	}
	println(max)
}

func task2(in io.Reader) {
	program := parse(in)
	var max int
	for _, phases := range permutations([]int{5, 6, 7, 8, 9}) {
		if out := run(program, phases, true); out > max {
			max = out
		}
	}
	println(max)
}

func run(program []int, phases []int, feedback bool) int {
	ampA, ampB, ampC, ampD, ampE := makeAmp(program), makeAmp(program), makeAmp(program), makeAmp(program), makeAmp(program)

	go ampA.run()
	go ampB.run()
	go ampC.run()
	go ampD.run()
	go ampE.run()

	ampA.in <- phases[0]
	ampA.in <- 0

	if !feedback {
		ampB.in <- phases[1]
		ampB.in <- ampA.waitHalt()
		ampC.in <- phases[2]
		ampC.in <- ampB.waitHalt()
		ampD.in <- phases[3]
		ampD.in <- ampC.waitHalt()
		ampE.in <- phases[4]
		ampE.in <- ampD.waitHalt()
		return ampE.waitHalt()
	}

	ampB.in <- phases[1]
	ampB.in <- <-ampA.out
	ampC.in <- phases[2]
	ampC.in <- <-ampB.out
	ampD.in <- phases[3]
	ampD.in <- <-ampC.out
	ampE.in <- phases[4]
	ampE.in <- <-ampD.out

	for {
		select {
		case ampA.in <- <-ampE.out:
			ampB.in <- <-ampA.out
			ampC.in <- <-ampB.out
			ampD.in <- <-ampC.out
			ampE.in <- <-ampD.out
		case out := <-ampE.halt:
			return out
		}
	}
}

func permutations(nums []int) [][]int {
	var res [][]int
	permute(nums, 0, &res)
	return res
}

func permute(nums []int, start int, res *[][]int) {
	if start == len(nums) {
		*res = append(*res, append([]int{}, nums...))
		return
	}

	for i := start; i < len(nums); i++ {
		nums[start], nums[i] = nums[i], nums[start]
		permute(nums, start+1, res)
		nums[start], nums[i] = nums[i], nums[start]
	}
}

type Intcode struct {
	program []int
	index   int
	in      chan int
	out     chan int
	lastout int
	halt    chan int
}

type Opcode int

const (
	Add Opcode = iota + 1
	Multiply
	Input
	Output
	JumpIfTrue
	JumpIfFalse
	LessThan
	Equals
	Halt Opcode = 99
)

type ParamMode int

const (
	Position ParamMode = iota
	Immediate
)

func (i *Intcode) waitHalt() int {
	for {
		select {
		case <-i.out:
		case o := <-i.halt:
			return o
		}
	}
}

func (i *Intcode) run() {
	for {
		op, pmode1, pmode2, pmode3 := i.loadOpcode()
		if op == Halt {
			close(i.out)
			i.halt <- i.lastout
			close(i.halt)
			return
		}

		i.step(op, pmode1, pmode2, pmode3)
	}
}

func (i *Intcode) step(op Opcode, pmode1, pmode2, pmode3 ParamMode) {
	switch op {
	case Add:
		param1 := i.loadParam(pmode1, i.index+1)
		param2 := i.loadParam(pmode2, i.index+2)
		i.storeParam(pmode3, i.index+3, param1+param2)
		i.index += 4

	case Multiply:
		param1 := i.loadParam(pmode1, i.index+1)
		param2 := i.loadParam(pmode2, i.index+2)
		i.storeParam(pmode3, i.index+3, param1*param2)
		i.index += 4

	case Input:
		i.storeParam(pmode1, i.index+1, <-i.in)
		i.index += 2

	case Output:
		param1 := i.loadParam(pmode1, i.index+1)
		i.lastout = param1
		i.out <- param1
		i.index += 2

	case JumpIfTrue:
		param1 := i.loadParam(pmode1, i.index+1)
		param2 := i.loadParam(pmode2, i.index+2)
		if param1 != 0 {
			i.index = param2
		} else {
			i.index += 3
		}

	case JumpIfFalse:
		param1 := i.loadParam(pmode1, i.index+1)
		param2 := i.loadParam(pmode2, i.index+2)
		if param1 == 0 {
			i.index = param2
		} else {
			i.index += 3
		}

	case LessThan:
		param1 := i.loadParam(pmode1, i.index+1)
		param2 := i.loadParam(pmode2, i.index+2)
		if param1 < param2 {
			i.storeParam(pmode3, i.index+3, 1)
		} else {
			i.storeParam(pmode3, i.index+3, 0)
		}
		i.index += 4

	case Equals:
		param1 := i.loadParam(pmode1, i.index+1)
		param2 := i.loadParam(pmode2, i.index+2)
		if param1 == param2 {
			i.storeParam(pmode3, i.index+3, 1)
		} else {
			i.storeParam(pmode3, i.index+3, 0)
		}
		i.index += 4
	}
}

func (i *Intcode) loadOpcode() (op Opcode, pmode1, pmode2, pmode3 ParamMode) {
	op = Opcode(i.program[i.index] % 100)
	pmode1 = ParamMode(i.program[i.index] / 100 % 10)
	pmode2 = ParamMode(i.program[i.index] / 1000 % 10)
	pmode3 = ParamMode(i.program[i.index] / 10000 % 10)
	return
}

func (i *Intcode) loadParam(mode ParamMode, index int) int {
	switch mode {
	case Position:
		return i.program[i.program[index]]
	case Immediate:
		return i.program[index]
	}
	return 0
}

func (i *Intcode) storeParam(mode ParamMode, index, value int) {
	switch mode {
	case Position:
		i.program[i.program[index]] = value
	case Immediate:
		i.program[index] = value
	}
}

func parse(in io.Reader) []int {
	scanner := bufio.NewScanner(in)
	scanner.Scan()

	strcode := strings.Split(scanner.Text(), ",")
	program := make([]int, len(strcode))
	for i, s := range strcode {
		program[i] = aoc.StrToInt(s)
	}

	return program
}

func makeAmp(program []int) *Intcode {
	programCopy := make([]int, len(program))
	copy(programCopy, program)

	amp := &Intcode{
		program: programCopy,
		in:      make(chan int, 1),
		out:     make(chan int, 1),
		halt:    make(chan int, 1),
	}

	return amp
}
