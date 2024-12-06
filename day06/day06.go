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
	fmt.Printf("AoC 2024 - Day 06\n")
	lines := readInput(input)
	solve(lines)
}

type direction struct {
	dx, dy int
}

func (d direction) String() string {
	switch d {
	case U:
		return "↑"
	case R:
		return "→"
	case D:
		return "↓"
	case L:
		return "←"
	default:
		return "o"
	}
}

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

var (
	U = direction{0, -1}
	D = direction{0, 1}
	L = direction{-1, 0}
	R = direction{1, 0}
)

func next(x, y int, d direction) (int, int) {
	return x + d.dx, y + d.dy
}

func printMap(cm charmap, visited map[int]direction) {
	for i := range cm.height * cm.width {
		if i != 0 && i%cm.height == 0 {
			fmt.Printf("\n")
		}

		if v, ok := visited[i]; ok {
			fmt.Printf("%v", v)
		} else {
			fmt.Printf("%c", cm.data[i])
		}
	}

	fmt.Printf("\n")
}

func rotate(d direction) direction {
	switch d {
	case U:
		return R
	case R:
		return D
	case D:
		return L
	case L:
		return U
	default:
		return direction{0, 0}
	}
}

// Checks for loops if starting from x,y in the given direction.
func runpath(m charmap, x, y int, d direction) (map[int][]direction, bool) {
	visited := make(map[int][]direction)

	// Run the path until exit the map. If reentered a visited map in a
	// previous direction, we got a loop.
	for {
		nx, ny := next(x, y, d)
		if !m.isValidPoint(nx, ny) {
			break
		}

		if m.data[m.idxFromPoint(nx, ny)] == '#' {
			d = rotate(d)
		} else {
			// Check if visited and mark it if not.
			if v, ok := visited[m.idxFromPoint(nx, ny)]; !ok {
				set := []direction{d}
				visited[m.idxFromPoint(nx, ny)] = set
			} else {
				for _, pd := range v {
					if pd == d {
						return visited, true
					}
				}
				v = append(v, d)
				visited[m.idxFromPoint(nx, ny)] = v
			}
			x, y = nx, ny
		}
	}

	return visited, false
}

// Part solutions.
func solve(lines []string) {
	cm := charmap{
		data:   []byte(strings.Join(lines, "")),
		height: len(lines),
		width:  len(lines[0]),
	}

	i := strings.IndexRune(string(cm.data), '^')
	fmt.Printf("  MAP: %v x %v\n", cm.width, cm.height)
	x, y := cm.pointFromIdx(i)
	currentdir := U

	fmt.Printf("Guard: (%v,%v)\n", x, y)
	visited, _ := runpath(cm, x, y, currentdir)

	fmt.Printf("Part 1: %v\n", len(visited))

	// Part 2
	// For each visited line, put a stone in the next position and run the rest
	// of the path to detect loops.
	possible := make(map[int]bool)

	x, y = cm.pointFromIdx(i)
	for k, v := range visited {
		if i == k {
			// Skip the starting point.
			continue
		}

		// Put a new rock.
		cm.data[k] = '#'

		// Only care about the first direction the cell is crossed.
		d := v[0]

		// run the path from the point before this new stone.
		x, y = cm.pointFromIdx(k)
		dp := rotate(rotate(d))
		x, y := next(x, y, dp)

		if _, loop := runpath(cm, x, y, d); loop {
			possible[k] = true
		}

		// Remove the rock.
		cm.data[k] = '.'
	}

	fmt.Printf("Part 2: %v\n", len(possible))
}
