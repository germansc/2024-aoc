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

	// Part 1.
	part1 := blink(stones, 25)
	fmt.Printf("Part 1: %v\n", part1)

	// Part 2.
	part2 := blink(stones, 75)
	fmt.Printf("Part 2: %v\n", part2)
}

// Cache of the next set of stones for each one.
var cache = make(map[string][]string)

func newstones(stone string) []string {
	if r, ok := cache[stone]; ok {
		return r
	}

	buff := []string{}
	vi, _ := strconv.Atoi(stone)
	if vi == 0 {
		buff = append(buff, "1")
	} else if len(stone)%2 == 0 {
		// Convert intermediate strings to ints to remove leading zeros.
		v1, _ := strconv.Atoi(stone[:len(stone)/2])
		v2, _ := strconv.Atoi(stone[len(stone)/2:])

		buff = append(buff, strconv.Itoa(v1))
		buff = append(buff, strconv.Itoa(v2))
	} else {
		buff = append(buff, strconv.Itoa(vi*2024))
	}

	cache[stone] = buff
	return buff
}

func blink(stones []string, count int) int {
	// Process each uniq stone, keeping track of its count.
	stoneQueue := make(map[string]int)
	for _, stone := range stones {
		stoneQueue[stone]++
	}

	fmt.Printf("\nStarting set: %v\n", stoneQueue)
	for i := range count {
		// Prepare a queue for the next iteration.
		next := make(map[string]int)

		// Process each uniq stone in the current queue.
		for k, v := range stoneQueue {
			ns := newstones(k)
			// Append the resulting stones, for each copy of the current k.
			for _, s := range ns {
				next[s] += v
			}
		}

		// Set the next iterantion queue.
		stoneQueue = next

		// Check the final count of stones. I don't even care what they are,
		// only how many of each there are.
		total := 0
		for _, c := range stoneQueue {
			total += c
		}
		fmt.Printf("Blink %v: %v stones\n", i+1, total)
	}

	// Check the final count of stones. I don't even care what they are,
	// only how many of each there are.
	total := 0
	for _, c := range stoneQueue {
		total += c
	}

	return total
}
