package main

import (
	"bufio"
	"fmt"
	"io"
	"math"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanBytes)

	sources := parse(in)

	var minLocation uint = math.MaxUint32
	for _, seed := range sources.Seeds {
		location := sources.mapSeedToLocation(seed)
		if location < minLocation {
			minLocation = location
		}
	}

	fmt.Println(minLocation)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanBytes)

	sources := parse(in)

	var minLocation uint = math.MaxUint32
	seedsIdx := 0

	for {
		if len(sources.Seeds) < seedsIdx+1 {
			break
		}

		seeds := sources.Seeds[seedsIdx : seedsIdx+2]
		for seed := seeds[0]; seed < seeds[0]+seeds[1]; seed++ {
			location := sources.mapSeedToLocation(seed)
			if location < minLocation {
				minLocation = location
			}
		}

		seedsIdx += 2
	}

	fmt.Println(minLocation)
}

type Range struct {
	SrcMin, SrcMax uint
	Diff           int
}

type Sources struct {
	Seeds                                                                                        []uint
	SoilsMap, FertilizersMap, WatersMap, LightsMap, TemperaturesMap, HumiditiesMap, LocationsMap []Range
}

func (s *Sources) mapSeedToLocation(seed uint) uint {
	soil := mapSource(seed, s.SoilsMap)
	fertilizer := mapSource(soil, s.FertilizersMap)
	water := mapSource(fertilizer, s.WatersMap)
	light := mapSource(water, s.LightsMap)
	temperature := mapSource(light, s.TemperaturesMap)
	humidity := mapSource(temperature, s.HumiditiesMap)
	location := mapSource(humidity, s.LocationsMap)

	return location
}

func mapSource(src uint, mapping []Range) uint {
	for _, dst := range mapping {
		if src >= dst.SrcMin && src <= dst.SrcMax {
			return uint(int(src) + dst.Diff)
		}
	}

	return src
}

func parse(in io.Reader) (sources Sources) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanBytes)

	seeds := parseNumbers(scanner, '\n')
	skipLine(scanner)
	skipLine(scanner)

	sources.Seeds = seeds
	sources.SoilsMap = parseMap(scanner)
	sources.FertilizersMap = parseMap(scanner)
	sources.WatersMap = parseMap(scanner)
	sources.LightsMap = parseMap(scanner)
	sources.TemperaturesMap = parseMap(scanner)
	sources.HumiditiesMap = parseMap(scanner)
	sources.LocationsMap = parseMap(scanner)

	return sources
}

func parseNumbers(scanner *bufio.Scanner, term byte) []uint {
	var ch byte
	numbers := []uint{}
	numParts := []byte{}

	parseNumber := func() {
		if len(numParts) == 0 {
			return
		}
		numbers = append(numbers, aoc.Atoui(numParts))
		numParts = numParts[:0]
	}

	for ch != term {
		if !scanner.Scan() {
			parseNumber()
			break
		}

		ch = scanner.Bytes()[0]

		if ch >= '0' && ch <= '9' {
			numParts = append(numParts, ch)
		}

		if ch == ' ' || ch == term {
			parseNumber()
		}
	}

	return numbers
}

func parseMap(scanner *bufio.Scanner) []Range {
	mapping := [][]uint{}
	for {
		numbers := parseNumbers(scanner, '\n')
		if len(numbers) == 0 {
			skipLine(scanner)
			break
		}
		mapping = append(mapping, numbers)
	}

	return buildRanges(mapping)
}

func buildRanges(mapping [][]uint) []Range {
	ranges := []Range{}
	for _, numbers := range mapping {
		dstMin := numbers[0]
		srcMin := numbers[1]
		length := numbers[2]
		ranges = append(ranges, Range{
			SrcMin: srcMin,
			SrcMax: srcMin + length,
			Diff:   int(dstMin) - int(srcMin),
		})
	}
	return ranges
}

func skipLine(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if scanner.Bytes()[0] == '\n' {
			break
		}
	}
}
