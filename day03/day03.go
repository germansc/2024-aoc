package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	input = os.Stdin

	// Regexp patterns
	regpart1 = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

	// This regexp matches a 'mul(xxx,xxx)', a 'do()' or a 'don't()'
	regpart2 = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)
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
	fmt.Printf("AoC 2024 - Day 03\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.
func solve(lines []string) {
	bulk := strings.Join(lines, "")
	part1 := 0

	// Get all mul(xxx,xxx) matches.
	matches := regpart1.FindAllStringSubmatch(bulk, -1)
	for _, match := range matches {
		part1 += score(match[1], match[2])
	}

	fmt.Printf("Part1: %d\n", part1)

	// Part 2
	part2 := 0

	// Get all matches, but only score those when the input is enabled.
	enabled := true
	matches = regpart2.FindAllStringSubmatch(bulk, -1)
	for _, match := range matches {
		switch match[0] {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			if enabled {
				part2 += score(match[1], match[2])
			}
		}
	}

	fmt.Printf("Part2: %d\n", part2)
}

func score(n1, n2 string) int {
	i1, _ := strconv.Atoi(string(n1))
	i2, _ := strconv.Atoi(string(n2))
	return i1 * i2
}
