package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(solve, solve)
}

type Elem struct {
	val  *uint
	list *[]Elem
}

type Order int

const (
	Left  Order = -1
	Right Order = 1
	Equal Order = 0
)

func solve(in io.Reader) {
	scanner := bufio.NewScanner(in)
	rightOrderIdxSum := 0

	div1 := parseLine([]byte("[[2]]"))
	div2 := parseLine([]byte("[[6]]"))
	packets := []Elem{div1, div2}

	idx := 1
	for scanner.Scan() {
		left := parseLine(scanner.Bytes())
		scanner.Scan()
		right := parseLine(scanner.Bytes())
		scanner.Scan()

		order := compare(left, right)
		if order == Right {
			rightOrderIdxSum += idx
		}

		idx += 1
		packets = append(packets, left, right)
	}

	fmt.Printf("Sum of right order indices: %d\n", rightOrderIdxSum)

	sort.Slice(packets, func(i, j int) bool {
		return compare(packets[i], packets[j]) == Right
	})

	dividersIdxProduct := 1
	for i, elem := range packets {
		if elem == div1 || elem == div2 {
			dividersIdxProduct *= i + 1
		}
	}

	fmt.Printf("Product of dividers indices: %d\n", dividersIdxProduct)
}

func parseLine(line []byte) Elem {
	elem, _ := parse(line)
	return (*elem)[0]
}

func parse(list []byte) (*[]Elem, int) {
	items := []Elem{}
	for i := 0; i < len(list); i++ {
		ch := list[i]
		if ch == ',' {
			continue
		}

		if ch == ']' {
			return &items, i + 1
		}

		item := Elem{}
		if ch == '[' {
			tmp, idx := parse(list[i+1:])
			i += idx
			item.list = tmp
		} else {
			num := []byte{ch}
			for i+1 < len(list) && list[i+1] >= '0' && list[i+1] <= '9' {
				i++
				num = append(num, list[i])
			}
			tmp := aoc.Atoui(num)
			item.val = &tmp
		}

		items = append(items, item)
	}

	return &items, 0
}

func compare(left, right Elem) Order {
	if left.val != nil && right.val != nil {
		if *left.val < *right.val {
			return Right
		}
		if *left.val > *right.val {
			return Left
		}
		return Equal
	}

	if left.list != nil && right.list != nil {
		for i := 0; i < len(*left.list) || i < len(*right.list); i++ {
			if i >= len(*left.list) {
				return Right
			}
			if i >= len(*right.list) {
				return Left
			}
			order := compare((*left.list)[i], (*right.list)[i])
			if order != Equal {
				return order
			}
		}
		return Equal
	}

	if left.val != nil && right.list != nil {
		left.list = &[]Elem{left}
		left.val = nil
		return compare(left, right)
	}

	if left.list != nil && right.val != nil {
		right.list = &[]Elem{right}
		right.val = nil
		return compare(left, right)
	}

	return Equal
}
