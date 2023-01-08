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
	intcode := parse(in)
	go intcode.run()
	intcode.in <- 1
	println(intcode.waitHalt())
}

func task2(in io.Reader) {
	intcode := parse(in)
	go intcode.run()
	intcode.in <- 2
	println(intcode.waitHalt())
}

type Intcode struct {
	program   []int
	index     int
	in        chan int
	out       chan int
	lastout   int
	reloffset int
	halt      chan int
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
	AdjustRelativeOffset
	Halt Opcode = 99
)

type ParamMode int

const (
	Position ParamMode = iota
	Immediate
	Relative
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

	case AdjustRelativeOffset:
		param1 := i.loadParam(pmode1, i.index+1)
		i.reloffset += param1
		i.index += 2
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
	idx := index

	switch mode {
	case Position:
		idx = i.program[index]
	case Relative:
		idx = i.program[index] + i.reloffset
	}

	if idx >= len(i.program) {
		return 0
	}

	return i.program[idx]
}

func (i *Intcode) storeParam(mode ParamMode, index, value int) {
	idx := index
	switch mode {
	case Position:
		idx = i.program[index]
	case Relative:
		idx = i.program[index] + i.reloffset
	}

	if idx >= len(i.program) {
		i.program = append(i.program, make([]int, idx+1-len(i.program))...)
	}

	i.program[idx] = value
}

func parse(in io.Reader) *Intcode {
	scanner := bufio.NewScanner(in)
	scanner.Scan()

	strcode := strings.Split(scanner.Text(), ",")
	program := make([]int, len(strcode))
	for i, s := range strcode {
		program[i] = aoc.StrToInt(s)
	}

	return &Intcode{
		program: program,
		in:      make(chan int, 1),
		out:     make(chan int, 1),
		halt:    make(chan int, 1),
	}
}
