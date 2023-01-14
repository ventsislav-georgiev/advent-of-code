package main

import (
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	intcode := aoc.ParseIntcode(in)
	go intcode.Run()
	intcode.In <- 1
	println(intcode.WaitHalt())
}

func task2(in io.Reader) {
	intcode := aoc.ParseIntcode(in)
	go intcode.Run()
	intcode.In <- 2
	println(intcode.WaitHalt())
}
