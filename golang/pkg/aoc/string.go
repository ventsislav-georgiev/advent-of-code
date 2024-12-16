package aoc

import (
	"bufio"
	"fmt"
	"strings"
)

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

func ParseNumber(scanner *bufio.Scanner, term byte) (uint, bool) {
	var ch byte
	var number uint
	numParts := []byte{}

	parseNumber := func() {
		if len(numParts) == 0 {
			return
		}

		number = Atoui(numParts)
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
			continue
		}

		if ch == term {
			parseNumber()
			break
		} else {
			return 0, false
		}
	}

	return number, true
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

func ListDelete(list []string, item string) ([]string, bool) {
	for i, listItem := range list {
		if listItem == item {
			return append(list[:i], list[i+1:]...), true
		}
	}
	return list, false
}

func JoinStrings(strs ...string) string {
	return strings.Join(strs, "")
}

func ToStr[T any](val T) string {
	var vAny interface{} = val
	switch v := vAny.(type) {
	case byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
