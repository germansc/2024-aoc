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
	regxmul   = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	regxdo    = regexp.MustCompile(`do`)
	regxdonot = regexp.MustCompile(`don't`)
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
	part1 := 0

	for _, line := range lines {
		part1 += score(line)
	}

	fmt.Printf("Part1: %d\n", part1)

	// Part 2
	// Join the lines in case there are do's and don't between lines
	bulk := strings.Join(lines, "")

	// Find the indexes of all keywords
	var dos []int
	var donts []int
	d1 := regxdo.FindAllStringIndex(bulk, -1)
	d2 := regxdonot.FindAllStringIndex(bulk, -1)

	// Get the starting index of each set.
	for _, d := range d1 {
		dos = append(dos, d[0])
	}

	for _, d := range d2 {
		donts = append(donts, d[0])
	}

	// Discards don'ts from dos
	dos = filter(dos, donts)

	// Discard disabled sections.
	ie := 0
	id := 0
	enabled := ""

	for {
		// Find the next index of a disabled section.
		i := 0
		for ; i < len(donts); i++ {
			if donts[i] > ie {
				break
			}
		}

		if i < len(donts) {
			id = donts[i]
			donts = donts[i:]
		} else {
			id = len(bulk)
		}

		// Append enabled data.
		enabled += bulk[ie:id]

		i = 0
		for ; i < len(dos); i++ {
			if dos[i] > id {
				break
			}
		}

		if i < len(dos) {
			ie = dos[i]
			dos = dos[i:]
		} else {
			break
		}
	}

	part2 := score(enabled)
	fmt.Printf("Part2: %d\n", part2)
}

func score(line string) int {
	score := 0
	// Get all matches.
	matches := regxmul.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		n1, _ := strconv.Atoi(string(match[1]))
		n2, _ := strconv.Atoi(string(match[2]))
		score += n1 * n2
	}
	return score
}

func filter(data, remove []int) []int {
	remmap := make(map[int]struct{})
	for _, i := range remove {
		remmap[i] = struct{}{}
	}

	// Filter the original slice
	result := []int{}
	for _, i := range data {
		if _, found := remmap[i]; !found {
			result = append(result, i)
		}
	}

	return result
}
