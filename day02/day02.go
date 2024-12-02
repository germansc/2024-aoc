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
	fmt.Printf("AoC 2024 - Day 02\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.
func solve(lines []string) {
	var reports [][]int

	// Generate the int slices of each report.
	for _, line := range lines {
		s := strings.Fields(line)
		ints := make([]int, len(s))

		for i, val := range s {
			num, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			ints[i] = num
		}

		reports = append(reports, ints)
	}

	// Part 1:
	var safeCount int
	for _, rep := range reports {
		safe, _ := checkReport(rep)
		if safe {
			safeCount++
		}
	}

	fmt.Printf("Part 1: %v\n", safeCount)

	// Part 2:
	safeCount = len(reports)
	for _, rep := range reports {
		safe, i := checkReport(rep)
		if !safe {
			// Try removing one of the conflicting levels, previous, actual or next.
			tmp0 := removeAt(rep, i-1) // Maybe the conflicts comes from a previous condition
			tmp1 := removeAt(rep, i)
			tmp2 := removeAt(rep, i+1)

			safe0, _ := checkReport(tmp0)
			safe1, _ := checkReport(tmp1)
			safe2, _ := checkReport(tmp2)
			if !safe0 && !safe1 && !safe2 {
				safeCount--
			}
		}
	}

	fmt.Printf("Part 2: %v\n", safeCount)
}

// Removes an element from the slice
func removeAt(slice []int, index int) []int {
	if index < 0 || index >= len(slice) {
		return slice
	}

	result := make([]int, 0, len(slice)-1)
	result = append(result, slice[:index]...)
	result = append(result, slice[index+1:]...)

	return result
}

// Validates the report, and returns the index of the first level of an unsafe
// comparison.
func checkReport(report []int) (bool, int) {
	prev := report[0]
	way := (report[0] - report[1]) > 0

	for i := 1; i < len(report); i++ {
		diff := prev - report[i]
		if diff == 0 || diff < -3 || diff > 3 || (diff > 0) != way {
			return false, (i - 1)
		}
		prev = report[i]
	}

	return true, 0
}
