package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	monkeys, _ := parse(in)
	solve(monkeys, 20, func(item int) int { return item / 3 })
}

func task2(in io.Reader) {
	monkeys, modLCM := parse(in)
	solve(monkeys, 10000, func(item int) int { return item % modLCM })
}

type Monkey struct {
	items  []int
	target func(int) *Monkey
	op     func(int) int
}

func solve(monkeys *[]*Monkey, rounds int, fn func(int) int) {
	activity := make([]int, len(*monkeys))

	for ; rounds > 0; rounds-- {
		for i, monkey := range *monkeys {
			for _, item := range monkey.items {
				activity[i]++
				item = fn(monkey.op(item))
				target := monkey.target(item)
				target.items = append(target.items, item)
			}
			monkey.items = []int{}
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(activity)))
	fmt.Println(activity[0] * activity[1])
}

func parse(in io.Reader) (*[]*Monkey, int) {
	monkeys := make([]*Monkey, 0, 8)
	modLCM := 1

	input, _ := io.ReadAll(in)
	format := `Monkey %d:
  Starting items: %s
  Operation: new = old %s %d
  Test: divisible by %d
	If true: throw to monkey %d
	If false: throw to monkey %d`

	for _, data := range strings.Split(string(input), "\n\n") {
		data = strings.NewReplacer(", ", ",", "* old", "^ 2").Replace(data)

		var items, op string
		var i, arg, mod, target, fallback int

		fmt.Sscanf(data, format, &i, &items, &op, &arg, &mod, &target, &fallback)
		modLCM *= mod

		monkey := &Monkey{
			items: make([]int, 0, 8),
		}
		for _, item := range strings.Split(items, ",") {
			monkey.items = append(monkey.items, aoc.Atoi([]byte(item)))
		}

		switch op[0] {
		case '+':
			monkey.op = func(item int) int { return item + arg }
		case '*':
			monkey.op = func(item int) int { return item * arg }
		case '^':
			monkey.op = func(item int) int { return item * item }
		}

		monkey.target = func(item int) *Monkey {
			if item%mod == 0 {
				return monkeys[target]
			} else {
				return monkeys[fallback]
			}
		}

		monkeys = append(monkeys, monkey)
	}

	return &monkeys, modLCM
}
