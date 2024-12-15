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
	start := point{}
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
				start = point{x, y}
			}
		}
	}

	// For part 2
	m2 := expandMap(&m)

	// Part 1
	robot := start
	for _, d := range moves {
		robot = move(&m, robot, d)
	}

	part1 := 0
	for k, v := range m.data {
		if v == 'O' {
			part1 += 100*k.y + k.x
		}
	}

	plotMap(&m)
	fmt.Printf("Part 1: %d\n", part1)

	// Part 2
	robot = point{start.x * 2, start.y}
	for _, d := range moves {
		robot = move(&m2, robot, d)
	}

	part2 := 0
	for k, v := range m2.data {
		if v == '[' {
			part2 += 100*k.y + k.x
		}
	}

	plotMap(&m2)
	fmt.Printf("Part 2: %d\n", part2)
}

func expandMap(m *heightmap) heightmap {
	result := heightmap{
		height: m.height,
		width:  m.width * 2,
		data:   make(map[point]rune),
	}

	for k, v := range m.data {
		switch v {
		case '#':
			result.data[point{2 * k.x, k.y}] = '#'
			result.data[point{2*k.x + 1, k.y}] = '#'
		case 'O':
			result.data[point{2 * k.x, k.y}] = '['
			result.data[point{2*k.x + 1, k.y}] = ']'
		case '@':
			result.data[point{2 * k.x, k.y}] = '@'
			result.data[point{2*k.x + 1, k.y}] = '.'
		case '.':
			result.data[point{2 * k.x, k.y}] = '.'
			result.data[point{2*k.x + 1, k.y}] = '.'
		}
	}

	return result
}

func move(m *heightmap, p point, direction rune) point {
	canMove := true
	diff := point{}
	newVal := make(map[point]rune)
	vertical := false
	switch direction {
	case '<':
		diff = point{-1, 0}
	case '>':
		diff = point{1, 0}
	case '^':
		diff = point{0, -1}
		vertical = true
	case 'v':
		diff = point{0, +1}
		vertical = true
	}

	// Grow a tree al all tiles to be moved as we cross sides of boxes
	visited := make(map[point]bool)
	queue := []point{p}
	newVal[p] = '.'
	visited[p] = true

	for len(queue) != 0 {
		cp := queue[0]
		queue = queue[1:]

		// Get the next neighbour
		np := point{cp.x + diff.x, cp.y + diff.y}
		newVal[np] = m.data[cp]

		// Exit condition.
		if m.data[np] == '#' {
			canMove = false
			break
		}

		if m.data[np] == '.' {
			continue
		}

		// Add more items what to add to the queue.
		queue = append(queue, np)
		visited[np] = true

		// Skip the vertical draggin of wide boxes if the move is not vertical,
		// or the boxes are not wide.
		if !vertical {
			continue
		}

		if m.data[np] == '[' {
			sibbling := point{np.x + 1, np.y}
			if visited[sibbling] {
				continue
			}

			newVal[sibbling] = '.'
			queue = append(queue, sibbling)
		}

		if m.data[np] == ']' {
			sibbling := point{np.x - 1, np.y}
			if visited[sibbling] {
				continue
			}

			newVal[sibbling] = '.'
			queue = append(queue, sibbling)
		}
	}

	if !canMove {
		return p
	}

	// Update the map characters
	for k, v := range newVal {
		m.data[k] = v
	}

	return point{p.x + diff.x, p.y + diff.y}
}

func plotMap(m *heightmap) {
	for y := range m.height {
		for x := range m.width {
			fmt.Printf("%c", m.data[point{x, y}])
		}
		fmt.Println()
	}
}
