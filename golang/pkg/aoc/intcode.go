package aoc

import (
	"bufio"
	"io"
	"strings"
)

type Intcode struct {
	In  chan int
	Out chan int

	program   map[int]int
	index     int
	reloffset int
	lastout   int
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

func (i *Intcode) SetProgram(index int, value int) {
	i.program[index] = value
}

func (i *Intcode) WaitHalt() int {
	for {
		select {
		case <-i.Out:
		case o := <-i.halt:
			return o
		}
	}
}

func (i *Intcode) Run() {
	for {
		op, pmode1, pmode2, pmode3 := i.loadOpcode()
		if op == Halt {
			close(i.Out)
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
		i.storeParam(pmode1, i.index+1, <-i.In)
		i.index += 2

	case Output:
		param1 := i.loadParam(pmode1, i.index+1)
		i.lastout = param1
		i.Out <- param1
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

	default:
		panic("unknown opcode")
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

	i.program[idx] = value
}

func ParseIntcode(in io.Reader) *Intcode {
	scanner := bufio.NewScanner(in)
	scanner.Scan()

	strcode := strings.Split(scanner.Text(), ",")
	program := make(map[int]int, len(strcode))
	for i, s := range strcode {
		program[i] = StrToInt(s)
	}

	return &Intcode{
		program: program,
		In:      make(chan int, 1),
		Out:     make(chan int, 1),
		halt:    make(chan int, 1),
	}
}
