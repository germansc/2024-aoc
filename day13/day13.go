package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
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
	fmt.Printf("AoC 2024 - Day 13\n")

	lines := readInput(input)

	solve(lines)
}

// Part solutions.

type machine struct {
	x1, y1 int
	x2, y2 int
	xp, yp int
}

func solve(lines []string) {
	machines := []machine{}
	mach := machine{}
	for k, line := range lines {
		switch k % 4 {
		case 0:
			// Read first button.
			if n, _ := fmt.Sscanf(line, "Button A: X%d, Y%d", &mach.x1, &mach.y1); n != 2 {
				panic("Couldn't parse first line")
			}

		case 1:
			// Read second button.
			if n, _ := fmt.Sscanf(line, "Button B: X%d, Y%d", &mach.x2, &mach.y2); n != 2 {
				panic("Couldn't parse second line")
			}

		case 2:
			// Read first line.
			if n, _ := fmt.Sscanf(line, "Prize: X=%d, Y=%d", &mach.xp, &mach.yp); n != 2 {
				panic("Couldn't parse third line")
			}
			machines = append(machines, mach)
		}
	}

	part1 := 0
	part2 := 0
	for _, m := range machines {
		s, found := findcrosspoint(m)
		if found && s.a < 100 && s.b < 100 {
			part1 += 3*s.a + 1*s.b
		}

		m.xp += int(math.Pow10(13))
		m.yp += int(math.Pow10(13))
		s, found = findcrosspoint(m)
		if found {
			part2 += 3*s.a + 1*s.b
		}
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func findcrosspoint(m machine) (struct{ a, b int }, bool) {
	num := m.x1*m.yp - m.xp*m.y1
	den := m.y2*m.x1 - m.x2*m.y1
	if den == 0 || num%den != 0 {
		return struct{ a, b int }{0, 0}, false
	}

	b := num / den
	a, q := (m.xp-m.x2*b)/m.x1, (m.xp-m.x2*b)%m.x1

	if q != 0 {
		return struct{ a, b int }{0, 0}, false
	}

	return struct{ a, b int }{a, b}, true
}
