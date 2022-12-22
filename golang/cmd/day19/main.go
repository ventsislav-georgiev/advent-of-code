package main

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	time := 24
	blueprints := parse(in, time)
	ans := solve(blueprints, true, time)

	sum := 0
	for _, r := range ans {
		sum += r
	}

	println(sum)
}

func task2(in io.Reader) {
	time := 32
	blueprints := parse(in, time)[:3]
	ans := solve(blueprints, false, time)

	product := 1
	for _, r := range ans {
		println(r)
		product *= r
	}

	println(product)
}

func solve(blueprints []Blueprint, quality bool, time int) []int {
	results := make(chan int, len(blueprints))
	var wg sync.WaitGroup

	for _, b := range blueprints {
		blueprint := b
		wg.Add(1)

		go func() {
			defer wg.Done()
			cache := map[Blueprint]int{}
			r := calc(blueprint, ore, time, cache)
			if quality {
				r *= blueprint.id
			}

			results <- r
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	ans := []int{}
	for r := range results {
		ans = append(ans, r)
	}

	return ans
}

func calc(b Blueprint, target uint8, maxTime int, cache map[Blueprint]int) int {
	create := false
	switch target {
	case ore:
		if b.resouces.ore >= b.cost.oreRobot.ore {
			b.resouces.ore -= b.cost.oreRobot.ore
			create = true
		}
	case clay:
		if b.resouces.ore >= b.cost.clayRobot.ore {
			b.resouces.ore -= b.cost.clayRobot.ore
			create = true
		}
	case obsidian:
		if b.resouces.ore >= b.cost.obsidianRobot.ore && b.resouces.clay >= b.cost.obsidianRobot.clay {
			b.resouces.ore -= b.cost.obsidianRobot.ore
			b.resouces.clay -= b.cost.obsidianRobot.clay
			create = true
		}
	case geode:
		if b.resouces.ore >= b.cost.geodeRobot.ore && b.resouces.obsidian >= b.cost.geodeRobot.obsidian {
			b.resouces.ore -= b.cost.geodeRobot.ore
			b.resouces.obsidian -= b.cost.geodeRobot.obsidian
			create = true
		}
	}

	b.resouces.ore += b.robots.ore
	b.resouces.clay += b.robots.clay
	b.resouces.obsidian += b.robots.obsidian
	b.resouces.geode += b.robots.geode

	if create {
		switch target {
		case ore:
			b.robots.ore += 1
		case clay:
			b.robots.clay += 1
		case obsidian:
			b.robots.obsidian += 1
		case geode:
			b.robots.geode += 1
		}
	}

	b.time -= 1
	if b.time == 0 {
		return b.resouces.geode
	}

	if r, ok := cache[b]; ok {
		return r
	}

	if b.robots.ore >= b.cost.geodeRobot.ore && b.robots.obsidian >= b.cost.geodeRobot.obsidian {
		r := calc(b, geode, maxTime, cache)
		cache[b] = r
		return r
	}

	r := 0
	if b.robots.obsidian > 0 {
		r = max(r, calc(b, geode, maxTime, cache))
	}
	if b.time > 2 && b.robots.ore < 4 {
		r = max(r, calc(b, ore, maxTime, cache))
	}
	if b.time > 3 && b.robots.clay < 8 {
		r = max(r, calc(b, clay, maxTime, cache))
	}
	if b.time > 2 && b.robots.clay > 0 && b.robots.obsidian < 8 {
		r = max(r, calc(b, obsidian, maxTime, cache))
	}

	cache[b] = r
	return r
}

func parse(in io.Reader, time int) []Blueprint {
	scanner := bufio.NewScanner(in)
	blueprints := []Blueprint{}
	for scanner.Scan() {
		id, o, c, ob1, ob2, g1, g2 := 0, 0, 0, 0, 0, 0, 0
		fmt.Sscanf(scanner.Text(), "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.", &id, &o, &c, &ob1, &ob2, &g1, &g2)

		blueprints = append(blueprints, Blueprint{
			id:       id,
			time:     time,
			robots:   RobotsAndResouces{ore: 1},
			resouces: RobotsAndResouces{},
			cost: Cost{
				oreRobot: OreRobotCost{
					ore: o,
				},
				clayRobot: ClayRobotCost{
					ore: c,
				},
				obsidianRobot: ObsidianRobotCost{
					ore:  ob1,
					clay: ob2,
				},
				geodeRobot: GeodeRobotCost{
					ore:      g1,
					obsidian: g2,
				},
			},
		})
	}
	return blueprints
}

func max(a ...int) int {
	m := 0
	for _, v := range a {
		if v > m {
			m = v
		}
	}
	return m
}

const (
	ore uint8 = iota
	clay
	obsidian
	geode
)

type Blueprint struct {
	id       int
	robots   RobotsAndResouces
	resouces RobotsAndResouces
	cost     Cost
	time     int
}

type RobotsAndResouces struct {
	ore      int
	clay     int
	obsidian int
	geode    int
}

type Cost struct {
	oreRobot      OreRobotCost
	clayRobot     ClayRobotCost
	obsidianRobot ObsidianRobotCost
	geodeRobot    GeodeRobotCost
}

type OreRobotCost struct {
	ore int
}

type ClayRobotCost struct {
	ore int
}

type ObsidianRobotCost struct {
	ore  int
	clay int
}

type GeodeRobotCost struct {
	ore      int
	obsidian int
}
