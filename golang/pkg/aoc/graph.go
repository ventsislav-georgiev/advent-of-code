package aoc

import (
	"image"
)

type Node interface {
	Pos() image.Point
	comparable
}

type Graph[T Node] struct {
	StartNode    T
	GetNextNodes func(T) []T
	GetWeight    func(T) int
	FilterNode   func(T) bool
}

type Path[T Node] struct {
	Node T
	Prev *Item[Path[T]]
	Dist int
}

func (g *Graph[T]) Dijkstra(destMatchPredicate func(a, b *Item[Path[T]]) bool) *Item[Path[T]] {
	start := Path[T]{Node: g.StartNode, Dist: 0}
	visited := map[T]struct{}{}
	itemsMap := make(map[T]*Item[Path[T]])
	queue := NewPriorityQueue[Path[T]]()

	getWeight := g.GetWeight
	if getWeight == nil {
		getWeight = func(_ T) int { return 1 }
	}

	enqueueOrUpdate := func(path Path[T]) {
		prevItem, found := itemsMap[path.Node]

		if found && prevItem.Val.Dist <= path.Dist {
			return
		}

		if found {
			prevItem.Val.Dist = path.Dist
			prevItem.Val.Prev = path.Prev
			queue.Update(prevItem)
			return
		}

		item := &Item[Path[T]]{
			Val: path,
			GetPriority: func() int {
				return path.Dist
			},
		}
		itemsMap[path.Node] = item
		queue.PushItem(item)
	}

	enqueueOrUpdate(start)

	for queue.Len() > 0 {
		item := queue.PopItem()
		path := item.Val

		for _, next := range g.GetNextNodes(path.Node) {
			if _, ok := visited[next]; ok {
				continue
			}

			visited[next] = struct{}{}

			if g.FilterNode != nil && !g.FilterNode(next) {
				continue
			}

			nextPath := Path[T]{
				Node: next,
				Prev: item,
				Dist: path.Dist + getWeight(next),
			}

			enqueueOrUpdate(nextPath)
		}
	}

	var bestPath *Item[Path[T]]
	for _, path := range itemsMap {
		if destMatchPredicate(bestPath, path) {
			bestPath = path
		}
	}

	return bestPath
}

type Pos struct {
	XY  image.Point
	Dir int
}

func (p Pos) Pos() image.Point {
	return p.XY
}

func (p Pos) MoveLR(dir int) Pos {
	switch p.Dir {
	case DirUp:
		switch dir {
		case DirLeft:
			return Pos{XY: p.XY.Add(MoveLeft), Dir: DirLeft}
		case DirRight:
			return Pos{XY: p.XY.Add(MoveRight), Dir: DirRight}
		}
	case DirDown:
		switch dir {
		case DirLeft:
			return Pos{XY: p.XY.Add(MoveRight), Dir: DirRight}
		case DirRight:
			return Pos{XY: p.XY.Add(MoveLeft), Dir: DirLeft}
		}
	case DirLeft:
		switch dir {
		case DirLeft:
			return Pos{XY: p.XY.Add(MoveDown), Dir: DirDown}
		case DirRight:
			return Pos{XY: p.XY.Add(MoveUp), Dir: DirUp}
		}
	case DirRight:
		switch dir {
		case DirLeft:
			return Pos{XY: p.XY.Add(MoveUp), Dir: DirUp}
		case DirRight:
			return Pos{XY: p.XY.Add(MoveDown), Dir: DirDown}
		}
	}

	panic("Only left and right directions are allowed")
}

func (p Pos) MoveForward() Pos {
	switch p.Dir {
	case DirUp:
		return Pos{XY: p.XY.Add(MoveUp), Dir: DirUp}
	case DirDown:
		return Pos{XY: p.XY.Add(MoveDown), Dir: DirDown}
	case DirLeft:
		return Pos{XY: p.XY.Add(MoveLeft), Dir: DirLeft}
	case DirRight:
		return Pos{XY: p.XY.Add(MoveRight), Dir: DirRight}
	}

	panic("Only up, down, left and right directions are allowed")
}
