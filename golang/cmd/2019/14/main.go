package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	recipes := parse(in)
	println(calculateOre(recipes, 1))
}

func task2(in io.Reader) {
	recipes := parse(in)
	ore := 1000000000000
	fuel := 1
	orePerFuel := calculateOre(recipes, fuel)

	for {
		requiredOre := calculateOre(recipes, fuel)
		if requiredOre <= ore {
			oreCount := (ore - requiredOre) / orePerFuel
			if oreCount > 1000000 {
				fuel += 1000000
			} else if oreCount > 100000 {
				fuel += 100000
			} else if oreCount > 10000 {
				fuel += 10000
			} else if oreCount > 1000 {
				fuel += 1000
			} else if oreCount > 100 {
				fuel += 100
			} else if oreCount > 10 {
				fuel += 10
			} else {
				fuel += 1
			}
		} else {
			fuel--
			break
		}
	}

	println(fuel)
}

type Chemical struct {
	Name   string
	Amount int
	Recipe map[string]int
}

type HashStack map[string]int

func (s *HashStack) Pop() (string, int) {
	for name, amount := range *s {
		delete(*s, name)
		return name, amount
	}
	panic("empty stack")
}

func calculateOre(recipes map[string]Chemical, fuel int) int {
	requiredOre := 0

	chemicals := HashStack{}
	excess := map[string]int{}
	chemicals["FUEL"] = fuel

	for len(chemicals) > 0 {
		name, amount := chemicals.Pop()
		amount -= excess[name]
		excess[name] = 0

		if amount <= 0 {
			continue
		}

		if name == "ORE" {
			requiredOre += amount
			continue
		}

		recipe := recipes[name]
		recipeAmount := recipe.Amount
		recipeCount := 1

		leftover := recipeAmount - amount
		for leftover < 0 {
			recipeCount++
			leftover = (recipeCount * recipeAmount) - amount
		}

		if leftover > 0 {
			excess[name] += leftover
		}

		for name, amount := range recipe.Recipe {
			chemicals[name] += amount * recipeCount
		}
	}

	return requiredOre
}

func parse(in io.Reader) map[string]Chemical {
	scanner := bufio.NewScanner(in)

	parseChemical := func(s string) Chemical {
		parts := strings.Split(strings.TrimSpace(s), " ")
		return Chemical{
			Name:   parts[1],
			Amount: aoc.StrToInt(parts[0]),
			Recipe: map[string]int{},
		}
	}

	recipes := map[string]Chemical{}
	for scanner.Scan() {
		line := scanner.Text()
		chemicalProcess := strings.Split(line, " => ")
		chemical := parseChemical(chemicalProcess[1])
		chemicals := strings.Split(chemicalProcess[0], ",")

		recipe := map[string]int{}
		for _, c := range chemicals {
			chemical := parseChemical(c)
			recipe[chemical.Name] = chemical.Amount
		}

		chemical.Recipe = recipe
		recipes[chemical.Name] = chemical
	}

	return recipes
}
