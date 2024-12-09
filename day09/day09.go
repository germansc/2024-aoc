package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
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
	fmt.Printf("AoC 2024 - Day 09\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.
func solve(lines []string) {
	// We'll work with only one line.
	line := lines[0]

	// Make an int value to represent free memory.
	maxID := len(line)/2 + 1

	mem := []int{}
	cursor := 0

	id := 0
	total := 0
	free := 0

	for k := range len(line) {
		blocks := int(line[k] - '0')
		total += blocks
		if k%2 != 0 {
			free += int(blocks)
			mem = append(mem, slices.Repeat[[]int]([]int{maxID}, blocks)...)
		} else {
			mem = append(mem, slices.Repeat[[]int]([]int{id}, blocks)...)
			id++
		}
		cursor += blocks
	}

	fmt.Printf("Total size: %v\n Free Size: %v\n", total, free)

	// frag the disk, actually...
	chksum := 0
	top := cursor - 1
	for k := range total - free {
		if mem[k] != maxID {
			chksum += mem[k] * k
			continue
		}

		for ; mem[top] == maxID; top-- {
		}

		mem[k] = mem[top]
		mem[top] = maxID
		chksum += mem[k] * k
	}

	fmt.Printf("Part 1: %v\n", chksum)
}
