package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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
	fmt.Printf("AoC 2024 - Day 21\n")
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

type robotstate struct {
	robot1, robot2, robot3 point
}

// Global grids and caches.
var (
	dcache = make(map[struct{ s, e point }]string)
	kcache = make(map[struct{ s, e point }]string)

	dpad = grid{
		height: 2,
		width:  3,
		data:   []byte(" ^A<v>"),
	}

	kpad = grid{
		height: 4,
		width:  3,
		data:   []byte("789456123 0A"),
	}
)

// Part solutions.
func solve(lines []string) {
	// Starting position for each robot
	ctx := robotstate{
		robot1: point{2, 3},
		robot2: point{2, 0},
		robot3: point{2, 0},
	}

	part1 := 0
	for _, code := range lines {
		seq := getSequence(&ctx, code)
		v, _ := strconv.Atoi(code[:len(code)-1])
		part1 += (v * len(seq))

		fmt.Println(code, ":", seq, v, len(seq))
	}

	fmt.Println("Part1:", part1)
}

func getSequence(state *robotstate, code string) string {
	var r1path, r2path, r3path string
	for _, k := range code {
		// Path to key.
		p := kpad.point(strings.IndexRune(string(kpad.data), k))
		r1path += bfs(&kpad, state.robot1, p, kcache)
		r1path += "A"
		state.robot1 = p
	}
	fmt.Println("Robot 1:", r1path)

	for _, k := range r1path {
		// Path to key.
		p := dpad.point(strings.IndexRune(string(dpad.data), k))
		r2path += bfs(&dpad, state.robot2, p, dcache)
		r2path += "A"
		state.robot2 = p
	}

	fmt.Println("Robot 2:", r2path)

	for _, k := range r2path {
		// Path to key.
		p := dpad.point(strings.IndexRune(string(dpad.data), k))
		r3path += bfs(&dpad, state.robot3, p, dcache)
		r3path += "A"
		state.robot3 = p
	}

	fmt.Println("Robot 3:", r3path)
	return r3path
}

// Bfs
func bfs(g *grid, start, end point, cache map[struct{ s, e point }]string) string {
	if v, ok := cache[struct{ s, e point }{start, end}]; ok {
		return v
	}

	seen := make(map[point]string)
	seen[start] = ""

	queue := []point{start}
	for len(queue) != 0 {
		p := queue[0]
		queue = queue[1:]
		if p == end {
			break
		}

		for k, diff := range []point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
			np := point{p.x + diff.x, p.y + diff.y}
			if !g.isValid(np) || g.data[g.index(np)] == ' ' {
				continue
			}

			if _, ok := seen[np]; ok {
				continue
			}

			move := "^>v<"
			queue = append(queue, np)
			seen[np] = seen[p] + string(move[k])
		}
	}

	cache[struct{ s, e point }{start, end}] = seen[end]
	return seen[end]
}
