package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"sync"
)

// GameOfLife - holds the state of the game
type GameOfLife struct {
	generation int
	board      map[int64](map[int64]bool)
	rwMutex    sync.RWMutex
	pushMutex  sync.Mutex
}

// GameOfLifeHandler - hold the game and multiplexer
type GameOfLifeHandler struct {
	mux        *http.ServeMux
	gameOfLife GameOfLife
}

// Game of life implements Handler interface
func (h *GameOfLifeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// Creates new GameOfLifeHandler
func NewGameOfLifeHandler(startCells [][2]int64) *GameOfLifeHandler {
	gameOfLife := GameOfLife{generation: 0,
		board: make(map[int64]map[int64]bool), rwMutex: sync.RWMutex{},
		pushMutex: sync.Mutex{}}
	for i := 0; i < len(startCells); i++ {
		gameOfLife.addCell(startCells[i][0], startCells[i][1])
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cell/status/", gameOfLife.getCellStatus)
	mux.HandleFunc("/generation/", gameOfLife.getGeneration)
	mux.HandleFunc("/cells/", gameOfLife.addCells)
	mux.HandleFunc("/generation/evolve/", gameOfLife.evolve)
	mux.HandleFunc("/reset/", gameOfLife.reset)

	gameOfLifeHandler := GameOfLifeHandler{mux: mux, gameOfLife: gameOfLife}

	return &gameOfLifeHandler
}

// Add a living cell to the game board
func (game *GameOfLife) addCell(x int64, y int64) {
	addToBoard(game.board, x, y)
}

// Add (x, y) to a board
func addToBoard(board map[int64]map[int64]bool, x int64, y int64) {
	ym, ok := board[x]
	if !ok {
		ym = make(map[int64]bool)
		board[x] = ym
	}
	ym[y] = true
}

// Check if cell (x, y) is alive
func (game *GameOfLife) isAlive(x int64, y int64) bool {
	ym, ok := game.board[x]
	if !ok {
		return false
	}
	alive, ok2 := ym[y]
	if !ok2 {
		return false
	}
	return alive
}

// Type used for creating json for /cell/status/ requests
type Alive struct {
	Alive bool `json:"alive"`
}

// Writes a response
func message(w http.ResponseWriter, message []byte, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if message != nil {
		w.Write(message)
	}
}

// Responsible to answer to /cell/status/
func (game *GameOfLife) getCellStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
	}
	xStr := r.URL.Query().Get("x")
	yStr := r.URL.Query().Get("y")
	x, errX := strconv.ParseInt(xStr, 10, 64)
	y, errY := strconv.ParseInt(yStr, 10, 64)
	if errX != nil {
		http.Error(w, errX.Error(), http.StatusBadRequest)
	}
	if errY != nil {
		http.Error(w, errY.Error(), http.StatusBadRequest)
	}
	game.rwMutex.RLock()
	alive, _ := json.Marshal(Alive{Alive: game.isAlive(x, y)})
	game.rwMutex.RUnlock()

	message(w, alive, http.StatusOK)
}

// Type used for creating json for /generation/ requests
type Generation struct {
	Generation int        `json:"generation"`
	Living     [][2]int64 `json:"living"`
}

// Responsible to answer to /generation/ requests
func (game *GameOfLife) getGeneration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
	}
	game.rwMutex.RLock()
	generation, _ := json.Marshal(Generation{Generation: game.generation, Living: game.getLiving()})
	game.rwMutex.RUnlock()
	message(w, generation, http.StatusOK)

}

// Returns all the living cells on the game board
func (game *GameOfLife) getLiving() [][2]int64 {
	living := make([][2]int64, 0)
	for x, ym := range game.board {
		for y, alive := range ym {
			if alive {
				s := [2]int64{x, y}
				living = append(living, s)
			}
		}
	}
	return living
}

// Type used to hold the points received with /cells/ requests
type Point struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

// Responsible to answer to /cells/ requests
func (game *GameOfLife) addCells(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	}

	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	s := make([]Point, 0)
	err1 := json.Unmarshal(bytes, &s)

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
	}

	game.pushMutex.Lock()
	game.rwMutex.Lock()
	for _, p := range s {
		game.addCell(p.X, p.Y)
	}
	game.rwMutex.Unlock()
	game.pushMutex.Unlock()

	message(w, nil, http.StatusCreated)
}

// Responsible to answer to /generation/evolve/ requests
func (game *GameOfLife) evolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	}

	// has to lock here for a little while
	game.pushMutex.Lock()
	game.rwMutex.RLock()

	newBoard := make(map[int64]map[int64]bool)

	for x, ym := range game.board {
		for y, alive := range ym {
			if alive {
				count := game.getLivingNeighbours(x, y)
				if count == 2 || count == 3 {
					addToBoard(newBoard, x, y)
				}

				//a dead cell with 3 alive neighbours is to be found only around living cells
				//so check the neighbours if this cell
				game.addBornCellsAround(newBoard, x, y)
			}
		}
	}
	game.rwMutex.RUnlock()

	// it is unwise to allow reading at this point, so lock again
	game.rwMutex.Lock()
	game.generation += 1
	game.board = newBoard
	game.rwMutex.Unlock()
	game.pushMutex.Unlock()

	message(w, nil, http.StatusNoContent)
}

// Returns the number of living neighbours around a cell. Counts only to 4 to be more efficient
func (game *GameOfLife) getLivingNeighbours(x int64, y int64) (count int) {
	var min int64 = math.MinInt64
	var max int64 = math.MaxInt64

	var minX, minY, maxX, maxY int64

	//clear corner cases
	if x > min {
		minX = x - 1
	}
	if x < max {
		maxX = x + 1
	}
	if y > min {
		minY = y - 1
	}
	if y < max {
		maxY = y + 1
	}

	count = 0

	for i := minX; i <= maxX; i++ {
		for j := minY; j <= maxY; j++ {
			if i == x && j == y {
				continue
			}
			if game.isAlive(i, j) {
				count += 1
				// no reason to check for more than 4 alive neighbours since the cell is overcrowded
				if count == 4 {
					return count
				}
			}
		}
	}
	return count
}

// Searches for places where cells have to be born and adds them to the new board
func (game *GameOfLife) addBornCellsAround(newBoard map[int64]map[int64]bool, x int64, y int64) {
	var min int64 = math.MinInt64
	var max int64 = math.MaxInt64

	var minX, minY, maxX, maxY int64

	//clear corner cases
	if x > min {
		minX = x - 1
	}
	if x < max {
		maxX = x + 1
	}
	if y > min {
		minY = y - 1
	}
	if y < max {
		maxY = y + 1
	}

	for i := minX; i <= maxX; i++ {
		for j := minY; j <= maxY; j++ {
			if i == x && j == y {
				// skip the center of the search
				continue
			}
			if !game.isAlive(i, j) {
				// dead cell found - count its neighbours
				count := game.getLivingNeighbours(i, j)
				if count == 3 {
					addToBoard(newBoard, i, j)
				}
			}
		}
	}
}

// Responsible to answer to /reset/ requests
func (game *GameOfLife) reset(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	}
	game.pushMutex.Lock()
	game.rwMutex.Lock()
	game.generation = 0
	game.board = make(map[int64]map[int64]bool)
	game.rwMutex.Unlock()
	game.pushMutex.Unlock()

	message(w, nil, http.StatusNoContent)
}
