package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"math"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	grid := parse(in)
	fmt.Println(leastStepsToAllKeys(grid))
}

func task2(in io.Reader) {
}

func leastStepsToAllKeys(grid [][]string) int {
	keysGraph := &KeysGraph{
		keysToKeys: map[string]map[string]*Node{},
	}

	keysXY := getKeysXY(grid)
	for keyChar, srcPos := range keysXY {
		// generate a new distance distGridForKey for every key
		distGridForKey := initGrid(grid, srcPos)

		// step through grid until complete to calculate distances
		for !distGridForKey.stepBFS() {
		}

		// update graph for all keys to all other keys
		for destKey, dest := range keysXY {
			if keyChar == destKey || destKey == "@" {
				// skip the entry point and the current key
				continue
			}

			destKeyNode := distGridForKey.grid[dest.Y][dest.X]
			keysGraph.addEdge(keyChar, destKeyNode)
		}
	}

	return keysGraph.dfsMinDistToAllKeys()
}

type Grid struct {
	grid  [][]*Node
	queue []image.Point
}

type Node struct {
	value      string
	distance   int
	keysFound  map[string]bool
	keysNeeded map[string]bool
	seen       bool
}

var directions = []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

func (g *Grid) stepBFS() (queueIsEmpty bool) {
	pos := g.queue[0]
	cell := g.grid[pos.Y][pos.X]
	cell.seen = true

	ch := cell.value[0]

	// key is found
	if ch >= 'a' && ch <= 'z' {
		cell.keysFound[cell.value] = true
	}

	// door is found
	if ch >= 'A' && ch <= 'Z' {
		keyForDoor := string(ch - 'A' + 'a')
		cell.keysNeeded[keyForDoor] = true
	}

	// push neighbors into queue
	for _, dir := range directions {
		neighborPos := pos.Add(dir)
		neighborCell := g.grid[neighborPos.Y][neighborPos.X]
		if neighborCell.seen || neighborCell.value == "#" {
			continue
		}

		// add to queue
		g.queue = append(g.queue, neighborPos)

		// update distance
		neighborCell.distance = cell.distance + 1

		// copy keysFound and keysNeeded
		for key := range cell.keysFound {
			neighborCell.keysFound[key] = true
		}
		for key := range cell.keysNeeded {
			neighborCell.keysNeeded[key] = true
		}
	}

	g.queue = g.queue[1:]
	isEmpty := len(g.queue) == 0
	if isEmpty {
		g.queue = nil
	}
	return isEmpty
}

func getKeysXY(grid [][]string) map[string]image.Point {
	keyCoords := make(map[string]image.Point)
	for y, rowSli := range grid {
		for x, cell := range rowSli {
			ch := int(cell[0])
			switch {
			case cell == "@":
				keyCoords["@"] = image.Point{x, y}
			case ch >= 'a' && ch <= 'z':
				keyCoords[cell] = image.Point{x, y}
			}
		}
	}
	return keyCoords
}

func initGrid(grid [][]string, src image.Point) *Grid {
	distanceGrid := make([][]*Node, len(grid))
	key := grid[src.Y][src.X]

	for y, row := range grid {
		distanceGrid[y] = make([]*Node, len(grid[0]))
		for x, cell := range row {
			distanceGrid[y][x] = &Node{
				value:      cell,
				distance:   math.MaxInt32,
				keysFound:  map[string]bool{key: true},
				keysNeeded: make(map[string]bool),
				seen:       false,
			}

			// "@" is not a key
			delete(distanceGrid[y][x].keysFound, "@")
		}
	}

	distanceGrid[src.Y][src.X].distance = 0
	distanceGrid[src.Y][src.X].seen = true
	queue := []image.Point{src}

	return &Grid{distanceGrid, queue}
}

type KeysGraph struct {
	keysToKeys map[string]map[string]*Node
}

func (graph *KeysGraph) addEdge(srcKey string, destKeyNode *Node) {
	if _, ok := graph.keysToKeys[srcKey]; !ok {
		graph.keysToKeys[srcKey] = map[string]*Node{}
	}

	destKey := destKeyNode.value
	graph.keysToKeys[srcKey][destKey] = destKeyNode
}

func (g *KeysGraph) dfsMinDistToAllKeys() int {
	var traverse func(string, map[string]bool) int
	cache := map[string]int{}

	keysToFind := []byte{}
	for i := 0; i < len(g.keysToKeys)-1; i++ {
		keysToFind = append(keysToFind, byte('a'+i))
	}

	traverse = func(srcChar string, keysFound map[string]bool) int {
		shortestDist := math.MaxInt32

		cacheKey := srcChar
		for _, key := range keysToFind {
			if !keysFound[string(key)] {
				cacheKey += string(key)
			}
		}

		if dist, found := cache[cacheKey]; found {
			return dist
		}

		// all keys found
		if len(keysFound) == len(keysToFind) {
			return 0
		}

		nextKeys := g.keysToKeys[srcChar]

		for destKey, destNode := range nextKeys {
			if keysFound[destKey] {
				continue
			}

			haveNeededKeys := true
			for neededKey := range destNode.keysNeeded {
				if !keysFound[neededKey] {
					haveNeededKeys = false
					break
				}
			}

			if !haveNeededKeys {
				continue
			}

			keysFound[destKey] = true

			distanceToEnd := destNode.distance + traverse(destKey, keysFound)
			if distanceToEnd < shortestDist {
				shortestDist = distanceToEnd
			}

			// backtrack
			delete(keysFound, destKey)
		}

		cache[cacheKey] = shortestDist
		return shortestDist
	}

	// start from entrance
	keysFound := map[string]bool{}
	return traverse("@", keysFound)
}

func parse(in io.Reader) [][]string {
	scanner := bufio.NewScanner(in)

	grid := make([][]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, strings.Split(line, ""))
	}

	return grid
}
