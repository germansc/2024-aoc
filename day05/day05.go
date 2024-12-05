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
	fmt.Printf("AoC 2024 - Day 05\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.
func solve(lines []string) {
	// Make a Map of al pages that must be printer after a specific one.
	mustPrior := make(map[int]map[int]struct{})
	var updates []string

	for k, line := range lines {
		if len(line) == 0 {
			updates = append(updates, lines[k+1:]...)
			break
		}

		var v1, v2 int
		if n, e := fmt.Sscanf(line, "%d|%d", &v1, &v2); n != 2 || e != nil {
			panic("Cant process rules!")
		}

		// make a ruleset for this value if it's new.
		if _, ok := mustPrior[v2]; !ok {
			mustPrior[v2] = make(map[int]struct{})
		}

		// Append the rule.
		mustPrior[v2][v1] = struct{}{}
	}

	// Process each value of the updates to see if a future value breaks a rule.
	var invalid [][]int
	part1 := 0
	for _, set := range updates {
		pages := makeIntArray(set)

		valid, _, _ := isValid(mustPrior, pages)
		if valid {
			part1 += pages[(len(pages)-1)/2]
		} else {
			invalid = append(invalid, pages)
		}
	}

	fmt.Printf("Part 1: %v\n", part1)

	// Part 2 // Could probably be made prettier with a recursive func.
	part2 := 0
	for _, set := range invalid {
		for valid, idx1, idx2 := isValid(mustPrior, set); !valid; valid, idx1, idx2 = isValid(mustPrior, set) {
			// Swap values, and recheck from idx1.
			v := set[idx2]
			set[idx2] = set[idx1]
			set[idx1] = v
		}

		// Now the set is valid.
		part2 += set[(len(set)-1)/2]
	}

	fmt.Printf("Part 2: %v\n", part2)
}

func makeIntArray(line string) []int {
	var result []int
	line = strings.TrimSpace(line)
	vals := strings.Split(line, ",")

	for _, v := range vals {
		v1, _ := strconv.Atoi(v)
		result = append(result, v1)
	}

	return result
}

func isValid(ruleset map[int]map[int]struct{}, pages []int) (bool, int, int) {
	valid := true
	var idx1, idx2 int
	for i := 0; i < len(pages)-1; i++ {
		p1 := pages[i]
		for j, p2 := range pages[i:] {
			if _, ok := ruleset[p1][p2]; ok {
				valid = false
				idx1, idx2 = i, i+j
				break
			}
		}
		if !valid {
			break
		}
	}

	return valid, idx1, idx2
}
