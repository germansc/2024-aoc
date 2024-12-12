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
	fmt.Printf("AoC 2024 - Day 12\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.

type point struct{ x, y int }

type charmap struct {
	data          []byte
	height, width int
}

func (c charmap) index(p point) int {
	return p.y*c.width + p.x
}

func (c charmap) isValid(p point) bool {
	if p.x < 0 || p.x > c.width-1 || p.y < 0 || p.y > c.height-1 {
		return false
	}

	return true
}

type plot struct {
	crop      byte
	cells     map[point]bool
	neighbors map[point]bool
}

func (p plot) perimeter() int {
	per := 0

	for k := range p.cells {
		per += 4
		neigh := []point{{k.x - 1, k.y}, {k.x + 1, k.y}, {k.x, k.y - 1}, {k.x, k.y + 1}}
		for _, n := range neigh {
			if _, ok := p.cells[n]; ok {
				per--
			}
		}
	}

	return per
}

func (p plot) corners() int {
	corners := 0

	// Num of sides = num of corners.
	for k := range p.cells {
		// Gather info about the 8 surrounding cells to count corners.
		_, ul := p.cells[point{k.x - 1, k.y - 1}]
		_, u := p.cells[point{k.x, k.y - 1}]
		_, ur := p.cells[point{k.x + 1, k.y - 1}]
		_, l := p.cells[point{k.x - 1, k.y}]
		_, r := p.cells[point{k.x + 1, k.y}]
		_, dl := p.cells[point{k.x - 1, k.y + 1}]
		_, d := p.cells[point{k.x, k.y + 1}]
		_, dr := p.cells[point{k.x + 1, k.y + 1}]

		// Top-left corner:
		if (l == u) && (!l || !ul) {
			corners++
		}

		// Top-right corner:
		if (r == u) && (!r || !ur) {
			corners++
		}

		// Bot-right corner:
		if (r == d) && (!r || !dr) {
			corners++
		}

		// bot-left corner:
		if (l == d) && (!l || !dl) {
			corners++
		}
	}

	return corners
}

func (p plot) String() string {
	return fmt.Sprintf("Region of %c plants | area: %v | perimeter: %v | corners: %v", p.crop, len(p.cells), p.perimeter(), p.corners())
}

func newPlot(start point, cm *charmap) plot {
	result := plot{
		crop:      cm.data[cm.index(start)],
		neighbors: make(map[point]bool),
		cells:     make(map[point]bool),
	}

	checked := make(map[point]bool)
	checked[start] = true

	queue := []point{start}
	for len(queue) != 0 {
		p := queue[0]
		queue = queue[1:]

		// Mark as a valid cell in the plot.
		result.cells[p] = true

		// Check neighbors
		neigh := []point{{p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y - 1}, {p.x, p.y + 1}}
		for _, np := range neigh {
			if !cm.isValid(np) || result.cells[np] || checked[np] {
				continue
			}

			if cm.data[cm.index(np)] != result.crop {
				checked[np] = true
				result.neighbors[np] = true
			} else {
				checked[np] = true
				queue = append(queue, np)
			}
		}
	}

	return result
}

func solve(lines []string) {
	cmap := charmap{
		data:   []byte(strings.Join(lines, "")),
		height: len(lines),
		width:  len(lines[0]),
	}

	pointToPlot := make(map[point]*plot)
	plots := []plot{}

	// Gather info on each plot.
	queue := []point{{0, 0}}
	for len(queue) != 0 {
		p := queue[0]
		queue = queue[1:]

		if _, ok := pointToPlot[p]; !ok {
			// Scan plot. Mark the cells as visited and the neighbors to the queue.
			plot := newPlot(p, &cmap)
			for k := range plot.cells {
				pointToPlot[k] = &plot
			}

			for k := range plot.neighbors {
				queue = append(queue, k)
			}

			plots = append(plots, plot)
			fmt.Println(plot)
		}
	}

	part1 := 0
	part2 := 0
	for _, p := range plots {
		part1 += len(p.cells) * p.perimeter()
		part2 += len(p.cells) * p.corners()
	}

	fmt.Printf("Part 1: %v\n", part1)
	fmt.Printf("Part 2: %v\n", part2)
}
