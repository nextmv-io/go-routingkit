//main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/nextmv-io/go-routingkit/routingkit"
)

func main() {
	mapFile := "data/map.osm.pbf"
	chFile := "data/map.ch"

	start := time.Now()
	client, err := routingkit.Wrapper(mapFile, chFile)
	if err != nil {
		fmt.Println("Error when creating client %", err)
	}
	fmt.Println("Load pbf and build or load ch took %v", time.Since(start))

	testFile(client)

	start = time.Now()
	p1 := []float64{-76.58749, 39.29971}
	p2 := []float64{-76.59735, 39.30587}

	points := make([][]float64, 1000)
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			points[i] = p1
		} else {
			points[i] = p2
		}

	}

	wg := sync.WaitGroup{}
	wg.Add(1000000)
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			go func(i, j int) {
				defer wg.Done()
				_ = client.Threaded(points[i], points[j])
			}(i, j)
		}
	}
	wg.Wait()
	fmt.Println("10000 queries took ", time.Since(start))
}

func testFile(client routingkit.Client) {
	// Prepare test data
	var file string
	flag.StringVar(
		&file,
		"testfile",
		"",
		"path to test file",
	)
	flag.Parse()
	if file == "" {
		fmt.Println("Skipping file test (no file specified; use -testfile <path>)")
		return
	}
	fmt.Println("Reading points from file: ", file)
	points := readPoints(file)[:100] //take the first 100 points only

	// Execute test query
	start := time.Now()
	for _, p1 := range points {
		for _, p2 := range points {
			_ = client.Distance(p1, p2)
			// fmt.Printf("%f,%f -> %f,%f : %f \n", p1[0], p1[1], p2[0], p2[1], dist)
		}
	}
	queries := len(points) * len(points)
	elapsed := time.Since(start)
	avg := time.Duration(float64(elapsed.Nanoseconds())/float64(queries)) * time.Nanosecond
	fmt.Printf("%d queries took %v (avg: %v)\n", queries, elapsed, avg)
}

func readPoints(path string) [][]float64 {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var b pointData
	err = json.Unmarshal(dat, &b)
	if err != nil {
		fmt.Printf("error deserializing: %s\n", string(dat))
		return [][]float64{}
	}
	return b.Points
}

type pointData struct {
	Points [][]float64 `json:"points"`
}
