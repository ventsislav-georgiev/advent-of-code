package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	graph := parse(in)
	updateEdges(graph)
	node := graph["AA"]
	totals := calcRoute(node, map[string]struct{}{}, 30, 0, Total{0, []string{}}, false)
	total := totals[0]
	fmt.Printf("Total: %d Path: %v\n", total.val, total.path)
}

func task2(in io.Reader) {
	graph := parse(in)
	updateEdges(graph)
	node := graph["AA"]
	totals := calcRoute(node, map[string]struct{}{}, 26, 0, Total{0, []string{}}, true)

	// reduce paths to speed up the calc
	tmptotals := []Total{}
	for _, total := range totals {
		if total.val > 1200 { // magic number found by trial and error (any less is slow, any more is giving 0 as result)
			tmptotals = append(tmptotals, total)
		}
	}
	totals = tmptotals

	max := 0
	indexes := [2]int{}
	for idx := 0; idx < len(totals); idx++ {
		for jdx := idx + 1; jdx < len(totals); jdx++ {
			use := true
			p1 := totals[idx].path
			p2 := totals[jdx].path

		test:
			for i := 0; i < len(p1); i++ {
				for j := 0; j < len(p2); j++ {
					if p1[i] == p2[j] && p1[i] != "AA" {
						use = false
						break test
					}
				}
			}

			if m := totals[idx].val + totals[jdx].val; use && m > max {
				max = m
				indexes[0] = idx
				indexes[1] = jdx
			}
		}
	}

	paths := []string{}
	paths = append(paths, totals[indexes[0]].path...)
	paths = append(paths, totals[indexes[1]].path...)
	fmt.Println(max, paths)
}

type Node struct {
	name  string
	val   int
	edges map[string]*Edge
}

type Edge struct {
	node  *Node
	steps int
}

type Total struct {
	val  int
	path []string
}

func calcRoute(node *Node, active map[string]struct{}, time, steps int, total Total, returnAll bool) []Total {
	time = (time - steps)
	if time <= 0 {
		return []Total{total}
	}

	totals := make([]Total, 0, len(node.edges)*2)

	for _, edge := range node.edges {
		if _, ok := active[edge.node.name]; ok || edge.node.val == 0 {
			continue
		}

		activecp := make(map[string]struct{})
		for k, v := range active {
			activecp[k] = v
		}
		activecp[node.name] = struct{}{}

		newtime := time
		if node.val > 0 {
			newtime -= 1
		}

		newtotal := Total{total.val + (node.val * newtime), append(total.path, node.name)}
		totals = append(totals, calcRoute(edge.node, activecp, newtime, edge.steps, newtotal, returnAll)...)
	}

	if len(totals) == 0 {
		if _, ok := active[node.name]; !ok && node.val > 0 {
			total.val += node.val * (time - 1)
			total.path = append(total.path, node.name)
		}

		return []Total{total}
	}

	if !returnAll {
		sort.Slice(totals, func(i, j int) bool {
			return totals[i].val > totals[j].val
		})
		return []Total{totals[0]}
	}

	return totals
}

func stepsToNode(src *Node, dest *Node) int {
	visited := map[string]bool{}
	q := []*Edge{{src, 0}}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		if cur.node == dest {
			return cur.steps
		}

		for _, edge := range cur.node.edges {
			if !visited[edge.node.name] {
				visited[edge.node.name] = true
				adjedge := &Edge{edge.node, cur.steps + edge.steps}
				q = append(q, adjedge)
			}
		}
	}

	return -1
}

func updateEdges(graph map[string]*Node) {
	nodes := []*Node{}
	for _, node := range graph {
		nodes = append(nodes, node)
	}

	edges := []map[string]*Edge{}
	for i := 0; i < len(nodes); i++ {
		edges = append(edges, map[string]*Edge{})
	}

	for i := 0; i < len(nodes); i++ {
		src := nodes[i]
		edges = append(edges, map[string]*Edge{})

		for j := i + 1; j < len(nodes); j++ {
			dest := nodes[j]
			if src.edges[dest.name] != nil {
				continue
			}

			steps := stepsToNode(src, dest)
			edges[i][dest.name] = &Edge{steps: steps, node: dest}
			edges[j][src.name] = &Edge{steps: steps, node: src}
		}
	}

	for i := 0; i < len(nodes); i++ {
		for name, edge := range edges[i] {
			nodes[i].edges[name] = edge
		}
	}
}

func parse(in io.Reader) map[string]*Node {
	graph := make(map[string]*Node)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.NewReplacer(", ", ",", "tunnels", "tunnel", "leads", "lead", "valves", "valve").Replace(scanner.Text())
		var srcName, connections string
		var flowRate int
		fmt.Sscanf(line, "Valve %s has flow rate=%d; tunnel lead to valve %s", &srcName, &flowRate, &connections)

		src, exists := graph[srcName]
		if !exists {
			src = &Node{name: srcName, val: 0, edges: map[string]*Edge{}}
			graph[srcName] = src
		}
		src.val = flowRate

		for _, destName := range strings.Split(connections, ",") {
			dest, exists := graph[destName]
			if !exists {
				dest = &Node{name: destName, val: 0, edges: map[string]*Edge{}}
				graph[destName] = dest
			}

			src.edges[destName] = &Edge{steps: 1, node: dest}
			dest.edges[srcName] = &Edge{steps: 1, node: src}
		}
	}

	return graph
}
