package aoc

import "bufio"

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

func Atoui(s []byte) uint {
	var n uint
	for _, ch := range s {
		ch -= '0'
		n = n*10 + uint(ch)
	}
	return n
}

func ListContains(list []string, item string) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}
	return false
}
