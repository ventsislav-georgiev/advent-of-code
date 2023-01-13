package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

var in *bytes.Reader

func init() {
	input := aoc.GetInput(2019, 11)
	b, _ := io.ReadAll(input)
	input.Close()

	in = bytes.NewReader(b)
}

func BenchmarkSolution(b *testing.B) {
	for i := 0; i < b.N; i++ {
		in.Seek(0, io.SeekStart)
		solve(in, 1)
	}
}
