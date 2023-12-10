package aoc

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"golang.org/x/exp/constraints"
)

func Exec(tasks ...func(io.Reader)) {
	_, file, _, _ := runtime.Caller(1)
	fileParts := strings.Split(file, "/")
	year := StrToInt(fileParts[len(fileParts)-3])
	dayStr := fileParts[len(fileParts)-2]
	task := flag.String("task", "", "")
	input := flag.String("input", "", "")
	flag.Parse()

	if *task != "" {
		taskNumber := StrToInt(*task)
		task := tasks[taskNumber-1]
		tasks = []func(io.Reader){task}
	}

	for _, task := range tasks {
		day := StrToInt(strings.TrimLeft(dayStr, "0"))
		in := GetInput(*input, year, day)
		if in != nil {
			defer in.Close()
		}
		task(in)
	}
}

func StrToInt(s string) int {
	return Atoi([]byte(s))
}

func Atoi(s []byte) int {
	s0 := s
	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
	}
	var n int
	for _, ch := range s {
		ch -= '0'
		n = n*10 + int(ch)
	}
	if s0[0] == '-' {
		n = -n
	}
	return n
}

func ToKeyXY(xy [2]int) uint64 {
	return ToKey(xy[0], xy[1])
}

func ToKey(x, y int) uint64 {
	return uint64(x)<<16 + uint64(y)
}

func Atoui(s []byte) uint {
	var n uint
	for _, ch := range s {
		ch -= '0'
		n = n*10 + uint(ch)
	}
	return n
}

func Reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func GetInput(input string, year, day int) io.ReadCloser {
	if input != "" {
		file, err := os.Open(input)
		if err != nil {
			panic(err)
		}

		return file
	}

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: os.Getenv("SESSION_KEY")})

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("status code: %d", resp.StatusCode))
	}

	return resp.Body
}

func ParseNumbers(scanner *bufio.Scanner, term byte) []uint {
	var ch byte
	numbers := []uint{}
	numParts := []byte{}

	parseNumber := func() {
		if len(numParts) == 0 {
			return
		}
		numbers = append(numbers, Atoui(numParts))
		numParts = numParts[:0]
	}

	for ch != term {
		if !scanner.Scan() {
			parseNumber()
			break
		}

		ch = scanner.Bytes()[0]

		if ch >= '0' && ch <= '9' {
			numParts = append(numParts, ch)
		}

		if ch == ' ' || ch == term {
			parseNumber()
		}
	}

	return numbers
}

func SplitBytes(str []byte, sep byte) [][]byte {
	out := [][]byte{}
	part := []byte{}

	for _, ch := range str {
		if ch == '\n' {
			continue
		}

		if ch != ' ' {
			part = append(part, ch)
		} else {
			out = append(out, part)
			part = []byte{}
		}
	}

	out = append(out, part)

	return out
}

func RemoveSpaces(str []byte) []byte {
	out := make([]byte, 0, len(str))

	for _, ch := range str {
		if ch != ' ' {
			out = append(out, ch)
		}
	}

	return out
}

func SkipLine(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if scanner.Bytes()[0] == '\n' {
			break
		}
	}
}

func LastIdx[T any](arr []T) int {
	return len(arr) - 1
}

func LastElement[T any](arr []T) T {
	return arr[len(arr)-1]
}

func Max[T constraints.Ordered](values ...T) T {
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func Min[T constraints.Ordered](values ...T) T {
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func LeastCommonDenominator(a, b uint) uint {
	return a * b / GreatestCommonDivisor(a, b)
}

func GreatestCommonDivisor(a, b uint) uint {
	if b == 0 {
		return a
	}

	return GreatestCommonDivisor(b, a%b)
}
