package main

import (
	"bufio"
	"container/ring"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	numbers := parse(in)
	ans := solve(numbers, 1, 1)
	println(ans)
}

func task2(in io.Reader) {
	numbers := parse(in)
	ans := solve(numbers, 10, 811589153)
	println(ans)
}

func solve(numbers []int, cycles int, decryptionKey int) int {
	r := ring.New(len(numbers))
	indexes := map[int]*ring.Ring{}
	zero := r

	for i, num := range numbers {
		if num == 0 {
			zero = r
		}

		r.Value = num * decryptionKey
		indexes[i] = r
		r = r.Next()
	}

	for ; cycles > 0; cycles-- {
		for i := 0; i < len(indexes); i++ {
			r := indexes[i].Prev()

			num := r.Unlink(1)
			nextIndex := num.Value.(int) % (len(indexes) - 1)

			r.Move(nextIndex).Link(num)
		}
	}

	return zero.Move(1000).Value.(int) + zero.Move(2000).Value.(int) + zero.Move(3000).Value.(int)
}

func parse(in io.Reader) []int {
	scanner := bufio.NewScanner(in)
	numbers := []int{}
	for scanner.Scan() {
		numbers = append(numbers, aoc.Atoi(scanner.Bytes()))
	}
	return numbers
}
