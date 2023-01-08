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
	solve(in, 1)
}

func task2(in io.Reader) {
	solve(in, 5)
}

func solve(in io.Reader, id int) {
	intcode := parse(in)
	go intcode.run()
	intcode.in <- id
	<-intcode.done
	println("halted")
}

type Intcode struct {
	program []int
	index   int
	in      chan int
	out     chan int
	done    chan struct{}
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

func (i *Intcode) run() {
	go func() {
		for {
			out, ok := <-i.out
			if !ok {
				i.done <- struct{}{}
				close(i.done)
				return
			}

			println(out)
		}
	}()

	for {
		op, pmode1, pmode2, pmode3 := i.loadOpcode()
		if op == Halt {
			close(i.out)
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
		done:    make(chan struct{}, 1),
	}
}
