package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
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
	fmt.Printf("AoC 2024 - Day 01\n")

	lines := readInput(input)

	solve(lines)
}

// Part solutions.

func solve(lines []string) {
	var list1, list2 []int

	for _, line := range lines {
		v := strings.Fields(line)
		v1, _ := strconv.Atoi(v[0])
		v2, _ := strconv.Atoi(v[1])

		list1 = append(list1, v1)
		list2 = append(list2, v2)
	}

	sort.Ints(list1)
	sort.Ints(list2)

	var distance int
	for k := range list1 {
		if list1[k] > list2[k] {
			distance += list1[k] - list2[k]
		} else {
			distance += list2[k] - list1[k]
		}
	}

	fmt.Printf("Part 1: %v\n", distance)

	// Part 2 - Reuse sorted lists.
	// Since the lists are sorted, we can compute a map of occurences by binary
	// searching the index of the next number.
	var uniq1 []int
	count1 := make(map[int]int)
	count2 := make(map[int]int)

	for i := 0; i < len(list1); {
		// Find the last occurrence of list1[i] using binary search
		end := sort.Search(len(list1), func(j int) bool { return list1[j] > list1[i] })
		count1[list1[i]] = end - i

		// Move to the next distinct value
		uniq1 = append(uniq1, list1[i])
		i = end
	}

	for i := 0; i < len(list2); {
		// Find the last occurrence of list2[i] using binary search
		end := sort.Search(len(list2), func(j int) bool { return list2[j] > list2[i] })
		count2[list2[i]] = end - i

		// Move to the next distinct value
		i = end
	}

	// Compute the similarity score:
	var simil int
	for _, item := range uniq1 {
		simil += item * count1[item] * count2[item]
	}

	fmt.Printf("Part 2: %v\n", simil)
}
