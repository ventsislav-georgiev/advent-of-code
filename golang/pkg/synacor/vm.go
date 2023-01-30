package synacor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	size     = 1 << 15
	memLimit = size + 8
)

type Opcode int

const (
	Halt Opcode = iota
	Set
	Push
	Pop
	Eq
	Gt
	Jmp
	Jt
	Jf
	Add
	Mult
	Mod
	And
	Or
	Not
	Rmem
	Wmem
	Call
	Ret
	Out
	In
	Noop
)

var Opcodes = map[Opcode]struct {
	name string
	args int
}{
	Halt: {"halt", 0},
	Set:  {"set", 2},
	Push: {"push", 1},
	Pop:  {"pop", 1},
	Eq:   {"eq", 3},
	Gt:   {"gt", 3},
	Jmp:  {"jmp", 1},
	Jt:   {"jt", 2},
	Jf:   {"jf", 2},
	Add:  {"add", 3},
	Mult: {"mult", 3},
	Mod:  {"mod", 3},
	And:  {"and", 3},
	Or:   {"or", 3},
	Not:  {"not", 2},
	Rmem: {"rmem", 2},
	Wmem: {"wmem", 2},
	Call: {"call", 1},
	Ret:  {"ret", 0},
	Out:  {"out", 1},
	In:   {"in", 1},
	Noop: {"noop", 0},
}

type VM struct {
	mem   []uint16
	idx   uint16
	in    chan uint16
	halt  chan struct{}
	cwd   string
	debug io.Writer
}

func (vm *VM) Run() {
	vm.listenInput()

	for {
		opcode := Opcode(*vm.operand())
		if opcode == Halt {
			vm.halt <- struct{}{}
			close(vm.halt)
			return
		}

		vm.step(opcode)
	}
}

func (vm *VM) WaitHalt() {
	<-vm.halt
}

func (vm *VM) SetMem(index int, value uint16) {
	vm.mem[index] = value
}

func (vm *VM) InputLine(line string) {
	defer fmt.Println(line)
	vm.inputLine(line)
}

func (vm *VM) inputLine(line string) {
	if line[0] == ':' {
		switch line[1] {
		case 'q':
			os.Exit(0)
		case 'r':
			var reg, val uint16
			fmt.Sscanf(line, ":r%d %d", &reg, &val)
			vm.mem[size+reg] = val
			return
		case 'i':
			for i := 0; i < 8; i++ {
				fmt.Printf("r%d: %d\n", i, vm.mem[size+i])
			}
			return
		case 'd':
			debugFilePath := vm.cwd + "/debug.txt"
			os.RemoveAll(debugFilePath)
			debugFile, _ := os.Create(debugFilePath)
			vm.debug = debugFile
			return
		}

		return
	}

	for _, c := range line {
		vm.in <- uint16(c)
	}
	vm.in <- 10 // newline
}

func (vm *VM) listenInput() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			switch text {
			case "n":
				text = "north"
			case "s":
				text = "south"
			case "e":
				text = "east"
			case "w":
				text = "west"
			}
			vm.inputLine(text)
		}
	}()
}

