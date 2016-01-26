package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestCreatingBoardWithSeed(t *testing.T) {
	testSrv := setUpServer([][2]int64{
		{0, 0},
		{1, 1},
		{3, 4},
		{-2, 5},
		{-19023482123, 5},
	})
	defer testSrv.Close()

	testTable := []struct {
		alive bool
		x, y  int64
	}{
		{x: 0, y: 0, alive: true},
		{x: 1, y: 1, alive: true},
		{x: 3, y: 4, alive: true},
		{x: -2, y: 5, alive: true},
		{x: -19023482123, y: 5, alive: true},
		{x: 0, y: -1, alive: false},
		{x: -100, y: 100, alive: false},
		{x: 100, y: 100, alive: false},
		{x: 55, y: 93, alive: false},
	}

	for _, testCase := range testTable {
		path := fmt.Sprintf("/cell/status/?x=%d&y=%d", testCase.x, testCase.y)
		cellURL := buildUrl(testSrv.URL, path)

		resp, err := http.Get(cellURL)

		if err != nil {
			t.Errorf("Error getting empty board: %s", err)
		}

		defer resp.Body.Close()

		cellStatus := &struct {
			Alive bool `json:"alive"`
		}{}

		respBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil && err != io.EOF {
			t.Errorf("Error reading from response: %s", err)
		}

		fmt.Println("Response body " + string(respBytes[:]))

		if err := json.Unmarshal(respBytes, cellStatus); err != nil {
			t.Errorf("Error decoding json: %s", err)
		}

		if testCase.alive != cellStatus.Alive {
			t.Errorf("Expected alive %t for (%d, %d) but it was %t. "+
				"JSON: %s", testCase.alive, testCase.x, testCase.y, cellStatus.Alive,
				string(respBytes))
		}
	}

}

func TestGetLiving(t *testing.T) {
	testSrv := setUpServer([][2]int64{
		{0, 0},
		{1, 1},
		{3, 4},
		{-2, 5},
		{-19023482123, 5},
	})
	defer testSrv.Close()

	url := buildUrl(testSrv.URL, "/generation/")
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	//todo: check the value with some code
	fmt.Println(string(respBytes[:]))
}

func TestReset(t *testing.T) {
	testSrv := setUpServer([][2]int64{
		{0, 0},
		{1, 1},
		{3, 4},
		{-2, 5},
		{-19023482123, 5},
	})
	defer testSrv.Close()

	url := buildUrl(testSrv.URL, "/generation/")
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	//todo: check the value with some code
	fmt.Println(string(respBytes[:]))

	url1 := buildUrl(testSrv.URL, "/reset/")
	resp1, err1 := http.Post(url1, "text/plain", nil)
	if err1 != nil {
		t.Fatal(err1.Error())
	}

	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204 but found %d", resp1.StatusCode)
	}

	url2 := buildUrl(testSrv.URL, "/generation/")
	resp2, err2 := http.Get(url2)
	if err2 != nil {
		t.Fatal(err2.Error())
	}

	defer resp2.Body.Close()

	respBytes2, err2 := ioutil.ReadAll(resp2.Body)

	//todo: check the value with some code
	fmt.Println(string(respBytes2[:]))
}

func TestAddCell(t *testing.T) {
	testSrv := setUpServer([][2]int64{
		{0, 0},
		{1, 1},
		{3, 4},
		{-2, 5},
		{-19023482123, 5},
	})
	defer testSrv.Close()

	url1 := buildUrl(testSrv.URL, "/cells/")
	var data = []byte(`[{"x": 42, "y": 43}, {"x": -5, "y": 10}]`)

	resp1, err1 := http.Post(url1, "application/json", bytes.NewBuffer(data))
	if err1 != nil {
		t.Fatal(err1.Error())
	}

	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201 but found %d", resp1.StatusCode)
	}

	url2 := buildUrl(testSrv.URL, "/generation/")
	resp2, err2 := http.Get(url2)
	if err2 != nil {
		t.Fatal(err2.Error())
	}

	defer resp2.Body.Close()

	respBytes2, err2 := ioutil.ReadAll(resp2.Body)

	//todo: check the value with some code
	fmt.Println(string(respBytes2[:]))
}

