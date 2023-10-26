package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

const (
	Width  = 40
	Height = 20
)

type Cell struct {
	Alive bool
}

type Grid [][]*Cell

func main() {
	rand.Seed(time.Now().UnixNano())

	grid := createGrid(Width, Height)
	initializeGrid(grid)

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, _ := tcell.NewScreen()
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	runGame(screen, grid)
}

func createGrid(width, height int) Grid {
	grid := make(Grid, height)
	for i := range grid {
		grid[i] = make([]*Cell, width)
	}
	return grid
}

func initializeGrid(grid Grid) {
	for y := range grid {
		for x := range grid[y] {
			cell := &Cell{
				Alive: rand.Intn(2) == 0,
			}
			grid[y][x] = cell
		}
	}
}

func runGame(screen tcell.Screen, grid Grid) {

	exitChan := make(chan struct{})

	go func() {
		for {
			ev := screen.PollEvent()
			switch e := ev.(type) {
			case *tcell.EventKey:
				if e.Key() == tcell.KeyRune && e.Rune() == 'q' {

					exitChan <- struct{}{}
					return
				}
			}
		}
	}()

	for {
		drawGrid(screen, grid)
		updateGrid(grid)

		select {
		case <-exitChan:
			return
		default:
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func updateGrid(grid Grid) {
	newGrid := createGrid(Width, Height)

	for y := range grid {
		for x := range grid[y] {
			cell := grid[y][x]
			aliveNeighbors := countAliveNeighbors(grid, x, y)

			newCell := &Cell{
				Alive: cell.Alive,
			}

			if cell.Alive {
				if aliveNeighbors < 2 || aliveNeighbors > 3 {
					newCell.Alive = false
				}
			} else {
				if aliveNeighbors == 3 {
					newCell.Alive = true
				}
			}

			newGrid[y][x] = newCell
		}
	}

	copyGrid(grid, newGrid)
}

func countAliveNeighbors(grid Grid, x, y int) int {
	count := 0
	neighbors := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, offset := range neighbors {
		nx, ny := x+offset[0], y+offset[1]

		if nx >= 0 && nx < Width && ny >= 0 && ny < Height {

			if grid[ny][nx] != nil && grid[ny][nx].Alive {
				count++
			}
		}
	}

	return count
}

func copyGrid(dest, src Grid) {
	for y := range dest {
		for x := range dest[y] {
			dest[y][x].Alive = src[y][x].Alive
		}
	}
}

func drawGrid(screen tcell.Screen, grid Grid) {
	screen.Clear()
	for y := range grid {
		for x := range grid[y] {
			cell := grid[y][x]
			if cell.Alive {
				screen.SetContent(x, y, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
			}
		}
	}
	screen.Show()
}
