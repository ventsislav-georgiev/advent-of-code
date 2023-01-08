package main

import (
	"bufio"
	"io"
	"math"
	"sort"

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
		visible := asteroid.GetVisible(asteroids)
		if len(visible) > maxVisible {
			maxVisible = len(visible)
			best = asteroid
		}
	}

	println("Detected:", maxVisible)
	println("X:", best.Pos.x, "Y:", best.Pos.y)
}

func task2(in io.Reader) {
	asteroids := parse(in)
	asteroid := Asteroid{Pos: Point{x: 37, y: 25}}
	visible := asteroid.GetVisible(asteroids)

	degrees90 := math.Pi / 2
	degrees360 := math.Pi * 2

	sort.Slice(visible, func(i, j int) bool {
		angle1 := *visible[i].Angle + degrees90
		angle2 := *visible[j].Angle + degrees90

		if angle1 < 0 {
			angle1 += degrees360
		}

		if angle2 < 0 {
			angle2 += degrees360
		}

		return angle1 < angle2
	})

	v := visible[199]
	println(v.Pos.x*100 + v.Pos.y)
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
	Pos   Point
	Angle *float64
}

func (a Asteroid) GetVisible(asteroids []Asteroid) []Asteroid {
	visible := []Asteroid{}

	for _, other := range asteroids {
		if a == other {
			continue
		}

		sightAngle := a.Pos.Angle(other.Pos)
		for _, blocker := range asteroids {
			if blocker == a || blocker == other {
				continue
			}

			if sightAngle == a.Pos.Angle(blocker.Pos) && a.Pos.Distance(blocker.Pos) < a.Pos.Distance(other.Pos) {
				goto next
			}
		}

		other.Angle = &sightAngle
		visible = append(visible, other)

	next:
	}

	return visible
}

func parse(in io.Reader) []Asteroid {
	scanner := bufio.NewScanner(in)
	asteroids := []Asteroid{}

	for y := 0; scanner.Scan(); y++ {
		for x, c := range scanner.Bytes() {
			if c != '#' {
				continue
			}

			asteroids = append(asteroids, Asteroid{Point{x, y}, nil})
		}
	}

	return asteroids
}
