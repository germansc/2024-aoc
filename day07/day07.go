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
	fmt.Printf("AoC 2024 - Day 06\n")
	lines := readInput(input)
	solve(lines)
}

type equation struct {
	result   int
	operands []int
}

// Part solutions.
func solve(lines []string) {
	// Prepare the sets.
	eq := []equation{}
	for _, line := range lines {
		e := equation{}
		i := strings.IndexRune(line, ':')
		e.result, _ = strconv.Atoi(line[:i])
		for _, k := range strings.Fields(line[i+1:]) {
			v, _ := strconv.Atoi(k)
			e.operands = append(e.operands, v)
		}

		eq = append(eq, e)
	}

	part1 := 0
	for i := range len(eq) {
		if isPossible(eq[i]) {
			fmt.Printf("%-50s: %s\n", fmt.Sprintf("%v: %v", i, eq[i]), "POSSIBLE")
			part1 += eq[i].result
		} else {
			fmt.Printf("%-50s: %s\n", fmt.Sprintf("%v: %v", i, eq[i]), "NOT POSSIBLE")
		}
	}

	fmt.Printf("Part 1: %v\n", part1)
}

// Recursive test of possible composition
func isPossible(e equation) bool {
	// Final condition
	l := len(e.operands)

	// Get the currently evaluated number for convenience.
	if l == 0 {
		panic("No operands")
	}
	n := e.operands[l-1]

	if l == 1 {
		return n == e.result
	}

	// Get the int division and remainder of the current result and the first
	// operand
	d, r := e.result/n, e.result%n

	if r == 0 && isPossible(equation{result: d, operands: e.operands[:l-1]}) {
		return true
	}

	return isPossible(equation{result: e.result - n, operands: e.operands[:l-1]})
}
