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
	fmt.Printf("AoC 2024 - Day 10\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.

type point struct{ x, y int }

type heightmap struct {
	data          map[point]int
	height, width int
}

func (c heightmap) isValid(p point) bool {
	if p.x < 0 || p.x > c.width-1 || p.y < 0 || p.y > c.height-1 {
		return false
	}

	return true
}

func solve(lines []string) {
	hmap := heightmap{
		data:   make(map[point]int),
		height: len(lines),
		width:  len(lines[0]),
	}

	startpoints := []point{}
	for y, line := range lines {
		for x, v := range line {
			h := byte(v) - '0'
			hmap.data[point{x, y}] = int(h)
			if h == 0 {
				startpoints = append(startpoints, point{x, y})
			}
		}
	}

	// Part 1. Disctint uphills paths from 0 to 9.
	part1 := 0
	for _, s := range startpoints {
		pc := pathCount(&hmap, s)
		part1 += pc
		fmt.Printf("Point %v has a score of %v\n", s, pc)
	}

	fmt.Printf("Part 1: %v\n", part1)
}

func pathCount(h *heightmap, start point) int {
	count := 0
	queue := []point{start}
	visited := make(map[point]struct{})

	for len(queue) != 0 {
		p := queue[0]
		queue = queue[1:]

		// Skip visited cells.
		if _, ok := visited[p]; ok {
			continue
		}

		visited[p] = struct{}{}
		// Check if we reach the end.
		if h.data[p] == 9 {
			count++
		}

		// get the next valid neighbors
		next := []point{{p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y - 1}, {p.x, p.y + 1}}

		// Insert valid points to the queue.
		for _, np := range next {
			if h.isValid(np) && h.data[np] == h.data[p]+1 {
				queue = append(queue, np)
			}
		}
	}

	return count
}
