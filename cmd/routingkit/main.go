//main.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nextmv-io/go-routingkit/routingkit"
)

func main() {

	mapFile := "data/map.osm.pbf"
	chFile := "data/map.ch"

	myClass := routingkit.NewClient()
	defer routingkit.DeleteClient(myClass)

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
