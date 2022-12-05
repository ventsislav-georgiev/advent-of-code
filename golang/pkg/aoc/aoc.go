package aoc

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Exec(tasks ...func(io.ReadCloser)) {
	day := flag.Int("day", 5, "")
	taskNumber := flag.Int("task", 2, "")
	flag.Parse()

	task := tasks[*taskNumber-1]
	in := getInput(*day)
	defer in.Close()
	task(in)
}

func getInput(day int) io.ReadCloser {
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
