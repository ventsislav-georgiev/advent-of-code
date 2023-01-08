package main

import (
	"bufio"
	"io"
	"math"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	asteroids := parse(in)

	var maxVisible int
	var best Asteroid
	for _, asteroid := range asteroids {
		visible := 0

		for _, other := range asteroids {
			if asteroid == other {
				continue
			}

			sightAngle := asteroid.Pos.Angle(other.Pos)

			var isBlocked bool
			for _, blocker := range asteroids {
				if blocker == asteroid || blocker == other {
					continue
				}

				if sightAngle == asteroid.Pos.Angle(blocker.Pos) && asteroid.Pos.Distance(blocker.Pos) < asteroid.Pos.Distance(other.Pos) {
					isBlocked = true
					break
				}
			}

			if !isBlocked {
				visible++
			}
		}

		if visible > maxVisible {
			maxVisible = visible
			best = asteroid
		}
	}

	println("Detected:", maxVisible)
	println("X:", best.Pos.x, "Y:", best.Pos.y)
}

func task2(in io.Reader) {
}

type Point struct {
	x int
	y int
}

func (p Point) Angle(other Point) float64 {
	return math.Atan2(float64(other.y-p.y), float64(other.x-p.x))
}

func (p Point) Distance(other Point) float64 {
	return math.Sqrt(math.Pow(float64(other.x-p.x), 2) + math.Pow(float64(other.y-p.y), 2))
}

type Asteroid struct {
	Pos Point
}

func parse(in io.Reader) []Asteroid {
	scanner := bufio.NewScanner(in)
	asteroids := []Asteroid{}

	for y := 0; scanner.Scan(); y++ {
		for x, c := range scanner.Bytes() {
			if c != '#' {
				continue
			}

			asteroids = append(asteroids, Asteroid{Point{x, y}})
		}
	}

	return asteroids
}
