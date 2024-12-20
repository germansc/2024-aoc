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
	fmt.Printf("AoC 2024 - Day 20\n")
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

func (g grid) point(i int) point {
	return point{i % g.width, i / g.width}
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
		height: len(lines),
		width:  len(lines[0]),
		data:   []byte(strings.Join(lines, "")),
	}

	start, end := g.point(strings.IndexRune(string(g.data), 'S')), g.point(strings.IndexRune(string(g.data), 'E'))
	steps := bfs(&g, start, end)

	part1 := 0
	cheatpoints := cheatAnalysis(steps, &g)
	for _, v := range cheatpoints {
		if v >= 100 {
			part1++
		}
	}

	fmt.Println("Part 1:", part1)

	part2 := part1
	cheatpoints = moaaarCheats(steps, &g)
	for _, v := range cheatpoints {
		if v >= 100 {
			part2++
		}
	}
	fmt.Println("Part 2:", part2)
}

// Bfs
func bfs(g *grid, start, end point) map[point]int {
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

	return seen
}

func cheatAnalysis(steps map[point]int, g *grid) map[struct{ s, e point }]int {
	cheats := make(map[struct{ s, e point }]int)

	// For each point of the path, check for neighbours at distance 2, that
	// have biggest delta steps than 2.
	for k, si := range steps {
		for _, np := range manhattanPoints(k, 2) {
			if !g.isValid(np) {
				continue
			}

			if se, ok := steps[np]; !ok || (se-si) <= 2 {
				continue
			} else {
				cheats[struct{ s, e point }{k, np}] = se - si - 2
			}
		}
	}

	return cheats
}

func moaaarCheats(steps map[point]int, g *grid) map[struct{ s, e point }]int {
	cheats := make(map[struct{ s, e point }]int)

	// Range over all points at manhattan distances 3-20, and compute the
	// skipping distance. No need to test points at 2, since we already did for
	// part 1.
	for k, si := range steps {
		for dist := 3; dist < 21; dist++ {
			for _, np := range manhattanPoints(k, dist) {
				if !g.isValid(np) {
					continue
				}

				if se, ok := steps[np]; !ok || (se-si) <= dist {
					continue
				} else {
					cheats[struct{ s, e point }{k, np}] = se - si - dist
				}
			}
		}
	}

	return cheats
}

func manhattanPoints(p point, dist int) []point {
	points := []point{}

	// For each possible x-offset from -distance to +distance
	for dx := -dist; dx <= dist; dx++ {
		// The remaining distance must be used for the y-offset
		dy := dist - abs(dx)

		if dy >= 0 {
			points = append(points, point{p.x + dx, p.y + dy})
			if dy != 0 {
				points = append(points, point{p.x + dx, p.y - dy})
			}
		}
	}

	return points
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}
