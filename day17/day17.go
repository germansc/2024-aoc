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
	output  []int
	mem     []int
}

func (c cpu) String() string {
	return fmt.Sprintf(
		"CPU STATUS:\nA = %d | B = %d | C = %d | PC = %d\nOUT = %v\n",
		c.a, c.b, c.c, c.pc, c.output)
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
	c.a = c.a >> operand
}

func (c *cpu) bdv() {
	operand := c.combo(c.mem[c.pc])
	c.pc++
	c.b = c.a >> operand
}

func (c *cpu) cdv() {
	operand := c.combo(c.mem[c.pc])
	c.pc++
	c.c = c.a >> operand
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
	c.output = append(c.output, (operand % 8))
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
	tmp := ""
	core := cpu{}
	fmt.Sscanf(lines[0], "Register A: %d", &core.a)
	fmt.Sscanf(lines[1], "Register B: %d", &core.b)
	fmt.Sscanf(lines[2], "Register C: %d", &core.c)
	fmt.Sscanf(lines[4], "Program: %s", &tmp)

	for _, v := range strings.Split(tmp, ",") {
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

	// Part 2:
	// The given input seems to run an inner program, that end in 5,5 to print
	// the B register value, and then loops by calling 3,0. The given example
	// has a similar structure, but prints the a register with 5,4.
	// The logic idea would be to run the iterations backwards, to get the
	// required inputs for the desired output of each byte of the program. The
	// ending loop also indicates that the A register must end with a value of
	// 0.

	// Input analysis
	// 2,4,1,5,7,5,1,6,0,3,4,2,5,5,3,0
	//
	// 1: 2,4 : b = a % 8
	// 2: 1,5 : b = b ^ 5  (flips bits 2 and 0)
	// 3: 7,5 : c = a >> b
	// 4: 1,6 : b = b ^ 6  (flips bits 2 and 1)
	// 5: 0,3 : a = a >> 3
	// 6: 4,2 : b = b ^ c
	// 7: 5,5 : out(b % 8)
	// 8: 3,0 : if a: goto 0

	// So, going backwards, to output the expected 0 in the last iteration:
	// 8: 'a' must end as 0 to end the loop.
	// 7: 'b' must end as a multiple of 8, to output 0 (3-lsb == 0).                  | b: 0b000
	// 6: 'c' must be (b % 8), to flips the bits leaving a multiple of 8              | b: 0b101 and c: 0b101
	// 5: 'a' comes from dividing by 8, so we can have up to 3-lsb set at this point. | a: 0bcccbbb >> 3 : Loop condition based on C
	// 4: 'b' will flip its bits 2 and 1.                                             | b: 0b101
	// 3: 'c' now comes from a >> b, and must be equal to (b ^ 6) & 0x7               | c: a >> 0b011 (3) | a = 0bcccbbb
	// 2: 'b' will flip its bits 2 and 0.                                             | b: 0b110
	// 1: 'b' = a & 0x7                                                               | a: 0bccc110

	// At the start, b comes from the last 3 bits of a, get's fliped on b2 and b0
	// And c comes from shifting a with b... and must endup as (b ^ 6) & 0x7:
	//   So: a in bits for the last iteration could be = cccbbb giving us
	//   control of the bbb value, the ccc value and L is the loop condition (0
	//   in the last iteration).
	//
	// The key point to minimize a, might be the shift c = a >> b.
	// To have complete control over C, we would need to shift a >> 3 bits.
	// This would fix the value of B needed for this shift, and at the same
	// time fixes the value of B what will be XOR to get the output, meaning we
	// can choose C based on the desired output. This strategy would leave a
	// composition of the A register as:
	// A = 0bcccbbb
	//
	// Now, as A is finally shift by 3 at step 5, the loop condition is
	// basically C... this means that the final step is a special case, where C
	// must be 0 to end the loop.
	//
	//
	// To wrap up:
	// * To shift by 3, giving us control of c, b allways reach the final ^ c as 0b101.
	// * From this I can determine C to get the desired value output by xoring it with 0b101.
	// * The final composition of A would be:
	//  -> a = 0bccc110 // Where a C != 0 means the loops keeps going.
	//
	// The final iteration is a special case, where C is fixed to 0 to end the
	// loop, this means that b must be 0 also on that iteration, this means
	// that b must be 0 also on that iteration, which also sets the final shift
	// for C to be a fixed 0b110 (6), and an initial value of A on that
	// iteration of 0b011. So, as a test, a starting condition of A = 0b11
	// shoud output a single 0. [CHECKED]
	//
	// This is where it gets kinda loopy. Having control of C with a fixed
	// shift of 3 in the first iteration, means that A ends the loop with the
	// same value of C (as stated).... but in the next iteration, B is obtained
	// from the 3-lsb of A, which are the C of the previous iteration......
	// This means that after setting a shift and a characer, the next shift is
	// fixed, and we must obey it. That means that now we kida have the start
	// and the end of the desired A... we have to figure out the shifts for the
	// rest of the values. Keeping in mind that A keeps getting shifted right,
	// the MSB of A should end with 0b11 in the last iteration, and the LSB of
	// A should start with 0b110 to have a first shift of 3. We can build the
	// inbetweens
	//
	// Let's rethink:
	//
	// So, to end with output '0', A must end the previous to last iteration as
	// 0b011
	// A = 0b011xxx
	// We only have 3 LSB to control (to start the next iteration with 0b11)
	// and we have to determine the shift and the value of C to get an output
	// of '3'
	//
	// Let's assume we don't want to shift, minimizing A.
	// A shift of 0 means a that 3LSB(C) = 3LSB(Bi)
	// B value of at the last XOR of 0b110
	// This needs a value of C of 0b101 to output the '3'.
	// But the initial Bi value of 0b110 for shift 0... can't be done.
	//
	// A shift of 1, Bs = 0b001 | Bi = 0b100 | Bx = 0b111 -> needs C = 0b100 but C is 0b011100 >> 1 NOPE
	// A shift of 2, Bs = 0b010 | Bi = 0b111 | Bx = 0b100 -> needs C = 0b111 but C is 0b011111 >> 2 YES!
	//
	// So, an A = 0b011111, outputs the desired '3,0'...
	// At the previous iteration, A would be
	// A = 0b011111xxx
	// We can keep making the same logic... Let's automate:
	//

	// This is specific for my input:
	type node struct {
		a int
		l int
	}

	queue := []node{{0, 1}}

	for len(queue) != 0 {
		n := queue[0]
		queue = queue[1:]

		if n.l > len(core.mem) {
			// We found a value of A that represents all core.mem elements.
			fmt.Println("Part 2:", n.a)
			break
		}

		for k := range 8 {
			// 8 possible values for the LSB of A at each iteration.
			candidate := n.a<<3 | k

			// Restart the CPU and run the program.
			core.pc = 0
			core.output = nil
			core.a = candidate
			core.RunProgram()

			if equalSlice(core.output, core.mem[len(core.mem)-n.l:]) {
				// Found a value that reproduce the last n.l elements.
				fmt.Printf("%d: Found candidate for a: %d (0b%s)\n", n.l, candidate, strconv.FormatInt(int64(candidate), 2))
				queue = append(queue, node{candidate, n.l + 1})
			}
		}
	}
}

func equalSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
