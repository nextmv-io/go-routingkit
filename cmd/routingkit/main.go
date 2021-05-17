//main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nextmv-io/go-routingkit/routingkit"
)

func main() {

	mapFile := "data/map.osm.pbf"
	chFile := "data/map.ch"

	myClass := routingkit.NewClient()
	defer routingkit.DeleteClient(myClass)

	testFile(myClass)

	start := time.Now()
	if _, err := os.Stat(chFile); os.IsNotExist(err) {
		myClass.Build_ch(mapFile, chFile)
		fmt.Println("Load pbf and build ch took %v", time.Since(start))
	} else {
		myClass.Load(mapFile, chFile)
		fmt.Println("Load pbf and ch took %v", time.Since(start))
	}

	start = time.Now()
	p1 := routingkit.NewPoint()
	p1.SetLon(-76.58749)
	p1.SetLat(39.29971)
	p2 := routingkit.NewPoint()
	p2.SetLon(-76.59735)
	p2.SetLat(39.30587)

	var number int64 = 100
	points := routingkit.NewPointVector(number)
	defer routingkit.DeletePointVector(points)
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			points.Set(i, p1)
		} else {
			points.Set(i, p2)
		}

	}
	//distances := myClass.Table(points, points)
	_ = myClass.Table(points, points)
	// for i := 0; i < int(distances.Size()); i++ {
	// 	fmt.Printf("%f, ", distances.Get(i))
	// }
	fmt.Println("100 queries took ", time.Since(start))

	// wg := sync.WaitGroup{}
	// wg.Add(100)
	// for i := 0; i < 100; i++ {
	// 	go func() {
	// 		defer wg.Done()
	// 		_ = myClass.Distance(-76.58749, 39.29971, -76.59735, 39.30587)
	// 	}()
	// }
	// wg.Wait()
	//fmt.Println("100 queries took %v", time.Since(start))
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
	points := readPoints(file)

	// Execute test query
	start := time.Now()
	pointVector := routingkit.NewPointVector(int64(len(points)))
	defer routingkit.DeletePointVector(pointVector)
	for i, point := range points {
		p := routingkit.NewPoint()
		p.SetLon(float32(point[0]))
		p.SetLat(float32(point[1]))
		pointVector.Set(i, p)
	}
	distances := client.Table(pointVector, pointVector)
	queries := len(points) * len(points)
	elapsed := time.Since(start)
	avg := time.Duration(float64(elapsed.Nanoseconds())/float64(queries)) * time.Nanosecond
	fmt.Printf("%d queries took %v (avg: %v)\n", queries, elapsed, avg)
	var total float32
	for i := 0; i < int(distances.Size()); i++ {
		total += distances.Get(i)
	}
	fmt.Printf("avg. distance: %f\n", total/float32(distances.Size()))
}

func readPoints(path string) [][2]float64 {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var b pointData
	err = json.Unmarshal(dat, &b)
	if err != nil {
		fmt.Printf("error deserializing: %s\n", string(dat))
		return [][2]float64{}
	}
	return b.Points
}

type pointData struct {
	Points [][2]float64 `json:"points"`
}
