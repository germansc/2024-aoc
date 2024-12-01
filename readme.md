# 2024 - Advent of Code

This repo contains my solution to the 2024 Advent of Code puzzles.

This year's language of choice will be go, to further explore its simplicity
and power for solving this kind of algorithmic challenges.

### Repository Structure

Each day's puzzle is organized in a separate directory with its own package.
The solution for each day is implemented to take input from `stdin` and print
the result to `stdout`. The directory for each day contains the following
structure:

```
dayXX/
├── dayXX.go         # Solution implementation
├── input.txt        # Raw input data for the puzzle
├── sample.txt       # Example input extracted from the puzzle description
└── puzzle.md        # The original puzzle description
```

### Running the solutions

To run a solution for a specific day, navigate to the corresponding directory
and run the solution using `go run`, piping the `input.txt` file as `stdin`.
For example, to run Day 1:

```bash
cd day01
go run . < input.txt
```

###### 2024 | germansc
