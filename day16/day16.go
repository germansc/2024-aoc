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

// Direction vectores, starting from East.
var dir = []point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func main() {
	fmt.Printf("AoC 2024 - Day 16\n")
	lines := readInput(input)
	solve(lines)
}

type point struct{ x, y int }

type charmap struct {
	data          []rune
	height, width int
}

func (c charmap) index(p point) int {
	return p.y*c.width + p.x
}

func (c charmap) point(i int) point {
	return point{i % c.width, i / c.width}
}

func (c charmap) isValid(p point) bool {
	if p.x < 0 || p.x > c.width-1 || p.y < 0 || p.y > c.height-1 {
		return false
	}

	return true
}

// Part solutions.
func solve(lines []string) {
	m := charmap{
		height: len(lines),
		width:  len(lines[0]),
		data:   []rune(strings.Join(lines, "")),
	}

	// Find the starting point.
	i := strings.IndexRune(string(m.data), 'S')
	start := m.point(i)

	// Part 1
	part1 := runmaze(&m, start)

	fmt.Println("Part 1:", part1)
}

func plotMap(m *charmap) {
	for y := range m.height {
		for x := range m.width {
			fmt.Printf("%c", m.data[m.index(point{x, y})])
		}
		fmt.Println()
	}
}

func runmaze(m *charmap, start point) int {
	cellscore := make(map[point]int)
	celldir := make(map[point]int)
	cellscore[start] = 0
	celldir[start] = 0

	queue := []point{start}

	for len(queue) != 0 {
		cp := queue[0]
		queue = queue[1:]

		// Get the next exploration points.
		for k := range 4 {
			// Skip 180° turns
			if k == 2 {
				continue
			}

			diff := dir[(celldir[cp]+k)%4]
			np := point{cp.x + diff.x, cp.y + diff.y}
			npscore := cellscore[cp] + 1 + 1000*(k%2)
			if m.data[m.index(np)] == '#' {
				continue
			}

			if cs, ok := cellscore[np]; !ok || npscore < cs {
				queue = append(queue, np)
				cellscore[np] = npscore
				celldir[np] = (celldir[cp] + k) % 4
			}
		}
	}

	// Check the best score at the exit point.
	i := strings.IndexRune(string(m.data), 'E')

	return cellscore[m.point(i)]
}
