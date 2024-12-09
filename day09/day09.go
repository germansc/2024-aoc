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

type block struct {
	start int
	size  int
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

	// For part 2.
	files := make(map[int]block)
	freeblocks := []block{}
	biggest := 0

	for k := range len(line) {
		blocks := int(line[k] - '0')
		total += blocks
		if k%2 != 0 {
			free += int(blocks)
			mem = append(mem, slices.Repeat[[]int]([]int{maxID}, blocks)...)
			freeblocks = append(freeblocks, block{cursor, blocks})
			if blocks > biggest {
				biggest = blocks
			}
		} else {
			mem = append(mem, slices.Repeat[[]int]([]int{id}, blocks)...)
			files[id] = block{start: cursor, size: blocks}
			id++
		}
		cursor += blocks
	}

	fmt.Printf("Total size: %v\n Free Size: %v\n", total, free)

	// Part 1
	// frag the disk, actually...
	mem1 := make([]int, len(mem))
	copy(mem1, mem)

	chksum := 0
	top := cursor - 1
	for k := range total - free {
		if mem1[k] != maxID {
			chksum += mem1[k] * k
			continue
		}

		for ; mem1[top] == maxID; top-- {
		}

		mem1[k] = mem1[top]
		mem1[top] = maxID
		chksum += mem1[k] * k
	}

	fmt.Printf("Part 1: %v\n", chksum)

	// Part 2
	// Relocate files from back to front, but only those that fit in previous
	// free blocks.
	mem2 := make([]int, len(mem))
	copy(mem2, mem)

	for id = maxID - 1; id > 0; id-- {
		fsize := files[id].size
		fstart := files[id].start
		if fsize > biggest {
			// No free-block of this size.
			continue
		}

		// Get the first block that can hold it.
		k := 0
		for ; k < len(freeblocks) && freeblocks[k].size < fsize; k++ {
		}

		// If there is no available block, or the available block comes after
		// this file, skip.
		if !(k < len(freeblocks)) || freeblocks[k].start > fstart {
			// Update the current biggest record
			if biggest > fsize {
				biggest = fsize - 1
			}
			continue
		}

		// Replace the block with the file.
		for jj := range fsize {
			mem2[freeblocks[k].start+jj] = id
			mem2[files[id].start+jj] = maxID
		}

		// update the freeblock record.
		if freeblocks[k].size == fsize {
			// remove the record.
			copy(freeblocks[k:], freeblocks[k+1:])
			freeblocks = freeblocks[:len(freeblocks)-1]
		} else {
			freeblocks[k].start += fsize
			freeblocks[k].size -= fsize
		}
	}

	chksum = 0
	for k, v := range mem2 {
		if v != maxID {
			chksum += v * k
		}
	}

	fmt.Printf("Part 2: %v\n", chksum)
}
