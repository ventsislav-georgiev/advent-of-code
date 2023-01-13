package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	moon1, moon2, moon3, moon4 := parse(in)
	allaxis := []bool{true, true, true}
	for i := 0; i < 1000; i++ {
		applyGravity(allaxis, &moon1, &moon2, &moon3, &moon4)
		applyVelocity(allaxis, &moon1, &moon2, &moon3, &moon4)
	}

	fmt.Println(energy(&moon1) + energy(&moon2) + energy(&moon3) + energy(&moon4))
}

func task2(in io.Reader) {
	moon1, moon2, moon3, moon4 := parse(in)
	initial := [4]Moon{moon1, moon2, moon3, moon4}

	xaxis := []bool{true, false, false}
	yaxis := []bool{false, true, false}
	zaxis := []bool{false, false, true}
	xperiod := findPeriod(xaxis, moon1, moon2, moon3, moon4, initial)
	yperiod := findPeriod(yaxis, moon1, moon2, moon3, moon4, initial)
	zperiod := findPeriod(zaxis, moon1, moon2, moon3, moon4, initial)

	fmt.Println(lcm(xperiod, yperiod, zperiod))
}

func findPeriod(axis []bool, moon1, moon2, moon3, moon4 Moon, inital [4]Moon) int {
	period := 0
	for {
		applyGravity(axis, &moon1, &moon2, &moon3, &moon4)
		applyVelocity(axis, &moon1, &moon2, &moon3, &moon4)
		period++

		if moon1 == inital[0] && moon2 == inital[1] && moon3 == inital[2] && moon4 == inital[3] {
			return period
		}
	}
}

func lcm(a, b, c int) int {
	return lcm2(a, lcm2(b, c))
}

func lcm2(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func applyGravity(axis []bool, moons ...*Moon) {
	for i := 0; i < len(moons); i++ {
		for j := i + 1; j < len(moons); j++ {
			applyGravityPair(axis, moons[i], moons[j])
		}
	}
}

func applyGravityPair(axis []bool, moon1, moon2 *Moon) {
	if axis[0] {
		if moon1.Pos.X < moon2.Pos.X {
			moon1.Vel.X++
			moon2.Vel.X--
		} else if moon1.Pos.X > moon2.Pos.X {
			moon1.Vel.X--
			moon2.Vel.X++
		}
	}

	if axis[1] {
		if moon1.Pos.Y < moon2.Pos.Y {
			moon1.Vel.Y++
			moon2.Vel.Y--
		} else if moon1.Pos.Y > moon2.Pos.Y {
			moon1.Vel.Y--
			moon2.Vel.Y++
		}
	}

	if axis[2] {
		if moon1.Pos.Z < moon2.Pos.Z {
			moon1.Vel.Z++
			moon2.Vel.Z--
		} else if moon1.Pos.Z > moon2.Pos.Z {
			moon1.Vel.Z--
			moon2.Vel.Z++
		}
	}
}

func applyVelocity(axis []bool, moons ...*Moon) {
	for _, moon := range moons {
		if axis[0] {
			moon.Pos.X += moon.Vel.X
		}
		if axis[1] {
			moon.Pos.Y += moon.Vel.Y
		}
		if axis[2] {
			moon.Pos.Z += moon.Vel.Z
		}
	}
}

func energy(moon *Moon) int {
	potential := abs(moon.Pos.X) + abs(moon.Pos.Y) + abs(moon.Pos.Z)
	kinetic := abs(moon.Vel.X) + abs(moon.Vel.Y) + abs(moon.Vel.Z)
	return potential * kinetic
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Moon struct {
	Pos Vec
	Vel Vec
}

type Vec struct {
	X, Y, Z int
}

func parse(in io.Reader) (Moon, Moon, Moon, Moon) {
	scanner := bufio.NewScanner(in)
	format := "<x=%d, y=%d, z=%d>"
	positions := []Vec{}

	for scanner.Scan() {
		var pos Vec
		fmt.Sscanf(scanner.Text(), format, &pos.X, &pos.Y, &pos.Z)
		positions = append(positions, pos)
	}

	return Moon{Pos: positions[0]}, Moon{Pos: positions[1]}, Moon{Pos: positions[2]}, Moon{Pos: positions[3]}
}
