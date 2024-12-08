package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	input = os.Stdin
)

type charmap struct {
	data   []byte
	height int
	width  int
}

func (m charmap) idxFromPoint(x, y int) int {
	if x >= m.width || y >= m.height {
		panic("")
	}

	idx := y*m.width + x
	if idx > m.width*m.height {
		fmt.Printf("WTF? (%v,%v) : %v\n", x, y, idx)
	}

	return y*m.width + x
}

func (m charmap) isValidPoint(x, y int) bool {
	return (x > -1 && x < m.width) && (y > -1 && y < m.height)
}

func (m charmap) pointFromIdx(i int) (int, int) {
	if i >= m.width*m.height {
		panic("")
	}

	return i % m.width, i / m.width
}

func readInput(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lines
}

func main() {
	fmt.Printf("AoC 2024 - Day 08\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.
func solve(lines []string) {
	m := charmap{
		data:   []byte(strings.Join(lines, "")),
		height: len(lines),
		width:  len(lines[0]),
	}

	fmt.Printf("Map %vx%v\n", m.width, m.height)

	nodes := make(map[int][]byte)
	antenas := make(map[byte][]int)

	// Get a list of all antennas positions.
	for y := range m.height {
		for x := range m.width {
			if c := m.data[m.idxFromPoint(x, y)]; c != '.' {
				antenas[c] = append(antenas[c], m.idxFromPoint(x, y))
			}
		}
	}

	// Iterate every pair of antenas and compute its antinodes.
	for c, points := range antenas {
		for i := range len(points) - 1 {
			for _, p2 := range points[i+1:] {
				p1 := points[i]
				x1, y1 := m.pointFromIdx(p1)
				x2, y2 := m.pointFromIdx(p2)

				dx := x2 - x1
				dy := y2 - y1

				// Check first antinode:
				if m.isValidPoint(x1-dx, y1-dy) {
					nodes[m.idxFromPoint(x1-dx, y1-dy)] = append(nodes[m.idxFromPoint(x1-dx, y1-dy)], c)
				}

				// Check second antinode:
				if m.isValidPoint(x2+dx, y2+dy) {
					nodes[m.idxFromPoint(x2+dx, y2+dy)] = append(nodes[m.idxFromPoint(x2+dx, y2+dy)], c)
				}
			}
		}
	}

	fmt.Printf("Part 1: %v\n", len(nodes))
}
