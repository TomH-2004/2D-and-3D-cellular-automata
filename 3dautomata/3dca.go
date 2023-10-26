package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	Width  = 20
	Height = 20
	Depth  = 20
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/ws", websocketHandler)

	http.Handle("/", r)
	fmt.Println("Server listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

var updateSpeed time.Duration = 250 * time.Millisecond

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	cells := make([][][]bool, Width)
	for x := 0; x < Width; x++ {
		cells[x] = make([][]bool, Height)
		for y := 0; y < Height; y++ {
			cells[x][y] = make([]bool, Depth)
			for z := 0; z < Depth; z++ {
				cells[x][y][z] = rand.Float32() < 0.5
			}
		}
	}

	gliderGunPattern(cells, 2, 2, 2)

	go func() {
		for {
			err := conn.WriteJSON(cells)
			if err != nil {
				fmt.Println(err)
				return
			}

			cells = update(cells)
			time.Sleep(updateSpeed)
		}
	}()

	for {
		var speedUpdate struct {
			Speed int `json:"speed"`
		}
		err := conn.ReadJSON(&speedUpdate)
		if err != nil {
			fmt.Println(err)
			return
		}

		updateSpeed = time.Duration(speedUpdate.Speed) * time.Millisecond
	}
}

func update(cells [][][]bool) [][][]bool {
	newCells := make([][][]bool, Width)
	for x := 0; x < Width; x++ {
		newCells[x] = make([][]bool, Height)
		for y := 0; y < Height; y++ {
			newCells[x][y] = make([]bool, Depth)
			for z := 0; z < Depth; z++ {
				count := countNeighbors(cells, x, y, z)
				if cells[x][y][z] {
					newCells[x][y][z] = count >= 4 && count <= 6
				} else {
					newCells[x][y][z] = count == 5
				}
			}
		}
	}
	return newCells
}

func countNeighbors(cells [][][]bool, x, y, z int) int {
	count := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				nx, ny, nz := x+dx, y+dy, z+dz
				if nx >= 0 && nx < Width && ny >= 0 && ny < Height && nz >= 0 && nz < Depth {
					if cells[nx][ny][nz] {
						count++
					}
				}
			}
		}
	}
	return count
}

func gliderGunPattern(cells [][][]bool, xOffset, yOffset, zOffset int) {
	gliderGun := [9][36]bool{

		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	}

	gliderGunWidth := len(gliderGun)
	gliderGunHeight := len(gliderGun[0])
	gliderGunDepth := len(cells[0][0])

	for x := 0; x < gliderGunWidth && x+xOffset < len(cells); x++ {
		for y := 0; y < gliderGunHeight && y+yOffset < len(cells[x+xOffset]); y++ {
			for z := 0; z < gliderGunDepth && z+zOffset < len(cells[x+xOffset][y+yOffset]); z++ {
				cells[x+xOffset][y+yOffset][z+zOffset] = gliderGun[x][y]
			}
		}
	}
}
