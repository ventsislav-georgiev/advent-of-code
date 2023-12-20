package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	relays, _, _ := parse(in)
	button := relays["button"]

	var signalsStack []RelaySignal
	var lo, hi int

	for buttonPresses := 1000; buttonPresses > 0; buttonPresses-- {
		signalsStack = button.OnSignal(RelaySignal{})
		lo++

		for len(signalsStack) > 0 {
			src := signalsStack[0]
			signalsStack = signalsStack[1:]

			relay := relays[src.RelayID]
			if relay == nil {
				continue
			}

			newSignals := relay.OnSignal(src)

			for _, s := range newSignals {
				if s.Signal {
					hi++
				} else {
					lo++
				}
			}

			signalsStack = append(signalsStack, newSignals...)
		}
	}

	fmt.Println(lo * hi)
}

func task2(in io.Reader) {
	relays, conjuctionRelays, _ := parse(in)
	button := relays["button"]

	var signalsStack []RelaySignal
	var buttonPresses int
	var requiredPresses uint = 1

	keyRelayDetectedCycles := map[string]struct{}{}
	keyRelays := map[string]struct{}{}
	for id, relay := range conjuctionRelays {
		if len(relay.Outputs) >= 4 {
			keyRelays[id] = struct{}{}
		}
	}

cycle:
	for {
		buttonPresses++
		signalsStack = button.OnSignal(RelaySignal{})

		for len(signalsStack) > 0 {
			src := signalsStack[0]
			signalsStack = signalsStack[1:]

			if len(keyRelayDetectedCycles) == len(keyRelays) {
				break cycle
			}

			relay := relays[src.RelayID]
			if relay == nil {
				continue
			}

			newSignals := relay.OnSignal(src)

			if _, isKey := keyRelays[relay.ID]; isKey {
				isLoOutput := !newSignals[0].Signal
				if isLoOutput {
					if _, hasCycles := keyRelayDetectedCycles[relay.ID]; !hasCycles {
						requiredPresses = aoc.LeastCommonDenominator(requiredPresses, uint(buttonPresses))
						keyRelayDetectedCycles[relay.ID] = struct{}{}
					}
				}
			}

			signalsStack = append(signalsStack, newSignals...)
		}
	}

	fmt.Println(requiredPresses)
}

const (
	RelayTypeConjunction = iota
	RelayTypeFlip
	RelayTypeBroadcast
	RelayTypeButton
)

type Relay struct {
	ID           string
	Outputs      []string
	Type         byte
	Signal       bool
	SourceStates map[string]bool
}

type RelaySignal struct {
	SrcID   string
	RelayID string
	Signal  bool
}

func (r *Relay) OnSignal(src RelaySignal) (toSend []RelaySignal) {
	toSend = []RelaySignal{}

	switch r.Type {
	case RelayTypeButton:
		for _, id := range r.Outputs {
			toSend = append(toSend, RelaySignal{
				SrcID:   r.ID,
				RelayID: id,
				Signal:  false,
			})
		}
	case RelayTypeBroadcast:
		for _, id := range r.Outputs {
			toSend = append(toSend, RelaySignal{
				SrcID:   r.ID,
				RelayID: id,
				Signal:  src.Signal,
			})
		}
	case RelayTypeConjunction:
		r.SourceStates[src.SrcID] = src.Signal
		useLoSignal := true
		for _, state := range r.SourceStates {
			if !state {
				useLoSignal = false
				break
			}
		}

		r.Signal = !useLoSignal
		for _, id := range r.Outputs {
			toSend = append(toSend, RelaySignal{
				SrcID:   r.ID,
				RelayID: id,
				Signal:  r.Signal,
			})
		}
	case RelayTypeFlip:
		if src.Signal {
			break
		}

		r.Signal = !r.Signal
		for _, id := range r.Outputs {
			toSend = append(toSend, RelaySignal{
				SrcID:   r.ID,
				RelayID: id,
				Signal:  r.Signal,
			})
		}
	}

	return
}

func parse(in io.Reader) (relays map[string]*Relay, conjuctionRelays map[string]*Relay, flipRelays map[string]*Relay) {
	relays = map[string]*Relay{}
	conjuctionRelays = map[string]*Relay{}
	flipRelays = map[string]*Relay{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " -> ")
		id := parts[0]
		relay := &Relay{}
		relay.Outputs = strings.Split(parts[1], ", ")

		switch id[0] {
		case '&':
			id = id[1:]
			relay.Signal = true
			relay.Type = RelayTypeConjunction
			relay.SourceStates = map[string]bool{}
			conjuctionRelays[id] = relay
		case '%':
			id = id[1:]
			relay.Type = RelayTypeFlip
			flipRelays[id] = relay
		default:
			relay.Type = RelayTypeBroadcast
		}

		relay.ID = id
		relays[id] = relay
	}

	for cID, conj := range conjuctionRelays {
		for id, relay := range relays {
			if aoc.ListContains(relay.Outputs, cID) {
				conj.SourceStates[id] = false
			}
		}
	}

	buttonRelay := &Relay{
		Type:    RelayTypeButton,
		Outputs: []string{"broadcaster"},
	}
	relays["button"] = buttonRelay

	return
}
