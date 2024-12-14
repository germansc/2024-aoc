package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	fmt.Printf("AoC 2024 - Day 14\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.

type point struct{ x, y int }

type robot struct {
	p point
	v point
}

func solve(lines []string) {
	width := 101
	height := 103
	robots := []robot{}

	// Get the input.
	for _, line := range lines {
		r := robot{}
		if n, _ := fmt.Sscanf(line, "p=%d,%d v=%d,%d", &r.p.x, &r.p.y, &r.v.x, &r.v.y); n != 4 {
			panic("Could not parse line!")
		}

		robots = append(robots, r)
	}

	// Simulate the 100 seconds.
	for range 100 {
		for k := range robots {
			robots[k].p.x = (robots[k].p.x + robots[k].v.x)
			robots[k].p.y = (robots[k].p.y + robots[k].v.y)
			if robots[k].p.x < 0 {
				robots[k].p.x += width
			}
			if robots[k].p.y < 0 {
				robots[k].p.y += height
			}
			if robots[k].p.x >= width {
				robots[k].p.x -= width
			}
			if robots[k].p.y >= height {
				robots[k].p.y -= height
			}
		}
	}

	m := make(map[point][]int)
	for k, r := range robots {
		m[r.p] = append(m[r.p], k)
	}

	plotMap(m, width, height)

	q1, q2, q3, q4 := 0, 0, 0, 0
	for p, count := range m {
		if p.x < (width-1)/2 {
			if p.y < (height-1)/2 {
				q1 += len(count)
			} else if p.y > (height-1)/2 {
				q3 += len(count)
			}
		} else if p.x > (width-1)/2 {
			if p.y < (height-1)/2 {
				q2 += len(count)
			} else if p.y > (height-1)/2 {
				q4 += len(count)
			}
		}
	}

	// Part 1.
	fmt.Printf("Part 1: %d\n", q1*q2*q3*q4)
}

func plotMap(m map[point][]int, w, h int) {
	fmt.Print("\033[H\033[2J")
	for y := range h {
		for x := range w {
			if id, ok := m[point{x, y}]; ok {
				fmt.Print(len(id))
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
