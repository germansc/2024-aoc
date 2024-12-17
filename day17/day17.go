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
	fmt.Printf("AoC 2024 - Day 17\n")
	lines := readInput(input)
	solve(lines)
}

// Part solutions.

type cpu struct {
	a, b, c int
	pc      int
	output  []byte
	mem     []int
}

func (c cpu) String() string {
	return fmt.Sprintf(
		"CPU STATUS:\nA = %d | B = %d | C = %d | PC = %d\nOUT = %q\n",
		c.a, c.b, c.c, c.pc, string(c.output))
}

func (c cpu) combo(v int) int {
	switch v {
	case 4:
		return c.a
	case 5:
		return c.b
	case 6:
		return c.c
	case 7:
		panic("combo 7?")
	default:
		return v
	}
}

func (c *cpu) adv() {
	operand := c.combo(c.mem[c.pc])
	c.pc++
	num := c.a
	den := 1 << operand
	c.a = num / den
}

func (c *cpu) bdv() {
	operand := c.combo(c.mem[c.pc])
	c.pc++
	num := c.a
	den := 1 << operand
	c.b = num / den
}

func (c *cpu) cdv() {
	operand := c.combo(c.mem[c.pc])
	c.pc++
	num := c.a
	den := 1 << operand
	c.c = num / den
}

func (c *cpu) bxl() {
	operand := c.mem[c.pc]
	c.pc++
	c.b ^= operand
}

func (c *cpu) bst() {
	operand := c.combo(c.mem[c.pc])
	c.pc++
	c.b = operand % 8
}

func (c *cpu) jnz() {
	operand := c.mem[c.pc]
	c.pc++
	if c.a != 0 {
		c.pc = operand
	}
}

func (c *cpu) bxc() {
	// Legacy
	c.pc++
	c.b ^= c.c
}

func (c *cpu) out() {
	operand := c.combo(c.mem[c.pc])
	c.pc++
	c.output = append(c.output, byte(operand%8))
}

func (c *cpu) RunProgram() {
	for c.pc < len(c.mem) {
		opcode := c.mem[c.pc]
		c.pc++

		// Run the op.
		switch opcode {
		case 0:
			c.adv()
		case 1:
			c.bxl()
		case 2:
			c.bst()
		case 3:
			c.jnz()
		case 4:
			c.bxc()
		case 5:
			c.out()
		case 6:
			c.bdv()
		case 7:
			c.cdv()

		default:
			panic("unk opcode")
		}
	}
}

func solve(lines []string) {
	core := cpu{}
	fmt.Sscanf(lines[0], "Register A: %d", &core.a)
	fmt.Sscanf(lines[1], "Register B: %d", &core.b)
	fmt.Sscanf(lines[2], "Register C: %d", &core.c)
	fmt.Sscanf(lines[4], "Program: %s", &core.output)
	for _, v := range strings.Split(string(core.output), ",") {
		i, _ := strconv.Atoi(v)
		core.mem = append(core.mem, i)
	}
	core.output = nil

	fmt.Println(core)
	fmt.Println(core.mem)

	core.RunProgram()
	fmt.Println(core)

	// Part 1:
	part1 := ""
	for _, v := range core.output {
		part1 += fmt.Sprintf(",%v", v)
	}
	if len(part1) > 0 {
		part1 = part1[1:]
	}

	fmt.Println("Part1: ", part1)
}
