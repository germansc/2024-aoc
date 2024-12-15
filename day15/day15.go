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
	fmt.Printf("AoC 2024 - Day 15\n")
	lines := readInput(input)
	solve(lines)
}

type point struct{ x, y int }

type heightmap struct {
	data          map[point]rune
	height, width int
}

func (c heightmap) isValid(p point) bool {
	if p.x < 0 || p.x > c.width-1 || p.y < 0 || p.y > c.height-1 {
		return false
	}

	return true
}

// Part solutions.
func solve(lines []string) {
	robot := point{}
	m := heightmap{
		width: len(lines[0]),
		data:  make(map[point]rune),
	}

	moves := ""
	for y, line := range lines {
		if line == "" {
			moves = strings.Join(lines[y+1:], "")
			m.height = y
			break
		}

		for x, v := range line {
			m.data[point{x, y}] = v
			if v == '@' {
				robot = point{x, y}
			}
		}
	}

	for _, d := range moves {
		robot = move(&m, robot, d)
	}

	// Part 1
	part1 := 0
	for k, v := range m.data {
		if v == 'O' {
			part1 += 100*k.y + k.x
		}
	}

	fmt.Printf("Part1: %d\n", part1)
}

func move(m *heightmap, p point, direction rune) point {
	cellsToSpace := 0
	switch direction {
	case '<':
		// Find an empty place to the left.
		for k := range p.x {
			if m.data[point{p.x - k, p.y}] == '#' {
				break
			}
			if m.data[point{p.x - k, p.y}] == '.' {
				cellsToSpace = k
				break
			}
		}
		if cellsToSpace == 0 {
			return p
		}

		// Shift all items in-between.
		for k := range cellsToSpace {
			m.data[point{p.x - cellsToSpace + k, p.y}] = m.data[point{p.x - cellsToSpace + k + 1, p.y}]
		}
		m.data[p] = '.'
		return point{p.x - 1, p.y}

	case '>':
		// Find an empty place to the left.
		for k := range m.width - p.x {
			if m.data[point{p.x + k, p.y}] == '#' {
				break
			}
			if m.data[point{p.x + k, p.y}] == '.' {
				cellsToSpace = k
				break
			}
		}
		if cellsToSpace == 0 {
			return p
		}

		// Shift all items in-between.
		for k := range cellsToSpace {
			m.data[point{p.x + cellsToSpace - k, p.y}] = m.data[point{p.x + cellsToSpace - k - 1, p.y}]
		}
		m.data[p] = '.'
		return point{p.x + 1, p.y}

	case 'v':
		// Find an empty place to the left.
		for k := range m.height - p.y {
			if m.data[point{p.x, p.y + k}] == '#' {
				break
			}
			if m.data[point{p.x, p.y + k}] == '.' {
				cellsToSpace = k
				break
			}
		}
		if cellsToSpace == 0 {
			return p
		}

		// Shift all items in-between.
		for k := range cellsToSpace {
			m.data[point{p.x, p.y + cellsToSpace - k}] = m.data[point{p.x, p.y + cellsToSpace - k - 1}]
		}
		m.data[p] = '.'
		return point{p.x, p.y + 1}

	case '^':
		// Find an empty place to the left.
		for k := range p.y {
			if m.data[point{p.x, p.y - k}] == '#' {
				break
			}
			if m.data[point{p.x, p.y - k}] == '.' {
				cellsToSpace = k
				break
			}
		}
		if cellsToSpace == 0 {
			return p
		}

		// Shift all items in-between.
		for k := range cellsToSpace {
			m.data[point{p.x, p.y - cellsToSpace + k}] = m.data[point{p.x, p.y - cellsToSpace + k + 1}]
		}

		m.data[p] = '.'
		return point{p.x, p.y - 1}
	}

	return point{}
}

func plotMap(m *heightmap) {
	for y := range m.height {
		for x := range m.width {
			fmt.Printf("%c", m.data[point{x, y}])
		}
		fmt.Println()
	}
}
