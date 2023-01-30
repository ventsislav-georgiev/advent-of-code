package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"time"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/synacor"
)

//go:embed challenge.bin
var challenge []byte

func main() {
	telecode := flag.String("telecode", "false", "calculate teleport code")
	vault := flag.String("vault", "false", "calculate vault path")
	flag.Parse()

	if *telecode == "true" {
		calcTeleportCode()
		return
	}

	if *vault == "true" {
		calcVaultPath()
		return
	}

	in := bytes.NewReader(challenge)
	vm := synacor.Parse(in)

	go func() {
		var cmdsPart1 = []string{
			"take tablet",
			"doorway",
			"north",
			"north",
			"bridge",
			"continue",
			"down",
			"east",
			"take empty lantern",
			"west",
			"west",
			"passage",
			"ladder",
			"west",
			"south",
			"north",
			"take can",
			"use can",
			"west",
			"ladder",
			"darkness",
			"use lantern",
			"continue",
			"west",
			"west",
			"west",
			"west",
			"north",
			"take red coin",
			"north",
			"east",
			"take concave coin",
			"down",
			"take corroded coin",
			"up",
			"west",
			"west",
			"take blue coin",
			"up",
			"take shiny coin",
			"down",
			"east",
			// 9+2*(5^2)+(7^3)-3 = 399
			"use blue coin",     // 9
			"use red coin",      // 2
			"use shiny coin",    // 5
			"use concave coin",  // 7
			"use corroded coin", // 3
			"north",
			"take teleporter",
			"use teleporter",
			"take business card",
			"take strange book",
		}

		for _, cmd := range cmdsPart1 {
			vm.InputLine(cmd)
		}

		time.Sleep(100 * time.Millisecond)
		fmt.Println("\n\nModifying r7...")
		vm.InputLine(":r7 25734")

		// 5483: set r0 4
		// 5486: set r1 1
		// 5489: call 6027
		// 5491: eq r1 r0 6
		// 5495: jf r1 5579
		vm.SetMem(5483+2, 6)
		vm.SetMem(5489, uint16(synacor.Noop))
		vm.SetMem(5490, uint16(synacor.Noop))

		vm.InputLine("use teleporter")

		var cmdsPart2 = []string{
			"north",
			"north",
			"north",
			"north",
			"north",
			"north",
			"north",
			"east",
			"take journal",
			"west",
			"north",
			"north",
			// vault
			"take orb",
			"north",
			"east",
			"east",
			"north",
			"west",
			"south",
			"east",
			"east",
			"west",
			"north",
			"north",
			"east",
			"vault",
			"take mirror",
			"use mirror",
		}

		for _, cmd := range cmdsPart2 {
			vm.InputLine(cmd)
		}
	}()

	vm.Run()
}

type CacheKey struct {
	reg0 uint16
	reg1 uint16
}

func calcTeleportCode() {
	fmt.Println("Calculating r7...")

	done := make(chan struct{}, 1)
	stepFn := func(min, max uint16) {
		if goto5489(min, max) {
			done <- struct{}{}
		}
	}

	for step := uint16(1); step <= 32768; step += 1024 {
		go stepFn(step, step+1024)
	}

	<-done
}

func goto5489(min, max uint16) bool {
	// 5483: set r0 4
	// 5486: set r1 1
	r0 := uint16(4)
	r1 := uint16(1)

	for r7 := min; r7 < max; r7++ {
		cache := make(map[CacheKey]uint16)

		// 5489: call 6027
		result := call6027(r0, r1, r7, cache)
		// 5491: eq r1 r0 6
		if result == 6 {
			fmt.Printf("done: %d\n", r7)
			return true
		}
	}

	return false
}

func call6027(reg0, reg1, reg7 uint16, cache map[CacheKey]uint16) (result uint16) {
	cacheKey := CacheKey{reg0, reg1}
	if result, ok := cache[cacheKey]; ok {
		return result
	}

	defer func() { cache[cacheKey] = result }()

	// 6027: jt r0(0) 6035
	if reg0 == 0 {
		// 6030: add r0(0) r1(153)  1
		// 6034: ret 6067
		reg0 = reg1 + 1
		return reg0
	}
	// 6035: jt r1(0) 6048
	if reg1 == 0 {
		// 6038: add r0(2) r0(2) 32767
		// 6042: set r1(0) r7(9)
		// 6045: call 6027
		reg0--
		reg1 = reg7
		return call6027(reg0, reg1, reg7, cache)
	}
	// 6048: push r0
	// 6050: add r1(9) r1(9) 32767
	reg1--
	// 6054: call 6027
	// 6067: ret 6056
	// 6056: set r1(153) r0(154)
	reg1 = call6027(reg0, reg1, reg7, cache)
	// 6059: pop r0(154)
	// 6061: add r0(1) r0(1) 32767
	// 6065: call 6027
	reg0--
	return call6027(reg0, reg1, reg7, cache)
}

