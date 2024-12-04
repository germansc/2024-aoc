package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	U = iota
	D
	L
	R
	UL
	UR
	DL
	DR
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
	fmt.Printf("AoC 2024 - Day 04\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.
func solve(lines []string) {
	width := len(lines[0])
	height := len(lines)

	part1 := part1(lines, width, height)
	part2 := part2(lines, width, height)

	fmt.Printf("Part 1: %v\n", part1)
	fmt.Printf("Part 2: %v\n", part2)
}

func part1(lines []string, width, height int) int {
	result := 0
	for y := range height {
		for x := range width {
			if lines[y][x] == 'X' {
				if extractWord(lines, width, height, x, y, 4, U) == "XMAS" {
					result++
				}
				if extractWord(lines, width, height, x, y, 4, D) == "XMAS" {
					result++
				}
				if extractWord(lines, width, height, x, y, 4, L) == "XMAS" {
					result++
				}
				if extractWord(lines, width, height, x, y, 4, R) == "XMAS" {
					result++
				}
				if extractWord(lines, width, height, x, y, 4, UL) == "XMAS" {
					result++
				}
				if extractWord(lines, width, height, x, y, 4, UR) == "XMAS" {
					result++
				}
				if extractWord(lines, width, height, x, y, 4, DL) == "XMAS" {
					result++
				}
				if extractWord(lines, width, height, x, y, 4, DR) == "XMAS" {
					result++
				}
			}
		}
	}

	return result
}

func part2(lines []string, width, height int) int {
	result := 0
	for y := range height {
		for x := range width {
			if lines[y][x] == 'A' {
				w1 := extractWord(lines, width, height, x-1, y+1, 3, UR)
				w2 := extractWord(lines, width, height, x-1, y-1, 3, DR)

				if (w1 == "MAS" || w1 == "SAM") && (w2 == "MAS" || w2 == "SAM") {
					result++
				}
			}
		}
	}

	return result
}

func extractWord(data []string, w, h, x, y int, length int, dir int) string {
	if x < 0 || x >= w || y < 0 || y >= h {
		return ""
	}

	result := make([]byte, 0, length)
	switch dir {
	case U:
		if y-(length-1) < 0 {
			return ""
		}
		for i := range length {
			result = append(result, data[y-i][x])
		}

	case D:
		if y+(length-1) >= h {
			return ""
		}
		for i := range length {
			result = append(result, data[y+i][x])
		}

	case L:
		if x-(length-1) < 0 {
			return ""
		}
		for i := range length {
			result = append(result, data[y][x-i])
		}

	case R:
		if x+(length-1) >= w {
			return ""
		}
		for i := range length {
			result = append(result, data[y][x+i])
		}

	case UL:
		if y-(length-1) < 0 || x-(length-1) < 0 {
			return ""
		}
		for i := range length {
			result = append(result, data[y-i][x-i])
		}

	case UR:
		if y-(length-1) < 0 || x+(length-1) >= w {
			return ""
		}
		for i := range length {
			result = append(result, data[y-i][x+i])
		}

	case DL:
		if y+(length-1) >= h || x-(length-1) < 0 {
			return ""
		}
		for i := range length {
			result = append(result, data[y+i][x-i])
		}

	case DR:
		if y+(length-1) >= h || x+(length-1) >= w {
			return ""
		}
		for i := range length {
			result = append(result, data[y+i][x+i])
		}
	}

	return string(result)
}
