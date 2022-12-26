package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

var testInput = `1=-0-2
12111
2=0=
21
2=01
111
20012
112
1=-1=
1-12
12
1=
122`

var testMode = false

func main() {
	aoc.Exec(task1, task1)
}

func task1(in io.Reader) {
	if testMode {
		in = strings.NewReader(testInput)
	}

	numbers := parse(in)

	var sum int
	for _, number := range numbers {
		sum += fromSNAFU(number)
	}

	println(sum)
	println(toSNAFU(sum))
}

var toMap = map[int]byte{
	-2: '=',
	-1: '-',
	0:  '0',
	1:  '1',
	2:  '2',
}

var fromMap = map[byte]int{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

func fromSNAFU(snafu []byte) int {
	num, power := 0, 1
	for i := len(snafu) - 1; i >= 0; i-- {
		num += power * fromMap[snafu[i]]
		power *= 5
		snafu = snafu[:i]
	}
	return num
}

func toSNAFU(num int) string {
	var snafu []byte
	for ; num > 0; num /= 5 {
		num += 2
		snafu = append(snafu, toMap[num%5-2])
	}
	return string(aoc.Reverse(snafu))
}

func parse(in io.Reader) [][]byte {
	var numbers [][]byte
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		numbers = append(numbers, scanner.Bytes())
	}
	return numbers
}
