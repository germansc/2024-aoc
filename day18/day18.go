package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	gridsize   = 71
	part1steps = 1024
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

type grid struct {
	data   []byte
	height int
	width  int
}

func (g grid) isValid(p point) bool {
	if p.x < 0 || p.y < 0 || p.x > g.width-1 || p.y > g.height-1 {
		return false
	}

	return true
}

func (g grid) index(p point) int {
	return p.x + p.y*g.width
}

func (g grid) plot() {
	for y := range g.height {
		for x := range g.width {
			char := g.data[g.index(point{x, y})]
			if char == 0 {
				char = '.'
			}
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
}

// Part solutions.
func solve(lines []string) {
	g := grid{
		height: gridsize,
		width:  gridsize,
		data:   make([]byte, gridsize*gridsize),
	}

	blocks := []point{}
	for _, line := range lines {
		p := point{}
		fmt.Sscanf(line, "%d,%d", &p.x, &p.y)
		blocks = append(blocks, p)
	}

	// Simulate the 1024 nanoseconds.
	for k := range part1steps {
		g.data[g.index(blocks[k])] = '#'
	}

	// Find the shortest path to the exit.
	part1 := bfs(&g, point{0, 0}, point{gridsize - 1, gridsize - 1})

	fmt.Println("Part 1:", part1)

	// Part 2. We know that at the 1024 byte, the path is still clear, we can
	// skip those first iteration, and make a binary search for the first map
	// where the end can't be reached.
	fmt.Println("Blocks:", len(blocks))
	blocks = blocks[part1steps:]
	fmt.Println("Blocks:", len(blocks))

	// Starting grid
	sg := grid{height: g.height, width: g.width, data: make([]byte, gridsize*gridsize)}
	copy(sg.data, g.data)

	last := false
	low, high, mid := 0, len(blocks)-1, 0
	for low <= high {
		mid = low + (high-low)/2

		// Simulate rocks up-to mid.
		copy(g.data, sg.data)
		fmt.Println("Trying iteration up to:", mid, blocks[mid])
		for k := range mid {
			g.data[g.index(blocks[k])] = '#'
		}

		if bfs(&g, point{0, 0}, point{gridsize - 1, gridsize - 1}) != 0 {
			fmt.Println("There's still a path at k:", mid)
			low = mid + 1
			last = true
		} else {
			fmt.Println("No path with k:", mid)
			high = mid - 1
			last = false
		}
	}

	// The last iteration might have a valid or invalid path. Use this to
	// determine which rock blocked the path.
	blocking := mid
	if last != true {
		blocking--
	}

	fmt.Println("Part 2:", blocks[blocking])
}

// Bfs
func bfs(g *grid, start, end point) int {
	seen := make(map[point]int)
	seen[start] = 0

	queue := []point{start}
	for len(queue) != 0 {
		p := queue[0]
		queue = queue[1:]
		if p == end {
			break
		}

		for _, diff := range []point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
			np := point{p.x + diff.x, p.y + diff.y}
			if !g.isValid(np) || g.data[g.index(np)] == '#' {
				continue
			}

			if _, ok := seen[np]; ok {
				continue
			}

			queue = append(queue, np)
			seen[np] = seen[p] + 1
		}
	}

	return seen[end]
}
