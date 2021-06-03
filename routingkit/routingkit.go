package routingkit

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
	rk "github.com/nextmv-io/go-routingkit/routingkit/internal/routingkit"
)

type Client interface {
	Threaded([]float64, []float64) float64
	Tables([][]float64, [][]float64) []float64
	Average() float64
	Query(float64, []float64, []float64) (float64, [][]float64)
}

func finalizer(client *rk.Client) {
	routingkit.DeleteClient(*client)
}

func Wrapper(mapFile, chFile string) (Client, error) {
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("could not find map file at %v", mapFile))
	}

	c := rk.NewClient()
	runtime.SetFinalizer(&c, finalizer)

	if _, err := os.Stat(chFile); os.IsNotExist(err) {
		c.Build_ch(mapFile, chFile)
	} else {
		c.Load(mapFile, chFile)
	}

	count := 100
	channel := make(chan int, count)
	for i := 0; i < count; i++ {
		channel <- i
	}

	return client{
		client:  c,
		channel: channel,
	}, nil
}

type client struct {
	client  rk.Client
	channel chan int
}

func (c client) Average() float64 {
	ints := make([]int, 10)
	vector := rk.NewIntVector(int64(len(ints)))
	defer rk.DeleteIntVector(vector)
	return c.client.Average(vector)
}

func (c client) Query(radius float64, from []float64, to []float64) (float64, [][]float64) {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	resp := c.client.Queryrequest(
		int(counter),
		float32(radius),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
	)
	wp := resp.GetWaypoints()
	waypoints := make([][]float64, wp.Size())
	for i := 0; i < len(waypoints); i++ {
		p := wp.Get(i)
		waypoints[i] = []float64{float64(p.GetLon()), float64(p.GetLat())}
	}

	return float64(resp.GetDistance()), waypoints
}

func (c client) Threaded(from []float64, to []float64) float64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()
	return float64(c.client.Threaded(
		int(counter),
		float32(from[0]),
		float32(from[1]),
		float32(to[0]),
		float32(to[1]),
	))
}

func (c client) Tables(sources [][]float64, targets [][]float64) []float64 {
	counter := <-c.channel
	defer func() {
		c.channel <- counter
	}()

	sourcesVector := rk.NewPointVector(int64(len(sources)))
	for i := 0; i < len(targets); i++ {
		t := rk.NewPoint()
		t.SetLon(float32(sources[i][0]))
		t.SetLat(float32(sources[i][1]))
		sourcesVector.Set(i, t)
	}

	targetsVector := rk.NewPointVector(int64(len(targets)))
	for i := 0; i < len(targets); i++ {
		t := rk.NewPoint()
		t.SetLon(float32(targets[i][0]))
		t.SetLat(float32(targets[i][1]))
		targetsVector.Set(i, t)
	}

	matrix := c.client.Table(counter, sourcesVector, targetsVector)
	numRows := matrix.Size()
	rows := make([]float64, numRows)
	for i := 0; i < int(numRows); i++ {
		col := matrix.Get(i)
		rows[i] = float64(col)
	}

	// defer func() {
	// 	routingkit.DeletePointVector(sourcesVector)
	// 	routingkit.DeletePointVector(targetsVector)
	// 	routingkit.DeleteMatrix(matrix)
	// }()

	return rows
}
