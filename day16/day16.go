package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"math"
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

// Direction vectores, starting from East.
var dir = []point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func main() {
	fmt.Printf("AoC 2024 - Day 16\n")
	lines := readInput(input)
	solve(lines)
}

type point struct{ x, y int }
type dpoint struct {
	p point
	d int
}

type charmap struct {
	data          []rune
	height, width int
}

func (c charmap) index(p point) int {
	return p.y*c.width + p.x
}

func (c charmap) point(i int) point {
	return point{i % c.width, i / c.width}
}

func (c charmap) isValid(p point) bool {
	if p.x < 0 || p.x > c.width-1 || p.y < 0 || p.y > c.height-1 {
		return false
	}

	return true
}

// Part solutions.
func solve(lines []string) {
	m := charmap{
		height: len(lines),
		width:  len(lines[0]),
		data:   []rune(strings.Join(lines, "")),
	}

	// Find the starting point.
	i := strings.IndexRune(string(m.data), 'S')
	start := m.point(i)

	// Part 1
	bestcost, seats := runmaze(&m, start)
	fmt.Println("Part 1:", bestcost)
	fmt.Println("Part 2:", seats)
}

func plotMap(m *charmap, modifier map[point]bool) {
	for y := range m.height {
		for x := range m.width {
			if modifier != nil && modifier[point{x, y}] {
				fmt.Printf("%c", 'O')
			} else {
				fmt.Printf("%c", m.data[m.index(point{x, y})])
			}
		}
		fmt.Println()
	}
}

type node struct {
	cost int
	p    point
	dir  int
	from dpoint
}

type PriorityQueue []*node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	item := x.(*node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func runmaze(m *charmap, start point) (int, int) {
	cellcost := make(map[dpoint]int)
	cellcost[dpoint{start, 0}] = 0
	cellfrom := make(map[dpoint][]dpoint)
	ep := m.point(strings.IndexRune(string(m.data), 'E'))
	bestcost := math.MaxInt
	endstate := dpoint{}

	pq := &PriorityQueue{&node{cost: 0, p: start, dir: 0}}
	heap.Init(pq)

	for pq.Len() != 0 {
		cnode := heap.Pop(pq).(*node)
		if cv, ok := cellcost[dpoint{cnode.p, cnode.dir}]; ok && cnode.cost > cv {
			continue
		}

		cellcost[dpoint{cnode.p, cnode.dir}] = cnode.cost
		if cnode.p == ep {
			if cnode.cost > bestcost {
				break
			}
			bestcost = cnode.cost
			endstate = dpoint{cnode.p, cnode.dir}
		}

		newrec := true
		for _, v := range cellfrom[dpoint{cnode.p, cnode.dir}] {
			if v == cnode.from {
				newrec = false
			}
		}

		if newrec {
			cellfrom[dpoint{cnode.p, cnode.dir}] = append(cellfrom[dpoint{cnode.p, cnode.dir}], cnode.from)
		}

		// Get the next exploration points from this node.
		for k := range 4 {
			// Skip 180° turns
			if k == 2 {
				continue
			}

			diff := dir[(cnode.dir+k)%4]
			np := point{cnode.p.x + diff.x, cnode.p.y + diff.y}

			// Skip if wall.
			if m.data[m.index(np)] == '#' {
				continue
			}

			// Skip already processed cells.
			cost := cnode.cost + 1 + 1000*(k%2)
			if v, ok := cellcost[dpoint{np, (cnode.dir + k) % 4}]; ok && cost > v {
				continue
			}

			heap.Push(pq,
				&node{cost: cost, p: np, dir: (cnode.dir + k) % 4,
					from: dpoint{cnode.p, cnode.dir}})
		}
	}

	// Part 2. Backtrack from the end position, counting visted cells.
	visited := make(map[dpoint]bool)
	queue := []dpoint{endstate}
	visited[endstate] = true

	seats := make(map[point]bool)
	seats[endstate.p] = true

	for len(queue) != 0 {
		cdp := queue[0]
		queue = queue[1:]

		if cdp.p == start {
			break
		}

		fmt.Println(cdp, " from:", cellfrom[cdp])
		for _, v := range cellfrom[cdp] {
			if _, ok := visited[v]; ok {
				continue
			}

			queue = append(queue, v)
			visited[v] = true
			seats[v.p] = true
		}
	}

	// Can be used to plot the final marked seats.
	// plotMap(m, seats)

	return bestcost, len(seats)
}
