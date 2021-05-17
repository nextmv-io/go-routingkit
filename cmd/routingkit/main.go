//main.go
package main

import (
	"fmt"
	"sync"
	"time"

	routingkit "github.com/rktest/routingkit"
)

func main() {
	start := time.Now()
	myClass := routingkit.NewClient()
	defer routingkit.DeleteClient(myClass)
	myClass.Load("map.osm.pbf", "map.ch")
	elapsed := time.Since(start)
	fmt.Println("Load pbf and ch took %v", elapsed)
	start = time.Now()
	var number int64 = 100
	v := routingkit.NewIntVector(number)
	defer routingkit.DeleteIntVector(v)
	for i := 0; i < 100; i++ {
		v.Set(i, i)
	}
	average := myClass.Average(v)
	fmt.Println(average)
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()
			_ = myClass.Int_get(39.29971, -76.58749, 39.30587, -76.59735)
		}()
	}
	wg.Wait()
	elapsed = time.Since(start)
	fmt.Println("100 queries took %v", elapsed)
}
