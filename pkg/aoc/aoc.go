package aoc

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetTask(tasks ...func(io.ReadCloser)) func(io.ReadCloser) {
	tID := flag.Int("task", 1, "")
	flag.Parse()

	var task func(io.ReadCloser)
	switch *tID {
	case 1:
		task = tasks[0]
	case 2:
		task = tasks[1]
	}

	return task
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
