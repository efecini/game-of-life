# Game of Life

A terminal [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) written in Go.

Each generation follows the classic rules:

- A live cell with 2 or 3 neighbors stays alive
- A dead cell with exactly 3 neighbors becomes alive
- Every other cell dies or stays empty

The board wraps at the edges, so patterns that walk off one side reappear on the other.

## Requirements

- [Go 1.22](https://go.dev/dl/) or later

## Run

```bash
git clone https://github.com/efecini/game-of-life.git
cd game-of-life
go run .
```

That starts a random board. Press **Ctrl+C** to stop.

## Browser

To watch the same animation in a browser:

```bash
go run . -web
```

Then open [http://127.0.0.1:8080](http://127.0.0.1:8080). The page can play/pause, step, load patterns, change speed, and let you paint cells with the mouse.

```bash
# Listen on a different address
go run . -web -addr 127.0.0.1:3000
```

## Examples

```bash
# A glider that walks across the grid
go run . -pattern glider

# Gosper glider gun (use a wider board)
go run . -pattern gosper -rows 30 -cols 50

# Pulsar oscillator
go run . -pattern pulsar

# Faster, denser random soup
go run . -fps 20 -density 0.35

# Run 200 generations then exit
go run . -gens 200
```

## Options

| Flag | Default | Description |
| --- | --- | --- |
| `-rows` | `24` | Grid height (minimum 3) |
| `-cols` | `40` | Grid width (minimum 3) |
| `-fps` | `10` | Frames per second |
| `-density` | `0.28` | Fill probability for `random` (0–1) |
| `-pattern` | `random` | Initial pattern |
| `-seed` | current time | Random seed |
| `-gens` | `0` | Generations to run (`0` means until Ctrl+C) |
| `-web` | `false` | Serve the browser UI |
| `-addr` | `127.0.0.1:8080` | HTTP listen address for `-web` |

Patterns: `random`, `glider`, `blinker`, `toad`, `beacon`, `pulsar`, `gosper`.
