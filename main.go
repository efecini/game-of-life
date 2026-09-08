package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

// Grid is a 2D board of cells. true means alive.
type Grid [][]bool

func newGrid(rows, cols int) Grid {
	g := make(Grid, rows)
	for i := range g {
		g[i] = make([]bool, cols)
	}
	return g
}

func (g Grid) rows() int { return len(g) }
func (g Grid) cols() int {
	if len(g) == 0 {
		return 0
	}
	return len(g[0])
}

func (g Grid) clone() Grid {
	out := newGrid(g.rows(), g.cols())
	for r := range g {
		copy(out[r], g[r])
	}
	return out
}

// neighbors counts living cells among the 8 surrounding cells.
// The board wraps at the edges (a torus), so patterns can travel off one
// side and reappear on the opposite side.
func (g Grid) neighbors(r, c int) int {
	rows, cols := g.rows(), g.cols()
	count := 0
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr := (r + dr + rows) % rows
			nc := (c + dc + cols) % cols
			if g[nr][nc] {
				count++
			}
		}
	}
	return count
}

// step applies Conway's Game of Life rules and returns the next generation:
//   - a live cell with 2 or 3 neighbors stays alive
//   - a dead cell with exactly 3 neighbors becomes alive
//   - every other cell is dead
func (g Grid) step() Grid {
	next := newGrid(g.rows(), g.cols())
	for r := 0; r < g.rows(); r++ {
		for c := 0; c < g.cols(); c++ {
			n := g.neighbors(r, c)
			alive := g[r][c]
			next[r][c] = n == 3 || (alive && n == 2)
		}
	}
	return next
}

func (g Grid) populateRandom(density float64, rng *rand.Rand) {
	for r := 0; r < g.rows(); r++ {
		for c := 0; c < g.cols(); c++ {
			g[r][c] = rng.Float64() < density
		}
	}
}

func (g Grid) stamp(r, c int, pattern [][]int) {
	rows, cols := g.rows(), g.cols()
	for pr := range pattern {
		for pc, v := range pattern[pr] {
			if v == 0 {
				continue
			}
			g[(r+pr+rows)%rows][(c+pc+cols)%cols] = true
		}
	}
}

var patterns = map[string][][]int{
	"glider": {
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 1},
	},
	"blinker": {
		{1, 1, 1},
	},
	"toad": {
		{0, 1, 1, 1},
		{1, 1, 1, 0},
	},
	"beacon": {
		{1, 1, 0, 0},
		{1, 1, 0, 0},
		{0, 0, 1, 1},
		{0, 0, 1, 1},
	},
	"pulsar": {
		{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
	},
	"gosper": {
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
}

func (g Grid) String() string {
	var b strings.Builder
	b.Grow(g.rows() * (g.cols()*2 + 1))
	for r := 0; r < g.rows(); r++ {
		for c := 0; c < g.cols(); c++ {
			if g[r][c] {
				b.WriteString("██")
			} else {
				b.WriteString("  ")
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	rows := flag.Int("rows", 24, "grid height")
	cols := flag.Int("cols", 40, "grid width")
	fps := flag.Int("fps", 10, "frames per second")
	density := flag.Float64("density", 0.28, "fill probability for random pattern (0-1)")
	pattern := flag.String("pattern", "random", "initial pattern: random, glider, blinker, toad, beacon, pulsar, gosper")
	seed := flag.Int64("seed", time.Now().UnixNano(), "random seed")
	generations := flag.Int("gens", 0, "number of generations to run (0 = infinite until Ctrl+C)")
	web := flag.Bool("web", false, "serve the browser UI instead of the terminal animation")
	addr := flag.String("addr", "127.0.0.1:8080", "HTTP listen address for -web")
	flag.Parse()

	if *web {
		serveWeb(*addr)
		return
	}

	if *rows < 3 || *cols < 3 {
		fmt.Fprintln(os.Stderr, "rows and cols must be at least 3")
		os.Exit(1)
	}

	grid := newGrid(*rows, *cols)
	name := strings.ToLower(*pattern)

	switch name {
	case "random":
		rng := rand.New(rand.NewSource(*seed))
		grid.populateRandom(*density, rng)
	default:
		p, ok := patterns[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown pattern %q; use random, glider, blinker, toad, beacon, pulsar, or gosper\n", *pattern)
			os.Exit(1)
		}
		r := (*rows - len(p)) / 2
		c := (*cols - len(p[0])) / 2
		if r < 0 {
			r = 0
		}
		if c < 0 {
			c = 0
		}
		grid.stamp(r, c, p)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	clear := "\033[H\033[2J"
	hide := "\033[?25l"
	show := "\033[?25h"
	fmt.Print(hide)
	defer fmt.Print(show)

	ticker := time.NewTicker(time.Second / time.Duration(*fps))
	defer ticker.Stop()

	gen := 0
	fmt.Print(clear)
	fmt.Printf("Game of Life  |  pattern=%s  |  %dx%d  |  gen %d  |  Ctrl+C to quit\n\n%s",
		name, *cols, *rows, gen, grid)

	for {
		select {
		case <-interrupt:
			fmt.Println()
			return
		case <-ticker.C:
			grid = grid.step()
			gen++
			fmt.Print(clear)
			fmt.Printf("Game of Life  |  pattern=%s  |  %dx%d  |  gen %d  |  Ctrl+C to quit\n\n%s",
				name, *cols, *rows, gen, grid)
			if *generations > 0 && gen >= *generations {
				return
			}
		}
	}
}
