package aoc

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Exec(tasks ...func(io.Reader)) {
	year := flag.Int("year", 2019, "")
	dayStr := flag.String("day", "14", "")
	taskNumber := flag.Int("task", 2, "")
	flag.Parse()

	task := tasks[*taskNumber-1]
	day := StrToInt(strings.TrimLeft(*dayStr, "0"))
	in := GetInput(*year, day)
	if in != nil {
		defer in.Close()
	}
	task(in)
}

func StrToInt(s string) int {
	return Atoi([]byte(s))
}

func Atoi(s []byte) int {
	s0 := s
	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
	}
	n := 0
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
	n := uint(0)
	for _, ch := range s {
		ch -= '0'
		n = n*10 + uint(ch)
	}
	return n
}

func Reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func GetInput(year, day int) io.ReadCloser {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: os.Getenv("SESSION_KEY")})

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	return resp.Body
}
