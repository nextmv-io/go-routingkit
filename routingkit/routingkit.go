package routingkit

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
	rk "github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
)

func finalizer(client *rk.Client) {
	routingkit.DeleteClient(*client)
}

func NewClient(mapFile, chFile string) (Client, error) {
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return Client{}, errors.New(fmt.Sprintf("could not find map file at %v", mapFile))
	}

	c := rk.NewClient()
	runtime.SetFinalizer(&c, finalizer)

	concurrentQueries := runtime.GOMAXPROCS(0)

	if _, err := os.Stat(chFile); os.IsNotExist(err) {
		c.Build_ch(concurrentQueries, mapFile, chFile)
	} else {
		c.Load(concurrentQueries, mapFile, chFile)
	}

	channel := make(chan int, concurrentQueries)
	for i := 0; i < concurrentQueries; i++ {
		channel <- i
	}

	return Client{
		client:     c,
		channel:    channel,
		snapRadius: 1000,
	}, nil
}

func (c *Client) SetSnapRadius(n float32) {
	c.snapRadius = n
}

type Client struct {
	client     rk.Client
	channel    chan int
	snapRadius float32
}

func (c Client) Route(from []float32, to []float32) (int64, [][]float32) {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	resp := c.client.Query(
		int(counter),
		float32(c.snapRadius),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
		true,
	)
	wp := resp.GetWaypoints()
	waypoints := make([][]float32, wp.Size())
	for i := 0; i < len(waypoints); i++ {
		p := wp.Get(i)
		waypoints[i] = []float32{float32(p.GetLon()), float32(p.GetLat())}
	}

	return resp.GetDistance(), waypoints
}

func (c Client) Distance(from []float32, to []float32) int64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	resp := c.client.Query(
		int(counter),
		c.snapRadius,
		from[0],
		from[1],
		to[0],
		to[1],
		false,
	)

	return resp.GetDistance()
}

type distanceMatrixRow struct {
	i         int
	distances []int64
}

func (c Client) Nearest(point []float32) []float32 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	p := rk.NewPoint()
	defer routingkit.DeletePoint(p)
	p.SetLon(float32(point[0]))
	p.SetLat(float32(point[1]))
	res := c.client.Nearest(counter, c.snapRadius, p)
	defer rk.DeletePoint(res)
	return []float32{res.GetLon(), res.GetLat()}
}

func (c Client) Matrix(sources [][]float32, targets [][]float32) [][]int64 {
	matrix := make([][]int64, len(sources))

	workers := make(chan struct{}, runtime.GOMAXPROCS(0))
	results := make(chan distanceMatrixRow)

	go func() {
		for i, source := range sources {
			workers <- struct{}{}
			go func(i int, source []float32) {
				distances := c.Distances(source, targets)
				results <- distanceMatrixRow{i, distances}
				<-workers
			}(i, source)
		}
	}()

	for range sources {
		matrixRow := <-results
		matrix[matrixRow.i] = matrixRow.distances
	}

	return matrix
}

func (c Client) Distances(source []float32, targets [][]float32) []int64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()

	s := rk.NewPoint()
	defer routingkit.DeletePoint(s)
	s.SetLon(float32(source[0]))
	s.SetLat(float32(source[1]))

	targetsVector := rk.NewPointVector(int64(len(targets)))
	defer routingkit.DeletePointVector(targetsVector)

	for i := 0; i < len(targets); i++ {
		t := rk.NewPoint()
		t.SetLon(float32(targets[i][0]))
		t.SetLat(float32(targets[i][1]))
		targetsVector.Set(i, t)
	}

	distanceVec := c.client.Distances(counter, float32(c.snapRadius), s, targetsVector)
	defer routingkit.DeleteLongIntVector(distanceVec)
	numDistances := distanceVec.Size()
	distances := make([]int64, numDistances)
	for i := 0; i < int(numDistances); i++ {
		col := distanceVec.Get(i)
		distances[i] = col
	}

	return distances
}
