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

func (c *Client) SetSnapRadius(n float64) {
	c.snapRadius = n
}

type Client struct {
	client     rk.Client
	channel    chan int
	snapRadius float64
}

func (c Client) Route(from []float64, to []float64) (float64, [][]float64) {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	resp := c.client.Queryrequest(
		int(counter),
		float32(c.snapRadius),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
		true,
	)
	wp := resp.GetWaypoints()
	waypoints := make([][]float64, wp.Size())
	for i := 0; i < len(waypoints); i++ {
		p := wp.Get(i)
		waypoints[i] = []float64{float64(p.GetLon()), float64(p.GetLat())}
	}

	return float64(resp.GetDistance()), waypoints
}

func (c Client) Distance(from []float64, to []float64) float64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	resp := c.client.Queryrequest(
		int(counter),
		float32(c.snapRadius),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
		false,
	)

	return float64(resp.GetDistance())
}

func (c Client) Distances(source []float64, targets [][]float64) []float64 {
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

	matrix := c.client.Distances(counter, s, targetsVector)
	defer routingkit.DeleteUnsignedVector(matrix)
	numRows := matrix.Size()
	rows := make([]float64, numRows)
	for i := 0; i < int(numRows); i++ {
		col := matrix.Get(i)
		rows[i] = float64(col)
	}

	return rows
}
