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
	total := calcRoute(node, map[string]struct{}{}, 30, 0, Total{0, []string{}})
	fmt.Printf("Total: %d Path: %v\n", total.val, total.path)
}

func task2(in io.Reader) {
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

func calcRoute(node *Node, active map[string]struct{}, time, steps int, total Total) Total {
	timeLeft := (time - steps)
	if timeLeft <= 0 {
		return total
	}

	totals := make([]Total, 0, len(node.edges)*2)

	for _, edge := range node.edges {
		if _, ok := active[edge.node.name]; ok {
			continue
		}

		if edge.node.val == 0 {
			continue
		}

		activeCopy := make(map[string]struct{})
		for k, v := range active {
			activeCopy[k] = v
		}
		activeCopy[node.name] = struct{}{}

		newTimeLeft := timeLeft
		if node.val > 0 {
			newTimeLeft -= 1
		}

		newTotal := Total{total.val + (node.val * newTimeLeft), append(total.path, fmt.Sprintf("%s(%d)", node.name, newTimeLeft))}
		resultTotal := calcRoute(edge.node, activeCopy, newTimeLeft, edge.steps, newTotal)
		if resultTotal.val > 0 {
			totals = append(totals, resultTotal)
		}
	}

	if len(totals) == 0 {
		if _, ok := active[node.name]; !ok && node.val > 0 {
			total.val += node.val * (timeLeft - 1)
			total.path = append(total.path, fmt.Sprintf("%s(%d)", node.name, timeLeft-1))
		}

		return total
	}

	sort.Slice(totals, func(i, j int) bool {
		return totals[i].val < totals[j].val
	})

	return totals[len(totals)-1]
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
