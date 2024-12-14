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

	startingRobots := make([]robot, len(robots))
	copy(startingRobots, robots)

	// Simulate the 100 seconds.
	for range 100 {
		for k := range robots {
			moveRobot(&robots[k], width, height)
		}
	}

	m := make(map[point]int)
	for _, r := range robots {
		m[r.p]++
	}

	q1, q2, q3, q4 := 0, 0, 0, 0
	for p, count := range m {
		if p.x < (width-1)/2 {
			if p.y < (height-1)/2 {
				q1 += count
			} else if p.y > (height-1)/2 {
				q3 += count
			}
		} else if p.x > (width-1)/2 {
			if p.y < (height-1)/2 {
				q2 += count
			} else if p.y > (height-1)/2 {
				q4 += count
			}
		}
	}

	// Part 2.
	// Positions of robots repeat after 10403 iterations.
	// Search for the time of least overlaps? assuming that most robots will be
	// forming the tree?
	least := len(startingRobots)
	frame := make(map[point]int)
	leastTime := 0
	for time := range 10403 {
		m2 := make(map[point]int)
		for k := range startingRobots {
			moveRobot(&startingRobots[k], width, height)
			m2[startingRobots[k].p]++
		}

		// Check if all robots are in different positions.
		if overlapping := len(startingRobots) - len(m2); overlapping < least {
			least = overlapping
			leastTime = time + 1

			// Clear the previous record
			for k := range frame {
				delete(frame, k)
			}

			// make a copy of the current frame.
			for k, v := range m2 {
				frame[k] = v
			}
		}
	}

	// Ok, this actually worked. And there was even no overlap in that frame,
	// although not all robots were part of the tree. Coincidence?
	plotMap(frame, width, height)

	// Print results.
	fmt.Printf("Part 1: %d\n", q1*q2*q3*q4)
	fmt.Printf("Part 2: The frame with the least overlap (%v) occurred after %v seconds.\n", least, leastTime)
}

func plotMap(m map[point]int, w, h int) {
	fmt.Print("\033[H\033[2J")
	for y := range h {
		for x := range w {
			if id, ok := m[point{x, y}]; ok {
				fmt.Print(id)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func moveRobot(r *robot, width, height int) {
	r.p.x = (r.p.x + r.v.x)
	r.p.y = (r.p.y + r.v.y)

	if r.p.x < 0 {
		r.p.x += width
	}

	if r.p.y < 0 {
		r.p.y += height
	}

	if r.p.x >= width {
		r.p.x -= width
	}

	if r.p.y >= height {
		r.p.y -= height
	}
}
