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
	impossible := []equation{}
	for i := range len(eq) {
		if isPossible(eq[i]) {
			fmt.Printf("%-50s: %s\n", fmt.Sprintf("%v: %v", i, eq[i]), "POSSIBLE")
			part1 += eq[i].result
		} else {
			fmt.Printf("%-50s: %s\n", fmt.Sprintf("%v: %v", i, eq[i]), "NOT POSSIBLE")
			impossible = append(impossible, eq[i])
		}
	}

	fmt.Printf("Part 1: %v\n\n", part1)

	// Part 2.
	// Check if the previously impossible counts can be made possible with concatenation.
	part2 := part1
	for i := range len(impossible) {
		if isPossibleWithCocat(impossible[i]) {
			fmt.Printf("%-50s: %s\n", fmt.Sprintf("%v: %v", i, impossible[i]), "POSSIBLE WITH CONCAT")
			part2 += impossible[i].result
		} else {
			fmt.Printf("%-50s: %s\n", fmt.Sprintf("%v: %v", i, impossible[i]), "NOT POSSIBLE WITH CONCAT")
		}
	}

	fmt.Printf("Part 2: %v\n", part2)
}

// Recursive test of possible composition '*' and '+'
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

// Recursive test of possible composition of '+', '+', '||'... same as before,
// but add an additional branch to check the concatenations.
func isPossibleWithCocat(e equation) bool {
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

	if r == 0 && isPossibleWithCocat(equation{result: d, operands: e.operands[:l-1]}) {
		return true
	}

	// Concatenation in reverse only would be an option if the current result
	// "ends" with the current value 'n'. So, make a similar check to the previous one.
	trimmedStr := strings.TrimSuffix(strconv.Itoa(e.result), strconv.Itoa(n))
	unconcatted, _ := strconv.Atoi(trimmedStr)

	// if 'unconcatted' is the same as 'e.result', then it 'e.result' didn't
	// had 'n' as a suffix and we can skip this branch.
	if e.result != unconcatted && isPossibleWithCocat(equation{result: unconcatted, operands: e.operands[:l-1]}) {
		return true
	}

	return isPossibleWithCocat(equation{result: e.result - n, operands: e.operands[:l-1]})
}