func calcVaultPath() {
	//  *  8  -  1
	//  4  * 11  *
	//  +  4  - 18
	// 22  -  9  *
	nodes := []*Node{
		{Val: 22},
		{Val: 9},
		{Val: 4},
		{Val: 18},
		{Val: 4},
		{Val: 11},
		{Val: 8},
		{Val: 1},
	}

	nodes[0].Edges = []Edge{
		{OpSub, nodes[1], [2]string{"east", "east"}},
		{OpSub, nodes[2], [2]string{"east", "north"}},
		{OpAdd, nodes[2], [2]string{"north", "east"}},
		{OpAdd, nodes[4], [2]string{"north", "north"}},
	}
	nodes[1].Edges = []Edge{
		{OpSub, nodes[0], [2]string{"west", "west"}},
		{OpSub, nodes[2], [2]string{"west", "north"}},
		{OpSub, nodes[3], [2]string{"north", "east"}},
		{OpMul, nodes[3], [2]string{"east", "north"}},
		{OpSub, nodes[5], [2]string{"north", "north"}},
	}
	nodes[2].Edges = []Edge{
		{OpSub, nodes[0], [2]string{"south", "west"}},
		{OpAdd, nodes[0], [2]string{"west", "south"}},
		{OpSub, nodes[1], [2]string{"south", "east"}},
		{OpSub, nodes[3], [2]string{"east", "east"}},
		{OpAdd, nodes[4], [2]string{"west", "north"}},
		{OpMul, nodes[4], [2]string{"north", "west"}},
		{OpSub, nodes[5], [2]string{"east", "north"}},
		{OpMul, nodes[5], [2]string{"north", "east"}},
		{OpMul, nodes[6], [2]string{"north", "north"}},
	}
	nodes[3].Edges = []Edge{
		{OpMul, nodes[1], [2]string{"south", "west"}},
		{OpSub, nodes[1], [2]string{"west", "south"}},
		{OpSub, nodes[2], [2]string{"west", "west"}},
		{OpSub, nodes[5], [2]string{"west", "north"}},
		{OpMul, nodes[5], [2]string{"north", "west"}},
		{OpMul, nodes[7], [2]string{"north", "north"}},
	}
	nodes[4].Edges = []Edge{
		{OpAdd, nodes[0], [2]string{"south", "south"}},
		{OpAdd, nodes[2], [2]string{"south", "east"}},
		{OpMul, nodes[2], [2]string{"east", "south"}},
		{OpMul, nodes[5], [2]string{"east", "east"}},
		{OpMul, nodes[6], [2]string{"north", "east"}},
	}
	nodes[5].Edges = []Edge{
		{OpSub, nodes[2], [2]string{"south", "west"}},
		{OpMul, nodes[2], [2]string{"west", "south"}},
		{OpSub, nodes[1], [2]string{"south", "south"}},
		{OpSub, nodes[3], [2]string{"south", "east"}},
		{OpMul, nodes[3], [2]string{"east", "south"}},
		{OpMul, nodes[4], [2]string{"west", "west"}},
		{OpMul, nodes[6], [2]string{"west", "north"}},
		{OpSub, nodes[6], [2]string{"north", "west"}},
		{OpSub, nodes[7], [2]string{"north", "east"}},
		{OpMul, nodes[7], [2]string{"east", "north"}},
	}
	nodes[6].Edges = []Edge{
		{OpMul, nodes[4], [2]string{"west", "south"}},
		{OpMul, nodes[5], [2]string{"south", "east"}},
		{OpSub, nodes[5], [2]string{"east", "south"}},
		{OpMul, nodes[2], [2]string{"south", "south"}},
		{OpSub, nodes[7], [2]string{"east", "east"}},
	}
	nodes[7].Edges = []Edge{
		{OpMul, nodes[3], [2]string{"south", "south"}},
		{OpMul, nodes[5], [2]string{"south", "west"}},
		{OpSub, nodes[5], [2]string{"west", "south"}},
		{OpSub, nodes[6], [2]string{"west", "west"}},
	}

	path := Path{Val: 22, Node: nodes[0]}
	queue := []*Path{&path}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		if p.Node.Val == 1 && p.Val == 30 {
			fmt.Println("found path")
			exits := []string{}
			for p != nil {
				exits = append(exits, p.Exit[1], p.Exit[0])
				p = p.Prev
			}
			for i := len(exits) - 1; i >= 0; i-- {
				fmt.Printf("\"%s\",\n", exits[i])
			}
			fmt.Println()
			return
		}

		if p.Node.Val == 22 && p.Val != 22 {
			continue
		}

		if p.Node.Val == 1 && p.Val != 30 {
			continue
		}

		for _, edge := range p.Node.Edges {
			var val int
			switch edge.Operation {
			case OpAdd:
				val = p.Val + edge.Next.Val
			case OpSub:
				val = p.Val - edge.Next.Val
			case OpMul:
				val = p.Val * edge.Next.Val
			}

			if val < 0 || val > 100 {
				continue
			}

			newPath := Path{
				Val:  val,
				Node: edge.Next,
				Exit: edge.Exit,
				Prev: p,
			}
			queue = append(queue, &newPath)
		}
	}
}

const (
	OpAdd = iota
	OpSub
	OpMul
)

type Node struct {
	Val   int
	Edges []Edge
}

type Edge struct {
	Operation int
	Next      *Node
	Exit      [2]string
}

type Path struct {
	Val  int
	Node *Node
	Exit [2]string
	Prev *Path
}
