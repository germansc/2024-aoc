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
	fmt.Printf("AoC 2024 - Day 11\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.

func solve(lines []string) {
	// Input should be one line.
	if (len(lines) != 1) || len(lines[0]) == 0 {
		panic("Invalid input")
	}

	stones := []string{}
	for _, v := range strings.Fields(lines[0]) {
		stones = append(stones, v)
	}

	// Part 1. Process 25 blinks.
	part1 := blink(stones, 25)
	fmt.Printf("Part 1: %v\n", len(part1))

	part2 := 0
	fmt.Printf("Part 2: %v\n", part2)
}

func blink(stones []string, count int) []string {
	fmt.Printf("%v\n", stones)
	for range count {
		buff := []string{}
		for _, v := range stones {
			// Apply the rules.
			vi, _ := strconv.Atoi(v)
			if vi == 0 {
				buff = append(buff, "1")
			} else if len(v)%2 == 0 {
				// Convert intermediate strings to ints to remove leading zeros.
				v1, _ := strconv.Atoi(v[:len(v)/2])
				v2, _ := strconv.Atoi(v[len(v)/2:])

				buff = append(buff, strconv.Itoa(v1))
				buff = append(buff, strconv.Itoa(v2))
			} else {
				buff = append(buff, strconv.Itoa(vi*2024))
			}
		}

		// fmt.Printf("> %v\n", buff)
		stones = buff
	}

	return stones
}
