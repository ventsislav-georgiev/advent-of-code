package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	monkeys := parse(in)
	fmt.Println(solve(monkeys, monkeys["root"]))
}

func task2(in io.Reader) {
	monkeys := parse(in)

	root := monkeys["root"]
	root.operation = '='
	humnValue := 1
	solve(monkeys, root)

	result := func() (int, int, bool) {
		return *root.leftValue, *root.rightValue, *root.leftValue == *root.rightValue
	}

	var multiplier float64
	leftValue, rightValue, isEqual := result()
	for !isEqual {
		multiplier = (float64(leftValue) / float64(rightValue))
		humnValue = int(float64(humnValue) * multiplier)
		monkeys["humn"].value = humnValue
		solve(monkeys, root)
		leftValue, rightValue, isEqual = result()
	}

	fmt.Println(humnValue)
}

func solve(monkeys map[string]*Monkey, monkey *Monkey) int {
	if monkey.operation == 0 {
		return monkey.value
	}

	left := monkeys[monkey.waiting[0]]
	right := monkeys[monkey.waiting[1]]

	result := -1
	leftValue := solve(monkeys, left)
	rightValue := solve(monkeys, right)

	switch monkey.operation {
	case '+':
		result = leftValue + rightValue
	case '-':
		result = leftValue - rightValue
	case '*':
		result = leftValue * rightValue
	case '/':
		result = leftValue / rightValue
	case '=':
		monkey.leftValue = &leftValue
		monkey.rightValue = &rightValue
		if leftValue == rightValue {
			result = 1
		}
	}

	return result
}

type Monkey struct {
	value     int
	operation byte
	waiting   []string

	leftValue  *int
	rightValue *int
}

func parse(in io.Reader) map[string]*Monkey {
	scanner := bufio.NewScanner(in)
	monkeys := map[string]*Monkey{}
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ": ")
		monkey := parts[0]

		if len(line) > 10 {
			var left, right string
			var operation byte
			fmt.Sscanf(parts[1], "%s %c %s", &left, &operation, &right)
			monkeys[monkey] = &Monkey{
				operation: operation,
				waiting:   []string{left, right},
			}
		} else {
			v, _ := strconv.Atoi(parts[1])
			monkeys[monkey] = &Monkey{value: v}
		}
	}
	return monkeys
}
