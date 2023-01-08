package main

import (
	"bufio"
	"bytes"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	planets := parse(in)
	var count int
	for _, p := range planets {
		for p.Orbiting != nil {
			count++
			p = p.Orbiting
		}
	}
	println(count)
}

func task2(in io.Reader) {
	planets := parse(in)
	src := planets["YOU"]
	dest := planets["SAN"]

	srcOrbiting := make(map[*Planet]int)
	distance := 0
	for src.Orbiting != nil {
		srcOrbiting[src.Orbiting] = distance
		src = src.Orbiting
		distance++
	}

	var totalDistance int
	for dest.Orbiting != nil {
		if distance, ok := srcOrbiting[dest.Orbiting]; ok {
			totalDistance += distance
			break
		}
		totalDistance++
		dest = dest.Orbiting
	}

	println(totalDistance)
}

type Planet struct {
	Orbiting *Planet
}

func parse(in io.Reader) map[string]*Planet {
	scanner := bufio.NewScanner(in)
	planets := make(map[string]*Planet)

	for scanner.Scan() {
		line := scanner.Bytes()
		sepIndex := bytes.IndexByte(line, ')')

		planetName := string(line[:sepIndex])
		orbitingName := string(line[sepIndex+1:])

		planet, ok := planets[planetName]
		if !ok {
			planet = &Planet{}
			planets[planetName] = planet
		}

		orbiting, ok := planets[orbitingName]
		if !ok {
			orbiting = &Planet{}
			planets[orbitingName] = orbiting
		}

		orbiting.Orbiting = planet
	}

	return planets
}
