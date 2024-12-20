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
	fmt.Printf("AoC 2024 - Day 19\n")

	lines := readInput(input)

	solve(lines)
}

// Global caches.
var cache = make(map[string]bool)
var combinationCache = make(map[string]int)

// Part solutions.
func solve(lines []string) {
	towels := strings.Fields(strings.ReplaceAll(lines[0], ",", ""))
	requests := lines[2:]

	possible := []string{}
	for _, r := range requests {
		if isPossible(r, towels) {
			possible = append(possible, r)
		}
	}

	fmt.Println("Part 1:", len(possible))

	// Part 2:
	part2 := 0
	for _, r := range possible {
		comb := possibleCombinations(r, towels)
		fmt.Println(r, "has combinations:", comb)
		part2 += comb
	}

	fmt.Println("Part 2:", part2)
}

func isPossible(request string, towels []string) bool {
	if request == "" {
		return true
	}

	if v, ok := cache[request]; ok {
		return v
	}

	for _, towel := range towels {
		if strings.HasPrefix(request, towel) {
			if isPossible(request[len(towel):], towels) {
				cache[request] = true
				return true
			}
		}
	}

	cache[request] = false
	return false
}

func possibleCombinations(request string, towels []string) int {
	if request == "" {
		return 1
	}

	if v, ok := combinationCache[request]; ok {
		return v
	}

	comb := 0
	for _, towel := range towels {
		if strings.HasPrefix(request, towel) {
			comb += possibleCombinations(request[len(towel):], towels)
		}
	}

	combinationCache[request] = comb
	return comb
}
