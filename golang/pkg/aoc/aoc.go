package aoc

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
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
