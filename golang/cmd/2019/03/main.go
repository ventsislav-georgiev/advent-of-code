package main

import (
	"bufio"
	"io"
	"math"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	wire1, wire2 := parse(in)

	dist := math.MaxInt32
	for _, l1 := range wire1 {
		for _, l2 := range wire2 {
			if !l1.Intersects(l2) {
				continue
			}

			intersection := l1.Intersection(l2)
			if intersection == (Point{0, 0}) {
				continue
			}

			newdist := intersection.Distance(Point{0, 0})
			if newdist < dist {
				dist = newdist
			}
		}
	}

	println(dist)
}

func task2(in io.Reader) {
	wire1, wire2 := parse(in)

	steps := math.MaxInt32
	for _, l1 := range wire1 {
		for _, l2 := range wire2 {
			if !l1.Intersects(l2) {
				continue
			}

			intersection := l1.Intersection(l2)
			if intersection == (Point{0, 0}) {
				continue
			}

			cur := l1.length + l2.length - l1.p2.Distance(intersection) - l2.p2.Distance(intersection)
			if cur < steps {
				steps = cur
			}
		}
	}

	println(steps)
}

type Point struct {
	x, y int
}

func (p1 Point) Distance(p2 Point) int {
	return int(math.Abs(float64(p1.x-p2.x)) + math.Abs(float64(p1.y-p2.y)))
}

type Line struct {
	p1, p2 Point
	length int
}

func (l Line) A() int {
	return l.p2.y - l.p1.y
}

func (l Line) B() int {
	return l.p1.x - l.p2.x
}

func (l Line) C() int {
	return l.A()*l.p1.x + l.B()*l.p1.y
}

func (l1 Line) Intersects(l2 Line) bool {
	return ccw(l1.p1, l2.p1, l2.p2) != ccw(l1.p2, l2.p1, l2.p2) && ccw(l1.p1, l1.p2, l2.p1) != ccw(l1.p1, l1.p2, l2.p2)
}

func (l1 Line) Intersection(l2 Line) Point {
	d := l1.A()*l2.B() - l2.A()*l1.B()
	if d == 0 {
		return Point{}
	}

	x := (l2.B()*l1.C() - l1.B()*l2.C()) / d
	y := (l1.A()*l2.C() - l2.A()*l1.C()) / d

	return Point{x, y}
}

func ccw(a, b, c Point) bool {
	return (c.y-a.y)*(b.x-a.x) > (b.y-a.y)*(c.x-a.x)
}

func parse(in io.Reader) ([]Line, []Line) {
	scanner := bufio.NewScanner(in)
	wire1 := []Line{{Point{0, 0}, Point{0, 0}, 0}}
	wire2 := []Line{{Point{0, 0}, Point{0, 0}, 0}}

	parseLine := func(wire *[]Line, line string) {
		for _, s := range strings.Split(line, ",") {
			var x, y int
			v := aoc.StrToInt(s[1:])
			switch s[0] {
			case 'U':
				y = v
			case 'D':
				y = -v
			case 'R':
				x = v
			case 'L':
				x = -v
			}

			last := (*wire)[len(*wire)-1]
			p1 := last.p2
			p2 := Point{last.p2.x + x, last.p2.y + y}
			length := last.length + p1.Distance(p2)
			*wire = append(*wire, Line{p1, p2, length})
		}
	}

	scanner.Scan()
	parseLine(&wire1, scanner.Text())
	scanner.Scan()
	parseLine(&wire2, scanner.Text())

	return wire1[1:], wire2[1:]
}
