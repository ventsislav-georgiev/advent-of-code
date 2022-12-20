package aoc

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Exec(tasks ...func(io.Reader)) {
	day := flag.Int("day", 17, "")
	taskNumber := flag.Int("task", 2, "")
	flag.Parse()

	task := tasks[*taskNumber-1]
	in := GetInput(*day)
	defer in.Close()
	task(in)
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

func GetInput(day int) io.ReadCloser {
	url := fmt.Sprintf("https://adventofcode.com/2022/day/%d/input", day)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: os.Getenv("SESSION_KEY")})

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	return resp.Body
}