func (vm *VM) step(opcode Opcode) {
	vm.printDebug("\t&%d: %s", vm.idx-1, Opcodes[opcode].name)

	switch opcode {
	case Set:
		a := vm.operand()
		b := vm.operand()
		*a = *b

	case Push:
		a := vm.operand()
		vm.mem = append(vm.mem, *a)

	case Pop:
		a := vm.operand()
		*a, vm.mem = vm.mem[vm.memLen()], vm.mem[:vm.memLen()]

	case Eq:
		a := vm.operand()
		b := vm.operand()
		c := vm.operand()
		if *b == *c {
			*a = 1
		} else {
			*a = 0
		}

	case Gt:
		a := vm.operand()
		b := vm.operand()
		c := vm.operand()
		if *b > *c {
			*a = 1
		} else {
			*a = 0
		}

	case Jmp:
		a := vm.operand()
		vm.idx = *a

	case Jt:
		a := vm.operand()
		b := vm.operand()
		if *a != 0 {
			vm.idx = *b
		}

	case Jf:
		a := vm.operand()
		b := vm.operand()
		if *a == 0 {
			vm.idx = *b
		}

	case Add:
		a := vm.operand()
		b := vm.operand()
		c := vm.operand()
		*a = (*b + *c) % size

	case Mult:
		a := vm.operand()
		b := vm.operand()
		c := vm.operand()
		*a = (*b * *c) % size

	case Mod:
		a := vm.operand()
		b := vm.operand()
		c := vm.operand()
		*a = *b % *c

	case And:
		a := vm.operand()
		b := vm.operand()
		c := vm.operand()
		*a = *b & *c

	case Or:
		a := vm.operand()
		b := vm.operand()
		c := vm.operand()
		*a = *b | *c

	case Not:
		a := vm.operand()
		b := vm.operand()
		*a = ^*b % size

	case Rmem:
		a := vm.operand()
		b := vm.operand()
		*a = vm.mem[*b]

	case Wmem:
		a := vm.operand()
		b := vm.operand()
		vm.mem[*a] = *b

	case Call:
		vm.idx, vm.mem = *vm.operand(), append(vm.mem, vm.idx)

	case Ret:
		vm.idx, vm.mem = vm.mem[vm.memLen()], vm.mem[:vm.memLen()]
		vm.printDebug(" %d", vm.idx)

	case Out:
		fmt.Printf("%c", *vm.operand())

	case In:
		a := vm.operand()
		*a = <-vm.in

	case Noop:
		break

	default:
		panic("unknown opcode")
	}

	vm.printDebug("\nop:")
}

func (vm *VM) memLen() int {
	return len(vm.mem) - 1
}

func (vm *VM) operand() *uint16 {
	defer func() { vm.idx++ }()
	param := vm.mem[vm.idx]
	if param >= size {
		vm.printDebug(" r%d(%d)", param%size, vm.mem[param])
		return &vm.mem[param]
	}
	vm.printDebug(" %2d", param)
	return &param
}

func (vm *VM) printDebug(format string, a ...any) {
	if vm.debug == nil {
		return
	}
	fmt.Fprintf(vm.debug, format, a...)
}

func Parse(in io.Reader) *VM {
	program := make([]uint16, memLimit)

	index := 0
	buffer := make([]uint8, 2)
	_, err := in.Read(buffer)
	for err == nil {
		program[index] = uint16(buffer[0]) | uint16(buffer[1])<<8
		index++
		_, err = in.Read(buffer)
	}

	cwd, _ := os.Getwd()
	if !strings.HasSuffix(cwd, "synacor") {
		cwd = filepath.Join(cwd, "golang/cmd/synacor")
	}

	return &VM{
		mem:  program,
		in:   make(chan uint16, 1),
		halt: make(chan struct{}, 1),
		cwd:  cwd,
	}
}

func Disasm(vm *VM) {
	lastOpcode := Halt
	var lastChar uint16
	for idx := 0; idx < 30051; {
		curOpcode := Opcode(vm.mem[idx])
		isOut := curOpcode == Out
		startOut := isOut && lastOpcode != Out
		endOut := !isOut && lastOpcode == Out
		lastOpcode = curOpcode

		opcode, ok := Opcodes[curOpcode]
		if !ok {
			idx++
			continue
		}

		if endOut && lastChar != 10 {
			fmt.Println()
		}

		if !isOut || startOut {
			fmt.Printf("%d: %s", idx, opcode.name)
		}

		if startOut {
			fmt.Printf(" ")
		}

		for argIdx := 1; argIdx <= opcode.args; argIdx++ {
			val := vm.mem[idx+argIdx]
			reg := false
			if val >= size {
				reg = true
				val = vm.mem[val]
			}
			if isOut {
				lastChar = val
				fmt.Printf("%c", val)
				continue
			}
			if reg {
				fmt.Printf(" r%d", vm.mem[idx+argIdx]%size)
				continue
			}
			fmt.Printf(" %d", val)
		}

		if !isOut {
			fmt.Println()
		}

		idx += opcode.args + 1
	}
}