func TestAddCellDupl(t *testing.T) {
	testSrv := setUpServer([][2]int64{
		{0, 0},
		{1, 1},
		{3, 4},
		{-2, 5},
		{-19023482123, 5},
	})
	defer testSrv.Close()

	url1 := buildUrl(testSrv.URL, "/cells/")
	var data = []byte(`[{"x": 42, "y": 43}, {"x": 3, "y": 4}]`)

	resp1, err1 := http.Post(url1, "application/json", bytes.NewBuffer(data))
	if err1 != nil {
		t.Fatal(err1.Error())
	}

	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201 but found %d", resp1.StatusCode)
	}

	url2 := buildUrl(testSrv.URL, "/generation/")
	resp2, err2 := http.Get(url2)
	if err2 != nil {
		t.Fatal(err2.Error())
	}

	defer resp2.Body.Close()

	respBytes2, err2 := ioutil.ReadAll(resp2.Body)

	//todo: check the value with some code
	fmt.Println(string(respBytes2[:]))
}

func TestNeighbours(t *testing.T) {
	gameOfLife := GameOfLife{generation: 0, board: make(map[int64]map[int64]bool), rwMutex: sync.RWMutex{}}
	gameOfLife.addCell(1, 2)
	gameOfLife.addCell(2, 3)

	count := gameOfLife.getLivingNeighbours(1, 2)
	if count != 1 {
		t.Errorf("EXpected count is one, found %d", count)
	}

	gameOfLife.addCell(math.MinInt64, math.MinInt64)
	gameOfLife.addCell(math.MaxInt64, math.MaxInt64)

	count = gameOfLife.getLivingNeighbours(math.MinInt64, math.MinInt64)
	if count != 0 {
		t.Errorf("EXpected count is 0, found %d", count)
	}

	count = gameOfLife.getLivingNeighbours(math.MaxInt64, math.MaxInt64)
	if count != 0 {
		t.Errorf("Expected count is 0, found %d", count)
	}
}

func TestEvolve(t *testing.T) {
	testSrv := setUpServer([][2]int64{
		{0, 0},
		{1, 1},
		{1, 2},
		{-2, 5},
		{-19023482123, 5},
	})
	defer testSrv.Close()

	url := buildUrl(testSrv.URL, "/generation/evolve/")
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204 but found %d", resp.StatusCode)
	}

	url1 := buildUrl(testSrv.URL, "/generation/")
	resp1, err1 := http.Get(url1)
	if err1 != nil {
		t.Fatal(err1.Error())
	}

	defer resp1.Body.Close()

	respBytes, _ := ioutil.ReadAll(resp1.Body)

	//todo: check the value with some code
	fmt.Println(string(respBytes[:]))
}

func TestAsyncEvolve(t *testing.T) {
	testSrv := setUpServer([][2]int64{
		{0, 0},
		{1, 1},
		{1, 2},
		{-2, 5},
		{-19023482123, 5},
	})
	defer testSrv.Close()


	url := buildUrl(testSrv.URL, "/generation/evolve/")

	for i := 0; i < 100; i++ {
		go func(url string) {
			resp, err := http.Post(url, "text/plain", nil)

			if err != nil {
				t.Fatal(err.Error())
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("Expected status 204 but found %d", resp.StatusCode)
			}
		}(url)
	}

	url1 := buildUrl(testSrv.URL, "/cells/")
	var data = []byte(`[{"x": 42, "y": 43}, {"x": -5, "y": 10}]`)

	for j := 0; j < 100; j++ {
		go func(url1 string) {
			resp1, err1 := http.Post(url1, "application/json", bytes.NewBuffer(data))
			if err1 != nil {
				t.Fatal(err1.Error())
			}

			defer resp1.Body.Close()

			if resp1.StatusCode != http.StatusCreated {
				t.Errorf("Expected status 201 but found %d", resp1.StatusCode)
			}
		}(url1)
	}
	time.Sleep(10 * time.Second)
}

func TestBadQuery(t *testing.T) {
	testSrv := setUpServer([][2]int64{
		{0, 0},
		{1, 1},
		{1, 2},
		{-2, 5},
		{-19023482123, 5},
	})
	defer testSrv.Close()

	url := buildUrl(testSrv.URL, "/wrong")
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 but found %d", resp.StatusCode)
	}
}

/* Utility functions */

func buildUrl(baseUrl, path string) string {
	return fmt.Sprintf("%s%s", baseUrl, path)
}

// Users of this function are resposible for calling Close() on the returned server.
// Failure to do so will result in leaked resources.
func setUpServer(cells [][2]int64) *httptest.Server {
	gofh := NewGameOfLifeHandler(cells)
	return httptest.NewServer(gofh)
}
